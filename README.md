# todo - CLI Task Manager in Go

A practical, production-quality command-line task manager written in Go with persistent local storage, intuitive subcommands via Cobra, and zero external database requirements.

---

## Features

* **Complete Task Lifecycle**: Add, list, show, edit, complete (`done`), revert (`undo`), delete, and clear tasks.
* **Persistent Local Storage**: Tasks are automatically persisted in standard JSON format across terminal sessions.
* **Predictable ID Management**: Incremental ID assignment that tracks `max(ID) + 1` to prevent collisions after task deletions.
* **Atomic File Writes**: Prevents file corruption during unexpected crashes or interruptions by using temporary files and atomic file replacement.
* **OS-Appropriate Default Directories**:
  * **Windows**: `%APPDATA%\todo\tasks.json`
  * **Linux / macOS**: `~/.todo/tasks.json`
* **Configurable Location**: Override the storage location via the `TODO_DATA_DIR` environment variable (ideal for testing, automation, or custom profiles).
* **Terminal Table Output**: Clean, tab-aligned task listing formatted using the standard library `text/tabwriter`.
* **Flexible Filtering & Sorting**: Filter tasks by `--all`, `--pending`, or `--completed`, and sort by `id`, `created`, or `updated`.
* **Safe CLI Interactions**: Interactive confirmation prompts for destructive operations (`delete`, `clear`) with a `--force` (`-f`) flag for automation.
* **Helpful Error Messages & Proper Exit Codes**: Clean error output to `stderr` with non-zero exit codes without spilling raw stack traces.

---

## Installation & Building

### Prerequisites

* Go 1.20 or newer installed ([Download Go](https://go.dev/dl/)).

### Building from Source

Clone or navigate into the repository directory and build the binary:

```bash
# Build executable in current directory
go build -o todo .

# On Windows:
go build -o todo.exe .
```

### Adding to PATH

To use `todo` from anywhere in your terminal:

* **Linux / macOS**:
  ```bash
  go install .
  # Ensure $GOPATH/bin or $HOME/go/bin is in your $PATH
  ```
  Or move the binary:
  ```bash
  sudo mv todo /usr/local/bin/
  ```

* **Windows (PowerShell)**:
  ```powershell
  go install .
  # Ensure $env:USERPROFILE\go\bin is in your System PATH
  ```

---

## Usage & Commands

### 1. Adding Tasks (`todo add`)

Creates a new task. The title must not be empty.

```bash
todo add "Learn Go concurrency"
todo add "Write unit tests"
```

Output:
```text
Task added successfully.

ID:    1
Title: Learn Go concurrency
```

### 2. Listing Tasks (`todo list`)

Displays tasks in a formatted terminal table:

```bash
todo list
```

Output:
```text
ID   STATUS   TASK
1    [ ]      Learn Go concurrency
2    [ ]      Write unit tests
```

#### Filtering & Sorting Options:

* `todo list --all` (`-a`): Show all tasks (default).
* `todo list --pending` (`-p`): Show only pending tasks.
* `todo list --completed` (`-c`): Show only completed tasks.
* `todo list --sort created` (`-s created`): Sort by creation timestamp.
* `todo list --sort updated` (`-s updated`): Sort by last updated timestamp.
* `todo list --sort id` (`-s id`): Sort by task ID (default).

### 3. Showing Details (`todo show`)

Displays full details, including status, creation time, and last updated time:

```bash
todo show 1
```

Output:
```text
Task #1

Title:      Learn Go concurrency
Status:     Pending
Created:    2026-09-10 13:50:00
Updated:    2026-09-10 13:50:00
```

### 4. Completing a Task (`todo done`)

Marks the task as completed:

```bash
todo done 1
```

Output:
```text
Task #1 marked as completed.
```

If the task is already completed, it informs you without modifying the update timestamp.

### 5. Reverting a Task (`todo undo`)

Marks a completed task back to pending:

```bash
todo undo 1
```

Output:
```text
Task #1 marked as pending.
```

### 6. Editing a Task Title (`todo edit`)

Modifies the title of an existing task and updates its timestamp:

```bash
todo edit 1 "Master Go concurrency and channels"
```

Output:
```text
Task #1 updated successfully.
```

### 7. Deleting a Task (`todo delete`)

Deletes a task by ID. Prompts for confirmation in interactive mode:

```bash
todo delete 1
```

Prompt:
```text
Delete task #1 "Master Go concurrency and channels"? [y/N]: y
Task #1 deleted successfully.
```

To skip the confirmation prompt (e.g., in CI or automation):

```bash
todo delete 1 --force
# or
todo delete 1 -f
```

### 8. Clearing All Tasks (`todo clear`)

Removes all tasks from storage:

```bash
todo clear
```

Prompt:
```text
Are you sure you want to delete all 5 tasks? [y/N]: y
All tasks cleared successfully.
```

Bypass confirmation:

```bash
todo clear --force
# or
todo clear -f
```

### 9. Version and Help

```bash
todo version
todo --version
todo --help
todo add --help
```

---

## Storage & Configuration

### Storage Path

Tasks are stored in a standard JSON file named `tasks.json`.

* **Windows**:
  `%APPDATA%\todo\tasks.json` (typically `C:\Users\<User>\AppData\Roaming\todo\tasks.json`)
* **Linux / macOS**:
  `$HOME/.todo/tasks.json`

The storage directory is created automatically with `0755` permissions if it does not exist.

### Environment Variable (`TODO_DATA_DIR`)

You can override the storage directory by setting the `TODO_DATA_DIR` environment variable:

```bash
# Linux / macOS
export TODO_DATA_DIR="/tmp/custom-todos"
todo list

# Windows (PowerShell)
$env:TODO_DATA_DIR = "C:\custom-todos"
todo list
```

---

## Project Structure

```text
todo/
├── cmd/
│   ├── root.go       # Root Cobra command, version, and global I/O bindings
│   ├── add.go        # 'todo add' command implementation
│   ├── list.go       # 'todo list' with tabwriter formatting & flags
│   ├── show.go       # 'todo show' command implementation
│   ├── edit.go       # 'todo edit' command implementation
│   ├── done.go       # 'todo done' command implementation
│   ├── undo.go       # 'todo undo' command implementation
│   ├── delete.go     # 'todo delete' with confirmation prompt & --force
│   ├── clear.go      # 'todo clear' with confirmation prompt & --force
│   ├── version.go    # 'todo version' command
│   └── cmd_test.go   # Integration unit tests for CLI commands
│
├── internal/
│   ├── task/
│   │   ├── task.go       # Task model, ID calculation, filtering, sorting
│   │   └── task_test.go  # Unit tests for domain logic
│   └── storage/
│       ├── storage.go    # Storage interface, custom errors, path resolution
│       ├── json.go       # JSON storage implementation with atomic write
│       └── json_test.go  # Unit tests for persistent storage
│
├── main.go           # CLI application entry point
├── go.mod            # Go module definitions
├── go.sum            # Checksum lockfile
└── README.md         # Documentation
```

---

## Testing & Quality

Run all unit tests across the entire codebase:

```bash
go test -v ./...
```

Run code formatting and static analysis:

```bash
gofmt -w .
go vet ./...
```

---

## License

MIT License. Feel free to modify and adapt for personal or commercial use.
