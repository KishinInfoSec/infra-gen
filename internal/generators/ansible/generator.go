package ansible

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/infra-gen/infra-gen/pkg/types"
)

// AnsiblePlaybook represents an Ansible playbook structure
type AnsiblePlaybook struct {
	Hosts       string                   `yaml:"hosts"`
	Name        string                   `yaml:"name"`
	Become      bool                     `yaml:"become,omitempty"`
	Vars        map[string]interface{}   `yaml:"vars,omitempty"`
	Environment map[string]interface{}   `yaml:"environment,omitempty"`
	Tasks       []AnsibleTask            `yaml:"tasks"`
}

// AnsibleTask represents an Ansible task
type AnsibleTask struct {
	Name    string                 `yaml:"name"`
	Module  string                 `yaml:"module,omitempty"`
	Package string                 `yaml:"name,omitempty"`
	State   string                 `yaml:"state,omitempty"`
	WithItems []interface{}       `yaml:"with_items,omitempty"`
	Vars    map[string]interface{} `yaml:"vars,omitempty"`
}

// AnsibleInventory represents an Ansible inventory structure
type AnsibleInventory struct {
	All struct {
		Children map[string]struct {
			Hosts map[string]map[string]interface{} `yaml:"hosts"`
		} `yaml:"children"`
		Vars map[string]interface{} `yaml:"vars,omitempty"`
	} `yaml:"all"`
}

// Generator implements Ansible playbook generation
type Generator struct{}

// NewGenerator creates a new Ansible generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GetTarget returns the target type
func (g *Generator) GetTarget() types.Target {
	return types.TargetAnsible
}

// Generate generates Ansible files from project config
func (g *Generator) Generate(config *types.ProjectConfig) ([]types.GeneratedFile, error) {
	if err := g.Validate(config); err != nil {
		return nil, err
	}

	files := []types.GeneratedFile{}

	// Generate main playbook
	playbookContent, err := g.generatePlaybook(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate playbook: %w", err)
	}

	files = append(files, types.GeneratedFile{
		Path:     "playbook.yml",
		Content:  playbookContent,
		Type:     types.TargetAnsible,
		Encoding: "utf-8",
	})

	// Generate inventory
	inventoryContent, err := g.generateInventory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate inventory: %w", err)
	}

	files = append(files, types.GeneratedFile{
		Path:     "inventory.yml",
		Content:  inventoryContent,
		Type:     types.TargetAnsible,
		Encoding: "utf-8",
	})

	// Generate requirements.yml if needed
	requirementsContent := g.generateRequirements(config)
	if requirementsContent != "" {
		files = append(files, types.GeneratedFile{
			Path:     "requirements.yml",
			Content:  requirementsContent,
			Type:     types.TargetAnsible,
			Encoding: "utf-8",
		})
	}

	return files, nil
}

// Validate validates the project config for Ansible generation
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

// generatePlaybook generates the main Ansible playbook
func (g *Generator) generatePlaybook(config *types.ProjectConfig) (string, error) {
	var builder strings.Builder
	
	builder.WriteString("---\n")
	builder.WriteString("- hosts: all\n")
	builder.WriteString("  become: true\n")
	builder.WriteString("  name: Deploy ")
	builder.WriteString(config.Name)
	builder.WriteString("\n")
	
	// Generate vars
	vars := g.generateVars(config)
	if len(vars) > 0 {
		builder.WriteString("  vars:\n")
		for key, value := range vars {
			builder.WriteString("    ")
			builder.WriteString(key)
			builder.WriteString(": ")
			builder.WriteString(fmt.Sprintf("%v", value))
			builder.WriteString("\n")
		}
	}
	
	// Generate tasks
	tasks := g.generateTasks(config)
	builder.WriteString("  tasks:\n")
	for _, task := range tasks {
		builder.WriteString("    - name: ")
		builder.WriteString(task.Name)
		builder.WriteString("\n")
		if task.Module != "" {
			builder.WriteString("      ")
			builder.WriteString(task.Module)
			builder.WriteString(":\n")
			if task.Package != "" {
				builder.WriteString("        name: ")
				builder.WriteString(task.Package)
				builder.WriteString("\n")
			}
			if task.State != "" {
				builder.WriteString("        state: ")
				builder.WriteString(task.State)
				builder.WriteString("\n")
			}
		}
	}

	return builder.String(), nil
}

// generateVars generates variables for the playbook
func (g *Generator) generateVars(config *types.ProjectConfig) map[string]interface{} {
	vars := make(map[string]interface{})
	
	// Add project-level variables
	for key, value := range config.Variables {
		vars[key] = value
	}

	// Add service variables
	for _, service := range config.Services {
		prefix := strings.ReplaceAll(service.Name, "-", "_")
		vars[fmt.Sprintf("%s_image", prefix)] = service.Image
		vars[fmt.Sprintf("%s_ports", prefix)] = g.extractPorts(service.Ports)
		vars[fmt.Sprintf("%s_volumes", prefix)] = g.extractVolumes(service.Volumes)
		vars[fmt.Sprintf("%s_enabled", prefix)] = service.Enabled
	}

	return vars
}

// generateTasks generates tasks for the playbook
func (g *Generator) generateTasks(config *types.ProjectConfig) []AnsibleTask {
	var tasks []AnsibleTask

	// Add common setup tasks
	tasks = append(tasks, AnsibleTask{
		Name: "Update package cache",
		Module: "apt",
		Package: "apt",
		State:   "present",
	})

	tasks = append(tasks, AnsibleTask{
		Name: "Install Docker",
		Module: "package",
		Package: "docker.io",
		State:   "present",
	})

	tasks = append(tasks, AnsibleTask{
		Name: "Start and enable Docker service",
		Module: "systemd",
		Package: "docker",
		State:   "started",
	})

	// Add service-specific tasks
	for _, service := range config.Services {
		if !service.Enabled {
			continue
		}

		// Create directories for volumes
		for _, volume := range service.Volumes {
			if volume.Source != "" && !strings.HasPrefix(volume.Source, "/") {
				tasks = append(tasks, AnsibleTask{
					Name:    fmt.Sprintf("Create directory for %s volume", volume.Source),
					Module:  "file",
					Package: volume.Source,
					State:   "directory",
				})
			}
		}

		// Pull Docker image
		if service.Image != "" {
			tasks = append(tasks, AnsibleTask{
				Name:    fmt.Sprintf("Pull %s Docker image", service.Name),
				Module:  "docker_image",
				Package: service.Image,
				State:   "present",
			})
		}
	}

	// Create docker-compose file
	tasks = append(tasks, AnsibleTask{
		Name:    "Create docker-compose.yml",
		Module:  "file",
		Package: "/opt/{{ project_name }}/docker-compose.yml",
		State:   "directory",
	})

	tasks = append(tasks, AnsibleTask{
		Name:    "Deploy services with Docker Compose",
		Module:  "docker_compose",
		Package: "/opt/{{ project_name }}",
		State:   "present",
	})

	return tasks
}

// generateInventory generates the Ansible inventory
func (g *Generator) generateInventory(config *types.ProjectConfig) (string, error) {
	inventory := AnsibleInventory{}

	// Create default groups
	inventory.All.Children = make(map[string]struct {
		Hosts map[string]map[string]interface{} `yaml:"hosts"`
	})

	// Add web servers group
	if g.hasServiceType(config.Services, "web", "frontend", "nginx") {
		inventory.All.Children["webservers"] = struct {
			Hosts map[string]map[string]interface{} `yaml:"hosts"`
		}{
			Hosts: map[string]map[string]interface{}{
				"webserver1": {
					"ansible_host": "{{ webserver_ip | default('127.0.0.1') }}",
					"ansible_user": "{{ ansible_user | default('ubuntu') }}",
				},
			},
		}
	}

	// Add database servers group
	if g.hasServiceType(config.Services, "database", "postgres", "mysql", "mongo") {
		inventory.All.Children["databases"] = struct {
			Hosts map[string]map[string]interface{} `yaml:"hosts"`
		}{
			Hosts: map[string]map[string]interface{}{
				"database1": {
					"ansible_host": "{{ database_ip | default('127.0.0.1') }}",
					"ansible_user": "{{ ansible_user | default('ubuntu') }}",
				},
			},
		}
	}

	// Add global variables
	inventory.All.Vars = map[string]interface{}{
		"project_name": config.Name,
		"environment":  config.Environment,
	}

	// Add project variables
	for key, value := range config.Variables {
		inventory.All.Vars[key] = value
	}

	yamlData, err := yaml.Marshal(inventory)
	if err != nil {
		return "", err
	}

	return string(yamlData), nil
}

// generateRequirements generates Ansible requirements if needed
func (g *Generator) generateRequirements(config *types.ProjectConfig) string {
	// For now, return empty. Can be extended to include collections/roles
	return ""
}

// Helper functions
func (g *Generator) extractPorts(ports []types.PortConfig) []string {
	var portStrings []string
	for _, port := range ports {
		if port.Host > 0 {
			portStrings = append(portStrings, fmt.Sprintf("%d:%d", port.Host, port.Container))
		} else {
			portStrings = append(portStrings, fmt.Sprintf("%d", port.Container))
		}
	}
	return portStrings
}

func (g *Generator) extractVolumes(volumes []types.VolumeConfig) []string {
	var volumeStrings []string
	for _, volume := range volumes {
		volumeStrings = append(volumeStrings, fmt.Sprintf("%s:%s", volume.Source, volume.Target))
	}
	return volumeStrings
}

func (g *Generator) hasServiceType(services []types.ServiceConfig, types ...string) bool {
	for _, service := range services {
		for _, t := range types {
			if strings.Contains(strings.ToLower(service.Type), strings.ToLower(t)) {
				return true
			}
		}
	}
	return false
}