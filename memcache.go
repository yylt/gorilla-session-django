package gorilla_session_django

import (
	"fmt"

	"github.com/memcachier/mc"
)

type Memcacher interface {
	Get(key string) (val string, err error)
	Set(key, val string, ttl uint32) (err error)
}

type memcache struct {
	cli *mc.Client
	cfg *MemCfg
}

type MemCfg struct {
	Endpoint string
	User     string
	Password string
}

func NewMemCli(cfg *MemCfg) (Memcacher, error) {
	cli := mc.NewMC(cfg.Endpoint, cfg.User, cfg.Password)
	if cli == nil {
		return nil, fmt.Errorf("create memcache client failed, config %v", cfg)
	}
	mem := &memcache{
		cli: cli,
		cfg: cfg,
	}
	return mem, nil
}

func (m *memcache) Get(key string) (val string, err error) {
	val, _, _, err = m.cli.Get(key)
	return
}

//Set key value
func (m *memcache) Set(key, val string, ttl uint32) (err error) {
	_, err = m.cli.Set(key, val, 0, ttl, 0)
	return err
}
