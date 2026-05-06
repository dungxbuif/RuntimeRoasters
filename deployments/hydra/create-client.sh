#!/bin/bash

echo "Wait for Hydra to be ready..."
until curl -s http://localhost:4445/health/ready; do
  sleep 2
done

echo "Registering Client App..."
docker exec rr-hydra \
    hydra create client \
    --endpoint http://localhost:4445/ \
    --id client-app \
    --secret client-secret \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --scope openid,offline_access,email \
    --redirect-uri http://localhost:3000/api/auth/callback

echo "Client App Registered Successfully."
