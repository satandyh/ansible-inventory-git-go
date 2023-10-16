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

// Task represents a task in the database
type Task struct {
	ID      int
	Command string
	State   string
}

// Trigger represents a trigger condition in the database
type Trigger struct {
	ID     int
	TaskID int
	Time   string // Assuming cron-like time description as a string
	State  string
}

type Conf struct {
	Branch string
	Target string
	Host string
	List bool
	Example bool
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

var genConfig = `
## EXAMPLE CONFIG
# repo with inventory file
repo_ssh_address: ssh://git@git.example.com:22/ansible/repo/inventory.git
# absolute path to private ssh key (should be accessible)
key_path: /example/dir/with/ssh_key/key_name
# branch name of repo with inventory file
branch: master
# relative path to inventory directory or inventory file (inventory.yaml) - it will be used with ansible command
target: inventory
`

func NewConfig() Conf {
	var c Conf

	// all env will look like AIG_SOMETHING form Ansible-Inventory-Git
	// for embedded use AIG_LEV1.VALUE
	viper.SetEnvPrefix("aig")

	// Defaults
	viper.SetDefault("Repo_ssh_address", "ssh://git@git.example.com:22/ansible/repo/inventory.git")
	viper.SetDefault("Key_path", "/example/dir/with/ssh_key/key_name")
	viper.SetDefault("Branch", "master")
	viper.SetDefault("Target", "inventory")

	//Flags
	c.confFlags = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	configFile := c.confFlags.StringP("config", "c", "", "Config file location. Supported formats {json,toml,yaml}. Default path {'$HOME/.ans-inv-git','.','./config','/opt/ans-inv-git'}/config.yaml")
	c.confFlags.StringP("host", "", "", "Output specific host info.")
	c.confFlags.StringP("list", "", "", "Output all hosts info.")
	//c.confFlags.BoolP("example", "e", false, "Insert example data into database. Default - false")
	generate := c.confFlags.BoolP("generate-config", "g", false, "Generate example config to stdout. Default - false")
	help := c.confFlags.BoolP("help", "h", false, "Print help message")

	//parse flags
	arg_err := c.confFlags.Parse(os.Args[1:])
	if arg_err != nil {
		log.Fatal().
			Err(arg_err).
			Str("module", "config").
			Msg("")
	}
	if *help {
		fmt.Println("Usage of ans-inv-git:")
		c.confFlags.PrintDefaults()
		os.Exit(0)
	}
	if *generate {
		fmt.Println(genConfig)
		os.Exit(0)
	}

	if len(*configFile) > 2 {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigName("config.yml")    // name of config file (without extension)
		viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/opt/ans-inv-git")   // path to look for the config file in
		viper.AddConfigPath("$HOME/.ans-inv-git") // call multiple times to add many search paths
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
	}

	// bind flags from pflags
	arg_bind_err := viper.BindPFlags(c.confFlags)
	if arg_bind_err != nil {
		log.Fatal().
			Err(arg_bind_err).
			Str("module", "config").
			Msg("")
	}

	// try to get values from env
	viper.AutomaticEnv()

	// get values from config
	file_read_err := viper.ReadInConfig()
	if file_read_err != nil {
		log.Fatal().
			Err(file_read_err).
			Str("module", "config").
			Msg("")
		//os.Exit(0)
	}

	// do all above and get our values
	dec_err := viper.Unmarshal(&c)
	if dec_err != nil {
		log.Fatal().
			Err(dec_err).
			Str("module", "config").
			Msg("")
	}

	return c
}
