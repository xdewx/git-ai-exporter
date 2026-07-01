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

	if sepIdx >= 0 {
		prologueLines = lines[:sepIdx]
		if sepIdx+1 < len(lines) {
			jsonStr = strings.Join(lines[sepIdx+1:], "\n")
		}
	} else {
		prologueLines = lines
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
			ranges := parseLineRanges(trimmed)
			for _, r := range ranges {
				if currentFile != "" {
					entries = append(entries, NoteEntry{
						File:      currentFile,
						LineStart: r[0],
						LineEnd:   r[1],
					})
				}
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

func parseLineRanges(s string) [][]int {
	parts := strings.Fields(s)
	for i := len(parts) - 1; i >= 0; i-- {
		ranges := parseRangeStr(parts[i])
		if len(ranges) > 0 {
			return ranges
		}
	}
	return nil
}

func parseRangeStr(s string) [][]int {
	var result [][]int
	for _, seg := range strings.Split(s, ",") {
		sep := strings.IndexByte(seg, '-')
		if sep <= 0 || sep >= len(seg)-1 {
			continue
		}

		start, end := 0, 0
		valid := true
		for _, r := range seg[:sep] {
			if !unicode.IsDigit(r) {
				valid = false
				break
			}
			start = start*10 + int(r-'0')
		}
		if !valid || start <= 0 {
			continue
		}
		for _, r := range seg[sep+1:] {
			if !unicode.IsDigit(r) {
				valid = false
				break
			}
			end = end*10 + int(r-'0')
		}
		if valid && end > 0 {
			result = append(result, []int{start, end})
		}
	}
	return result
}
