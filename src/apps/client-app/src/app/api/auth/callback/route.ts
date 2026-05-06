import { AUTH_PARAMS, OIDC_CONFIG } from "@/constants/auth";
import { ENV } from "@/constants/env";
import { APP_ROUTES } from "@/constants/routes";
import { buildFormUrlEncoded, getCallbackUrl, joinPaths } from "@/lib/utils/url";
import axios from "axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  const searchParams = req.nextUrl.searchParams;
  const code = searchParams.get("code");
  const error = searchParams.get("error");
  const errorDescription = searchParams.get("error_description");

  console.log('[CALLBACK] received', {
    url: req.nextUrl.toString(),
    hasCode: !!code,
    error,
    errorDescription,
  });

  if (!code) {
    console.error('[CALLBACK] no code — Hydra error:', error, errorDescription);
    return NextResponse.redirect(new URL(APP_ROUTES.LOGIN, req.url));
  }

  try {
    const callbackUrl = getCallbackUrl();
    const hydraUrl = ENV.HYDRA_PUBLIC_URL;

    console.log('[CALLBACK] exchanging code for token', { hydraUrl, callbackUrl, clientId: OIDC_CONFIG.CLIENT_ID });

    const payload = {
      grant_type: OIDC_CONFIG.GRANT_TYPE,
      code,
      redirect_uri: callbackUrl,
      client_id: OIDC_CONFIG.CLIENT_ID,
      client_secret: OIDC_CONFIG.CLIENT_SECRET,
    };

    const response = await axios.post(
      joinPaths(hydraUrl, '/oauth2/token'),
      buildFormUrlEncoded(payload),
      { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
    );

    const tokens = response.data;

    if (tokens.error) {
      console.error('[CALLBACK] token exchange error:', tokens);
      return NextResponse.redirect(new URL(`${APP_ROUTES.LOGIN}?error=token_exchange_failed`, req.url));
    }

    console.log('[CALLBACK] token exchange success → redirecting to HOME with token');

    const successUrl = new URL(APP_ROUTES.HOME, req.url);
    successUrl.searchParams.set(AUTH_PARAMS.ACCESS_TOKEN, tokens.access_token);
    return NextResponse.redirect(successUrl);

  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    console.error('[CALLBACK] exception:', errorMessage);
    return NextResponse.redirect(new URL(`${APP_ROUTES.LOGIN}?error=callback_failed`, req.url));
  }
}
