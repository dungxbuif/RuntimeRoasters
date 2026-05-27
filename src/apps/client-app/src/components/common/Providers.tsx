'use client';

import { AuthProvider } from '@/lib/auth';
import { CasbinProvider } from '@/lib/auth/casbin';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { useState } from 'react';

export default function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,
        retry: 1,
      },
    },
  }));

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <CasbinProvider>
          {children}
        </CasbinProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}
