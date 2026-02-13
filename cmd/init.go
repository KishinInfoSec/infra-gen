package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kishininfosec/infra-gen/infra-gen/internal/presets"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [preset]",
	Short: "Initialize a new project from a preset",
	Long: `Initialize a new project infrastructure configuration using a preset template.
Available presets include web-app, microservice, database, and more.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		presetID := args[0]
		projectName, _ := cmd.Flags().GetString("name")
		environment, _ := cmd.Flags().GetString("environment")
		outputDir, _ := cmd.Flags().GetString("output")

		if projectName == "" {
			fmt.Println("Error: --name is required")
			os.Exit(1)
		}

		// Create preset manager
		presetManager := presets.NewManager()

		// Check if preset exists
		preset, err := presetManager.GetPreset(presetID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Create project from preset
		config, err := presetManager.CreateProjectFromPreset(presetID, projectName, environment)
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}

		// Create output directory if needed
		if outputDir != "" {
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		// Save project config
		configFile := filepath.Join(outputDir, "infra-gen.yml")
		err = presetManager.SaveProject(config, configFile)
		if err != nil {
			fmt.Printf("Error saving project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project '%s' initialized successfully!\n", projectName)
		fmt.Printf("Configuration saved to: %s\n", configFile)
		fmt.Printf("Preset: %s - %s\n", preset.Name, preset.Description)
		fmt.Printf("Services: %d\n", len(config.Services))
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  infra-gen generate docker\n")
		fmt.Printf("  infra-gen generate ansible\n")
		fmt.Printf("  infra-gen generate terraform\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringP("name", "n", "", "Project name (required)")
	initCmd.Flags().StringP("environment", "e", "development", "Environment (development, staging, production)")
	initCmd.Flags().StringP("output", "o", "", "Output directory (default: current directory)")
}
