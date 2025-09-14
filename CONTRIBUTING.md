# 贡献指南 / Contribution Guidelines

[//]: # (感谢您对VANTUN项目的兴趣！我们欢迎各种形式的贡献。)
[//]: # (Thank you for your interest in the VANTUN project! We welcome contributions of all kinds.)

## 如何贡献 / How to Contribute

### 报告问题 / Reporting Issues

[//]: # (如果您发现了bug或有功能建议，请在GitHub上创建一个issue。)
[//]: # (If you find a bug or have a feature suggestion, please create an issue on GitHub.)

在提交issue之前，请确保：
<br>Before submitting an issue, please ensure:

- [//]: # (您使用的是最新版本的代码)
  [//]: # (You are using the latest version of the code)
- [//]: # (您已经查阅了现有的issue，确认没有重复的报告)
  [//]: # (You have checked existing issues to confirm there is no duplicate report)
- [//]: # (您提供了详细的复现步骤和环境信息)
  [//]: # (You have provided detailed reproduction steps and environment information)

### 代码贡献 / Code Contributions

[//]: # (我们欢迎代码贡献！在开始之前，请遵循以下步骤：)
[//]: # (We welcome code contributions! Before you start, please follow these steps:)

1. [//]: # (Fork项目仓库)
   [//]: # (Fork the project repository)
2. [//]: # (创建您的特性分支)
   [//]: # (Create your feature branch)
3. [//]: # (提交您的修改)
   [//]: # (Commit your changes)
4. [//]: # (推送到您的分支)
   [//]: # (Push to your branch)
5. [//]: # (创建一个Pull Request)
   [//]: # (Create a Pull Request)

```bash
git clone https://github.com/your-username/vantun.git
cd vantun
git checkout -b feature/your-feature-name
# Make your changes
git commit -m "Add some feature"
git push origin feature/your-feature-name
```

### 代码规范 / Code Standards

[//]: # (为了保持代码的一致性，请遵循以下规范：)
[//]: # (To maintain code consistency, please follow these guidelines:)

- [//]: # (使用Go语言的最佳实践)
  [//]: # (Use Go language best practices)
- [//]: # (编写清晰、有意义的提交信息)
  [//]: # (Write clear and meaningful commit messages)
- [//]: # (为公共函数和类型添加注释)
  [//]: # (Add comments for public functions and types)
- [//]: # (确保所有测试通过)
  [//]: # (Ensure all tests pass)

```bash
go fmt ./...
go vet ./...
go test ./...
```

### 开发环境设置 / Development Environment Setup

[//]: # (要设置开发环境，请确保您已安装：)
[//]: # (To set up the development environment, ensure you have installed:)

- Go 1.21 or higher
- [//]: # (必要的依赖项会通过go.mod自动管理)
  [//]: # (Necessary dependencies will be automatically managed through go.mod)

```bash
go mod tidy
```

### 测试 / Testing

[//]: # (在提交代码之前，请确保所有测试都通过：)
[//]: # (Before submitting code, please ensure all tests pass:)

```bash
go test -v ./...
```

[//]: # (对于新功能，请添加相应的测试用例。)
[//]: # (For new features, please add corresponding test cases.)

### 文档 / Documentation

[//]: # (如果您的更改影响了用户使用方式，请相应更新文档。)
[//]: # (If your changes affect how users use the project, please update the documentation accordingly.)

---

[//]: # (感谢您的贡献！)
[//]: # (Thank you for your contribution!)