package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Starters []starter
	Rules    []rule
}

type starter struct {
	Targets    []string
	StartGroup *startGroupAction `yaml:"start_group"`
}

type rule struct {
	Targets []target
	Action  *action
}

type target struct {
	URL   string `yaml:"url,omitempty"`
	Xpath string `yaml:"xpath,omitempty"`
}

type action struct {
	DownloadFile *downloadFileAction `yaml:"download_file,omitempty"`
	StartGroup   *startGroupAction   `yaml:"start_group,omitempty"`
	ProcessLink  *processLinkAction  `yaml:"process_link,omitempty"`
}

type downloadFileAction struct {
	Target          target
	DirPattern      string `yaml:"dir_pattern"`
	FilenamePattern string `yaml:"filename_pattern"`
}

type startGroupAction struct {
	DirPattern   string `yaml:"dir_pattern"`
	GroupName    string `yaml:"group_name"`
	ResetCounter bool   `yaml:"reset_counter"`
}

type processLinkAction struct {
	Target target
	Final  bool
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
