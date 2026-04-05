# Step 1: Intake

*Skill: job-application | Planning step 1*

## Purpose

Gather and summarize resume, performance review, and preferences into a single intake document. This is the **single source of truth** for job matching AND application materials (resume, cover letter, ATS form answers). Use performance review wisely: pick achievements most relevant to target roles and job description.

## Retries & Human Interaction

When inputs are missing or invalid, use [retry-config.yaml](retry-config.yaml): evaluate condition → prompt user (per case retries) → on_exhausted: abort / infer / ask_paste.

## Inputs

| Source | Path | Required |
|--------|------|----------|
| Resume | `.claude/skills/job-application/resources/Resume-plain-text.md` or `.txt` or user paste | Yes |
| Performance review | `.claude/skills/job-application/resources/performance_review.md` or user paste | Yes — use wisely |
| Preferences | `references/preferences.yaml` | No (ask or infer) |

## Output Sections

Write the following into this file (or a linked `01-intake-output.md`):

### Resume Summary
- Years of experience
- Key technologies and frameworks
- Notable projects and outcomes
- Education and certifications
- Companies and roles (chronological)

### Performance Review Highlights
- Strengths (as noted by reviewer)
- Growth areas or goals
- Quantifiable achievements mentioned
- Areas to emphasize in applications

### Application-Ready Highlights (for resume, cover letter, and ATS form fields)

**Quantification formula:** [Action Verb] + [What You Did] + [Quantified Result] — resumes with measurable outcomes get ~40% more callbacks.

- **Quantifiable wins** — At least 3–5 metrics across: time, money, people, volume, quality, scale, efficiency (e.g. PRs, events/month, latency P99, cost %, team size)
- **Lead stories** — 2–3 anecdotes in PAR format (Problem–Action–Result) for cover letters and "most interesting project" forms
- **Keywords to emphasize** — From performance review that match typical target roles

### Preferences
- Target roles (e.g., "Senior Backend Engineer", "Staff Platform Engineer")
- Location / remote (e.g., "Remote US", "Hybrid NYC")
- Companies to target or avoid
- Level (Senior, Staff, Principal)
- Tech stack preferences (e.g., "Go, Kubernetes preferred")

### Match-Ready Profile (Synthesis)
A short paragraph combining the above for quick reference during matching and application generation:
> "X years experience in Y. Strong in A, B, C. Performance review highlights D, E. Targeting F roles, G location. Prefers H tech. Application-ready: [key story], [quantification]."

### Application Form Readiness (for ATS/Greenhouse/Lever custom questions)

**Why-[Company] Bank (adaptable template):**
- Mission fit: [e.g. "make money work for everyone"]
- Product fit: [e.g. "I use Monzo; early salary, pots, transparent UX"]
- Tech fit: [e.g. "Go, Kafka, K8s — aligns with my distributed systems work"]

**Most Interesting Project (PAR format — problem, action, result):**
- Project: [e.g. "Courier race condition RCA"]
- Objective: [e.g. "Root-cause intermittent production panic"]
- Why interesting: [e.g. "60 goroutines, non-deterministic, cold-start edge case"]
- What I learned: [e.g. "Concurrent map writes; singleton patterns; -race under load"]
- Links: [GitHub, blog, internal doc — if shareable]

**Location answer:** [e.g. "Remote UK" / "London office" / "Cardiff"]

## Checklist

- [ ] Resume read and summarized
- [ ] Performance review read and used wisely
- [ ] Application-ready highlights extracted (3–5 quantifications, 2–3 PAR stories)
- [ ] Application form readiness filled (why-company template, interesting project PAR, location)
- [ ] Preferences captured or inferred
- [ ] Match-ready profile written
- [ ] Intake saved to this file or `01-intake-output.md`

## Examples

See [01-Intake-examples.md](01-Intake-examples.md) for filled examples.
