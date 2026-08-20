//go:build !windows

package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamOutputOSSpecificReapsPagerOnGenerateError(t *testing.T) {
	tempDir := t.TempDir()
	marker := filepath.Join(tempDir, "pager-finished")
	pager := filepath.Join(tempDir, "pager.sh")

	script := "#!/bin/sh\nsleep 1\nprintf done > \"$PAGER_MARKER\"\n"
	require.NoError(t, os.WriteFile(pager, []byte(script), 0700))
	t.Setenv("PAGER", pager)
	t.Setenv("PAGER_MARKER", marker)

	generateErr := errors.New("generation failed")
	err := streamOutputOSSpecific("test", func(*os.File) error {
		return generateErr
	})
	require.ErrorIs(t, err, generateErr)

	contents, err := os.ReadFile(marker)
	require.NoError(t, err, "pager must finish before streamOutputOSSpecific returns")
	require.Equal(t, "done", string(contents))
}
