import { Configuration, FrontendApi } from "@ory/client";
import { ENV } from "@/constants/env";

const kratos = new FrontendApi(
  new Configuration({
    basePath: ENV.KRATOS_PUBLIC_URL,
    baseOptions: {
      withCredentials: true,
    },
  })
);

export default kratos;
