package git

import "errors"

var (
	// errReadOnly 表示 GitBackend 处于只读模式，不支持写操作。
	errReadOnly = errors.New("git backend is read-only")

	// errPathEscape 表示请求路径试图穿越 baseDir 边界。
	errPathEscape = errors.New("path escapes base directory")

	// errNotFound 表示 tree 中不存在目标路径。
	errNotFound = errors.New("path not found")
)
