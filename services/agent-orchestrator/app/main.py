import re
from collections import Counter
from typing import Any, Literal

from fastapi import FastAPI
from pydantic import BaseModel, field_validator, model_validator

from .contracts import OrchestratorState, build_orchestrator_state, validate_modules
from .graph import make_compiled_graph, record_module_output

app = FastAPI(title="agent-orchestrator", version="0.1.0")

ModuleName = Literal["neuro", "argentum", "kairos"]


class OrchestrateRequest(BaseModel):
    goal: str
    modules: list[ModuleName] | None = None

    @field_validator("goal")
    @classmethod
    def goal_nonempty(cls, value: str) -> str:
        if not value.strip():
            raise ValueError("goal must be a non-empty string")
        return value

    @model_validator(mode="after")
    def normalize_modules(self) -> "OrchestrateRequest":
        if self.modules is None:
            self.modules = list(["neuro", "argentum", "kairos"])
        else:
            self.modules = validate_modules(list(self.modules))
        return self


class ModuleExecutionRecord(BaseModel):
    module: ModuleName
    status: Literal["completed", "skipped", "failed"]
    output: dict[str, Any]


class ExecutionSummary(BaseModel):
    module_status: dict[str, Literal["pending", "completed", "skipped", "failed"]]
    ordered_outputs: list[ModuleExecutionRecord]
    completed_modules: list[ModuleName]
    skipped_modules: list[ModuleName]

    @classmethod
    def from_state(cls, state: dict[str, Any]) -> "ExecutionSummary":
        ordered_outputs = [
            ModuleExecutionRecord(
                module=entry["module"],
                status=entry["status"],
                output=dict(entry["output"]),
            )
            for entry in state["ordered_outputs"]
        ]
        return cls(
            module_status=dict(state["module_status"]),
            ordered_outputs=ordered_outputs,
            completed_modules=list(state["completed_modules"]),
            skipped_modules=list(state["skipped_modules"]),
        )


class OrchestrateResponse(BaseModel):
    goal: str
    modules: list[ModuleName]
    outputs: dict[str, Any]
    execution_summary: ExecutionSummary


POSITIVE_WORDS = {
    "aligned",
    "clear",
    "confident",
    "efficient",
    "focused",
    "green",
    "progress",
    "ready",
    "stable",
    "success",
    "verified",
}

RISK_WORDS = {
    "delay",
    "blocked",
    "fragile",
    "risk",
    "uncertain",
    "urgent",
    "volatile",
}

SCHEDULING_WORDS = {
    "audit",
    "build",
    "deploy",
    "launch",
    "monitor",
    "plan",
    "review",
    "stabilize",
}


def _tokens(goal: str) -> list[str]:
    return re.findall(r"[a-z0-9']+", goal.lower())


def _keyword_slice(tokens: list[str], limit: int = 4) -> list[str]:
    counts = Counter(token for token in tokens if len(token) > 2)
    return [token for token, _ in counts.most_common(limit)]


def _module_available(state: OrchestratorState, module: ModuleName) -> bool:
    return module in state["modules"]


def _neuro_node(state: OrchestratorState) -> OrchestratorState:
    if not _module_available(state, "neuro"):
        return state
    tokens = _tokens(state["goal"])
    positive = sum(1 for token in tokens if token in POSITIVE_WORDS)
    risk = sum(1 for token in tokens if token in RISK_WORDS)
    confidence = round(min(1.0, 0.58 + abs(positive - risk) * 0.08 + len(set(tokens)) * 0.01), 4)
    output = {
        "label": "positive" if positive >= risk else "negative",
        "confidence": confidence,
        "keywords": _keyword_slice(tokens),
        "signals": {
            "positive_hits": float(positive),
            "risk_hits": float(risk),
            "lexical_diversity": round(len(set(tokens)) / max(1, len(tokens)), 4),
        },
        "summary": "Goal framing suggests "
        + ("forward momentum" if positive >= risk else "active risk management")
        + ".",
    }
    return record_module_output(state, "neuro", output)


def _argentum_node(state: OrchestratorState) -> OrchestratorState:
    if not _module_available(state, "argentum"):
        return state
    tokens = _tokens(state["goal"])
    exposure = min(1.0, len(tokens) / 14.0)
    urgency = 1.0 if any(token in {"urgent", "asap", "fast"} for token in tokens) else 0.2
    risk_score = round(min(1.0, 0.28 + exposure * 0.35 + urgency * 0.18), 4)
    band = "low" if risk_score < 0.33 else "medium" if risk_score < 0.66 else "high"
    decision = "approve" if band == "low" else "review" if band == "medium" else "escalate"
    output = {
        "risk_score": risk_score,
        "band": band,
        "decision": decision,
        "driver_scores": {
            "exposure": round(exposure, 4),
            "urgency": round(urgency, 4),
            "goal_length": float(len(tokens)),
        },
        "controls": [
            "review funding assumptions" if band != "low" else "standard approval controls",
            "confirm downstream capacity" if urgency > 0.2 else "track normal operating cadence",
        ],
        "rationale": f"Goal length and urgency indicate a {band} risk posture.",
        "recommended_limit": round(max(0.0, 1_000_000 * (1.0 - risk_score) * 0.75), 2),
    }
    return record_module_output(state, "argentum", output)


def _kairos_node(state: OrchestratorState) -> OrchestratorState:
    if not _module_available(state, "kairos"):
        return state
    tokens = _tokens(state["goal"])
    objective_type = next((token for token in tokens if token in SCHEDULING_WORDS), "plan")
    step_count = min(5, max(3, len(tokens) // 2 or 3))
    steps = [
        {
            "at_offset_hours": index * 4,
            "action": f"{objective_type} phase {index + 1}",
            "kind": "execution" if index else "setup",
        }
        for index in range(step_count)
    ]
    output = {
        "objective_type": objective_type,
        "steps": steps,
        "summary": f"Sequenced {step_count} steps for '{state['goal']}'.",
        "confidence": round(min(1.0, 0.6 + step_count * 0.05), 4),
        "critical_path_hours": float(step_count * 4),
        "buffer_hours": 4.0,
    }
    return record_module_output(state, "kairos", output)


def _build_graph(modules: list[ModuleName]) -> Any:
    node_map = {
        "neuro": ("neuro", _neuro_node),
        "argentum": ("argentum", _argentum_node),
        "kairos": ("kairos", _kairos_node),
    }
    node_sequence = [node_map[module] for module in modules if module in node_map]
    return make_compiled_graph(node_sequence)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "agent-orchestrator"}


@app.post("/orchestrate", response_model=OrchestrateResponse)
def orchestrate(payload: OrchestrateRequest) -> OrchestrateResponse:
    modules = list(payload.modules or [])
    state: OrchestratorState = build_orchestrator_state(payload.goal, modules)
    final = _build_graph(modules).invoke(state)
    summary = ExecutionSummary.from_state(final["execution_summary"])
    return OrchestrateResponse(
        goal=payload.goal,
        modules=modules,
        outputs=final["outputs"],
        execution_summary=summary,
    )
