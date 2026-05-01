from dataclasses import dataclass
from typing import Any, Callable

from .contracts import OrchestratorState, build_orchestrator_state


NodeFn = Callable[[OrchestratorState], OrchestratorState]


@dataclass(frozen=True)
class GraphNode:
    name: str
    run: NodeFn


class PlaceholderLangGraph:
    """
    Minimal LangGraph-style executor.

    TODO: Replace with real `langgraph` StateGraph wiring once dependency and
    runtime orchestration decisions are finalized.
    """

    def __init__(self, nodes: list[GraphNode]) -> None:
        self._nodes = nodes

    def invoke(self, state: OrchestratorState) -> OrchestratorState:
        goal = state.get("goal", "")
        modules = state.get("modules")
        outputs = state.get("outputs")
        execution_summary = state.get("execution_summary")
        if not isinstance(modules, list):
            raise ValueError("orchestrator state requires a module list")
        if not isinstance(outputs, dict):
            raise ValueError("orchestrator state requires an outputs mapping")
        if not isinstance(execution_summary, dict):
            raise ValueError("orchestrator state requires an execution summary")

        state = build_orchestrator_state(goal, modules, outputs)
        for node in self._nodes:
            state = node.run(state)
        return state


def record_module_output(state: OrchestratorState, module: str, output: dict[str, Any]) -> OrchestratorState:
    state["outputs"][module] = output
    summary = state["execution_summary"]
    summary["module_status"][module] = "completed"
    summary["ordered_outputs"].append({"module": module, "status": "completed", "output": output})
    summary["completed_modules"].append(module)
    return state


def mark_module_failed(state: OrchestratorState, module: str, reason: str) -> OrchestratorState:
    summary = state["execution_summary"]
    summary["module_status"][module] = "failed"
    summary["ordered_outputs"].append(
        {"module": module, "status": "failed", "output": {"error": reason}}
    )
    return state
