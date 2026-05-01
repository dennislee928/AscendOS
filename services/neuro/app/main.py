import re
from collections import Counter
from dataclasses import dataclass
from typing import Callable, Literal


class SimpleApp:
    def __init__(self, title: str, version: str) -> None:
        self.title = title
        self.version = version
        self.routes: dict[tuple[str, str], Callable[..., object]] = {}

    def get(self, path: str) -> Callable[[Callable[..., object]], Callable[..., object]]:
        def decorator(func: Callable[..., object]) -> Callable[..., object]:
            self.routes[("GET", path)] = func
            return func

        return decorator

    def post(self, path: str, response_model: object | None = None) -> Callable[[Callable[..., object]], Callable[..., object]]:
        def decorator(func: Callable[..., object]) -> Callable[..., object]:
            self.routes[("POST", path)] = func
            return func

        return decorator


app = SimpleApp(title="neuro", version="0.1.0")


@dataclass(frozen=True, slots=True)
class NeuroAnalyzeRequest:
    text: str

    def __post_init__(self) -> None:
        if not isinstance(self.text, str) or not self.text.strip():
            raise ValueError("text must be a non-empty string")


@dataclass(frozen=True, slots=True)
class NeuroAnalyzeResponse:
    label: Literal["positive", "negative"]
    confidence: float
    embedding_preview: list[float]
    keywords: list[str]
    signals: dict[str, float]
    summary: str


POSITIVE_WORDS = {
    "achieve",
    "calm",
    "clear",
    "confident",
    "efficient",
    "enable",
    "excellent",
    "fast",
    "good",
    "great",
    "improve",
    "positive",
    "progress",
    "resolve",
    "stable",
    "success",
    "support",
}

NEGATIVE_WORDS = {
    "alert",
    "blocked",
    "delay",
    "difficult",
    "failure",
    "error",
    "fragile",
    "risk",
    "slow",
    "stress",
    "uncertain",
    "urgent",
    "weak",
    "worry",
}

URGENCY_WORDS = {"urgent", "asap", "now", "critical", "immediately", "fast"}
STOPWORDS = {
    "a",
    "an",
    "and",
    "are",
    "as",
    "be",
    "but",
    "for",
    "from",
    "if",
    "in",
    "into",
    "is",
    "it",
    "of",
    "on",
    "or",
    "our",
    "the",
    "to",
    "with",
    "we",
    "will",
}


def _tokenize(text: str) -> list[str]:
    return re.findall(r"[a-z0-9']+", text.lower())


def _keywords(tokens: list[str], limit: int = 5) -> list[str]:
    counts = Counter(token for token in tokens if len(token) > 2 and token not in STOPWORDS)
    return [token for token, _ in counts.most_common(limit)]


def _embedding_preview(tokens: list[str], positive_hits: int, negative_hits: int) -> list[float]:
    total = max(1, len(tokens))
    unique = len(set(tokens))
    urgency_hits = sum(1 for token in tokens if token in URGENCY_WORDS)
    avg_length = sum(len(token) for token in tokens) / total
    balance = (positive_hits - negative_hits) / max(1, positive_hits + negative_hits)
    return [
        round(total / 20.0, 4),
        round(unique / total, 4),
        round((positive_hits + 1) / (negative_hits + 1), 4),
        round(min(1.0, urgency_hits / total), 4),
        round(min(1.0, avg_length / 12.0 + (balance + 1.0) / 4.0), 4),
    ]


def _clamp(value: float, low: float, high: float) -> float:
    return max(low, min(high, value))


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "neuro"}


@app.post("/neuro/analyze")
def analyze(payload: NeuroAnalyzeRequest) -> NeuroAnalyzeResponse:
    tokens = _tokenize(payload.text)
    positive_hits = sum(1 for token in tokens if token in POSITIVE_WORDS)
    negative_hits = sum(1 for token in tokens if token in NEGATIVE_WORDS)
    urgency_hits = sum(1 for token in tokens if token in URGENCY_WORDS)
    total = max(1, len(tokens))
    unique = len(set(tokens))

    balance = (positive_hits - negative_hits) / max(1, positive_hits + negative_hits)
    confidence = _clamp(
        0.55
        + abs(balance) * 0.25
        + min(0.15, unique / total * 0.1 + urgency_hits / total * 0.05),
        0.0,
        1.0,
    )
    label = "positive" if balance >= 0 else "negative"
    keywords = _keywords(tokens)
    signals = {
        "positive_hits": float(positive_hits),
        "negative_hits": float(negative_hits),
        "urgency_hits": float(urgency_hits),
        "lexical_diversity": round(unique / total, 4),
        "sentiment_balance": round(balance, 4),
    }
    if keywords:
        summary = f"{label.title()} tone with focus on {', '.join(keywords[:3])}."
    else:
        summary = f"{label.title()} tone with no dominant topical keywords."
    return NeuroAnalyzeResponse(
        label=label,
        confidence=round(confidence, 4),
        embedding_preview=_embedding_preview(tokens, positive_hits, negative_hits),
        keywords=keywords,
        signals=signals,
        summary=summary,
    )
