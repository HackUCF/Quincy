package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HackUCF/quincy/common/types"
)

func makeTestSvc(checkName string) *types.Service {
	return &types.Service{
		ServiceTemplate: types.ServiceTemplate{
			ServiceSpec: types.ServiceSpec{
				Name:      "http",
				CheckName: types.CheckName(checkName),
			},
			BoxName: "testbox",
			Host:    "10.0.1.1",
			TeamNum: 1,
		},
	}
}

func addStubToPath(t *testing.T, name, content string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		t.Fatalf("write stub %q: %v", name, err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestDumpService_roundtrip(t *testing.T) {
	svc := makeTestSvc("stub")
	svc.User = &types.User{Username: "alice", Password: "secret"}

	f, err := dumpService(svc)
	if err != nil {
		t.Fatalf("dumpService: %v", err)
	}
	defer os.Remove(f.Name())

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var got types.Service
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != svc.Name {
		t.Errorf("Name = %q, want %q", got.Name, svc.Name)
	}
	if got.User == nil || got.User.Username != "alice" {
		t.Errorf("User = %+v, want alice", got.User)
	}
	if got.BoxName != "testbox" {
		t.Errorf("BoxName = %q, want testbox", got.BoxName)
	}
}

func TestRunCheck_pass(t *testing.T) {
	addStubToPath(t, "stub-pass", "#!/bin/sh\nexit 0\n")

	out, err := runCheck(makeTestSvc("stub-pass"), 5*time.Second)
	if err != nil {
		t.Fatalf("runCheck: %v", err)
	}
	if out.Status != Pass {
		t.Errorf("Status = %v, want Pass", out.Status)
	}
}

func TestRunCheck_fail(t *testing.T) {
	addStubToPath(t, "stub-fail", "#!/bin/sh\nexit 1\n")

	out, err := runCheck(makeTestSvc("stub-fail"), 5*time.Second)
	if err != nil {
		t.Fatalf("runCheck: %v", err)
	}
	if out.Status != Fail {
		t.Errorf("Status = %v, want Fail", out.Status)
	}
}

func TestRunCheck_timeout(t *testing.T) {
	addStubToPath(t, "stub-slow", "#!/bin/sh\nsleep 10\n")

	_, err := runCheck(makeTestSvc("stub-slow"), 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestRunCheck_readsJSONArg(t *testing.T) {
	addStubToPath(t, "stub-read", "#!/bin/sh\ncat \"$1\" > /dev/null && exit 0\n")

	out, err := runCheck(makeTestSvc("stub-read"), 5*time.Second)
	if err != nil {
		t.Fatalf("runCheck: %v", err)
	}
	if out.Status != Pass {
		t.Errorf("Status = %v, want Pass", out.Status)
	}
}

func TestRunCheck_capturesStdout(t *testing.T) {
	addStubToPath(t, "stub-echo", "#!/bin/sh\necho 'hello'\nexit 0\n")

	out, err := runCheck(makeTestSvc("stub-echo"), 5*time.Second)
	if err != nil {
		t.Fatalf("runCheck: %v", err)
	}
	if got := out.Stdout.String(); got == "" {
		t.Error("expected stdout captured, got empty string")
	}
}
