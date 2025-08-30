package util

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/zan8in/gologger"
)

// BuildLoaderFile 编译生成的Go文件为可执行文件
func BuildLoaderFile(goFileName, exeFileName string) error {
	if goFileName == "" || exeFileName == "" {
		return fmt.Errorf("文件名不能为空")
	}

	fileName := filepath.Base(exeFileName)
	// 编译参数
	buildFlags := "-w -s -H windowsgui"
	buildMode := "-trimpath"
	args := []string{"build", "-o", fileName, "-ldflags", buildFlags, buildMode}

	cmd := exec.Command("go", args...)
	// 设置工作目录为Go文件所在目录
	workingDir := filepath.Dir(goFileName)
	cmd.Dir = workingDir

	// 捕获输出用于错误诊断
	output, err := cmd.CombinedOutput()
	if err != nil {
		gologger.Error().Msgf("捆绑文件编译失败: %v\n输出: %s", err, string(output))
		return fmt.Errorf("编译失败: %w", err)
	}

	gologger.Info().Msgf("捆绑文件编译成功: %s", exeFileName)
	return nil
}
