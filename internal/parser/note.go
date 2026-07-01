package parser

import (
	"encoding/json"
	"strings"
	"unicode"
)

func ParseNote(noteContent string) ParsedNote {
	trimmed := strings.TrimSpace(noteContent)
	if trimmed == "" {
		return ParsedNote{Entries: nil, Sessions: nil}
	}

	lines := strings.Split(noteContent, "\n")
	sepIdx := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "---" {
			sepIdx = i
			break
		}
	}

	var prologueLines []string
	var jsonStr string

	if sepIdx > 0 {
		prologueLines = lines[:sepIdx]
	}
	if sepIdx >= 0 && sepIdx+1 < len(lines) {
		jsonStr = strings.Join(lines[sepIdx+1:], "\n")
	}

	var entries []NoteEntry
	currentFile := ""
	for _, line := range prologueLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := leadingSpaces(line)
		if indent == 0 {
			currentFile = trimmed
		} else {
			m := parseLineRange(trimmed)
			if m != nil && currentFile != "" {
				entries = append(entries, NoteEntry{
					File:      currentFile,
					LineStart: m[0],
					LineEnd:   m[1],
					IsAI:      isAISession(trimmed),
				})
			}
		}
	}

	var sessions map[string]any
	if jsonStr != "" {
		var parsed struct {
			Sessions map[string]any `json:"sessions"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil && parsed.Sessions != nil {
			sessions = parsed.Sessions
		}
	}

	return ParsedNote{
		Entries:  entries,
		Sessions: sessions,
	}
}

func isAISession(s string) bool {
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return true
	}
	if strings.HasPrefix(parts[0], "h_") {
		return false
	}
	return true
}

func leadingSpaces(s string) int {
	count := 0
	for _, r := range s {
		if r == ' ' || r == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}

func parseLineRange(s string) []int {
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return nil
	}

	rangePart := parts[len(parts)-1]
	sep := strings.IndexByte(rangePart, '-')
	if sep < 0 {
		return nil
	}

	start, end := 0, 0
	for _, r := range rangePart[:sep] {
		if !unicode.IsDigit(r) {
			return nil
		}
		start = start*10 + int(r-'0')
	}
	for _, r := range rangePart[sep+1:] {
		if !unicode.IsDigit(r) {
			return nil
		}
		end = end*10 + int(r-'0')
	}

	if start <= 0 || end <= 0 {
		return nil
	}
	return []int{start, end}
}
