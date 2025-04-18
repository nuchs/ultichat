package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun_PrintsHeadline(t *testing.T) {
	service := "foo"
	msg := "bar"
	time := "2025-04-17T14:46:07.300505921+01:00"
	level := "info"
	input := fmt.Sprintf(`{"service":"%s","msg":"%s","time":"%s","level":"%s"}`, service, msg, time, level)
	scanner := bufio.NewScanner(strings.NewReader(input + "\n"))
	var out bytes.Buffer
	var errOut bytes.Buffer
	args := []string{"logfmt"}

	err := run(scanner, &out, &errOut, args)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, service) ||
		!strings.Contains(result, msg) ||
		!strings.Contains(result, time) ||
		!strings.Contains(result, level) {
		t.Errorf("expected structured log output, got: %q", result)
	}
}

func TestRun_FiltersByService(t *testing.T) {
	input := `{"service":"foo","msg":"bar"}` + "\n" + `{"service":"baz","msg":"qux"}`
	scanner := bufio.NewScanner(strings.NewReader(input + "\n"))
	var out bytes.Buffer
	var errOut bytes.Buffer
	args := []string{"logfmt", "-service=foo"}

	err := run(scanner, &out, &errOut, args)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !strings.Contains(out.String(), "foo") || strings.Contains(out.String(), "baz") {
		t.Errorf("expected only service=foo output, got: %q", out.String())
	}
}

func TestRun_PrintsAdditionalFields(t *testing.T) {
	input := `{"service":"foo","msg":"bar","extra1":"val1","extra2":42}`
	scanner := bufio.NewScanner(strings.NewReader(input + "\n"))
	var out bytes.Buffer
	var errOut bytes.Buffer
	args := []string{"logfmt"}

	err := run(scanner, &out, &errOut, args)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "extra1 = val1") || !strings.Contains(result, "extra2 = 42") {
		t.Errorf("expected additional fields in output, got: %q", result)
	}
}

func TestRun_PrintsRawLineOnUnmarshalError(t *testing.T) {
	input := "notjson\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	var out bytes.Buffer
	var errOut bytes.Buffer
	args := []string{"logfmt"}

	err := run(scanner, &out, &errOut, args)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "notjson") {
		t.Errorf("expected raw line for unmarshal error, got: %q", result)
	}
}
