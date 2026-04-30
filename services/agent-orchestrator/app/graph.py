from dataclasses import dataclass
from typing import Any, Callable, TypedDict


class OrchestratorState(TypedDict):
    goal: str
    modules: list[str]
    outputs: dict[str, Any]


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
        for node in self._nodes:
            state = node.run(state)
        return state

