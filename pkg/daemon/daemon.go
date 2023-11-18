package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/kardianos/service"
)

func New(name, description string, daemon Daemon) (service.Service, error) {
	return service.New(daemon, &service.Config{
		Name:        name,
		DisplayName: name,
		Description: description,
		Option:      options(),
	})
}

type Daemon interface {
	service.Interface
}

func Periodic(interval time.Duration, runner func(service.Logger) error) Daemon {
	return &periodicDaemon{interval: interval, runner: runner}
}

type periodicDaemon struct {
	interval time.Duration
	runner   func(service.Logger) error

	lock   sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func (dmn *periodicDaemon) Start(s service.Service) error {
	dmn.lock.Lock()
	if dmn.ctx != nil {
		dmn.lock.Unlock()
		return nil
	}
	dmn.ctx, dmn.cancel = context.WithCancel(context.Background())
	dmn.lock.Unlock()

	go func() {
		logger, err := s.Logger(nil)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-time.After(dmn.interval):
				if err := dmn.runner(logger); err != nil {
					logger.Error(err)
				}
			case <-dmn.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (dmn *periodicDaemon) Stop(service.Service) error {
	dmn.lock.Lock()
	defer dmn.lock.Unlock()
	if dmn.ctx != nil {
		dmn.cancel()
		dmn.ctx, dmn.cancel = nil, nil
	}
	return nil
}
