package main

import (
	"errors"
	"strings"
)

// ParseLink 解析秒传链接
func ParseLink(Link string) (ContentMD5, SliceMD5, FileName, ContentLength string, Err error) {
	Fields := strings.Split(Link, "#")
	if len(Fields) < 4 {
		Err = errors.New(`不是有效的秒传链接！`)
		return
	}
	ContentMD5 = Fields[0]
	SliceMD5 = Fields[1]
	ContentLength = Fields[2]
	FileName = Fields[3]
	return
}
