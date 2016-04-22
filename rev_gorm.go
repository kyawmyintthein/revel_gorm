// The command line tool for running Revel apps.
package main

import (
	"flag"
	"fmt"
	"github.com/agtorre/gocolorize"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"text/template"
	"time"
)

// Cribbed from the genius organization of the "go" command.
type Command struct {
	Run func(cmd *Command, args []string)
	UsageLine, Short, Long string
	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

var commands = []*Command{
	cmdDBSetup,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if runtime.GOOS == "windows" {
		gocolorize.SetPlain(true)
	}
	fmt.Fprintf(os.Stdout, gocolorize.NewColor("blue").Paint(header))
	flag.Usage = func() { usage(1) }
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 || args[0] == "help" {
		if len(args) == 1 {
			usage(0)
		}
		if len(args) > 1 {
			for _, cmd := range commands {
				if cmd.Name() == args[1] {
					tmpl(os.Stdout, helpTemplate, cmd)
					return
				}
			}
		}
		usage(2)
	}

	// Commands use panic to abort execution when something goes wrong.
	// Panics are logged at the point of error.  Ignore those.
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(LoggedError); !ok {
				// This panic was not expected / logged.
				panic(err)
			}
			os.Exit(1)
		}
	}()

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Run(cmd, args[1:])
			return
		}
	}

	ColorLog("[ERRO] unknown command %q\nRun 'revel help' for usage.\n", args[0])
	os.Exit(2)
}

func usage(exitCode int) {
	tmpl(os.Stderr, usageTemplate, commands)
	os.Exit(exitCode)
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}



// templates
const header = `~
~ revel_gorm! https://github.com/kyawmyintthein/revel_gorm
~
`


const usageTemplate = `usage: revel_gorm command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "revel_gorm help [command]" for more information.
`


var helpTemplate = `usage: revel_gorm {{.UsageLine}}
{{.Long}}
`

