# [RR-2] Core Framework - Config & Logger

- **Goal:** Build a shared infrastructure for configuration and logging.
- **Business Value:** Unify configuration and logging across the entire system to simplify operations and debugging.
- **Priority:** `HIGH`

## 📝 Description
Implement a dynamic configuration loading module (Viper) and a structured Logger (Uber Zap) to standardize operations for 10+ services.

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Dynamically Loading Configuration from .env
- **Given:** A `.env` file containing the variable `APP_PORT=8080`.
- **When:** I start the service and call the `config.LoadConfig()` function.
- **Then:** The value `8080` must be correctly mapped to the Go configuration struct.

### Scenario 2: Structured Logging in Development Environment
- **Given:** The service is running in `development` mode.
- **When:** I call `logger.Info("Hello World")`.
- **Then:** The log must be displayed in **Console** format with colors and caller information (source code line).

### Scenario 3: Structured Logging in Production Environment
- **Given:** The service is running in `production` mode.
- **When:** I call `logger.Info("Hello World")`.
- **Then:** The log must be displayed in standard **JSON** format for log management tools to parse.
