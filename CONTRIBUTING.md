# 贡献指南

感谢您对 File Binder 项目的关注！我们欢迎所有形式的贡献。

## 🤝 如何贡献

### 报告Bug

1. 确保bug没有被之前报告过，请查看 [Issues](https://github.com/2js56/file-binder/issues)
2. 如果没有找到相关issue，请 [创建新的issue](https://github.com/2js56/file-binder/issues/new)
3. 在issue中请包含：
   - 详细的bug描述
   - 重现步骤
   - 期望的行为
   - 实际发生的行为
   - 操作系统和Go版本
   - 相关的错误日志

### 建议功能

我们欢迎新功能建议！请：

1. 先查看现有的 [Issues](https://github.com/2js56/file-binder/issues) 确保没有重复
2. 创建一个新的issue，标题以"Feature Request:"开头
3. 详细描述建议的功能和使用场景
4. 解释为什么这个功能对项目有价值

### 提交代码

1. **Fork项目**
   ```bash
   git clone https://github.com/2js56/file-binder.git
   cd file-binder
   ```

2. **创建特性分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **编写代码**
   - 遵循现有的代码风格
   - 添加必要的注释
   - 确保代码通过所有测试

4. **测试**
   ```bash
   go test ./...
   go build cmd/binder/main.go
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "Add: your feature description"
   ```

6. **推送到GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建Pull Request**

## 📝 代码风格

### Go代码规范

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 使用有意义的变量和函数名
- 添加适当的注释，特别是公共函数

### 提交信息规范

使用清晰的提交信息格式：

```
类型: 简短描述

详细描述（可选）

关闭的issue（可选）
```

**类型示例：**
- `Add:` 新功能
- `Fix:` Bug修复
- `Update:` 更新现有功能
- `Remove:` 删除功能
- `Docs:` 文档更新
- `Style:` 代码格式化
- `Refactor:` 代码重构
- `Test:` 测试相关

**示例：**
```
Add: 支持多种文件格式捆绑

- 增加对图片文件的支持
- 添加文件类型检测
- 更新文档说明

Closes #123
```

## 🛡️ 安全考虑

由于本项目涉及安全研究，请在贡献时注意：

1. **不要提交恶意代码**
2. **确保所有功能都是为了教育目的**
3. **在pull request中说明安全影响**
4. **遵循负责任的披露原则**

## 🧪 测试指南

在提交代码前，请确保：

1. **单元测试通过**
   ```bash
   go test ./...
   ```

2. **编译成功**
   ```bash
   go build cmd/binder/main.go
   ```

3. **功能测试**
   - 测试基本的文件捆绑功能
   - 验证生成的文件能正常运行
   - 检查错误处理

## 📚 开发环境设置

1. **安装Go 1.20+**
2. **克隆项目**
   ```bash
   git clone https://github.com/2js56/file-binder.git
   cd file-binder
   ```

3. **安装依赖**
   ```bash
   go mod download
   ```

4. **运行测试**
   ```bash
   go test ./...
   ```

## 📋 项目结构

```
├── cmd/binder/          # 主程序入口
├── pkg/
│   ├── binder/          # 核心捆绑逻辑
│   ├── encode/          # 加密模块
│   ├── loader/          # 模板加载器
│   └── util/            # 工具函数
├── .github/             # GitHub Actions
├── docs/                # 文档
└── tests/               # 测试文件
```

## ❓ 疑问解答

如果您有任何疑问，可以：

1. 查看 [FAQ](docs/FAQ.md)
2. 搜索现有的 [Issues](https://github.com/2js56/file-binder/issues)
3. 创建新的issue提问

## 📄 许可证

通过贡献代码，您同意您的贡献将在 [MIT License](LICENSE) 下许可。

---

再次感谢您的贡献！🎉
