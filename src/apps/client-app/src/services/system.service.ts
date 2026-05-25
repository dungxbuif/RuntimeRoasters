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

export const systemService = {
  getLogisticsStatus: async () => {
    const response = await api.get<SystemStatus>('/v1/logistics/system/status');
    return response.data;
  },

  seedLogistics: async (force: boolean = false) => {
    const response = await api.post<SeedResult>('/v1/logistics/system/seed', { force });
    return response.data;
  },

  getRetailStatus: async () => {
    const response = await api.get<SystemStatus>('/v1/retail/system/status');
    return response.data;
  },

  seedRetail: async (force: boolean = false) => {
    const response = await api.post<SeedResult>('/v1/retail/system/seed', { force });
    return response.data;
  },

  // Add more services as they are implemented
  checkAllStatus: async () => {
    // For now logistics and retail
    try {
      console.debug('[SystemService] Checking status for logistics and retail...');
      const [logistics, retail] = await Promise.all([
        systemService.getLogisticsStatus(),
        systemService.getRetailStatus()
      ]);
      console.debug('[SystemService] Logistics:', logistics.seeded, 'Retail:', retail.seeded);
      return {
        allSeeded: logistics.seeded && retail.seeded,
        services: { logistics, retail }
      };
    } catch (err) {
      console.error('[SystemService] Failed to check system status:', err);
      return { allSeeded: false, services: {} }; // Show modal if check fails to be safe
    }
  },

  seedAll: async () => {
    const results = await Promise.allSettled([
      systemService.seedLogistics(),
      systemService.seedRetail()
    ]);
    return results;
  }
};
