/**
 * Casbin model definition for the frontend.
 * Optimized for synchronous evaluation (enforceSync) in React components.
 * 
 * Matcher Logic:
 * - ADMIN role has full access to everything.
 * - Other roles must match the subject in the policy directly or through inheritance.
 * - Resources match using keyMatch (RESTful paths with *) and regexMatch.
 */
export const CASBIN_MODEL = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.sub == 'ADMIN') || ((r.sub == p.sub || g(r.sub, p.sub)) && (keyMatch(r.obj, p.obj) || regexMatch(r.obj, p.obj)) && regexMatch(r.act, p.act))
`;
