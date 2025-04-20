# Contributing to Todo TUI

Thank you for considering contributing to Todo TUI! This document outlines the process for contributing to this project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone. This includes:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

Violations of these guidelines may result in comment moderation, account restriction, or removal from the project. Report unacceptable behavior by contacting the project maintainers.

## Development Setup

### Prerequisites

- Go 1.18 or later
- Terminal with TUI support
- Git

### Local Development

1. Install dependencies:
   ```bash
   go mod download
   ```
2. Run the application:
   ```bash
   go run ./cmd/tui-todo/main.go
   ```
3. For debugging mode with additional logs:
   ```bash
   DEBUG=true go run ./cmd/tui-todo/main.go
   ```

## Project Structure

- `cmd/tui-todo/` - Main application entry point
- `internal/` - Private application code
  - `ui/` - TUI components and views
  - `repository/` - Data storage and management
  - `keys/` - Keyboard mapping definitions
  - `i18n/` - Internationalization support
  - `os-operations/` - OS-specific functionality
  - `logger/` - The logger used in the project
  - `service/` - Contains the app and ui logic
  - `styling/` - All the shared styling
  - `version/` - Version control and update checks
  - `models/` - The models in use by the repository
- `docs/` - Documentation and images
- `pkg/` - Reusable public packages (if applicable)

## Development Workflow

### 1. Fork and Clone

1. Fork the [Todo TUI repository](https://github.com/martijnspitter/tui-todo)
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/tui-todo.git
   cd tui-todo
   ```
3. Add the original repository as an upstream remote:
   ```bash
   git remote add upstream https://github.com/martijnspitter/tui-todo.git
   ```

### 2. Create a Branch

Create a branch for your work:

```bash
git checkout -b branch-type/description
```

Branch naming convention:

- `feature/add-new-feature`
- `bugfix/fix-issue-123`
- `docs/update-installation-guide`
- `refactor/optimize-task-storage`

### 3. Make Your Changes

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write tests for new functionality
- Ensure existing tests pass: `go test ./...`
- Format your code with `go fmt ./...`

### 4. Commit Your Changes

We use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.

Format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style/formatting changes (no code change)
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or modifying tests
- `chore`: Changes to build process or auxiliary tools

Examples:

```
feat(tasks): add recurring task functionality

- Tasks can now be set to recur daily, weekly, or monthly
- Added UI controls for setting recurrence
- Updated database schema

Closes #42
```

```
fix: prevent crash when opening empty task list

The application would crash when opening a project with no tasks.
This commit adds a check to prevent the crash and displays a message
instead.
```

### 5. Push Changes and Create Pull Request

1. Push your branch to your fork:

   ```bash
   git push origin your-branch-name
   ```

2. Go to the [Todo TUI repository](https://github.com/martijnspitter/tui-todo) and create a Pull Request.

3. Fill out the PR template with all relevant information, including:
   - Summary of changes
   - Related issue numbers
   - Screenshots (for UI changes)
   - Any breaking changes
   - Checklist of completed items

## Issue Management

Before creating a new issue:

1. Search existing issues to avoid duplicates
2. Use the appropriate issue template:
   - Bug report: For reporting application problems
   - Feature request: For suggesting new functionality
   - Question: For general inquiries

Issues are triaged weekly by maintainers and labeled appropriately to track status and priority.

When reporting bugs, please include:

- Todo TUI version
- Operating system
- Terminal emulator
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable

## Code Review Process

- A maintainer will review your PR
- Address any requested changes
- Once approved, a maintainer will merge your PR

### Code Review Criteria

PRs are evaluated based on:

- Functionality: Does the code work as expected?
- Tests: Are there sufficient tests covering the changes?
- Style: Does it follow our code conventions?
- Performance: Are there any performance concerns?
- Documentation: Are changes properly documented?
- Compatibility: Does it maintain backward compatibility?

## Testing

- Write tests for new functionality
- Run the test suite before submitting a PR:
  ```bash
  go test -v ./...
  ```
- For UI components, include manual testing steps in your PR description

## Documentation

Update relevant documentation when:

- Adding new features
- Changing existing functionality
- Fixing bugs that affect user experience

Documentation includes:

- Code comments
- README.md updates
- Help text within the application
- Screenshots for UI changes

## Release Process

Only project maintainers can create releases. The process is:

1. Changes are accumulated on the `main` branch
2. Maintainers create a `release/vX.Y.Z` branch when ready
3. Final testing and version bumping occurs on this branch
4. Maintainers tag the release and merge back to `main`

## Questions?

If you have any questions, please:

1. Check the [GitHub Discussions](https://github.com/martijnspitter/tui-todo/discussions) for similar questions
2. Open a new discussion for general inquiries
3. Open an issue for specific problems
4. Reach out to the maintainers for urgent matters

Thank you for contributing to Todo TUI!
