package git

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/go-git/go-billy/v6/memfs"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/go-git/go-git/v6/storage/memory"
)

// GitBackend 基于 Git 仓库 tree 的只读 filesystem.Backend 实现。
type GitBackend struct {
	cfg Config

	once    sync.Once
	initErr error

	repo       *git.Repository
	commit     *object.Commit
	baseTree   *object.Tree
	modifiedAt string
}

var _ filesystem.Backend = (*GitBackend)(nil)

// NewGitBackend 创建 GitBackend；配置校验在构造阶段完成，仓库 clone 在首次只读操作时懒加载。
func NewGitBackend(cfg Config) (*GitBackend, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.Transport == "" {
		cfg.Transport = inferTransport(cfg.URL)
	}
	return &GitBackend{cfg: cfg}, nil
}

func (b *GitBackend) ensureReady(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b.once.Do(func() {
		if b.repo != nil {
			b.initErr = b.resolveCommitAndTree()
			return
		}
		b.initErr = b.initRepo()
	})
	return b.initErr
}

func (b *GitBackend) initRepo() error {
	transportType := b.cfg.Transport
	if transportType == "" {
		transportType = inferTransport(b.cfg.URL)
	}

	if transportType == TransportLocal {
		path := b.cfg.URL
		path = strings.TrimPrefix(path, "file://")

		repo, err := git.PlainOpen(path)
		if err != nil {
			return fmt.Errorf("git backend: open local repo: %w", err)
		}
		b.repo = repo
		return b.resolveCommitAndTree()
	}

	storer := memory.NewStorage()
	fs := memfs.New()
	opts, err := cloneOptions(b.cfg)
	if err != nil {
		return err
	}
	repo, err := git.Clone(storer, fs, opts)
	if err != nil {
		return fmt.Errorf("git backend: clone failed: %w", err)
	}
	b.repo = repo
	return b.resolveCommitAndTree()
}

func (b *GitBackend) resolveCommitAndTree() error {
	commit, err := resolveCommit(b.repo, b.cfg.Ref)
	if err != nil {
		return err
	}
	b.commit = commit
	b.modifiedAt = commit.Committer.When.UTC().Format(time.RFC3339Nano)

	rootTree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf("git backend: resolve tree: %w", err)
	}

	baseDir := cleanRepoPath(b.cfg.BaseDir)
	if baseDir == "" {
		b.baseTree = rootTree
		return nil
	}

	entry, err := findTreeEntry(b.repo.Storer, rootTree, baseDir)
	if err != nil {
		return fmt.Errorf("git backend: baseDir %q not found: %w", baseDir, err)
	}
	if entry.Mode != filemode.Dir {
		return fmt.Errorf("git backend: baseDir %q is not a directory", baseDir)
	}
	tree, err := object.GetTree(b.repo.Storer, entry.Hash)
	if err != nil {
		return fmt.Errorf("git backend: load baseDir tree: %w", err)
	}
	b.baseTree = tree
	return nil
}

func resolveCommit(repo *git.Repository, ref RefConfig) (*object.Commit, error) {
	if ref.Commit != "" {
		hash := plumbing.NewHash(ref.Commit)
		commit, err := repo.CommitObject(hash)
		if err != nil {
			return nil, fmt.Errorf("git backend: commit %q not found: %w", ref.Commit, err)
		}
		return commit, nil
	}
	if ref.Tag != "" {
		tagRef, err := repo.Tag(ref.Tag)
		if err != nil {
			return nil, fmt.Errorf("git backend: tag %q not found: %w", ref.Tag, err)
		}
		obj, err := repo.TagObject(tagRef.Hash())
		if err == nil {
			return repo.CommitObject(obj.Target)
		}
		return repo.CommitObject(tagRef.Hash())
	}
	if ref.Branch != "" {
		name := plumbing.NewBranchReferenceName(ref.Branch)
		reference, err := repo.Reference(name, true)
		if err != nil {
			return nil, fmt.Errorf("git backend: branch %q not found: %w", ref.Branch, err)
		}
		return repo.CommitObject(reference.Hash())
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("git backend: resolve HEAD: %w", err)
	}
	return repo.CommitObject(head.Hash())
}

func findTreeEntry(store storer.EncodedObjectStorer, root *object.Tree, repoPath string) (*object.TreeEntry, error) {
	repoPath = cleanRepoPath(repoPath)
	if repoPath == "" {
		return nil, errNotFound
	}

	parts := strings.Split(repoPath, "/")
	current := root

	for i, part := range parts {
		var found *object.TreeEntry
		for idx := range current.Entries {
			entry := current.Entries[idx]
			if entry.Name == part {
				found = &entry
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("%w: %s", errNotFound, repoPath)
		}
		if i == len(parts)-1 {
			return found, nil
		}
		if found.Mode != filemode.Dir {
			return nil, fmt.Errorf("%w: %s is not a directory", errNotFound, repoPath)
		}
		next, err := object.GetTree(store, found.Hash)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return nil, errNotFound
}

type resolvedNode struct {
	entry    *object.TreeEntry
	tree     *object.Tree
	repoPath string
	virtual  string
}

func (b *GitBackend) resolveNode(virtualPath string) (*resolvedNode, error) {
	repoPath, err := virtualToRepo(b.cfg.BaseDir, virtualPath)
	if err != nil {
		return nil, err
	}

	if repoPath == "" || normalizePath(virtualPath) == "/" {
		return &resolvedNode{
			tree:     b.baseTree,
			repoPath: cleanRepoPath(b.cfg.BaseDir),
			virtual:  "/",
		}, nil
	}

	base := cleanRepoPath(b.cfg.BaseDir)
	var searchTree *object.Tree
	var searchPath string

	if base == "" {
		searchTree = b.baseTree
		searchPath = repoPath
	} else {
		searchTree = b.baseTree
		searchPath = strings.TrimPrefix(repoPath, base+"/")
		if searchPath == repoPath && repoPath != base {
			return nil, errPathEscape
		}
	}

	entry, err := findTreeEntry(b.repo.Storer, searchTree, searchPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %s", normalizePath(virtualPath))
	}

	node := &resolvedNode{
		entry:    entry,
		repoPath: repoPath,
		virtual:  normalizePath(virtualPath),
	}
	if entry.Mode == filemode.Dir {
		tree, err := object.GetTree(b.repo.Storer, entry.Hash)
		if err != nil {
			return nil, err
		}
		node.tree = tree
	}
	return node, nil
}

func (b *GitBackend) readBlob(entry *object.TreeEntry) (string, int64, error) {
	blob, err := b.repo.BlobObject(entry.Hash)
	if err != nil {
		return "", 0, err
	}
	reader, err := blob.Reader()
	if err != nil {
		return "", 0, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", 0, err
	}
	return string(data), int64(len(data)), nil
}

type fileRecord struct {
	virtual  string
	repoPath string
	content  string
	size     int64
}

func (b *GitBackend) walkFiles(rootTree *object.Tree, repoPrefix, virtualPrefix string, fn func(fileRecord) error) error {
	for _, entry := range rootTree.Entries {
		repoPath := cleanRepoPath(repoPrefix, entry.Name)
		virtualPath := joinVirtual(virtualPrefix, entry.Name)

		if entry.Mode == filemode.Dir {
			subTree, err := object.GetTree(b.repo.Storer, entry.Hash)
			if err != nil {
				return err
			}
			if err := b.walkFiles(subTree, repoPath, virtualPath, fn); err != nil {
				return err
			}
			continue
		}

		content, size, err := b.readBlob(&entry)
		if err != nil {
			return err
		}
		if err := fn(fileRecord{
			virtual:  virtualPath,
			repoPath: repoPath,
			content:  content,
			size:     size,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (b *GitBackend) subtreeAt(virtualPath string) (*object.Tree, string, string, error) {
	virtualPath = normalizePath(virtualPath)
	if virtualPath == "/" {
		return b.baseTree, cleanRepoPath(b.cfg.BaseDir), "/", nil
	}

	node, err := b.resolveNode(virtualPath)
	if err != nil {
		return nil, "", "", err
	}
	if node.tree == nil {
		return nil, "", "", fmt.Errorf("not a directory: %s", virtualPath)
	}
	return node.tree, node.repoPath, node.virtual, nil
}

// LsInfo 列出目录下的直接子项。
func (b *GitBackend) LsInfo(ctx context.Context, req *filesystem.LsInfoRequest) ([]filesystem.FileInfo, error) {
	if err := b.ensureReady(ctx); err != nil {
		return nil, err
	}

	node, err := b.resolveNode(req.Path)
	if err != nil {
		return nil, err
	}
	if node.tree == nil {
		return nil, fmt.Errorf("not a directory: %s", normalizePath(req.Path))
	}

	result := make([]filesystem.FileInfo, 0, len(node.tree.Entries))
	for _, entry := range node.tree.Entries {
		info := filesystem.FileInfo{
			Path:       entry.Name,
			IsDir:      entry.Mode == filemode.Dir,
			ModifiedAt: b.modifiedAt,
		}
		if !info.IsDir {
			blob, err := b.repo.BlobObject(entry.Hash)
			if err == nil {
				info.Size = blob.Size
			}
		}
		result = append(result, info)
	}
	return result, nil
}

// Read 读取文件内容，支持按行 offset/limit。
func (b *GitBackend) Read(ctx context.Context, req *filesystem.ReadRequest) (*filesystem.FileContent, error) {
	if err := b.ensureReady(ctx); err != nil {
		return nil, err
	}

	node, err := b.resolveNode(req.FilePath)
	if err != nil {
		return nil, err
	}
	if node.entry == nil || node.entry.Mode == filemode.Dir {
		return nil, fmt.Errorf("file not found: %s", normalizePath(req.FilePath))
	}

	content, _, err := b.readBlob(node.entry)
	if err != nil {
		return nil, err
	}
	return sliceContent(content, req.Offset, req.Limit), nil
}

// GrepRaw 在文件内容中搜索正则 pattern。
func (b *GitBackend) GrepRaw(ctx context.Context, req *filesystem.GrepRequest) ([]filesystem.GrepMatch, error) {
	if err := b.ensureReady(ctx); err != nil {
		return nil, err
	}
	if req.Pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	pattern := req.Pattern
	if req.CaseInsensitive {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	searchVirtual := "/"
	if req.Path != "" {
		searchVirtual = normalizePath(req.Path)
	}

	subTree, repoPrefix, virtualPrefix, err := b.subtreeAt(searchVirtual)
	if err != nil {
		return nil, err
	}

	var files []fileRecord
	if err := b.walkFiles(subTree, repoPrefix, virtualPrefix, func(rec fileRecord) error {
		files = append(files, rec)
		return nil
	}); err != nil {
		return nil, err
	}

	files = filterFiles(files, req)
	var matches []filesystem.GrepMatch
	for _, file := range files {
		fileMatches := grepFile(file, re, req)
		matches = append(matches, fileMatches...)
	}
	return matches, nil
}

// GlobInfo 按 glob 模式匹配文件。
func (b *GitBackend) GlobInfo(ctx context.Context, req *filesystem.GlobInfoRequest) ([]filesystem.FileInfo, error) {
	if err := b.ensureReady(ctx); err != nil {
		return nil, err
	}

	baseVirtual := "/"
	if req.Path != "" {
		baseVirtual = normalizePath(req.Path)
	}

	subTree, repoPrefix, virtualPrefix, err := b.subtreeAt(baseVirtual)
	if err != nil {
		return nil, err
	}

	isAbsolutePattern := strings.HasPrefix(req.Pattern, "/")
	var result []filesystem.FileInfo

	err = b.walkFiles(subTree, repoPrefix, virtualPrefix, func(rec fileRecord) error {
		var matchPath string
		var resultPath string

		if isAbsolutePattern {
			matchPath = rec.virtual
			resultPath = rec.virtual
		} else {
			if baseVirtual == "/" {
				matchPath = strings.TrimPrefix(rec.virtual, "/")
			} else {
				matchPath = strings.TrimPrefix(rec.virtual, baseVirtual+"/")
				if matchPath == rec.virtual {
					matchPath = strings.TrimPrefix(rec.virtual, "/")
				}
			}
			resultPath = matchPath
		}

		matched, err := doublestar.Match(req.Pattern, matchPath)
		if err != nil {
			return fmt.Errorf("invalid glob pattern: %w", err)
		}
		if matched {
			result = append(result, filesystem.FileInfo{
				Path:       resultPath,
				IsDir:      false,
				Size:       rec.size,
				ModifiedAt: b.modifiedAt,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Write 暂不支持，返回只读错误。
func (b *GitBackend) Write(ctx context.Context, req *filesystem.WriteRequest) error {
	if err := b.ensureReady(ctx); err != nil {
		return err
	}
	return errReadOnly
}

// Edit 暂不支持，返回只读错误。
func (b *GitBackend) Edit(ctx context.Context, req *filesystem.EditRequest) error {
	if err := b.ensureReady(ctx); err != nil {
		return err
	}
	return errReadOnly
}

func sliceContent(content string, offset, limit int) *filesystem.FileContent {
	if offset < 1 {
		offset = 1
	}
	lineOffset := offset - 1

	if lineOffset == 0 && limit <= 0 {
		return &filesystem.FileContent{Content: content}
	}

	if lineOffset == 0 && limit > 0 {
		lineCount := strings.Count(content, "\n") + 1
		if lineCount <= limit {
			return &filesystem.FileContent{Content: content}
		}
	}

	start := 0
	for i := 0; i < lineOffset; i++ {
		idx := strings.IndexByte(content[start:], '\n')
		if idx == -1 {
			return &filesystem.FileContent{}
		}
		start += idx + 1
	}

	if limit <= 0 {
		return &filesystem.FileContent{Content: content[start:]}
	}

	end := start
	for i := 0; i < limit; i++ {
		idx := strings.IndexByte(content[end:], '\n')
		if idx == -1 {
			return &filesystem.FileContent{Content: content[start:]}
		}
		end += idx + 1
	}
	return &filesystem.FileContent{Content: content[start : end-1]}
}

func filterFiles(files []fileRecord, req *filesystem.GrepRequest) []fileRecord {
	if req.Glob == "" && req.FileType == "" {
		return files
	}

	filtered := make([]fileRecord, 0, len(files))
	for _, file := range files {
		if req.Glob != "" {
			matchPath := strings.TrimPrefix(file.virtual, "/")
			matched, err := doublestar.Match(req.Glob, filepath.Base(matchPath))
			if err != nil || !matched {
				if !strings.Contains(req.Glob, "/") && !strings.Contains(req.Glob, "**") {
					continue
				}
				matched, err = doublestar.Match(req.Glob, matchPath)
				if err != nil || !matched {
					continue
				}
			}
		}
		if req.FileType != "" && !matchFileType(fileExt(file.virtual), req.FileType) {
			continue
		}
		filtered = append(filtered, file)
	}
	return filtered
}

func grepFile(file fileRecord, re *regexp.Regexp, req *filesystem.GrepRequest) []filesystem.GrepMatch {
	if req.EnableMultiline {
		return grepMultiline(file, re)
	}
	return grepSingleLine(file, re)
}

func grepSingleLine(file fileRecord, re *regexp.Regexp) []filesystem.GrepMatch {
	lines := strings.Split(file.content, "\n")
	matches := make([]filesystem.GrepMatch, 0)
	for i, line := range lines {
		if re.MatchString(line) {
			matches = append(matches, filesystem.GrepMatch{
				Path:    file.virtual,
				Line:    i + 1,
				Content: line,
			})
		}
	}
	return matches
}

func grepMultiline(file fileRecord, re *regexp.Regexp) []filesystem.GrepMatch {
	matches := make([]filesystem.GrepMatch, 0)
	indices := re.FindAllStringIndex(file.content, -1)
	lines := strings.Split(file.content, "\n")
	for _, match := range indices {
		startLine := 1 + strings.Count(file.content[:match[0]], "\n")
		endLine := 1 + strings.Count(file.content[:match[1]], "\n")
		for lineNum := startLine; lineNum <= endLine && lineNum <= len(lines); lineNum++ {
			matches = append(matches, filesystem.GrepMatch{
				Path:    file.virtual,
				Line:    lineNum,
				Content: lines[lineNum-1],
			})
		}
	}
	return matches
}

func matchFileType(ext, fileType string) bool {
	typeMap := map[string][]string{
		"go": {"go"},
		"md": {"markdown", "md", "mdown", "mdwn", "mdx", "mkd", "mkdn"},
		"py": {"py", "pyi"},
		"js": {"cjs", "js", "jsx", "mjs", "vue"},
		"ts": {"cts", "mts", "ts", "tsx"},
	}
	if exts, ok := typeMap[fileType]; ok {
		for _, e := range exts {
			if ext == e {
				return true
			}
		}
	}
	return ext == fileType
}

func fileExt(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
