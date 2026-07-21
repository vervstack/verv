package renamer

import (
	"bytes"

	"go.vervstack.ru/verv/internal/envpatterns"
)

func ReplaceProjectNameShort(src []byte, newName string) []byte {
	b := make([]byte, len(src))
	copy(b, src)

	bigName := bytes.ReplaceAll(bytes.ToUpper([]byte(newName)), []byte{'-'}, []byte{'_'})

	b = bytes.ReplaceAll(
		b,
		[]byte(envpatterns.ProjNameCapsPattern),
		bigName,
	)

	smallName := bytes.ReplaceAll(bytes.ToLower([]byte(newName)), []byte{'-'}, []byte{'_'})

	b = bytes.ReplaceAll(b,
		[]byte(envpatterns.ProjNamePattern),
		smallName,
	)

	return b
}
