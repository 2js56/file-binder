package binder

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zan8in/goflags"
	"github.com/zan8in/gologger"
)

// 1.定义参数结构体
type (
	Options struct {
		Payload string
		File    string
		Output  string
	}
)

// 2.解析参数
func ParseOptions() *Options {
	options := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`binder`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Payload, "payload", "p", "", "input a payload.exe"),
		flagSet.StringVarP(&options.File, "file", "f", "", "input a normal file"),
		flagSet.StringVarP(&options.Output, "output", "o", "result", "save output path"),
	)

	flagSet.Parse()
	// 解析参数是否有值
	err := options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("Program exiting: %s\n", err)
	}
	return options
}

var (
	errNoFile1 = errors.New("no input payload.exe file")
	errNoFile2 = errors.New("no input normal file")
	errNoInput = errors.New("no input,please binder.exe -h")
)

// validateOptions 验证输入参数的有效性
func (options *Options) validateOptions() error {
	if options.File == "" && options.Payload == "" {
		return errNoInput
	}

	if options.Payload == "" {
		return errNoFile1
	}

	if options.File == "" {
		return errNoFile2
	}

	// 验证文件扩展名
	if err := options.validateFileExtensions(); err != nil {
		return err
	}

	// 验证输出路径
	if err := options.validateOutputPath(); err != nil {
		return err
	}

	return nil
}

// validateFileExtensions 验证文件扩展名
func (options *Options) validateFileExtensions() error {
	payloadExt := strings.ToLower(filepath.Ext(options.Payload))
	if payloadExt != ".exe" {
		return fmt.Errorf("payload文件必须是.exe格式，当前: %s", payloadExt)
	}

	return nil
}

// validateOutputPath 验证输出路径
func (options *Options) validateOutputPath() error {
	if options.Output == "" {
		options.Output = "result" // 默认输出目录
	}

	// 检查输出路径是否包含非法字符
	if strings.ContainsAny(options.Output, `<>:"|?*`) {
		return fmt.Errorf("输出路径包含非法字符: %s", options.Output)
	}

	return nil
}
