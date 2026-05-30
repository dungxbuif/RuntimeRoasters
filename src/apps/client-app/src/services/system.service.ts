import api from '@/lib/axios';

export interface SystemStatus {
  seeded: boolean;
  service_name: string;
  record_counts: Record<string, number>;
  last_seeded_at?: string;
}

export interface SeedResult {
  success: boolean;
  message: string;
  records_created: number;
}

export interface UnifiedStatus {
  all_seeded: boolean;
  services_status: Record<string, boolean>;
}

export const systemService = {
  getSystemStatus: async () => {
    const response = await api.get<UnifiedStatus>('/v1/auth/system/status');
    return response.data;
  },

  seedSystem: async (force: boolean = false) => {
    const response = await api.post<SeedResult>('/v1/auth/seed', { force });
    return response.data;
  },

  getLogisticsStatus: async () => {
    const response = await api.get<SystemStatus>('/v1/logistics/system/status');
    return response.data;
  },

  seedRetail: async (force: boolean = false) => {
    const response = await api.post<SeedResult>('/v1/retail/system/seed', { force });
    return response.data;
  },

  checkAllStatus: async () => {
    try {
      console.debug('[SystemService] Checking unified status via auth service...');
      const status = await systemService.getSystemStatus();
      console.debug('[SystemService] Unified Status:', status.all_seeded, status.services_status);
      return {
        allSeeded: status.all_seeded,
        services: status.services_status
      };
    } catch (err) {
      console.error('[SystemService] Failed to check unified system status:', err);
      return { allSeeded: false, services: {} };
    }
  },

  seedAll: async () => {
    try {
      console.debug('[SystemService] Triggering unified seed via auth service...');
      const result = await systemService.seedSystem(true);
      return [result];
    } catch (err) {
      console.error('[SystemService] Unified seed failed:', err);
      throw err;
    }
  }
};
