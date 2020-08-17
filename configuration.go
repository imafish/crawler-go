package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Startpages []startpage
	Rules      []Rule
}

type startpage struct {
	URL   string
	Group group
}

type group struct {
	DirPattern string `yaml:"dir_pattern"`
	GroupBy    string `yaml:"group_by"`
}

// Rule represents an execution rule
type Rule struct {
	Name        string
	Config      config
	Matches     []match
	Actions     []action
	PostActions []postAction `yaml:"post_actions"`
}

type config struct {
	Webdriver string
}

type match struct {
	URL   string
	Title string
}

type action struct {
	Targets      []target
	DownloadFile *downloadFileAction `yaml:"download_file,omitempty"`
	GrabText     *grabTextAction     `yaml:"grab_text,omitempty"`
	ProcessLink  *processLinkAction  `yaml:"process_link,omitempty"`
}

type target struct {
	Xpath string `yaml:"xpath,omitempty"`
}

type downloadFileAction struct {
	Target          target
	DirPattern      string `yaml:"dir_pattern"`
	FilenamePattern string `yaml:"filename_pattern"`
	Filters         []filter
}

type processLinkAction struct {
	Target   target
	Final    bool
	NewGroup bool   `yaml:"new_group"`
	Group    *group `yaml:"group,omitempty"`
}

type grabTextAction struct {
	Target          target
	DirPattern      string `yaml:"dir_pattern"`
	FilenamePattern string `yaml:"filename_pattern"`
}

type filter struct {
	filesize string `yaml:"filesize,omitempty"`
}

type postAction struct {
	Join *join `yaml:"join,omitempty"`
	Zip  *zip  `yaml:"zip,omitempty"`
}

type join struct {
	filename       string
	removeOriginal bool   `yaml:"remove_original"`
	joinText       string `yaml:"join_text,omitempty"`
}

type zip struct {
	filename       string
	removeOriginal bool `yaml:"remove_original"`
	includeDir     bool `yaml:"include_dir"`
}

func createConfigurationFromYaml(ymlPath string) (*configuration, error) {
	f, err := os.Open(ymlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &configuration{}

	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(false)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
