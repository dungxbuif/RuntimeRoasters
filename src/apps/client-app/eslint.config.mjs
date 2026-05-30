import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    rules: {
      "no-restricted-syntax": [
        "error",
        {
          selector: "JSXOpeningElement[name.name='CasbinGuard'] > JSXAttribute[name.name='obj'] > Literal",
          message: "Do not use magic strings in CasbinGuard's 'obj' prop. Use constants from AUTH_RESOURCES."
        },
        {
          selector: "JSXOpeningElement[name.name='CasbinGuard'] > JSXAttribute[name.name='act'] > Literal",
          message: "Do not use magic strings in CasbinGuard's 'act' prop. Use constants from AUTH_ACTIONS."
        }
      ]
    }
  },
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
]);

export default eslintConfig;
