# Deploy Secret Posture — Decision Doc

**Created:** 2026-11-04
**Scope:** how the Firebase Hosting deploy in `.github/workflows/website.yml` authenticates, what the current posture risks, and the two ways off it.

## Current posture

- The Website workflow's `deploy` job uses the repository secret
  `FIREBASE_SERVICE_ACCOUNT_LARS_SOFTWARE` — a **long-lived SA JSON key** for
  `github-website-deploy@lars-software.iam.gserviceaccount.com` (no expiry;
  Google issues these with `validUntilTime = 9999-12-31`).
- Without that secret every CI deploy fails at the final step and the live site
  **silently freezes** — this happened 2026-07-26 → 2026-09-17, which is why the
  posture is documented instead of tribal knowledge.
- The shared GitHub-deploy SA pattern is used by other LarsArtmann repos too
  (templ-components runs the same SA shape). Do not mint keys for the old
  `firebase-adminsdk` SA: it is over GCP's 10-key cap and `keys.create` fails
  with `FAILED_PRECONDITION`.
- IAM on the `lars-software` project is only visible to `lartyhd@gmail.com`
  (not `lars@helpless.ai`).

## Option A (recommended): Workload Identity Federation

GitHub's OIDC token exchanged for short-lived GCP credentials via WIF. The
repo secret disappears entirely; nothing long-lived is stored on GitHub.

One-time setup (owner action, `lartyhd@gmail.com`):

1. Create a WIF provider for the repo:
   `gcloud iam workload-identity-pool providers create-oidc github --location=global --workload-identity-pool=github --issuer-uri=https://token.actions.githubusercontent.com --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository"`
2. Allow the repo principal:
   `gcloud iam service-accounts add-iam-policy-binding github-website-deploy@lars-software.iam.gserviceaccount.com --role=roles/iam.workloadIdentityUser --member="principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github/attribute.repository/LarsArtmann/go-atomic-write"`
3. Swap the deploy step to `google-github-actions/auth` with
   `workload_identity_provider: projects/…/providers/github` + `service_account: github-website-deploy@…` (SHA-pinned action), then remove the JSON secret.

## Option B (fallback): rotation cadence

If WIF is deferred: mint a fresh key quarterly, update the secret, delete the
old key. Rotation without deletion is security theater — the deleted key is
the actual win.

## Stale-key cleanup (owner action, GCP console)

The `firebase-adminsdk-dwv0a` SA sat at 12 keys (July 2026), above the 10-key
cap; ~7 of those are stale keys from July. Delete unused ones via
`gcloud iam service-accounts keys list --iam-account=firebase-adminsdk-dwv0a@lars-software.iam.gserviceaccount.com`
(GCP console works too; IAM on the project is only visible to `lartyhd@gmail.com`).
Deleting the stale keys also un-blocks future `keys.create` calls on that SA.

## Related

- The local-deploy recipe (no secret needed) is in `CONTRIBUTING.md` → Website notes.
- `firebase-adminsdk` key-cap gotcha: `AGENTS.md` → Gotchas.
