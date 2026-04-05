# Cover Letter — Backend Engineer III, Monzo

**Kanishk Mittal**
kanishkable@gmail.com · 07899007443 · linkedin.com/in/kanishkable

---

Hiring Team, Monzo

I'm applying for the Backend Engineer III role at Monzo. I've been using Monzo for several years and I'm applying because the engineering problems are real — not just at the product layer but at the infrastructure level: resilient distributed systems, event-driven communication at scale, and on-call ownership with meaningful consequences. That's exactly the work I've been doing at DailyPay.

In my current role, I was paged on a production panic in our Go communications service. The error was non-deterministic and had been intermittent for weeks. I traced it to a concurrent map write inside the Schema Registry client's HTTP header initialisation — triggered only during cold-starts when 60 goroutines hit `NewSerializer` simultaneously. The fix was a singleton pattern with `sync.Once`, and I verified it with Go's `-race` detector under simulated load before merging. The panic stopped. That kind of investigation — systematic, evidence-based, verified — is what I bring to on-call and production work. I believe it maps directly to how Monzo's squads approach reliability.

Beyond that specific debugging win, the broader shape of my work at DailyPay matches Monzo's stack closely: Go, Kafka (30M+ events/month), gRPC, PostgreSQL, Kubernetes, Docker, Terraform, and GitHub Actions CI/CD — all in production. My most recent large piece of work was a 4-repository cross-service pipeline connecting `courier` → `dp-protos` → `data-router` → Amplitude, including Protobuf schema design, Kafka topic management, and event transformation (1,892 additions). I owned it end-to-end following a senior's departure, from scoping to deployment. On the review side, I've conducted 190+ code reviews across 10+ teammates in nine months — not LGTM reviews, but ones that caught proto field mapping bugs, flagged Terraform permission scope, and proposed architectural patterns. I mention this because Monzo's squads model means engineers do a lot of this kind of lateral work, and I'm comfortable with it.

I'm aware my total experience is ~21 months, which is on the lower end for L40. I think the relevant proxy isn't tenure — it's whether the output and ownership match the level, and I believe mine does. I'd welcome the chance to discuss that. I'm particularly interested in the Platform or Payments collective, where the event-driven and distributed systems work maps most directly to my background.

Kanishk Mittal
