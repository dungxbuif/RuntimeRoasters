import { NextRequest, NextResponse } from "next/server";
import { hydraAdmin } from "@/lib/ory/hydra";

export async function POST(req: NextRequest) {
  try {
    const { login_challenge, subject } = await req.json();

    console.log('[LOGIN/ACCEPT] received', { login_challenge, subject });

    if (!login_challenge || !subject) {
      console.warn('[LOGIN/ACCEPT] missing params');
      return NextResponse.json({ error: "Missing login_challenge or subject" }, { status: 400 });
    }

    const { data } = await hydraAdmin.acceptOAuth2LoginRequest({
      loginChallenge: login_challenge,
      acceptOAuth2LoginRequest: {
        subject: subject,
        remember: true,
        remember_for: 3600,
      },
    });

    console.log('[LOGIN/ACCEPT] Hydra accepted → redirect_to', data.redirect_to);
    return NextResponse.json({ redirect_to: data.redirect_to });
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    console.error('[LOGIN/ACCEPT] error:', errorMessage);
    return NextResponse.json({ error: "Failed to accept login request" }, { status: 500 });
  }
}
