# Contributing to Todo TUI

Thank you for considering contributing to Todo TUI! This document outlines the process for contributing to this project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

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

3. Fill out the PR template with all relevant information.

### 6. Code Review Process

- A maintainer will review your PR
- Address any requested changes
- Once approved, a maintainer will merge your PR

## Testing

- Write tests for new functionality
- Run the test suite before submitting a PR:
  ```bash
  go test -v ./...
  ```

## Documentation

Update relevant documentation when:

- Adding new features
- Changing existing functionality
- Fixing bugs that affect user experience

## Release Process

Only project maintainers can create releases. The process is:

1. Changes are accumulated on the `main` branch
2. Maintainers create a `release/vX.Y.Z` branch when ready
3. Final testing and version bumping occurs on this branch
4. Maintainers tag the release and merge back to `main`

## Questions?

If you have any questions, please open an issue or reach out to the maintainers.

Thank you for contributing to Todo TUI!
