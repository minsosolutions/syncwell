package agent

import (
	"strings"
	"testing"
)

const okEnvelope = `{"type":"result","subtype":"success","is_error":false,"result":"done",
  "permission_denials":[],"total_cost_usd":0.01,
  "structured_output":{"customer":"ebbahus","run_at":"2026-09-16T09:00:00Z","items":[
    {"key":"k","title":"t","severity":"high","evidence":[{"kind":"record","source":"email","path":"email/x.md","quote":"q"}]}]}}`

func TestParseReadsTheStructuredOutput(t *testing.T) {
	res, err := Parse([]byte(okEnvelope))
	if err != nil {
		t.Fatal(err)
	}
	if res.Report.Customer != "ebbahus" || len(res.Report.Items) != 1 {
		t.Fatalf("report did not come through: %+v", res.Report)
	}
}

// Verified in #2: a run denied at every turn still reports subtype "success" with
// is_error false. Success without structured output is a failed Run, not an empty one.
func TestParseRejectsSuccessWithNothingInIt(t *testing.T) {
	_, err := Parse([]byte(`{"subtype":"success","is_error":false,
	  "permission_denials":[{"tool_name":"Read","tool_use_id":"t1"}]}`))
	if err == nil {
		t.Fatal("a Run that produced no findings object was treated as success")
	}
	if !strings.Contains(err.Error(), "Read") {
		t.Errorf("the denial that caused it is not in the error: %v", err)
	}
}

// Verified in #2: --max-turns exhaustion exits 1 with result null.
func TestParseNamesTurnExhaustion(t *testing.T) {
	_, err := Parse([]byte(`{"subtype":"error_max_turns","is_error":true,"result":null}`))
	if err == nil || !strings.Contains(err.Error(), "error_max_turns") {
		t.Fatalf("turn exhaustion was not reported clearly: %v", err)
	}
}

func TestParseSurfacesDenialsOnASuccessfulRun(t *testing.T) {
	res, err := Parse([]byte(`{"subtype":"success","is_error":false,
	  "permission_denials":[{"tool_name":"Read","tool_use_id":"t1"}],
	  "structured_output":{"customer":"c","items":[]}}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Denials) != 1 || res.Denials[0].ToolName != "Read" {
		t.Fatalf("denials were swallowed: %+v", res.Denials)
	}
}
