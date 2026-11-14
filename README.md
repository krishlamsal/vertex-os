# vertex-os

![version](https://img.shields.io/badge/version-1.0.0-blue)
![license](https://img.shields.io/github/license/krishlamsal/vertex-os)
![Go version](https://img.shields.io/badge/Go%20version-1.26.5-green)

`vertex-os` is **not** an operating system but an executable that runs on top of a terminal OS as a user interface. It is designed to be used with minimal inputs (e.g. a 3x4 keypad).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [License](#license)

## Installation

### 1. Download the Binary

Download the latest release binary and place it in your system's `PATH` (e.g., `/usr/local/bin` on Linux/macOS or a dedicated tools folder on Windows).

## Usage

### Apps

`vertex-os` comes with:

- [x] Launcher



#### Utilities

  - [x] Calculator
  - [ ] Notes
  - [ ] Calendar
  - [ ] Camera

#### Time and Weather

  - [x] Clock
  - [x] Stopwatch
  - [ ] Timer
  - [ ] Weather

#### Communication

  - [ ] Browser
  - [ ] Phone
  - [ ] Phonebook
  - [ ] Messages

#### Files

  - [ ] Image Viewer
  - [ ] Text Editor
  - [ ] Audio Player
  - [ ] Files

#### Advanced

  - [ ] Help
  - [ ] Settings
  - [ ] System Monitor
  - [ ] Terminal

### Keybinds

`vertex-os` is primarily designed for this 3x4 layout:

```

┌───────┐
│ 1 2 3 │
│ 4 5 6 │
│ 7 8 9 │
│ * 0 # │
└───────┘

```

#### While navigating

- **2**, **4**, **6** and **8** function as arrow keys:

```

    ↑
    2
← 4   6 → 
    8
    ↓
    
```

- **7** and **9** are for app drawer paging

```

prev. page ←← 7   9 →→ next page

```

- `5` functions as an enter/OK key
- `*` functions as a home/app drawer button

#### While typing

`vertex-os` uses an **on-screen keyboard** to type characters. Navigate around the keyboard using **2**, **4**, **6** and **8** and use `5` to select the key.

Whilst typing numbers (e.g. a PIN or phone number), you can simply use the number keys to input and use `#` to enter

## License

[MIT © Krish Lamsal.](LICENSE)
