# Impact Analysis

- **Affected Components**: 
  - `auth-service` (User creation)
  - `logistics-service` (Location & Fleet seeding)
  - `trace-service` (Historical Elasticsearch events insertion)
  - Demo initialization scripts (`deployments/reset-demo-state.sh`)
  - Demo documentation (`docs/product/GUIDE.md`)

- **Business Impact**: 
  - Essential for the Showcase objective. Reviewers and stakeholders need to see the Traceability and Finance capabilities immediately upon system launch without spending 10-15 minutes manually generating data through the UI.
  
- **Risk Level**: 
  - **Medium**. Seeding historical logic can be complex because it requires mimicking the exact sequence of Domain Events (Event Sourcing) with logically staggered timestamps (e.g., T-5 days, T-4 days) to prevent breaking data integrity rules in the read models.
