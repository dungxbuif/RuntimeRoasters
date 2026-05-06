import { API_ENDPOINTS } from "@/constants/api";
import { ENV } from "@/constants/env";
import api from "@/lib/axios";
import kratos from "@/lib/ory/kratos";
import { LoginFlow, Session, UpdateLoginFlowBody } from "@ory/client";

class AuthService {
  async createLoginFlow(returnTo?: string, loginChallenge?: string): Promise<LoginFlow> {
    const { data } = await kratos.createBrowserLoginFlow({
      returnTo: returnTo ?? `${ENV.APP_URL}/`,
      loginChallenge,
    });
    return data;
  }

  async getSession(): Promise<Session> {
    const { data } = await kratos.toSession();
    return data;
  }

  async submitLoginFlow(flowId: string, body: UpdateLoginFlowBody): Promise<Session> {
    const { data } = await kratos.updateLoginFlow({
      flow: flowId,
      updateLoginFlowBody: body,
    });
    return data.session as Session;
  }

  async acceptHydraLogin(loginChallenge: string, subject: string): Promise<{ redirect_to: string }> {
    const { data } = await api.post(API_ENDPOINTS.AUTH.LOGIN_ACCEPT, {
      login_challenge: loginChallenge,
      subject,
    });
    return data;
  }

  async acceptHydraConsent(consentChallenge: string): Promise<{ redirect_to: string }> {
    const { data } = await api.post(API_ENDPOINTS.AUTH.CONSENT_ACCEPT, {
      consent_challenge: consentChallenge,
    });
    return data;
  }
}

export const authService = new AuthService();
