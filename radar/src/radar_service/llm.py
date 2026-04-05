"""Unified LLM client for radar. ApplyPilot-style: Gemini, OpenAI, local endpoint.

Provides get_client() -> LLMClient with ask(prompt, temperature, max_tokens).
Used by SmartExtract and any code needing simple prompt/response.
"""

import logging
import os
import time

import httpx

from radar_service.config import Settings, detect_llm_provider

log = logging.getLogger(__name__)

_MAX_RETRIES = 5
_TIMEOUT = 120
_RATE_LIMIT_BASE_WAIT = 10

_GEMINI_COMPAT_BASE = "https://generativelanguage.googleapis.com/v1beta/openai"
_GEMINI_NATIVE_BASE = "https://generativelanguage.googleapis.com/v1beta"


class _GeminiCompatForbidden(Exception):
    """Sentinel: Gemini OpenAI-compat returned 403."""
    def __init__(self, response: httpx.Response) -> None:
        self.response = response
        super().__init__(f"Gemini compat 403: {response.text[:200]}")


class LLMClient:
    """Thin LLM client: OpenAI-compat and native Gemini. Same as ApplyPilot."""

    def __init__(self, base_url: str, model: str, api_key: str) -> None:
        self.base_url = base_url
        self.model = model
        self.api_key = api_key
        self._client = httpx.Client(timeout=_TIMEOUT)
        self._use_native_gemini = False
        self._is_gemini = base_url.startswith(_GEMINI_COMPAT_BASE)

    def _chat_native_gemini(
        self, messages: list[dict], temperature: float, max_tokens: int
    ) -> str:
        contents = []
        system_parts = []
        for msg in messages:
            role = msg["role"]
            text = msg.get("content", "")
            if role == "system":
                system_parts.append({"text": text})
            elif role == "user":
                contents.append({"role": "user", "parts": [{"text": text}]})
            elif role == "assistant":
                contents.append({"role": "model", "parts": [{"text": text}]})

        payload = {
            "contents": contents,
            "generationConfig": {"temperature": temperature, "maxOutputTokens": max_tokens},
        }
        if system_parts:
            payload["systemInstruction"] = {"parts": system_parts}

        url = f"{_GEMINI_NATIVE_BASE}/models/{self.model}:generateContent"
        resp = self._client.post(
            url, json=payload,
            headers={"Content-Type": "application/json"},
            params={"key": self.api_key},
        )
        resp.raise_for_status()
        return resp.json()["candidates"][0]["content"]["parts"][0]["text"]

    def _chat_compat(
        self, messages: list[dict], temperature: float, max_tokens: int
    ) -> str:
        headers = {"Content-Type": "application/json"}
        if self.api_key:
            headers["Authorization"] = f"Bearer {self.api_key}"
        payload = {
            "model": self.model,
            "messages": messages,
            "temperature": temperature,
            "max_tokens": max_tokens,
        }
        resp = self._client.post(
            f"{self.base_url}/chat/completions",
            json=payload,
            headers=headers,
        )
        if resp.status_code == 403 and self._is_gemini:
            raise _GeminiCompatForbidden(resp)
        resp.raise_for_status()
        return resp.json()["choices"][0]["message"]["content"]

    def chat(
        self,
        messages: list[dict],
        temperature: float = 0.0,
        max_tokens: int = 4096,
    ) -> str:
        for attempt in range(_MAX_RETRIES):
            try:
                if self._use_native_gemini:
                    return self._chat_native_gemini(messages, temperature, max_tokens)
                return self._chat_compat(messages, temperature, max_tokens)
            except _GeminiCompatForbidden:
                log.warning("Gemini compat 403, switching to native API")
                self._use_native_gemini = True
                return self._chat_native_gemini(messages, temperature, max_tokens)
            except httpx.HTTPStatusError as exc:
                if exc.response.status_code in (429, 503) and attempt < _MAX_RETRIES - 1:
                    wait = min(_RATE_LIMIT_BASE_WAIT * (2 ** attempt), 60)
                    log.warning("LLM rate limited, retry in %ds", wait)
                    time.sleep(wait)
                    continue
                raise
            except httpx.TimeoutException:
                if attempt < _MAX_RETRIES - 1:
                    time.sleep(min(_RATE_LIMIT_BASE_WAIT * (2 ** attempt), 60))
                    continue
                raise
        raise RuntimeError("LLM request failed after all retries")

    def ask(self, prompt: str, temperature: float = 0.0, max_tokens: int = 4096) -> str:
        """Single user prompt -> response. Used by SmartExtract."""
        return self.chat([{"role": "user", "content": prompt}], temperature, max_tokens)


_instance: LLMClient | None = None


def get_client() -> LLMClient:
    """Return singleton LLMClient. Uses GEMINI > OPENAI > LLM_URL."""
    global _instance
    if _instance is None:
        base_url, model, api_key = detect_llm_provider()
        log.info("LLM provider: %s model: %s", base_url, model)
        _instance = LLMClient(base_url, model, api_key)
    return _instance


def get_crawl4ai_llm_config():
    """Return crawl4ai LLMConfig for the detected provider. Used by job_search."""
    from crawl4ai import LLMConfig

    base_url, model, api_key = detect_llm_provider()
    if "generativelanguage" in base_url:
        provider = f"gemini/{model}"
    elif "api.openai.com" in base_url:
        provider = f"openai/{model}"
    else:
        provider = f"openai/{model}"
        return LLMConfig(provider=provider, api_token=api_key or "no-token", base_url=base_url)
    return LLMConfig(provider=provider, api_token=api_key)