package listener

import (
	"net"

	"fmt"

	"github.com/dmi3midd/grpcsso/internal/config"
)

type Listener struct {
	cfg *config.ServerConfig
}

func NewListener(cfg *config.ServerConfig) *Listener {
	return &Listener{cfg: cfg}
}

func (l *Listener) Listen() (net.Listener, error) {
	op := "listener.Listen"

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", l.cfg.Host, l.cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lis, nil
}
