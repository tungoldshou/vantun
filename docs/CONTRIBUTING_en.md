# Contribution Guidelines

Thank you for your interest in the VANTUN project! We welcome contributions of all kinds.

## How to Contribute

### Reporting Issues

If you find a bug or have a feature suggestion, please create an issue on GitHub.

Before submitting an issue, please ensure:

- You are using the latest version of the code
- You have checked existing issues to confirm there is no duplicate report
- You have provided detailed reproduction steps and environment information

### Code Contributions

We welcome code contributions! Before you start, please follow these steps:

1. Fork the project repository
2. Create your feature branch
3. Commit your changes
4. Push to your branch
5. Create a Pull Request

```bash
git clone https://github.com/your-username/vantun.git
cd vantun
git checkout -b feature/your-feature-name
# Make your changes
git commit -m "Add some feature"
git push origin feature/your-feature-name
```

### Code Standards

To maintain code consistency, please follow these guidelines:

- Use Go language best practices
- Write clear and meaningful commit messages
- Add comments for public functions and types
- Ensure all tests pass

```bash
go fmt ./...
go vet ./...
go test ./...
```

### Development Environment Setup

To set up the development environment, ensure you have installed:

- Go 1.21 or higher
- Necessary dependencies will be automatically managed through go.mod

```bash
go mod tidy
```

### Testing

Before submitting code, please ensure all tests pass:

```bash
go test -v ./...
```

For new features, please add corresponding test cases.

### Documentation

If your changes affect how users use the project, please update the documentation accordingly.

---

Thank you for your contribution!