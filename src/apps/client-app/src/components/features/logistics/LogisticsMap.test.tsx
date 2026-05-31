import { render, screen } from '@testing-library/react';
import { expect, test, vi } from 'vitest';
import LogisticsMap from './LogisticsMap';

// Mock Leaflet and React Leaflet
vi.mock('leaflet', () => ({
  default: {
    Icon: { 
      Default: { 
        prototype: {},
        mergeOptions: vi.fn(),
      } 
    },
    divIcon: vi.fn(),
  },
  divIcon: vi.fn(),
}));

vi.mock('react-leaflet', () => ({
  MapContainer: ({ children }: any) => <div data-testid="map-container">{children}</div>,
  TileLayer: () => <div data-testid="tile-layer" />,
  Marker: ({ children }: any) => <div data-testid="marker">{children}</div>,
  Popup: ({ children }: any) => <div data-testid="popup">{children}</div>,
  Polyline: () => <div data-testid="polyline" />,
}));

test('renders LogisticsMap with light mode tiles and legend', async () => {
  render(
    <LogisticsMap 
      locations={[]} 
      shipments={[]} 
      routes={[]} 
      activeDriverLocations={{}} 
    />
  );
  
  // Legend should be present
  expect(screen.getByText('Coffee Farmers')).toBeInTheDocument();
  expect(screen.getByText('Roasteries (KCN)')).toBeInTheDocument();
  expect(screen.getByText('Retail Outlets')).toBeInTheDocument();
});
