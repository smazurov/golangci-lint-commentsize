package commentsize

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("commentsize", newPlugin)
}

type settings struct {
	MaxLines int `json:"max-lines"`
}

type plugin struct {
	maxLines int
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[settings](conf)
	if err != nil {
		return nil, err
	}
	return &plugin{maxLines: s.MaxLines}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{New(p.maxLines)}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
