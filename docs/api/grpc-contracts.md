# gRPC Contracts

Runtime Roasters uses **Protocol Buffers (Protobuf)** as the source of truth for all internal service-to-service communication.

## 📁 Source of Truth
All `.proto` definitions are located in the root `/api/` directory:
- `api/runtime/auth/v1/`: Identity and Authorization.
- `api/runtime/farm/v1/`: Plantation and Harvest logic.
- `api/runtime/warehouse/v1/`: Inventory management.
- `api/runtime/retail/v1/`: Order and Transaction logic.
- ...and so on.

## 🛠️ Code Generation
We use **Buf** to manage plugins and generate Go code.
- Configuration: `api/buf.yaml`, `api/buf.gen.yaml`
- Command: `task proto` (generates code into `src/pkg/api/` or directly into service folders depending on config).

## 🛡️ Best Practices
- **Breaking Changes:** Never remove or rename fields; use `reserved` if necessary.
- **Documentation:** Use comments in `.proto` files to document RPCs and fields; they will be extracted into the generated code and Swagger specs.
