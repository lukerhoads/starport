package plugins

import (
	"context"
	"reflect"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

// NOTE: Go plugins were added in 1.8, so using them is questionable
// unless we can build the plugins with a greater go version and somehow
// grab the exposed functions with that build.

// Plugin manager
type Manager struct {
	ChainId string `yaml:"chain_id"`

	cmdPlugins  []CmdPlugin
	hookPlugins []HookPlugin
}

func NewManager(chainId string) Manager {
	return Manager{
		ChainId: chainId,
	}
}

// Run will go through all of the steps:
// - Downloading
// - Building
// - Cache .so files
// - Execution and Injection
func (m *Manager) Run(ctx context.Context, cfg chaincfg.Config, rootCommand *cobra.Command) error {
	// Check for change in config contents since last
	// Don't check for remote package changes, as theoretically we want it
	// to be up to the user to reload the plugins.
	configChanged, err := pluginsChanged(cfg, m.ChainId)
	if err != nil {
		return err
	}

	if configChanged {
		return nil
	}

	if err := m.pull(ctx, cfg); err != nil {
		return err
	}

	// Build
	if err := m.build(ctx, cfg); err != nil {
		return err
	}

	// Extraction
	if err := m.extractPlugins(ctx, rootCommand, cfg); err != nil {
		return err
	}

	// Injection
	if err := m.inject(ctx, rootCommand); err != nil {
		return err
	}

	return nil
}

// Check if plugin-specified configuration is different from downloaded plugins
// For now, ONLY CHECKS DIRECTORY NAMES
// This is not adequate, because one could delete files from directories
func pluginsChanged(cfg chaincfg.Config, chainId string) (bool, error) {
	var configChanged bool
	var configPluginNames []string
	var fileConfigNames []string

	for _, plug := range cfg.Plugins {
		plugId := getPluginId(plug)
		configPluginNames = append(configPluginNames, plugId)
	}

	dst, err := formatPluginHome(chainId, "")
	if err != nil {
		return false, err
	}

	pluginDirs, err := listDirs(dst)
	if err != nil {
		return false, err
	}

	for _, plugDir := range pluginDirs {
		fileConfigNames = append(fileConfigNames, plugDir.Name())
	}

	return !reflect.DeepEqual(configPluginNames, fileConfigNames), nil
}
