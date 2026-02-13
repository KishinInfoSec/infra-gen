package types

import "time"

// ProjectType represents different types of projects that can be generated
type ProjectType string

const (
	ProjectTypeWebApp       ProjectType = "web-app"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeDatabase     ProjectType = "database"
	ProjectTypeML           ProjectType = "ml"
	ProjectTypeInfrastructure ProjectType = "infrastructure"
)

// Target represents the infrastructure target
type Target string

const (
	TargetDocker    Target = "docker"
	TargetAnsible   Target = "ansible"
	TargetTerraform Target = "terraform"
)

// ProjectConfig holds the configuration for a project
type ProjectConfig struct {
	Name        string            `yaml:"name"`
	Type        ProjectType       `yaml:"type"`
	Description string            `yaml:"description,omitempty"`
	Version     string            `yaml:"version,omitempty"`
	Environment string            `yaml:"environment,omitempty"`
	Services    []ServiceConfig   `yaml:"services"`
	Variables   map[string]string `yaml:"variables,omitempty"`
	CreatedAt   time.Time         `yaml:"created_at"`
	UpdatedAt   time.Time         `yaml:"updated_at"`
}

// ServiceConfig represents a single service in the project
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"`
	Image       string            `yaml:"image,omitempty"`
	Ports       []PortConfig      `yaml:"ports,omitempty"`
	Volumes     []VolumeConfig    `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Enabled     bool              `yaml:"enabled"`
}

// PortConfig represents a port mapping
type PortConfig struct {
	Host      int    `yaml:"host,omitempty"`
	Container int    `yaml:"container"`
	Protocol  string `yaml:"protocol,omitempty"`
}

// VolumeConfig represents a volume mapping
type VolumeConfig struct {
	Source   string `yaml:"source"`
	Target   string `yaml:"target"`
	ReadOnly bool   `yaml:"read_only,omitempty"`
	Type     string `yaml:"type,omitempty"`
}

// Generator interface for different infrastructure generators
type Generator interface {
	Generate(config *ProjectConfig) ([]GeneratedFile, error)
	Validate(config *ProjectConfig) error
	GetTarget() Target
}

// GeneratedFile represents a generated infrastructure file
type GeneratedFile struct {
	Path     string `yaml:"path"`
	Content  string `yaml:"content"`
	Type     Target `yaml:"type"`
	Encoding string `yaml:"encoding,omitempty"`
}

// Preset represents a project preset
type Preset struct {
	ID          string              `yaml:"id"`
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Category    string              `yaml:"category"`
	Services    []PresetService     `yaml:"services"`
	Variables   map[string]string   `yaml:"variables,omitempty"`
	Tags        []string            `yaml:"tags,omitempty"`
}

// PresetService represents a service in a preset
type PresetService struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"`
	Description string            `yaml:"description"`
	Image       string            `yaml:"image,omitempty"`
	Ports       []PortConfig      `yaml:"ports,omitempty"`
	Volumes     []VolumeConfig    `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Optional    bool              `yaml:"optional,omitempty"`
}