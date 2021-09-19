//HOWTOS:
//$ go get -u github.com/spf13/cobra

package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	version bool
	debug   bool
	host    string

	rootCmd = &cobra.Command{
		Use:     "jcli",
		Short:   "A cli-tool for jocker",
		Long:    `JCli is the reference cli-tool for interacting with jocker-engine`,
		Version: "0.0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// StringVarP(p *string, name, shorthand string, value string, usage string)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "Daemon socket to connect to: tcp://[host]:[port][path] or unix://[/path/to/socket]")
	rootCmd.AddCommand(newContainerCommand())
	//rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")

	//rootCmd.AddCommand(addCmd)
	//rootCmd.AddCommand(initCmd)
}

func initConfig() {
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := os.UserHomeDir()
	//	cobra.CheckErr(err)

	//	// Search config in home directory with name ".cobra" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigType("yaml")
	//	viper.SetConfigName(".cobra")
	//}

	//viper.AutomaticEnv()

	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
}
