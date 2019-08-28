package watcher

import (
	"udns/logger"

	"github.com/fsnotify/fsnotify"
)

var (
	w *fsnotify.Watcher

	watcherMap     = map[string]func() error{}
	watchingEvents = []fsnotify.Op{
		fsnotify.Create,
		fsnotify.Rename,
		fsnotify.Remove,
		fsnotify.Write,
	}
	stopEvt    = make(chan struct{}, 1)
	stoppedEvt = make(chan struct{}, 1)
)

func init() {
	nw, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("watcher", err)
	}
	w = nw
	go watching()
}

// Remove remove one handler
func Remove(file string) error {
	return w.Remove(file)
}

// On add handler
func On(file string, handler func() error) {
	logger.Infof("watcher", "add watcher %s", file)
	if err := w.Add(file); err != nil {
		logger.Error("watcher", err)
		return
	}
	watcherMap[file] = handler
}

// Stop stop the watching goroutine
func Stop() {
	stopEvt <- struct{}{}
	<-stoppedEvt
	logger.Info("watcher", "stopped")
}

func watching() {
	defer func() {
		stoppedEvt <- struct{}{}
	}()
	defer w.Close()

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				logger.Error("watcher", "failed to read watch event")
				return
			}
			logger.Infof("watcher", "got event: %v", event)
			handler, exists := watcherMap[event.Name]
			if !exists {
				continue
			}
			for _, evt := range watchingEvents {
				if event.Op&evt == evt {
					if evt != fsnotify.Write {
						logger.Infof("watcher", "add watcher %s", event.Name)
						if err := w.Add(event.Name); err != nil {
							logger.Error("watcher", err)
							continue
						}
					}
					if err := handler(); err != nil {
						logger.Errorf("watcher", "handle %s error, %v", event.Name, err)
					}
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			logger.Errorf("watcher", "got watcher error %v", err)
		case <-stopEvt:
			logger.Warn("watcher", "got stop event")
			return
		}
	}
}
