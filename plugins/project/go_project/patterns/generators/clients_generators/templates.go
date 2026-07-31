package clients_generators

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/grpc/conn.go.pattern
	grpcConnPattern  string
	grpcConnTemplate *template.Template

	//go:embed templates/redis/conn.go.pattern
	redisConnPattern  string
	redisConnTemplate *template.Template

	//go:embed templates/sqldb/conn.go.pattern
	sqlConnPattern  string
	sqlConnTemplate *template.Template

	//go:embed templates/sqldb/sqlite.go.pattern
	sqliteDriverPattern  string
	sqliteDriverTemplate *template.Template

	//go:embed templates/telegram/conn.go.pattern
	telegramConnPattern  string
	telegramConnTemplate *template.Template

	//go:embed templates/postgres/conn.go.pattern
	postgresConnPattern  string
	postgresConnTemplate *template.Template

	//go:embed templates/postgres/driver.go.pattern
	postgresDriverPattern  string
	postgresDriverTemplate *template.Template

	//go:embed templates/postgres/instances.go.pattern
	postgresInstancesPattern  string
	postgresInstancesTemplate *template.Template
)

//nolint:gochecknoinits // one-time compile of embedded templates into package-level *template.Template values
func init() {
	grpcConnTemplate = template.Must(
		template.New("grpc_conn").
			Parse(grpcConnPattern))

	redisConnTemplate = template.Must(
		template.New("redis_conn").
			Parse(redisConnPattern))

	sqlConnTemplate = template.Must(
		template.New("sql_conn").
			Parse(sqlConnPattern))

	sqliteDriverTemplate = template.Must(
		template.New("sqlite_driver").
			Parse(sqliteDriverPattern))

	telegramConnTemplate = template.Must(
		template.New("telegram_conn").
			Parse(telegramConnPattern))

	postgresConnTemplate = template.Must(
		template.New("postgres_conn").
			Parse(postgresConnPattern))

	postgresDriverTemplate = template.Must(
		template.New("postgres_driver").
			Parse(postgresDriverPattern))

	postgresInstancesTemplate = template.Must(
		template.New("postgres_instances").
			Parse(postgresInstancesPattern))
}
