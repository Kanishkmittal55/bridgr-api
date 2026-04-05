# Step 1: Intake — Examples

*Skill: job-application | Planning step 1*

## Example: Filled Intake (real-world aligned)

```markdown
### Resume Summary
- 6 years backend/platform engineering
- Tech: Go, Python, Kubernetes, PostgreSQL, AWS
- Led migration of 50 services to k8s; 40% cost reduction
- BS Computer Science, Stanford

### Performance Review Highlights
- Strengths: System design, mentorship, delivery
- Goal: Grow into staff-level scope
- Highlight: "Consistently delivers under ambiguity"

### Application-Ready Highlights

**Quantification formula:** [Action Verb] + [What] + [Result]

- **Quantifiable wins:** 142 PRs, 190 reviews, 30M events/month, P99 latency 45ms→12ms, 18% efficiency gain (daily-go lib)
- **Lead stories (PAR):**
  1. Courier panic RCA: Problem — intermittent panic in Go comms service. Action — traced to concurrent map write in Schema Registry client (60 goroutines, cold-start). Result — singleton fix, -race verified, panic stopped.
  2. Schema Registry migration: 4-repo cross-service (provider_account_created, earnings_invalidation); full integration tests.
- **Keywords:** Go, Kafka, Kubernetes, gRPC, on-call, production ownership

### Preferences
- Roles: Senior Backend Engineer, Staff Platform Engineer
- Location: Remote US or hybrid SF
- Target: Series B–D startups, avoid FAANG

### Match-Ready Profile
6Y backend/platform. Go, K8s, PostgreSQL, AWS. Led large migration. Staff-track. Remote US. Application-ready: courier RCA, 30M events pipeline, 18% efficiency gain.

### Application Form Readiness (for ATS/Greenhouse/Lever custom questions)

**Why-[Company] Bank (adaptable template):**
- Mission fit: make money work for everyone; transparent banking
- Product fit: I use the app daily; early salary, pots, UX
- Tech fit: Go, Kafka, K8s — aligns with my distributed systems and on-call work

**Most Interesting Project (PAR):**
- Project: Courier race condition RCA
- Objective: Root-cause intermittent production panic in Go comms service
- Why interesting: 60 goroutines, non-deterministic, cold-start edge case; concurrent map write
- What I learned: Singleton patterns with sync.Once; -race detector under load; systematic RCA
- Links: [internal post-mortem if shareable]

**Location answer:** Remote UK / London office
```
