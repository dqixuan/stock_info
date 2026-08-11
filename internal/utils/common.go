package utils

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// GetInitials 获取单个词语的拼音首字母
func GetInitials(word string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.FirstLetter

	pinyinList := pinyin.Pinyin(word, args)

	var sb strings.Builder

	for _, py := range pinyinList {
		if len(py) > 0 {
			sb.WriteString(py[0])
		}
	}

	return strings.ToUpper(sb.String())
}

// GetInitialsList 批量获取中文词语拼音首字母
func GetInitialsList(words []string) []string {
	result := make([]string, 0, len(words))

	for _, word := range words {
		result = append(result, GetInitials(word))
	}

	return result
}
