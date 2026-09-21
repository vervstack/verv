package cases_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/utils/cases"
)

func Test_ToPascal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"single word", "name", "Name"},
		{"kebab case", "my-service-name", "MyServiceName"},
		{"snake case", "my_service_name", "MyServiceName"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, cases.ToPascal(tc.input))
		})
	}
}
