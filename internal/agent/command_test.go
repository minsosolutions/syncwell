package agent

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func argsOf(t *testing.T) (args []string, settings string) {
	t.Helper()
	cmd := Command(t.Context(), Options{
		CustomerDir:   "/data/customers/ebbahus",
		ReferencesDir: "/data/references",
		Prompt:        Prompt("ebbahus", "2026-09-17T09:00:00Z", []string{"email/x.md"}, []string{"product/y.md"}),
	})
	if cmd.Dir != "/data/customers/ebbahus" {
		t.Fatalf("the Run does not start in the Customer Directory: %q", cmd.Dir)
	}
	for i, a := range cmd.Args {
		if a == "--settings" && i+1 < len(cmd.Args) {
			settings = cmd.Args[i+1]
		}
	}
	return cmd.Args, settings
}

// ADR-0001: an Agent has no path to another Customer Directory, structurally. These flags
// are that structure — see docs/research/claude-headless-contract.md for what each enforces.
func TestCommandConfinesTheAgent(t *testing.T) {
	args, settings := argsOf(t)
	joined := strings.Join(args, " ")

	for _, want := range []string{
		"--tools Read",              // the whole tool list: no Bash, so no subprocess and no network
		"--permission-mode dontAsk", // nothing that would prompt can run
		"--permission-prompts none", // and nothing waits on a human to allow it
		"--bare",                    // a Customer Directory cannot configure the Agent reading it
		"--strict-mcp-config",       // no MCP servers
		"--output-format json",      // one envelope on stdout
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q from the Run's flags", want)
		}
	}
	if strings.Contains(joined, "--add-dir") {
		t.Error("--add-dir grants write and loads skills from the tree; use additionalDirectories")
	}
	if strings.Contains(joined, "bypassPermissions") {
		t.Error("the Run must never bypass permissions")
	}

	var s struct {
		Permissions struct {
			BlockReads            bool     `json:"blockReadsOutsideWorkingDirectories"`
			AdditionalDirectories []string `json:"additionalDirectories"`
			Deny                  []string `json:"deny"`
		} `json:"permissions"`
	}
	if err := json.Unmarshal([]byte(settings), &s); err != nil {
		t.Fatalf("--settings is not JSON: %v", err)
	}
	if !s.Permissions.BlockReads {
		t.Error("blockReadsOutsideWorkingDirectories is off; the fence is then only a prompt")
	}
	if len(s.Permissions.AdditionalDirectories) != 1 || s.Permissions.AdditionalDirectories[0] != "/data/references" {
		t.Errorf("References are not the one extra directory: %v", s.Permissions.AdditionalDirectories)
	}
	if !slices.Contains(s.Permissions.Deny, "Edit(//data/references/**)") {
		t.Errorf("References are writable: %v — the // prefix is required for an absolute path", s.Permissions.Deny)
	}
}

func TestCommandCarriesTheWorkingSetInThePrompt(t *testing.T) {
	args, _ := argsOf(t)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "email/x.md") || !strings.Contains(joined, "product/y.md") {
		t.Error("the Agent was given no file list; with Read as its only tool it cannot find anything")
	}
}

func TestCommandHidesTheManifestFromTheAgent(t *testing.T) {
	_, settings := argsOf(t)
	if !strings.Contains(settings, `Read(state/manifest.json)`) {
		t.Error("the Manifest is readable by the Agent; reconciliation detail belongs to Go (#3, #10)")
	}
}
