package git

import (
	"strconv"
	"strings"
)

type DiffStat struct {
	SHA          string
	AddedLines   int
	DeletedLines int
}

func (r *Runner) Numstat(branch, since, until string, count int) (map[string]DiffStat, error) {
	args := []string{"log", "--numstat", "--format=commit %H"}
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
		args = append(args, "--max-count="+strconv.Itoa(count))
	}

	out, err := r.Run(args...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]DiffStat)
	var currentSHA string

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "commit ") {
			currentSHA = strings.TrimPrefix(line, "commit ")
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 3 || currentSHA == "" {
			continue
		}

		added, errA := strconv.Atoi(parts[0])
		deleted, errD := strconv.Atoi(parts[1])
		if errA != nil || errD != nil {
			continue
		}

		stat := result[currentSHA]
		stat.SHA = currentSHA
		stat.AddedLines += added
		stat.DeletedLines += deleted
		result[currentSHA] = stat
	}

	return result, nil
}
