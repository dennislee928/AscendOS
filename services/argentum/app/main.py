import math
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


app = SimpleApp(title="argentum", version="0.1.0")


@dataclass(frozen=True, slots=True)
class RiskEvaluateRequest:
    amount: float
    tenor_days: int
    counterparty_tier: int
    purpose: str | None = None
    collateral_coverage: float | None = None
    payment_history_score: float | None = None

    def __post_init__(self) -> None:
        if self.amount <= 0:
            raise ValueError("amount must be greater than zero")
        if self.tenor_days <= 0:
            raise ValueError("tenor_days must be greater than zero")
        if not 1 <= self.counterparty_tier <= 5:
            raise ValueError("counterparty_tier must be between 1 and 5")
        if self.collateral_coverage is not None and not 0.0 <= self.collateral_coverage <= 2.0:
            raise ValueError("collateral_coverage must be between 0.0 and 2.0")
        if self.payment_history_score is not None and not 0.0 <= self.payment_history_score <= 1.0:
            raise ValueError("payment_history_score must be between 0.0 and 1.0")


@dataclass(frozen=True, slots=True)
class RiskEvaluateResponse:
    risk_score: float
    band: Literal["low", "medium", "high"]
    decision: Literal["approve", "review", "escalate"]
    driver_scores: dict[str, float]
    controls: list[str]
    rationale: str
    recommended_limit: float


PURPOSE_RISK = {
    "acquisition": 0.16,
    "bridge": 0.18,
    "inventory": 0.08,
    "m&a": 0.2,
    "operations": 0.05,
    "refinance": 0.06,
    "restructuring": 0.14,
    "speculative": 0.22,
}


def _clamp(value: float, low: float, high: float) -> float:
    return max(low, min(high, value))


def _classify_band(score: float) -> Literal["low", "medium", "high"]:
    if score < 0.33:
        return "low"
    if score < 0.66:
        return "medium"
    return "high"


def _decision_from_band(band: Literal["low", "medium", "high"], tier: int) -> Literal["approve", "review", "escalate"]:
    if band == "low" and tier <= 2:
        return "approve"
    if band == "medium":
        return "review"
    return "escalate"


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "argentum"}


@app.post("/argentum/evaluate-risk")
def evaluate_risk(payload: RiskEvaluateRequest) -> RiskEvaluateResponse:
    amount_score = _clamp(math.log10(payload.amount + 10.0) / 7.0, 0.0, 1.0)
    tenor_score = _clamp(payload.tenor_days / 720.0, 0.0, 1.0)
    tier_score = _clamp((payload.counterparty_tier - 1) / 4.0, 0.0, 1.0)
    collateral_score = 0.0
    if payload.collateral_coverage is not None:
        collateral_score = _clamp(1.0 - min(payload.collateral_coverage, 1.5) / 1.5, 0.0, 1.0)
    history_score = 0.0
    if payload.payment_history_score is not None:
        history_score = _clamp(1.0 - payload.payment_history_score, 0.0, 1.0)

    purpose_score = 0.0
    if payload.purpose:
        lowered = payload.purpose.lower()
        purpose_score = max((weight for key, weight in PURPOSE_RISK.items() if key in lowered), default=0.0)

    score = _clamp(
        round(
            (0.24 * amount_score)
            + (0.18 * tenor_score)
            + (0.34 * tier_score)
            + (0.1 * collateral_score)
            + (0.08 * history_score)
            + (0.06 * purpose_score),
            4,
        ),
        0.0,
        1.0,
    )
    band = _classify_band(score)
    decision = _decision_from_band(band, payload.counterparty_tier)
    driver_scores = {
        "amount": round(amount_score, 4),
        "tenor": round(tenor_score, 4),
        "counterparty_tier": round(tier_score, 4),
        "collateral": round(collateral_score, 4),
        "payment_history": round(history_score, 4),
        "purpose": round(purpose_score, 4),
    }

    controls: list[str] = []
    if band != "low":
        controls.append("review covenant headroom")
    if payload.counterparty_tier >= 4:
        controls.append("require senior credit approval")
    if payload.collateral_coverage is None or payload.collateral_coverage < 1.0:
        controls.append("obtain additional collateral or guarantees")
    if payload.payment_history_score is None or payload.payment_history_score < 0.7:
        controls.append("confirm recent payment performance")
    if payload.tenor_days > 365:
        controls.append("add periodic monitoring checkpoints")

    rationale_parts = [
        f"amount contributes {driver_scores['amount']:.2f}",
        f"tenor contributes {driver_scores['tenor']:.2f}",
        f"counterparty tier contributes {driver_scores['counterparty_tier']:.2f}",
    ]
    if payload.purpose:
        rationale_parts.append(f"purpose '{payload.purpose}' contributes {driver_scores['purpose']:.2f}")
    if payload.collateral_coverage is not None:
        rationale_parts.append(f"collateral offsets risk by {1.0 - driver_scores['collateral']:.2f}")
    rationale = "; ".join(rationale_parts)
    recommended_limit = round(payload.amount * (1.0 - score) * 0.85, 2)

    return RiskEvaluateResponse(
        risk_score=score,
        band=band,
        decision=decision,
        driver_scores=driver_scores,
        controls=controls,
        rationale=rationale,
        recommended_limit=recommended_limit,
    )
