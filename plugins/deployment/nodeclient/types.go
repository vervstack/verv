package nodeclient

type Platform struct {
	Version string
	Os      string
	Arch    string
}

type RegistryCredentials struct {
	Username    string
	Password    string
	RegistryUrl string
}

type DeployRequest struct {
	ServiceName string
	Environment string
	Image       string
}

type versionResponse struct {
	Version string `json:"version"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
}

type listRegistriesResponse struct {
	Instances []registryInstance `json:"instances"`
}

type registryInstance struct {
	Name string `json:"name"`
}

type credentialsResponse struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	RegistryUrl      string `json:"registryUrl"`
	RegistryUrlSnake string `json:"registry_url"`
}

type emptyBody struct{}

type createDeployBody struct {
	ServiceName  string       `json:"service_name"`
	Environment  string       `json:"environment"`
	Vervonomicon vervonomicon `json:"vervonomicon"`
}

type vervonomicon struct {
	Image string `json:"image"`
}
