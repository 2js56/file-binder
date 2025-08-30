package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/derian/binder/pkg/loader"
	"github.com/zan8in/gologger"
)

// 1.判断文件是否存在
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("文件 %s 不存在", filePath)
	}
	return err == nil, err
}

// GenGoFileParams 生成Go文件的参数结构
type GenGoFileParams struct {
	LoaderName     string
	FileName       string
	RandomKey      string
	EncryptPayload string
	EncryptFile    string
	ResultDir      string
}

// GenGoFile 生成go文件 bin文件 模板文件 生成的文件路径
func GenGoFile(params GenGoFileParams) (string, error) {
	// 根据loaderName获取指定的shellcode模板内容
	loaderTemplate, exists := loader.Modules[params.LoaderName]
	if !exists {
		return "", fmt.Errorf("模板 %s 不存在", params.LoaderName)
	}

	loaderStr := string(loaderTemplate)

	// 替换模板文件内容
	// 处理文件名中的反斜杠转义问题
	escapedFileName := strings.ReplaceAll(params.FileName, "\\", "\\\\")

	replacements := map[string]string{
		"fileName":       escapedFileName,
		"randomKey":      params.RandomKey,
		"encryptPayload": params.EncryptPayload,
		"encryptFile":    params.EncryptFile,
	}

	for placeholder, value := range replacements {
		loaderStr = strings.ReplaceAll(loaderStr, placeholder, value)
	}

	// 使用随机go名字
	randomString, err := GenerateRandomString(10)
	if err != nil {
		return "", fmt.Errorf("生成随机文件名失败: %w", err)
	}

	outputFileName := randomString + ".go"
	resultFilePath := filepath.Join(params.ResultDir, outputFileName)

	if err := os.WriteFile(resultFilePath, []byte(loaderStr), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	gologger.Info().Msgf("文件生成成功: %s\n", resultFilePath)
	return resultFilePath, nil
}

// CleanupOldGoFiles 清理目录中的旧Go文件
func CleanupOldGoFiles(dir string) error {
	pattern := filepath.Join(dir, "*.go")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("查找Go文件失败: %w", err)
	}

	for _, file := range matches {
		if err := os.Remove(file); err != nil {
			gologger.Warning().Msgf("删除文件失败: %s, 错误: %v", file, err)
			// 继续删除其他文件，不因单个文件失败而停止
		} else {
			gologger.Info().Msgf("已删除旧文件: %s", file)
		}
	}

	return nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
