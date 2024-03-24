package analysis

import (
	"fmt"
	"strings"

	"edulsp/lsp"
)

type State struct {
	// Map of file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic
	for row, line := range strings.Split(text, "\n") {
		if idx := strings.Index(line, "VS Code"); idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Please make sure we use a good language in this video",
			})
		}

		if idx := strings.Index(line, "Neovim"); idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: 4,
				Source:   "Common Sense",
				Message:  "Great choice :)",
			})
		}
	}

	return diagnostics
}

func (s *State) OpenDocument(document, text string) []lsp.Diagnostic {
	s.Documents[document] = text

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(document, text string) []lsp.Diagnostic {
	s.Documents[document] = text

	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// In real life, this would look up the type in our type analysis code ...

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			Message: lsp.Message{
				RPC: "2.0",
			},
			ID: &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Charaters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// In real life, this would look up the definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			Message: lsp.Message{
				RPC: "2.0",
			},
			ID: &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:     position.Line - 1,
					Charater: 0,
				},
				End: lsp.Position{
					Line:     position.Line - 1,
					Charater: 0,
				},
			},
		},
	}
}

func (s *State) CodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]

	var actions []lsp.CodeAction
	for row, line := range strings.Split(text, "\n") {
		if idx := strings.Index(line, "VS Code"); idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}

	response := lsp.CodeActionResponse{
		Response: lsp.Response{
			Message: lsp.Message{
				RPC: "2.0",
			},
			ID: &id,
		},
		Result: actions,
	}

	return response
}

func (s *State) Completion(id int, uri string) lsp.CompletionResponse {
	// Ask your static analysis tools to figure out good completions
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos. Don't forget to like & subscribe to streamers using it :)",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			Message: lsp.Message{
				RPC: "2.0",
			},
			ID: &id,
		},
		Result: items,
	}

	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:     line,
			Charater: start,
		},
		End: lsp.Position{
			Line:     line,
			Charater: end,
		},
	}
}
