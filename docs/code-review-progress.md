# Code review / dev-profile audit — progress log

Autonomous overnight run, continuing the same pass that just finished on `~/verv/Velez`
(see that repo's `docs/code-review-progress.md` for the full methodology writeup). This is a
Cobra CLI tool, not a service — `31-go-layering.md` (transport/service/repository) mostly doesn't
apply; the relevant packs are `00-core.md`, `10-workflow.md`, `30-go.md`, `32-go-testing.md`,
`20-architecture.md`.

Legend: **L** linter-enforced · **C** custom-linter candidate · **F** language/build-forced ·
**P** deliberate preference, recorded not enforced.

Status key per section: ⏳ queued · 🔄 running · ✅ done (verified green) · ⚠️ done with open
questions for the user.

Note (from this repo's own recent history, relevant to this audit): a fix already landed here this
session as part of the Velez matreshka-migration work — `gen_servers.go:21`'s error-string casing
(commit `c11a0f9`) and a version bump to `v0.0.49` (pushed to `master`). Not re-litigated here.

---

## 1. `cmd/` + `internal/` — ⏳ queued

## 2. `plugins/project/` — ⏳ queued

## 3. `plugins/tools/` + `tests/` (code-quality read-through) — ⏳ queued

---

## Proposed rule changes awaiting your one-word approval

(Populated as audits surface genuine profile gaps.)
