package agent

import (
	"testing"
	"time"

	"github.com/HackUCF/quincy/common/types"
)

func TestGetTimeout_serviceOverride(t *testing.T) {
	cfg := &AgentConfig{LoopTime: 60}
	svc := &types.Service{
		ServiceTemplate: types.ServiceTemplate{
			ServiceSpec: types.ServiceSpec{Timeout: 10.0},
		},
	}
	if got := getTimeout(cfg, svc); got != 10*time.Second {
		t.Errorf("got %v, want 10s", got)
	}
}

func TestGetTimeout_loopTimeFallback(t *testing.T) {
	cfg := &AgentConfig{LoopTime: 20}
	svc := &types.Service{} // Timeout == 0
	if got := getTimeout(cfg, svc); got != 20*time.Second {
		t.Errorf("got %v, want 20s", got)
	}
}

func TestGetTimeout_defaultFallback(t *testing.T) {
	cfg := &AgentConfig{LoopTime: 3} // < 5 → default 30s
	svc := &types.Service{}
	if got := getTimeout(cfg, svc); got != 30*time.Second {
		t.Errorf("got %v, want 30s", got)
	}
}

func TestGetTimeout_loopTimeExactly5(t *testing.T) {
	cfg := &AgentConfig{LoopTime: 5}
	svc := &types.Service{}
	if got := getTimeout(cfg, svc); got != 5*time.Second {
		t.Errorf("got %v, want 5s", got)
	}
}

func TestMakeScore_fields(t *testing.T) {
	svc := &types.Service{
		ServiceTemplate: types.ServiceTemplate{
			ServiceSpec: types.ServiceSpec{Name: "http"},
			BoxName:     "box1",
			TeamNum:     3,
		},
	}
	out := &scriptOutput{Status: Pass}
	out.Stdout.WriteString("check output")

	score := makeScore(svc, out)

	if score.ServiceName != "http" {
		t.Errorf("ServiceName = %q, want http", score.ServiceName)
	}
	if score.BoxName != "box1" {
		t.Errorf("BoxName = %q, want box1", score.BoxName)
	}
	if score.TeamNum != 3 {
		t.Errorf("TeamNum = %d, want 3", score.TeamNum)
	}
	if score.Status != Pass {
		t.Errorf("Status = %v, want Pass", score.Status)
	}
	if score.Message == "" {
		t.Error("Message is empty")
	}
}
