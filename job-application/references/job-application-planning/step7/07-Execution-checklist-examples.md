# Step 7: Execution Checklist — Examples

*Skill: job-application | Planning step 7*

## Example: Full Run

```
1. Intake + search strategy loaded
2. WebSearch × 3 → 12 URLs
3. Dedupe → 8 new
4. Process 8 → 8 reports (5 learn/, 3 execute/)
5. Summary created
6. Phase B: Generate resume + cover letter + (if §10) application-answers for 3 execute/ jobs
7. Outputs: monzo-*, gitlab-*, cuvva-* (resume + cover letter each; application-answers when form fields present)
```

## Example: Outputs Only (Replay Phase B)

```
1. List execute/*.md → monzo-backend-engineer-iii.md, gitlab-intermediate-backend-engineer-golang.md
2. For each: check if outputs/{base}-resume.md AND {base}-cover-letter.md exist
   - Monzo: both exist → SKIP
   - GitLab: both exist → SKIP
   - New job X: missing → generate both
3. Only generate for jobs missing outputs
```

## Example: Application Only

```
Input: Job URL https://job-boards.greenhouse.io/monzo/jobs/6635595
1. Grep discover/*.md for "6635595" → monzo-backend-engineer-iii.md
2. Load ONLY: that job report + base resume + performance review (do NOT read summary)
3. Compare resume to job report (strengths, gaps, reframe)
4. Generate resume → monzo-backend-engineer-iii-resume.md
5. Generate cover letter (using resume + report §5,§6,§7,§8) → monzo-backend-engineer-iii-cover-letter.md
6. If job report §10 (Application Form Fields) exists: generate application-answers → monzo-backend-engineer-iii-application-answers.md
   (adapt intake Application Form Readiness to "What attracted you to Monzo?", "Most interesting project?")
(No skip — user explicitly requested this job)
```
