from typing import Any, Callable

from langgraph.graph import END, START, StateGraph

from .contracts import OrchestratorState


NodeFn = Callable[[OrchestratorState], OrchestratorState]


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


def make_compiled_graph(node_sequence: list[tuple[str, NodeFn]]) -> Any:
    if not node_sequence:
        raise ValueError("at least one node is required")

    builder: StateGraph = StateGraph(OrchestratorState)
    for name, fn in node_sequence:
        builder.add_node(name, fn)

    names = [name for name, _ in node_sequence]
    builder.add_edge(START, names[0])
    for i in range(len(names) - 1):
        builder.add_edge(names[i], names[i + 1])
    builder.add_edge(names[-1], END)

    return builder.compile()
