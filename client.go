package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// BaiduClient 百度客户端
type BaiduClient struct {
	Cookie string
}

// RequestContext 要发送请求的上下文
type RequestContext struct {
	URL       string
	Method    string
	Cookie    string
	Referer   string
	Body      string
	UserAgent string
}

// DoRequest 发送请求
func DoRequest(Ctx *RequestContext) (Result []byte, Err error) {
	var Body io.Reader
	if Ctx.Method == "" {
		Ctx.Method = http.MethodGet
	} else if Ctx.Method == http.MethodPost {
		Body = bytes.NewReader([]byte(Ctx.Body))
	}
	Req, Err := http.NewRequest(Ctx.Method, Ctx.URL, Body)
	if Err != nil {
		log.Fatalln("Create Request Failed:", *Ctx)
	}
	if Ctx.Cookie != "" {
		Req.Header.Add("Cookie", Ctx.Cookie)
	}
	if Ctx.Referer != "" {
		Req.Header.Add("Referer", Ctx.Referer)
	}
	if Ctx.UserAgent != "" {
		Req.Header.Add("User-Agent", Ctx.UserAgent)
	} else {
		Req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0")
	}
	Req.Header.Add("Accept-Language", "zh-CN,en;q=0.9,ja;q=0.7,zh;q=0.6,zh-TW;q=0.4,zh-HK;q=0.3,en-US;q=0.1")
	Req.Header.Add("Accept", "*/*")
	Req.Header.Add("Connection", "close")
	Resp, Err := http.DefaultClient.Do(Req)
	if Err != nil {
		log.Println("DoRequest Failed:", Err)
		return nil, Err
	}
	defer Resp.Body.Close()
	return ioutil.ReadAll(Resp.Body)
}

// NewBaiduClient 创建一个百度客户端对象
func NewBaiduClient(Bduss string) *BaiduClient {
	return &BaiduClient{
		Cookie: `BDUSS=` + Bduss + `;`,
	}
}

// CheckCookie 测试Cookie是否有效
func (B *BaiduClient) CheckCookie() bool {
	Buffer, Err := DoRequest(&RequestContext{
		URL:    "http://tieba.baidu.com/dc/common/tbs",
		Cookie: B.Cookie,
	})
	if Err != nil {
		log.Println("CheckCookie Failed:", Err)
		return false
	}
	Result := &struct {
		IsLogin int `json:"is_login"`
	}{}
	Err = json.Unmarshal(Buffer, Result)
	if Err != nil {
		log.Println("CheckCookie Failed While Decoding Buffer:", Err, string(Buffer))
		return false
	}
	return Result.IsLogin == 1
}

// RapidUpload 秒传接口
func (B *BaiduClient) RapidUpload(Path, ContentMD5, SliceMD5, ContentLength string) (Success bool, Message string) {
	UploadURL := `http://pan.baidu.com/api/rapidupload?clienttype=6&version=2.0.0.3`
	Body := fmt.Sprintf(`path=%s&content-md5=%s&slice-md5=%s&content-length=%s`, Path, ContentMD5, SliceMD5, ContentLength)
	Buffer, Err := DoRequest(&RequestContext{
		URL:       UploadURL,
		Method:    http.MethodPost,
		Cookie:    B.Cookie,
		UserAgent: `netdisk;2.0.0.3;PC;PC-Windows;10.0.16299;uploadplugin`,
		Body:      Body,
	})
	Result := &struct {
		Errno int `json:"errno"`
		Info  struct {
			Path string `json:"path"`
		} `json:"info"`
	}{}
	Err = json.Unmarshal(Buffer, Result)
	if Err != nil {
		log.Println("RapidUpload Failed While Decoding:", Err, string(Buffer))
		return false, "解析返回值时出错：" + string(Buffer)
	}
	if Result.Errno == 0 {
		return true, Result.Info.Path
	}
	switch Result.Errno {
	case 404:
		Message = "链接已失效"
	case -6:
		Message = "BDUSS已失效"
	case -8:
		Message = "已存在同名文件"
	default:
		Message = "转存失败：" + fmt.Sprint(Result.Errno)
	}
	return false, Message
}
