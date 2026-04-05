# Step 2: Search Strategy

*Skill: job-application | Planning step 2*

## Purpose

Define how and where to search for jobs. Choose between browser automation (LinkedIn, Workday, job boards) and WebSearch. Each query produces candidate URLs for the execution stage.

## Retries & Human Interaction

Use [retry-config.yaml](retry-config.yaml) when intake missing or target_roles empty: prompt → retry → on_exhausted.

## Inputs

- Intake from [step1/01-intake.md](../step1/01-intake.md) (target roles, preferences)
- `references/preferences.yaml` → `matching` (if `strict_no_skill_gaps: true`, see below)
- Available tools: browser (MCP or agent-browser), WebSearch

### Strict no skill gaps (discovery bias)

When `matching.strict_no_skill_gaps` is **true**:

- Derive **verified skills** from resume (and performance review per `gap_evidence_scope`) before writing queries.
- Draft search queries that **stack-match**: role + **verified** technologies + location (e.g. `Backend Engineer Go Kafka PostgreSQL Remote UK`), not technologies you aspire to learn.
- Prefer sources/snippets that list requirements you already evidence; avoid chasing listings whose visible requirements are mostly outside verified skills (note in strategy: "skipped URL pattern: X").
- Expect fewer high-match URLs per run — that is intentional.

## Output Sections

### Search Queries (3–5 per role)

For each target role, generate 3–5 search queries. Include:
- Role name variations
- Tech stack keywords from resume
- Location filters
- Source-specific phrasing (LinkedIn vs Indeed vs generic)

### Sources & Method

| Source | Method | Notes |
|--------|--------|-------|
| LinkedIn Jobs | browser or WebSearch | WebSearch for quick discovery |
| Workday | browser | Often requires browser for ATS pages |
| Indeed / Glassdoor | WebSearch or browser | WebSearch often sufficient |

### Common ATS Application Form Questions (for job report §10)

Many Greenhouse, Lever, Workday jobs include custom questions. When fetching a job page, check for:
- "What attracted you to [Company]?" / "Why [Company]?" — adapt from intake why-company template
- "Most interesting project?" — use intake PAR (objective, why interesting, what learned, links)
- "Where would you like to be based?" — from intake Preferences
- "Why this role?" — focus on tasks/responsibilities vs company mission

### Context for Execution

- **Max URLs per run: 10–12** (cost control; see EFFICIENCY.md)
- **Max WebSearch queries: 3–5** per run
- **Sequential fetches** — avoid parallel WebFetch + WebSearch
- Prioritize: direct apply links > aggregators > company pages

## Checklist

- [ ] 3–5 queries per target role
- [ ] Method assigned per source
- [ ] Strategy saved to this file

## Examples

See [02-Search-strategy-examples.md](02-Search-strategy-examples.md) for filled examples.
