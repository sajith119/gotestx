package gotestx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	Version     = "1.0.0"
	ToolName    = "GoTestX"
	Author      = "Entiqon Project Team"
	Description = "Go Test eXtended tool with coverage support"
)

// commandRunner is used to create exec.Cmd.
// In production, it defaults to exec.Command, but tests can override it.
var commandRunner = exec.Command

// getGOOS is used to fetch the current runtime.GOOS.
// In production, it defaults to runtime.GOOS, but tests can override it.
var getGOOS = func() string { return runtime.GOOS }

func usage(w io.Writer) {
	_, _ = fmt.Fprintf(w, `%s v%s
%s
Author: %s

Usage: %s [options] [packages]

Options:
  -c, --with-coverage   Run tests with coverage report generation (coverage.out)
  -o, --open-coverage   Open coverage report in browser (macOS only, implies -c)
  -q, --quiet           Suppress info messages (only errors and test output shown)
  -C, --clean           Suppress 'no test files' lines for cleaner output
  -h, --help            Show this help
  -v, --version         Show version info
`, ToolName, Version, Description, Author, ToolName)
}

func versionInfo(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s\n\n%s\n%s\nAuthor: %s\nVersion: %s\nProcessor: %s (%s)\n",
		Version, ToolName, Description, Author, Version, runtime.GOARCH, runtime.GOOS)
}

// Run executes gotestx with given args, stdout/stderr, returns exit code
func Run(args []string, stdout, stderr io.Writer) int {
	var withCoverage, openCoverage, quiet, clean bool
	var packages []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			usage(stdout)
			return 0
		case arg == "-v" || arg == "--version":
			versionInfo(stdout)
			return 0
		case arg == "-c" || arg == "--with-coverage":
			withCoverage = true
		case arg == "-o" || arg == "--open-coverage":
			openCoverage = true
		case arg == "-q" || arg == "--quiet":
			quiet = true
		case arg == "-C" || arg == "--clean":
			clean = true
		case strings.HasPrefix(arg, "-"):
			// Handle combined short flags (-cqC)
			flags := arg[1:]
			for _, f := range flags {
				switch f {
				case 'c':
					withCoverage = true
				case 'o':
					openCoverage = true
				case 'q':
					quiet = true
				case 'C':
					clean = true
				default:
					_, _ = fmt.Fprintf(stderr, "Error: Unknown short option: -%c\n", f)
					usage(stderr)
					return 2
				}
			}
		default:
			packages = append(packages, arg)
		}
	}

	// Default ./...
	if len(packages) == 0 {
		packages = []string{"./..."}
	}

	// Validate
	for i, pkg := range packages {
		if strings.Contains(pkg, "...") {
			continue
		}
		st, err := os.Stat(pkg)
		if err != nil || !st.IsDir() {
			_, _ = fmt.Fprintf(stderr, "Error: Package path '%s' does not exist.\n", pkg)
			return 1
		}
		matches, _ := filepath.Glob(filepath.Join(pkg, "*.go"))
		if len(matches) == 0 {
			subMatches, _ := filepath.Glob(filepath.Join(pkg, "**", "*.go"))
			if len(subMatches) > 0 {
				if !quiet {
					_, _ = fmt.Fprintf(stdout, "Info: No Go files in '%s', using subpackages instead (%s/...)\n", pkg, pkg)
				}
				packages[i] = pkg + "/..."
			} else {
				_, _ = fmt.Fprintf(stderr, "Error: No Go files found in '%s'.\n", pkg)
				return 1
			}
		}
	}

	if openCoverage && !withCoverage {
		withCoverage = true
	}

	if openCoverage && getGOOS() != "darwin" {
		_, _ = fmt.Fprintln(stderr, "Error: --open-coverage is only supported on macOS.")
		return 1
	}

	// Run tests
	var cmd *exec.Cmd
	if withCoverage {
		if !quiet {
			_, _ = fmt.Fprintf(stdout, "Running tests with coverage across: %s\n", strings.Join(packages, " "))
		}
		args := append([]string{"test", "-coverprofile=coverage.out", "-covermode=atomic"}, packages...)
		cmd = commandRunner("go", args...)
	} else {
		if !quiet {
			_, _ = fmt.Fprintf(stdout, "Running tests normally across: %s\n", strings.Join(packages, " "))
		}
		args := append([]string{"test"}, packages...)
		cmd = commandRunner("go", args...)
	}

	var buf bytes.Buffer
	if clean {
		cmd.Stdout = &buf
		cmd.Stderr = &buf
	} else if quiet {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		cmd.Stdout = stdout
		cmd.Stderr = stderr
	}

	if err := cmd.Run(); err != nil {
		if !quiet {
			_, _ = fmt.Fprintf(stderr, "Error: go test failed: %v\n", err)
		}
		return 1
	}

	if clean {
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "[no test files]") {
				continue
			}
			_, _ = fmt.Fprintln(stdout, line)
		}
	}

	if withCoverage {
		if !quiet {
			_, _ = fmt.Fprintln(stdout, "Coverage report saved as coverage.out")
			_, _ = fmt.Fprintln(stdout, "Run 'go tool cover -html=coverage.out' to view it")
		}
		if openCoverage {
			if !quiet {
				_, _ = fmt.Fprintln(stdout, "Opening coverage report in browser...")
			}
			openCmd := commandRunner("go", "tool", "cover", "-html=coverage.out")
			openCmd.Stdout = stdout
			openCmd.Stderr = stderr
			if err := openCmd.Run(); err != nil {
				_, _ = fmt.Fprintf(stderr, "Error: failed to open coverage report: %v\n", err)
				return 1
			}
		}
	}

	return 0
}
