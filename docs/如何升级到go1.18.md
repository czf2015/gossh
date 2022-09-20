<!-- 参考：https://zhuanlan.zhihu.com/p/482133602 -->
昨天晚上（2022.03.15）go 1.18正式版发布，go1.18为我们带来了 泛型、模糊测试、工作空间 等好多新特性，是大家期待比较高的一个go版本。接下来就让我们升级到全新的go1.18版本吧。

go官方提供了一个 go 的版本管理工具 golang.org/dl

使用方法是：

```sh
go install golang.org/dl/go1.18@latest  # 安装下载程序
go1.18 download                         # 下载并安装
go1.18 version                          # 检查 Go 版本
```
但是这个库是从 http://dl.google.com/go 下载Go语言的安装包，由于国内网络问题，基本下载不下来。

为了方便自己的使用， 因此fork了这个仓库，并进行了修改，更改下载地址为国内的清华大学源，方便国内下载安装。

仓库地址 https://github.com/SunJary/dl

这里目前有3个分支：

原始分支master： 下载地址为 https://dl.google.com/go，国内基本不可用
默认分支（清华大学源）： 下载地址为 https://mirrors.ustc.edu.cn/golang/ ，国内下载速度较快
mirror分支：下载地址为 https://gomirrors.org/，go官方镜像地址，国内下载速度也较慢，另外不支持sha256文件验证
使用说明
基本跟 golang.org/dl 的使用方式一致

命令行下运行 
```sh
go install http://github.com/SunJary/dl/go<version>@ustc
```

原文：
```sh
go install http://github.com/SunJary/dl/go1.18@ustc # 使用清华大学源下载 ustc
更新：2022年4月13日发布了 go1.18.1 版本，使用@USTC标签无法下载最新版本的go1.18.1，应改为@latest

go install github.com/SunJary/dl/go1.18@latest     # 默认使用清华大学源下载 ustc

go install github.com/SunJary/dl/go1.18@ustc
go: downloading github.com/SunJary/dl v0.0.0-20220316004807-aba2b4152c32
```
也可以手动指定其他分支比如
```sh
go install http://github.com/SunJary/dl/go1.18@mirror # 使用 https://gomirrors.org/ 进行下载，不推荐

go1.18 download
ubuntu@VM-8-7-ubuntu:~$ go1.18 download
Download From 'https://mirrors.ustc.edu.cn/golang/go1.18.linux-amd64.tar.gz'
Downloaded   0.0% (     8192 / 141702072 bytes) ...
Downloaded  15.7% ( 22265856 / 141702072 bytes) ...
Downloaded  24.5% ( 34668544 / 141702072 bytes) ...
Downloaded  33.8% ( 47849472 / 141702072 bytes) ...
Downloaded  42.5% ( 60162048 / 141702072 bytes) ...
Downloaded  51.6% ( 73146368 / 141702072 bytes) ...
Downloaded  60.4% ( 85565440 / 141702072 bytes) ...
Downloaded  69.4% ( 98402304 / 141702072 bytes) ...
Downloaded  77.4% (109658112 / 141702072 bytes) ...
Downloaded  87.2% (123543552 / 141702072 bytes) ...
Downloaded  96.1% (136167424 / 141702072 bytes) ...
Downloaded 100.0% (141702072 / 141702072 bytes)
Unpacking /home/ubuntu/sdk/go1.18/go1.18.linux-amd64.tar.gz ...
Success. You may now run 'go1.18'
```
这时我们可以验证一下我们的go安装是否成功，运行 go1.18 version
```sh
ubuntu@VM-8-7-ubuntu:~$ go1.18 version
go version go1.18 linux/amd64
```
到这里我们就成功将go1.18安装到了我们的GOPATH/bin目录下。并且我们可以在GOPATH/bin 目录下找到刚刚安装的go1.18的可执行文件。

将go1.18设置为默认go版本
每次运行都要运行 go1.18 xxx 有点麻烦，我们可以让go1.18作为我们的默认版本。

每次运行都要运行 go1.18 xxx 有点麻烦，我们可以让go1.18作为我们的默认版本。

查看GOROOT路径
```sh
go1.18 env GOROOT
```
C:\Users\你的用户名\sdk\go1.18
我们可以在 GOROOT\bin 目录下找到 go 的可执行文件，现在我们只需要将其添加到操作系统的环境变量当中就可以全局运行1.18版本的go了。

Windows下修改环境变量
将环境变量中原来的PATH中的go安装目录改为 C:\Users\你的用户名\sdk\go1.18\bin 即可

Linux下修改环境变量
在这之前我们安装的go是在 /usr/local/go 目录下

并在/etc/profile中每次开机export GOROOT等变量，通过这种方法将go命令添加到环境变量中，以便后期使用
```sh
export GOROOT=/usr/local/go
export GOPATH=/home/ubuntu/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
现在我们改成
```sh
export GOPATH=/home/ubuntu/go         # 先定义GOPATH，这里视自己的情况而定
export PATH=$PATH:${GOPATH}/bin       # 将 GOPATH 加入到环境变量，这样就能执行刚刚安装的go1.18了

export GOROOT=$(go1.18 env GOROOT)    # 设置GOROOT
export PATH=$PATH:${GOROOT}/bin       # 将新的GOROOT加入环境变量
export GOPROXY=https://goproxy.cn
ubuntu@VM-8-7-ubuntu:~$ go version
go version go1.18 linux/amd64
```
这样子，go1.18就成了我们的默认版本



原文链接：

[Go 升级到1.18 - TryDoTop](​blog.trydo.top/p/316)
