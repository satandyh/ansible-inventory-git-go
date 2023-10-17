package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	config "github.com/satandyh/ansible-inventory-git-go/internal/config"
	logging "github.com/satandyh/ansible-inventory-git-go/internal/logger"
)

// Global vars for logs
var logConfig = logging.LogConfig{
	ConsoleLoggingEnabled: true,
	EncodeLogsAsJson:      true,
	FileLoggingEnabled:    false,
	Directory:             "./data",
	Filename:              "ans-inv-git.log",
	MaxSize:               10,
	MaxBackups:            7,
	MaxAge:                7,
	LogLevel:              6,
}

var logger = logging.Configure(logConfig)

func main() {

	conf := config.NewConfig()

	// temp dir for download git inventory
	workDir, tmp_err := os.MkdirTemp("", "")
	if tmp_err != nil {
		logger.Fatal().
			Err(tmp_err).
			Str("module", "main").
			Msg("Can't create temp folder.")
		os.Exit(1)
	}
	// remove directory before exit
	defer os.RemoveAll(workDir)

	//clone repo to temp directory
	cl_err := gitClone(conf.Git.Repo_ssh_address, conf.Git.Key_path, conf.Git.Branch, workDir)
	if cl_err != nil {
		logger.Fatal().
			Err(cl_err).
			Str("module", "main").
			Msg("Can't clone repo.")
		os.Exit(1)
	}

	// exec ansible
	targetPath := filepath.Join(workDir, conf.Git.Target)

	cmd := exec.Command("ansible-inventory")
	cmd.Dir = workDir
	cmd.Env = append(cmd.Environ(), "ANSIBLE_INVENTORY_ENABLED=host_list,auto,yaml,ini,toml")

	if len(conf.Host) == 0 { // we have some "host" arg
		cmd.Args = append(cmd.Args, "--list", "-i", targetPath)

	} else { // no any "host" arg - just list all
		cmd.Args = append(cmd.Args, "--host", conf.Host, "-i", targetPath)
	}

	output, out_err := cmd.Output()
	if out_err != nil {
		logger.Fatal().
			Err(out_err).
			Str("module", "main").
			Msg("Can't execute ansible-inventory.")
		os.Exit(1)
	}

	// to stdout
	fmt.Println(string(output))

}
