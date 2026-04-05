"""Job extraction strategy for Crawl4ai."""

from pydantic import BaseModel, Field

from radar_service.config import Settings
from crawl4ai import JsonCssExtractionStrategy, LLMConfig, LLMExtractionStrategy

# Indeed UK/US job card selectors (from mosaic-provider-jobcards DOM)
INDEED_CSS_SCHEMA = {
    "name": "indeed_jobs",
    "baseSelector": "div.job_seen_beacon",
    "fields": [
        {
            "name": "job_key",
            "selector": "h2.jobTitle a[data-jk]",
            "type": "attribute",
            "attribute": "data-jk",
        },
        {
            "name": "title",
            "selector": "h2.jobTitle span",
            "type": "text",
        },
        {
            "name": "company",
            "selector": "span[data-testid='company-name']",
            "type": "text",
        },
        {
            "name": "location",
            "selector": "div[data-testid='text-location']",
            "type": "text",
        },
        {
            "name": "job_url",
            "selector": "h2.jobTitle a[href]",
            "type": "attribute",
            "attribute": "href",
        },
    ],
}

# Indeed job detail page: full description (viewjob page). Layout varies by locale/A-B tests.
INDEED_JOB_DETAIL_SCHEMA = {
    "name": "job_detail",
    "baseSelector": "body",
    "fields": [
        {
            "name": "description",
            "selector": "#jobDescriptionText",
            "type": "text",
        },
        {
            "name": "description_fallback",
            "selector": "div.jobsearch-jobDescriptionText",
            "type": "text",
        },
    ],
}


class JobExtractionItem(BaseModel):
    """Single job for LLM extraction."""

    job_url: str = Field(description="Full URL to the job posting")
    title: str = Field(description="Job title")
    company: str = Field(description="Company name")
    requirements: list[str] = Field(
        default_factory=list,
        description="Key requirements/skills from the listing",
    )


class JobExtractionSchema(BaseModel):
    """Root schema for job list extraction."""

    jobs: list[JobExtractionItem] = Field(
        description="List of job listings found on the page",
    )


def get_indeed_css_extraction_strategy() -> JsonCssExtractionStrategy:
    """Return JsonCssExtractionStrategy for Indeed job cards (no LLM)."""
    return JsonCssExtractionStrategy(schema=INDEED_CSS_SCHEMA, verbose=False)


def get_indeed_job_detail_extraction_strategy() -> JsonCssExtractionStrategy:
    """Return JsonCssExtractionStrategy for Indeed job detail page (full description)."""
    return JsonCssExtractionStrategy(schema=INDEED_JOB_DETAIL_SCHEMA, verbose=False)

def get_job_extraction_strategy() -> LLMExtractionStrategy:
    """Return LLMExtractionStrategy using ApplyPilot-style LLM config."""
    from radar_service.llm import get_crawl4ai_llm_config

    llm_config = get_crawl4ai_llm_config()
    return LLMExtractionStrategy(
        llm_config=llm_config,
        schema=JobExtractionSchema.model_json_schema(),
        extraction_type="schema",
        instruction="""Extract all job listings from this page. Each job has:
- job_url: full URL to the job posting (must be absolute, e.g. https://...)
- title: job title
- company: company name
- requirements: list of key skills/requirements from the listing (can be empty)

Return a JSON object with a "jobs" array. Do not include jobs that are ads or duplicates.""",
        chunk_token_threshold=4000,
        overlap_rate=0.1,
        apply_chunking=True,
        input_format="markdown",
        extra_args={"temperature": 0.1, "max_tokens": 2000},
        verbose=False,
    )
