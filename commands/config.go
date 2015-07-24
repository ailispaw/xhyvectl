package commands

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	UUID_FILE_NAME = ".uuid"
	MAC_FILE_NAME  = ".mac"
)

type Config struct {
	Type       string     `yaml:"type"`
	BootConfig BootConfig `yaml:"boot"`
	SSHConfig  SSHConfig  `yaml:"ssh"`
	VMConfig   VMConfig   `yaml:"vm,omitempty"`
}

type BootConfig struct {
	Kernel      string `yaml:"kernel,omitempty"`
	Initrd      string `yaml:"initrd,omitempty"`
	CommandLine string `yaml:"cmdline,omitempty"`

	UserBoot   string `yaml:"userboot,omitempty"`
	BootVolume string `yaml:"uservolume,omitempty"`
	KernelENV  string `yaml:"kernelenv,omitempty"`
}

type SSHConfig struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password,omitempty"`
	KeyPath  string `yaml:"key-path,omitempty"`
}

type VMConfig struct {
	ACPI   bool     `yaml:"acpi,omitempty"`
	CPUs   int      `yaml:"cpus,omitempty"`
	Memory string   `yaml:"memory,omitempty"`
	HDD    []string `yaml:"hdd,omitempty"`
	Net    []string `yaml:"net,omitempty"`
	PCI    []string `yaml:"pci,omitempty"`
	LPC    []string `yaml:"lpc,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	path = os.ExpandEnv(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
