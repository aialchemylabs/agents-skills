---
name: configuring-litellm-entra-sso
description: Configure, repair, rotate, audit, or document direct Microsoft Entra OIDC SSO for a self-hosted LiteLLM AI Gateway. Use for single-tenant app registration, exact /sso/callback redirects, Entra app roles and user assignment, protected client-secret wiring, Docker Compose configuration, UI hardening, licensing checks, or SSO verification. Do not use for SAML, SCIM provisioning, Cloudflare Access, or API-client authentication.
compatibility: Requires access to the LiteLLM deployment, Microsoft Entra administration, an HTTPS gateway hostname, and current LiteLLM and Microsoft documentation.
metadata:
  author: aialchemylabs
  version: "0.1.0"
  domain: ai-gateway
  category: identity-and-access
  keywords: litellm,microsoft-entra,oidc,sso,docker-compose,ai-gateway
  validated-with: litellm-1.93.0
---

# Configure LiteLLM Microsoft Entra SSO

Configure Microsoft Entra directly as LiteLLM's identity provider for human Admin UI access. Keep machine clients on scoped LiteLLM virtual keys.

## Preserve these architecture decisions

- Use LiteLLM's native Microsoft OIDC integration.
- Do not make `prototypes-aialchemylabs` an identity provider. Treat a future `aialchemylabs-auth` service as a deferred alternative only when the user explicitly requests that architecture.
- Decide and confirm the production HTTPS hostname and ingress path before creating redirect-dependent configuration.
- Allow a transport tunnel or reverse proxy, but do not add Cloudflare Access when direct Entra SSO is the requested access control.
- Bind the origin to loopback or a private container network. Do not expose the origin independently of the approved ingress.
- Use Entra only for human UI login. Require scoped LiteLLM virtual keys for local, container, and remote API clients.
- Keep the LiteLLM master key for break-glass administration. Never distribute it to applications.
- Do not configure SAML or SCIM. A provisioning warning in Entra is expected because OIDC sign-in and SCIM lifecycle provisioning are separate features.

## Check current state before changing it

1. Inspect the deployed LiteLLM version, license state, Compose or orchestrator configuration, current hostname, ingress, and health.
2. Identify existing Cloudflare Access applications, reverse-proxy auth, or other identity layers that would create double authentication.
3. Confirm the tenant, allowed email domain, application display name, intended administrators, viewers, and internal users.
4. Confirm the exact public base URL without a trailing slash. Derive the callback as `<base-url>/sso/callback`.
5. Read the current official sources before mutation:
   - [LiteLLM Admin UI SSO](https://docs.litellm.ai/docs/proxy/admin_ui_sso)
   - [LiteLLM proxy configuration settings](https://docs.litellm.ai/docs/proxy/config_settings)
   - [Microsoft application objects and service principals](https://learn.microsoft.com/en-us/entra/identity-platform/app-objects-and-service-principals)
   - [Microsoft redirect URI configuration](https://learn.microsoft.com/en-us/entra/identity-platform/how-to-add-redirect-uri)
   - [Microsoft app roles](https://learn.microsoft.com/en-us/entra/identity-platform/howto-add-app-roles-in-apps)
6. Report the licensing boundary before rollout. LiteLLM currently documents SSO as free for up to five users from v1.76.0; require an Enterprise license or written confirmation from LiteLLM before exceeding the documented free allowance.

Do not make paid inference requests while discovering or validating configuration.

## Obtain action-time confirmation

Ask immediately before any browser, CLI, or Graph action that:

- creates a client secret or other credential;
- grants persistent application access or tenant consent;
- assigns users or groups to privileged app roles;
- deletes an existing credential, Cloudflare Access application, or access policy.

A prior general approval is not a substitute when the credential or persistent grant is about to be created. Continue read-only inspection and local preparation while waiting.

## Create the Entra application

1. Open **Microsoft Entra ID > App registrations > New registration**. Do not create this OIDC integration from **Browse Microsoft Entra App Gallery**.
2. Enter the approved display name.
3. Select **Accounts in this organizational directory only** for a single-tenant application.
4. Select the **Web** platform and enter exactly `<base-url>/sso/callback`.
5. Register the application and record only:
   - Application (client) ID
   - Directory (tenant) ID
   - Application object ID when automation needs it
6. Verify the minimum delegated Microsoft Graph permission required by the current LiteLLM Microsoft user-info flow. Prefer `User.Read`; do not add broad directory or application permissions speculatively.

An App Registration automatically creates a service principal in the tenant. Manage sign-in assignment under **Enterprise applications**. If it is not visible, clear list filters, search by Application ID, and allow for portal propagation. Do not create a duplicate gallery application.

## Define LiteLLM app roles

Create these enabled roles under **App registrations > App roles** with **Users/Groups** as the allowed member type. The role **Value** must match exactly:

| Display name | Value | Access |
|---|---|---|
| LiteLLM Proxy Admin | `proxy_admin` | Full platform administration |
| LiteLLM Proxy Admin Viewer | `proxy_admin_viewer` | Read-only global administration |
| LiteLLM Internal User | `internal_user` | Normal user and self-service key access subject to team permissions |

Use a distinct generated role ID for each role. Do not use Entra directory roles such as Global Administrator as a substitute for application roles.

Then open the corresponding **Enterprise application**:

1. Set **Assignment required?** to **Yes**.
2. Assign only approved users or groups under **Users and groups**.
3. Give full administrators `proxy_admin`; use least privilege for everyone else.
4. Verify each assignment by user object ID, application role ID, and resulting role name.

Entra emits assigned app-role values in the token's `roles` claim. LiteLLM treats those values as authoritative on sign-in.

## Create and store the client secret

After action-time confirmation:

1. Create a client secret with the shortest operationally practical expiry and record its expiry date.
2. Copy the one-time **Value**, not the Secret ID.
3. Never print the value, paste it into chat, pass it as a command-line argument, or commit it.
4. Feed the value through standard input to a repository helper or secret manager.
5. Store the runtime copy as a mode-0600 Docker secret or equivalent protected secret.
6. Mirror it to `.env` only when the repository's vault workflow explicitly requires that copy. Ensure `.env` is ignored and mode 0600 before writing.
7. Make secret updates atomic. Preserve the existing credential until the replacement is deployed and verified, then request confirmation before deletion.

Never put a secret value in `.env.example`, Compose YAML, documentation, shell history, logs, diffs, or process arguments.

## Configure LiteLLM

Set these non-secret values for the gateway:

```text
PROXY_BASE_URL=https://gateway.example.com
MICROSOFT_CLIENT_ID=<application-client-id>
MICROSOFT_TENANT=<directory-tenant-id>
ALLOWED_EMAIL_DOMAINS=example.com
AUTO_REDIRECT_UI_LOGIN_TO_SSO=true
NO_DOCS=true
NO_OPENAPI=true
NO_REDOC=true
ROOT_REDIRECT_URL=/ui/
```

Treat automatic SSO redirect as optional. Keep it enabled for an Entra-only UI unless the user wants the login chooser. Preserve `/fallback/login` for break-glass master-key access.

Mount the client secret without interpolating it into rendered Compose output:

```yaml
services:
  llm-gateway:
    environment:
      PROXY_BASE_URL: ${LITELLM_PUBLIC_URL:?set LITELLM_PUBLIC_URL}
      MICROSOFT_CLIENT_ID: ${MICROSOFT_CLIENT_ID:?set MICROSOFT_CLIENT_ID}
      MICROSOFT_TENANT: ${MICROSOFT_TENANT:?set MICROSOFT_TENANT}
      ALLOWED_EMAIL_DOMAINS: ${ALLOWED_EMAIL_DOMAINS:?set ALLOWED_EMAIL_DOMAINS}
      AUTO_REDIRECT_UI_LOGIN_TO_SSO: ${AUTO_REDIRECT_UI_LOGIN_TO_SSO:-true}
      NO_DOCS: "true"
      NO_OPENAPI: "true"
      NO_REDOC: "true"
      ROOT_REDIRECT_URL: /ui/
    secrets:
      - litellm_microsoft_client_secret
    command:
      - /bin/sh
      - -ec
      - |
        export MICROSOFT_CLIENT_SECRET="$$(cat /run/secrets/litellm_microsoft_client_secret)"
        exec litellm --config /app/config.yaml --port 4000 --host 0.0.0.0

secrets:
  litellm_microsoft_client_secret:
    file: ./.secrets/litellm-microsoft-client-secret
```

Adapt service names and startup arguments to the repository. Do not replace existing hardening, health checks, resource limits, or immutable image references.

If removing Cloudflare Access, remove only the application and policies for the approved LiteLLM hostname. Retain the tunnel and DNS route when they are the selected HTTPS transport. Verify no unrelated Cloudflare resources are affected.

## Deploy safely

1. Validate configuration without revealing resolved secrets.
2. Check syntax and scan both unstaged and staged diffs for the actual secret value without printing matches.
3. Recreate the gateway; a restart alone may retain the previous Compose environment or secret.
4. Wait for the configured health check to become healthy.
5. Inspect logs for SSO initialization and errors without dumping environment variables.
6. Update operational documentation with the base URL, callback, roles, assignment model, secret rotation method, licensing threshold, and break-glass route.

## Verify without inference

Verify locally and through the public hostname:

1. `/health/liveliness` returns the expected healthy response.
2. An unauthenticated `/v1/models` request returns `401` from LiteLLM both locally and publicly.
3. `/sso/key/generate` redirects to `login.microsoftonline.com/<tenant-id>/...` and uses the exact registered callback.
4. `/.well-known/litellm-ui-config` reports SSO configured and, when enabled, automatic SSO redirect.
5. `/` redirects to `/ui/`.
6. `/docs`, `/openapi.json`, and `/redoc` return `404` when public schema surfaces are disabled.
7. Unauthenticated key, user, spend, model, and management API requests return `401`.
8. A nested `/ui/...` route may return a static frontend shell. Confirm it contains no model aliases or private data and that its API calls remain authenticated.
9. In a private browser session, an assigned administrator signs in through Entra, returns through `/sso/callback`, and receives `proxy_admin`.
10. Verify an unassigned user is denied when assignment is required.
11. Verify a scoped LiteLLM virtual key still works for an approved API client without relying on an Entra browser session. Avoid a paid model call; prefer authenticated model listing or another non-inference endpoint.

Record exact status codes, redirect hosts, role results, LiteLLM version, and license state. Do not claim completion until the remote SSO flow and API boundary have both been verified.

## Troubleshoot narrowly

- `AADSTS50011`: compare the runtime callback and registered Web redirect character-for-character, including scheme, host, path, and trailing slash.
- Application absent from Enterprise applications: search by Application ID and remove filters; the service principal is created from the App Registration.
- OIDC page shows a SCIM provisioning warning: ignore it unless lifecycle provisioning is separately required.
- User authenticates but lacks the intended role: verify Enterprise Application assignment, exact role value, token `roles` claim, and a fresh login after clearing cached sessions.
- UI redirects to the wrong host: correct `PROXY_BASE_URL`, recreate the gateway, and verify the generated authorization request.
- Static nested UI page returns `200`: inspect whether it is only the frontend shell, then verify all backing management APIs return `401`.
- More than five SSO users are required: stop and resolve LiteLLM Enterprise licensing before onboarding additional users.
