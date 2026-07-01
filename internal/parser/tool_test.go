package parser

import "testing"

func TestGetToolModel_Nil(t *testing.T) {
	if s := GetToolModel(nil); s != "unknown" {
		t.Fatalf("expected unknown, got %s", s)
	}
}

func TestGetToolModel_Empty(t *testing.T) {
	if s := GetToolModel(map[string]any{}); s != "unknown" {
		t.Fatalf("expected unknown, got %s", s)
	}
}

func TestGetToolModel_ToolAndModel(t *testing.T) {
	sessions := map[string]any{
		"s1": map[string]any{
			"agent_id": map[string]any{
				"tool":  "openai",
				"model": "gpt-4",
			},
		},
	}
	if s := GetToolModel(sessions); s != "openai/gpt-4" {
		t.Fatalf("expected openai/gpt-4, got %s", s)
	}
}

func TestGetToolModel_ToolOnly(t *testing.T) {
	sessions := map[string]any{
		"s1": map[string]any{
			"agent_id": map[string]any{
				"tool": "claude",
			},
		},
	}
	if s := GetToolModel(sessions); s != "claude" {
		t.Fatalf("expected claude, got %s", s)
	}
}

func TestGetToolModel_FirstSessionWins(t *testing.T) {
	sessions := map[string]any{
		"s1": map[string]any{
			"agent_id": map[string]any{
				"tool":  "openai",
				"model": "gpt-4",
			},
		},
		"s2": map[string]any{
			"agent_id": map[string]any{
				"tool":  "claude",
				"model": "opus",
			},
		},
	}
	if s := GetToolModel(sessions); s != "openai/gpt-4" {
		t.Fatalf("expected first session's tool, got %s", s)
	}
}

func TestGetToolModel_NoAgentID(t *testing.T) {
	sessions := map[string]any{
		"s1": map[string]any{
			"other": "data",
		},
	}
	if s := GetToolModel(sessions); s != "unknown" {
		t.Fatalf("expected unknown, got %s", s)
	}
}

func TestCalculateAiAdditions(t *testing.T) {
	entries := []NoteEntry{
		{File: "a.go", LineStart: 1, LineEnd: 5, IsAI: true},     // 5 AI lines
		{File: "a.go", LineStart: 10, LineEnd: 12, IsAI: true},   // 3 AI lines
		{File: "b.go", LineStart: 1, LineEnd: 1},                 // 1 human line
	}
	if total := CalculateAiAdditions(entries); total != 8 {
		t.Fatalf("expected 8 AI lines, got %d", total)
	}
	if total := CalculateHumanAdditions(entries); total != 1 {
		t.Fatalf("expected 1 human line, got %d", total)
	}
}

func TestCalculateAiAdditions_AllHuman(t *testing.T) {
	entries := []NoteEntry{
		{File: "a.go", LineStart: 1, LineEnd: 5},
		{File: "a.go", LineStart: 10, LineEnd: 15},
	}
	if total := CalculateAiAdditions(entries); total != 0 {
		t.Fatalf("expected 0 AI lines for human entries, got %d", total)
	}
	if total := CalculateHumanAdditions(entries); total != 11 {
		t.Fatalf("expected 11 human lines, got %d", total)
	}
}

func TestCalculateAiAdditions_Empty(t *testing.T) {
	if total := CalculateAiAdditions(nil); total != 0 {
		t.Fatalf("expected 0, got %d", total)
	}
}
