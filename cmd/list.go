package cmd

import (
	"fmt"
	"strings"

	"github.com/kishininfosec/infra-gen/infra-gen/internal/presets"
	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List available presets and project information",
	Long: `List available project presets, categories, or current project information.
Use 'presets' to see all available presets, 'categories' to see preset categories,
or 'project' to see current project details.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listType := "presets"
		if len(args) > 0 {
			listType = args[0]
		}

		presetManager := presets.NewManager()

		switch listType {
		case "presets":
			listPresets(presetManager)
		case "categories":
			listCategories(presetManager)
		case "project":
			listProject(cmd)
		default:
			fmt.Printf("Unknown list type: %s\n", listType)
			fmt.Println("Available types: presets, categories, project")
		}
	},
}

func listPresets(presetManager *presets.Manager) {
	presets := presetManager.ListPresets()

	if len(presets) == 0 {
		fmt.Println("No presets available")
		return
	}

	fmt.Println("Available Project Presets:")
	fmt.Println(strings.Repeat("=", 50))

	categories := make(map[string][]types.Preset)
	for _, preset := range presets {
		categories[preset.Category] = append(categories[preset.Category], preset)
	}

	for category, categoryPresets := range categories {
		fmt.Printf("\n%s:\n", category)
		for _, preset := range categoryPresets {
			fmt.Printf("  %-12s - %s\n", preset.ID, preset.Description)
			if len(preset.Tags) > 0 {
				fmt.Printf("      Tags: %s\n", strings.Join(preset.Tags, ", "))
			}
		}
	}

	fmt.Printf("\nTotal: %d presets across %d categories\n", len(presets), len(categories))
}

func listCategories(presetManager *presets.Manager) {
	presets := presetManager.ListPresets()

	if len(presets) == 0 {
		fmt.Println("No presets available")
		return
	}

	categories := make(map[string]int)
	for _, preset := range presets {
		categories[preset.Category]++
	}

	fmt.Println("Available Categories:")
	fmt.Println(strings.Repeat("=", 30))

	for category, count := range categories {
		fmt.Printf("  %-20s (%d presets)\n", category, count)
	}

	fmt.Printf("\nTotal: %d categories\n", len(categories))
}

func listProject(cmd *cobra.Command) {
	configFile, _ := cmd.Flags().GetString("config")

	presetManager := presets.NewManager()
	config, err := presetManager.LoadProject(configFile)
	if err != nil {
		fmt.Printf("Error loading project: %v\n", err)
		return
	}

	fmt.Printf("Project Information:\n")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Name:        %s\n", config.Name)
	fmt.Printf("Type:        %s\n", config.Type)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Version:     %s\n", config.Version)
	fmt.Printf("Created:     %s\n", config.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", config.UpdatedAt.Format("2006-01-02 15:04:05"))

	if len(config.Services) > 0 {
		fmt.Printf("\nServices (%d):\n", len(config.Services))
		for _, service := range config.Services {
			status := "Enabled"
			if !service.Enabled {
				status = "Disabled"
			}
			fmt.Printf("  %s %-15s (%s)\n", status, service.Name, service.Type)
		}
	}

	if len(config.Variables) > 0 {
		fmt.Printf("\nVariables (%d):\n", len(config.Variables))
		for key, value := range config.Variables {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Flags
	listCmd.Flags().StringP("config", "c", "infra-gen.yml", "Project configuration file (for 'list project')")
}
