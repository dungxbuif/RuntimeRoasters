import { ENV } from "@/constants/env";
import { Configuration, OAuth2Api } from "@ory/client";

export const hydraPublic = new OAuth2Api(
  new Configuration({
    basePath: ENV.HYDRA_PUBLIC_URL,
    baseOptions: {
      withCredentials: true,
    },
  })
);

export const hydraAdmin = new OAuth2Api(
  new Configuration({
    basePath: ENV.HYDRA_ADMIN_URL,
  })
);
