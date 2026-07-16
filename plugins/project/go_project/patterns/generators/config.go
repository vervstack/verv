package generators

import (
	"strings"

	"go.vervstack.ru/verv/internal/utils/cases"
)

type KeyValue struct {
	Key   string
	Value string
}

var nameReplacer = strings.NewReplacer(
	" ", "_",
	"-", "_")

func NormalizeResourceName(in string) string {
	in = nameReplacer.Replace(in)
	return cases.SnakeToPascal(in)
}
