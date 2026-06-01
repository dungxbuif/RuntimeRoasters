import { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import { authService } from '@/services/auth.service';
import * as casbin from 'casbin-core';
import { useAuth } from './AuthProvider';
import { CASBIN_MODEL } from '@/constants/casbin';

// Helper to handle ESM/CJS interop for casbin-core
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getCasbinCore = (lib: any) => {
  if (lib.newEnforcer) return lib;
  if (lib.default && lib.default.newEnforcer) return lib.default;
  return lib;
};

interface CasbinContextType {
  enforcer: casbin.Enforcer | null;
  can: (action: string, resource: string) => boolean;
  isLoading: boolean;
}

const CasbinContext = createContext<CasbinContextType>({
  enforcer: null,
  can: () => false,
  isLoading: true,
});

export const CasbinProvider = ({ children }: { children: ReactNode }) => {
  const { user, isAuthenticated } = useAuth();
  const [enforcer, setEnforcer] = useState<casbin.Enforcer | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initCasbin = async () => {
      setIsLoading(true);
      if (isAuthenticated && user?.role) {
        try {
          const { policies } = await authService.getPolicies();
          
          const core = getCasbinCore(casbin);
          const m = new core.Model(CASBIN_MODEL);
          const a = new core.MemoryAdapter(policies.join('\n'));
          const e = await core.newEnforcer(m, a);

          // Register keyMatch and regexMatch functions (they are synchronous)
          e.addFunction('keyMatch', core.Util.keyMatchFunc);
          e.addFunction('regexMatch', core.Util.regexMatchFunc);

          setEnforcer(e);
        } catch (error) {
          console.error('[Casbin] Failed to initialize Casbin policies:', error);
        } finally {
          setIsLoading(false);
        }
      } else {
        setEnforcer(null);
        setIsLoading(false);
      }
    };

    initCasbin();
  }, [isAuthenticated, user?.role]);

  const can = (action: string, resource: string): boolean => {
    if (!enforcer || !user) return false;
    // We enforce with the user's role as the subject
    return enforcer.enforceSync(user.role, resource, action);
  };

  return (
    <CasbinContext.Provider value={{ enforcer, can, isLoading }}>
      {children}
    </CasbinContext.Provider>
  );
};

export const useCasbin = () => useContext(CasbinContext);

export const PermissionGuard = ({
  action,
  resource,
  children,
  fallback = null,
}: {
  action: string;
  resource: string;
  children: ReactNode;
  fallback?: ReactNode;
}) => {
  const { can, enforcer } = useCasbin();
  
  if (!enforcer) {
    return <>{fallback}</>;
  }

  if (can(action, resource)) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
};
