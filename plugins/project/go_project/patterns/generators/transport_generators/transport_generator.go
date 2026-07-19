package transport_generators

import (
	"text/template"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/rw"
)

func GenerateServerManager() ([]byte, error) {
	return execute(serverManagerTemplate)
}

func GenerateGrpcServer() ([]byte, error) {
	return execute(grpcServerTemplate)
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

func execute(t *template.Template) ([]byte, error) {
	out := &rw.RW{}

	err := t.Execute(out, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating "+t.Name())
	}

	return out.Bytes(), nil
}
