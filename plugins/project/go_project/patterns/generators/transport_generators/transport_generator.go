package transport_generators

import (
	"text/template"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/rw"
)

func GenerateServerManager() ([]byte, error) {
	return execute(serverManagerTemplate)
}

func GenerateGrpcServer() []byte {
	res := make([]byte, len(grpcServerFile))
	copy(res, grpcServerFile)

	return res
}

func GenerateHttpServer() ([]byte, error) {
	return []byte(httpServerPattern), nil
}

func GenerateTelegramListener() ([]byte, error) {
	return execute(telegramListenerTemplate)
}

func GenerateTelegramVersionHandler() ([]byte, error) {
	return execute(telegramVersionHandlerTemplate)
}

type gatewayMuxArgs struct {
	FullProjPath string
}

func GenerateGatewayMux(fullProjPath string) ([]byte, error) {
	args := gatewayMuxArgs{FullProjPath: fullProjPath}

	return executeWithData(gatewayMuxTemplate, args)
}

func execute(t *template.Template) ([]byte, error) {
	return executeWithData(t, nil)
}

func executeWithData(t *template.Template, data any) ([]byte, error) {
	out := &rw.RW{}

	err := t.Execute(out, data)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating "+t.Name())
	}

	return out.Bytes(), nil
}
