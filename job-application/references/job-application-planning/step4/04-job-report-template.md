# Step 4: Job Report Template

*Skill: job-application | Planning step 4*

## Purpose

Define the per-URL report structure.

## Retries & Human Interaction

[retry-config.yaml](retry-config.yaml): no human interaction cases (0 retries). Each processed job gets one report with match score, company insights, ready-to-apply status, and nudges. Reports feed resume and cover letter generation.

## Report Location

| Track | Path | Contents |
|-------|------|----------|
| **Discover** | `job-application/discover/{sanitized-company}-{sanitized-title}.md` | All reports |
| **Learn** | `job-application/learn/` | Jobs to pursue — symlinks |
| **Execute** | `job-application/execute/` | Ready to apply now — symlinks |

Based on **Ready to Apply?**: Yes → execute/; Not yet/Needs work → learn/.

When `preferences.yaml` → `matching.strict_no_skill_gaps` is **true**, **Yes** (and thus execute/) is allowed **only** if **§6 Gaps** has **no** Medium or High severity rows (see §6). Seniority or visa gaps may still justify **Not yet** — document in §4.

Sanitize filenames: lowercase, replace spaces with `-`, remove special chars. Example: `acme-inc-senior-backend-engineer.md`

## Template Structure

1. **URL** — Job posting link
2. **Summary** — Title, company, location, source
3. **Match Score** (0.0–1.0) — One-line interpretation ([Estimate], [Data] tags; add "Projections are estimates, not guarantees.")
4. **Ready to Apply?** — Yes | Not yet | Needs work (with reason in one sentence). Must align with §6 severities when strict matching is on (see preferences `matching.strict_no_skill_gaps`).
5. **Company** — Work style, popularity, culture (used for cover letter §5)
6. **Gaps (Resume vs Requirements)** — Table: Gap | **Severity** (High / Medium / Low) | Notes. Map every **required** or **core** JD skill to evidenced resume (and performance review if `gap_evidence_scope: resume_plus_review`) or log a gap. **Low** = optional / clearly transferable only.
7. **Keywords to Highlight** — In resume and cover letter
8. **Nudges (Boost Interview Chance)** — Tasks with expected impact
9. **Red Flags / Yellow Flags** — If any
10. **Application Form Fields** — Custom ATS questions (if visible on job page). List: question text, required/optional, char limit if shown. Common: "What attracted you to [Company]?", "Most interesting project?", "Where based?"

## Nudge Types

| Type | Example |
|------|---------|
| Skill gap | Complete X project / course — unlocks Y more roles |
| Resume tweak | Add Z keyword — +N% match |
| Cover letter | Highlight A from performance review; lead with courier RCA |
| Timing | Apply within 48h — early applicant advantage |
| Network | Find referral — 3x response rate |

## Checklist

- [ ] Template structure understood
- [ ] Nudge types clear
- [ ] §10 Application Form Fields captured when visible on job page
- [ ] Symlinks added to learn/ or execute/ per Ready to Apply? (execute/ only if strict mode allows — no Medium/High §6 gaps)
- [ ] Output dir `job-application/discover/` created when first report written

## Examples

See [04-Job-report-template-examples.md](04-Job-report-template-examples.md).
