# Gondo — Terminal Todo Manager

Gondo is a fast, keyboard-driven terminal todo manager built in Go.
It focuses on hierarchical task management, allowing you to organize tasks in a tree structure with a clean and responsive TUI.

---

## ✨ Features

* 📁 **Hierarchical Tasks**

  * Create nested tasks (tasks inside tasks)
  * Collapse / expand composite tasks

* ⚡ **Keyboard-Driven Workflow**

  * No mouse needed
  * Fast navigation and editing

* 🎯 **Task Operations**

  * Add, insert, push, delete, rename
  * Toggle completion state

* 🧠 **Smart Rendering**

  * Scrollable view (vertical + horizontal)
  * Styled output using a TUI framework

* 💾 **Persistence**

  * Tasks are saved locally
  * Supports global and project-specific lists

---

## 🚀 Installation

```bash
git clone https://github.com/OmarAbuAbdullah69/gondo
cd gondo
go build -o gondo
```

---

## ▶️ Usage

### Run with default local list

```bash
./gondo
```

### Use global todo list

```bash
./gondo -g
```

### Show help

```bash
./gondo -h
```

### Use custom directory

```bash
./gondo /path/to/project
```

---

## 🎮 Key Bindings

| Key     | Action                    |
| ------- | ------------------------- |
| ↑ / ↓   | Navigate tasks            |
| ← / →   | Horizontal scroll         |
| `A`     | Add main task             |
| `a`     | Add task next to selected |
| `p`     | Push task under selected  |
| `i`     | Insert after selected     |
| `I`     | Insert before selected    |
| `space` | Toggle completion         |
| `f`     | Fold / unfold task        |
| `d`     | Delete task               |
| `r`     | Rename task               |
| `R`     | Rename list               |
| `esc`   | Cancel input / exit       |
| `q`     | Quit                      |

---

## 🧱 Project Structure

```
gondo/
├── main.go          # Entry point
├── internal/
│   ├── core/        # Application logic (state + behavior)
│   ├── task/        # Task system (interfaces + implementations)
│   └── tui/         # Terminal UI layer
```

---

## 🧩 Architecture Overview

Gondo is split into three main parts:

### 1. Core

Handles:

* Task manipulation
* Navigation logic
* State management

### 2. Task System

Defines:

* Task interface (`Tasker`)
* Simple tasks
* Composite (nested) tasks

### 3. TUI

Responsible for:

* Rendering UI
* Handling keyboard input
* Display updates

---

## 📌 Design Goals

* Keep everything **keyboard-first**
* Maintain **simple but powerful task hierarchy**
* Build a **clean separation between logic and UI**
* Stay **lightweight and fast**

---

## ⚠️ Known Limitations

* Heavy use of global state (planned refactor)
* No undo/redo system yet
* No search/filtering
* Performance may degrade with very large task lists

---

## 🔮 Future Improvements

* Undo / redo support
* Task search and filtering
* Better state management (remove globals)
* Improved scrolling and rendering performance
* Configurable key bindings

---

## 🤝 Contributing

Contributions are welcome!

You can help by:

* Refactoring core architecture
* Adding features (search, undo, etc.)
* Improving UI/UX
* Fixing bugs

---

## 📜 License

MIT License

---

## 💡 Inspiration

Gondo is inspired by the need for a simple, fast, and structured todo system that works entirely inside the terminal.

