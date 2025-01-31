package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Item represents a single service or product.
// The Quantity field can be left unset in the YAML and later updated
// via command-line input.
type Item struct {
	Description string  `yaml:"description"`
	UnitPrice   float64 `yaml:"unit_price"`
	Quantity    int     `yaml:"quantity,omitempty"`
}

// Config holds invoice-related configuration.
type Config struct {
	Sender struct {
		Name    string `yaml:"name"`
		City    string `yaml:"city"`
		Address string `yaml:"address"`
		RegNr   string `yaml:"reg_nr"`
		Phone   string `yaml:"phone"`
	} `yaml:"sender"`
	BillTo struct {
		Name    string   `yaml:"name"`
		Address []string `yaml:"address"`
	} `yaml:"bill_to"`
	ProjectName string `yaml:"project_name"`
	Payment     struct {
		BIC     string `yaml:"bic"`
		IBAN    string `yaml:"iban"`
		Address string `yaml:"address"`
	} `yaml:"payment"`
	Items []Item `yaml:"items"`
}

// LoadConfig reads the YAML configuration from the given filename.
// It first checks for a config file in the XDG config directory under the
// "invo" subfolder. If found, that file is used; otherwise, it falls back
// to the provided filename.
func LoadConfig(filename string) (*Config, error) {
	xdgConfigFile := getXDGConfigPath("invo", "config.yaml")
	if fileExists(xdgConfigFile) {
		filename = xdgConfigFile
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// getXDGConfigPath constructs the full path for a file located in the XDG
// configuration directory. It uses the XDG_CONFIG_HOME environment variable if
// set; otherwise, it defaults to "$HOME/.config".
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
