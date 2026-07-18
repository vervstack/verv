package clients_generators

import (
	"text/template"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/rw"
)

func GenerateGRPCConn() ([]byte, error) {
	return execute(grpcConnTemplate)
}

func GenerateRedisConn() ([]byte, error) {
	return execute(redisConnTemplate)
}

func GenerateSQLConn() ([]byte, error) {
	return execute(sqlConnTemplate)
}

func GeneratePostgresDriver() ([]byte, error) {
	return execute(postgresDriverTemplate)
}

func GenerateSqliteDriver() ([]byte, error) {
	return execute(sqliteDriverTemplate)
}

func GenerateTelegramConn() ([]byte, error) {
	return execute(telegramConnTemplate)
}

func execute(t *template.Template) ([]byte, error) {
	out := &rw.RW{}

	err := t.Execute(out, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating "+t.Name())
	}

	return out.Bytes(), nil
}
