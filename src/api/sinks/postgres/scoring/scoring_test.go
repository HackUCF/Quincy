package scoring

import "testing"

func TestRound(t *testing.T) {
	cases := []struct {
		f    float64
		n    int
		want float64
	}{
		{3.14159, 2, 3.14},
		{100.0, 0, 100.0},
		{0.0, 2, 0.0},
		{2.005, 2, 2.01},
	}
	for _, c := range cases {
		got := round(c.f, c.n)
		if got != c.want {
			t.Errorf("round(%v, %d) = %v, want %v", c.f, c.n, got, c.want)
		}
	}
}

func TestGetResult_zero(t *testing.T) {
	r := getResult(0, 0)
	if r.TotalChecks != 0 || r.UptimePercent != 0 || r.ChecksFailed != 0 {
		t.Errorf("expected all zeros, got %+v", r)
	}
}

func TestGetResult_allPass(t *testing.T) {
	r := getResult(10, 10)
	if r.UptimePercent != 100.0 {
		t.Errorf("uptime = %v, want 100.0", r.UptimePercent)
	}
	if r.ChecksFailed != 0 {
		t.Errorf("failed = %d, want 0", r.ChecksFailed)
	}
	if r.TotalChecks != 10 {
		t.Errorf("total = %d, want 10", r.TotalChecks)
	}
}

func TestGetResult_half(t *testing.T) {
	r := getResult(10, 5)
	if r.UptimePercent != 50.0 {
		t.Errorf("uptime = %v, want 50.0", r.UptimePercent)
	}
	if r.ChecksFailed != 5 {
		t.Errorf("failed = %d, want 5", r.ChecksFailed)
	}
}

func TestGetResult_allFail(t *testing.T) {
	r := getResult(8, 0)
	if r.UptimePercent != 0.0 {
		t.Errorf("uptime = %v, want 0.0", r.UptimePercent)
	}
	if r.ChecksPassed != 0 {
		t.Errorf("passed = %d, want 0", r.ChecksPassed)
	}
}
