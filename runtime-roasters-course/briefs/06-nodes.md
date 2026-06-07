# Module 6: When Nodes Whisper

### Teaching Arc
- **Metaphor:** The "City at Night". From high up, you see the whole city — traffic moving between buildings, lights turning on and off. Our "Topology" is the city map of our services.
- **Opening hook:** How do we see the whole system at once? How do we know if one service is "lonely" or if the traffic is flowing correctly?
- **Key insight:** "Observability" means being able to see what's happening inside a complex system from the outside.
- **"Why should I care?":** This helps you debug "ghost problems" where everything looks fine but nothing is working. You'll know how to ask AI for "Distributed Tracing" or "Health Checks".

### Code Snippets (pre-extracted)

File: src/apps/client-app/src/services/topology.service.ts (lines 10-25)
```typescript
class TopologyService {
  async getPublicConfig(): Promise<TopologyConfig> {
    const res = await api.get(API_ENDPOINTS.TRACE.PUBLIC_TOPOLOGY_CONFIG);
    return res.data.config;
  }

  async getHistory(): Promise<TopologyHistoryEntry[]> {
    const res = await api.get(API_ENDPOINTS.TRACE.TOPOLOGY_HISTORY);
    return res.data.history;
  }
}
```

### Interactive Elements

- [x] **Topology Mesh** — An interactive map where services "glow" when they receive an event.
- [x] **Layer Toggle** — A toggle to switch between "Public View" (just services) and "Admin View" (IPs, Latency, Error rates).
- [x] **Group Chat Animation** — "Behind the scenes" whispers: [SigNoz: "I see a slow request!", Jaeger: "It's coming from Warehouse.", SigNoz: "Found it. Latency is 500ms."].
- [x] **Quiz** — 4 questions, style: architecture. "Which service is the 'Traffic Controller' that watches everyone else?"

### Reference Files to Read

- `references/interactive-elements.md` → "Topology Mesh", "Layer Toggle"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** Transparent Proof (Individual data)
- **Next module:** Steering the Roaster (Conclusion)
- **Tone/style notes:** Cinematic, "Big Picture", "Matrix" vibe.
