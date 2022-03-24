/*
Copyright © 2022 Prince Merluza <prince.merluza@gmail.com>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/PrinceMerluza/devcenter-content-linter/blueprintrepo"
	"github.com/PrinceMerluza/devcenter-content-linter/config"
	"github.com/PrinceMerluza/devcenter-content-linter/linter"
	"github.com/PrinceMerluza/devcenter-content-linter/logger"
	"github.com/PrinceMerluza/devcenter-content-linter/transform_data"
	"github.com/PrinceMerluza/devcenter-content-linter/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	isRemoteRepo bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gc-linter repo-path --config config.json",
	Short: "Valdiates content for the Genesys Cloud Developer Center",
	Long: `The gc-linter is a CLI tool which validates the structure, format, and required files 
of different Genesys Cloud developer center content. 

Examples of this content are: blueprints.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initViperConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		blueprintrepo.UseRepo(repoPath, isRemoteRepo)

		results := validateContent(blueprintrepo.GetWorkingPath())
		resultsJsonB, err := json.Marshal(results)
		if err != nil {
			logger.Fatal(err)
		}

		utils.Render(string(resultsJsonB))
	},
	Args: cobra.ExactArgs(1),
}

func validateContent(repoPath string) *linter.ValidationResult {
	validationData := &linter.ValidationData{
		ContentPath: repoPath,
		RuleData:    config.LoadedRuleSet,
	}

	result, err := validationData.Validate()
	if err != nil {
		logger.Fatal(err)
	}

	return result
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Initizalize the viper config
// Viper config is a required flag
func initViperConfig() error {
	if cfgFile == "" {
		return errors.New("config file is required")
	}
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Error reading config file: ", err)
		return err
	}

	logger.Info("Using config file: ", viper.ConfigFileUsed())

	// Set the config data
	if err := viper.Unmarshal(&config.LoadedRuleSet); err != nil {
		logger.Fatal(err)
	}

	return nil
}

func init() {
	cobra.OnInitialize()

	// Flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file that defines the type of content")
	rootCmd.MarkFlagRequired("config")

	rootCmd.PersistentFlags().BoolVarP(&logger.LoggingEnabled, "enable-logging", "l", false, "enable logging")
	rootCmd.PersistentFlags().BoolVarP(&isRemoteRepo, "remote", "r", false, "if the repo-path is an HTTP URL")

	rootCmd.PersistentFlags().StringVarP(&transform_data.TemplateFile, "transform", "t", "", "Provide a Go template file for transforming output data")

	logger.InitLogger()
}
