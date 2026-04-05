# Kanishk Mittal

**Go Backend Engineer — event-driven microservices · Kafka · Kubernetes · gRPC**

| | |
|---|---|
| Email | kanishkable@gmail.com |
| Phone | 07899007443 |
| LinkedIn | linkedin.com/in/kanishkable |
| Website | millionairethinking.co.uk |

---

## Summary

Backend engineer with production Go experience at a fintech (DailyPay, Belfast), building event-driven microservices across Kafka, gRPC, Protobuf, and distributed data pipelines handling 30M+ events/month. On-call ownership from day one — root-caused a Go race condition in a concurrent Schema Registry client under production load, implementing a singleton fix verified with `-race`. Delivered cross-service systems spanning four repositories end-to-end, including architecture, implementation, testing, and deployment.

---

## Experience

### Software Developer I — DailyPay, Belfast
**Jul 2025 – Present**

*Go · Kafka · gRPC · PostgreSQL · Protobuf · Docker · Terraform · GitHub Actions*

- **Courier race condition RCA & fix:** Root-caused a production panic to a concurrent map write in the Schema Registry client's HTTP header initialisation, triggered during cold-starts across 60 parallel goroutines. Implemented a singleton pattern fix; verified with Go's `-race` detector under load.
- **Amplitude analytics pipeline (end-to-end):** Led a 4-repository cross-service initiative spanning `courier`, `dp-protos`, `data-router`, and Amplitude — 1,892 additions including gRPC integration, Protobuf schema design, Kafka topic management, and event transformation, serving 30M+ events/month.
- **On-call engineer:** Monitored production systems, diagnosed incidents from first alert, proposed and implemented remediations; participated in load testing to identify API breaking points.
- **Iterable marketing integration:** Integrated Iterable platform via gRPC; built comprehensive event testing suite validating schemas and idempotency across the communications pipeline.
- **Golangci-lint CI:** Integrated 4–5 linters across 52 files (236 additions); SMSO service decommissioned (3,730 lines removed), reducing operational surface.
- **End-to-end ownership:** Took over feature delivery following senior departure — budget planning, task breakdown, technical investigation, and implementation across Kafka-based communications system.

### Backend Engineer (Placement) — DailyPay, Belfast
**Sep 2023 – Aug 2024**

*Go · Python · PostgreSQL · Docker · AWS · Kafka · GitHub Actions*

- **Schema Registry migration:** Executed cross-repo migration for `provider_account_created` and `earnings_invalidation` events spanning `core` and `ledger` — production schema changes coordinated across teams.
- **Multi-currency API support:** Implemented currency handling across GrossEarnings, TipsDetailed, and TipsEarningsProviderAccountDetails endpoints, enabling DailyPay's international expansion to Canada.
- **ETL pipelines (v1–v3):** Designed and shipped Gross Earnings and Available Balance pipelines using Go, PostgreSQL, Docker, and AWS; contributed to Ruby monolith → Go microservices migration.
- **V3 Earnings Job Status endpoint:** Delivered full OAPI spec, handler, and integration test suite.
- **Internal libraries & CI/CD:** Built reusable `daily-go` library components, improving system efficiency by 18%; created GitHub Actions pipelines and PostgreSQL cron scheduling scripts.

**Hackathon (DailyPay):** Winner (£500) — Led 3-person team; AR-based customer onboarding prototype recognised by CEO.

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
| **Messaging & Infra** | Kafka, gRPC, Protobuf, Schema Registry, Kubernetes, Docker, Terraform, AWS, GitHub Actions, Prometheus |
| **Databases** | PostgreSQL, Neo4j |
| **Practices** | Distributed systems, event-driven microservices, on-call, CI/CD, Agile/Kanban |

---

## Projects

### Bare Metal Kubernetes Cluster
- Built 4-node K8s cluster (3× Raspberry Pi + custom CPU) with Helm, k9s, nginx load balancing, 1.2TB NAS
- Automated deployments with reusable manifests; networking via ngrok and Glinet firewall
- [Video Demo →](https://millionairethinking.co.uk)

### Patent Knowledge Graph
- Scraping + NLP pipeline to map chemical patent relationships; semantic search via vector embeddings + graph traversal
- Python, Scrapy, Neo4j, LLMs
- [Video Demo →](https://millionairethinking.co.uk)
