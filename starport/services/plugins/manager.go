package plugins

import (
	"context"
	"fmt"

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
func (m *Manager) RunAll(ctx context.Context, cfg chaincfg.Config, rootCommand *cobra.Command) error {
	if len(cfg.Plugins) == 0 {
		return nil
	}

	if err := m.PullBuild(ctx, cfg); err != nil {
		return err
	}

	// Extraction
	if err := m.extractPlugins(ctx, rootCommand, cfg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) PullBuild(ctx context.Context, cfg chaincfg.Config) error {
	if len(cfg.Plugins) == 0 {
		return nil
	}

	// Check for change in config contents since last
	// Don't check for remote package changes, as theoretically we want it
	// to be up to the user to reload the plugins.
	// configChanged, err := PluginsChanged(cfg, m.ChainId)
	// if err != nil {
	// 	return err
	// }

	// if !configChanged {
	// 	return nil
	// }

	fmt.Println("pulling")
	// Pull
	if err := m.pull(ctx, cfg); err != nil {
		return err
	}

	fmt.Println("building")
	// Build
	if err := m.build(ctx, cfg); err != nil {
		return err
	}

	return nil
}
