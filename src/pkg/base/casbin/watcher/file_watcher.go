package watcher

import (
	"github.com/casbin/casbin/v3"
	"github.com/fsnotify/fsnotify"
	"RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

// WatchCasbinFiles monitors changes to model and policy files and reloads the enforcer automatically.
func WatchCasbinFiles(enforcer *casbin.SyncedEnforcer, modelPath string, policyPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.GetLogger().Error("Failed to create casbin file watcher", zap.Error(err))
		return
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// We care about write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					log := logger.GetLogger()
					if event.Name == modelPath {
						log.Info("Casbin model file changed, reloading...", zap.String("path", modelPath))
						// For model changes, we need to reload the model from file
						if err := enforcer.LoadModel(); err != nil {
							log.Error("Failed to reload casbin model", zap.Error(err))
						}
					} else if event.Name == policyPath {
						log.Info("Casbin policy file changed, reloading...", zap.String("path", policyPath))
						// For policy changes in CSV, we load into a temp enforcer and sync to the main one
						tempEnforcer, err := casbin.NewEnforcer(enforcer.GetModel(), policyPath)
						if err != nil {
							log.Error("Failed to load policy from CSV", zap.Error(err))
							continue
						}
						
						enforcer.ClearPolicy()
						policies, err := tempEnforcer.GetPolicy()
						if err == nil {
							_, _ = enforcer.AddPolicies(policies)
						}
						
						groupingPolicies, err := tempEnforcer.GetGroupingPolicy()
						if err == nil {
							_, _ = enforcer.AddGroupingPolicies(groupingPolicies)
						}
						
						if err := enforcer.SavePolicy(); err != nil {
							log.Error("Failed to save synced policies to DB", zap.Error(err))
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.GetLogger().Error("Casbin watcher error", zap.Error(err))
			}
		}
	}()

	err = watcher.Add(modelPath)
	if err != nil {
		logger.GetLogger().Error("Failed to watch casbin model file", zap.Error(err), zap.String("path", modelPath))
	}

	if policyPath != "" {
		err = watcher.Add(policyPath)
		if err != nil {
			logger.GetLogger().Error("Failed to watch casbin policy file", zap.Error(err), zap.String("path", policyPath))
		}
	}
}
