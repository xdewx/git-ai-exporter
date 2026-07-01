package parser

import (
	"testing"
)

func TestParseNote_Empty(t *testing.T) {
	n := ParseNote("")
	if n.Entries != nil || n.Sessions != nil {
		t.Fatal("expected nil entries and sessions for empty input")
	}

	n = ParseNote("  \n  \n")
	if n.Entries != nil || n.Sessions != nil {
		t.Fatal("expected nil entries and sessions for whitespace-only input")
	}
}

func TestParseNote_PrologueOnly(t *testing.T) {
	input := `src/main.go
  func foo 1-5
  func bar 10-15
src/utils.go
  type Config 1-3
`

	n := ParseNote(input)
	if len(n.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(n.Entries))
	}

	e0 := n.Entries[0]
	if e0.File != "src/main.go" || e0.LineStart != 1 || e0.LineEnd != 5 {
		t.Fatalf("unexpected entry 0: %+v", e0)
	}

	e1 := n.Entries[1]
	if e1.File != "src/main.go" || e1.LineStart != 10 || e1.LineEnd != 15 {
		t.Fatalf("unexpected entry 1: %+v", e1)
	}

	e2 := n.Entries[2]
	if e2.File != "src/utils.go" || e2.LineStart != 1 || e2.LineEnd != 3 {
		t.Fatalf("unexpected entry 2: %+v", e2)
	}

	if n.Sessions != nil {
		t.Fatal("expected nil sessions for prologue-only note")
	}
}

func TestParseNote_WithSessions(t *testing.T) {
	input := `src/main.go
  func foo 1-5
---
{"sessions": {"s1": {"agent_id": {"tool": "openai", "model": "gpt-4"}}}}
`
	n := ParseNote(input)
	if len(n.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(n.Entries))
	}
	if n.Entries[0].File != "src/main.go" {
		t.Fatalf("expected file src/main.go, got %s", n.Entries[0].File)
	}

	if n.Sessions == nil {
		t.Fatal("expected non-nil sessions")
	}
	s1, ok := n.Sessions["s1"]
	if !ok {
		t.Fatal("expected session s1")
	}
	m, ok := s1.(map[string]any)
	if !ok {
		t.Fatal("session value is not map")
	}
	agent, ok := m["agent_id"].(map[string]any)
	if !ok {
		t.Fatal("agent_id not found or not map")
	}
	if agent["tool"] != "openai" || agent["model"] != "gpt-4" {
		t.Fatalf("unexpected agent: %+v", agent)
	}
}

func TestParseNote_MalformedJSON(t *testing.T) {
	input := `src/main.go
  func foo 1-5
---
{invalid json}
`
	n := ParseNote(input)
	if len(n.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(n.Entries))
	}
	if n.Sessions != nil {
		t.Fatal("expected nil sessions for malformed JSON")
	}
}

func TestParseNote_NoSeparator(t *testing.T) {
	input := `src/main.go
  func foo 1-5
`

	n := ParseNote(input)
	if len(n.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(n.Entries))
	}
}

func TestParseNote_OnlyJSON(t *testing.T) {
	input := `---
{"sessions": {"s1": {"agent_id": {"tool": "openai", "model": "gpt-4"}}}}
`
	n := ParseNote(input)
	if len(n.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(n.Entries))
	}
	if n.Sessions == nil {
		t.Fatal("expected non-nil sessions")
	}
}

func TestParseNote_NoSessionsKey(t *testing.T) {
	input := `---
{"other": {"key": "value"}}
`
	n := ParseNote(input)
	if n.Sessions != nil {
		t.Fatal("expected nil sessions when JSON has no sessions key")
	}
}

func TestParseNote_InvalidLineRange(t *testing.T) {
	input := `src/main.go
  invalid-range
  foo 1-abc
  bar abc-5
  baz 1-2 extra
`
	n := ParseNote(input)
	// Only the last line should produce a valid entry
	if len(n.Entries) != 1 {
		t.Fatalf("expected 1 valid entry, got %d", len(n.Entries))
	}
	if n.Entries[0].LineStart != 1 || n.Entries[0].LineEnd != 2 {
		t.Fatalf("unexpected entry: %+v", n.Entries[0])
	}
}

func TestParseNote_CommaSeparatedRanges(t *testing.T) {
	input := `src/main.go
  s_xxx::t_yyy 10-20,30-35,40-42
`
	n := ParseNote(input)
	if len(n.Entries) != 3 {
		t.Fatalf("expected 3 entries (3 ranges), got %d", len(n.Entries))
	}
	if n.Entries[0].LineStart != 10 || n.Entries[0].LineEnd != 20 {
		t.Fatalf("unexpected first entry: %+v", n.Entries[0])
	}
	if n.Entries[1].LineStart != 30 || n.Entries[1].LineEnd != 35 {
		t.Fatalf("unexpected second entry: %+v", n.Entries[1])
	}
	if n.Entries[2].LineStart != 40 || n.Entries[2].LineEnd != 42 {
		t.Fatalf("unexpected third entry: %+v", n.Entries[2])
	}
	total := CalculateAiAdditions(n.Entries)
	if total != 20 {
		t.Fatalf("expected 20 ai additions (11+6+3), got %d", total)
	}
}
