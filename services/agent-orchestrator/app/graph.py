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
        if "goal" not in state or not state["goal"]:
            raise ValueError("orchestrator state requires a non-empty goal")
        if "modules" not in state or not isinstance(state["modules"], list):
            raise ValueError("orchestrator state requires a module list")
        if "outputs" not in state or not isinstance(state["outputs"], dict):
            raise ValueError("orchestrator state requires an outputs mapping")

        # Work on a shallow copy so nodes cannot accidentally mutate caller-owned state.
        state = {
            "goal": state["goal"],
            "modules": list(state["modules"]),
            "outputs": dict(state["outputs"]),
        }
        for node in self._nodes:
            state = node.run(state)
        return state
