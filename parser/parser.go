package parser

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Parser fields' definitions and the regular expression pattern to extract them
type Parser struct {
	Pattern string `yaml:"pattern"`
	Fields  []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"fields"`
}

// New init a Parser
func New(fname string) (*Parser, error) {
	var p Parser
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, &p)
	return &p, err
}

// Parse analyze line
func (p *Parser) Parse(line string) ([]float64, error) {
	var values []float64
	var err error

	re := regexp.MustCompile(p.Pattern)
	results := re.FindStringSubmatch(line)
	if len(results) < 1 || len(results)-1 != len(p.Fields) {
		err = errors.New("The pattern does not define the same num. of fields as required")
		return nil, err
	}

	for _, part := range results[1:] {
		v, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}
	return values, nil
}
