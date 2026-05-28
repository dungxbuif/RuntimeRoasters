import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { CasbinProvider, useCasbin } from '../casbin';
import { useAuth } from '../AuthProvider';
import { authService } from '@/services/auth.service';
import { AUTH_RESOURCES, AUTH_ACTIONS } from '@/constants/resources';
import { AuthContextType } from '../types';

// Mock dependencies
vi.mock('../AuthProvider', () => ({
  useAuth: vi.fn(),
}));

vi.mock('@/services/auth.service', () => ({
  authService: {
    getPolicies: vi.fn(),
  },
}));

// Create a mock consumer component to test the context
const TestConsumer = () => {
  const { can, enforcer } = useCasbin();
  
  return (
    <div>
      <div data-testid="enforcer-status">{enforcer ? 'ready' : 'null'}</div>
      <div data-testid="can-read-dashboard">{can(AUTH_ACTIONS.READ, '/v1/dashboard') ? 'yes' : 'no'}</div>
      <div data-testid="can-write-users">{can(AUTH_ACTIONS.WRITE, '/v1/users') ? 'yes' : 'no'}</div>
    </div>
  );
};

describe('CasbinProvider and useCasbin', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockAuthContext = (overrides: Partial<AuthContextType> = {}): AuthContextType => ({
    user: null,
    isAuthenticated: false,
    isLoading: false,
    token: null,
    login: vi.fn(),
    logout: vi.fn(),
    refreshSession: vi.fn(),
    ...overrides,
  });

  it('renders children and defaults to null enforcer when unauthenticated', async () => {
    vi.mocked(useAuth).mockReturnValue(mockAuthContext());

    render(
      <CasbinProvider>
        <TestConsumer />
      </CasbinProvider>
    );

    expect(screen.getByTestId('enforcer-status')).toHaveTextContent('null');
    expect(screen.getByTestId('can-read-dashboard')).toHaveTextContent('no');
  });

  it('initializes enforcer and evaluates policies when authenticated', async () => {
    vi.mocked(useAuth).mockReturnValue(mockAuthContext({
      user: { id: 'u1', role: 'ADMIN', email: 'admin@test.com', name: 'Admin' },
      isAuthenticated: true,
    }));

    vi.mocked(authService.getPolicies).mockResolvedValue({
      policies: [
        `p, ADMIN, /v1/*, ${AUTH_ACTIONS.MANAGE}`,
        `p, ADMIN, ${AUTH_RESOURCES.SYSTEM_SERVICE}, ${AUTH_ACTIONS.MANAGE}`
      ]
    });

    render(
      <CasbinProvider>
        <TestConsumer />
      </CasbinProvider>
    );

    // Wait for async Casbin initialization
    await waitFor(() => {
      expect(screen.getByTestId('enforcer-status')).toHaveTextContent('ready');
    });

    // Check policy enforcement (ADMIN has access to everything)
    expect(screen.getByTestId('can-read-dashboard')).toHaveTextContent('yes');
    expect(screen.getByTestId('can-write-users')).toHaveTextContent('yes');
  });

  it('handles API errors gracefully during policy fetch', async () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    
    vi.mocked(useAuth).mockReturnValue(mockAuthContext({
      user: { id: 'u1', role: 'ADMIN', email: 'admin@test.com', name: 'Admin' },
      isAuthenticated: true,
    }));

    vi.mocked(authService.getPolicies).mockRejectedValue(new Error('Network error'));

    render(
      <CasbinProvider>
        <TestConsumer />
      </CasbinProvider>
    );

    // Give it a moment to fail
    await new Promise(r => setTimeout(r, 100));

    expect(screen.getByTestId('enforcer-status')).toHaveTextContent('null');
    expect(screen.getByTestId('can-read-dashboard')).toHaveTextContent('no');
    expect(consoleSpy).toHaveBeenCalledWith('[Casbin] Failed to initialize Casbin policies:', expect.any(Error));
    
    consoleSpy.mockRestore();
  });
  
  it('returns false for unauthorized actions even when initialized', async () => {
    vi.mocked(useAuth).mockReturnValue(mockAuthContext({
      user: { id: 'u2', role: 'STORE_MGR', email: 'store@test.com', name: 'Store Mgr' },
      isAuthenticated: true,
    }));

    vi.mocked(authService.getPolicies).mockResolvedValue({
      policies: [
        `p, STORE_MGR, /v1/dashboard, ${AUTH_ACTIONS.READ}`
      ]
    });

    render(
      <CasbinProvider>
        <TestConsumer />
      </CasbinProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('enforcer-status')).toHaveTextContent('ready');
    });

    // Can read dashboard but cannot write users
    expect(screen.getByTestId('can-read-dashboard')).toHaveTextContent('yes');
    expect(screen.getByTestId('can-write-users')).toHaveTextContent('no');
  });

  it('enforces FARM_MANAGER policies correctly with path matching', async () => {
    vi.mocked(useAuth).mockReturnValue(mockAuthContext({
      user: { id: 'm1', role: 'FARM_MANAGER', email: 'manager@test.com', name: 'Farm Mgr' },
      isAuthenticated: true,
    }));

    vi.mocked(authService.getPolicies).mockResolvedValue({
      policies: [
        `p, FARM_MANAGER, /v1/inventory/*, ${AUTH_ACTIONS.READ}`,
        `p, FARM_MANAGER, /v1/dashboard, ${AUTH_ACTIONS.READ}`
      ]
    });

    render(
      <CasbinProvider>
        <TestConsumer />
      </CasbinProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('enforcer-status')).toHaveTextContent('ready');
    });

    expect(screen.getByTestId('can-read-dashboard')).toHaveTextContent('yes');
    expect(screen.getByTestId('can-write-users')).toHaveTextContent('no');
  });
});
