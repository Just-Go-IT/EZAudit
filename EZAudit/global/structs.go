package global

import "regexp"

type ConfigStruct struct {
	OS              string    `json:"os"`
	Verbosity       int       `json:"verbosity"`
	MaxOutputLength int       `json:"maxOutputLength"`
	Censor          []string  `json:"censor"`
	Commands        []Command `json:"commands"`
}

type Command struct {
	Name   string `json:"name"`
	Steps  []Step `json:"steps"`
	Index  int
	Status int
}

type Step struct {
	ModuleName       string                 `json:"module"`
	AllowFailure     bool                   `json:"allowFailure"`
	DontSaveArtifact bool                   `json:"dontSaveArtifact"`
	UsePipe          bool                   `json:"usePipe"`
	Parameters       map[string]interface{} `json:"parameter"`
	Comparison       string                 `json:"comparison"`
	ExpectedValue    string                 `json:"expected"`
	Censor           []string               `json:"censor"`
	OnSuccess        *Step                  `json:"onSuccess"`
	OnFailure        *Step                  `json:"onFailure"`
	NeedsElevation   bool                   `json:"needsElevation"`
	Path             Path
	Module           Module
	Regex            []regexp.Regexp
	GetNext          *Step
}

type Path struct {
	CommandName  string
	CommandIndex int
	StepIndex    int
	ExactPath    string
}

type ReportStruct struct {
	Summary Summary `json:"auditSummary"`
	General General `json:"general"`
	Audits  []Audit `json:"audit"`
}

type General struct {
	Date     string `json:"date"`
	Runtime  string `json:"executionTime"`
	OS       string `json:"os"`
	Admin    bool   `json:"admin"`
	User     string `json:"user"`
	UserName string `json:"userName"`
}

type Audit struct {
	Name       string      `json:"name"`
	ID         int         `json:"auditID"`
	AuditSteps []AuditStep `json:"steps"`
	Status     string      `json:"status"`
}

type AuditStep struct {
	ID          int        `json:"auditStepID"`
	Expected    string     `json:"expected"`
	Comparison  string     `json:"comparison"`
	RealValue   string     `json:"realValue"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	OnSuccess   *AuditStep `json:"onSuccess,omitempty"`
	OnFailure   *AuditStep `json:"onFailure,omitempty"`
}

type Summary struct {
	PassedPercentage string `json:"passedPercentage"`
	Passed           int    `json:"passed"`
	Failed           int    `json:"failed"`
	Errors           int    `json:"errors"`
}
