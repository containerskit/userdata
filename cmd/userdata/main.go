package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/containerskit/userdata/pkg/userdata"
)

const usage = `
userdata -path /path/to/userdata/file

Read userdata yaml file and replicates its content on the filesystem.

Example:

files:
- mode: 0644
  path: file1
  text: |
    content of file1
- mode: 0755
  type: dir
  path: dir1/file1
  text: |
    content of dir1/file1
- mode: 0400
  path: ssh
  link: /root/.ssh/authorized_keys
  text: |
    ssh public key

will produce the following:

cwd
├── dir1
│   └── file1
├── file1
└── ssh

Arguments:

  -help print this message and exist
`

type args struct {
	path *string
}

var argv = args{
	path: flag.String("path", "", "userdata file path"),
}

func init() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
	}
}

func run() error {
	bin, err := ioutil.ReadFile(*argv.path)
	if err != nil {
		return fmt.Errorf("unable to read config: %v", err)
	}

	cfg := userdata.Config{}
	if err := yaml.Unmarshal(bin, &cfg); err != nil {
		return fmt.Errorf("unable to parse config: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current working directory: %v", err)
	}

	return userdata.Apply(cwd, &cfg)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
