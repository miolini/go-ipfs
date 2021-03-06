package fsrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/facebookgo/atomicfile"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/util"
	logging "github.com/ipfs/go-ipfs/vendor/QmXJkcEXB6C9h6Ytb6rrUTFU56Ro62zxgrbxTT3dgjQGA8/go-log"
)

var log = logging.Logger("fsrepo")

// ReadConfigFile reads the config from `filename` into `cfg`.
func ReadConfigFile(filename string, cfg interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("Failure to decode config: %s", err)
	}
	return nil
}

// WriteConfigFile writes the config from `cfg` into `filename`.
func WriteConfigFile(filename string, cfg interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0775)
	if err != nil {
		return err
	}

	f, err := atomicfile.New(filename, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	return encode(f, cfg)
}

// encode configuration with JSON
func encode(w io.Writer, value interface{}) error {
	// need to prettyprint, hence MarshalIndent, instead of Encoder
	buf, err := config.Marshal(value)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// Load reads given file and returns the read config, or error.
func Load(filename string) (*config.Config, error) {
	// if nothing is there, fail. User must run 'ipfs init'
	if !util.FileExists(filename) {
		return nil, errors.New("ipfs not initialized, please run 'ipfs init'")
	}

	var cfg config.Config
	err := ReadConfigFile(filename, &cfg)
	if err != nil {
		return nil, err
	}

	// tilde expansion on datastore path
	// TODO why is this here??
	cfg.Datastore.Path, err = util.TildeExpansion(cfg.Datastore.Path)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
