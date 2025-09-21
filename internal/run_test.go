package gotestx

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func init() {
	// In tests, echo the command instead of executing it.
	// This way, stdout contains the simulated command, useful for assertions.
	commandRunner = func(name string, args ...string) *exec.Cmd {
		all := append([]string{name}, args...)
		return exec.Command("echo", strings.Join(all, " "))
	}
}

func run(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	code = Run(args, &outBuf, &errBuf)
	return outBuf.String(), errBuf.String(), code
}

func TestHelp(t *testing.T) {
	out, _, code := run(t, "-h")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage in output, got %q", out)
	}
}

func TestVersion(t *testing.T) {
	out, _, code := run(t, "-v")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, Version) {
		t.Errorf("expected version %s, got %q", Version, out)
	}
}

func TestUnknownFlag(t *testing.T) {
	_, err, code := run(t, "-z")
	if code == 0 {
		t.Errorf("expected non-zero exit code")
	}
	if !strings.Contains(err, "Unknown short option") {
		t.Errorf("expected unknown option error, got %q", err)
	}
}

func TestWithCoverage(t *testing.T) {
	out, _, code := run(t, "-c", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage message in output, got %q", out)
	}
}

func TestClean(t *testing.T) {
	out, _, code := run(t, "-C", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean mode to suppress 'no test files', got %q", out)
	}
}

func TestCleanWithCoverage(t *testing.T) {
	out, _, code := run(t, "-cC", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage message in output, got %q", out)
	}
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean mode to suppress 'no test files', got %q", out)
	}
}

func TestQuietMode(t *testing.T) {
	out, _, code := run(t, "-q", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if strings.Contains(out, "Running") {
		t.Errorf("quiet mode should suppress output, got %q", out)
	}

	out, _, code = run(t, "-q")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if strings.Contains(out, "Running") {
		t.Errorf("quiet mode should suppress output, got %q", out)
	}
}

func TestOpenCoverage(t *testing.T) {
	out, err, code := run(t, "-o", "./")
	if runtime.GOOS != "darwin" {
		// On non-macOS systems, this should error
		if code == 0 {
			t.Errorf("expected non-zero exit code on non-macOS, got %d", code)
		}
		if !strings.Contains(err, "only supported on macOS") {
			t.Errorf("expected macOS-only error, got %q", err)
		}
	} else {
		// On macOS, it should succeed
		if code != 0 {
			t.Errorf("expected exit 0, got %d", code)
		}
		if !strings.Contains(out, "coverage") {
			t.Errorf("expected coverage message in output, got %q", out)
		}
		if !strings.Contains(out, "Opening coverage report") {
			t.Errorf("expected 'Opening coverage report' message, got %q", out)
		}
	}
}

func TestCombinedFlagsQuiet(t *testing.T) {
	out, _, code := run(t, "-cq", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if strings.Contains(out, "Running") || strings.Contains(out, "Opening coverage report") {
		t.Errorf("quiet mode should suppress output, got %q", out)
	}
}

func TestCombinedFlagsWithOpenCoverage(t *testing.T) {
	out, err, code := run(t, "-co", "./")
	if runtime.GOOS != "darwin" {
		if code == 0 {
			t.Errorf("expected non-zero exit code on non-macOS, got %d", code)
		}
		if !strings.Contains(err, "only supported on macOS") {
			t.Errorf("expected macOS-only error, got %q", err)
		}
	} else {
		if code != 0 {
			t.Errorf("expected exit 0, got %d", code)
		}
		if !strings.Contains(out, "Opening coverage report") {
			t.Errorf("expected 'Opening coverage report' message, got %q", out)
		}
	}
}

func TestInvalidPackagePath(t *testing.T) {
	out, errStr, code := run(t, "./does-not-exist")
	if code == 0 {
		t.Errorf("expected non-zero exit code for invalid package path")
	}
	if !strings.Contains(errStr, "does not exist") {
		t.Errorf("expected 'does not exist' error, got %q", errStr)
	}
	if out != "" {
		t.Errorf("expected no stdout, got %q", out)
	}
}

func TestPackageWithSubpackages(t *testing.T) {
	dir := t.TempDir()

	// subdir with a .go file
	sub := filepath.Join(dir, "subpkg")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	src := []byte(`package subpkg; func Foo() {}`)
	if err := os.WriteFile(filepath.Join(sub, "dummy.go"), src, 0644); err != nil {
		t.Fatal(err)
	}

	out, errStr, code := run(t, dir)
	if code != 0 {
		t.Errorf("expected exit 0, got %d (stderr=%q)", code, errStr)
	}
	if !strings.Contains(out, "using subpackages") {
		t.Errorf("expected info about subpackages, got %q", out)
	}
}

func TestPackageWithoutGoFiles(t *testing.T) {
	dir := t.TempDir() // no files inside

	out, errStr, code := run(t, dir)
	if code == 0 {
		t.Errorf("expected non-zero exit code for no go files")
	}
	if !strings.Contains(errStr, "No Go files found") {
		t.Errorf("expected error about no go files, got %q", errStr)
	}
	if out != "" {
		t.Errorf("expected no stdout, got %q", out)
	}
}

func TestOpenCoverageGuardForced(t *testing.T) {
	old := getGOOS
	defer func() { getGOOS = old }()

	// Force non-macOS
	getGOOS = func() string { return "linux" }
	_, errStr, code := run(t, "-o", "./")
	if code == 0 {
		t.Errorf("expected non-zero exit code")
	}
	if !strings.Contains(errStr, "only supported on macOS") {
		t.Errorf("expected macOS-only error, got %q", errStr)
	}

	// Force macOS
	getGOOS = func() string { return "darwin" }
	out, _, code := run(t, "-o", "./")
	if code != 0 {
		t.Errorf("expected exit 0 on darwin, got %d", code)
	}
	if !strings.Contains(out, "Opening coverage report") {
		t.Errorf("expected 'Opening coverage report' message, got %q", out)
	}
}

func TestGoTestFailure(t *testing.T) {
	old := commandRunner
	defer func() { commandRunner = old }()

	// Force the command to fail
	commandRunner = func(name string, args ...string) *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "1")
		}
		return exec.Command("false") // "false" exits 1 on Unix
	}

	_, errStr, code := run(t, "./")
	if code == 0 {
		t.Errorf("expected non-zero exit code when go test fails")
	}
	if !strings.Contains(errStr, "Error: go test failed") {
		t.Errorf("expected go test failed error, got %q", errStr)
	}
}

func TestOpenCoverageFailure(t *testing.T) {
	old := commandRunner
	defer func() { commandRunner = old }()

	calls := 0
	commandRunner = func(name string, args ...string) *exec.Cmd {
		calls++
		// First call is "go test", second is "go tool cover"
		if calls == 1 {
			// success
			return exec.Command("true")
		}
		// fail cover
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "1")
		}
		return exec.Command("false")
	}

	// Force macOS path so it doesn't exit early
	oldGOOS := getGOOS
	getGOOS = func() string { return "darwin" }
	defer func() { getGOOS = oldGOOS }()

	_, errStr, code := run(t, "-o", "./")
	if code == 0 {
		t.Errorf("expected non-zero exit code when coverage report fails to open")
	}
	if !strings.Contains(errStr, "failed to open coverage report") {
		t.Errorf("expected failed to open coverage report error, got %q", errStr)
	}
}

func TestCleanMode(t *testing.T) {
	// Simulate `go test` output that includes a "no test files" line.
	oldRunner := commandRunner
	defer func() { commandRunner = oldRunner }()

	commandRunner = func(name string, args ...string) *exec.Cmd {
		// Force an echo with both "no test files" and "ok" outputs
		script := "echo '?   github.com/entiqon/db/token/types [no test files]'; echo 'ok   github.com/entiqon/db/token/join 0.01s'"
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", script)
		}
		return exec.Command("sh", "-c", script)
	}

	out, errStr, code := run(t, "-C", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d (stderr=%q)", code, errStr)
	}

	// The [no test files] line should be filtered
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean mode to suppress 'no test files', got %q", out)
	}

	// The "ok ..." line should still be present
	if !strings.Contains(out, "ok   github.com/entiqon/db/token/join") {
		t.Errorf("expected ok line in output, got %q", out)
	}
}
