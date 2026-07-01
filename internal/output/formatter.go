package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xdewx/git-ai-exporter/internal/parser"
)

type CommitOutput struct {
	SHA            string                       `json:"sha"`
	Author         string                       `json:"author"`
	AuthorEmail    string                       `json:"authorEmail"`
	Message        string                       `json:"message"`
	CommittedAt    string                       `json:"committedAt"`
	HumanAdditions int                          `json:"humanAdditions"`
	AiAdditions    int                          `json:"aiAdditions"`
	AiAccepted     int                          `json:"aiAccepted"`
	DiffAddedLines int                          `json:"diffAddedLines"`
	DiffDelLines   int                          `json:"diffDeletedLines"`
	ToolBreakdown  map[string]parser.ToolModelStat `json:"toolBreakdown"`
}

type ReportOutput struct {
	OriginUrl   string         `json:"originUrl"`
	ProjectName string         `json:"projectName"`
	Branch      string         `json:"branch"`
	Hostname    string         `json:"hostname"`
	Commits     []CommitOutput `json:"commits"`
}

func ToCommitOutput(c parser.ParsedCommit) CommitOutput {
	return CommitOutput{
		SHA:            c.SHA,
		Author:         c.Author,
		AuthorEmail:    c.AuthorEmail,
		Message:        c.Message,
		CommittedAt:    c.CommittedAt,
		HumanAdditions: c.HumanAdditions,
		AiAdditions:    c.AiAdditions,
		AiAccepted:     c.AiAccepted,
		DiffAddedLines: c.DiffAddedLines,
		DiffDelLines:   c.DiffDelLines,
		ToolBreakdown:  c.ToolBreakdown,
	}
}

func FormatJSON(report ReportOutput) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FormatPretty(report ReportOutput) string {
	var b strings.Builder
	for _, c := range report.Commits {
		b.WriteString(fmt.Sprintf("%s  %s  %s\n", shortSHA(c.SHA), c.CommittedAt[:10], truncate(c.Message, 60)))
		b.WriteString(fmt.Sprintf("      Human: %+5d  AI: %+5d  AI%%: %d%%\n",
			c.HumanAdditions, c.AiAdditions, aiPercent(c)))
		if len(c.ToolBreakdown) > 0 {
			for tool, stat := range c.ToolBreakdown {
				b.WriteString(fmt.Sprintf("      Tool: %s (AI: %d)\n", tool, stat.AiAdditions))
			}
		}
	}
	return b.String()
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func aiPercent(c CommitOutput) int {
	total := c.HumanAdditions + c.AiAdditions
	if total == 0 {
		return 0
	}
	return (c.AiAdditions * 100) / total
}
