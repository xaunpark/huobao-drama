package fixed

import (
	"embed"
	"strings"
)

//go:embed *.txt
var fixedFS embed.FS

// Get returns the fixed format instruction by prompt type key.
// These are the JSON schema / output format requirements
// that MUST always be appended to prompts regardless of user customization.
func Get(promptType string) string {
	filename := promptType + ".txt"
	b, err := fixedFS.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(b), "\r\n", "\n")
}
