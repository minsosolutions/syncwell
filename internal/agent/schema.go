package agent

// FindingsSchema is passed to --json-schema so Go parses a validated object rather than a
// fenced block scraped out of prose. Draft-07: newer drafts are rejected.
//
// Absence evidence is deliberately absent from this schema. An Agent flags a record with
// awaiting_reply and Go states the absence (#16).
const FindingsSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["customer", "run_at", "items"],
  "additionalProperties": false,
  "properties": {
    "customer": {"type": "string"},
    "run_at": {"type": "string"},
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["key", "title", "severity", "summary", "action", "evidence"],
        "additionalProperties": false,
        "properties": {
          "key": {"type": "string", "pattern": "^[a-z0-9-]+$"},
          "title": {"type": "string"},
          "severity": {"type": "string", "enum": ["high", "medium", "low"]},
          "summary": {"type": "string"},
          "action": {"type": "string"},
          "evidence": {
            "type": "array",
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["kind", "source", "path", "locator", "at", "quote"],
              "additionalProperties": false,
              "properties": {
                "kind": {"type": "string", "enum": ["record"]},
                "source": {"type": "string", "enum": ["slack", "email", "transcript", "zammad"]},
                "path": {"type": "string"},
                "locator": {"type": "string"},
                "at": {"type": "string"},
                "quote": {"type": "string"},
                "awaiting_reply": {"type": "boolean"}
              }
            }
          },
          "proposal": {
            "type": "object",
            "required": ["kind", "subject", "body"],
            "additionalProperties": false,
            "properties": {
              "kind": {"type": "string", "enum": ["email"]},
              "subject": {"type": "string"},
              "body": {"type": "string"}
            }
          }
        }
      }
    }
  }
}`
