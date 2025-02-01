// Package config provides functionality for loading and managing invoice configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Item represents a single service or product in an invoice.
type Item struct {
	Description string  `yaml:"description"`
	UnitPrice   float64 `yaml:"unit_price"`
	Quantity    int     `yaml:"quantity,omitempty"`
}

// SenderInfo contains information about the invoice sender.
type SenderInfo struct {
	Name    string `yaml:"name"`
	City    string `yaml:"city"`
	Address string `yaml:"address"`
	RegNr   string `yaml:"reg_nr"`
	Phone   string `yaml:"phone"`
}

// BillToInfo contains information about the invoice recipient.
type BillToInfo struct {
	Name    string   `yaml:"name"`
	Address []string `yaml:"address"`
}

// PaymentInfo contains banking and payment details.
type PaymentInfo struct {
	BIC     string `yaml:"bic"`
	IBAN    string `yaml:"iban"`
	Address string `yaml:"address"`
}

// Config holds all invoice-related configuration.
type Config struct {
	Sender      SenderInfo  `yaml:"sender"`
	BillTo      BillToInfo  `yaml:"bill_to"`
	ProjectName string      `yaml:"project_name"`
	Payment     PaymentInfo `yaml:"payment"`
	Items       []Item      `yaml:"items"`
}

// LoadConfig reads the YAML configuration from the specified file.
// It first checks for a config file in the XDG config directory under the
// "invo" subfolder. If found, that file is used; otherwise, it falls back
// to the provided filename.
func LoadConfig(filename string) (*Config, error) {
	configPath, err := resolveConfigPath(filename)
	if err != nil {
		return nil, fmt.Errorf("resolving config path: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

// validate performs basic validation of the configuration.
func (c *Config) validate() error {
	if c.Sender.Name == "" {
		return fmt.Errorf("sender name is required")
	}
	if c.BillTo.Name == "" {
		return fmt.Errorf("bill to name is required")
	}
	if len(c.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	return nil
}

// resolveConfigPath determines the actual config file path to use.
func resolveConfigPath(filename string) (string, error) {
	xdgConfigFile := getXDGConfigPath("invo", "config.yaml")
	if fileExists(xdgConfigFile) {
		return xdgConfigFile, nil
	}
	return filename, nil
}

// getXDGConfigPath constructs the full path for a file in the XDG config directory.
func getXDGConfigPath(subfolder, fileName string) string {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			xdgConfig = filepath.Join(home, ".config")
		}
	}
	return filepath.Join(xdgConfig, subfolder, fileName)
}

// fileExists checks if the file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
