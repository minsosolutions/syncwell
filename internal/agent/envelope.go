package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/minsosolutions/syncwell/internal/output"
)

// Denial is one tool call Claude Code refused. A Run can be denied at every turn and still
// report success, so denials are surfaced rather than inferred from the exit code.
type Denial struct {
	ToolName  string `json:"tool_name"`
	ToolUseID string `json:"tool_use_id"`
}

type Result struct {
	Report  output.Report
	Denials []Denial
	CostUSD float64
}

type envelope struct {
	Subtype          string          `json:"subtype"`
	IsError          bool            `json:"is_error"`
	Result           *string         `json:"result"`
	StructuredOutput json.RawMessage `json:"structured_output"`
	Denials          []Denial        `json:"permission_denials"`
	CostUSD          float64         `json:"total_cost_usd"`
	APIErrorStatus   any             `json:"api_error_status"`
}

// Parse reads the one JSON object `claude -p --output-format json` writes to stdout.
//
// Success is subtype "success" *and* a structured output: a Run whose every tool call was
// denied still reports success, and turn exhaustion reports a null result. Neither is a Run
// we can believe.
func Parse(stdout []byte) (Result, error) {
	var e envelope
	if err := json.Unmarshal(stdout, &e); err != nil {
		return Result{}, fmt.Errorf("the Agent did not return a result envelope: %w", err)
	}
	if e.Subtype != "success" || e.IsError {
		return Result{}, fmt.Errorf("the Agent failed: %s%s", e.Subtype, denialSuffix(e.Denials))
	}
	if len(e.StructuredOutput) == 0 {
		return Result{}, fmt.Errorf("the Agent returned no findings object%s", denialSuffix(e.Denials))
	}
	var r output.Report
	if err := json.Unmarshal(e.StructuredOutput, &r); err != nil {
		return Result{}, fmt.Errorf("findings object does not fit the schema: %w", err)
	}
	return Result{Report: r, Denials: e.Denials, CostUSD: e.CostUSD}, nil
}

func denialSuffix(d []Denial) string {
	if len(d) == 0 {
		return ""
	}
	var names []string
	for _, x := range d {
		names = append(names, x.ToolName)
	}
	return " (denied: " + strings.Join(names, ", ") + ")"
}
