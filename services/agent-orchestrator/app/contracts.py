from typing import Any, Literal, TypedDict

ModuleName = str
ALLOWED_MODULES = ("neuro", "argentum", "kairos")
ModuleStatus = Literal["pending", "completed", "skipped", "failed"]


class ExecutionEntry(TypedDict):
    module: str
    status: ModuleStatus
    output: dict[str, Any]


class ExecutionSummaryState(TypedDict):
    module_status: dict[str, ModuleStatus]
    ordered_outputs: list[ExecutionEntry]
    completed_modules: list[str]
    skipped_modules: list[str]


class OrchestratorState(TypedDict):
    goal: str
    modules: list[str]
    outputs: dict[str, Any]
    execution_summary: ExecutionSummaryState


def validate_modules(modules: list[str]) -> list[str]:
    if not modules:
        raise ValueError("modules must include at least one service")
    if len(modules) != len(set(modules)):
        raise ValueError("modules must not contain duplicates")
    invalid = [module for module in modules if module not in ALLOWED_MODULES]
    if invalid:
        raise ValueError(f"modules contains unsupported services: {', '.join(invalid)}")
    return list(modules)


def build_orchestrator_state(
    goal: str,
    modules: list[str],
    outputs: dict[str, Any] | None = None,
) -> OrchestratorState:
    if not goal:
        raise ValueError("orchestrator state requires a non-empty goal")
    module_status = {module: "pending" for module in ALLOWED_MODULES}
    skipped_modules = [module for module in ALLOWED_MODULES if module not in modules]
    for module in skipped_modules:
        module_status[module] = "skipped"
    return {
        "goal": goal,
        "modules": validate_modules(modules),
        "outputs": dict(outputs or {}),
        "execution_summary": {
            "module_status": module_status,
            "ordered_outputs": [],
            "completed_modules": [],
            "skipped_modules": skipped_modules,
        },
    }
