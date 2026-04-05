"""Configuration from environment variables."""

from pydantic_settings import BaseSettings, SettingsConfigDict
import os

class Settings(BaseSettings):
    """App settings loaded from env."""

    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    gemini_api_key: str | None = None
    openai_api_key: str | None = None
    local_llm_url: str | None = None # pilot used llm_url
    anthropic_api_key: str | None = None
    llm_provider: str = "openai"
    llm_model: str = "gpt-4o"

    # SmartExtract: when True, skip headful browser retry (safe for Docker, no X server)
    smartextract_headless_only: bool = True

def detect_llm_provider() -> tuple[str, str, str]:
    gemini = os.environ.get("GEMINI_API_KEY", "")
    openai = os.environ.get("OPENAI_API_KEY", "")
    local = os.environ.get("LLM_URL", "")
    model = os.environ.get("LLM_MODEL", "")

    if gemini and not local:
        return (
            "https://generativelanguage.googleapis.com/v1beta/openai",
            model or "gemini-2.0-flash",
            gemini,
        )
    if openai and not local:
        return (
            "https://api.openai.com/v1",
            model or "gpt-4o-mini",
            openai,
        )
    if local:
        return (
            local.rstrip("/"),
            model or "local-model",
            os.environ.get("LLM_API_KEY", ""),
        )

    raise RuntimeError(
        "No LLM provider configured. Set GEMINI_API_KEY, OPENAI_API_KEY, or LLM_URL."
    )