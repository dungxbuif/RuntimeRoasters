# Module 2: The Cast of Services

### Teaching Arc
- **Metaphor:** A "Professional Relay Race". In a relay, no single person runs the whole track. One runner finishes their leg and hands the baton to the next. Our services do the same with your coffee's data.
- **Opening hook:** Runtime Roasters isn't one big app — it's a team of specialized services working together. Who are these "actors" and what are their roles?
- **Key insight:** "Microservices" means breaking a big problem into small, manageable jobs.
- **"Why should I care?":** Knowing who does what helps you tell AI exactly where to add a new feature — like "Add a loyalty discount to the Retail service," not "Add it to the app."

### Code Snippets (pre-extracted)

File: deployments/docker-compose.dev.yaml (lines 50-70 - representative)
```yaml
  farm-service:
    build: ../src/apps/farm-service
    ports: ["8081:8080"]
  
  warehouse-service:
    build: ../src/apps/warehouse-service
    ports: ["8082:8080"]

  retail-service:
    build: ../src/apps/retail-service
    ports: ["8084:8080"]
```

### Interactive Elements

- [x] **Group chat animation** — actors: [Farm, Warehouse, Retail, Trace]. Message flow: [Farm: "I have 500kg!", Warehouse: "I'll send a driver.", Retail: "I need stock!", Trace: "I'm watching everyone."].
- [x] **Architecture Diagram** — A simplified version of the system topology showing the 5 core services.
- [x] **Pattern Cards** — 5 cards, one for each service with an icon and a "One-Sentence Job Description".
- [x] **Quiz** — 4 questions, style: tracing. "If a customer wants to see where their beans were roasted, which service provides that proof?"

### Reference Files to Read

- `references/interactive-elements.md` → "Group Chat Animation", "Pattern Cards", "Architecture Diagrams"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** The Digital Birth (Individual data creation)
- **Next module:** The Logistics Pulse (How they talk: Kafka & Events)
- **Tone/style notes:** Focus on the "Teamwork" aspect. Use distinct icons for each service.
