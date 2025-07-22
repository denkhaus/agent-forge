package schema

import (
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/denkhaus/agentforge/internal/types"
)

// componentParser is a private implementation of types.ComponentParser interface.
type componentParser struct{}

// NewComponentParser creates a new component parser.
func NewComponentParser() types.ComponentParser {
	return &componentParser{}
}

// ParseComponent parses a component manifest from YAML content.
func (p *componentParser) ParseComponent(content []byte) (interface{}, error) {
	log.Info("Parsing component manifest")
	
	// First, parse just the base structure to determine the kind
	var base BaseComponent
	if err := yaml.Unmarshal(content, &base); err != nil {
		return nil, fmt.Errorf("failed to parse base component: %w", err)
	}
	
	log.Info("Detected component kind", zap.String("kind", string(base.Kind)))
	
	// Parse the specific component type based on kind
	switch base.Kind {
	case KindTool:
		var tool Tool
		if err := yaml.Unmarshal(content, &tool); err != nil {
			return nil, fmt.Errorf("failed to parse tool component: %w", err)
		}
		if err := tool.Validate(); err != nil {
			return nil, fmt.Errorf("tool validation failed: %w", err)
		}
		return &tool, nil
		
	case KindPrompt:
		var prompt Prompt
		if err := yaml.Unmarshal(content, &prompt); err != nil {
			return nil, fmt.Errorf("failed to parse prompt component: %w", err)
		}
		if err := prompt.Validate(); err != nil {
			return nil, fmt.Errorf("prompt validation failed: %w", err)
		}
		return &prompt, nil
		
	case KindAgent:
		var agent Agent
		if err := yaml.Unmarshal(content, &agent); err != nil {
			return nil, fmt.Errorf("failed to parse agent component: %w", err)
		}
		if err := agent.Validate(); err != nil {
			return nil, fmt.Errorf("agent validation failed: %w", err)
		}
		return &agent, nil
		
	default:
		return nil, fmt.Errorf("unsupported component kind: %s", base.Kind)
	}
}

// ParseComponentFromReader parses a component manifest from an io.Reader.
func (p *componentParser) ParseComponentFromReader(reader io.Reader) (interface{}, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}
	return p.ParseComponent(content)
}

// ParseComponentFromString parses a component manifest from a string.
func (p *componentParser) ParseComponentFromString(content string) (interface{}, error) {
	return p.ParseComponent([]byte(content))
}

// SerializeComponent serializes a component to YAML.
func (p *componentParser) SerializeComponent(component interface{}) ([]byte, error) {
	log.Info("Serializing component to YAML")
	
	data, err := yaml.Marshal(component)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize component: %w", err)
	}
	
	return data, nil
}

// SerializeComponentToString serializes a component to a YAML string.
func (p *componentParser) SerializeComponentToString(component interface{}) (string, error) {
	data, err := p.SerializeComponent(component)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ValidateComponentKind validates that a component kind is supported.
func (p *componentParser) ValidateComponentKind(kind string) error {
	if !IsValidKind(kind) {
		validKinds := make([]string, len(ValidComponentKinds()))
		for i, k := range ValidComponentKinds() {
			validKinds[i] = string(k)
		}
		return fmt.Errorf("invalid component kind '%s', must be one of: %s", 
			kind, strings.Join(validKinds, ", "))
	}
	return nil
}

// GetComponentInfo extracts basic information from a component.
func (p *componentParser) GetComponentInfo(component interface{}) (*types.ComponentInfo, error) {
	var base *BaseComponent
	
	switch c := component.(type) {
	case *Tool:
		base = &c.BaseComponent
	case *Prompt:
		base = &c.BaseComponent
	case *Agent:
		base = &c.BaseComponent
	default:
		return nil, fmt.Errorf("unsupported component type: %T", component)
	}
	
	return &types.ComponentInfo{
		APIVersion:  base.APIVersion,
		Kind:        string(base.Kind),
		Name:        base.Metadata.Name,
		Namespace:   base.Metadata.Namespace,
		Version:     base.Metadata.Version,
		Description: base.Metadata.Description,
		Author:      base.Metadata.Author,
		License:     base.Metadata.License,
		Tags:        base.Metadata.Tags,
		Categories:  base.Metadata.Categories,
		Stability:   string(base.Metadata.Stability),
		Maturity:    string(base.Metadata.Maturity),
	}, nil
}

