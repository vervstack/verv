package imagebuild

import "strings"

func ImageRef(registryUrl, serviceName, tag string) string {
	return registryHost(registryUrl) + "/" + serviceName + ":" + tag
}

func registryHost(registryUrl string) string {
	host := strings.TrimPrefix(registryUrl, "https://")
	host = strings.TrimPrefix(host, "http://")

	return strings.TrimRight(host, "/")
}
