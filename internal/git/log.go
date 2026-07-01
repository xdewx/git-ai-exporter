package git

import (
	"fmt"
	"strings"
	"time"
)

const recordSep = "\u27D0"
const fieldSep = "\u25A3"

type CommitRaw struct {
	SHA         string
	Author      string
	AuthorEmail string
	Message     string
	CommittedAt string
	NoteContent string
}

func (r *Runner) LogCommits(count int, branch, since, until string) ([]CommitRaw, error) {
	args := []string{"log", "--notes=ai",
		"--format=" + fieldSep + "%H" + fieldSep + "%an" + fieldSep + "%ae" + fieldSep + "%s" + fieldSep + "%aI" + fieldSep + "%N" + recordSep,
	}

	if branch != "" {
		args = append(args, branch)
	}
	if since != "" {
		args = append(args, "--since="+since)
	}
	if until != "" {
		args = append(args, "--until="+until)
	}
	if count > 0 {
		args = append(args, "--max-count="+fmtCount(count))
	}

	out, err := r.Run(args...)
	if err != nil {
		return nil, err
	}

	var commits []CommitRaw
	for _, record := range strings.Split(out, recordSep) {
		record = strings.TrimSpace(record)
		if record == "" || !strings.HasPrefix(record, fieldSep) {
			continue
		}

		parts := strings.Split(strings.TrimPrefix(record, fieldSep), fieldSep)
		if len(parts) < 5 {
			continue
		}

		c := CommitRaw{
			SHA:         parts[0],
			Author:      parts[1],
			AuthorEmail: parts[2],
			Message:     parts[3],
			CommittedAt: parts[4],
		}
		if len(parts) > 5 {
			c.NoteContent = strings.Join(parts[5:], fieldSep)
		}

		if c.SHA != "" {
			commits = append(commits, c)
		}
	}

	return commits, nil
}

func (r *Runner) WaitForNote(timeout time.Duration) bool {
	before, _ := r.Run("rev-parse", "refs/notes/ai")
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		after, err := r.Run("rev-parse", "refs/notes/ai")
		if err == nil && after != "" && after != before {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}



func fmtCount(n int) string {
	if n == 0 {
		return "0"
	}
	return fmt.Sprintf("%d", n)
}
