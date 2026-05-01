<script lang="ts">
  import ModuleLanding from '../../lib/components/ModuleLanding.svelte';

  const layers = [
    {
      label: 'Frontend shell',
      detail: 'The Svelte route presents the launchpad surface that operators use to inspect module readiness.',
    },
    {
      label: 'Operational services',
      detail: 'Dashboard signals sit on top of the persisted Go services and richer Python module outputs.',
    },
    {
      label: 'Observability',
      detail: 'Readiness and rollout checks align with the Prometheus and Grafana stack from the phase 7 layer.',
    },
    {
      label: 'Release automation',
      detail: 'Deployment metadata, quality gates, and contract checks decide when the dashboard can move forward.',
    },
  ];

  const callouts = [
    {
      label: 'Launch status',
      detail: 'Ready for review as a Phase 6 marketing surface with live metadata instead of a placeholder.',
    },
    {
      label: 'Stack fit',
      detail: 'Maps cleanly to the shared launchpad catalog exported from packages/ui and the release gates beneath it.',
    },
    {
      label: 'Operator value',
      detail: 'Best used to verify module health, rollout order, and whether the platform is safe to promote.',
    },
  ];
</script>

<ModuleLanding
  title="Dashboard Module"
  summary="Command-center entry point for module health, rollout readiness, and the current operating state of the product."
  phase="Phase 6"
  domain="Operations"
  status="Launch-ready"
  launch="Surface aligned with shared launchpad routes"
  ctaHref="/"
  ctaLabel="Back to launchpad"
  {layers}
  {callouts}
/>
