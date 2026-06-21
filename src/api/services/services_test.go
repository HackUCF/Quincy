package services

import (
	"context"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
)

func minimalServicesConfig() *config.APIConfigSpec {
	return &config.APIConfigSpec{
		NumTeams: 2,
		Boxes: []config.BoxSpec{
			{
				Name: "box1",
				Host: "10.0.{}.1",
				Services: []types.ServiceSpec{
					{Name: "http", CheckName: "stub"},
					{Name: "ssh", CheckName: "stub"},
				},
			},
		},
	}
}

func setupTeamRange(cfg *config.APIConfigSpec) {
	config.TeamRange = make([]types.TeamNum, 0, cfg.NumTeams)
	for i := types.TeamNum(1); i <= cfg.NumTeams; i++ {
		config.TeamRange = append(config.TeamRange, i)
	}
}

func TestInitServices_count(t *testing.T) {
	cfg := minimalServicesConfig()
	setupTeamRange(cfg)
	t.Cleanup(resetForTest)

	if err := InitServices(cfg); err != nil {
		t.Fatalf("InitServices: %v", err)
	}

	// 1 box × 2 services × 2 teams = 4
	if servicesLen != 4 {
		t.Errorf("servicesLen = %d, want 4", servicesLen)
	}
}

func TestGetNext_noCredentials(t *testing.T) {
	cfg := minimalServicesConfig()
	setupTeamRange(cfg)
	t.Cleanup(resetForTest)

	if err := InitServices(cfg); err != nil {
		t.Fatalf("InitServices: %v", err)
	}

	svc, err := GetNext(context.Background(), cfg, nil) // nil db: safe when UserList == ""
	if err != nil {
		t.Fatalf("GetNext: %v", err)
	}
	if svc.User != nil {
		t.Error("expected no user for service without UserList")
	}
}

func TestGetNext_roundRobin(t *testing.T) {
	cfg := minimalServicesConfig()
	setupTeamRange(cfg)
	t.Cleanup(resetForTest)

	if err := InitServices(cfg); err != nil {
		t.Fatalf("InitServices: %v", err)
	}

	type key struct {
		name    types.ServiceName
		box     types.BoxName
		teamNum types.TeamNum
	}
	seen := make(map[key]bool)
	for i := 0; i < int(servicesLen); i++ {
		svc, err := GetNext(context.Background(), cfg, nil)
		if err != nil {
			t.Fatalf("GetNext at i=%d: %v", i, err)
		}
		seen[key{svc.Name, svc.BoxName, svc.TeamNum}] = true
	}

	if len(seen) != int(servicesLen) {
		t.Errorf("one cycle covered %d unique services, want %d", len(seen), servicesLen)
	}
}

func TestGetNext_wrapsAround(t *testing.T) {
	cfg := minimalServicesConfig()
	setupTeamRange(cfg)
	t.Cleanup(resetForTest)

	if err := InitServices(cfg); err != nil {
		t.Fatalf("InitServices: %v", err)
	}

	for i := 0; i < int(servicesLen)*3; i++ {
		if _, err := GetNext(context.Background(), cfg, nil); err != nil {
			t.Fatalf("GetNext at i=%d: %v", i, err)
		}
	}
}
