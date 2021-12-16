package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitRepoRootDir(t *testing.T) {
	repoRoot := GetGitRepoRootDir()
	info, err := os.Stat(repoRoot + "/.git")
	assert.NoError(t, err)
	assert.Equal(t, true, info.IsDir())
}
