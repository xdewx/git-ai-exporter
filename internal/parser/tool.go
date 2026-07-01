package parser

import "fmt"

type AgentID struct {
	Tool  string `json:"tool"`
	Model string `json:"model"`
}

func GetToolModel(sessions map[string]any) string {
	if sessions == nil {
		return "unknown"
	}

	for _, v := range sessions {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		agentRaw, ok := m["agent_id"]
		if !ok {
			continue
		}
		agentMap, ok := agentRaw.(map[string]any)
		if !ok {
			continue
		}

		tool, _ := agentMap["tool"].(string)
		model, _ := agentMap["model"].(string)

		if tool != "" && model != "" {
			return fmt.Sprintf("%s/%s", tool, model)
		}
		if tool != "" {
			return tool
		}
	}

	return "unknown"
}

func CalculateAiAdditions(entries []NoteEntry) int {
	total := 0
	for _, e := range entries {
		if e.IsAI {
			total += e.LineEnd - e.LineStart + 1
		}
	}
	return total
}

func CalculateHumanAdditions(entries []NoteEntry) int {
	total := 0
	for _, e := range entries {
		if !e.IsAI {
			total += e.LineEnd - e.LineStart + 1
		}
	}
	return total
}
