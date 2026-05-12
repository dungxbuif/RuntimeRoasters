'use client';

import { KratosForm } from '@/components/auth/KratosForm';
import { LoginCard } from '@/components/auth/LoginCard';
import { AUTH_PARAMS } from '@/constants/auth';
import { APP_ROUTES } from '@/constants/routes';
import { useAcceptHydraLogin, useLoginFlow, useSubmitLogin } from '@/hooks/useAuthFlow';
import { useAuth } from '@/lib/auth';
import { testId, e2eSelectors } from '@/lib/utils/test-id';
import { UpdateLoginFlowBody } from '@ory/client';
import { useRouter, useSearchParams } from 'next/navigation';
import React, { Suspense, useRef, useEffect } from 'react';

const LOG = (...args: unknown[]) => console.log('%c[AUTH]', 'color:#38bdf8;font-weight:bold', ...args);

function LoginContent() {
  const { isAuthenticated } = useAuth();
  const searchParams = useSearchParams();
  const router = useRouter();
  const loginChallenge = searchParams.get(AUTH_PARAMS.LOGIN_CHALLENGE);

  LOG('LoginContent mount', { loginChallenge, isAuthenticated });

  // If already authenticated and no challenge, go home
  useEffect(() => {
    if (isAuthenticated && !loginChallenge) {
      router.push(APP_ROUTES.HOME);
    }
  }, [isAuthenticated, loginChallenge, router]);

  const { data: flow, error: flowError, refetch } = useLoginFlow(undefined, loginChallenge || undefined);
  const submitLogin = useSubmitLogin();
  const acceptHydra = useAcceptHydraLogin();
  const acceptHydraRef = useRef(acceptHydra);
  
  // Use effect to update ref safely
  useEffect(() => {
    acceptHydraRef.current = acceptHydra;
  }, [acceptHydra]);

  // Guard: only handle a given flowError once to prevent re-entrancy
  const flowErrorHandled = useRef(false);

  LOG('flow state', { flowId: flow?.id, hasError: !!flowError, errorId: (flowError as { response?: { data?: { error?: { id: string } } } })?.response?.data?.error?.id });

  // Handle cases where flow is missing but there's no explicit error yet (e.g. empty response)
  useEffect(() => {
    if (!flow && !submitLogin.isPending && !acceptHydra.isPending && !flowError) {
      const timer = setTimeout(() => {
        LOG('Flow is missing for too long, checking session...');
        import('@/services/auth.service').then(({ authService }) => {
          authService.getSession().then((session) => {
            if (session.identity?.id && loginChallenge) {
              LOG('Found existing session, completing Hydra flow...');
              acceptHydraRef.current.mutateAsync({
                challenge: loginChallenge,
                subject: session.identity.id,
              }).then(({ redirect_to }) => {
                window.location.href = redirect_to;
              });
            } else if (session.identity?.id) {
              router.push(APP_ROUTES.DASHBOARD.USERS);
            }
          }).catch(() => {
            LOG('No session found, flow creation might be failed');
          });
        });
      }, 2000); // Wait 2s before fallback
      return () => clearTimeout(timer);
    }
  }, [flow, flowError, loginChallenge, router, submitLogin.isPending, acceptHydra.isPending]);

  const handleLogin = async (body: UpdateLoginFlowBody) => {
    LOG('handleLogin: submitting to Kratos', { flowId: flow?.id, loginChallenge });
    try {
      const session = await submitLogin.mutateAsync({
        flowId: flow?.id as string,
        body
      });

      LOG('handleLogin: Kratos returned session directly', { identityId: session.identity?.id });

      // When loginChallenge was NOT in the flow, Kratos returns the session directly.
      // Accept Hydra login manually and follow the redirect.
      if (loginChallenge && session.identity?.id) {
        LOG('handleLogin: accepting Hydra login', { loginChallenge, subject: session.identity.id });
        const { redirect_to } = await acceptHydra.mutateAsync({
          challenge: loginChallenge,
          subject: session.identity.id,
        });
        LOG('handleLogin: Hydra accept → redirect_to', redirect_to);
        if (redirect_to) {
          window.location.href = redirect_to;
          return;
        }
      }

      LOG('handleLogin: no challenge, pushing to DASHBOARD');
      router.push(APP_ROUTES.DASHBOARD.USERS);
    } catch (err: unknown) {
      console.error('[AUTH] handleLogin error:', err);
      if (err && typeof err === 'object' && 'response' in err) {
        const axiosError = err as { response: { status?: number; data: { ui?: unknown; redirect_browser_to?: string; error?: { id: string } } } };

        LOG('handleLogin: error response', {
          status: axiosError.response?.status,
          errorId: axiosError.response?.data?.error?.id,
          redirect_browser_to: axiosError.response?.data?.redirect_browser_to,
          hasUi: !!axiosError.response?.data?.ui,
        });

        // Kratos accepted OAuth2 login and signals a browser redirect (422 browser_location_change_required).
        // This is the normal success path when loginChallenge is embedded in the flow.
        if (axiosError.response?.data?.redirect_browser_to) {
          LOG('handleLogin: following Kratos redirect (422 normal path) →', axiosError.response.data.redirect_browser_to);
          window.location.href = axiosError.response.data.redirect_browser_to;
          return;
        }

        // Form validation error — re-fetch flow to show updated UI
        if (axiosError.response?.data?.ui) {
          LOG('handleLogin: form validation error, refetching flow');
          refetch();
        }
      }
    }
  };

  useEffect(() => {
    // Guard: prevent re-entrancy. acceptHydra.mutateAsync changing mutation state would
    // otherwise re-trigger this effect on every render → infinite loop.
    if (!flowError || flowErrorHandled.current) return;
    flowErrorHandled.current = true;

    const err = flowError as { response?: { status?: number; data?: { error?: { id: string }, redirect_browser_to?: string } } };
    const errorId = err.response?.data?.error?.id;
    const redirectTo = err.response?.data?.redirect_browser_to;
    const status = err.response?.status;

    LOG('flowError effect fired', { errorId, status, redirectTo, loginChallenge });

    if (errorId === 'session_already_available') {
      LOG('flowError: session already available');
      if (loginChallenge) {
        LOG('flowError: has challenge → fetching session to accept Hydra login');
        import('@/services/auth.service').then(({ authService }) => {
          authService.getSession().then((session) => {
            LOG('flowError: got session', { identityId: session.identity?.id });
            if (session.identity?.id) {
              acceptHydraRef.current.mutateAsync({
                challenge: loginChallenge,
                subject: session.identity.id,
              }).then(({ redirect_to }) => {
                LOG('flowError: Hydra accepted → redirect_to', redirect_to);
                window.location.href = redirect_to;
              });
            }
          }).catch((sessionErr) => {
            console.error('[AUTH] getSession error:', sessionErr);
            alert("Lỗi không lấy được Session từ Kratos. Vui lòng xóa cookie Kratos và thử lại.");
          });
        });
      } else {
        LOG('flowError: no challenge → pushing DASHBOARD');
        router.push(APP_ROUTES.DASHBOARD.USERS);
      }
    }
    if (redirectTo) {
      LOG('flowError: Kratos redirect →', redirectTo);
      window.location.href = redirectTo;
    } else if (status === 422 && errorId === 'browser_location_change_required' && redirectTo) {
      LOG('flowError: 422 browser_location_change_required →', redirectTo);
      window.location.href = redirectTo;
    } else {
      LOG('flowError: unhandled error', err.response?.data);
    }
  }, [flowError, router, loginChallenge]);

  if (!flow) return (
    <div className="flex items-center justify-center min-h-screen bg-[#f7f9fb]">
      <div className="text-primary font-black uppercase italic animate-pulse tracking-[0.5em] text-xs">
        Initializing Secure Terminal...
      </div>
    </div>
  );

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-[#f7f9fb]">
      <LoginCard title="System Access" subtitle="Please authenticate to manage cluster nodes">
        <KratosForm 
          flow={flow} 
          onSubmit={handleLogin} 
          isLoading={submitLogin.isPending || acceptHydra.isPending} 
        />
      </LoginCard>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <LoginContent />
    </Suspense>
  );
}
