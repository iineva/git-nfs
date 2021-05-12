package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type callback func(s os.Signal, done func())

var (
	_callbacks     = []callback{}
	_callbacksLock = sync.RWMutex{}
)

func init() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			os.Kill,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		select {
		case s := <-c:
			doCallback(s)
		}
	}()
}

func doCallback(s os.Signal) {
	_callbacksLock.RLock()
	defer _callbacksLock.RUnlock()
	wg := sync.WaitGroup{}
	for _, c := range _callbacks {
		wg.Add(1)
		go c(s, wg.Done)
	}
	wg.Wait()
	os.Exit(0)
}

func AddTermCallback(cb callback) {
	_callbacksLock.Lock()
	defer _callbacksLock.Unlock()
	_callbacks = append(_callbacks, cb)
}
