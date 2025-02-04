package cosmosutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

// ChainGenesis represents the stargate genesis file
type ChainGenesis struct {
	AppState struct {
		Auth struct {
			Accounts []struct {
				Address string `json:"address"`
			} `json:"accounts"`
		} `json:"auth"`
	} `json:"app_state"`
}

// HasAccount check if account exist into the genesis account
func (g ChainGenesis) HasAccount(address string) bool {
	for _, account := range g.AppState.Auth.Accounts {
		if account.Address == address {
			return true
		}
	}
	return false
}

// ParseGenesis parse ChainGenesis object from a genesis file
func ParseGenesis(genesisPath string) (genesis ChainGenesis, err error) {
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return genesis, errors.New("cannot open genesis file: " + err.Error())
	}
	return genesis, json.Unmarshal(genesisFile, &genesis)
}

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := ParseGenesis(genesisPath)
	if err != nil {
		return false, err
	}
	return genesis.HasAccount(addr), nil
}

// GenesisAndHashFromURL fetches the genesis from the given url and returns its content along with the sha256 hash.
func GenesisAndHashFromURL(ctx context.Context, url string) (genesis []byte, hash string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	genesis, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(genesis)); err != nil {
		return nil, "", err
	}

	hexHash := hex.EncodeToString(h.Sum(nil))

	return genesis, hexHash, nil
}
