package tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CompareLongStrings(t *testing.T, expected, actual []byte) (eq bool) {
	t.Helper()

	expectedReader := bytes.NewReader(expected)
	actualReader := bytes.NewReader(actual)

	// 800 because of IDE parsing abilities
	const defaultBatchSizeBytes = 800

	for {
		expectedSlice := make([]byte, defaultBatchSizeBytes)
		actualSlice := make([]byte, defaultBatchSizeBytes)

		expLen, expErr := expectedReader.Read(expectedSlice)
		actLen, actErr := actualReader.Read(actualSlice)

		if expErr == io.EOF && actErr == io.EOF {
			break
		}

		expectedSlice = expectedSlice[:expLen]
		actualSlice = actualSlice[:actLen]

		eq = assert.Equal(t, string(expectedSlice), string(actualSlice))
		if !eq {
			return false
		}

		require.NoError(t, actErr)
		require.NoError(t, expErr)
	}

	return true
}
