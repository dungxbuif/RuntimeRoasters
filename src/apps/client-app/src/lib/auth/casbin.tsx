import { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import { authService } from '@/services/auth.service';
import * as casbin from 'casbin-core';
import { useAuth } from './AuthProvider';
import { CASBIN_MODEL } from '@/constants/casbin';

interface CasbinContextType {
  enforcer: casbin.Enforcer | null;
  can: (action: string, resource: string) => boolean;
}

const CasbinContext = createContext<CasbinContextType>({
  enforcer: null,
  can: () => false,
});

export const CasbinProvider = ({ children }: { children: ReactNode }) => {
  const { user, isAuthenticated } = useAuth();
  const [enforcer, setEnforcer] = useState<casbin.Enforcer | null>(null);

  useEffect(() => {
    const initCasbin = async () => {
      if (isAuthenticated && user?.role) {
        try {
          const { policies } = await authService.getPolicies();
          
          const m = new casbin.Model(CASBIN_MODEL);
          const a = new casbin.MemoryAdapter(policies.join('\n'));
          const e = await casbin.newEnforcer(m, a);

          // Register keyMatch and regexMatch functions (they are synchronous)
          e.addFunction('keyMatch', casbin.Util.keyMatchFunc);
          e.addFunction('regexMatch', casbin.Util.regexMatchFunc);

          setEnforcer(e);
        } catch (error) {
          console.error('[Casbin] Failed to initialize Casbin policies:', error);
        }
      } else {
        setEnforcer(null);
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
    <CasbinContext.Provider value={{ enforcer, can }}>
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
