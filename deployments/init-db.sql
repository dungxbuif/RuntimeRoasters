-- ----------------------------------------------------------------------
-- RuntimeRoasters Database Initialization
-- ----------------------------------------------------------------------

-- Create Databases
CREATE DATABASE identity_db;  -- Ory Kratos
CREATE DATABASE hydra_db;     -- Ory Hydra
CREATE DATABASE auth_db;      -- Auth Service
CREATE DATABASE audit_db;     -- Audit Service
CREATE DATABASE demo_db;      -- Demo Service
CREATE DATABASE farm_db;      -- Farm Service
CREATE DATABASE warehouse_db; -- Warehouse Service
CREATE DATABASE retail_db;    -- Retail Service
CREATE DATABASE logistics_db; -- Logistics Service
CREATE DATABASE payment_db;   -- Payment Service
CREATE DATABASE trace_db;     -- Trace Service

-- Apply Schemas (Simplified approach for dev)
-- For a cleaner production approach, use a proper migration tool.
-- In this dev setup, we'll pipe the schema files during container init.
