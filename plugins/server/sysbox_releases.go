package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.redsock.ru/rerrors"
)

const (
	sysboxReleasesUrl = "https://api.github.com/repos/nestybox/sysbox/releases?per_page=100"
	githubAcceptType  = "application/vnd.github+json"
	vervUserAgent     = "verv"
)

var sysboxTagPattern = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)

type sysboxAsset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type sysboxReleaseDto struct {
	TagName      string        `json:"tag_name"`
	IsDraft      bool          `json:"draft"`
	IsPrerelease bool          `json:"prerelease"`
	Assets       []sysboxAsset `json:"assets"`
}

type sysboxRelease struct {
	// Version is the tag without the leading "v".
	Version string
	Assets  []sysboxAsset

	semver [3]int
}

// FetchSysboxVersions returns the numbered sysbox release versions without the leading "v", newest first.
func FetchSysboxVersions(ctx context.Context) ([]string, error) {
	releases, err := fetchSysboxReleases(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching sysbox releases")
	}

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		versions = append(versions, release.Version)
	}

	return versions, nil
}

func fetchSysboxReleases(ctx context.Context) ([]sysboxRelease, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sysboxReleasesUrl, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error building request")
	}

	req.Header.Set("Accept", githubAcceptType)
	req.Header.Set("User-Agent", vervUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "error calling github")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, rerrors.Wrap(errUnexpectedStatus, fmt.Sprintf("unexpected status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading response")
	}

	releases, err := parseSysboxReleases(body)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing releases")
	}

	return releases, nil
}

// parseSysboxReleases keeps published releases tagged vX.Y.Z and orders them by numeric semver, newest first;
// the API order is not trusted.
func parseSysboxReleases(body []byte) ([]sysboxRelease, error) {
	var dtos []sysboxReleaseDto

	err := json.Unmarshal(body, &dtos)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding releases json")
	}

	releases := make([]sysboxRelease, 0, len(dtos))

	for _, dto := range dtos {
		if dto.IsDraft || dto.IsPrerelease {
			continue
		}

		semver, ok := parseSysboxTag(dto.TagName)
		if !ok {
			continue
		}

		release := sysboxRelease{
			Version: strings.TrimPrefix(dto.TagName, "v"),
			Assets:  dto.Assets,
			semver:  semver,
		}

		releases = append(releases, release)
	}

	sort.Slice(releases, func(i, j int) bool {
		return isSemverGreater(releases[i].semver, releases[j].semver)
	})

	return releases, nil
}

func parseSysboxTag(tag string) ([3]int, bool) {
	match := sysboxTagPattern.FindStringSubmatch(tag)
	if match == nil {
		return [3]int{}, false
	}

	var semver [3]int

	for idx := range semver {
		part, err := strconv.Atoi(match[idx+1])
		if err != nil {
			return [3]int{}, false
		}

		semver[idx] = part
	}

	return semver, true
}

func isSemverGreater(a, b [3]int) bool {
	for idx := range a {
		if a[idx] != b[idx] {
			return a[idx] > b[idx]
		}
	}

	return false
}

func normalizeSysboxVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

// selectSysboxRelease returns the release for version, or the newest one when version is empty.
func selectSysboxRelease(releases []sysboxRelease, version string) (sysboxRelease, error) {
	if len(releases) == 0 {
		return sysboxRelease{}, rerrors.Wrap(errSysboxNoReleases)
	}

	wanted := normalizeSysboxVersion(version)
	if wanted == "" {
		return releases[0], nil
	}

	for _, release := range releases {
		if release.Version == wanted {
			return release, nil
		}
	}

	return sysboxRelease{}, rerrors.Wrap(errSysboxVersionUnknown, "version: "+wanted)
}

// findSysboxDebUrl matches by suffix because asset names are inconsistent across releases
// (sysbox-ce_0.7.1.linux_amd64.deb vs sysbox-ce_0.6.6-0.linux_amd64.deb).
func findSysboxDebUrl(release sysboxRelease, arch string) (string, error) {
	suffix := ".linux_" + arch + ".deb"

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			return asset.Url, nil
		}
	}

	return "", rerrors.Wrap(errSysboxAssetMissing, "version: "+release.Version+", arch: "+arch)
}
