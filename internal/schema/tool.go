package schema

import "fmt"

// ToolType represents the type of tool implementation.
type ToolType string

const (
	ToolTypeMCPServer ToolType = "mcp-server"
	ToolTypeFunction  ToolType = "function"
	ToolTypeWebhook   ToolType = "webhook"
	ToolTypeGRPC      ToolType = "grpc"
	ToolTypeHTTP      ToolType = "http"
)

// ToolRuntime represents the runtime environment for the tool.
type ToolRuntime string

const (
	RuntimeGo     ToolRuntime = "go"
	RuntimePython ToolRuntime = "python"
	RuntimeNode   ToolRuntime = "node"
	RuntimeRust   ToolRuntime = "rust"
	RuntimeJava   ToolRuntime = "java"
	RuntimeDocker ToolRuntime = "docker"
)

// ToolParameter represents a parameter definition for a tool.
type ToolParameter struct {
	Name        string      `yaml:"name" json:"name" validate:"required"`
	Type        string      `yaml:"type" json:"type" validate:"required,oneof=string number boolean object array"`
	Description string      `yaml:"description" json:"description" validate:"required"`
	Required    bool        `yaml:"required" json:"required"`
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []string    `yaml:"enum,omitempty" json:"enum,omitempty"`
	Pattern     string      `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MinLength   *int        `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength   *int        `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Minimum     *float64    `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum     *float64    `yaml:"maximum,omitempty" json:"maximum,omitempty"`
}

// ToolFunction represents a function exposed by the tool.
type ToolFunction struct {
	Name        string          `yaml:"name" json:"name" validate:"required"`
	Description string          `yaml:"description" json:"description" validate:"required"`
	Parameters  []ToolParameter `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Returns     *ToolParameter  `yaml:"returns,omitempty" json:"returns,omitempty"`
	Examples    []ToolExample   `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// ToolExample represents an example usage of a tool function.
type ToolExample struct {
	Name        string                 `yaml:"name" json:"name" validate:"required"`
	Description string                 `yaml:"description" json:"description"`
	Input       map[string]interface{} `yaml:"input" json:"input"`
	Output      interface{}            `yaml:"output" json:"output"`
}

// ToolDependency represents a dependency required by the tool.
type ToolDependency struct {
	Name    string `yaml:"name" json:"name" validate:"required"`
	Version string `yaml:"version" json:"version"`
	Type    string `yaml:"type" json:"type" validate:"oneof=system package library service"`
	Source  string `yaml:"source,omitempty" json:"source,omitempty"`
}

// ToolConfiguration represents configuration options for the tool.
type ToolConfiguration struct {
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Files       []string          `yaml:"files,omitempty" json:"files,omitempty"`
	Secrets     []string          `yaml:"secrets,omitempty" json:"secrets,omitempty"`
}

// ToolSpec defines the specification for Tool components.
type ToolSpec struct {
	// Type of tool implementation
	Type ToolType `yaml:"type" json:"type" validate:"required"`
	
	// Runtime environment
	Runtime ToolRuntime `yaml:"runtime" json:"runtime" validate:"required"`
	
	// Entry point for the tool
	EntryPoint string `yaml:"entryPoint" json:"entryPoint" validate:"required"`
	
	// Command line arguments
	Args []string `yaml:"args,omitempty" json:"args,omitempty"`
	
	// Working directory
	WorkingDir string `yaml:"workingDir,omitempty" json:"workingDir,omitempty"`
	
	// Functions exposed by this tool
	Functions []ToolFunction `yaml:"functions" json:"functions" validate:"required,min=1"`
	
	// Dependencies required by the tool
	Dependencies []ToolDependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	
	// Configuration options
	Configuration ToolConfiguration `yaml:"configuration,omitempty" json:"configuration,omitempty"`
	
	// Health check configuration
	HealthCheck *ToolHealthCheck `yaml:"healthCheck,omitempty" json:"healthCheck,omitempty"`
	
	// Resource requirements
	Resources *ToolResources `yaml:"resources,omitempty" json:"resources,omitempty"`
	
	// Security context
	Security *ToolSecurity `yaml:"security,omitempty" json:"security,omitempty"`
}

// ToolHealthCheck defines health check configuration.
type ToolHealthCheck struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	Path            string `yaml:"path,omitempty" json:"path,omitempty"`
	IntervalSeconds int    `yaml:"intervalSeconds,omitempty" json:"intervalSeconds,omitempty"`
	TimeoutSeconds  int    `yaml:"timeoutSeconds,omitempty" json:"timeoutSeconds,omitempty"`
	FailureThreshold int   `yaml:"failureThreshold,omitempty" json:"failureThreshold,omitempty"`
}

// ToolResources defines resource requirements and limits.
type ToolResources struct {
	Requests *ToolResourceSpec `yaml:"requests,omitempty" json:"requests,omitempty"`
	Limits   *ToolResourceSpec `yaml:"limits,omitempty" json:"limits,omitempty"`
}

// ToolResourceSpec defines specific resource specifications.
type ToolResourceSpec struct {
	CPU    string `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty" json:"memory,omitempty"`
	Disk   string `yaml:"disk,omitempty" json:"disk,omitempty"`
}

// ToolSecurity defines security context for the tool.
type ToolSecurity struct {
	RunAsUser         *int64   `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	RunAsGroup        *int64   `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`
	ReadOnlyRootFS    bool     `yaml:"readOnlyRootFS,omitempty" json:"readOnlyRootFS,omitempty"`
	AllowedCapabilities []string `yaml:"allowedCapabilities,omitempty" json:"allowedCapabilities,omitempty"`
	DroppedCapabilities []string `yaml:"droppedCapabilities,omitempty" json:"droppedCapabilities,omitempty"`
}

// Tool represents a complete Tool component manifest.
type Tool struct {
	BaseComponent `yaml:",inline"`
	Spec          ToolSpec `yaml:"spec" json:"spec" validate:"required"`
}

// NewTool creates a new Tool with default values.
func NewTool(name, version string) *Tool {
	return &Tool{
		BaseComponent: BaseComponent{
			APIVersion: APIVersion,
			Kind:       KindTool,
			Metadata: BaseMetadata{
				Name:      name,
				Version:   version,
				Stability: StabilityExperimental,
				Maturity:  MaturityAlpha,
			},
		},
	}
}

// Validate performs Tool-specific validation.
func (t *Tool) Validate() error {
	if err := t.BaseComponent.Validate(); err != nil {
		return err
	}
	
	if t.Kind != KindTool {
		return fmt.Errorf("invalid kind for Tool: %s", t.Kind)
	}
	
	if len(t.Spec.Functions) == 0 {
		return fmt.Errorf("tool must define at least one function")
	}
	
	// Validate function names are unique
	functionNames := make(map[string]bool)
	for _, fn := range t.Spec.Functions {
		if functionNames[fn.Name] {
			return fmt.Errorf("duplicate function name: %s", fn.Name)
		}
		functionNames[fn.Name] = true
	}
	
	return nil
}