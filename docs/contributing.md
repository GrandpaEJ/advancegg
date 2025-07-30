# Contributing to AdvanceGG

We welcome contributions to AdvanceGG! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your feature or bug fix
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Setting Up Your Development Environment

```bash
# Clone your fork
git clone https://github.com/GrandpaEJ/advancegg.git
cd advancegg

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run examples to verify everything works
cd examples
go run circle.go
```

## Code Style

### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Use Go's built-in error handling patterns

### Example:

```go
// NewContext creates a new drawing context with the specified dimensions.
// The context is initialized with a white background and black foreground color.
func NewContext(width, height int) *Context {
    // Implementation here
}
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./pkg/advancegg
```

### Writing Tests

- Write unit tests for new functions
- Include edge cases in your tests
- Use table-driven tests when appropriate
- Test both success and error cases

Example test:

```go
func TestNewContext(t *testing.T) {
    tests := []struct {
        name   string
        width  int
        height int
        want   bool
    }{
        {"valid dimensions", 100, 100, true},
        {"zero width", 0, 100, false},
        {"zero height", 100, 0, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dc := NewContext(tt.width, tt.height)
            got := dc != nil
            if got != tt.want {
                t.Errorf("NewContext() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use Go's standard documentation format
- Include examples in documentation when helpful

### User Documentation

- Update relevant documentation in the `docs/` directory
- Add examples for new features
- Update the API reference for new functions

## Submitting Changes

### Pull Request Process

1. Ensure your code follows the project's style guidelines
2. Add or update tests for your changes
3. Update documentation as needed
4. Ensure all tests pass
5. Create a clear pull request description

### Pull Request Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] New tests added for new functionality
- [ ] Manual testing completed

## Documentation
- [ ] Code comments updated
- [ ] User documentation updated
- [ ] API reference updated
```

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

- Go version
- Operating system
- Minimal code example that reproduces the issue
- Expected behavior
- Actual behavior
- Any error messages

### Feature Requests

When requesting features, please include:

- Use case description
- Proposed API (if applicable)
- Examples of how the feature would be used
- Any alternative solutions considered

## Code of Conduct

Please be respectful and constructive in all interactions. We want to maintain a welcoming environment for all contributors.

## Questions?

If you have questions about contributing, feel free to:

- Open an issue for discussion
- Start a discussion in the GitHub Discussions section
- Reach out to the maintainers

Thank you for contributing to AdvanceGG!
