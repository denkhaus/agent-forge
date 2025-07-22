package database

import (
	"context"
	"fmt"

	"github.com/denkhaus/agentforge/internal/database/ent/localconfig"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// configService provides configuration management operations.
type configService struct {
	client DatabaseClient
}

// NewConfigService creates a new config service.
func NewConfigService(client DatabaseClient) ConfigService {
	return &configService{
		client: client,
	}
}

// GetConfig retrieves a configuration value by key.
func (cs *configService) GetConfig(ctx context.Context, key string) (string, error) {
	config, err := cs.client.GetEnt().LocalConfig.Query().
		Where(localconfig.Key(key)).
		Only(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get configuration: %w", err)
	}
	return config.Value, nil
}

// SetConfig sets a configuration value.
func (cs *configService) SetConfig(ctx context.Context, key, value, configType string) error {
	// Try to update existing configuration
	existing, err := cs.client.GetEnt().LocalConfig.Query().
		Where(localconfig.Key(key)).
		Only(ctx)
	
	if err == nil {
		// Update existing
		_, err = cs.client.GetEnt().LocalConfig.UpdateOne(existing).
			SetValue(value).
			SetType(localconfig.Type(configType)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update configuration: %w", err)
		}
	} else {
		// Create new
		_, err = cs.client.GetEnt().LocalConfig.Create().
			SetID(uuid.New().String()).
			SetKey(key).
			SetValue(value).
			SetType(localconfig.Type(configType)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create configuration: %w", err)
		}
	}
	
	log.Info("Configuration updated", zap.String("key", key))
	return nil
}

// GetBoolConfig retrieves a boolean configuration value by key.
func (cs *configService) GetBoolConfig(ctx context.Context, key string) (bool, error) {
	value, err := cs.GetConfig(ctx, key)
	if err != nil {
		return false, err
	}
	return value == "true", nil
}

// GetIntConfig retrieves an integer configuration value by key.
func (cs *configService) GetIntConfig(ctx context.Context, key string) (int, error) {
	value, err := cs.GetConfig(ctx, key)
	if err != nil {
		return 0, err
	}
	
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return 0, fmt.Errorf("failed to parse integer config: %w", err)
	}
	return result, nil
}

// SetBoolConfig sets a boolean configuration value.
func (cs *configService) SetBoolConfig(ctx context.Context, key string, value bool) error {
	strValue := "false"
	if value {
		strValue = "true"
	}
	return cs.SetConfig(ctx, key, strValue, "BOOLEAN")
}

// SetIntConfig sets an integer configuration value.
func (cs *configService) SetIntConfig(ctx context.Context, key string, value int) error {
	strValue := fmt.Sprintf("%d", value)
	return cs.SetConfig(ctx, key, strValue, "INTEGER")
}