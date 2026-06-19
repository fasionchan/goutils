package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVirtualToRepo(t *testing.T) {
	repoPath, err := virtualToRepo("users/tom", "/projects/todos")
	require.NoError(t, err)
	require.Equal(t, "users/tom/projects/todos", repoPath)

	virtual := repoToVirtual("users/tom", "users/tom/projects/todos")
	require.Equal(t, "/projects/todos", virtual)
}

func TestNormalizePath(t *testing.T) {
	require.Equal(t, "/", normalizePath(""))
	require.Equal(t, "/foo/bar", normalizePath("foo/bar"))
}
