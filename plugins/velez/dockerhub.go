package velez

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.redsock.ru/rerrors"
)

const (
	dockerHubTagsURL = "https://hub.docker.com/v2/repositories/vervstack/velez/tags?page_size=3&ordering=last_updated"
	latestTag        = "latest"
	maxTags          = 3
)

var (
	errNoTagsFound            = rerrors.New("no tags found for vervstack/velez on docker hub")
	errDockerHubRequestFailed = rerrors.New("docker hub request failed")
)

type tagsResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

// FetchLatestTags returns up to the 3 most recently updated tags of
// vervstack/velez from Docker Hub's public API, skipping a literal "latest"
// entry if present. It is a pure data-fetch function — no interactive UI.
func FetchLatestTags(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dockerHubTagsURL, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error building docker hub request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "error calling docker hub")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading docker hub response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, rerrors.Wrap(errDockerHubRequestFailed, fmt.Sprintf("status %d", resp.StatusCode))
	}

	tags, err := parseTags(body)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing docker hub response")
	}

	return tags, nil
}

// parseTags extracts up to maxTags tag names from a Docker Hub tags-list
// response body, in the order Docker Hub returns them, skipping a literal
// "latest" entry if present.
func parseTags(body []byte) ([]string, error) {
	var parsed tagsResponse

	err := json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding docker hub tags response")
	}

	tags := make([]string, 0, maxTags)

	for _, r := range parsed.Results {
		if r.Name == latestTag {
			continue
		}

		tags = append(tags, r.Name)

		if len(tags) == maxTags {
			break
		}
	}

	if len(tags) == 0 {
		return nil, rerrors.Wrap(errNoTagsFound)
	}

	return tags, nil
}
