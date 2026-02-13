package docker

import (
	"fmt"
	"strings"

	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
)

// Generator implements Docker Compose generation
type Generator struct{}

// NewGenerator creates a new Docker Compose generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GetTarget returns the target type
func (g *Generator) GetTarget() types.Target {
	return types.TargetDocker
}

// Generate generates Docker Compose files from project config
func (g *Generator) Generate(config *types.ProjectConfig) ([]types.GeneratedFile, error) {
	if err := g.Validate(config); err != nil {
		return nil, err
	}

	// Generate docker-compose.yml
	yamlContent, err := g.generateComposeYAML(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate compose YAML: %w", err)
	}

	// Generate .env file if there are environment variables
	envContent := g.generateEnvFile(config)

	files := []types.GeneratedFile{
		{
			Path:     "docker-compose.yml",
			Content:  yamlContent,
			Type:     types.TargetDocker,
			Encoding: "utf-8",
		},
	}

	if envContent != "" {
		files = append(files, types.GeneratedFile{
			Path:     ".env",
			Content:  envContent,
			Type:     types.TargetDocker,
			Encoding: "utf-8",
		})
	}

	return files, nil
}

// Validate validates the project config for Docker Compose generation
func (g *Generator) Validate(config *types.ProjectConfig) error {
	var errors types.ValidationErrors

	if config.Name == "" {
		errors.Add("name", "project name is required", config.Name)
	}

	if len(config.Services) == 0 {
		errors.Add("services", "at least one service is required", len(config.Services))
	}

	// Validate each service
	for i, service := range config.Services {
		if service.Name == "" {
			errors.Add(fmt.Sprintf("services[%d].name", i), "service name is required", service.Name)
		}
		if service.Type == "" {
			errors.Add(fmt.Sprintf("services[%d].type", i), "service type is required", service.Type)
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// generateComposeYAML generates the YAML content for docker-compose.yml
func (g *Generator) generateComposeYAML(config *types.ProjectConfig) (string, error) {
	var builder strings.Builder

	builder.WriteString("services:\n")

	// Generate services
	for _, service := range config.Services {
		if !service.Enabled {
			continue
		}

		builder.WriteString("  ")
		builder.WriteString(service.Name)
		builder.WriteString(":\n")

		if service.Image != "" {
			builder.WriteString("    image: ")
			builder.WriteString(service.Image)
			builder.WriteString("\n")
		}

		// Generate ports
		if len(service.Ports) > 0 {
			builder.WriteString("    ports:\n")
			for _, port := range service.Ports {
				builder.WriteString("      - \"")
				if port.Host > 0 {
					builder.WriteString(fmt.Sprintf("%d", port.Host))
					builder.WriteString(":")
				}
				builder.WriteString(fmt.Sprintf("%d", port.Container))
				builder.WriteString("\"")
				if port.Protocol != "" {
					builder.WriteString(" # ")
					builder.WriteString(port.Protocol)
				}
				builder.WriteString("\n")
			}
		}

		// Generate volumes
		if len(service.Volumes) > 0 {
			builder.WriteString("    volumes:\n")
			for _, volume := range service.Volumes {
				builder.WriteString("      - ")
				builder.WriteString(volume.Source)
				builder.WriteString(":")
				builder.WriteString(volume.Target)
				if volume.ReadOnly {
					builder.WriteString(":ro")
				}
				builder.WriteString("\n")
			}
		}

		// Generate environment variables
		if len(service.Environment) > 0 {
			builder.WriteString("    environment:\n")
			for key, value := range service.Environment {
				builder.WriteString("      ")
				builder.WriteString(key)
				builder.WriteString(": ")
				builder.WriteString(value)
				builder.WriteString("\n")
			}
		}

		// Generate dependencies
		if len(service.DependsOn) > 0 {
			builder.WriteString("    depends_on:\n")
			for _, dep := range service.DependsOn {
				builder.WriteString("      - ")
				builder.WriteString(dep)
				builder.WriteString("\n")
			}
		}

		builder.WriteString("\n")
	}

	// Generate volumes
	hasVolumes := false
	for _, service := range config.Services {
		if len(service.Volumes) > 0 {
			hasVolumes = true
			break
		}
	}

	if hasVolumes {
		builder.WriteString("volumes:\n")
		for _, service := range config.Services {
			for _, volume := range service.Volumes {
				if volume.Type == "volume" {
					builder.WriteString("  ")
					builder.WriteString(volume.Source)
					builder.WriteString(":\n")
				}
			}
		}
	}

	return builder.String(), nil
}

// generateEnvFile generates .env file content
func (g *Generator) generateEnvFile(config *types.ProjectConfig) string {
	var envVars []string

	// Add project-level variables
	for key, value := range config.Variables {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Add service-level environment variables that should be external
	for _, service := range config.Services {
		for key, value := range service.Environment {
			// Only include variables that look like they should be external
			if strings.Contains(strings.ToUpper(key), "PASSWORD") ||
				strings.Contains(strings.ToUpper(key), "SECRET") ||
				strings.Contains(strings.ToUpper(key), "KEY") ||
				strings.Contains(strings.ToUpper(key), "TOKEN") {
				envVars = append(envVars, fmt.Sprintf("%s_%s=%s", strings.ToUpper(service.Name), key, value))
			}
		}
	}

	if len(envVars) == 0 {
		return ""
	}

	return strings.Join(envVars, "\n") + "\n"
}
