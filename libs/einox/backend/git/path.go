package git

import (
	"path/filepath"
	"strings"
)

// normalizePath 规范化虚拟路径，与 InMemoryBackend 一致。
func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

// cleanRepoPath 将仓库内路径统一为 POSIX 分隔符。
func cleanRepoPath(parts ...string) string {
	joined := filepath.Join(parts...)
	joined = filepath.ToSlash(joined)
	return strings.Trim(joined, "/")
}

// virtualToRepo 将对外虚拟路径映射为仓库 tree 内路径。
func virtualToRepo(baseDir, virtual string) (string, error) {
	virtual = normalizePath(virtual)
	rel := strings.TrimPrefix(virtual, "/")
	if rel == "." {
		rel = ""
	}

	base := cleanRepoPath(baseDir)
	repoPath := cleanRepoPath(base, rel)

	if base != "" && repoPath != "" && repoPath != base {
		if !strings.HasPrefix(repoPath, base+"/") {
			return "", errPathEscape
		}
	}

	return repoPath, nil
}

// repoToVirtual 将仓库内路径转换为对外虚拟路径。
func repoToVirtual(baseDir, repoPath string) string {
	base := cleanRepoPath(baseDir)
	repoPath = cleanRepoPath(repoPath)

	if base == "" {
		if repoPath == "" {
			return "/"
		}
		return normalizePath("/" + repoPath)
	}

	if repoPath == base {
		return "/"
	}
	if strings.HasPrefix(repoPath, base+"/") {
		return normalizePath("/" + strings.TrimPrefix(repoPath, base+"/"))
	}
	return normalizePath("/" + repoPath)
}

// joinVirtual 拼接虚拟路径片段。
func joinVirtual(baseVirtual, name string) string {
	baseVirtual = normalizePath(baseVirtual)
	if baseVirtual == "/" {
		return normalizePath("/" + name)
	}
	return normalizePath(baseVirtual + "/" + name)
}
