# 贡献指南

感谢您对VANTUN项目的兴趣！我们欢迎各种形式的贡献，包括但不限于代码、文档、测试和功能建议。

## 如何贡献

### 报告问题

如果您发现了bug或有功能建议，请在GitHub上创建一个issue。

在提交issue之前，请确保：
- 您使用的是最新版本的代码
- 您已经查阅了现有的issue，确认没有重复的报告
- 您提供了详细的复现步骤和环境信息

### 代码贡献

我们欢迎代码贡献！在开始之前，请遵循以下步骤：

1. Fork项目仓库到您的GitHub账户
2. 克隆您的fork到本地开发环境
3. 创建您的特性分支（建议使用`feature/your-feature-name`命名）
4. 进行开发并提交您的修改
5. 推送到您的分支
6. 创建一个Pull Request到VANTUN主仓库

```bash
git clone https://github.com/your-username/vantun.git
cd vantun
git checkout -b feature/your-feature-name
# 进行您的修改
# Make your changes
git commit -m "Add some feature"
git push origin feature/your-feature-name
```

### 代码规范

为了保持代码的一致性，请遵循以下规范：
- 使用Go语言的最佳实践和惯用法
- 编写清晰、有意义的提交信息，遵循[Conventional Commits](https://www.conventionalcommits.org/zh-hans/v1.0.0/)规范
- 为公共函数和类型添加注释，说明其用途、参数和返回值
- 确保所有测试通过，并为新功能添加相应的测试用例

在提交代码前，请运行以下命令确保代码质量：

```bash
go fmt ./...              # 格式化代码
go vet ./...              # 静态分析代码
go test ./...             # 运行测试
```

### 开发环境设置

要设置开发环境，请确保您已安装：
- Go 1.21 or higher
- Git

获取项目依赖：

```bash
go mod tidy
```

### 测试

在提交代码之前，请确保所有测试都通过：

```bash
go test -v ./...
```

对于新功能，请添加相应的测试用例，确保代码质量和功能正确性。

VANTUN使用多种类型的测试：
- 单元测试：测试单个函数或组件
- 集成测试：测试多个组件协同工作
- 压力测试：验证系统在高负载下的表现

### 文档

如果您的更改影响了用户使用方式，请相应更新文档。

包括但不限于：
- README.md
- DEMOGUIDE.md
- CHANGELOG.md
- 代码注释

---

感谢您的贡献！您的努力将帮助VANTUN变得更好。