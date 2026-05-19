# Technical Plan - [RR-3] Core Framework - Base App & Error Handling

## 🎯 Strategy
- Build a base app for application lifecycle management.
- Standardize error format according to RFC 9457.

## 🛠️ Implementation Steps
- [ ] Implement `pkg/base/app.go`.
    - Manage signals and wait groups for graceful shutdown.
    - Setup basic Health Check mechanism.
- [ ] Implement `pkg/errs/problem.go` to map Go errors to the Problem JSON struct.

## 🧪 Verification
- [ ] Verify the Graceful Shutdown mechanism using the kill command.
- [ ] Mock errors and verify the response body follows the RFC 9457 standard.
