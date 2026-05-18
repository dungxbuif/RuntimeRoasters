import { ArchitectureDiagramCanvas } from '@/components/features/architecture-topology/ArchitectureTopology';

export default function SystemTopologyPage() {
  return (
    <div className="h-full bg-white overflow-auto p-6">
      <ArchitectureDiagramCanvas />
    </div>
  );
}
