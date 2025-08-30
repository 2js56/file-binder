# File Binder

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://www.microsoft.com/windows)

一个用Go语言编写的文件捆绑工具，支持将可执行文件与普通文件捆绑为单一可执行文件。

> ⚠️ **免责声明：本项目仅供安全研究和教育学习使用，请遵守当地法律法规，不得用于任何恶意目的。**

## 项目结构

```
binder-src/
├── cmd/binder/          # 主程序入口
├── pkg/
│   ├── binder/          # 核心捆绑逻辑
│   │   ├── banner.go    # 横幅显示
│   │   ├── config.go    # 配置常量
│   │   ├── options.go   # 命令行参数解析
│   │   └── runner.go    # 主要运行逻辑
│   ├── encode/          # 加密模块
│   │   └── encode.go    # AES加密实现
│   ├── loader/          # 模板加载器
│   │   ├── loader.go    # 模板管理
│   │   └── module/      # 模板文件
│   └── util/            # 工具函数
│       ├── build.go     # 编译工具
│       └── file.go      # 文件操作
└── result/              # 输出目录
```

## ✨ 功能特性

- 🔐 **AES-GCM加密** - 使用军用级加密算法保护文件内容
- 📦 **文件捆绑** - 将可执行文件与普通文件合并为单一PE文件
- 🛡️ **免杀技术** - 内置多种反检测和混淆技术
- 🎭 **行为模拟** - 模拟正常程序行为，降低检测概率
- 🔧 **自动编译** - 自动生成并编译Go源码
- 📁 **动态路径** - 避免硬编码路径特征
- ⚡ **异步执行** - 时间和空间分离敏感操作
- 🎯 **静默运行** - 无控制台窗口，完全隐蔽执行

## 🛠️ 技术实现

### 免杀技术
- **反调试检测** - 检测常见调试器进程（x64dbg、OllyDbg、IDA、WinDbg）
- **行为延迟** - 随机延迟执行，避免沙箱检测
- **动态路径构建** - 运行时生成随机路径，每次执行路径不同
- **函数名混淆** - 使用无害的函数名称（cleanupTemp、setFileAttribs）
- **分离执行** - 异步执行敏感操作，时间和空间分离
- **环境适应** - 智能路径选择，支持跨盘符操作
- **静默运行** - 完全无输出，避免控制台痕迹

### 加密保护
- **AES-GCM** - 256位密钥，认证加密
- **随机密钥** - 每次生成唯一32字节密钥
- **完整性验证** - 防止文件被篡改

### 代码优化
- **模块化设计** - 清晰的包结构和依赖关系
- **错误处理** - 完善的静默错误处理机制
- **资源管理** - 自动清理临时文件和资源
- **编译优化** - 使用编译参数减小文件体积

## 🚀 快速开始

### 编译

```bash
# 克隆项目
git clone https://github.com/2js56/file-binder.git
cd file-binder

# 编译
go build -o binder.exe cmd/binder/main.go
```

### 使用示例

```bash
# 基本用法
./binder.exe -p payload.exe -f document.pdf

# 指定输出目录
./binder.exe -p payload.exe -f image.jpg -o custom_output

# 查看帮助
./binder.exe -h
```

### 参数说明

| 参数 | 长参数 | 说明 | 默认值 |
|------|--------|------|--------|
| `-p` | `--payload` | 指定payload可执行文件路径 | 必需 |
| `-f` | `--file` | 指定要捆绑的普通文件路径 | 必需 |
| `-o` | `--output` | 指定输出目录 | `result` |

## 📋 系统要求

- **Go版本**: 1.20+
- **操作系统**: Windows (目标平台)
- **架构**: x64
- **权限**: 普通用户权限即可

## 🔧 工作原理

### 构建阶段
1. **文件读取** - 读取payload和普通文件内容
2. **AES加密** - 使用32字节随机密钥加密两个文件
3. **模板嵌入** - 将加密数据嵌入Go模板代码
4. **代码生成** - 生成包含解密逻辑的Go源码
5. **自动编译** - 编译为Windows PE可执行文件

### 运行阶段
1. **环境检测** - 反调试检查和行为模拟
2. **路径生成** - 动态创建随机目录和文件路径
3. **文件解密** - AES-GCM解密payload和普通文件
4. **异步执行** - 先打开普通文件，再异步执行payload
5. **痕迹清理** - 自动清理临时文件和进程

## 📁 输出文件说明

生成的可执行文件运行时会：
- **环境检测**: 检查调试器，模拟正常程序行为
- **路径随机化**: 在用户目录创建随机子目录
- **文件释放**: 解密并释放普通文件到当前目录
- **自动打开**: 使用系统默认程序打开普通文件
- **payload执行**: 异步解密并执行payload到随机路径
- **痕迹清理**: 自动清理临时文件和自身副本

### 动态路径示例
- **Payload路径**: `C:\Users\[user]\abc123\def45678.exe`
- **临时文件**: `C:\Users\[user]\AppData\Local\Temp\randomname.tmp`
- **普通文件**: `.\document.pdf` (当前目录)

## 🛡️ 安全声明

本项目采用多层安全措施：
- **加密存储** - 文件内容使用AES-GCM加密
- **动态执行** - 运行时动态生成路径和文件名
- **内存安全** - 自动资源清理，防止内存泄露
- **反检测** - 内置多种反静态和动态分析技术

## 📦 依赖项

```go
require (
    github.com/zan8in/goflags v0.0.0-20230204144650-0745934af58a
    github.com/zan8in/gologger v0.0.0-20220917062627-c34a83c0a373
)
```

## 🤝 贡献指南

欢迎提交Issue和Pull Request：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## ⚠️ 免责声明

**本工具仅供网络安全研究和教育学习使用。**

使用本工具的用户需要：
- 遵守所在地区的法律法规
- 仅在授权的测试环境中使用
- 不得用于任何非法或恶意目的
- 承担使用本工具产生的所有后果

作者不对任何滥用行为承担责任。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 [GitHub Issue](https://github.com/2js56/file-binder/issues)
- 发送邮件到: 2js56666@gmail.com

---

⭐ **如果这个项目对你有帮助，请给个Star！**
