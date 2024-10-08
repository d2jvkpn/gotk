package gotk

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

type Command struct {
	App         string       `json:"app"`
	Project     *viper.Viper `json:"project"`
	Subcommands []Subcommand `json:"subcommands"`
}

func NewCommand(app string, projects ...*viper.Viper) (command *Command) {
	command = &Command{
		App:         app,
		Subcommands: make([]Subcommand, 0),
	}

	if len(projects) == 0 {
		command.Project = viper.New()
	} else {
		command.Project = projects[0]
	}

	return command
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

func (self *Command) UsageTemplate() {
	var (
		text  string
		templ *template.Template
	)

	text = `usage:
- {{.App}} [command]

commands: {{range .Subcommands}}
- {{.Name}}: {{.Help}}{{end}}
`

	templ, _ = template.New("usage").Parse(text)
	_ = templ.Execute(os.Stderr, self)

	if meta := self.Project.GetStringMap("meta"); len(meta) > 0 {
		fmt.Printf("\nmeta:\n%s\n", BuildInfoText(meta, "  "))
	}
}

func (self *Command) Usage() {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("usage:\n  %s: [command]\n\n", self.App))

	builder.WriteString("commands:\n")
	for _, v := range self.Subcommands {
		builder.WriteString(fmt.Sprintf("  %s: %s\n", v.Name, v.Help))
	}
	builder.WriteString("\n")

	if meta := self.Project.GetStringMap("meta"); len(meta) > 0 {
		builder.WriteString("meta:\n")
		builder.WriteString(BuildInfoText(meta, "  "))
	}

	fmt.Println(builder.String())
}

func (self *Command) UpdateMeta(mp map[string]any) {
	meta := self.Project.GetStringMap("meta")

	for k, v := range mp {
		meta[k] = v
	}
}
