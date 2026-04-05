#!/usr/bin/env -S uv run python
"""Test FindJobs RPC: read resume PDF, call FindJobs, print results.

Usage:
  make radar-test-find-jobs RESUME=/path/to/resume.pdf
  uv run python scripts/test_find_jobs.py /path/to/resume.pdf

Requires: gRPC server running (make radar-run)
"""

import sys

from pypdf import PdfReader

import grpc
from radar.services.job_search.v1 import service_reads_pb2, service_definition_pb2_grpc


def extract_resume_text(pdf_path: str) -> str:
    reader = PdfReader(pdf_path)
    return "\n".join(page.extract_text() or "" for page in reader.pages)


def main():
    if len(sys.argv) < 2:
        print("Usage: test_find_jobs.py <resume.pdf> [search_query] [location]")
        sys.exit(1)

    pdf_path = sys.argv[1]
    search_query = sys.argv[2] if len(sys.argv) > 2 else "software engineer"
    location = sys.argv[3] if len(sys.argv) > 3 else ""

    print(f"Reading resume from {pdf_path}...")
    resume_text = extract_resume_text(pdf_path)
    print(f"Extracted {len(resume_text)} chars from resume")

    channel = grpc.insecure_channel("localhost:50051")
    stub = service_definition_pb2_grpc.JobSearchServiceStub(channel)

    req = service_reads_pb2.FindJobsRequest(
        resume_text=resume_text,
        target_roles=["Software Engineer", "Product Manager"],
        search_query=search_query,
        location=location,
        max_results=5,
    )

    print(f"\nCalling FindJobs (query={search_query}, location={location})...")
    resp = stub.FindJobs(req)

    print(f"\nFound {len(resp.jobs)} jobs:\n")
    for i, job in enumerate(resp.jobs, 1):
        print(f"{i}. {job.title} @ {job.company}")
        print(f"   URL: {job.job_url}")
        if job.requirements:
            print(f"   Requirements: {job.requirements[:3]}...")
        print()


if __name__ == "__main__":
    main()
