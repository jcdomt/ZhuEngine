package router

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

func SiteAutoRun(sites []*Site) {
	for _, v := range sites {
		if v.Config.Autorun {
			var err error
			if !isInSameDir(v.Config.Exec, "./") {
				strs := strings.Split(v.Config.Exec, "/")
				execName := strs[len(strs)-1]
				parentDir := getParentDir(v.Config.Exec)

				fmt.Println("Changing directory to:", parentDir)

				var cmd *exec.Cmd
				if runtime.GOOS == "windows" {
					cmd = exec.Command("cmd.exe", "/C", "start", execName)
				} else {
					cmd = exec.Command("./" + execName)
				}

				cmd.Dir = parentDir
				err := cmd.Run()
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				cmd := exec.Command(v.Config.Exec)
				err := cmd.Run()
				if err != nil {
					fmt.Println("Error:", err)
				}
			}

			if err != nil {
				log.Error(v.Config.Exec, "启动失败：", err)
			} else {
				log.Info(v.Config.Exec, "启动成功")
			}
		}
	}
}

// isInSameDir 判断两个路径是否在同一文件夹
func isInSameDir(dir1, dir2 string) bool {
	// 规范化路径
	dir1 = normalizePath(dir1)
	dir2 = normalizePath(dir2)

	// 获取父目录
	parentDir1 := getParentDir(dir1)
	parentDir2 := getParentDir(dir2)

	// 比较父目录是否相同
	return parentDir1 == parentDir2
}

// normalizePath 规范化路径，移除多余的路径分隔符
func normalizePath(path string) string {
	parts := strings.Split(path, "/")
	var normalizedParts []string
	for _, part := range parts {
		if part != "" && part != "." {
			normalizedParts = append(normalizedParts, part)
		}
	}
	return strings.Join(normalizedParts, "/")
}

// getParentDir 获取路径的父目录
func getParentDir(path string) string {
	if path == "/" {
		return "/"
	}
	parts := strings.Split(path, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}
