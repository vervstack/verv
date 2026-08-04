package config_generators

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

// stubSkeletonYAML is a defensive fallback used only when the caller has no
// resolved config.yaml content to embed. In the normal verv init/tidy flow
// this should not happen: PrepareConfigFolder.Do (plugins/project/actions/
// go_actions/config.go) always runs generateConfigYamlFile before calling
// GenerateConfigFolder, which guarantees a config.yaml entry — with real
// content on `tidy`, or an empty-but-present AppConfig marshaled to yaml on
// a from-scratch `init` — exists in the in-memory folder tree by the time
// the skeleton is generated.
var stubSkeletonYAML = []byte("app_info: {}\n")

type skeletonGenArgs struct {
	GenTag string
}

// newGenerateConfigSkeleton returns the generated skeleton.go wrapper (which
// //go:embeds the yaml) and the skeleton.yaml file itself. skeleton.yaml's
// content is configYamlBytes — the project's already-resolved, up-to-date
// config.yaml content — copied byte-for-byte, or stubSkeletonYAML if
// configYamlBytes is empty.
func newGenerateConfigSkeleton(configYamlBytes []byte) func() (goFile, yamlFile *folder.Folder, err error) {
	return func() (goFile, yamlFile *folder.Folder, err error) {
		content := configYamlBytes
		if len(content) == 0 {
			content = stubSkeletonYAML
		}

		args := skeletonGenArgs{
			GenTag: defaultGenTag,
		}

		buf := &rw.RW{}

		err = configSkeletonTemplate.Execute(buf, args)
		if err != nil {
			return nil, nil, rerrors.Wrap(err, "error executing config skeleton template")
		}

		goFile = &folder.Folder{
			Name:    patterns.ConfigSkeletonGoFileName,
			Content: buf.Bytes(),
		}

		yamlFile = &folder.Folder{
			Name:    patterns.ConfigSkeletonYamlFileName,
			Content: content,
		}

		return goFile, yamlFile, nil
	}
}
