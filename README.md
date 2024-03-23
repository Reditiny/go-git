# go-git
本库尝试使用 go 语言实现一个简易 git 工具，供自学使用。
内容参考了 Thibault Polge 的教程 https://wyag.thb.lt/ 其中的实现为 python 版本。

# 项目结构
cmd 中为 cobra 实现的命令行操作
common 中为主要的运行的方法

# 内容介绍
go-git 复刻了 git 的版本控制原理，实现上更为简单，是一个玩具级别的版本控制工具。
目前支持 init、hash-object、cat-file、log、ls-tree、tag、show-ref、ls-file、checkout、add、rm、commit 等命令的基本功能。
仅就已支持的功能而言 go-git 可以与 git 无缝衔接使用，因此不做过多陈述。若想自己实现建议参考上面的链接。
