package gotk

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

type Command struct {
	App         string       `json:"app"`
	Project     *viper.Viper `json:"project"`
	Subcommands []Subcommand `json:"subcommands"`
}

func NewCommand(app string) (command *Command) {
	return &Command{
		App:         app,
		Project:     viper.New(),
		Subcommands: make([]Subcommand, 0),
	}
}

func (self *Command) ProjectFromBts(bts []byte) (err error) {
	var meta map[string]any

	// _Project.ReadConfig(strings.NewReader(str))
	if err = self.Project.ReadConfig(bytes.NewReader(bts)); err != nil {
		return err
	}

	meta = BuildInfo()
	meta["app_name"] = self.Project.GetString("app_name")
	meta["app_version"] = self.Project.GetString("app_version")
	self.Project.Set("meta", meta)

	return nil
}

type Subcommand struct {
	Name string         `json:"name"`
	Help string         `json:"help"`
	Run  func([]string) `json:"-"`
}

func (self *Command) Execute(args []string) {
	var (
		cmd        string
		subcommand *Subcommand
	)

	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		cmd = "help"
	} else {
		cmd = args[0]
	}

	if cmd == "help" {
		self.Usage()
		os.Exit(0)
	} else if subcommand = self.Find(cmd); subcommand != nil {
		subcommand.Run(args[1:])
	} else {
		self.Usage()
		os.Exit(1)
	}
}

func (self *Command) AddCmd(name, help string, run func([]string)) *Command {
	self.Subcommands = append(self.Subcommands, Subcommand{
		Name: name,
		Help: help,
		Run:  run,
	})

	return self
}

func (self *Command) Find(name string) *Subcommand {
	for i := range self.Subcommands {
		if self.Subcommands[i].Name == name {
			return &self.Subcommands[i]
		}
	}

	return nil
}

func (self *Command) Usage() {
	var (
		text  string
		templ *template.Template
	)

	text = `usage:
- {{.App}} [command]

commands: {{range .Subcommands}}
- {{.Name}}: {{.Help}}
{{end}}`

	templ, _ = template.New("usage").Parse(text)
	_ = templ.Execute(os.Stderr, self)

	fmt.Printf("\nmeta:\n%s\n", BuildInfoText(self.Project.GetStringMap("meta")))
}
