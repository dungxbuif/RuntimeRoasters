import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import StoreDashboardPage from './page';
import { TestProviders } from '@/tests/utils';
import { retailService } from '@/services/retail.service';

// Mock the retail service
vi.mock('@/services/retail.service', () => ({
  retailService: {
    listStores: vi.fn(),
    listOrders: vi.fn(),
  }
}));

describe('StoreDashboardPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders "No Stores Assigned" when there are no stores', async () => {
    (retailService.listStores as unknown as import('vitest').Mock).mockResolvedValue([]);
    (retailService.listOrders as unknown as import('vitest').Mock).mockResolvedValue([]);

    render(
      <TestProviders>
        <StoreDashboardPage />
      </TestProviders>
    );

    expect(screen.getByText(/Retail Dashboard/i)).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.getByText(/No Stores Assigned/i)).toBeInTheDocument();
    });
  });

  it('renders stores and active orders correctly', async () => {
    const mockStores = [
      { id: 'store-1', name: 'Hoan Kiem Store', location: 'Hanoi', is_active: true }
    ];
    const mockOrders = [
      { id: 'ord-12345678', store_id: 'store-1', status: 'IN_TRANSIT', total_amount: 150000, created_at: new Date().toISOString() }
    ];

    (retailService.listStores as unknown as import('vitest').Mock).mockResolvedValue(mockStores);
    (retailService.listOrders as unknown as import('vitest').Mock).mockResolvedValue(mockOrders);

    render(
      <TestProviders>
        <StoreDashboardPage />
      </TestProviders>
    );

    // Wait for store to be loaded and selected
    await waitFor(() => {
      expect(screen.getByText(/Hoan Kiem Store/i)).toBeInTheDocument();
    });

    // Check if the order is rendered in the DataGrid
    await waitFor(() => {
      expect(screen.getAllByText(/ORD-1234/i).length).toBeGreaterThan(0);
      // Use a more relaxed matcher for currency as Intl.NumberFormat can output non-breaking spaces
      expect(screen.getByText(/150[.,\s]*000/i)).toBeInTheDocument();
      expect(screen.getAllByText(/IN TRANSIT/i).length).toBeGreaterThan(0);
    });
  });
});
