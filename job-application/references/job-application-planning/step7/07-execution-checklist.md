# Step 7: Execution Checklist

*Skill: job-application | Planning step 7*

## Purpose

Final planning step. Defines Phase A (discover + report) and Phase B (application materials). Run modes: Discover only, Full, or Application only.

## Retries & Human Interaction

Use [retry-config.yaml](retry-config.yaml) for pre_exec, phase_a, phase_b_per_job: condition → prompt → on_exhausted (abort / infer / skip_job).

## Pre-Execution Checklist

- [ ] **Preferences**: Load `references/preferences.yaml` (defaults); override with user-provided params (including `matching.strict_no_skill_gaps`)
- [ ] **Step 1**: Intake complete (resume, performance review, application-ready highlights, application form readiness)
- [ ] **Step 2**: Search strategy defined
- [ ] **Step 3**: URL deduplication ready
- [ ] **Step 4**: Job report template reviewed
- [ ] **Step 5**: Application strategy reviewed
- [ ] **Step 6**: Resume and cover letter templates reviewed
- [ ] Resume at `.claude/skills/job-application/resources/Resume-plain-text.md`
- [ ] Performance review at `.claude/skills/job-application/resources/performance_review.md`

---

## Phase A — Discover & Report

1. **Load context**
   - Intake (step 1)
   - Search strategy (step 2)

2. **Discover URLs**
   - WebSearch (max 3–5 queries)
   - Optionally: browser `agent-browser --session job-application open "https://linkedin.com/jobs/..."`; snapshot -i; extract links
   - Collect candidate URLs

3. **Deduplicate**
   - Check Memory / visited-urls.txt
   - Filter to NEW URLs only

4. **Process each URL** (max 10–12)
   - Fetch (sequential); skip domains that returned 429/403
   - Match to intake (resume + performance review); enforce §6 severities and `matching` from preferences
   - Write report to `job-application/discover/{company}-{title}.md`
   - Add symlink to learn/ or execute/ per Ready to Apply? **If `matching.strict_no_skill_gaps` is true:** symlink to **execute/** only when §6 has no **Medium** or **High** gaps; otherwise **learn/** only
   - Store URL in Memory: `npx @claude-flow/cli memory store --key "job-application-urls:{hash}" --value "{url}|{date}" --namespace job-application-urls`
   - Append to `references/job-application-planning/visited-urls.txt`

5. **Summarize**
   - Create `job-application/summary.md` (use [Estimate], [Data] tags; add "Projections are estimates, not guarantees.")
   - If `strict_no_skill_gaps`: add subsection **Zero-gap matches (strict)** — jobs that symlinked to execute/ under the gate; if none, say so clearly
   - Update discover/learn/execute READMEs

---

## Phase B — Application Materials (for execute/ jobs)

**Re-run safe:** Before generating for each job, check if outputs already exist. If both `{base}-resume.md` and `{base}-cover-letter.md` exist → SKIP (preserve user edits). If either missing → generate both.

6. **For each job in execute/** (or user-specified job URL/path, max 3–5):

   **Dedupe check (re-run safe):**
   - Base name = report filename without `.md` (e.g. `monzo-backend-engineer-iii`)
   - If `job-application/outputs/{base}-resume.md` AND `{base}-cover-letter.md` both exist → SKIP this job
   - Else → generate (or regenerate if one is missing)

   **CRITICAL: Minimal data loading (Application only mode)**
   - If job URL given: `grep -l "EXACT_JOB_URL" .claude/skills/job-application/job-application/discover/*.md` → read ONLY the matching file
   - Do NOT read summary.md or other job reports
   - Load: job report + base resume + performance review (3 files only)

   **For each job (if not skipped):**
   - Compare base resume to job report: strengths, gaps, reframe opportunities, keywords placement
   - Generate tailored resume → `job-application/outputs/{company}-{title}-resume.md`
   - Generate cover letter using: generated resume + job report §5 (Company), §6 (Gaps), §7 (Keywords), §8 (Nudges) → `job-application/outputs/{company}-{title}-cover-letter.md`

   **Application form answers (if job report §10 exists):**
   - If `outputs/{base}-application-answers.md` exists → SKIP (re-run safe)
   - Else: Generate field-by-field answers from intake Application Form Readiness + job report §5, §7, §8 → `job-application/outputs/{company}-{title}-application-answers.md`
   - Use PAR format for "interesting project"; adapt why-company template to job report §5; keep under char limits

---

## Run Modes

| Mode | Input | Phase A | Phase B |
|------|-------|---------|---------|
| **Discover only** | Intake + preferences | ✓ | — |
| **Full** | Intake + preferences | ✓ | ✓ (for execute/; skip jobs with existing outputs) |
| **Outputs only** | Intake (from prior run) | — | ✓ (all execute/; skip jobs with existing outputs) |
| **Application only** | Job URL or report path + intake | — | ✓ (that job only; no skip — user explicitly requested) |

---

## Validation

- [ ] Reports in discover/; symlinks in learn/execute/
- [ ] Summary created
- [ ] Resume, cover letter, and (when §10 exists) application-answers in outputs/ (if Phase B run)

## Examples

See [07-Execution-checklist-examples.md](07-Execution-checklist-examples.md).
