package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.redsock.ru/rerrors"
)

const fetchTimeout = 30 * time.Second

func fetchText(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error calling server")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", rerrors.Wrap(errUnexpectedStatus, fmt.Sprintf("unexpected status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading response")
	}

	return string(body), nil
}
