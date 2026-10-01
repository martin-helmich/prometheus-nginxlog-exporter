package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/log"
)

// ConfigWatcher watches a configuration file for changes
type ConfigWatcher struct {
	watcher    *fsnotify.Watcher
	logger     *log.Logger
	configFile string
	onChange   func()
}

// NewConfigWatcher creates a new configuration watcher
func NewConfigWatcher(logger *log.Logger, configFile string, onChange func()) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &ConfigWatcher{
		watcher:    watcher,
		logger:     logger,
		configFile: configFile,
		onChange:   onChange,
	}

	if err := watcher.Add(configFile); err != nil {
		watcher.Close()
		return nil, err
	}

	go w.watch()

	return w, nil
}

// watch monitors the config file for changes
func (w *ConfigWatcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.logger.Infof("config file changed, triggering reload")
				w.onChange()
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.logger.Errorf("error watching config file: %v", err)
		}
	}
}

// Close stops watching the configuration file
func (w *ConfigWatcher) Close() error {
	return w.watcher.Close()
}
