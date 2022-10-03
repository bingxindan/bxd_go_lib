package config

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type YamlReader struct {
	path string
}

func NewYamlReader(path string) (FileReader, error) {
	return &YamlReader{path}, nil
}

func (s *YamlReader) Read(ctx context.Context) (map[string]string, error) {
	config := map[string]string{}

	source, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, errors.WithMessage(err, "ReadFile yaml file failed")
	}

	content := map[string]interface{}{}
	err = yaml.Unmarshal(source, &content)
	if err != nil {
		return nil, errors.WithMessage(err, "Unmarshal config yaml file failed")
	}

	Walk(ctx, "", content, &config)

	return config, nil
}

func (s *YamlReader) Close() error {
	return nil
}
