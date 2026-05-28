# Root Cause Analysis

**Symptom**: `newEnforcer is not a function` error when a manager logs in for the first time, requiring a double click to bypass.
**Actual Cause**: In `casbin.tsx`, the named import `import { newEnforcer } from 'casbin.js'` fails because Turbopack/Next.js handles CommonJS exports strictly, meaning `newEnforcer` is undefined.
