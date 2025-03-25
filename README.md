# Todo TUI

A powerful, terminal-based todo application built with Go. Manage your tasks efficiently without leaving the command line.

![Todo TUI Screenshot](docs/images/screenshot.png)

## Features

- 📋 Simple and intuitive terminal UI
- 🏷️ Tag support for organizing related tasks
- 🔄 Multiple status views (Open, Doing, Done, Archived)
- 🚩 Priority levels (Low, Medium, High)
- 📅 Due date support
- 🔍 Filtering and searching capabilities
- ⌨️ Keyboard-driven interfaceo

## Installation

### macOS and Linux

```bash
# Download the appropriate version for your system
curl -L -o tui-todo.tar.gz https://github.com/martijnspitter/tui-todo/releases/latest/download/tui-todo-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m).tar.gz

# Extract
tar -xzf tui-todo.tar.gz

# Run installer
./install.sh

# Clean up
rm tui-todo.tar.gz tui-todo install.sh
```

### Windows

1. Download the [latest Windows release](https://github.com/martijnspitter/tui-todo/releases/latest/download/tui-todo-windows-amd64.zip)
2. Extract the ZIP file
3. Right-click on `install.ps1` and select "Run with PowerShell"
4. Open a new PowerShell window and run `tui-todo`

## Installation

### macOS and Linux

```bash
# Download the appropriate version for your system
curl -L -o tui-todo.tar.gz https://github.com/martijnspitter/tui-todo/releases/latest/download/tui-todo-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m).tar.gz

# Extract
tar -xzf tui-todo.tar.gz

# Run installer
./install.sh

# Clean up
rm tui-todo.tar.gz todo install.sh

# Now you can use the app from anywhere
todo
```

### Windows

1. Download the [latest Windows release](https://github.com/martijnspitter/tui-todo/releases/latest/download/tui-todo-windows-amd64.zip)
2. Extract the ZIP file
3. Right-click on `install.ps1` and select "Run with PowerShell"
4. Open a new PowerShell window and run `todo`

## Usage

| Key               | Action                   |
| ----------------- | ------------------------ |
| Tab/Right         | Next field/item          |
| Shift+Tab/Left    | Previous field/item      |
| Enter             | Select/Confirm           |
| Esc               | Go back/Cancel           |
| Ctrl+C            | Quit application         |
| Ctrl+N            | Create new todo          |
| Ctrl+E            | Edit selected todo       |
| Ctrl+D            | Delete selected todo     |
| Ctrl+Space/Ctrl+S | Advance todo status      |
| 1                 | Switch to Open todos     |
| 2                 | Switch to Doing todos    |
| 3                 | Switch to Done todos     |
| 4                 | Switch to Archived todos |
| 5                 | Switch to New todo pane  |

## Collaboration

We welcome contributions to make Todo TUI even better! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) guide for details on:

- Our development workflow
- Commit message guidelines
- Pull request process
- Code style standards

For a quick start:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using conventional commits
4. Push to your branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

We look forward to your contributions!

### Development Workflow

1. **Make your changes** and ensure they adhere to the project's style and standards.

2. **Run tests** to make sure everything works:

   ```bash
   go test -v ./...
   ```

3. **Commit your changes** with a descriptive message:

   ```bash
   git commit -m "Add feature: your feature description"
   ```

4. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request** by navigating to your fork on GitHub and clicking "New Pull Request".

### Pull Request Guidelines

- Keep PRs focused on a single feature or bugfix
- Include a clear description of the changes and their purpose
- Ensure all tests pass
- Update documentation if necessary

### Code Style

- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Include comments for complex logic
- Format code with `go fmt` before submitting

### Issue Reporting

Found a bug or have a feature request? Please [create an issue](https://github.com/martijnspitter/tui-todo/issues/new) with:

- A clear description of the bug or feature
- Steps to reproduce (for bugs)
- Expected vs. actual behavior (for bugs)
- Any relevant screenshots or logs

### Building from Source

```bash
# Clone the repository
git clone https://github.com/martijnspitter/tui-todo.git
cd tui-todo

# Build the application
go build -o todo ./cmd/tui-todo

# Run the application
./todo
```

### Running Tests

```bash
go test -v ./...
```

### Development Setup

To set up a development environment:

```bash
# Clone the repository
git clone https://github.com/martijnspitter/tui-todo.git
cd tui-todo

# Install dependencies
go mod download

# Run the application in development mode
go run ./cmd/tui-todo/main.go
```

## Configuration

Todo TUI stores its data in:

- Linux: `~/.local/share/tui-todo/todo.sql`
- macOS: `~/Library/Application Support/tui-todo/todo.sql`
- Windows: `%APPDATA%\tui-todo\todo.sql`

## Screenshots

![Task List View](docs/images/task-list.png)
_Main task view with different status tabs_

![Task Edit View](docs/images/task-edit.png)
_Editing a task with tags and due date_

## License

MIT License - See [LICENSE](LICENSE) for details.
