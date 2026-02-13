package cmd

import (
	"fmt"

	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/ansible"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/docker"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/generators/terraform"
	"github.com/kishininfosec/infra-gen/infra-gen/internal/presets"
	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate project configuration",
	Long: `Validate the current project configuration for Docker Compose, Ansible, and Terraform
generation. Checks for required fields, service configurations, and potential issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		target, _ := cmd.Flags().GetString("target")

		// Load project configuration
		presetManager := presets.NewManager()
		config, err := presetManager.LoadProject(configFile)
		if err != nil {
			fmt.Printf("❌ Error loading project config: %v\n", err)
			return
		}

		// Validate project configuration
		err = presetManager.ValidateProject(config)
		if err != nil {
			fmt.Printf("❌ Project validation failed:\n%v\n", err)
			return
		}

		fmt.Printf("Project configuration is valid\n")
		fmt.Printf("Project: %s (%s)\n", config.Name, config.Type)
		fmt.Printf("Services: %d\n", len(config.Services))

		// Validate specific targets
		validators := map[string]types.Generator{
			"docker":    docker.NewGenerator(),
			"ansible":   ansible.NewGenerator(),
			"terraform": terraform.NewGenerator(),
		}

		allValid := true
		if target == "all" {
			for name, validator := range validators {
				if err := validator.Validate(config); err != nil {
					fmt.Printf("%s validation failed: %v\n", name, err)
					allValid = false
				} else {
					fmt.Printf("%s configuration is valid\n", name)
				}
			}
		} else if validator, exists := validators[target]; exists {
			if err := validator.Validate(config); err != nil {
				fmt.Printf("%s validation failed: %v\n", target, err)
				allValid = false
			} else {
				fmt.Printf("%s configuration is valid\n", target)
			}
		} else {
			fmt.Printf("Unknown target: %s\n", target)
			fmt.Println("Available targets: all, docker, ansible, terraform")
			return
		}

		if allValid {
			fmt.Println("\nAll validations passed! Ready to generate infrastructure.")
		} else {
			fmt.Println("\nSome validations failed. Please fix issues before generating.")
		}

		// Show warnings and recommendations
		fmt.Println("\nRecommendations:")
		showRecommendations(config)
	},
}

func showRecommendations(config *types.ProjectConfig) {
	// Check for common issues
	for _, service := range config.Services {
		if service.Image == "" && service.Enabled {
			fmt.Printf("  WARNING: Service '%s' is enabled but has no Docker image specified\n", service.Name)
		}

		if len(service.Ports) == 0 && service.Type == "frontend" {
			fmt.Printf("  WARNING: Frontend service '%s' has no ports specified\n", service.Name)
		}

		if len(service.Volumes) == 0 && service.Type == "database" {
			fmt.Printf("  WARNING: Database service '%s' has no persistent volumes\n", service.Name)
		}
	}

	// Check for security concerns
	for key := range config.Variables {
		if containsSensitiveKeywords(key) {
			fmt.Printf("  SECURITY: Variable '%s' contains sensitive data - consider using environment variables\n", key)
		}
	}

	for _, service := range config.Services {
		for key := range service.Environment {
			if containsSensitiveKeywords(key) {
				fmt.Printf("  SECURITY: Service '%s' environment variable '%s' contains sensitive data\n", service.Name, key)
			}
		}
	}
}

func containsSensitiveKeywords(key string) bool {
	key = string([]byte(key))
	key = string([]byte(key))
	sensitiveKeywords := []string{"password", "secret", "key", "token", "auth"}

	for _, keyword := range sensitiveKeywords {
		if len(key) >= len(keyword) {
			for i := 0; i <= len(key)-len(keyword); i++ {
				if key[i:i+len(keyword)] == keyword {
					return true
				}
			}
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Flags
	validateCmd.Flags().StringP("config", "c", "infra-gen.yml", "Project configuration file")
	validateCmd.Flags().StringP("target", "t", "all", "Target to validate (all, docker, ansible, terraform)")
}
