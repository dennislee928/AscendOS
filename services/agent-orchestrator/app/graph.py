from dataclasses import dataclass
from typing import Callable

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
        if not isinstance(modules, list):
            raise ValueError("orchestrator state requires a module list")
        if not isinstance(outputs, dict):
            raise ValueError("orchestrator state requires an outputs mapping")

        state = build_orchestrator_state(goal, modules, outputs)
        for node in self._nodes:
            state = node.run(state)
        return state
