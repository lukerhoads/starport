package pluginsrpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	gogetter "github.com/hashicorp/go-getter"
)

// MUST BE RAN BEFORE BUILD
func (m *Manager) Pull(ctx context.Context) error {
	fmt.Println("🤏 Pulling plugins...")

	pluginHome, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	outputDir := path.Join(pluginHome, "cached")
	dir, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, d := range dir {
		os.RemoveAll(path.Join(outputDir, d.Name()))
	}

	for _, cfgPlugin := range m.Config.Plugins {
		// Seperate individual plugins by ID
		plugId := getPluginId(cfgPlugin)

		// Check GOPATH for plugin?

		// Get plugin home
		dst, err := formatPluginHome(m.ChainId, plugId)
		if err != nil {
			return err
		}

		_, err = os.Stat(dst)
		if err == nil {
			err = os.RemoveAll(dst)
			if err != nil {
				return err
			}
		}

		if err := download(cfgPlugin.Repo, "", dst); err != nil {
			return err
		}
	}

	return nil
}

func download(repo string, subdir string, dst string) error {
	url := "git::https://" + repo + ".git"
	// url = repo
	if subdir != "" {
		url += ("//" + subdir)
	}

	// Not cancelling for some noob reason
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client := gogetter.Client{
		Ctx:  ctx,
		Src:  url,
		Dst:  dst,
		Mode: gogetter.ClientModeAny,
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}
