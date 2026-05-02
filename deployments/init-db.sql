-- Khởi tạo Database cho từng Microservice
CREATE DATABASE demo_db;
CREATE DATABASE process_db;
CREATE DATABASE warehouse_db;
CREATE DATABASE logistics_db;
CREATE DATABASE retail_db;
CREATE DATABASE payment_db;
CREATE DATABASE audit_db;
CREATE DATABASE webhook_db;

-- Chú thích: Trace Service dùng Elasticsearch nên không cần SQL DB.
