package parser

type NoteEntry struct {
	File      string
	LineStart int
	LineEnd   int
}

type ParsedNote struct {
	Entries  []NoteEntry
	Sessions map[string]any
}

type ToolModelStat struct {
	AiAdditions int `json:"aiAdditions"`
	AiAccepted  int `json:"aiAccepted"`
}

type ParsedCommit struct {
	SHA            string
	Author         string
	AuthorEmail    string
	Message        string
	CommittedAt    string
	HumanAdditions int
	AiAdditions    int
	AiAccepted     int
	DiffAddedLines int
	DiffDelLines   int
	ToolBreakdown  map[string]ToolModelStat
}
