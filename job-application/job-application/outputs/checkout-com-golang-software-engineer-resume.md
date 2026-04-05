# Kanishk Mittal
**Golang Software Engineer** | Go • Kafka • PostgreSQL • Kubernetes • Payments

| | |
|-|-|
| Email | kanishkable@gmail.com |
| Phone | 07899007443 |
| LinkedIn | [linkedin.com/in/kanishkable](https://linkedin.com/in/kanishkable) |
| Website | millionairethinking.co.uk |

---

## Summary

Backend engineer with production Go experience building high-throughput, event-driven systems in the payments and financial services space. Delivered a cross-service analytics pipeline processing 30M+ events/month, diagnosed and fixed a race condition in a concurrent 60-goroutine system, and shipped multi-currency payment features used by thousands of workers daily. Seeking to bring payments infrastructure expertise to Checkout.com's Payment Performance team.

---

## Experience

### Software Developer I — DailyPay, Belfast
**Jul 2025 – Present**
*Go, Kafka, PostgreSQL, Terraform, gRPC, Docker, Datadog*

- Delivered end-to-end email-to-Amplitude analytics pipeline across 4 repositories (courier → dp-protos → data-router → Amplitude): ~1,900 lines added, gRPC integration, Protobuf schema design, Kafka topic management, and event transformation — full unit + integration test coverage
- Diagnosed and fixed production courier panic caused by race condition in Schema Registry client under 60-goroutine cold-start; implemented singleton pattern; verified with `go test -race`
- Integrated Iterable marketing platform via gRPC; built event testing suite validating schemas and idempotency across 30M+ events/month
- On-call engineer: monitored production systems, triaged panics, proposed and implemented remediations from incident start
- Integrated golangci-lint CI across 52 files (4–5 linters); decommissioned SMSO service (3,730 lines removed)
- Resolved critical production bugs under pressure; participated in load testing real-world APIs to identify breaking points

### Data Engineer (Placement) — DailyPay, Belfast
**Sep 2023 – Aug 2024**
*Go, Python, PostgreSQL, Docker, AWS, Kafka*

- Designed and shipped ETL pipelines (v1–v3) for Gross Earnings and Available Balance, powering payroll calculations for tens of thousands of workers
- Implemented multi-currency API support for GrossEarnings, TipsDetailed, and TipsEarningsProviderAccountDetails endpoints — enabled company's international expansion to Canada
- Built reusable internal libraries (daily-go) and automation tooling, improving pipeline efficiency by 18%
- Built V3 Earnings Job Status endpoint with full OpenAPI spec and integration tests; executed Schema Registry migrations across core and ledger repositories
- Created CI/CD pipelines with GitHub Actions; wrote Linux scripts for PostgreSQL cron job scheduling; Terraform migration to shared monolith module

**Hackathon (DailyPay):** Winner (£500) — Led 3-person team to build AR-based customer onboarding prototype; recognized by CEO

---

## Education

**B.Eng. Computer Systems Engineering** — Brunel University, London
*Sep 2021 – Jun 2025*
- First Class Honours | Graham Hawke Award — Best Final Year Project in EEE
- Dissertation: Knowledge Graph-Augmented LLMs for Patent Analysis

---

## Skills

| Category | Technologies |
|----------|--------------|
| **Languages** | Go (production), Python, JavaScript, SQL, Bash |
| **Infrastructure** | Docker, Kubernetes, Terraform, AWS, GitHub Actions, Kafka, Prometheus, Datadog |
| **Databases** | PostgreSQL, Neo4j, Snowflake |
| **Practices** | Event-driven architecture, microservices, gRPC/Protobuf, CI/CD, on-call, OAPI spec |
| **Other** | LLM fine-tuning, Knowledge Graphs, LeetCode 150+ |

---

## Projects

**Bare Metal Kubernetes Cluster**
Built 4-node K8s cluster (3× Raspberry Pi + custom CPU) with Helm, k9s, nginx load balancing, 1.2TB NAS. Automated deployments with reusable manifests; tunneling via ngrok.

**Patent Knowledge Graph** (Dissertation)
NLP + semantic search pipeline to map chemical patent relationships using Python, Scrapy, Neo4j, vector embeddings, and LLMs. Awarded Best Final Year Project (EEE faculty).
