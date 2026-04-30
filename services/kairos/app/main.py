from datetime import UTC, datetime, timedelta
from typing import List

from fastapi import FastAPI
from pydantic import BaseModel, Field

app = FastAPI(title="kairos", version="0.1.0")


class BuildTimelineRequest(BaseModel):
    objective: str = Field(min_length=1, description="Goal to schedule")
    window_hours: int = Field(gt=0, le=168, description="Scheduling window in hours")
    step_count: int = Field(gt=0, le=10, description="Number of placeholder steps")


class TimelineStep(BaseModel):
    at: datetime
    action: str


class BuildTimelineResponse(BaseModel):
    objective: str
    steps: List[TimelineStep]


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "kairos"}


@app.post("/kairos/build-timeline", response_model=BuildTimelineResponse)
def build_timeline(payload: BuildTimelineRequest) -> BuildTimelineResponse:
    start = datetime.now(tz=UTC)
    delta = max(1, payload.window_hours // payload.step_count)
    steps = [
        TimelineStep(at=start + timedelta(hours=i * delta), action=f"{payload.objective} :: step {i + 1}")
        for i in range(payload.step_count)
    ]
    # TODO: Replace placeholder spacing with constraint-aware scheduling logic.
    return BuildTimelineResponse(objective=payload.objective, steps=steps)

