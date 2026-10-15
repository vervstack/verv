package nodeclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.redsock.ru/rerrors"
)

const authorizationHeader = "Grpc-Metadata-Authorization"

type Client struct {
	baseUrl    string
	secret     string
	httpClient *http.Client
}

func New(url, secret string) *Client {
	return &Client{
		baseUrl:    strings.TrimRight(url, "/"),
		secret:     secret,
		httpClient: &http.Client{},
	}
}

func (c *Client) Version(ctx context.Context) (Platform, error) {
	var response versionResponse
	err := c.do(ctx, http.MethodGet, "/api/version", nil, &response)
	if err != nil {
		return Platform{}, rerrors.Wrap(err, "error requesting version")
	}

	platform := Platform{
		Version: response.Version,
		Os:      response.Os,
		Arch:    response.Arch,
	}

	return platform, nil
}

func (c *Client) ListRegistryNames(ctx context.Context) ([]string, error) {
	var response listRegistriesResponse
	body := emptyBody{}
	err := c.do(ctx, http.MethodPost, "/api/container-registry/list", body, &response)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing registries")
	}

	names := make([]string, 0, len(response.Instances))
	for _, instance := range response.Instances {
		names = append(names, instance.Name)
	}

	return names, nil
}

func (c *Client) RegistryCredentials(ctx context.Context, name string) (RegistryCredentials, error) {
	path := "/api/container-registry/" + url.PathEscape(name) + "/credentials"

	var response credentialsResponse
	err := c.do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return RegistryCredentials{}, rerrors.Wrap(err, "error requesting registry credentials")
	}

	registryUrl := response.RegistryUrl
	if registryUrl == "" {
		registryUrl = response.RegistryUrlSnake
	}

	credentials := RegistryCredentials{
		Username:    response.Username,
		Password:    response.Password,
		RegistryUrl: registryUrl,
	}

	return credentials, nil
}

func (c *Client) CreateDeploy(ctx context.Context, req DeployRequest) error {
	body := createDeployBody{
		ServiceName:  req.ServiceName,
		Environment:  req.Environment,
		Vervonomicon: vervonomicon{Image: req.Image},
	}

	err := c.do(ctx, http.MethodPost, "/api/service/deploy/create", body, nil)
	if err != nil {
		return rerrors.Wrap(err, "error creating deploy")
	}

	return nil
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return rerrors.Wrap(err, "error encoding request body")
		}
		reader = bytes.NewReader(raw)
	}

	request, err := http.NewRequestWithContext(ctx, method, c.baseUrl+path, reader)
	if err != nil {
		return rerrors.Wrap(err, "error building request")
	}
	request.Header.Set(authorizationHeader, c.secret)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return rerrors.Wrap(err, "error sending request")
	}
	defer closeQuietly(response.Body)

	err = checkStatus(response.StatusCode)
	if err != nil {
		return rerrors.Wrap(err)
	}

	if out == nil {
		return nil
	}

	err = json.NewDecoder(response.Body).Decode(out)
	if err != nil {
		return rerrors.Wrap(err, "error decoding response body")
	}

	return nil
}

func checkStatus(code int) error {
	if code >= http.StatusOK && code < http.StatusMultipleChoices {
		return nil
	}

	details := fmt.Sprintf("status: %d", code)
	if code == http.StatusUnauthorized || code == http.StatusForbidden {
		return rerrors.Wrap(ErrUnauthorized, details)
	}

	return rerrors.Wrap(ErrUnexpectedStatus, details)
}

func closeQuietly(closer io.Closer) {
	_ = closer.Close()
}
