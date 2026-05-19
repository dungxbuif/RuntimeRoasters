# REST API Specifications

Runtime Roasters exposes a unified RESTful interface via the **KrakenD API Gateway**.

## 🚪 API Gateway (KrakenD)
- **Host:** `http://localhost:8081`
- **Configuration:** `deployments/krakend/krakend.json`
- **Features:** JWT Validation, Rate Limiting, Request Transformation, Aggregation.

## 📑 Interactive Documentation (Swagger/OpenAPI)
Each service generates its own Swagger specification. The Gateway can also export a consolidated OpenAPI spec.

### How to view:
1.  Start the services: `task dev`
2.  The Frontend Portal (Client App) includes an embedded Swagger UI at `/dashboard/api-docs`.
3.  Alternatively, you can import the raw JSON specs from `src/apps/{service}/docs/swagger.json` into Postman or Insomnia.

## 🔐 Authentication
Most endpoints require a **Bearer JWT**.
- **Header:** `Authorization: Bearer <token>`
- **Obtaining a token:** Use the `/v1/auth/login` flow via the Frontend or call the Kratos/Hydra APIs directly.
