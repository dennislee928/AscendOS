from dataclasses import dataclass
from datetime import UTC, datetime, timedelta
from typing import Callable, Literal


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


app = SimpleApp(title="kairos", version="0.1.0")


@dataclass(frozen=True, slots=True)
class BuildTimelineRequest:
    objective: str
    window_hours: int
    step_count: int
    priority: int = 3
    constraints: list[str] = None  # type: ignore[assignment]

    def __post_init__(self) -> None:
        if not isinstance(self.objective, str) or not self.objective.strip():
            raise ValueError("objective must be a non-empty string")
        if self.window_hours <= 0 or self.window_hours > 168:
            raise ValueError("window_hours must be between 1 and 168")
        if self.step_count <= 0 or self.step_count > 10:
            raise ValueError("step_count must be between 1 and 10")
        if self.priority < 1 or self.priority > 5:
            raise ValueError("priority must be between 1 and 5")
        if self.constraints is None:
            object.__setattr__(self, "constraints", [])
        elif not isinstance(self.constraints, list):
            raise ValueError("constraints must be a list of strings")


@dataclass(frozen=True, slots=True)
class TimelineStep:
    at: datetime
    action: str
    kind: Literal["setup", "execution", "review", "buffer"]
    duration_hours: float
    notes: str


@dataclass(frozen=True, slots=True)
class BuildTimelineResponse:
    objective: str
    steps: list[TimelineStep]
    summary: str
    confidence: float
    critical_path_hours: float
    buffer_hours: float


def _normalize_objective(objective: str) -> list[str]:
    return [token for token in objective.lower().split() if token]


def _phase_templates(objective_tokens: list[str]) -> list[tuple[str, str, str]]:
    joined = " ".join(objective_tokens)
    if any(token in joined for token in ("deploy", "release", "launch")):
        return [
            ("setup", "prepare release scope", "align dependencies and release criteria"),
            ("execution", "deploy change set", "apply the main change in the window"),
            ("review", "verify post-release signals", "check health, errors, and user impact"),
        ]
    if any(token in joined for token in ("audit", "review", "inspect")):
        return [
            ("setup", "collect evidence", "assemble the records and controls to inspect"),
            ("execution", "sample target items", "review a representative subset"),
            ("review", "summarize exceptions", "capture follow-up items and ownership"),
        ]
    if any(token in joined for token in ("build", "implement", "ship")):
        return [
            ("setup", "define implementation slice", "confirm the smallest valuable increment"),
            ("execution", "implement core changes", "focus on the highest-leverage work"),
            ("review", "validate integration", "verify the result against the objective"),
        ]
    return [
        ("setup", "frame the objective", "clarify scope and success criteria"),
        ("execution", "work the primary task", "make measurable progress on the goal"),
        ("review", "confirm outcomes", "close the loop and capture follow-up work"),
    ]


def _generate_steps(payload: BuildTimelineRequest) -> list[TimelineStep]:
    start = datetime.now(tz=UTC)
    tokens = _normalize_objective(payload.objective)
    templates = _phase_templates(tokens)
    step_count = min(payload.step_count, 10)
    buffer_hours = max(1.0, round(payload.window_hours * (0.12 + 0.03 * (5 - payload.priority)), 2))
    available_hours = max(1.0, payload.window_hours - buffer_hours)
    duration_hours = round(available_hours / step_count, 2)
    cadence = max(1.0, available_hours / max(1, step_count - 1))

    steps: list[TimelineStep] = []
    for index in range(step_count):
        kind, action, note = templates[index % len(templates)]
        if index == step_count - 1 and step_count > 1:
            kind = "review" if kind != "review" else kind
            note = f"{note}; include buffer and sign-off"
        steps.append(
            TimelineStep(
                at=start + timedelta(hours=round(index * cadence, 2)),
                action=f"{payload.objective}: {action}",
                kind=kind,
                duration_hours=duration_hours,
                notes=note,
            )
        )
    return steps


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "kairos"}


@app.post("/kairos/build-timeline")
def build_timeline(payload: BuildTimelineRequest) -> BuildTimelineResponse:
    steps = _generate_steps(payload)
    buffer_hours = round(max(1.0, payload.window_hours * (0.12 + 0.03 * (5 - payload.priority))), 2)
    critical_path_hours = round(max(1.0, payload.window_hours - buffer_hours), 2)
    confidence = round(min(1.0, 0.58 + min(0.25, len(payload.constraints) * 0.05) + payload.priority * 0.03), 4)
    summary = (
        f"Planned {len(steps)} steps across a {payload.window_hours}-hour window "
        f"with {buffer_hours} hours reserved as buffer."
    )
    return BuildTimelineResponse(
        objective=payload.objective,
        steps=steps,
        summary=summary,
        confidence=confidence,
        critical_path_hours=critical_path_hours,
        buffer_hours=buffer_hours,
    )
