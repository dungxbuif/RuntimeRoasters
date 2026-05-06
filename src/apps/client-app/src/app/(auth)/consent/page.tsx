import React from 'react';
import { redirect } from 'next/navigation';
import { hydraAdmin } from "@/lib/ory/hydra";
import { AUTH_PARAMS, TRUSTED_CLIENTS } from '@/constants/auth';

interface ConsentPageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

/**
 * Consent Page - Server Component
 * Thực hiện logic Auto-Accept cho các Trusted Clients ngầm định (Behind the scenes)
 */
export default async function ConsentPage({ searchParams }: ConsentPageProps) {
  const params = await searchParams;
  const consentChallenge = params[AUTH_PARAMS.CONSENT_CHALLENGE] as string;

  if (!consentChallenge) {
    console.error('[CONSENT] Missing consent_challenge in URL');
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-950 text-white">
        <p className="font-mono text-error uppercase tracking-widest">Error: Invalid Authorization Request</p>
      </div>
    );
  }

  console.log('[CONSENT] processing challenge', consentChallenge.substring(0, 20) + '...');

  // redirect() from next/navigation throws a special RedirectError internally.
  // It MUST be called outside try/catch — otherwise the catch block swallows it
  // and the redirect never happens, showing the error UI instead.
  let redirectUrl: string | undefined;
  let consentRequest;

  try {
    const response = await hydraAdmin.getOAuth2ConsentRequest({ consentChallenge });
    consentRequest = response.data;

    console.log('[CONSENT] client:', consentRequest.client?.client_id, 'scopes:', consentRequest.requested_scope);

    const isTrusted = TRUSTED_CLIENTS.includes(consentRequest.client?.client_id || '');

    if (isTrusted) {
      const { data: acceptResponse } = await hydraAdmin.acceptOAuth2ConsentRequest({
        consentChallenge,
        acceptOAuth2ConsentRequest: {
          grant_scope: consentRequest.requested_scope,
          grant_access_token_audience: consentRequest.requested_access_token_audience,
          session: {
            id_token: { role: "ADMIN" },
            access_token: { role: "ADMIN" },
          },
          remember: true,
          remember_for: 3600,
        },
      });

      console.log('[CONSENT] accepted → redirect_to', acceptResponse.redirect_to?.substring(0, 60) + '...');
      redirectUrl = acceptResponse.redirect_to;
    }
  } catch (error) {
    console.error('[CONSENT] Hydra error:', error);
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-950 text-white">
        <p className="font-mono text-error uppercase tracking-widest">Identity Gateway Timeout</p>
      </div>
    );
  }

  // Called OUTSIDE try/catch so the RedirectError propagates correctly
  if (redirectUrl) {
    redirect(redirectUrl);
  }

  // Nếu không thuộc diện Trusted (Fallback UI - hiếm khi xảy ra trong project này)
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 text-white p-8">
      <div className="max-w-md w-full bg-slate-900 p-8 rounded-3xl border border-slate-800 shadow-2xl">
        <h1 className="text-2xl font-black font-headline uppercase italic mb-4">Manual Consent Required</h1>
        <p className="text-slate-400 text-sm mb-8">The client <strong>{consentRequest?.client?.client_id}</strong> is requesting access to your identity.</p>
        <div className="space-y-4">
           {/* UI nút đồng ý/từ chối thủ công nếu cần */}
           <p className="text-xs text-slate-600 italic">Auto-Accept logic failed or client is untrusted.</p>
        </div>
      </div>
    </div>
  );
}
