package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Root is the top-level config file shape.
type Root struct {
	Devices []Device `yaml:"devices"`
}

// Device describes a wake target shown in the UI.
type Device struct {
	MAC             string `yaml:"mac"`
	IP              string `yaml:"ip,omitempty"`
	Name            string `yaml:"name"`
	NameOriginal    string `yaml:"name_original,omitempty"`
	Description     string `yaml:"description,omitempty"`
	Spec            string `yaml:"spec,omitempty"`
	Image           string `yaml:"image"`
	BackgroundColor string `yaml:"background_color"`
	ForegroundColor string `yaml:"foreground_color"`
}

// Load reads and parses the YAML file at path.
func Load(path string) (*Root, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var root Root
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	for i := range root.Devices {
		if err := root.Devices[i].validate(); err != nil {
			return nil, fmt.Errorf("device %d: %w", i, err)
		}
	}
	return &root, nil
}

func (d *Device) validate() error {
	if d.MAC == "" {
		return fmt.Errorf("mac is required")
	}
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.Image == "" {
		return fmt.Errorf("image is required")
	}
	if d.BackgroundColor == "" {
		return fmt.Errorf("background_color is required")
	}
	if d.ForegroundColor == "" {
		return fmt.Errorf("foreground_color is required")
	}
	return nil
}
