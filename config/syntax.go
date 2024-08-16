package config

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

func parseIniConfigSyntax(filepath string) (*ini.File, error) {
	str, err := parseMain(filepath)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("./conf/runtime.ini", []byte(str), 0666)
	if err != nil {
		return nil, err
	}
	inicfg, err := ini.Load("./conf/runtime.ini")
	return inicfg, err
}

func parseMain(filepath string) (string, error) {
	// 读入文件
	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	file_str := string(b)

	file_str, err = parseIniInclude(file_str, filepath)
	if err != nil {
		return "", err
	}

	return file_str, nil
}

func parseIniInclude(origin string, path string) (string, error) {
	str_arr := strings.Split(origin, "\n")
	str := origin

	for i, token := range str_arr {
		index := strconv.Itoa(i)
		token = strings.TrimSuffix(token, "\r\n")
		token = strings.TrimSuffix(token, "\r")
		token = strings.TrimSuffix(token, "\n")

		if strings.HasPrefix(token, "#include ") {
			re := regexp.MustCompile(`#include\s+(.+)`)

			// 查找匹配项
			matches := re.FindStringSubmatch(token)

			if len(matches) > 1 {
				// 提取出文件路径
				filepath := matches[1]

				temp_str, err := parseMain(filepath)
				if err != nil {
					return "", err
				}
				str = strings.Replace(str, token, temp_str, 1)
			} else {
				return "", errors.New(path + ": " + index + "行：未找到文件路径")
			}
		}
	}

	return str, nil
}
