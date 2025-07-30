# Contributing to TechIRCd

We love your input! We want to make contributing to TechIRCd as easy and transparent as possible.

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and single-purpose

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## Bug Reports

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening)

## Feature Requests

We welcome feature requests! Please provide:

- Clear description of the feature
- Use case or motivation
- Possible implementation approach (if you have ideas)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
