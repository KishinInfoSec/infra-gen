package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/kishininfosec/infra-gen/infra-gen/pkg/types"
)

//go:embed embedded/*.yaml
var embeddedTemplates embed.FS

// TemplateManager manages embedded and custom templates
type TemplateManager struct {
	embedded map[string]types.Preset
	custom   map[string]types.Preset
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		embedded: make(map[string]types.Preset),
		custom:   make(map[string]types.Preset),
	}

	// Load embedded templates
	if err := tm.loadEmbeddedTemplates(); err != nil {
		// Log error but continue with empty templates
		fmt.Printf("Warning: Failed to load embedded templates: %v\n", err)
	}

	return tm
}

// loadEmbeddedTemplates loads templates from embedded filesystem
func (tm *TemplateManager) loadEmbeddedTemplates() error {
	return fs.WalkDir(embeddedTemplates, "embedded", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		// Parse YAML content to get preset
		// For now, we'll create simple presets
		preset := types.Preset{
			ID:          strings.TrimSuffix(d.Name(), ".yaml"),
			Name:        tm.generatePresetName(d.Name()),
			Description: tm.generatePresetDescription(d.Name()),
			Category:    tm.getCategoryFromPath(path),
		}

		tm.embedded[preset.ID] = preset
		return nil
	})
}

// GetPreset returns a preset by ID
func (tm *TemplateManager) GetPreset(id string) (*types.Preset, error) {
	if preset, exists := tm.embedded[id]; exists {
		return &preset, nil
	}

	if preset, exists := tm.custom[id]; exists {
		return &preset, nil
	}

	return nil, fmt.Errorf("preset '%s' not found", id)
}

// ListPresets returns all available presets
func (tm *TemplateManager) ListPresets() []types.Preset {
	var presets []types.Preset

	for _, preset := range tm.embedded {
		presets = append(presets, preset)
	}

	for _, preset := range tm.custom {
		presets = append(presets, preset)
	}

	return presets
}

// ListPresetsByCategory returns presets filtered by category
func (tm *TemplateManager) ListPresetsByCategory(category string) []types.Preset {
	var presets []types.Preset

	for _, preset := range tm.embedded {
		if preset.Category == category {
			presets = append(presets, preset)
		}
	}

	for _, preset := range tm.custom {
		if preset.Category == category {
			presets = append(presets, preset)
		}
	}

	return presets
}

// LoadCustomTemplate loads a custom template from file
func (tm *TemplateManager) LoadCustomTemplate(path string) error {
	// This would load templates from filesystem
	// For now, return not implemented
	return fmt.Errorf("custom template loading not implemented yet")
}

// Helper functions
func (tm *TemplateManager) generatePresetName(filename string) string {
	name := strings.TrimSuffix(filename, ".yaml")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Capitalize first letter
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}

func (tm *TemplateManager) generatePresetDescription(filename string) string {
	base := strings.TrimSuffix(filename, ".yaml")

	descriptions := map[string]string{
		"web-app":        "Basic web application with frontend, backend, and database",
		"microservice":   "Microservice architecture with API gateway and services",
		"database":       "Database services (PostgreSQL, MySQL, MongoDB)",
		"ml":             "Machine learning pipeline with model serving",
		"infrastructure": "Core infrastructure components (monitoring, logging)",
	}

	if desc, exists := descriptions[base]; exists {
		return desc
	}

	return fmt.Sprintf("%s project template", base)
}

func (tm *TemplateManager) getCategoryFromPath(path string) string {
	// Extract category from path structure
	// For now, return generic category
	if strings.Contains(path, "web") {
		return "Web Applications"
	} else if strings.Contains(path, "database") {
		return "Databases"
	} else if strings.Contains(path, "ml") {
		return "Machine Learning"
	} else if strings.Contains(path, "infrastructure") {
		return "Infrastructure"
	}

	return "General"
}
