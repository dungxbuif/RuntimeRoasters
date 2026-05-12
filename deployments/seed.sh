#!/bin/sh

# ===========================================================================
# RuntimeRoasters — System Seeder
#
# Mục đích: Tự động hóa việc tạo dữ liệu mẫu (Admin User, OAuth2 Client)
# để PO/Dev có thể sử dụng ngay sau khi hạ tầng được dựng lên.
# ===========================================================================

echo "🌱 Starting System Seeder..."

# 1. Chờ Kratos & Hydra sẵn sàng
echo "⏳ Waiting for Ory Kratos & Hydra to be healthy..."
until curl -s http://rr-kratos:4433/health/ready && curl -s http://rr-hydra:4445/health/ready; do
  echo "Still waiting..."
  sleep 3
done

# 2. Seed Admin User vào Ory Kratos
echo "👤 Provisioning Admin User in Kratos..."
# Kiểm tra xem user đã tồn tại chưa bằng cách search email qua Admin API
USER_EXISTS=$(curl -s http://rr-kratos:4434/admin/identities?credentials_identifier=admin@runtimeroasters.com)

if [ "$USER_EXISTS" = "[]" ] || [ -z "$USER_EXISTS" ]; then
  curl -X POST \
    -H "Content-Type: application/json" \
    -d @/etc/config/kratos/seed-admin.json \
    http://rr-kratos:4434/admin/identities
  echo "✅ Admin User Created."
else
  echo "ℹ️ Admin User already exists, skipping."
fi

# 3. Seed OAuth2 Client vào Ory Hydra
echo "🔑 Provisioning OAuth2 Client in Hydra..."
# Hydra create oauth2-client sẽ báo lỗi nếu client_id đã tồn tại, nên chúng ta dùng import (Upsert)
cat /etc/config/hydra/client-app.json | \
  docker exec -i rr-hydra hydra import oauth2-client --endpoint http://localhost:4445/

echo "✅ Client App Registered."

echo "✨ Seeding Completed Successfully."
