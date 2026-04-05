# Step 3: URL Deduplication

*Skill: job-application | Planning step 3*

## Purpose

Ensure each job URL is processed only once.

## Retries & Human Interaction

[retry-config.yaml](retry-config.yaml): no human interaction cases (0 retries). Use Memory DB (namespace `job-application-urls`) as primary. Maintain `visited-urls.txt` as fallback.

## Memory DB

- **Namespace:** `job-application-urls`
- **Key format:** `job-application-urls:{hash}` (first 16 chars of SHA-256 of normalized URL)
- **Store after processing:**
  ```bash
  npx @claude-flow/cli memory store --key "job-application-urls:{hash}" --value "{url}|{date}" --namespace job-application-urls
  ```
- **Fallback path:** `references/job-application-planning/visited-urls.txt`

## URL Normalization

Before hashing: lowercase, remove fragment (#...), trailing slash, tracking params (utm_*, ref); sort query params for consistent hashing.

## Integration

For each candidate URL: normalize → check Memory (search for hash) or grep visited-urls.txt → if found SKIP → else process, store, append.

## Examples

See [03-Url-deduplication-examples.md](03-Url-deduplication-examples.md).
