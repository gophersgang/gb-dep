package subcommands

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Command interface {
	Run(args []string, log *log.Logger)
	Usage() string
}

type Subcommands struct {
	name, desc string
	cmds       map[string]info
}

func New(name, desc string) *Subcommands {
	c := &Subcommands{
		name: name,
		desc: desc,
		cmds: make(map[string]info),
	}

	c.Register("help", "Describe a subcommand", help{c})

	return c
}

func (c *Subcommands) Run(args []string, log *log.Logger) {
	if len(args) == 0 {
		args = []string{"help"}
	}

	for cmd, f := range c.cmds {
		if args[0] == cmd {
			f.cmd.Run(args[1:], log)
			return
		}
	}

	if os.Args[1][0] == '-' {
		log.Fatalf("flag provided but not defined '%s'", os.Args[1])
	}

	log.Fatalf("Unknown subcommand '%s'", os.Args[1])
}

func (c *Subcommands) Register(name, desc string, cmd Command) {
	c.cmds[name] = info{desc: desc, cmd: cmd}
}

type info struct {
	desc string
	cmd  Command
}

type help struct {
	c *Subcommands
}

func (h help) Run(args []string, log *log.Logger) {
	if len(args) == 1 {
		c, ok := h.c.cmds[args[0]]
		if !ok {
			log.Fatalf("unknown help topic '%s'", args[0])
		}

		usage := c.cmd.Usage()
		if usage == "" {
			log.Fatalf("Unknown help topic '%s'", args[0])
		}
		log.Fatal(usage)
	}

	var commandsDesc string
	width := 0
	for cmd, _ := range h.c.cmds {
		if len(cmd) > width {
			width = len(cmd)
		}
	}

	for cmd, _ := range h.c.cmds {
		commandsDesc = fmt.Sprintf("%s\n    %s%s%s",
			commandsDesc,
			cmd,
			strings.Repeat(" ", width+4-len(cmd)),
			h.c.cmds[cmd].desc)
	}
	log.Fatalf("%s: %s\n\nCommands:%s", h.c.name, h.c.desc, commandsDesc)
}

func (h help) Usage() string {
	return "help describes a help topic"
}
