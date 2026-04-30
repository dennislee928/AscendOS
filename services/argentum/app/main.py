from fastapi import FastAPI
from pydantic import BaseModel, Field

app = FastAPI(title="argentum", version="0.1.0")


class RiskEvaluateRequest(BaseModel):
    amount: float = Field(gt=0, description="Requested amount")
    tenor_days: int = Field(gt=0, description="Duration in days")
    counterparty_tier: int = Field(ge=1, le=5, description="1 is strongest, 5 is weakest")


class RiskEvaluateResponse(BaseModel):
    risk_score: float
    band: str


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "argentum"}


@app.post("/argentum/evaluate-risk", response_model=RiskEvaluateResponse)
def evaluate_risk(payload: RiskEvaluateRequest) -> RiskEvaluateResponse:
    base = payload.amount / 100_000
    tenor_factor = payload.tenor_days / 365
    tier_factor = payload.counterparty_tier / 5
    score = min(1.0, round((0.5 * base) + (0.2 * tenor_factor) + (0.3 * tier_factor), 4))
    if score < 0.35:
        band = "low"
    elif score < 0.7:
        band = "medium"
    else:
        band = "high"
    # TODO: Replace deterministic placeholder with policy/rules engine + model-assisted scoring.
    return RiskEvaluateResponse(risk_score=score, band=band)

