# Git Account Manager (Go Edition)

这是一个使用 Go 语言和 Fyne 框架写的轻量级跨平台 Git 账户管理工具。

## ✨ 特性
*   **极致轻量**: 打包后的单文件可执行程序仅 **10MB-20MB**。
*   **原生性能**: 启动极快，内存占用极低。
*   **跨平台**: 完美支持 Windows, macOS, Linux。
*   **功能完整**: 支持多账户管理、SSH Key 绑定、一键切换。

## 🛠️ 构建指南

### 前置要求
*   [Go](https://go.dev/dl/) 1.20 或更高版本
*   **必须安装 C 编译器**: Fyne 依赖 GPU 渲染，需要 CGO。
    *   Windows 用户推荐安装 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) (安装时选择默认选项即可)。
    *   安装完成后，在终端运行 `gcc --version` 确认安装成功。

### 1. 运行开发版
```bash
go run main.go
```

### 2. 打包为可执行文件 (本地编译)

确保已安装 TDM-GCC，然后在项目根目录运行：

**Windows**:
```bash
go build -ldflags="-H windowsgui -s -w" -o GitAccountManager.exe main.go
```
*   `-H windowsgui`: 隐藏运行时背后的黑色命令行窗口。
*   `-s -w`: 去除调试信息和符号表，这是**减小体积的关键参数**。
*   打包后体积预期：**约 15MB**。

**macOS / Linux**:
```bash
go build -ldflags="-s -w" -o GitAccountManager main.go
```

### 3. 一键跨平台打包 (推荐方案)
如果你不想在 Windows 上安装 GCC，或者需要打包 Mac/Linux 版本，最简单的方法是使用 Docker + `fyne-cross`。

1.  安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)。
2.  安装构建工具:
    ```bash
    go install github.com/fyne-io/fyne-cross/v2/cmd/fyne-cross@latest
    ```
3.  执行打包命令:
    ```bash
    # 打包 Windows (无需本地 GCC)
    fyne-cross windows -arch=amd64

    # 打包 Linux
    fyne-cross linux -arch=amd64

    # 打包 macOS
    fyne-cross darwin -arch=amd64
    ```
    构建结果会生成在 `fyne-cross/bin` 目录下。

## 📦 目录结构
```
git-account-manager-go/
├── internal/
│   ├── gitops/    # Git 操作封装
│   └── storage/   # 配置文件存取
├── main.go        # UI 主程序
├── go.mod         # 依赖定义
└── README.md      # 说明文档
```
