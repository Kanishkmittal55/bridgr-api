# Step 6: Resume and Cover Letter Templates

*Skill: job-application | Planning step 6*

## Purpose

Define output structure for tailored resume, cover letter, and ATS application form answers.

## Retries & Human Interaction

[retry-config.yaml](retry-config.yaml): no human interaction cases (0 retries).

---

## Resume

### Output Location

```
job-application/outputs/{sanitized-company}-{sanitized-title}-resume.md
```

### Sections (in order)

1. Headline — Role focus + key tech
2. Summary — 2–3 sentences targeting this role
3. Experience — Chronological, 3–5 bullets per role, quantified
4. Education
5. Skills — Job-critical first
6. Projects (optional)

### ATS

- Standard headers; keywords from job report; 1 page target

---

## Cover Letter

### Output Location

```
job-application/outputs/{sanitized-company}-{sanitized-title}-cover-letter.md
```

### Structure

1. **Opening** — Role + company; why them (job report §5)
2. **Lead story** — One strong example (Nudges §8, application-ready highlights)
3. **Fit** — How experience maps to role; address gaps briefly
4. **Close** — Eager to discuss; call to action

### Sources

- **Intake** — Application-ready highlights, match-ready profile
- **Job report §5** — Company culture, work style
- **Job report §6** — Gaps to preempt
- **Job report §7** — Keywords to highlight
- **Job report §8** — Nudges (e.g. lead with courier RCA)

### Length

- 3–4 short paragraphs; no hedging language

---

## Application Form Answers

### Output Location

```
job-application/outputs/{sanitized-company}-{sanitized-title}-application-answers.md
```

### Structure

Field-by-field answers keyed by question text (from job report §10). Format:

```markdown
## [Question text]

[Answer, respecting char limit if shown]
```

### Common Fields

- **What attracted you to [Company]?** — Mission, product, tech fit (intake + job report §5)
- **Most interesting project?** — PAR: objective, why interesting, what learned, links
- **Where would you like to be based?** — From intake Preferences

### Length

- Many ATS fields cap ~1000 chars; be concise

---

## Checklist

- [ ] Resume template clear
- [ ] Cover letter template clear
- [ ] Application answers template clear (when §10 exists)

## Examples

See [06-Resume-and-cover-letter-templates-examples.md](06-Resume-and-cover-letter-templates-examples.md).
