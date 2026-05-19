# Technical Plan - [RR-2] Core Framework - Config & Logger

## 🎯 Implementation Strategy
- Use Viper to manage configuration from multiple sources (env, file).
- Use Zap as a high-performance logger supporting structured logging.

## 🛠️ Implementation Steps
- [ ] Implement `pkg/config/config.go` using `spf13/viper`.
- [ ] Implement `pkg/logger/logger.go` using `uber-go/zap`.
    - Support Log Level via the `LOG_LEVEL` environment variable.
    - Support Console format for Dev and JSON for Prod.
- [ ] Create `.env.example` as a template.

## 🧪 Verification
- [ ] Unit test for the config module.
- [ ] Verify log format on the terminal in Dev/Prod environments.
