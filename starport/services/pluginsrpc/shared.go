package pluginsrpc

import (
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	plugintypes "github.com/tendermint/starport/starport/services/pluginsrpc/types"
)

var (
	// Blacklisted strings in an RPC log broadcast output
	excludedSet = []string{
		"waiting for RPC",
		"plugin started",
		"plugin process",
		"using plugin",
		"plugin address",
		"plugin exited",
		"plugin server",
		"plugin: path",
		"starting plugin",
	}

	// PluginLogger that filters unnecessary debug messages
	pluginLogger = hclog.New(&hclog.LoggerOptions{
		Output:      hclog.DefaultOutput,
		Level:       hclog.Trace,
		DisableTime: true,
		Exclude: func(level hclog.Level, msg string, args ...interface{}) bool {
			for _, excluded := range excludedSet {
				if strings.Contains(msg, excluded) {
					return true
				}
			}

			return false
		},
	})
)

// Basic handshake config used for all RPC connections
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// BasePluginMap represents what the extractor expects to see available in an RPC plugin.
var BasePluginMap = map[string]plugin.Plugin{
	"command_map": &plugintypes.CommandMapperPlugin{},
	"hook_map":    &plugintypes.HookMapperPlugin{},
}

// Extracted representation of a command module.
type ExtractedCommandModule struct {
	ModuleName string
	PluginDir  string

	ParentCommand []string
	Name          string
	Usage         string
	ShortDesc     string
	LongDesc      string
	NumArgs       int
	Exec          func(*cobra.Command, []string) error
}

// Extracted representation of a hook module.
type ExtractedHookModule struct {
	ModuleName string
	PluginDir  string

	ParentCommand []string
	Name          string
	HookType      string
	PreRun        func(*cobra.Command, []string) error
	PostRun       func(*cobra.Command, []string) error
}

// PluginState represents the states a plugin can be in
type PluginState uint32

const (
	Undefined PluginState = iota
	Configured
	Downloaded
	Built
)

func PluginStateFromString(state string) PluginState {
	switch state {
	case "configured":
		return Configured
	case "downloaded":
		return Downloaded
	case "built":
		return Built
	}

	return Undefined
}

func (p PluginState) String() string {
	switch p {
	case Configured:
		return "configured"
	case Downloaded:
		return "downloaded"
	case Built:
		return "built"
	}

	return "undefined"
}
