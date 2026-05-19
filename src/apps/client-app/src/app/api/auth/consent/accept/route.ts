import { hydraAdmin } from "@/lib/ory/hydra";
import { resolveIdentityClaims } from "@/lib/ory/identity-admin";
import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    const { consent_challenge } = await req.json();

    console.log('[CONSENT/ACCEPT] received', { consent_challenge });

    if (!consent_challenge) {
      console.warn('[CONSENT/ACCEPT] missing consent_challenge');
      return NextResponse.json({ error: "Missing consent_challenge" }, { status: 400 });
    }

    const { data: consentRequest } = await hydraAdmin.getOAuth2ConsentRequest({
      consentChallenge: consent_challenge,
    });
    const claims = await resolveIdentityClaims(consentRequest.subject);

    console.log('[CONSENT/ACCEPT] consent request', {
      client_id: consentRequest.client?.client_id,
      requested_scope: consentRequest.requested_scope,
      role: claims.role,
    });

    const { data: acceptResponse } = await hydraAdmin.acceptOAuth2ConsentRequest({
      consentChallenge: consent_challenge,
      acceptOAuth2ConsentRequest: {
        grant_scope: consentRequest.requested_scope,
        grant_access_token_audience: consentRequest.requested_access_token_audience,
        session: {
          id_token: claims,
          access_token: claims,
        },
      },
    });

    console.log('[CONSENT/ACCEPT] Hydra accepted → redirect_to', acceptResponse.redirect_to);
    return NextResponse.json({ redirect_to: acceptResponse.redirect_to });
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    console.error('[CONSENT/ACCEPT] error:', errorMessage);
    return NextResponse.json({ error: "Failed to accept consent" }, { status: 500 });
  }
}
