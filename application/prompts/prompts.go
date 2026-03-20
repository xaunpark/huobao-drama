package prompts

import (
	"embed"
	"strings"
)

//go:embed *.txt
var promptFS embed.FS

// Get returns the prompt template by name.
// It loads the .txt file from the embedded filesystem.
func Get(name string) string {
	b, err := promptFS.ReadFile(name)
	if err != nil {
		return ""
	}
	// Normalize line endings
	return strings.ReplaceAll(string(b), "\r\n", "\n")
}
