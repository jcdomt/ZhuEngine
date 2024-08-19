package config

import (
	"errors"
	"fmt"
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
	err = os.WriteFile("./conf/runtime.ini", []byte("# 这是配置解析系统运行时产生的文件\n"+str), 0666)
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

	tokens := getIniTokens(file_str)
	//str := origin

	for i, token := range tokens {
		index := strconv.Itoa(i + 1)

		if strings.HasPrefix(token, "#include ") {
			file_str, err = parseIniInclude(token, file_str)
			if err != nil {
				return "", errors.New(filepath + ": " + index + "行：" + err.Error())
			}
		}

		if strings.HasPrefix(token, "#schedule ") {
			file_str, err = parseIniScheduleTable(token, file_str, tokens, i)
			if err != nil {
				return "", errors.New(filepath + ": " + index + "行：" + err.Error())
			}
		}
	}
	return file_str, nil
}

func parseIniInclude(token string, origin string) (string, error) {
	re := regexp.MustCompile(`#include\s+(.+)`)
	str := origin

	// 查找匹配项
	matches := re.FindStringSubmatch(token)

	if len(matches) > 1 {
		// 提取出文件路径
		filepath := matches[1]

		temp_str, err := parseMain(filepath)
		if err != nil {
			return "", errors.New("文件路径解析错误：" + err.Error())
		}
		str = strings.Replace(origin, token, temp_str, 1)
	} else {
		return "", errors.New("未找到文件路径")
	}

	return str, nil
}

// 解析调度器配置表
func parseIniScheduleTable(head_token string, origin string, tokens []string, index int) (string, error) {
	table_name := strings.Fields(head_token)[1]
	table_content := ""
	for i := index + 1; i < len(tokens); i++ {
		token := tokens[i]
		if token == "#table_end" {
			break
		}
		a := strings.Fields(token)
		// a[0] 是 #
		if a[1] == "ip" && a[2] == "weight" {
			// 表头，无所谓
			continue
		}
		table_content = table_content + fmt.Sprintf("%s?%s,", a[1], a[2])

	}
	SetScheduleTable(table_name, strings.TrimRight(table_content, ","))
	return origin, nil
}

func getIniTokens(str string) []string {
	str_arr := strings.Split(str, "\n")
	tokens := make([]string, 0)

	for _, token := range str_arr {
		token = strings.TrimSuffix(token, "\r\n")
		token = strings.TrimSuffix(token, "\r")
		token = strings.TrimSuffix(token, "\n")
		tokens = append(tokens, token)
	}

	return tokens
}
