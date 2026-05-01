from hashlib import sha256
from typing import List, Literal

from fastapi import FastAPI
from pydantic import BaseModel, ConfigDict, Field

app = FastAPI(title="neuro", version="0.1.0")


class NeuroAnalyzeRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    text: str = Field(min_length=1, description="Raw input text for placeholder analysis")


class NeuroAnalyzeResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    label: Literal["positive", "negative"]
    confidence: float = Field(ge=0.0, le=1.0)
    embedding_preview: List[float] = Field(min_length=5, max_length=5)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "neuro"}


@app.post("/neuro/analyze", response_model=NeuroAnalyzeResponse)
def analyze(payload: NeuroAnalyzeRequest) -> NeuroAnalyzeResponse:
    digest = sha256(payload.text.encode("utf-8")).hexdigest()
    score = int(digest[:8], 16) / 0xFFFFFFFF
    label = "positive" if score >= 0.5 else "negative"
    embedding_preview = [round(int(digest[i : i + 2], 16) / 255, 4) for i in range(0, 10, 2)]
    # TODO: Replace hash-based placeholder with a real model inference pipeline.
    return NeuroAnalyzeResponse(label=label, confidence=round(score, 4), embedding_preview=embedding_preview)
