# SSO and authentication

**Document:** `integrations/sso-and-authentication.md`
**Version:** 5.0 · **Last updated:** 2026-08-14 · **Owner:** Engineering
**Applies to:** Grantly 2026.8 and later

Grantly has two separate authentication surfaces: **staff** (caseworkers, reviewers, administrators
inside the tenant's organisation) and **applicants** (the public portal). They are configured
independently and never share a session.

## Staff authentication

| Method | Protocol | Notes |
| --- | --- | --- |
| SAML 2.0 | SP-initiated | Most common with municipalities |
| OIDC | Authorization code + PKCE | Recommended for new setups |
| Local accounts | Email and password + TOTP | Fallback, discouraged |

A tenant can have one SAML or one OIDC provider, plus local accounts for break-glass access. Mixing
SAML and OIDC in the same tenant is not supported.

### SAML configuration

| Setting | Value |
| --- | --- |
| ACS URL | `https://{tenant}.grantly.se/auth/saml/acs` |
| Entity ID | `https://{tenant}.grantly.se/auth/saml/metadata` |
| NameID format | `emailAddress` (required) |
| Signing | Assertions must be signed; encryption optional |
| Clock skew tolerance | 3 minutes |

Required attributes: `email`, `given_name`, `family_name`. Optional: `groups`.

### OIDC configuration

Discovery is read from the provider's `.well-known/openid-configuration`. Claims used: `email`,
`given_name`, `family_name`, `groups`. Since release 2026.8, `acr_values` can be set per tenant to
require a specific authentication level — used by tenants that need step-up for decision makers.

## Roles and group mapping

| Role | Can |
| --- | --- |
| Reviewer | See and score assigned applications |
| Caseworker | Manage applications, return, put on hold |
| Decision maker | Issue approvals and rejections, backdate `decision_date` |
| Programme administrator | Configure programmes, bulk actions, reviewer assignment |
| Finance | Payment plans, runs, reconciliation |
| Tenant administrator | Users, integrations, retention settings, legal holds |

Roles are assigned in Grantly, or derived from the `groups` claim via a mapping table configured per
tenant. Where a mapping exists, it wins: a role set by hand is overwritten at the next sign-in.

> **Note**
> A user with no mapped group and no local role gets no access at all, not read-only access. This is
> the usual cause of "SSO works but the user sees an empty screen".

## Just-in-time provisioning

Staff accounts are created on first successful sign-in. Deprovisioning is not automatic: removing a
user from the identity provider stops them signing in, but their Grantly account stays active until a
tenant administrator deactivates it. Deactivated accounts are deleted after 24 months, see
[Data retention](../legal/data-retention.md).

SCIM is not supported.

## Staff sessions

| Setting | Value | Configurable |
| --- | --- | --- |
| Idle timeout | 60 minutes | 15–240 minutes per tenant |
| Absolute lifetime | 12 hours | No |
| Single logout (SAML) | Supported | On/off |
| Concurrent sessions | 10 | No |

## Applicant authentication

Applicants do not use the tenant's identity provider. Available methods:

| Method | Identifier | Notes |
| --- | --- | --- |
| Email and password | Email | TOTP optional, chosen by the applicant |
| BankID | Personal identity number | Via Minso's BankID connection |
| Organisation account | Organisation number | Members invited by email |

Applicant session behaviour — the 30-minute idle timeout and the warning dialogue — is documented in
[Portal for applicants](../product/portal-for-applicants.md#sessions-and-timeouts).

Password rules: minimum 12 characters, checked against a breached-password list, no forced rotation.
Reset links are valid for 60 minutes and single use. Five failed sign-ins lock the account for 15
minutes.

## API clients

API clients authenticate separately with OAuth 2.0 client credentials and are unaffected by staff
SSO. See [API overview](api-overview.md). An API client keeps working when the administrator who
created it leaves; rotate secrets on staff turnover.

## Troubleshooting

| Symptom | Usual cause |
| --- | --- |
| Loop back to the sign-in page | Clock skew over 3 minutes |
| "User not found" | NameID is not `emailAddress` |
| Empty screen after sign-in | No mapped group |
| Works for some users only | Group mapping missing for one group |
| Logged out of the wrong tenant | Two tenants, one browser — fixed in 2026.6.1 |

## Related documents

- [API overview](api-overview.md)
- [Portal for applicants](../product/portal-for-applicants.md)
- [Incident escalation](../runbooks/incident-escalation.md)
