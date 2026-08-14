package middleware_generators

import (
	_ "embed"
)

var (
	//go:embed templates/cookie_names.go.pattern
	cookieNamesFile []byte

	//go:embed templates/cookie_annotator.go.pattern
	cookieAnnotatorFile []byte

	//go:embed templates/cookie_response.go.pattern
	cookieResponseFile []byte

	//go:embed templates/csrf_interceptor.go.pattern
	csrfInterceptorFile []byte

	//go:embed templates/request_scheme_annotator.go.pattern
	requestSchemeAnnotatorFile []byte
)
