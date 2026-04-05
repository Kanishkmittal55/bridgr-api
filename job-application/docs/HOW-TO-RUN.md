# How to Run Job Application Skill

Unified job search + application materials.

---

## Prerequisites

1. Resume: `.claude/skills/job-application/resources/Resume-plain-text.md`
2. Performance review: `.claude/skills/job-application/resources/performance_review.md` (use wisely)
3. Memory (optional): `npx @claude-flow/cli memory init`

---

## Run Modes

### Full Run (Discover + Application Materials)

```
Run full job application. Resume at .claude/skills/job-application/resources/Resume-plain-text.md.
Target: Backend Engineer Go remote UK.
```

### Discover Only (Reports + Summary)

```
Run job application Phase A only. Find jobs, generate reports. No resume/cover letter.
Target: Backend Engineer Go remote UK.
```

### Outputs Only (Replay Phase B — Resume + Cover Letter for execute/)

**Re-run safe:** Skips jobs that already have resume and cover letter in outputs/. Generates resume, cover letter, and (when §10 exists) application-answers for execute/ jobs missing outputs.

```
Run job application outputs only. Generate resume and cover letter for all execute/ jobs that don't have outputs yet.
```

### Application Only (One Job)

**Minimal data loading:** Load ONLY the job report for the pasted URL. Do NOT read summary.md.

```
Generate resume and cover letter for Monzo Backend Engineer III.
Job URL: https://job-boards.greenhouse.io/monzo/jobs/6635595
```

Or with report path:

```
Generate application materials for job at job-application/execute/monzo-backend-engineer-iii.md
```

*Note: Application only (single job) does NOT skip — user explicitly requested that job. Use "Outputs only" for safe batch replay.*

---

## Outputs

| Path | Content |
|------|---------|
| `job-application/discover/*.md` | All job reports |
| `job-application/learn/` | Jobs to pursue |
| `job-application/execute/` | Ready to apply |
| `job-application/summary.md` | Top matches |
| `job-application/outputs/*-resume.md` | Tailored resumes |
| `job-application/outputs/*-cover-letter.md` | Tailored cover letters |
| `job-application/outputs/*-application-answers.md` | ATS form answers (when job has custom questions) |

---

## Migration

- **job-application** is the unified skill for job discovery and application materials.
- Resources (resume, performance review) live in `.claude/skills/job-application/resources/`.
- Intake now includes "Application-ready highlights" and "Application form readiness" (why-company, interesting project PAR) — extract from performance review for better resume/cover letter and ATS form answers.
