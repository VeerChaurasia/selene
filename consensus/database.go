package consensus
import (
	"errors"
	"github.com/BlocSoc-iitr/selene/config"
	"os"
	"path/filepath"
)
type Database interface {
	New(config *config.Config) (Database, error)
	SaveCheckpoint(checkpoint []byte) error
	LoadCheckpoint() ([]byte, error)
}
type FileDB struct {
	DataDir           string
	defaultCheckpoint []byte
}
func (f *FileDB) New(config *config.Config) (Database, error) {
	if config.DataDir == "" {
		return nil, errors.New("data dir not in config")
	}
	return &FileDB{
		dataDir:           config.DataDir,
		defaultCheckpoint: config.DefaultCheckpoint,
	}, nil
}
func (f *FileDB) SaveCheckpoint(checkpoint []byte) error {
	err := os.MkdirAll(f.DataDir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(f.DataDir, "checkpoint"), checkpoint, 0644)
}
func (f *FileDB) LoadCheckpoint() ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(f.DataDir, "checkpoint"))
	if err != nil {
		if os.IsNotExist(err) {
			return f.defaultCheckpoint, nil
		}
		return nil, err
	}
	if len(data) == 32 {
		return data, nil
	}
	return f.defaultCheckpoint, nil
}
type ConfigDB struct {
	checkpoint []byte
}
func (c *ConfigDB) New(config *config.Config) (Database, error) {
	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = config.DefaultCheckpoint
	}
	return &ConfigDB{
		checkpoint: checkpoint,
	}, nil
}
func (c *ConfigDB) SaveCheckpoint(checkpoint []byte) error {
	return nil
}
func (c *ConfigDB) LoadCheckpoint() ([]byte, error) {
	return c.checkpoint, nil
}
