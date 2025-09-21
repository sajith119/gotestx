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
	// In tests, simulate go test instead of running it for real.
	// If coverage is requested, emit a fake coverage summary line
	// so quiet mode can pick it up.
	commandRunner = func(name string, args ...string) *exec.Cmd {
		all := append([]string{name}, args...)
		joined := strings.Join(all, " ")
		if strings.Contains(joined, "-coverprofile=coverage.out") {
			// Simulate go test with coverage output
			return exec.Command("echo", "ok   mypkg   0.01s  coverage: 75.0% of statements")
		}
		// Default: just echo the command
		return exec.Command("echo", joined)
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

func TestCleanView(t *testing.T) {
	out, _, code := run(t, "-V", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean-view mode to suppress 'no test files', got %q", out)
	}
}

func TestCleanViewWithCoverage(t *testing.T) {
	out, _, code := run(t, "-cV", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "coverage") {
		t.Errorf("expected coverage message in output, got %q", out)
	}
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean-view mode to suppress 'no test files', got %q", out)
	}
}

func TestQuietMode(t *testing.T) {
	out, _, code := run(t, "-q", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	// Quiet mode should not show "Running ..." but must show summary
	if strings.Contains(out, "Running") {
		t.Errorf("quiet mode should suppress verbose info, got %q", out)
	}
	if !strings.Contains(out, "Tests finished successfully") &&
		!strings.Contains(out, "coverage:") {
		t.Errorf("quiet mode should show summary or coverage, got %q", out)
	}
}

func TestOpenCoverage(t *testing.T) {
	out, err, code := run(t, "-o", "./")
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
	if strings.Contains(out, "Running") {
		t.Errorf("quiet mode should suppress verbose info, got %q", out)
	}
	if !strings.Contains(out, "coverage:") {
		t.Errorf("quiet+coverage should show coverage summary, got %q", out)
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

	getGOOS = func() string { return "linux" }
	_, errStr, code := run(t, "-o", "./")
	if code == 0 {
		t.Errorf("expected non-zero exit code")
	}
	if !strings.Contains(errStr, "only supported on macOS") {
		t.Errorf("expected macOS-only error, got %q", errStr)
	}

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

	commandRunner = func(name string, args ...string) *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "1")
		}
		return exec.Command("false")
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
		if calls == 1 {
			return exec.Command("true")
		}
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "1")
		}
		return exec.Command("false")
	}

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

func TestCleanViewMode(t *testing.T) {
	oldRunner := commandRunner
	defer func() { commandRunner = oldRunner }()

	commandRunner = func(name string, args ...string) *exec.Cmd {
		script := "echo '?   github.com/entiqon/db/token/types [no test files]'; echo 'ok   github.com/entiqon/db/token/join 0.01s'"
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", script)
		}
		return exec.Command("sh", "-c", script)
	}

	out, errStr, code := run(t, "-V", "./")
	if code != 0 {
		t.Errorf("expected exit 0, got %d (stderr=%q)", code, errStr)
	}
	if strings.Contains(out, "[no test files]") {
		t.Errorf("expected clean-view mode to suppress 'no test files', got %q", out)
	}
	if !strings.Contains(out, "ok   github.com/entiqon/db/token/join") {
		t.Errorf("expected ok line in output, got %q", out)
	}
}

func TestQuietModeFailure(t *testing.T) {
	old := commandRunner
	defer func() { commandRunner = old }()

	// Force the command to fail
	commandRunner = func(name string, args ...string) *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "1")
		}
		return exec.Command("false") // exits 1 on Unix
	}

	stdout, stderr, code := run(t, "-q", "./")
	if code == 0 {
		t.Errorf("expected non-zero exit code in quiet mode failure")
	}
	if stdout != "" {
		t.Errorf("expected no stdout in quiet failure, got %q", stdout)
	}
	if !strings.Contains(stderr, "❌ Tests failed") {
		t.Errorf("expected quiet mode failure message in stderr, got %q", stderr)
	}

	_, errStr, _ := run(t, "-q", "./does-not-exist")
	if !strings.Contains(errStr, "❌ Tests failed") {
		t.Errorf("expected quiet mode failure message, got %q", errStr)
	}
}

func TestDefaultPackageIsDotDotDot(t *testing.T) {
	out, _, code := run(t) // no args
	if code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "./...") {
		t.Errorf("expected default ./... in command, got %q", out)
	}
}
