import { render, screen } from '@testing-library/react';
import { expect, test, vi } from 'vitest';
import { ArchitectureDiagramCanvas } from './ArchitectureTopology';

// Mock dependencies
vi.mock('reactflow', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="react-flow">{children}</div>,
  Controls: () => <div data-testid="controls" />,
  Handle: () => <div />,
  Position: { Left: 'left', Right: 'right' },
  MarkerType: { ArrowClosed: 'arrowclosed' },
}));

vi.mock('@/services/topology.service', () => ({
  topologyService: {
    getPublicConfig: vi.fn().mockResolvedValue({ nodes: [], edges: [], flows: [] }),
    getPublicHistory: vi.fn().mockResolvedValue([]),
    getTicket: vi.fn().mockResolvedValue('ticket'),
    publicSocketURL: vi.fn().mockReturnValue('ws://localhost'),
    privateSocketURL: vi.fn().mockReturnValue('ws://localhost'),
  },
}));

test('renders ArchitectureDiagramCanvas with fixed frames layout', () => {
  render(<ArchitectureDiagramCanvas />);
  
  // It should render the flow container
  expect(screen.getByTestId('react-flow')).toBeInTheDocument();
  expect(screen.getByText('History')).toBeInTheDocument();
  expect(screen.getByText('Flow')).toBeInTheDocument();
});
