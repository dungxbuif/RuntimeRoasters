'use client';

import { KratosForm } from '@/components/auth/KratosForm';
import { LoginCard } from '@/components/auth/LoginCard';
import { AUTH_PARAMS } from '@/constants/auth';
import { APP_ROUTES } from '@/constants/routes';
import { useAcceptHydraLogin, useLoginFlow, useSubmitLogin } from '@/hooks/useAuthFlow';
import { UpdateLoginFlowBody } from '@ory/client';
import { useRouter, useSearchParams } from 'next/navigation';
import React, { Suspense, useRef } from 'react';

const LOG = (...args: unknown[]) => console.log('%c[AUTH]', 'color:#38bdf8;font-weight:bold', ...args);

function LoginContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const loginChallenge = searchParams.get(AUTH_PARAMS.LOGIN_CHALLENGE);

  LOG('LoginContent mount', { loginChallenge, url: typeof window !== 'undefined' ? window.location.href : '' });

  // When loginChallenge is present, Kratos uses it to know where to redirect after auth.
  // Passing returnTo = current login URL would cause Kratos to redirect back here.
  const { data: flow, error: flowError, refetch } = useLoginFlow(undefined, loginChallenge || undefined);
  const submitLogin = useSubmitLogin();
  const acceptHydra = useAcceptHydraLogin();
  // Stable ref so the useEffect below doesn't list acceptHydra as a dep (would cause a loop)
  const acceptHydraRef = useRef(acceptHydra);
  acceptHydraRef.current = acceptHydra;

  // Guard: only handle a given flowError once to prevent re-entrancy
  const flowErrorHandled = useRef(false);

  LOG('flow state', { flowId: flow?.id, hasError: !!flowError, errorId: (flowError as any)?.response?.data?.error?.id });

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

      LOG('handleLogin: no challenge, pushing to HOME');
      router.push(APP_ROUTES.HOME);
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

  React.useEffect(() => {
    // Guard: prevent re-entrancy. acceptHydra.mutateAsync changing mutation state would
    // otherwise re-trigger this effect on every render → infinite loop.
    if (!flowError || flowErrorHandled.current) return;
    flowErrorHandled.current = true;

    const err = flowError as any;
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
        LOG('flowError: no challenge → pushing HOME');
        router.push(APP_ROUTES.HOME);
      }
    } else if (redirectTo) {
      LOG('flowError: Kratos redirect →', redirectTo);
      window.location.href = redirectTo;
    } else if (status === 422 && errorId === 'browser_location_change_required') {
      LOG('flowError: 422 browser_location_change_required →', redirectTo);
      window.location.href = redirectTo;
    } else {
      LOG('flowError: unhandled error', err.response?.data);
    }
  // acceptHydra intentionally omitted from deps — mutateAsync is stable, and including
  // the mutation object would re-trigger this effect on every state change → infinite loop.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [flowError, router, loginChallenge]);

  if (!flow) return (
    <div className="flex items-center justify-center min-h-screen bg-slate-950">
      <div className="text-primary font-black uppercase italic animate-pulse tracking-[0.5em]">
        Initializing Secure Terminal...
      </div>
    </div>
  );

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-[radial-gradient(circle_at_top_right,_var(--tw-gradient-stops))] from-slate-900 via-slate-950 to-blue-900/20">
      <LoginCard title="System Login" subtitle="">
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
