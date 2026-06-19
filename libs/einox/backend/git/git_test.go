package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/memfs"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/storage/memory"
	"github.com/stretchr/testify/require"
)

type testRepo struct {
	repo         *git.Repository
	mainCommit   plumbing.Hash
	branchCommit plumbing.Hash
	tagCommit    plumbing.Hash
}

func newTestRepo(t *testing.T) *testRepo {
	t.Helper()

	storer := memory.NewStorage()
	fs := memfs.New()
	repo, err := git.Init(storer, git.WithWorkTree(fs))
	require.NoError(t, err)

	writeFile(t, fs, "README.md", "root readme\n")
	writeFile(t, fs, "users/tom/projects/todos/README.md", "hello\nworld\nline3\n")
	writeFile(t, fs, "users/tom/projects/todos/main.go", "package main\nfunc main() {}\n")
	writeFile(t, fs, "users/alice/note.txt", "alice secret\n")

	w, err := repo.Worktree()
	require.NoError(t, err)
	_, err = w.Add(".")
	require.NoError(t, err)

	mainCommit, err := w.Commit("main commit", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", mainCommit, nil)
	require.NoError(t, err)

	writeFile(t, fs, "branch.txt", "branch content\n")
	_, err = w.Add("branch.txt")
	require.NoError(t, err)
	branchCommit, err := w.Commit("feature commit", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)

	require.NoError(t, repo.Storer.SetReference(
		plumbing.NewHashReference(plumbing.NewBranchReferenceName("feature"), branchCommit),
	))

	ref, err := repo.Reference(plumbing.NewBranchReferenceName("feature"), true)
	require.NoError(t, err)
	require.Equal(t, branchCommit, ref.Hash())

	return &testRepo{
		repo:         repo,
		mainCommit:   mainCommit,
		branchCommit: branchCommit,
		tagCommit:    mainCommit,
	}
}

func writeFile(t *testing.T, fs billy.Filesystem, name, content string) {
	t.Helper()
	require.NoError(t, fs.MkdirAll(filepath.Dir(name), 0o755))
	f, err := fs.Create(name)
	require.NoError(t, err)
	_, err = f.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

func newTestBackend(t *testing.T, repo *git.Repository, cfg Config) *GitBackend {
	t.Helper()
	if cfg.Transport == "" {
		cfg.Transport = TransportLocal
	}
	if cfg.URL == "" {
		cfg.URL = "file://local"
	}
	backend, err := NewGitBackend(cfg)
	require.NoError(t, err)
	backend.repo = repo
	require.NoError(t, backend.resolveCommitAndTree())
	return backend
}

func TestRead(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{
		BaseDir: "users/tom",
	})

	content, err := backend.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/projects/todos/README.md",
	})
	require.NoError(t, err)
	require.Equal(t, "hello\nworld\nline3\n", content.Content)

	partial, err := backend.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/projects/todos/README.md",
		Offset:   2,
		Limit:    1,
	})
	require.NoError(t, err)
	require.Equal(t, "world", partial.Content)

	_, err = backend.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/missing.txt",
	})
	require.Error(t, err)
}

func TestLsInfo(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{BaseDir: "users/tom"})

	rootItems, err := backend.LsInfo(t.Context(), &filesystem.LsInfoRequest{Path: "/"})
	require.NoError(t, err)
	names := make([]string, 0, len(rootItems))
	for _, item := range rootItems {
		names = append(names, item.Path)
	}
	require.Contains(t, names, "projects")

	subItems, err := backend.LsInfo(t.Context(), &filesystem.LsInfoRequest{Path: "/projects/todos"})
	require.NoError(t, err)
	require.Len(t, subItems, 2)
}

func TestGrepRaw(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{BaseDir: "users/tom"})

	matches, err := backend.GrepRaw(t.Context(), &filesystem.GrepRequest{
		Pattern: "hello",
	})
	require.NoError(t, err)
	require.NotEmpty(t, matches)

	filtered, err := backend.GrepRaw(t.Context(), &filesystem.GrepRequest{
		Pattern: "package",
		Glob:    "*.go",
	})
	require.NoError(t, err)
	require.Len(t, filtered, 1)
}

func TestGlobInfo(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{BaseDir: "users/tom"})

	files, err := backend.GlobInfo(t.Context(), &filesystem.GlobInfoRequest{
		Pattern: "**/*.go",
		Path:    "/projects",
	})
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Equal(t, "todos/main.go", files[0].Path)
}

func TestBaseDirAndPathEscape(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{BaseDir: "users/tom"})

	// baseDir 隔离：虚拟路径映射到 baseDir 子树，无法读取仓库根下 users/alice
	_, err := backend.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/users/alice/note.txt",
	})
	require.Error(t, err)
}

func TestRefResolution(t *testing.T) {
	tr := newTestRepo(t)

	byCommit := newTestBackend(t, tr.repo, Config{
		Ref:     RefConfig{Commit: tr.mainCommit.String()},
		BaseDir: "users/tom",
	})
	content, err := byCommit.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/projects/todos/README.md",
	})
	require.NoError(t, err)
	require.Contains(t, content.Content, "hello")

	byTag := newTestBackend(t, tr.repo, Config{
		Ref:     RefConfig{Tag: "v1.0.0"},
		BaseDir: "users/tom",
	})
	_, err = byTag.Read(t.Context(), &filesystem.ReadRequest{FilePath: "/projects/todos/README.md"})
	require.NoError(t, err)

	byBranch := newTestBackend(t, tr.repo, Config{
		Ref:     RefConfig{Branch: "feature"},
		BaseDir: "",
	})
	readBranch, err := byBranch.Read(t.Context(), &filesystem.ReadRequest{FilePath: "/branch.txt"})
	require.NoError(t, err)
	require.Equal(t, "branch content\n", readBranch.Content)
}

func TestReadOnly(t *testing.T) {
	tr := newTestRepo(t)
	backend := newTestBackend(t, tr.repo, Config{BaseDir: "users/tom"})

	require.ErrorIs(t, backend.Write(t.Context(), &filesystem.WriteRequest{
		FilePath: "/projects/todos/new.txt",
		Content:  "nope",
	}), errReadOnly)

	require.ErrorIs(t, backend.Edit(t.Context(), &filesystem.EditRequest{
		FilePath:  "/projects/todos/README.md",
		OldString: "hello",
		NewString: "bye",
	}), errReadOnly)
}

func TestConfigValidation(t *testing.T) {
	_, err := NewGitBackend(Config{})
	require.Error(t, err)

	_, err = NewGitBackend(Config{
		URL:       "git@github.com:org/repo.git",
		Transport: TransportSSH,
	})
	require.Error(t, err)

	_, err = NewGitBackend(Config{
		URL: "https://github.com/org/repo.git",
		Ref: RefConfig{Branch: "main", Tag: "v1"},
	})
	require.Error(t, err)
}

func TestAuthDoesNotLeakCredentials(t *testing.T) {
	auth := buildHTTPAuth(AuthConfig{
		Username: "user",
		Password: "super-secret-password",
		Token:    "ghp_super_secret_token",
	})
	require.NotNil(t, auth)

	basic, ok := auth.(*githttp.BasicAuth)
	require.True(t, ok)
	require.Equal(t, "user", basic.Username)
	require.Equal(t, "ghp_super_secret_token", basic.Password)
}

func TestLocalBareRepoOpen(t *testing.T) {
	tr := newTestRepo(t)
	dir := t.TempDir()
	bareDir := filepath.Join(dir, "repo.git")
	_, err := git.PlainInit(bareDir, true)
	require.NoError(t, err)

	remote, err := tr.repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{bareDir},
	})
	require.NoError(t, err)
	require.NoError(t, remote.Push(&git.PushOptions{RemoteName: "origin"}))

	backend, err := NewGitBackend(Config{
		URL:       "file://" + bareDir,
		Transport: TransportLocal,
		BaseDir:   "users/tom",
	})
	require.NoError(t, err)

	content, err := backend.Read(t.Context(), &filesystem.ReadRequest{
		FilePath: "/projects/todos/README.md",
	})
	require.NoError(t, err)
	require.Contains(t, content.Content, "hello")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
