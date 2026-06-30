# M1 — per-org logical DB experiment, findings

The goal: figure out whether Quack itself can enforce per-org isolation, so a token issued for `org_test_alpha` can never read `org_test_beta` data, regardless of what SQL the client sends. The fallback (M2) is enforcement at the Olive translation layer, which the design doc already commits to.

This file records what we observed in the spike.

## What Quack actually exposes for auth

Two global settings:

```
quack_authentication_function   default: quack_check_token
quack_authorization_function    default: quack_nop_authorization
```

Both point at function names. We override them with custom DuckDB macros and observe behavior.

### Authentication callback

Signature: 3 VARCHAR args → BOOLEAN. Probed semantics:

```
arg slot │ value observed                                       │ interpretation
─────────┼──────────────────────────────────────────────────────┼─────────────────────
   a     │ 32-char uppercase hex                                │ server-side session
         │   (e.g. 241724834749CDD66688B73EC415F464)             │ or token derivation
   b     │ client's token (the one carried by CREATE SECRET)    │ authoritative input
   c     │ non-empty, length > 10                               │ not yet identified
```

Behavior verified:

- Setting `quack_authentication_function = 'm1_authn'` does take effect — the custom macro is consulted on every connection.
- A token NOT present in our `token_orgs` table → "Authentication failed".
- A token present in the table → query proceeds.

→ Per-token gating is solid; M1 lite (allow / deny at the gate) works.

### Authorization callback

Signature: 2 VARCHAR args → BOOLEAN. Probed semantics:

```
arg slot │ value observed                                       │ interpretation
─────────┼──────────────────────────────────────────────────────┼─────────────────────
   x     │ 32-char uppercase hex                                │ same shape as authn
         │   (same per-session opaque ID, probably)             │ slot a
   y     │ contains the SQL string ('SELECT', etc.)             │ the actual statement
```

→ The full SQL is handed to the authz callback. This means we *can* inspect what's being asked.

## Where the model breaks for per-org row-level isolation

To enforce "token alpha can only see org_test_alpha", the authz callback would need to:

1. Know the calling token's allowed org.
2. Confirm the SQL it received only touches that org.

Step (2) is feasible: parse `y` for `WHERE org_id = 'X'` and check `X` against the allowed org.

Step (1) is the wall:

- The authz callback does **not** receive the client's token. It only gets an opaque session hex.
- DuckDB macros cannot have side effects (no `INSERT`, no session variables). So authn cannot record `(session_hex → org)` for authz to look up.
- There is no built-in way to carry state from authn into authz with macros alone.

Workarounds we considered:

- **A C++ / Python UDF that records state.** Possible, but moves us out of the "small spike" zone and adds a custom-extension dependency that has to ship with the image.
- **Encode the org in the SQL contract: every query must say `org_id = 'X'`, and authz simply rejects anything missing that filter.** This catches sloppy queries but does not stop a malicious caller from inserting any org they like — without knowing the caller's true org, authz has nothing to compare against.
- **One Quack server per org** (or per shard). Authentication scopes naturally because each pod only knows about one org's data. Heavy ops cost (100 orgs → 100 pods, or 10 sharded pods); the operational story is real, the security story is clean.

## What this tells us about the design

- **M1 as "Quack enforces row-level isolation"**: not practical with the current callback model. Achievable only by writing a custom server-side state extension, or by running per-org servers.
- **M1 lite ("Quack gates the connection by per-org token, content filtering lives elsewhere")**: works today. Per-org tokens, custom authn macro, ~10 lines of SQL setup. This is what we'd ship if we wanted layered defense.
- **M2 ("Olive is the only enforcer")**: still the recommendation in the design doc. It does not depend on Quack capabilities.

The combination "M1 lite + M2" gives defense in depth: even if Olive's WHERE rewrite has a bug, the per-org token at the Quack gate limits blast radius to that org's data only if isolation is also at the storage layer. With shared views across orgs (our current setup), per-org tokens at the gate only protect against credentials being stolen by a different tenant, not against query mistakes from the same tenant.

## Approaches considered + verdicts

Two server-side enforcement models were on the table:

- **Row-level isolation** — one DB, views over all orgs, server filters rows per caller. Verdict: **not workable with current Quack callbacks** because authz does not receive the caller's token; it cannot validate the SQL against an unknown org.

- **DB-level isolation (one DB per org)** — separate ATTACHed databases per org, per-org tokens routed to per-org databases. Verdict: **not operationally sustainable**. Adding a new org means restarting Quack with a new ATTACH list. Schema evolution touches N databases. Boot cost scales linearly with org count. Backups multiply. The model is correct in principle but a maintenance bomb.

Both paths fail. There is no cheap way for Quack to enforce per-org isolation on its own.

## Recommendation

For the production design:

1. **Olive is the canonical isolation layer (M2 — the original design)**. All SQL going to Quack already has `WHERE org_id = $caller.org_id` injected at the typed endpoint. SQL passthrough is admin-only. Quack just executes whatever it is told.
2. **Per-org Quack tokens (M1 lite) as defense-in-depth**, not correctness. Each Olive tenant uses a different Quack token stored alongside that tenant's other credentials. If Olive is compromised for one tenant, the stolen token only reaches that tenant's data once a server-side mechanism actually scopes to that tenant — without such scoping, the token gates the front door but does not narrow what's visible behind it. So this layer is useful only when paired with one of the unworkable approaches above; documenting it here for completeness, not as a recommendation.
3. **Skip full per-org isolation inside Quack** until either the callback model exposes caller identity to authz, or we are willing to write a custom server-side extension that maintains `(session → org)` state.

## Future-only paths (do not attempt now)

- **Sharded Quack** — N pods, each scoped to a subset of orgs. Solves correctness, fails on ops at our growth profile.
- **Custom Quack auth extension** — C++/Rust extension storing `(session → org)`. Multi-month effort. Only justified if we have a hard regulatory requirement that mandates server-side enforcement.

## Scripts in this folder

- `scripts/m1-server.sh` — boots a Quack server with custom authn + authz macros.
- `scripts/m1-probe-arg.sh <a|b|c>` — used to identify which authn arg slot holds the client token.

Both are kept in the repo as evidence; they are not part of the main `make all` flow.
