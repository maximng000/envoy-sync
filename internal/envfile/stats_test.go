package envfile

import (
	"testing"
)

func statsEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "APP_DEBUG", Value: ""},
		{Key: "APP_NAME", Value: "duplicate"},
	}
}

func TestGatherStats_Total(t *testing.T) {
	s := GatherStats(statsEntries())
	if s.Total != 7 {
		t.Errorf("expected Total=7, got %d", s.Total)
	}
}

func TestGatherStats_Secrets(t *testing.T) {
	s := GatherStats(statsEntries())
	if s.Secrets != 2 {
		t.Errorf("expected Secrets=2, got %d", s.Secrets)
	}
}

func TestGatherStats_Empty(t *testing.T) {
	s := GatherStats(statsEntries())
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
}

func TestGatherStats_Duplicates(t *testing.T) {
	s := GatherStats(statsEntries())
	if s.Duplicates != 1 {
		t.Errorf("expected Duplicates=1, got %d", s.Duplicates)
	}
}

func TestGatherStats_Prefixes(t *testing.T) {
	s := GatherStats(statsEntries())
	if s.Prefixes["APP"] != 3 {
		t.Errorf("expected APP prefix count=3, got %d", s.Prefixes["APP"])
	}
	if s.Prefixes["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", s.Prefixes["DB"])
	}
}

func TestTopPrefixes_Order(t *testing.T) {
	s := GatherStats(statsEntries())
	top := TopPrefixes(s, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top prefixes, got %d", len(top))
	}
	if top[0] != "APP" {
		t.Errorf("expected top prefix APP, got %s", top[0])
	}
}

func TestGatherStats_EmptyInput(t *testing.T) {
	s := GatherStats([]Entry{})
	if s.Total != 0 || s.Secrets != 0 || s.Empty != 0 {
		t.Error("expected all zeros for empty input")
	}
}
