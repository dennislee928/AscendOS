from pathlib import Path
import sys
import unittest

AGENT_ORCHESTRATOR_ROOT = Path(__file__).resolve().parents[1]
SERVICES_ROOT = AGENT_ORCHESTRATOR_ROOT.parent


def _load_module(service_name: str):
    service_root = SERVICES_ROOT / service_name
    sys.path.insert(0, str(service_root))
    try:
        for key in list(sys.modules):
            if key == "app" or key.startswith("app."):
                del sys.modules[key]
        import app.main as module  # type: ignore

        return module
    finally:
        sys.path.remove(str(service_root))


neuro_module = _load_module("neuro")
argentum_module = _load_module("argentum")
kairos_module = _load_module("kairos")
orchestrator_module = _load_module("agent-orchestrator")


class OrchestratorContractSmokeTest(unittest.TestCase):
    def test_health_and_request_response_paths(self) -> None:
        self.assertEqual(neuro_module.healthz(), {"status": "ok", "service": "neuro"})
        self.assertEqual(argentum_module.healthz(), {"status": "ok", "service": "argentum"})
        self.assertEqual(kairos_module.healthz(), {"status": "ok", "service": "kairos"})
        self.assertEqual(orchestrator_module.healthz(), {"status": "ok", "service": "agent-orchestrator"})

        neuro_response = neuro_module.analyze(
            neuro_module.NeuroAnalyzeRequest(
                text="We need a clear and stable plan to resolve the issue and support progress.",
            )
        )
        self.assertIn(neuro_response.label, {"positive", "negative"})
        self.assertTrue(neuro_response.summary)
        self.assertEqual(len(neuro_response.embedding_preview), 5)

        argentum_response = argentum_module.evaluate_risk(
            argentum_module.RiskEvaluateRequest(
                amount=250000,
                tenor_days=180,
                counterparty_tier=3,
                purpose="operations refinance",
                collateral_coverage=1.1,
                payment_history_score=0.82,
            )
        )
        self.assertIn(argentum_response.band, {"low", "medium", "high"})
        self.assertIn(argentum_response.decision, {"approve", "review", "escalate"})
        self.assertTrue(argentum_response.controls)
        self.assertGreater(argentum_response.recommended_limit, 0)

        kairos_response = kairos_module.build_timeline(
            kairos_module.BuildTimelineRequest(
                objective="Deploy the release and verify post-launch stability",
                window_hours=24,
                step_count=4,
                priority=4,
                constraints=["avoid weekend work", "keep a rollback checkpoint"],
            )
        )
        self.assertGreaterEqual(len(kairos_response.steps), 4)
        self.assertTrue(kairos_response.summary)
        self.assertGreater(kairos_response.critical_path_hours, 0)

        orchestrator_response = orchestrator_module.orchestrate(
            orchestrator_module.OrchestrateRequest(
                goal="Coordinate a clear deployment plan with risk review and timeline follow-up",
                modules=["neuro", "argentum", "kairos"],
            )
        )
        self.assertEqual(orchestrator_response.modules, ["neuro", "argentum", "kairos"])
        self.assertEqual(list(orchestrator_response.outputs.keys()), ["neuro", "argentum", "kairos"])

        summary = orchestrator_response.execution_summary
        self.assertEqual(summary.completed_modules, ["neuro", "argentum", "kairos"])
        self.assertEqual(summary.skipped_modules, [])
        self.assertEqual([entry.module for entry in summary.ordered_outputs], ["neuro", "argentum", "kairos"])
        self.assertEqual(summary.module_status["neuro"], "completed")
        self.assertEqual(summary.module_status["argentum"], "completed")
        self.assertEqual(summary.module_status["kairos"], "completed")


if __name__ == "__main__":
    unittest.main()
