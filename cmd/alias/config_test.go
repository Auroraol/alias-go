package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "alias.toml")

	configContent := `[aliases]
c = "cargo"
ll = "ls -la"

[aliases.cls]
bash = "clear"
zsh = "clear"

[aliases.multi]
bash = ["echo", "hello", "world"]
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test simple alias
	if value, exists := config.GetAliasForShell("c", "bash"); !exists {
		t.Error("Expected alias 'c' to exist for bash")
	} else if value != "cargo" {
		t.Errorf("Expected alias 'c' to be 'cargo', got '%s'", value)
	}

	// Test shell-specific alias
	if value, exists := config.GetAliasForShell("cls", "bash"); !exists {
		t.Error("Expected alias 'cls' to exist for bash")
	} else if value != "clear" {
		t.Errorf("Expected alias 'cls' to be 'clear', got '%s'", value)
	}

	// Test non-existent shell
	if _, exists := config.GetAliasForShell("cls", "fish"); exists {
		t.Error("Expected alias 'cls' to not exist for fish")
	}
}

func TestGetAllAliasesForShell(t *testing.T) {
	config := &AliasConfig{
		Aliases: map[string]interface{}{
			"c":  "cargo",
			"ll": "ls -la",
			"cls": map[string]interface{}{
				"bash": "clear",
				"zsh":  "clear",
			},
		},
	}

	bashAliases := config.GetAllAliasesForShell("bash")

	expectedBashAliases := map[string]string{
		"c":   "cargo",
		"ll":  "ls -la",
		"cls": "clear",
	}

	if len(bashAliases) != len(expectedBashAliases) {
		t.Errorf("Expected %d bash aliases, got %d", len(expectedBashAliases), len(bashAliases))
	}

	for name, expectedValue := range expectedBashAliases {
		if actualValue, exists := bashAliases[name]; !exists {
			t.Errorf("Expected bash alias '%s' to exist", name)
		} else if actualValue != expectedValue {
			t.Errorf("Expected bash alias '%s' to be '%s', got '%s'", name, expectedValue, actualValue)
		}
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	_, err := LoadConfig("/non/existent/path")
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}

func TestEscapeFunctions(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		fn       func(string) string
	}{
		{"hello world", "hello world", escapeSingleQuotes},
		{"hello'world", "hello'\"'\"'world", escapeSingleQuotes},
		{"it's working", "it'\"'\"'s working", escapeSingleQuotes},
		{"hello world", "hello world", escapeDoubleQuotes},
		{"hello\"world", "hello\\\"world", escapeDoubleQuotes},
	}

	for _, tc := range testCases {
		result := tc.fn(tc.input)
		if result != tc.expected {
			t.Errorf("Expected %s('%s') = '%s', got '%s'",
				getFunctionName(tc.fn), tc.input, tc.expected, result)
		}
	}
}

func getFunctionName(fn func(string) string) string {
	// This is a simple way to identify the function for test output
	test := fn("'")
	if test == "'\"'\"'" {
		return "escapeSingleQuotes"
	}
	return "escapeDoubleQuotes"
}

func TestIsSimpleCommand(t *testing.T) {
	testCases := []struct {
		command  string
		expected bool
	}{
		{"ls", true},
		{"cargo", true},
		{"ls -la", false},
		{"echo hello", false},
		{"cat file.txt | grep pattern", false},
		{"echo hello && echo world", false},
		{"echo hello; echo world", false},
	}

	for _, tc := range testCases {
		result := isSimpleCommand(tc.command)
		if result != tc.expected {
			t.Errorf("Expected isSimpleCommand('%s') = %v, got %v", tc.command, tc.expected, result)
		}
	}
}
