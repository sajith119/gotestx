# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [v1.0.0] - 2025-09-20

### Added
- Initial extraction of **Entiqon CLI** as a standalone module (`github.com/entiqon/cli`).
- Introduced the first Go-native binary:
    - **GoTestX** — Go Test eXtended tool with coverage support.
- Integrated existing CLI helpers (Git, release, coverage, Docker) as Bash utilities, with migration to Go planned in future versions.

### Notes
- Consolidates CLI history from the `entiqon` monorepo.
- Establishes baseline for future migration of Bash tools to Go-native binaries.
