// got from
// https://github.com/RobinUS2/indispenso/blob/master/conf.go

package config

import (
	"fmt"
	"os"

	logging "github.com/satandyh/ansible-inventory-git-go/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Conf struct {
	Host    string
	List    bool
	Example bool
	Git     struct {
		Repo_ssh_address string
		Key_path         string
		Branch           string
		Target           string
	}
	confFlags *pflag.FlagSet
}

// logger
var logConfig = logging.LogConfig{
	ConsoleLoggingEnabled: true,
	EncodeLogsAsJson:      true,
	FileLoggingEnabled:    false,
	Directory:             ".",
	Filename:              "log.log",
	MaxSize:               10,
	MaxBackups:            1,
	MaxAge:                1,
	LogLevel:              6,
}

var log = logging.Configure(logConfig)

var genConfig = `## EXAMPLE CONFIG
# repo with inventory file
repo_ssh_address: ssh://git@git.example.com:22/ansible/repo/inventory.git
# absolute path to private ssh key (should be accessible)
key_path: /example/dir/with/ssh_key/key_name
# branch name of repo with inventory file
branch: master
# relative path to inventory directory or inventory file (inventory.yaml) - it will be used with ansible command
target: inventory`

var standaloneUsage = `Standalone Usage: ans-inv-git_darwin [-c CONFIG_FILE | --config[=]CONFIG_FILE] [--host HOST] [--list] [-g | --generate-config] [-h | --help]

Options:`

var ansibleUsage = `
Usage with Ansible:
  1. Check your ansible.cfg file: script statement should present in enable_plugins option.
  2. Place app and it's config (named the same as app file but with .yaml) somewhere and remember path.
  3. Use next command to check that all works:
    ansible -i /some/folder/ans-inv-git lovely_host -m ping
  4. Use ansible as you always do:
    ansible-playbook -i /some/folder/ans-inv-git --diff plays/lovely_play.yml -l lovely_host`

func NewConfig() Conf {
	var c Conf

	//default config path
	def_conf, _ := os.Executable()
	def_conf = def_conf + ".yaml"

	// all env will look like AIG_SOMETHING form Ansible-Inventory-Git
	// for embedded use AIG_LEV1.VALUE
	viper.SetEnvPrefix("aig")

	// Defaults
	viper.SetDefault("Git.Branch", "master")
	viper.SetDefault("Git.Target", "inventory")

	//Flags
	c.confFlags = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	configFile := c.confFlags.StringP("config", "c", "", "Set the absolute path to the config file.\nBy default it use work directory and application filename with \".yaml\" at the end.\nExample, if application filename is \"ans-git-inv\", then default config filename\nwill be \"ans-git-inv.yaml\" in the same directory.")
	c.confFlags.String("host", "", "Output specific host info.")
	c.confFlags.Bool("list", true, "Output all hosts info.")
	generate := c.confFlags.BoolP("generate-config", "g", false, "Generate example config to stdout.")
	help := c.confFlags.BoolP("help", "h", false, "Print help message.")

	//parse flags
	arg_err := c.confFlags.Parse(os.Args[1:])
	if arg_err != nil {
		log.Fatal().
			Err(arg_err).
			Str("module", "config").
			Msg("Can't parse cmd args.")
		os.Exit(1)
	}
	if *help {
		fmt.Println(standaloneUsage)
		c.confFlags.PrintDefaults()
		fmt.Println(ansibleUsage)
		os.Exit(0)
	}
	if *generate {
		fmt.Println(genConfig)
		os.Exit(0)
	}

	if len(*configFile) > 2 {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(def_conf)
		//viper.SetConfigName("config")             // name of config file (without extension)
		//viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
		//viper.AddConfigPath("/opt/ans-inv-git")   // path to look for the config file in
		//viper.AddConfigPath("$HOME/.ans-inv-git") // call multiple times to add many search paths
		//viper.AddConfigPath("./config")
		//viper.AddConfigPath(".")
	}

	// bind flags from pflags
	arg_bind_err := viper.BindPFlags(c.confFlags)
	if arg_bind_err != nil {
		log.Fatal().
			Err(arg_bind_err).
			Str("module", "config").
			Msg("Internal viper binding error.")
		os.Exit(1)
	}

	// try to get values from env
	viper.AutomaticEnv()

	// get values from config
	file_read_err := viper.ReadInConfig()
	if file_read_err != nil {
		log.Fatal().
			Err(file_read_err).
			Str("module", "config").
			Msg("Can't get values from config file.")
		os.Exit(1)
	}

	// do all above and get our values
	dec_err := viper.Unmarshal(&c)
	if dec_err != nil {
		log.Fatal().
			Err(dec_err).
			Str("module", "config").
			Msg("Internal viper unmarshal error.")
		os.Exit(1)
	}

	dec_git_err := viper.Unmarshal(&c.Git)
	if dec_git_err != nil {
		log.Fatal().
			Err(dec_git_err).
			Str("module", "config").
			Msg("Another internal viper unmarshal error.")
		os.Exit(1)
	}

	return c
}
