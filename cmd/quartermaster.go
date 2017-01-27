package cmd

import (
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/pschwartz/quartermaster/helpers"
	"github.com/pschwartz/quartermaster/lib"
)

var quartermasterCmd = &cobra.Command{
	Use:   "quartermaster",
	Short: "An irc bot for async standups",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := InitializeConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

var quartermasterCmdV *cobra.Command

var (
	cfgFile    string
	logLevel   string
	buildWatch bool
)

func Execute() {
	quartermasterCmd.SetGlobalNormalizationFunc(helpers.NormalizeQuartermasterFlags)

	quartermasterCmd.SilenceUsage = true

	AddCommands()

	if c, err := quartermasterCmd.ExecuteC(); err != nil {
		if helpers.IsUserError(err) {
			c.Println("")
			c.Println(c.UsageString())
		}

		os.Exit(-1)
	}
}

func AddCommands() {
	quartermasterCmd.AddCommand(botCmd)
	quartermasterCmd.AddCommand(getCmd)
	quartermasterCmd.AddCommand(versionCmd)
}

func init() {
	quartermasterCmdV = quartermasterCmd

	quartermasterCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "", "set logging level (default is Error)")
	quartermasterCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.quartermaster.yaml)")

	// Set bash-completion
	validConfigFilenames := []string{"yaml", "yml"}
	quartermasterCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

func InitializeConfig(subCmdVs ...*cobra.Command) error {
	if err := lib.LoadGlobalConfig(cfgFile); err != nil {
		return err
	}

	for _, cmdV := range append([]*cobra.Command{quartermasterCmdV}, subCmdVs...) {
		initializeFlags(cmdV)
	}

	err := helpers.CreateLogger(logLevel)
	if err != nil {
		return err
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	return nil

}

func initializeFlags(cmd *cobra.Command) {
	persFlagKeys := []string{"verbose", "logFile"}
	flagKeys := []string{}

	for _, key := range persFlagKeys {
		setValueFromFlag(cmd.PersistentFlags(), key)
	}
	for _, key := range flagKeys {
		setValueFromFlag(cmd.Flags(), key)
	}
}

func setValueFromFlag(flags *pflag.FlagSet, key string) {
	if flagChanged(flags, key) {
		f := flags.Lookup(key)
		viper.Set(key, f.Value.String())
	}
}

func flagChanged(flags *pflag.FlagSet, key string) bool {
	flag := flags.Lookup(key)
	if flag == nil {
		return false
	}
	return flag.Changed
}
