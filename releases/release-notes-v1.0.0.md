# Release Notes — v1.0.0

## 🚀 Overview
This release establishes **Entiqon CLI** as a standalone module, extracted from the main Entiqon monorepo.  
It provides a unified toolkit to support development, testing, release automation, and operational workflows.

---

## ✨ Features

### Go-native binaries
- **GoTestX** — Go Test eXtended tool with coverage support:
  - `-c, --with-coverage`: generate coverage report (`coverage.out`).
  - `-o, --open-coverage`: open coverage report in browser (macOS only).
  - `-q, --quiet`: suppress info messages.
  - `-C, --clean`: filter out `? … [no test files]` lines.
  - Supports combining short flags (e.g. `-cq`, `-coq`, `-cC`).
  - Smart package detection (`./pkg` → `./pkg/...` when only subpackages contain Go files).
  - Deterministic test suite with mocked runner for full coverage.
  - Supersedes the legacy Bash helpers `run-tests.sh` and `open-coverage.sh`.

### Bash-based helpers (planned migration to Go)
- **Git & Release automation**:
  - `gcpr` — create GitHub Pull Requests.
  - `gce` — extract commit history.
  - `gcr` — generate release notes.
  - `gct` — create and sign tags.
  - `gsux` — stash/unstash workflow utility.
  - `gcch` — changelog helper.
- **Docker**:
  - `ddc` — deploy Docker container.

---

## 🛠 CI/CD
- Workflow renamed to **“CLI Build & Test”** for clarity.
- Runs on GitHub Actions with:
  - Go stable setup.
  - Test execution and coverage enforcement.
  - Upload of coverage reports to Codecov.
- Enforces **80% minimum coverage**.

---

## 📚 Documentation
- Added project-level `README.md` for GoTestX with badges.
- Added `CHANGELOG.md` (Keep a Changelog format, Semantic Versioning).
- Release notes prepared for v1.0.0.

---

## 📝 Notes
- This release consolidates the CLI history from the Entiqon monorepo into a dedicated repository.
- Introduces **GoTestX** as the first **Go-native binary**, replacing `run-tests.sh` and `open-coverage.sh` with cross-platform support.
- Other utilities (`gcpr`, `gce`, `gcr`, `gct`, `gsux`, `gcch`, `ddc`) remain **Bash-based**, with migration to Go planned for future releases.
- CI/CD pipeline established under the name **“CLI Build & Test”**, with Codecov integration and enforced coverage thresholds.

