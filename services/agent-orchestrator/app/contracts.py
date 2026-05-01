from typing import Any, TypedDict

ModuleName = str
ALLOWED_MODULES = ("neuro", "argentum", "kairos")


class OrchestratorState(TypedDict):
    goal: str
    modules: list[str]
    outputs: dict[str, Any]


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
    return {
        "goal": goal,
        "modules": validate_modules(modules),
        "outputs": dict(outputs or {}),
    }
