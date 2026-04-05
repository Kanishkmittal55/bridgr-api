# Step 5: Application Strategy

*Skill: job-application | Planning step 5*

## Purpose

Define when and how to generate tailored resume, cover letter, and ATS form answers.

## Retries & Human Interaction

[retry-config.yaml](retry-config.yaml): no human interaction cases (0 retries). Uses intake (step 1) and job report structure (step 4). Planning for resume, cover letter, and application form fields is done together.

## Inputs

- Intake from [step1/01-intake.md](../step1/01-intake.md) — application-ready highlights
- Job report template from [step4/04-job-report-template.md](../step4/04-job-report-template.md)

## When to Generate Application Materials

- **Full run:** For each job in execute/ (Ready to Apply? Yes). Skip jobs that already have both resume + cover letter in outputs/.
- **Outputs-only run:** For all jobs in execute/; skip those with existing outputs. Replay Phase B safely.
- **Application-only run:** For user-specified job URL or report path (no skip — user explicitly requested).
- **Max per run:** 3–5 resumes + cover letters (cost control)
- **Application form answers:** When job report §10 exists, generate `{base}-application-answers.md`. Skip if file exists.

## Resume Analysis (per job)

Before generating, compare base resume to job report:
- **Strengths** — Bullets/keywords that already match
- **Gaps to address** — Missing keywords, experience not framed for role
- **Reframe opportunities** — Bullets to reorder/rewrite; projects to elevate
- **Quantification additions** — Numbers from performance review or job-report nudges
- **Keywords placement** — Where each keyword from job report §7 should appear (headline, bullets, skills)

## Resume Tailoring Principles

- **Headline** — Role focus + key tech from job keywords
- **Experience** — Lead with most relevant bullets; use application-ready quantifications
- **Skills** — Job-critical first
- **Gaps** — Address or reframe per job report §6

## Cover Letter Principles

- **Lead story** — Use Nudges §8 (e.g. "lead with courier RCA"); pick from application-ready highlights
- **Company fit** — Reference job report §5 (culture, work style)
- **Gaps** — Preempt briefly if needed (YoE, domain)
- **Length** — 3–4 short paragraphs; no hedging

## Application Form Answers (Greenhouse/Lever/Workday)

- **Source:** Intake Application Form Readiness + job report §5 (Company), §7 (Keywords), §8 (Nudges)
- **"Why [Company]?"** — Adapt intake why-company template to job report §5; mission, product, tech fit
- **"Most interesting project?"** — Use intake PAR; objective, why interesting, what learned, links; keep under char limit
- **"Where based?"** — From intake Preferences
- **Length** — Many ATS fields cap ~1000 chars; be concise

## Checklist

- [ ] When to generate defined (execute/ jobs or user-specified)
- [ ] Resume tailoring principles clear
- [ ] Cover letter principles clear
- [ ] Application form answer principles clear (when §10 exists)

## Examples

See [05-Application-strategy-examples.md](05-Application-strategy-examples.md).
