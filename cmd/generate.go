package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/ansible"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/docker"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/terraform"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/presets"
	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [target]",
	Short: "Generate infrastructure configurations",
	Long: `Generate infrastructure configurations for Docker Compose, Ansible, or Terraform
based on the current project configuration. Use 'all' to generate all targets.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		outputDir, _ := cmd.Flags().GetString("output")
		target := "all"

		if len(args) > 0 {
			target = args[0]
		}

		// Load project configuration
		presetManager := presets.NewManager()
		config, err := presetManager.LoadProject(configFile)
		if err != nil {
			fmt.Printf("Error loading project config: %v\n", err)
			os.Exit(1)
		}

		// Validate project config
		err = presetManager.ValidateProject(config)
		if err != nil {
			fmt.Printf("Validation error: %v\n", err)
			os.Exit(1)
		}

		// Create output directory
		if outputDir != "" {
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		// Generate configurations
		generatedFiles := 0
		targets := []types.Target{}

		switch target {
		case "all":
			targets = []types.Target{types.TargetDocker, types.TargetAnsible, types.TargetTerraform}
		case "docker":
			targets = []types.Target{types.TargetDocker}
		case "ansible":
			targets = []types.Target{types.TargetAnsible}
		case "terraform":
			targets = []types.Target{types.TargetTerraform}
		default:
			fmt.Printf("Unknown target: %s\n", target)
			os.Exit(1)
		}

		for _, t := range targets {
			var generator types.Generator

			switch t {
			case types.TargetDocker:
				generator = docker.NewGenerator()
			case types.TargetAnsible:
				generator = ansible.NewGenerator()
			case types.TargetTerraform:
				generator = terraform.NewGenerator()
			}

			files, err := generator.Generate(config)
			if err != nil {
				fmt.Printf("Error generating %s: %v\n", t, err)
				continue
			}

			// Write files
			for _, file := range files {
				filePath := filepath.Join(outputDir, file.Path)

				// Create directory if needed
				dir := filepath.Dir(filePath)
				if dir != "." {
					err = os.MkdirAll(dir, 0755)
					if err != nil {
						fmt.Printf("Error creating directory %s: %v\n", dir, err)
						continue
					}
				}

				err = os.WriteFile(filePath, []byte(file.Content), 0644)
				if err != nil {
					fmt.Printf("Error writing file %s: %v\n", filePath, err)
					continue
				}

				fmt.Printf("Generated: %s\n", filePath)
				generatedFiles++
			}
		}

		if generatedFiles > 0 {
			fmt.Printf("\nGenerated %d files for project '%s'\n", generatedFiles, config.Name)
		} else {
			fmt.Printf("No files generated\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Flags
	generateCmd.Flags().StringP("config", "c", "infra-gen.yml", "Project configuration file")
	generateCmd.Flags().StringP("output", "o", "", "Output directory (default: current directory)")
}
