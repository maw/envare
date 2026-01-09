package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadWithPath(t *testing.T) {
	// Create a temp .env file
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "FOO=bar\nBAZ=qux\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	env := read(envPath)

	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got FOO=%s", env["FOO"])
	}
	if env["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got BAZ=%s", env["BAZ"])
	}
}

func TestReadWithQuotedValues(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := `SINGLE='single quoted'
DOUBLE="double quoted"
SPACES=value with spaces
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	env := read(envPath)

	if env["SINGLE"] != "single quoted" {
		t.Errorf("expected SINGLE='single quoted', got SINGLE=%s", env["SINGLE"])
	}
	if env["DOUBLE"] != "double quoted" {
		t.Errorf("expected DOUBLE='double quoted', got DOUBLE=%s", env["DOUBLE"])
	}
	if env["SPACES"] != "value with spaces" {
		t.Errorf("expected SPACES='value with spaces', got SPACES=%s", env["SPACES"])
	}
}

func TestReadWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := `PASSWORD=p@ss$word!
URL=https://example.com?foo=bar&baz=qux
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	env := read(envPath)

	if env["PASSWORD"] != "p@ss$word!" {
		t.Errorf("expected PASSWORD=p@ss$word!, got PASSWORD=%s", env["PASSWORD"])
	}
	if env["URL"] != "https://example.com?foo=bar&baz=qux" {
		t.Errorf("expected URL with query params, got URL=%s", env["URL"])
	}
}

// Integration tests that run the built binary
func TestBinarySingleArg(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("HELLO=world\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath, "HELLO")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Single arg without -c should output just the value (quoted)
	expected := "world"
	if strings.TrimSpace(string(output)) != expected {
		t.Errorf("expected %q, got %q", expected, string(output))
	}
}

func TestBinarySingleArgCaptureMode(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("HELLO=world\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath, "-c", "HELLO")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// With -c flag, should output KEY=value format
	expected := "HELLO=world"
	if strings.TrimSpace(string(output)) != expected {
		t.Errorf("expected %q, got %q", expected, string(output))
	}
}

func TestBinaryMultipleArgs(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath, "FOO", "BAZ")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Multiple args should output KEY=value format for each, in order
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "FOO=bar" {
		t.Errorf("expected first line FOO=bar, got %s", lines[0])
	}
	if lines[1] != "BAZ=qux" {
		t.Errorf("expected second line BAZ=qux, got %s", lines[1])
	}
}

func TestBinaryAllVars(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("ALPHA=one\nBETA=two\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// No args should output all vars
	out := string(output)
	if !strings.Contains(out, "ALPHA=one") {
		t.Errorf("expected output to contain ALPHA=one, got %s", out)
	}
	if !strings.Contains(out, "BETA=two") {
		t.Errorf("expected output to contain BETA=two, got %s", out)
	}
}

func TestBinaryShellEscaping(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	// Value with spaces and special chars that need escaping
	if err := os.WriteFile(envPath, []byte("MSG=hello world\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath, "MSG")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Should be shell-escaped (quoted)
	expected := "'hello world'"
	if string(output) != expected {
		t.Errorf("expected %q, got %q", expected, string(output))
	}
}

func TestBinaryMissingVar(t *testing.T) {
	binary := buildTestBinary(t)
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cmd := exec.Command(binary, "-f", envPath, "NONEXISTENT")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Missing var should output empty (with quotes for empty string)
	expected := "''"
	if string(output) != expected {
		t.Errorf("expected %q for missing var, got %q", expected, string(output))
	}
}

func buildTestBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "envare")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = "."
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, output)
	}
	return binary
}
