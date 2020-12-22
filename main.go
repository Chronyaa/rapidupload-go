package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// 拖动文件时，计算秒传链接，结果输出到result.txt中
	if len(os.Args) > 1 {
		ResultFile, Err := os.Create("result.txt")
		if Err != nil {
			log.Fatalln("Cannot create result.txt", Err)
		}
		defer ResultFile.Close()
		for i := 1; i < len(os.Args); i++ {
			Link, Err := GenerateRapidUploadLink(os.Args[i])
			if Err != nil {
				log.Fatalln("Error occurred while generating md5 of", os.Args[i], Err)
			}
			fmt.Fprintln(ResultFile, Link)
		}
		return
	}
	// 尝试读取bduss
	Buffer, Err := ioutil.ReadFile("bduss.txt")
	Bduss := ""
	if Err != nil {
		fmt.Println("bduss不存在，请先输入BDUSS：")
		fmt.Println("首先打开pan.baidu.com，并登录；")
		fmt.Println("Firefox浏览器：按下Shift+F9，在弹出的开发者工具中找到'名称'为BDUSS的那一行，双击该行的'值'，按Ctrl+C复制")
		fmt.Println("Chrome浏览器：按下Ctrl+Shift+I，在弹出的开发者工具中点击上方横条的Application，然后点击左侧的Storage下的Cookies，按左侧的▲展开，点击pan.baidu.com那一行，找到'Name'为BDUSS的那一行，双击该行的'Value'，按Ctrl+C复制")
		fmt.Println("Edge浏览器：目前Edge浏览器用的是Chromium内核，步骤与Chrome浏览器类似，如果弹出的开发者工具位于右侧，顶部横条可能不会显示Application，请点开顶部横条中的⏩图标寻找")
		fmt.Println("其它浏览器：许多国产浏览器使用的也是Chromium内核，因此步骤也与Chrome浏览器类似，找不到请使用Windows自带的Edge浏览器")
		for {
			fmt.Scanln(&Bduss)
			if !NewBaiduClient(Bduss).CheckCookie() {
				fmt.Println("输入的BDUSS有误！请认真查看上述的教程并重新输入！")
			} else {
				ioutil.WriteFile("bduss.txt", []byte(Bduss), 0644)
				break
			}
		}
	} else {
		Bduss = string(Buffer)
		if !NewBaiduClient(Bduss).CheckCookie() {
			fmt.Println("Bduss有误！请删除bduss.txt并重新打开本程序！")
			return
		}
	}
	Path := ""
	Scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("网盘路径（按回车确认，留空时为根目录）：")
		if Scanner.Scan() {
			Path = Scanner.Text()
		} else {
			return
		}
		if Path == "" {
			Path = "/"
		}
		if Path[len(Path)-1] != '/' {
			Path += "/"
		}
		if Path[0] != '/' {
			fmt.Println("路径如果需要自定义请以/开头")
		} else {
			break
		}
	}
	fmt.Println("秒传链接，一行一个，按回车确认：")
	for {
		Link := ""
		if Scanner.Scan() {
			Link = Scanner.Text()
		} else {
			return
		}
		ContentMD5, SliceMD5, FileName, ContentLength, Err := ParseLink(Link)
		if Err != nil {
			fmt.Println(Err)
			continue
		}
		Success, Message := NewBaiduClient(Bduss).RapidUpload("/"+Path+FileName, ContentMD5, SliceMD5, ContentLength)
		if Success {
			fmt.Println("转存成功 =>", Message)
		} else {
			fmt.Println("转存失败：", Message)
		}
	}
}
