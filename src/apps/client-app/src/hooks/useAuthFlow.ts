import { useMutation, useQuery } from "@tanstack/react-query";
import { authService } from "@/services/auth.service";
import { UpdateLoginFlowBody } from "@ory/client";

export const useLoginFlow = (returnTo?: string, loginChallenge?: string) => {
  return useQuery({
    queryKey: ['kratos-login-flow', returnTo, loginChallenge],
    queryFn: () => authService.createLoginFlow(returnTo, loginChallenge),
    staleTime: 0, // Always fresh for auth flows
    retry: false, // Don't retry on 400 errors
  });
};

export const useSubmitLogin = () => {
  return useMutation({
    mutationFn: ({ flowId, body }: { flowId: string, body: UpdateLoginFlowBody }) => 
      authService.submitLoginFlow(flowId, body),
  });
};

export const useAcceptHydraLogin = () => {
  return useMutation({
    mutationFn: ({ challenge, subject }: { challenge: string, subject: string }) =>
      authService.acceptHydraLogin(challenge, subject),
  });
};
