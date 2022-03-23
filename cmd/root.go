/*
Copyright © 2022 Prince Merluza <prince.merluza@gmail.com>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/PrinceMerluza/devcenter-content-linter/config"
	"github.com/PrinceMerluza/devcenter-content-linter/linter"
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

		// If the repoPath is a URL, clone the repo
		if isRemoteRepo {
			tmpPath, err := utils.CloneRepoTemp(repoPath)
			if err != nil {
				log.Fatal(err)
			}
			repoPath = tmpPath
		}

		results := validateContent(repoPath)
		printResults(results)
		ExportJsonResult(results, "result.json")
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
		log.Fatal(err)
	}

	return result
}

func printResults(finalResult *linter.ValidationResult) {
	for _, result := range *finalResult.SuccessResults {
		fmt.Printf("\n--- SUCCESS --\n")

		fmt.Printf("%s \nLevel: %s \nDescription: %s",
			result.Id, result.Level, result.Description)

		if result.FileHighlights != nil {
			for _, fileHighlight := range *result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n\t%v \n", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
		fmt.Println()
	}

	for _, result := range *finalResult.FailureResults {
		fmt.Printf("\n--- FAILED --\n")

		fmt.Printf("%s \nLevel: %s \nDescription: %s",
			result.Id, result.Level, result.Description)

		if result.FileHighlights != nil {
			for _, fileHighlight := range *result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n%v \n", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
		fmt.Println()
	}
}

func ExportJsonResult(finalResult *linter.ValidationResult, filename string) error {
	data, err := json.Marshal(finalResult)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
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
		fmt.Println("Error reading config file:", err)
		return err
	}

	fmt.Println("Using config file: ", viper.ConfigFileUsed())

	// Set the config data
	if err := viper.Unmarshal(&config.LoadedRuleSet); err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	cobra.OnInitialize()

	// Flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file that defines the type of content")
	rootCmd.MarkFlagRequired("config")

	rootCmd.PersistentFlags().BoolVarP(&isRemoteRepo, "remote", "r", false, "if the repo-path is an HTTP URL")
}
