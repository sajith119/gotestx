# GoTestX

> Go Test eXtended tool with coverage support

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev) <a href="https://github.com/entiqon/gotestx/releases"><img src="https://img.shields.io/github/v/release/entiqon/gotestx" alt="Latest Release" /></a>
[![Build Status](https://github.com/entiqon/gotestx/actions/workflows/ci.yml/badge.svg)](https://github.com/entiqon/gotestx/actions)
[![Codecov](https://codecov.io/gh/entiqon/gotestx/branch/main/graph/badge.svg)](https://codecov.io/gh/entiqon/gotestx)

GoTestX extends the standard [`go test`](https://pkg.go.dev/cmd/go#hdr-Test_packages) command with a simpler, more versatile interface.
It adds optional coverage reporting, quiet mode, and clean output filtering — while remaining fully compatible with `go test`.

---

## ✨ Features

* **Coverage mode** (`-c`): generates `coverage.out` with `-covermode=atomic`.
* **Open coverage** (`-o`): opens the HTML coverage report in a browser (macOS only).
* **Quiet mode** (`-q`): suppresses info logs, only shows test results and errors.
* **Clean mode** (`-C`): removes `? … [no test files]` lines for cleaner output.
* **Flag combinations**: short flags can be combined (e.g. `-cq`, `-coq`, `-cC`).
* **Smart package detection**:

    * Expands `./pkg` → `./pkg/...` if root has no Go files but subpackages do.
    * Reports errors if a path doesn’t exist or has no Go files.

---

## 🚀 Installation

Install directly via GitHub:

```bash
go install github.com/entiqon/gotestx@latest
```

Check installation:

```bash
gotestx -v
```

---

## 📦 Usage

```bash
gotestx [options] [packages]
```

Options:

```
  -c, --with-coverage   Run tests with coverage report generation (coverage.out)
  -o, --open-coverage   Open coverage report in browser (macOS only, implies -c)
  -q, --quiet           Suppress info messages (only errors and test output shown)
  -C, --clean           Suppress 'no test files' lines for cleaner output
  -h, --help            Show this help
  -v, --version         Show version info
```

---

## 🧪 Examples

Run tests for all packages:

```bash
gotestx
```

Run tests with coverage:

```bash
gotestx -c ./...
```

Run quietly with coverage:

```bash
gotestx -cq ./...
```

Run with coverage and open report in browser (macOS):

```bash
gotestx -o ./...
```

Run with clean output (no `[no test files]` lines):

```bash
gotestx -C ./...
```

Combined:

```bash
gotestx -cCq ./...
```

---

## 🛠 Development

Clone the repository:

```bash
git clone https://github.com/entiqon/gotestx.git
cd gotestx
```

Build:

```bash
go build -o gotestx .
```

Test:

```bash
go test ./internal/... -v
```

---

## 📄 License

Part of the [Entiqon Project](https://github.com/entiqon).
Licensed under the MIT License.
