"""JobSearchService gRPC servicer implementation."""

from radar.services.job_search.v1 import (
    models_pb2,
    service_definition_pb2,
    service_definition_pb2_grpc,
    service_reads_pb2,
    service_writes_pb2,
)

from radar_service.job_search.find_jobs import find_jobs


class JobSearchServicer(service_definition_pb2_grpc.JobSearchServiceServicer):
    """Implements JobSearchService RPCs."""

    def Health(self, request, context):
        return service_definition_pb2.HealthResponse(status="ok")

    def FindJobs(self, request, context):
        jobs_data = find_jobs(
            resume_text=request.resume_text,
            target_roles=list(request.target_roles),
            search_query=request.search_query,
            location=request.location,
            max_results=request.max_results or 10,
            run_uuid=(request.run_uuid or "").strip(),
        )
        jobs = [
            models_pb2.JobInfo(
                job_url=j["job_url"],
                title=j["title"],
                company=j["company"],
                requirements=j.get("requirements", []),
            )
            for j in jobs_data
        ]
        return service_reads_pb2.FindJobsResponse(jobs=jobs)

    def AnalyzeJob(self, request, context):
        # TODO: Implement using browser-use Agent
        return service_reads_pb2.AnalyzeJobResponse(
            fit_score=0.0,
            summary="",
            extracted_requirements=[],
            recommendation="",
        )

    def ApplyToJob(self, request, context):
        # TODO: Implement using browser-use Agent
        return service_writes_pb2.ApplyToJobResponse(
            success=False,
            result="",
            error="Not implemented",
        )
