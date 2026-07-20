package version

import (
	"encoding/json"
	"net/http"

	v "github.com/hashicorp/go-version"
	"github.com/rs/zerolog/log"
)

func CanUpdate() (string, bool) {
	resp, err := http.Get("https://api.github.com/repos/vervstack/verv/releases/latest")
	if err != nil {
		return "", false
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body after getting latest release")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return "", false
	}

	tag := m["tag_name"]
	if tag == "" {
		return "", false
	}
	tagStr, _ := tag.(string)
	if tagStr == "" {
		return "", false
	}

	originVersion, err := v.NewVersion(tagStr)
	if err != nil {
		return "", false
	}

	localVersion, err := v.NewVersion(version)
	if err != nil {
		return "", false
	}

	return tagStr, originVersion.GreaterThan(localVersion)
}
