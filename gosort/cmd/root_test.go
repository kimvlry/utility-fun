package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestSortFromFile(t *testing.T) {
	// Create temporary file with unsorted lines
	tmpFile, err := os.CreateTemp("", "sorttest")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(tmpFile.Name())

	content := "banana\t2\napple\t10\ncherry\t1\n"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	err = tmpFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Capture output
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"-k", "2", "-n", tmpFile.Name()})

	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := out.String()
	want := "cherry\t1\nbanana\t2\napple\t10\n"
	if got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}
