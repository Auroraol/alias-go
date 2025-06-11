package alias

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// AliasConfig represents the structure of the alias configuration file
type AliasConfig struct {
	Aliases map[string]interface{} `toml:"aliases"`
}

// AliasValue represents different types of alias values
type AliasValue interface{}

// LoadConfig loads the alias configuration from the specified file
func LoadConfig(configPath string) (*AliasConfig, error) {
	if configPath == "" {
		return &AliasConfig{Aliases: make(map[string]interface{})}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AliasConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Aliases == nil {
		config.Aliases = make(map[string]interface{})
	}

	return &config, nil
}

// GetAliasForShell returns the alias value for a specific shell
func (c *AliasConfig) GetAliasForShell(aliasName, shell string) (string, bool) {
	alias, exists := c.Aliases[aliasName]
	if !exists {
		return "", false
	}

	switch v := alias.(type) {
	case string:
		// Simple string alias, works for all shells
		return v, true
	case []interface{}:
		// Array alias, join with spaces
		var parts []string
		for _, part := range v {
			if str, ok := part.(string); ok {
				parts = append(parts, str)
			}
		}
		if len(parts) > 0 {
			return joinStrings(parts, " "), true
		}
	case map[string]interface{}:
		// Shell-specific alias
		if shellAlias, exists := v[shell]; exists {
			switch shellV := shellAlias.(type) {
			case string:
				return shellV, true
			case []interface{}:
				var parts []string
				for _, part := range shellV {
					if str, ok := part.(string); ok {
						parts = append(parts, str)
					}
				}
				if len(parts) > 0 {
					return joinStrings(parts, " "), true
				}
			}
		}
	}

	return "", false
}

// GetAllAliasesForShell returns all aliases applicable to a specific shell
func (c *AliasConfig) GetAllAliasesForShell(shell string) map[string]string {
	result := make(map[string]string)

	for aliasName := range c.Aliases {
		if value, exists := c.GetAliasForShell(aliasName, shell); exists {
			result[aliasName] = value
		}
	}

	return result
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
