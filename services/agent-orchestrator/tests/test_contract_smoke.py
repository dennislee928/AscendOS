from pathlib import Path
import sys
import unittest

SERVICE_ROOT = Path(__file__).resolve().parents[1]
if str(SERVICE_ROOT) not in sys.path:
    sys.path.insert(0, str(SERVICE_ROOT))

from app.contracts import build_orchestrator_state, validate_modules  # noqa: E402
from app.graph import GraphNode, PlaceholderLangGraph  # noqa: E402


class OrchestratorContractSmokeTest(unittest.TestCase):
    def test_validate_modules_rejects_invalid_or_duplicate_inputs(self) -> None:
        with self.assertRaises(ValueError):
            validate_modules([])

        with self.assertRaises(ValueError):
            validate_modules(["neuro", "neuro"])

        with self.assertRaises(ValueError):
            validate_modules(["neuro", "unknown"])

        self.assertEqual(validate_modules(["neuro", "kairos"]), ["neuro", "kairos"])

    def test_state_builder_and_graph_preserve_input_state(self) -> None:
        graph = PlaceholderLangGraph(
            [
                GraphNode(name="plan", run=lambda state: state),
                GraphNode(
                    name="append-output",
                    run=lambda state: {
                        **state,
                        "outputs": {**state["outputs"], "routed": state["modules"]},
                    },
                ),
            ]
        )
        state = build_orchestrator_state(
            "prepare a follow-up",
            ["neuro", "kairos"],
            {"seed": {"status": "kept"}},
        )

        result = graph.invoke(state)

        self.assertEqual(
            state,
            {
                "goal": "prepare a follow-up",
                "modules": ["neuro", "kairos"],
                "outputs": {"seed": {"status": "kept"}},
            },
        )
        self.assertEqual(result["goal"], "prepare a follow-up")
        self.assertEqual(result["modules"], ["neuro", "kairos"])
        self.assertEqual(result["outputs"]["seed"], {"status": "kept"})
        self.assertEqual(result["outputs"]["routed"], ["neuro", "kairos"])


if __name__ == "__main__":
    unittest.main()
