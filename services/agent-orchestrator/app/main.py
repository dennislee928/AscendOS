from typing import Any, Literal

from fastapi import FastAPI
from pydantic import BaseModel, ConfigDict, Field, field_validator

from .contracts import build_orchestrator_state, validate_modules
from .graph import GraphNode, OrchestratorState, PlaceholderLangGraph

app = FastAPI(title="agent-orchestrator", version="0.1.0")

ModuleName = Literal["neuro", "argentum", "kairos"]


class OrchestrateRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    goal: str = Field(min_length=1, description="User goal to route through module agents")
    modules: list[ModuleName] = Field(
        default_factory=lambda: ["neuro", "argentum", "kairos"],
        description="Requested modules in execution order",
    )

    @field_validator("modules")
    @classmethod
    def validate_modules(cls, modules: list[ModuleName]) -> list[ModuleName]:
        return validate_modules(modules)


class OrchestrateResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    goal: str
    modules: list[ModuleName]
    outputs: dict[str, Any]


def _plan_node(state: OrchestratorState) -> OrchestratorState:
    state["outputs"]["plan"] = {
        "goal": state["goal"],
        "module_order": state["modules"],
    }
    # TODO: Use an LLM planner node (prompt + guardrails + schema validation).
    return state


def _neuro_node(state: OrchestratorState) -> OrchestratorState:
    if "neuro" in state["modules"]:
        state["outputs"]["neuro"] = {
            "summary": f"placeholder neuro analysis for '{state['goal']}'",
            "confidence": 0.66,
        }
        # TODO: Call neuro service endpoint over internal RPC/HTTP and map result.
    return state


def _argentum_node(state: OrchestratorState) -> OrchestratorState:
    if "argentum" in state["modules"]:
        state["outputs"]["argentum"] = {
            "decision": "review",
            "reason": "placeholder risk evaluation pending real data",
        }
        # TODO: Call argentum service endpoint and apply policy checks.
    return state


def _kairos_node(state: OrchestratorState) -> OrchestratorState:
    if "kairos" in state["modules"]:
        state["outputs"]["kairos"] = {
            "next_action": "schedule follow-up",
            "eta_hours": 24,
        }
        # TODO: Call kairos service endpoint for timeline generation.
    return state


def _build_graph() -> PlaceholderLangGraph:
    # TODO: Replace linear node list with conditional/branching LangGraph edges.
    nodes = [
        GraphNode(name="plan", run=_plan_node),
        GraphNode(name="neuro", run=_neuro_node),
        GraphNode(name="argentum", run=_argentum_node),
        GraphNode(name="kairos", run=_kairos_node),
    ]
    return PlaceholderLangGraph(nodes=nodes)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "agent-orchestrator"}


@app.post("/orchestrate", response_model=OrchestrateResponse)
def orchestrate(payload: OrchestrateRequest) -> OrchestrateResponse:
    state: OrchestratorState = build_orchestrator_state(payload.goal, list(payload.modules))
    final = _build_graph().invoke(state)
    return OrchestrateResponse(goal=payload.goal, modules=payload.modules, outputs=final["outputs"])
