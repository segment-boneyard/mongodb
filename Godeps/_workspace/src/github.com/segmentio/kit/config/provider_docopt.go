package config

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/segmentio/kit/schema"
	"github.com/tj/docopt"
)

const docOptTemplate = `
	Usage: {{.ServiceName}}
		{{range .Options}}[{{.Command}} {{.CommandKey}}]
		{{end}}

	Options:
		-h --help	Show help information.
		--version	Show version information.
		{{range .Options}}{{.Command}} {{.Description}}
		{{end}}
`

type docOptProvider struct {
	args map[string]interface{}
}

type doctOptConfig struct {
	ServiceName string
	Options     []docOptOption
}

type docOptOption struct {
	Command     string
	CommandKey  string
	Description string
}

func (d *docOptProvider) Setup(service schema.Service) error {
	var err error
	usage, err := d.generateUsage(service)
	if err != nil {
		return err
	}
	d.args, err = docopt.Parse(usage, nil, true, service.Version, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *docOptProvider) transformKey(key string) string {
	return fmt.Sprintf("--%s", strings.ToLower(strings.Replace(key, ".", "-", -1)))
}

func (d *docOptProvider) Get(val schema.ConfigValue) interface{} {
	return d.args[d.transformKey(val.Key)]
}

// generateUsage uses a template to generate
// the configuration usage for docopt
func (d *docOptProvider) generateUsage(service schema.Service) (string, error) {
	var out bytes.Buffer
	c := doctOptConfig{
		ServiceName: service.Name,
		Options:     []docOptOption{},
	}

	for _, v := range service.Config {
		if v.Description == "" {
			v.Description = "<Description Missing>"
		}
		c.Options = append(c.Options, docOptOption{
			Command:     d.transformKey(v.Key),
			Description: v.Description,
			CommandKey:  v.Key,
		})
	}

	tmpl, err := template.New("docopt").Parse(docOptTemplate)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&out, c); err != nil {
		return "", err
	}

	return out.String(), nil
}
