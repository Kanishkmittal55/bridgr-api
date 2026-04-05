# Job Application — Efficiency & Cost Control

Use when runs cost too much or for budget-friendly execution.

---

## Cost Drivers

| Driver | Impact |
|--------|--------|
| WebSearch | ~$0.01–0.03 per search |
| WebFetch | Page content = tokens |
| Extended thinking | 4+ min = $$$ |
| Phase B (resume/CL) | Per-job generation |

---

## Hard Caps

| Cap | Value |
|-----|-------|
| WebSearch per run | Max 3–5 queries |
| URLs to process | Max 10–12 |
| Resume + cover letter per run | Max 3–5 jobs |

---

## Execution Rules

1. **Deduplicate first** — Check Memory / visited-urls before any fetch
2. **Sequential fetches** — No parallel WebFetch + WebSearch
3. **Skip 429/403 sites**
4. **Phase A only** for discovery — skip Phase B if you only need reports
5. **Application only** — Skip Phase A if you have reports; generate for 1 job

---

## Budget Prompt

```
Run job application Phase A only. BUDGET MODE:
- Max 3 WebSearch, 10 URLs
- Sequential fetches
- Target: [your target]
```
