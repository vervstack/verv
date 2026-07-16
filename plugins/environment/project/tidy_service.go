package project

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/compose"
	"go.vervstack.ru/verv/internal/envpatterns"
	"go.vervstack.ru/verv/internal/utils/copier"
)

var ErrNoProjectComposePattern = rerrors.New("")

func (e *ProjEnv) tidyService() error {
	srcPattern, ok := e.globalComposePatternManager.Patterns[envpatterns.ProjNamePattern]
	if !ok {
		return ErrNoProjectComposePattern
	}

	var pattern compose.Pattern
	err := copier.Copy(srcPattern, &pattern)
	if err != nil {
		return rerrors.Wrap(err, "error coping proj pattern")
	}

	// TODO
	//for _, s := range e.Config.Servers {
	//port := s.GetPort()

	//pattern.ContainerDefinition.Ports = append(
	//	pattern.ContainerDefinition.Ports,
	//	fmt.Sprintf("%d:%d", e.globalPortManager.GetNextPort(port, e.projName+"_"+s.GetName()), port),
	//)
	//}

	e.Compose.AppendService(e.projName, pattern.ContainerDefinition)

	return nil
}
