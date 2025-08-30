package main

import (
	"github.com/derian/binder/pkg/binder"
	"github.com/zan8in/gologger"
)

var (
	VERSION string
	CONTENT string
)

func main() {

	// 1.输出banner信息
	binder.ShowBanner(VERSION, CONTENT)

	// 2.解析参数
	options := binder.ParseOptions()

	// 3.开始运行捆绑器
	if err := binder.Run(options); err != nil {
		gologger.Fatal().Msg(err.Error())
	}
}
