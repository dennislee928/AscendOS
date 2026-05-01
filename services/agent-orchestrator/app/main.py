import re
from collections import Counter
from dataclasses import dataclass
from typing import Any, Callable, Literal

from .contracts import build_orchestrator_state, validate_modules
from .graph import GraphNode, OrchestratorState, PlaceholderLangGraph, record_module_output


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


app = SimpleApp(title="agent-orchestrator", version="0.1.0")

ModuleName = Literal["neuro", "argentum", "kairos"]


@dataclass(frozen=True, slots=True)
class OrchestrateRequest:
    goal: str
    modules: list[ModuleName] | None = None

    def __post_init__(self) -> None:
        if not isinstance(self.goal, str) or not self.goal.strip():
            raise ValueError("goal must be a non-empty string")
        if self.modules is None:
            object.__setattr__(self, "modules", ["neuro", "argentum", "kairos"])
        else:
            object.__setattr__(self, "modules", validate_modules(list(self.modules)))


@dataclass(frozen=True, slots=True)
class ModuleExecutionRecord:
    module: ModuleName
    status: Literal["completed", "skipped", "failed"]
    output: dict[str, Any]


@dataclass(frozen=True, slots=True)
class ExecutionSummary:
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


@dataclass(frozen=True, slots=True)
class OrchestrateResponse:
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


def _build_graph(modules: list[ModuleName]) -> PlaceholderLangGraph:
    node_map = {
        "neuro": GraphNode(name="neuro", run=_neuro_node),
        "argentum": GraphNode(name="argentum", run=_argentum_node),
        "kairos": GraphNode(name="kairos", run=_kairos_node),
    }
    nodes = [node_map[module] for module in modules]
    return PlaceholderLangGraph(nodes=nodes)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "agent-orchestrator"}


@app.post("/orchestrate")
def orchestrate(payload: OrchestrateRequest) -> OrchestrateResponse:
    state: OrchestratorState = build_orchestrator_state(payload.goal, list(payload.modules or []))
    final = _build_graph(payload.modules or []).invoke(state)
    summary = ExecutionSummary.from_state(final["execution_summary"])
    return OrchestrateResponse(
        goal=payload.goal,
        modules=list(payload.modules or []),
        outputs=final["outputs"],
        execution_summary=summary,
    )
