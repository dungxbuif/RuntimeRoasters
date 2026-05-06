import { ENV } from "@/constants/env";
import { API_ENDPOINTS } from "@/constants/api";

/**
 * Join URL paths and handle slashes automatically
 */
export const joinPaths = (...parts: string[]): string => {
  return parts
    .map((part, index) => {
      if (index === 0) {
        return part.replace(/\/+$/, "");
      }
      return part.replace(/^\/+/, "").replace(/\/+$/, "");
    })
    .filter((part) => part.length > 0)
    .join("/");
};

export const getCallbackUrl = () => {
  return joinPaths(ENV.APP_URL, API_ENDPOINTS.AUTH.CALLBACK);
};

export const buildFormUrlEncoded = (data: Record<string, string>) => {
  const params = new URLSearchParams();
  Object.entries(data).forEach(([key, value]) => {
    params.append(key, value);
  });
  return params.toString();
};
