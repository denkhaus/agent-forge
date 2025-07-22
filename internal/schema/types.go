// Package schema defines AgentForge component schemas and validation.
package schema

import (
	"fmt"
	"time"
)

// APIVersion represents the AgentForge API version.
const APIVersion = "forge.dev/v1"

// ComponentKind represents the type of AgentForge component.
type ComponentKind string

const (
	KindTool   ComponentKind = "Tool"
	KindPrompt ComponentKind = "Prompt"
	KindAgent  ComponentKind = "Agent"
)

// ValidComponentKinds returns all valid component kinds.
func ValidComponentKinds() []ComponentKind {
	return []ComponentKind{KindTool, KindPrompt, KindAgent}
}

// IsValidKind checks if a component kind is valid.
func IsValidKind(kind string) bool {
	for _, validKind := range ValidComponentKinds() {
		if string(validKind) == kind {
			return true
		}
	}
	return false
}

// Stability represents the stability level of a component.
type Stability string

const (
	StabilityExperimental Stability = "experimental"
	StabilityBeta         Stability = "beta"
	StabilityStable       Stability = "stable"
	StabilityDeprecated   Stability = "deprecated"
)

// Maturity represents the maturity level of a component.
type Maturity string

const (
	MaturityAlpha  Maturity = "alpha"
	MaturityBeta   Maturity = "beta"
	MaturityStable Maturity = "stable"
	MaturityMature Maturity = "mature"
)

// BaseMetadata contains common metadata fields for all component types.
// This follows Kubernetes object metadata patterns.
type BaseMetadata struct {
	// Name is the unique identifier within the repository
	Name string `yaml:"name" json:"name" validate:"required,dns1123label"`
	
	// Namespace for component organization (optional, defaults to "default")
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" validate:"omitempty,dns1123label"`
	
	// Version is the semantic version of the component
	Version string `yaml:"version" json:"version" validate:"required,semver"`
	
	// Description provides a human-readable description
	Description string `yaml:"description" json:"description" validate:"required,min=10,max=500"`
	
	// Author or organization that created the component
	Author string `yaml:"author" json:"author" validate:"required"`
	
	// License using SPDX license identifier
	License string `yaml:"license" json:"license" validate:"required"`
	
	// Homepage URL for the component (optional)
	Homepage string `yaml:"homepage,omitempty" json:"homepage,omitempty" validate:"omitempty,url"`
	
	// Documentation URL (optional)
	Documentation string `yaml:"documentation,omitempty" json:"documentation,omitempty" validate:"omitempty,url"`
	
	// Repository URL for source code (optional)
	Repository string `yaml:"repository,omitempty" json:"repository,omitempty" validate:"omitempty,url"`
	
	// Labels for categorization and search (Kubernetes-style)
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	
	// Annotations for additional metadata (Kubernetes-style)
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	
	// Tags for searchability
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	
	// Categories for hierarchical organization
	Categories []string `yaml:"categories,omitempty" json:"categories,omitempty"`
	
	// Keywords for search optimization
	Keywords []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`
	
	// Stability level of the component
	Stability Stability `yaml:"stability" json:"stability" validate:"required,oneof=experimental beta stable deprecated"`
	
	// Maturity level of the component
	Maturity Maturity `yaml:"maturity" json:"maturity" validate:"required,oneof=alpha beta stable mature"`
	
	// ForgeVersion specifies minimum AgentForge version required
	ForgeVersion string `yaml:"forgeVersion" json:"forgeVersion" validate:"required"`
	
	// Platforms supported by this component
	Platforms []string `yaml:"platforms,omitempty" json:"platforms,omitempty"`
	
	// CreationTimestamp when the component was created (auto-generated)
	CreationTimestamp *time.Time `yaml:"creationTimestamp,omitempty" json:"creationTimestamp,omitempty"`
}

// BaseComponent represents the common structure for all AgentForge components.
// This follows Kubernetes object structure patterns.
type BaseComponent struct {
	// APIVersion specifies the API version
	APIVersion string `yaml:"apiVersion" json:"apiVersion" validate:"required"`
	
	// Kind specifies the component type
	Kind ComponentKind `yaml:"kind" json:"kind" validate:"required"`
	
	// Metadata contains the component metadata
	Metadata BaseMetadata `yaml:"metadata" json:"metadata" validate:"required"`
}

// Validate performs basic validation on the base component.
func (bc *BaseComponent) Validate() error {
	if bc.APIVersion != APIVersion {
		return fmt.Errorf("invalid apiVersion: expected %s, got %s", APIVersion, bc.APIVersion)
	}
	
	if !IsValidKind(string(bc.Kind)) {
		return fmt.Errorf("invalid kind: %s", bc.Kind)
	}
	
	if bc.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}
	
	if bc.Metadata.Version == "" {
		return fmt.Errorf("metadata.version is required")
	}
	
	return nil
}

// GetFullName returns the full name including namespace.
func (bc *BaseComponent) GetFullName() string {
	if bc.Metadata.Namespace != "" {
		return fmt.Sprintf("%s/%s", bc.Metadata.Namespace, bc.Metadata.Name)
	}
	return bc.Metadata.Name
}

// GetLabel returns a label value by key.
func (bc *BaseComponent) GetLabel(key string) string {
	if bc.Metadata.Labels == nil {
		return ""
	}
	return bc.Metadata.Labels[key]
}

// SetLabel sets a label value.
func (bc *BaseComponent) SetLabel(key, value string) {
	if bc.Metadata.Labels == nil {
		bc.Metadata.Labels = make(map[string]string)
	}
	bc.Metadata.Labels[key] = value
}

// GetAnnotation returns an annotation value by key.
func (bc *BaseComponent) GetAnnotation(key string) string {
	if bc.Metadata.Annotations == nil {
		return ""
	}
	return bc.Metadata.Annotations[key]
}

// SetAnnotation sets an annotation value.
func (bc *BaseComponent) SetAnnotation(key, value string) {
	if bc.Metadata.Annotations == nil {
		bc.Metadata.Annotations = make(map[string]string)
	}
	bc.Metadata.Annotations[key] = value
}