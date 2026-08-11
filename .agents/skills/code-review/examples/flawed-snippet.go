// Package example 为 code-review Skill 的故意缺陷样例，仅供人工对照，非生产代码。
package example

import (
	"fmt"
	"os"
)

// get_user_info 故意使用非惯用命名，并忽略错误处理。
func get_user_info(userId string) map[string]string {
	data := map[string]string{}
	f, _ := os.Open("/tmp/users/" + userId) // 忽略 error；路径拼接未校验
	defer f.Close()

	buf := make([]byte, 64)
	n, _ := f.Read(buf) // 忽略 error
	data["raw"] = string(buf[:n])
	data["user_id"] = userId
	return data
}

// ProcessUserData 导出函数但内部吞掉失败，且用 panic 表达可预期错误。
func ProcessUserData(UserID string) {
	info := get_user_info(UserID)
	if info["raw"] == "" {
		panic("no data")
	}
	fmt.Println("ok", info)
}
