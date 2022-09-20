<!-- 参考：https://zhuanlan.zhihu.com/p/453462046 -->
Ubuntu安装go/升级go版本【转载】
Mr.Hachi
Mr.Hachi
​
安徽工业大学 计算机技术硕士
​关注他
15 人赞同了该文章
1.若系统之前存在旧版本的go，无则跳过此步骤

1、sudo rm -rf /usr/local/go
2、sudo apt-get remove golang
3、sudo apt-get remove golang-go
4、sudo apt-get autoremove
2. 获取安装包
#wget 后面的下载链接请去golang官网(https://golang.google.cn/dl/)获取你想下载的对应go版本
sudo wget https://golang.google.cn/dl/go1.18.5.linux-amd64.tar.gz
# 解压文件
sudo tar xfz go1.18.5.linux-amd64.tar.gz -C /usr/local
3.设置环境变量

打开：

sudo vim /etc/profile
将以下内容追加到文件末尾

export GOROOT=/usr/local/go
export GOPATH=$HOME/gowork
export GOBIN=$GOPATH/bin
export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
输入以下命令保存

:wq
4. 使环境变量生效

 source /etc/profile
如果只是这样做，在关闭终端后，重新打开环境变量又会失效，除了重新启动系统之外，可以在用户根目录的.bashrc

cd ~
sudo vim .bashrc
在文件末尾加入如下命令

source /etc/profile
5. 查看环境是否搭建成功
go env

6.开启GO111MOUDLE和更改GOPROXY

go env -w GOPROXY="https://goproxy.cn"
go env -w GO111MODULE=on
7.vscode安装go插件一些问题代理问题
1）代理问题

Mr.Hachi：The "gopls" command is not available. Run "go get -v golang.org/x/tools/gopls" to install.【转载】
13 赞同 · 3 评论文章
2）没有安装gcc

Ubuntu如何安装最新版安装gcc
​www.jianshu.com/p/9d7e83928a08

参考博客：

Ubuntu下升级Golang版本_nonoli287的博客-CSDN博客_ubuntu升级golang
​blog.csdn.net/nonoli287/article/details/109255152

https://www.xiaoheidiannao.com/227783.html
​www.xiaoheidiannao.com/227783.html
ubuntu16.04安装指定版本的Go环境 - Go语言中文网 - Golang中文社区
​studygolang.com/articles/18790
GO111MODULE的设置（及GOPROXY） - pu369com - 博客园
​www.cnblogs.com/pu369/p/12068645.html