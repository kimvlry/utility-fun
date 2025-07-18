package main

import (
	"os/exec"
	"strings"
	"testing"
)

func runWithArgs(t *testing.T, args ...string) (stdout string, err error) {
	t.Helper()

	bin := "./ntpnow_test_bin"
	build := exec.Command("go", "build", "-o", bin)
	if err := build.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	defer func(command *exec.Cmd) {
		err := command.Run()
		if err != nil {
			t.Fatalf("failed to run command: %v", err)
		}
	}(exec.Command("rm", bin))

	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestTimeoutTooShort(t *testing.T) {
	out, err := runWithArgs(t, "-timeout=1ms")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(out, "timeout") {
		t.Errorf("expected timeout error message, got: %s", out)
	}
}

func TestInvalidServer(t *testing.T) {
	out, err := runWithArgs(t, "-server=invalid.doesnotexist.example.com")
	if err == nil {
		t.Fatal("expected DNS resolution error, got nil")
	}
	if !strings.Contains(out, "Error querying") && !strings.Contains(out, "timeout") {
		t.Errorf("unexpected error output: %s", out)
	}
}

func TestValidQuery(t *testing.T) {
	out, err := runWithArgs(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "UTC:") {
		t.Errorf("expected UTC output, got: %s", out)
	}
	if !strings.Contains(out, "Local:") {
		t.Errorf("expected Local output, got: %s", out)
	}
}
