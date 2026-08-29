package cmd

import (
	"os"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	"github.com/spf13/viper"
)

func TestAPISecretsResolveFromEnvWhenAbsentFromConfigFile(t *testing.T) {
	t.Setenv("QU_SINKS_POSTGRES_USERNAME", "ENVUSER")
	t.Setenv("QU_SINKS_POSTGRES_PASSWORD", "ENVPASS")
	t.Setenv("QU_SINKS_OTEL_BASIC_AUTH", "ENVAUTH")

	// run the real hook so this exercises rootCmd's actual env wiring
	viper.Reset()
	if err := rootCmd().PersistentPreRunE(nil, nil); err != nil {
		t.Fatal(err)
	}
	v := viper.GetViper()
	f, _ := os.CreateTemp(t.TempDir(), "*.yaml")
	// the real k8s config.yaml shape: no username/password/basic_auth
	f.WriteString("num_teams: 1\nsinks:\n  postgres:\n    host: pg\n    database: quincy-dev\n  otel:\n    endpoint: http://oo\n")
	f.Close()
	v.SetConfigFile(f.Name())
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	var cfg config.APIConfigSpec
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	t.Logf("pg.Username=%q pg.Password=%q otel.BasicAuth=%q", cfg.Sinks.PGConfig.Username, cfg.Sinks.PGConfig.Password, cfg.Sinks.OTelConfig.BasicAuth)
	t.Logf("file values preserved: pg.Host=%q pg.Database=%q NumTeams=%d", cfg.Sinks.PGConfig.Host, cfg.Sinks.PGConfig.Database, cfg.NumTeams)
	if cfg.Sinks.PGConfig.Username != "ENVUSER" || cfg.Sinks.PGConfig.Password != "ENVPASS" || cfg.Sinks.OTelConfig.BasicAuth != "ENVAUTH" {
		t.Fatal("env secrets did not land")
	}
	if cfg.Sinks.PGConfig.Host != "pg" || cfg.NumTeams != 1 {
		t.Fatal("file values clobbered")
	}
}
