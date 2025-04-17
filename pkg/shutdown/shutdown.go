package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Shutdown a graceful shutdown interface
type Shutdown interface {
	// WithSignals add signals into hook
	WithSignals(signals ...syscall.Signal) Shutdown
	// Close register shutdown handles
	Close(funcs ...func())

	ignore()
}

// shutdown a graceful shutdown implementation
type shutdown struct {
	ch chan os.Signal
}

var _ Shutdown = (*shutdown)(nil)

func (sd shutdown) WithSignals(signals ...syscall.Signal) Shutdown {
	for _, s := range signals {
		signal.Notify(sd.ch, s)
	}
	return sd
}

func (sd shutdown) Close(funcs ...func()) {
	select {
	case <-sd.ch:
	}
	signal.Stop(sd.ch)

	for _, apply := range funcs {
		apply()
	}
}

func (sd shutdown) ignore() {
	panic("implement me")
}

var (
	once     sync.Once
	instance Shutdown
)

// New create a Shutdown instance
func New() Shutdown {
	once.Do(func() {
		instance = &shutdown{
			ch: make(chan os.Signal, 1),
		}
	})

	return instance
}
