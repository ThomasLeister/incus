package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxd/shared/logging"
	"github.com/lxc/lxd/shared/version"
)

// Global variables
var debug bool
var verbose bool

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type cmdGlobal struct {
	flagHelp    bool
	flagVersion bool

	flagLogFile    string
	flagLogDebug   bool
	flagLogSyslog  bool
	flagLogTrace   []string
	flagLogVerbose bool
}

func (c *cmdGlobal) Run(cmd *cobra.Command, args []string) error {
	// Set logging global variables
	debug = c.flagLogVerbose
	verbose = c.flagLogDebug

	// Setup logger
	syslog := ""
	if c.flagLogSyslog {
		syslog = "lxd"
	}

	log, err := logging.GetLogger(syslog, c.flagLogFile, c.flagLogVerbose, c.flagLogDebug, nil)
	if err != nil {
		return err
	}
	logger.Log = log

	return nil
}

func main() {
	// daemon command (main)
	daemonCmd := cmdDaemon{}
	app := daemonCmd.Command()
	app.SilenceUsage = true

	// Workaround for main command
	app.Args = cobra.ArbitraryArgs

	// Global flags
	globalCmd := cmdGlobal{}
	daemonCmd.global = &globalCmd
	app.PersistentPreRunE = globalCmd.Run
	app.PersistentFlags().BoolVar(&globalCmd.flagVersion, "version", false, "Print version number")
	app.PersistentFlags().BoolVarP(&globalCmd.flagHelp, "help", "h", false, "Print help")
	app.PersistentFlags().StringVar(&globalCmd.flagLogFile, "logfile", "", "Path to the log file"+"``")
	app.PersistentFlags().StringArrayVar(&globalCmd.flagLogTrace, "trace", []string{}, "Log tracing targets"+"``")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogDebug, "debug", "d", false, "Show all debug messages")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogVerbose, "verbose", "v", false, "Show all information messages")

	// Version handling
	app.SetVersionTemplate("{{.Version}}\n")
	app.Version = version.Version

	// activateifneeded sub-command
	activateifneededCmd := cmdActivateifneeded{global: &globalCmd}
	app.AddCommand(activateifneededCmd.Command())

	// callhook sub-command
	callhookCmd := cmdCallhook{global: &globalCmd}
	app.AddCommand(callhookCmd.Command())

	// import sub-command
	importCmd := cmdImport{global: &globalCmd}
	app.AddCommand(importCmd.Command())

	// migratedumpsuccess sub-command
	migratedumpsuccessCmd := cmdMigratedumpsuccess{global: &globalCmd}
	app.AddCommand(migratedumpsuccessCmd.Command())

	// netcat sub-command
	netcatCmd := cmdNetcat{global: &globalCmd}
	app.AddCommand(netcatCmd.Command())

	// shutdown sub-command
	shutdownCmd := cmdShutdown{global: &globalCmd}
	app.AddCommand(shutdownCmd.Command())

	// sql sub-command
	sqlCmd := cmdSql{global: &globalCmd}
	app.AddCommand(sqlCmd.Command())

	// waitready sub-command
	waitreadyCmd := cmdWaitready{global: &globalCmd}
	app.AddCommand(waitreadyCmd.Command())

	// Run the main command and handle errors
	err := app.Execute()
	if err != nil {
		os.Exit(1)
	}
}
