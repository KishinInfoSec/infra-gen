package presets

import (
	"fmt"
	"os"
	"time"

	"github.com/kishininfosec/infra-gen/infra-gen/internal/templates"
	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
	"gopkg.in/yaml.v3"
)

// Manager manages project presets
type Manager struct {
	templateManager *templates.TemplateManager
}

// NewManager creates a new preset manager
func NewManager() *Manager {
	return &Manager{
		templateManager: templates.NewTemplateManager(),
	}
}

// CreateProjectFromPreset creates a project configuration from a preset
func (m *Manager) CreateProjectFromPreset(presetID string, projectName string, environment string) (*types.ProjectConfig, error) {
	preset, err := m.templateManager.GetPreset(presetID)
	if err != nil {
		return nil, err
	}

	config := &types.ProjectConfig{
		Name:        projectName,
		Type:        m.getProjectTypeFromPreset(presetID),
		Description: preset.Description,
		Environment: environment,
		Version:     "1.0.0",
		Services:    make([]types.ServiceConfig, 0),
		Variables:   make(map[string]string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Get preset with service data
	presetWithServices, err := m.loadPresetTemplateData(presetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load preset template data: %w", err)
	}

	// Convert preset services to project services
	for _, presetService := range presetWithServices.Services {
		service := types.ServiceConfig{
			Name:        presetService.Name,
			Type:        presetService.Type,
			Image:       presetService.Image,
			Ports:       presetService.Ports,
			Volumes:     presetService.Volumes,
			Environment: make(map[string]string),
			DependsOn:   []string{},
			Enabled:     !presetService.Optional,
		}

		// Copy environment variables
		for key, value := range presetService.Environment {
			service.Environment[key] = value
		}

		config.Services = append(config.Services, service)
	}

	// Copy preset variables
	for key, value := range preset.Variables {
		config.Variables[key] = value
	}

	return config, nil
}

// ListPresets returns all available presets
func (m *Manager) ListPresets() []types.Preset {
	return m.templateManager.ListPresets()
}

// ListPresetsByCategory returns presets filtered by category
func (m *Manager) ListPresetsByCategory(category string) []types.Preset {
	return m.templateManager.ListPresetsByCategory(category)
}

// GetPreset returns a specific preset
func (m *Manager) GetPreset(id string) (*types.Preset, error) {
	return m.templateManager.GetPreset(id)
}

// ValidateProject validates a project configuration
func (m *Manager) ValidateProject(config *types.ProjectConfig) error {
	var errors types.ValidationErrors

	if config.Name == "" {
		errors.Add("name", "project name is required", config.Name)
	}

	if config.Type == "" {
		errors.Add("type", "project type is required", config.Type)
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

	// Check for duplicate service names
	serviceNames := make(map[string]bool)
	for _, service := range config.Services {
		if serviceNames[service.Name] {
			errors.Add("services", fmt.Sprintf("duplicate service name: %s", service.Name), service.Name)
		}
		serviceNames[service.Name] = true
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// SaveProject saves a project configuration to a file
func (m *Manager) SaveProject(config *types.ProjectConfig, filePath string) error {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	return os.WriteFile(filePath, yamlData, 0644)
}

// LoadProject loads a project configuration from a file
func (m *Manager) LoadProject(filePath string) (*types.ProjectConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}

	var config types.ProjectConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal project config: %w", err)
	}

	return &config, nil
}

// loadPresetTemplateData loads the detailed template data for a preset
func (m *Manager) loadPresetTemplateData(presetID string) (*types.Preset, error) {
	// For now, return basic preset data based on ID
	switch presetID {
	case "web-app":
		return &types.Preset{
			Services: []types.PresetService{
				{
					Name:        "frontend",
					Type:        "frontend",
					Description: "Frontend web application",
					Image:       "nginx:alpine",
					Ports:       []types.PortConfig{{Container: 80, Protocol: "tcp"}},
					Environment: map[string]string{"REACT_APP_API_URL": "http://api:8080"},
					Optional:    false,
				},
				{
					Name:        "api",
					Type:        "api",
					Description: "Backend API server",
					Image:       "node:18-alpine",
					Ports:       []types.PortConfig{{Container: 8080, Protocol: "tcp"}},
					Environment: map[string]string{"NODE_ENV": "development", "DB_HOST": "database"},
					Optional:    false,
				},
				{
					Name:        "database",
					Type:        "database",
					Description: "PostgreSQL database",
					Image:       "postgres:15",
					Ports:       []types.PortConfig{{Container: 5432, Protocol: "tcp"}},
					Volumes:     []types.VolumeConfig{{Source: "db_data", Target: "/var/lib/postgresql/data", Type: "volume"}},
					Environment: map[string]string{"POSTGRES_DB": "webapp", "POSTGRES_USER": "admin"},
					Optional:    false,
				},
			},
		}, nil
	case "microservice":
		return &types.Preset{
			Services: []types.PresetService{
				{
					Name:        "api-gateway",
					Type:        "gateway",
					Description: "API Gateway",
					Image:       "nginx:alpine",
					Ports:       []types.PortConfig{{Container: 80, Protocol: "tcp"}},
					Environment: map[string]string{"UPSTREAM_SERVICE1": "service1:8081"},
					Optional:    false,
				},
				{
					Name:        "service1",
					Type:        "microservice",
					Description: "First microservice",
					Image:       "node:18-alpine",
					Ports:       []types.PortConfig{{Container: 8081, Protocol: "tcp"}},
					Environment: map[string]string{"SERVICE_NAME": "service1"},
					Optional:    false,
				},
				{
					Name:        "redis",
					Type:        "cache",
					Description: "Redis cache",
					Image:       "redis:7-alpine",
					Ports:       []types.PortConfig{{Container: 6379, Protocol: "tcp"}},
					Volumes:     []types.VolumeConfig{{Source: "redis_data", Target: "/data", Type: "volume"}},
					Optional:    false,
				},
			},
		}, nil
	case "database":
		return &types.Preset{
			Services: []types.PresetService{
				{
					Name:        "postgres",
					Type:        "database",
					Description: "PostgreSQL database",
					Image:       "postgres:15",
					Ports:       []types.PortConfig{{Container: 5432, Protocol: "tcp"}},
					Volumes:     []types.VolumeConfig{{Source: "postgres_data", Target: "/var/lib/postgresql/data", Type: "volume"}},
					Environment: map[string]string{"POSTGRES_DB": "myapp", "POSTGRES_USER": "postgres"},
					Optional:    false,
				},
				{
					Name:        "mysql",
					Type:        "database",
					Description: "MySQL database",
					Image:       "mysql:8.0",
					Ports:       []types.PortConfig{{Container: 3306, Protocol: "tcp"}},
					Volumes:     []types.VolumeConfig{{Source: "mysql_data", Target: "/var/lib/mysql", Type: "volume"}},
					Environment: map[string]string{"MYSQL_DATABASE": "myapp", "MYSQL_USER": "mysql"},
					Optional:    true,
				},
			},
		}, nil
	default:
		return &types.Preset{Services: []types.PresetService{}}, nil
	}
}

// Helper functions
func (m *Manager) getProjectTypeFromPreset(presetID string) types.ProjectType {
	switch presetID {
	case "web-app":
		return types.ProjectTypeWebApp
	case "microservice":
		return types.ProjectTypeMicroservice
	case "database":
		return types.ProjectTypeDatabase
	case "ml":
		return types.ProjectTypeML
	case "infrastructure":
		return types.ProjectTypeInfrastructure
	default:
		return types.ProjectTypeWebApp // Default fallback
	}
}
