# Kanishk Mittal

**Go Backend Engineer — microservices · CI/CD · PostgreSQL · distributed systems**

| | |
|---|---|
| Email | kanishkable@gmail.com |
| Phone | 07899007443 |
| LinkedIn | linkedin.com/in/kanishkable |
| Website | millionairethinking.co.uk |

---

## Summary

Backend engineer with production Go experience at DailyPay (Belfast), building microservices across event-driven pipelines, Kafka, gRPC, and PostgreSQL. Conducted 190+ code reviews across 10+ teammates in nine months — technically substantive, catching bugs and proposing architectural patterns. Initiated and owned golangci-lint CI integration across 52 files; delivered cross-service pipelines end-to-end across four repositories. Remote-first by default (Belfast); fluent in async collaboration and documentation-driven development.

---

## Experience

### Software Developer I — DailyPay, Belfast
**Jul 2025 – Present**

*Go · Kafka · gRPC · PostgreSQL · Protobuf · Docker · GitHub Actions · CI/CD*

- **Golangci-lint CI initiative:** Integrated 4–5 linters across 52 files (236 additions), improving code quality standards across the communications codebase — self-directed, no ticket required.
- **Cross-service pipeline (end-to-end):** Led a 4-repository initiative spanning `courier`, `dp-protos`, `data-router`, and Amplitude — 1,892 additions including Protobuf schema design, Kafka topic management, gRPC integration, and event transformation serving 30M+ events/month.
- **Code review (190+ MRs):** Reviewed pull requests across 10+ teammates — caught proto field mapping bugs, flagged Terraform permission scope, proposed record processor patterns and `assert.Eventually` for test simplification.
- **Courier race condition RCA & fix:** Root-caused a production panic to a concurrent map write in the Schema Registry client triggered during cold-starts across 60 goroutines; implemented singleton fix verified with `-race`.
- **Testing:** Built comprehensive event testing suites validating schemas and idempotency; wrote integration tests for Kafka consumers, gRPC handlers, and SQS message processors.
- **On-call engineer:** Monitored production systems, diagnosed incidents, proposed remediations; participated in load testing to identify API breaking points.

### Backend Engineer (Placement) — DailyPay, Belfast
**Sep 2023 – Aug 2024**

*Go · Python · PostgreSQL · Docker · AWS · Kafka · GitHub Actions*

- **V3 Earnings Job Status endpoint:** Delivered full OAPI spec, handler, and integration test suite (Go, PostgreSQL).
- **Schema Registry migration:** Cross-repo migration for `provider_account_created` and `earnings_invalidation` events spanning `core` and `ledger`, with full integration testing.
- **ETL pipelines (v1–v3):** Designed and shipped Gross Earnings and Available Balance pipelines (Go, PostgreSQL, Docker, AWS); contributed to Ruby → Go microservices migration.
- **Multi-currency API:** Implemented currency support across three earnings endpoints, enabling international expansion to Canada.
- **CI/CD pipelines:** Created GitHub Actions workflows; built reusable `daily-go` library components improving system efficiency by 18%.

**Hackathon (DailyPay):** Winner (£500) — Led 3-person team; prototype recognised by CEO.

---

## Education

**B.Eng. Computer Systems Engineering** — Brunel University, London
*Sep 2021 – Jun 2025*

First Class Honours | Graham Hawke Award — Best Final Year Project in EEE
Dissertation: Knowledge Graph-Augmented LLMs for Patent Analysis

---

## Skills

| Category | Technologies |
|----------|--------------|
| **Languages** | Go, Python, SQL, Bash, JavaScript |
| **Infra & Tooling** | Kafka, gRPC, Protobuf, Docker, Kubernetes, Terraform, AWS, GitHub Actions, Prometheus |
| **Databases** | PostgreSQL, Neo4j |
| **Practices** | Code review (190+ MRs), CI/CD, integration testing, async-remote, microservices, event-driven architecture |

---

## Projects

### Bare Metal Kubernetes Cluster
- 4-node K8s cluster (3× Raspberry Pi + custom CPU) with Helm, k9s, nginx, 1.2TB NAS — public repo
- Automated deployments with reusable manifests; networking via ngrok and Glinet firewall
- [Video Demo →](https://millionairethinking.co.uk)

### Patent Knowledge Graph
- Scraping + NLP pipeline; semantic search via vector embeddings + graph traversal
- Python, Scrapy, Neo4j, LLMs
- [Video Demo →](https://millionairethinking.co.uk)
