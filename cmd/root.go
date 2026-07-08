package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/xdewx/git-ai-exporter/internal/git"
	"github.com/xdewx/git-ai-exporter/internal/output"
	"github.com/xdewx/git-ai-exporter/internal/parser"
	"github.com/xdewx/git-ai-exporter/internal/reporter"
)

var (
	repoDir        string
	branch         string
	count          int
	since          string
	until          string
	outFmt         string
	push           bool
	pushURL        string
	pushToken      string
	hostname       string
	installHook    bool
	guardMode      bool
	noGuard        bool
	uninstallGuard bool
	detach         bool
)

var (
	doUpdate bool
)

var rootCmd = &cobra.Command{
	Use:     "git-ai-exporter",
	Short:   "Export git-ai commit statistics",
	Version: Version,
	Long: `Parse git-ai notes from a git repository and output structured data.

By default, outputs JSON to stdout. Use --push to send data to a git-ai-dashboard instance.
Use --install-hook to set up automatic push on every commit.`,
	SilenceUsage: true,
	RunE:         runExport,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&repoDir, "repo", "r", ".", "Git repository path")
	rootCmd.Flags().StringVarP(&branch, "branch", "b", "", "Target branch (default: current branch)")
	rootCmd.Flags().IntVarP(&count, "count", "n", 0, "Max commits to parse (0 = all)")
	rootCmd.Flags().StringVar(&since, "since", "", "Start date (e.g. 2026-01-01 or 7 days ago)")
	rootCmd.Flags().StringVar(&until, "until", "", "End date")
	rootCmd.Flags().StringVar(&outFmt, "output", "json", "Output format: json or pretty")
	rootCmd.Flags().BoolVar(&push, "push", false, "Push results to dashboard")
	rootCmd.Flags().StringVar(&pushURL, "url", "", "Dashboard collect API URL (required with --push)")
	rootCmd.Flags().StringVar(&pushToken, "token", "", "API token (required with --push)")
	rootCmd.Flags().StringVar(&hostname, "hostname", defaultHostname(), "Client hostname identifier")
	rootCmd.Flags().BoolVar(&installHook, "install-hook", false, "Install post-commit hook and exit")
	rootCmd.Flags().BoolVar(&guardMode, "guard", false, "Run guard daemon (keeps git-ai alive)")
	rootCmd.Flags().BoolVar(&noGuard, "no-guard", false, "Skip guard installation with --install-hook")
	rootCmd.Flags().BoolVar(&uninstallGuard, "uninstall-guard", false, "Remove guard auto-start service")
	rootCmd.Flags().BoolVar(&doUpdate, "update", false, "Update to latest version from GitHub")
	rootCmd.Flags().BoolVar(&detach, "detach", false, "Run push in background and exit immediately")
}

func defaultHostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func runExport(_ *cobra.Command, _ []string) error {
	if guardMode {
		return runGuard()
	}

	if uninstallGuard {
		return doUninstallGuard()
	}

	if installHook {
		absDir, err := resolvePath(repoDir)
		if err != nil {
			return fmt.Errorf("resolve repo path: %w", err)
		}
		r := git.NewRunner(absDir)
		if err := doInstallHook(r); err != nil {
			return err
		}
		if !noGuard {
			if err := doInstallGuard(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: guard installation failed: %v\n", err)
			}
		}
		return nil
	}

	if doUpdate {
		return doUpdateFn()
	}

	if detach {
		return doDetach()
	}

	if push {
		absDir, err := resolvePath(repoDir)
		if err != nil {
			return fmt.Errorf("resolve repo path: %w", err)
		}
		r := git.NewRunner(absDir)
		if pushURL == "" {
			if v, err := getGitConfig(r, "hooks.ai-exporter-url"); err == nil {
				pushURL = v
			}
		}
		if pushToken == "" {
			if v, err := getGitConfig(r, "hooks.ai-exporter-token"); err == nil {
				pushToken = v
			}
		}
		if pushURL == "" || pushToken == "" {
			return fmt.Errorf("--url and --token are required with --push (or set hooks.ai-exporter-url and hooks.ai-exporter-token via git config)")
		}
		if err := r.CheckDaemon(); err != nil {
			return err
		}
	}

	absDir, err := resolvePath(repoDir)
	if err != nil {
		return fmt.Errorf("resolve repo path: %w", err)
	}

	r := git.NewRunner(absDir)

	currentBranch, err := getCurrentBranch(r)
	if err != nil {
		return fmt.Errorf("detect branch: %w", err)
	}
	if branch == "" {
		branch = currentBranch
	}

	commits, err := r.LogCommits(count, branch, since, until)
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found")
		return nil
	}

	incomplete := false
	if push && commits[0].NoteContent == "" {
		fmt.Fprintln(os.Stderr, "Waiting for git-ai daemon to process latest commit...")
		if r.WaitForNote(commits[0].SHA, 3*time.Minute) {
			fmt.Fprintln(os.Stderr, "git-ai daemon processed the commit")
		} else {
			fmt.Fprintln(os.Stderr, "Warning: git-ai daemon not running or timed out, data may be incomplete")
			incomplete = true
		}
		commits, err = r.LogCommits(count, branch, since, until)
		if err != nil {
			return fmt.Errorf("git log: %w", err)
		}
	}

	st, err := r.Numstat(branch, since, until, count)
	if err != nil {
		return fmt.Errorf("git numstat: %w", err)
	}

	var parsedCommits []parser.ParsedCommit
	for _, c := range commits {
		note := parser.ParseNote(c.NoteContent)
		aiAdd := parser.CalculateAiAdditions(note.Entries)
		toolKey := parser.GetToolModel(note.Sessions)

		totalAdd := 0
		if s, ok := st[c.SHA]; ok {
			totalAdd = s.AddedLines
		} else {
			totalAdd = aiAdd
		}

		humanAdd := totalAdd - aiAdd
		if humanAdd < 0 {
			humanAdd = 0
		}

		toolBreakdown := make(map[string]parser.ToolModelStat)
		if toolKey != "unknown" {
			toolBreakdown[toolKey] = parser.ToolModelStat{
				AiAdditions: aiAdd,
				AiAccepted:  aiAdd,
			}
		}

		parsedCommits = append(parsedCommits, parser.ParsedCommit{
			SHA:            c.SHA,
			Author:         c.Author,
			AuthorEmail:    c.AuthorEmail,
			Message:        c.Message,
			CommittedAt:    c.CommittedAt,
			HumanAdditions: humanAdd,
			AiAdditions:    aiAdd,
			AiAccepted:     aiAdd,
			DiffAddedLines: totalAdd,
			DiffDelLines:   st[c.SHA].DeletedLines,
			ToolBreakdown:  toolBreakdown,
		})
	}

	sort.Slice(parsedCommits, func(i, j int) bool {
		return parsedCommits[i].CommittedAt > parsedCommits[j].CommittedAt
	})

	outputCommits := make([]output.CommitOutput, len(parsedCommits))
	for i, c := range parsedCommits {
		outputCommits[i] = output.ToCommitOutput(c)
	}

	originURL, err := getOriginURL(r)
	if err != nil {
		originURL = "unknown"
	}

	report := output.ReportOutput{
		OriginUrl:   originURL,
		ProjectName: extractProjectName(originURL),
		Branch:      branch,
		Hostname:    hostname,
		Incomplete:  incomplete,
		Commits:     outputCommits,
	}

	if push {
		result, err := reporter.Push(pushURL, pushToken, hostname, report)
		if err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Pushed %d commits, projectId: %s\n", result.Imported, result.ProjectID)
	}

	switch outFmt {
	case "pretty":
		fmt.Print(output.FormatPretty(report))
	default:
		jsonOut, err := output.FormatJSON(report)
		if err != nil {
			return fmt.Errorf("format output: %w", err)
		}
		fmt.Println(jsonOut)
	}

	return nil
}
