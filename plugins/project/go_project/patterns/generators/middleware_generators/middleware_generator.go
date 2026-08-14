package middleware_generators

func GenerateCookieNames() []byte {
	return copyBytes(cookieNamesFile)
}

func GenerateCookieAnnotator() []byte {
	return copyBytes(cookieAnnotatorFile)
}

func GenerateCookieResponse() []byte {
	return copyBytes(cookieResponseFile)
}

func GenerateCSRFInterceptor() []byte {
	return copyBytes(csrfInterceptorFile)
}

func GenerateRequestSchemeAnnotator() []byte {
	return copyBytes(requestSchemeAnnotatorFile)
}

func copyBytes(src []byte) []byte {
	res := make([]byte, len(src))
	copy(res, src)

	return res
}
