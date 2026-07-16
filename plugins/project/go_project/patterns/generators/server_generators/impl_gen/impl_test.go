package impl_gen

import (
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/tests/mocks"
)

func TestGenerateImpl(t *testing.T) {
	projMock := mocks.NewIProjectMock(t)
	projMock.GetNameMock.Expect().Return("test_Proj")

	folders := &folder.Folder{
		Name: "",
		Inner: []*folder.Folder{
			{
				Name: "api",
				Inner: []*folder.Folder{
					{
						Name: patterns.GRPCServer,
						Inner: []*folder.Folder{
							patterns.ProtoContract.Copy(),
						},
					},
				},
			},
		},
	}

	projMock.GetFolderMock.
		Expect().
		Return(folders)

	cfg := &rscliconfig.RsCliConfig{
		Env: rscliconfig.Project{PathToServerDefinition: "api"},
	}

	out, err := GenerateImpl(cfg, projMock)
	require.NoError(t, err)
	_ = out
}
