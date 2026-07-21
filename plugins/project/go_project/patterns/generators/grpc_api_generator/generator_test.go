package grpc_api_generator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/tests/project_mock"
)

func Test_GenerateServiceApiProto(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		shortName string
		expected  string
	}{
		"basic": {
			shortName: "myservice",
			expected: `syntax = "proto3";

package myservice_api;

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";
import "npm.proto";

option go_package = "/myservice_api";
option (npm_package) = "@myservice/api";

service myserviceAPI {
  rpc Version(Version.Request) returns (Version.Response) {
    option (google.api.http) = {
      get: "/api/version"
    };
  };
}

message Version {
  message Request {}

  message Response {
    string version = 1;
    google.protobuf.Timestamp client_timestamp = 2;
  }
}
`,
		},
		"different_name": {
			shortName: "billing",
			expected: `syntax = "proto3";

package billing_api;

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";
import "npm.proto";

option go_package = "/billing_api";
option (npm_package) = "@billing/api";

service billingAPI {
  rpc Version(Version.Request) returns (Version.Response) {
    option (google.api.http) = {
      get: "/api/version"
    };
  };
}

message Version {
  message Request {}

  message Response {
    string version = 1;
    google.protobuf.Timestamp client_timestamp = 2;
  }
}
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			proj := project_mock.GetMockProject(t)

			proj.Name = tc.shortName

			f, err := GenerateServiceApiProto(proj)
			require.NoError(t, err)

			require.Equal(t, tc.shortName+".proto", f.Name)
			require.Equal(t, tc.expected, string(f.Content))
		})
	}
}
