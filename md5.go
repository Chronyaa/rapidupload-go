package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SliceSize 分片大小
const SliceSize = 256 * 1024

// ComputeMD5 计算MD5并返回字符串表示
func ComputeMD5(Reader io.Reader) (ContentMD5, SliceMD5 string) {
	Hasher := md5.New()
	io.CopyN(Hasher, Reader, SliceSize)
	SliceMD5 = strings.ToUpper(hex.EncodeToString(Hasher.Sum(nil)))
	io.CopyBuffer(Hasher, Reader, make([]byte, SliceSize))
	ContentMD5 = strings.ToUpper(hex.EncodeToString(Hasher.Sum(nil)))
	return
}

// GenerateRapidUploadLink 生成秒传链接
func GenerateRapidUploadLink(FileName string) (Link string, Err error) {
	File, Err := os.Open(FileName)
	if Err != nil {
		return
	}
	defer File.Close()
	Stat, _ := File.Stat()
	ContentMD5, SliceMD5 := ComputeMD5(File)
	return ContentMD5 + "#" + SliceMD5 + "#" + fmt.Sprint(Stat.Size()) + "#" + filepath.Base(FileName), nil
}
