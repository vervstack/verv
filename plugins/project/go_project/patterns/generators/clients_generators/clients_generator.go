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

func GeneratePostgresConn() ([]byte, error) {
	return execute(postgresConnTemplate)
}

// PostgresInstance describes a single configured Postgres data source used to
// render the per-instance Connect/Migrate function pair in instances.go.
type PostgresInstance struct {
	PascalName string
}

// PostgresInstancesArgs is the template data for instances.go, listing every
// currently configured Postgres data source.
type PostgresInstancesArgs struct {
	Instances []PostgresInstance
}

func GeneratePostgresInstances(args PostgresInstancesArgs) ([]byte, error) {
	return executeWith(postgresInstancesTemplate, args)
}

func execute(t *template.Template) ([]byte, error) {
	return executeWith(t, nil)
}

func executeWith(t *template.Template, data any) ([]byte, error) {
	out := &rw.RW{}

	err := t.Execute(out, data)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating "+t.Name())
	}

	return out.Bytes(), nil
}
