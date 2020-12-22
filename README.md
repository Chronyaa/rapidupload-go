# RapidUpload-Go

`RapidUpload-Go`用于使用/生成秒传链接，使用`Go`编写。

本地部署，无需担心泄漏BDUSS。

下载即用，无需额外安装软件。

`Windows Defender`可能会提示不明来源应用，如不放心请自行下载`Golang`环境编译项目。

## 使用方法

### 计算秒传链接

Windows GUI：拖动要计算的文件放到程序图标上，之后会在要计算的文件所在的文件夹下产生一个`result.txt`文件，多个文件时，`result.txt`中每行一个秒传链接。

Terminal：`./rapidupload-go <file1> <file2> ...`

对于大文件需要等待一段时间，在机械盘上约80M/秒。

### 保存秒传链接

打开程序，第一次打开程序时需要根据提示输入BDUSS，之后BDUSS会在程序所在的文件夹下产生一个`bduss.txt`的文件保存BDUSS。

然后根据提示输入要保存的路径与秒传链接即可。

## LICENSE

本项目基于`MIT`许可证。
