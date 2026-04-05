---
name: "job-application"
description: "Unified skill: discover jobs, match to resume/performance review, generate reports (discover/learn/execute), and produce tailored resume + cover letter for apply-ready jobs. Optional strict_no_skill_gaps (preferences.yaml) limits execute/ to roles with no Medium/High skill gaps vs evidenced resume/review. Single intake feeds matching and application materials. Use for end-to-end job search and application prep."
---

# Job Application

## Overview

Job Application is a unified skill combining job discovery, matching, and application materials. Two-stage flow: **Planning** (steps 1–7 in `references/job-application-planning/`) and **Execution** (Phase A: discover + report; Phase B: resume + cover letter for execute/ jobs).

Single intake (resume, performance review, application-ready highlights) is the source of truth for both job matching and resume/cover letter generation.

## Prerequisites

- Resume: `.claude/skills/job-application/resources/Resume-plain-text.md` (or .txt)
- Performance review: `.claude/skills/job-application/resources/performance_review.md` (use wisely)
- Preferences (default params): `references/preferences.yaml` — fill with your defaults; user overrides apply
- Memory init (for URL dedup): `npx @claude-flow/cli memory init`

## What This Skill Does

1. **Intake** — Resume, performance review, application-ready highlights (quantifiable wins, lead stories), application form readiness (why-company, interesting project PAR)
2. **Discover** — WebSearch/browser → candidate URLs
3. **Report** — Per-URL reports with match score, company, gaps, nudges → discover/learn/execute
4. **Application materials** — Tailored resume, cover letter, and (when form fields exist) ATS application answers for execute/ jobs (or user-specified job)

---

## Parameter Resolution (Orchestrator)

1. **Load** `references/preferences.yaml` — default params (target_roles, location, paths, mode, **matching**).
2. **Override** with user-provided params in the message (e.g. "Target: X", "Phase A only", "Job URL: Y", "Strict no skill gaps: on").
3. **Apply** resolved params to intake, search strategy, and execution.

See [references/preferences.yaml](references/preferences.yaml) for param list and resolution order.

### Matching: strict no skill gaps

When `matching.strict_no_skill_gaps` is **true** (or the user says e.g. "only zero-gap jobs", "no skill gaps"):

1. **Build a verified skills set** before discovery — from the resume (Experience tech lines, Skills section, projects). If `gap_evidence_scope` is `resume_plus_review`, you may also include technologies **explicitly named** in the performance review with evidence of real use (not generic praise).
2. **Discovery** — Prefer postings whose **required** stack is a subset of verified skills. Deprioritize or skip URLs whose titles/snippets advertise **must-have** tech not in that set. Search queries should mirror the verified stack (e.g. title + "Go" + "Kafka" + location), not aspirational stacks.
3. **Job report §6 (Gaps)** — For every **required** or **core day-one** skill in the JD: either map it to verified evidence or list it as a gap with severity **High** (missing required) or **Medium** (core implied but absent). **Low** = nice-to-have / clearly transferable only.
4. **§4 Ready to Apply?** — **Yes → execute/** only if §6 has **no** rows with severity **Medium** or **High** (only Low or empty). Otherwise **Not yet** or **Needs work** → **learn/** only (same as today for non-ready jobs).
5. **summary.md** — Include a section **Zero-gap matches (strict)** listing jobs that cleared the gate. If none: state that explicitly and suggest relaxing `strict_no_skill_gaps` or expanding verified skills honestly.

**Honesty:** Strict mode reduces false "apply now" jobs; it does not invent matching roles. Some markets may return few or zero execute-ready listings in one run.

---

## Orchestration Flow

**Step transitions, retries, mode branches:** [references/orchestrator.yaml](references/orchestrator.yaml). Per-step retry configs: `stepN/retry-config.yaml`.

**Planning** (step files):
1. [01-intake](references/job-application-planning/step1/01-intake.md) → 2. [02-search-strategy](references/job-application-planning/step2/02-search-strategy.md) → 3–6 (specs) → 7. [07-execution-checklist](references/job-application-planning/step7/07-execution-checklist.md)

**Execution:** Phase A (discover → report → summary) then Phase B (resume + cover letter + application-answers). Mode selects branch; see orchestrator.

### Browser Integration (Phase A)

```bash
agent-browser --session job-application open "https://www.linkedin.com/jobs/..."
agent-browser --session job-application snapshot -i
agent-browser --session job-application get text @e1
```
See `.claude/skills/browser/SKILL.md` for full commands.

---

## Run Modes

| Mode | Command | Output |
|------|---------|--------|
| **Discover only** | "Run job application Phase A only" | Reports + summary |
| **Full** | "Run full job application" | Reports + resume + cover letter + application-answers (skip execute/ jobs that already have outputs). With `strict_no_skill_gaps`, Phase B runs only for jobs with no Medium/High skill gaps. |
| **Outputs only** | "Run job application outputs only" | Resume + cover letter for all execute/ jobs missing outputs (replay Phase B) |
| **Application only** | "Generate resume and cover letter for job URL X" | Resume + cover letter for that job |

**Phase B re-run safe:** Jobs with both `-resume.md` and `-cover-letter.md` in outputs/ are skipped. Application-answers generated only when §10 exists and file missing. New execute/ jobs get outputs; existing ones preserved.

---

## Output Files

| Path | Content |
|------|---------|
| `job-application/discover/` | All job reports |
| `job-application/learn/` | Jobs to pursue (symlinks) |
| `job-application/execute/` | Ready to apply (symlinks) |
| `job-application/summary.md` | Top matches, gaps, nudges |
| `job-application/outputs/{company}-{title}-resume.md` | Tailored resume |
| `job-application/outputs/{company}-{title}-cover-letter.md` | Tailored cover letter |
| `job-application/outputs/{company}-{title}-application-answers.md` | ATS form field answers (when job report §10 exists) |

---

## Data Flow

```
Intake (step 1)
    ├──→ Search strategy, URL dedup
    ├──→ Job matching (resume + performance review)
    └──→ Resume/CL generation (application-ready highlights)

Job report (step 4)
    └──→ Resume/CL generation (Company §5, Nudges §8, Keywords §7, Gaps §6)
```

---

## Honesty Protocol

- Use [Estimate], [Assumption], [Data] tags in reports and summaries
- Add: "Projections are estimates, not guarantees."
- When WebSearch/browser unavailable: mark [Knowledge-Based — verify independently], reduce confidence
- Company work style/popularity: cite source or mark [Research needed]

---

## Related Skills

- **browser** — LinkedIn, Workday navigation

---

## Troubleshooting

### Memory not initialized
```bash
npx @claude-flow/cli memory init
```
Use `visited-urls.txt` fallback until Memory is ready.

### Browser session fails
Fall back to WebSearch for job discovery.

### Duplicate URLs
Normalize URL (lowercase, strip trailing slash, fragment) before hashing. See step 3 examples.
