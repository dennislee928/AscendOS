<script lang="ts">
  import ModuleLanding from '../../lib/components/ModuleLanding.svelte';

  const layers = [
    {
      label: 'Revenue workflow',
      detail: 'Plan changes, checkout paths, and entitlement updates are presented as a single billing flow.',
    },
    {
      label: 'Provider wiring',
      detail: 'The module narrative ties into the canonical provider defaults and environment resolution used by deployment scripts.',
    },
    {
      label: 'Durable state',
      detail: 'Billing outcomes are expected to follow the same persisted service pattern that now backs the Go modules.',
    },
    {
      label: 'Release automation',
      detail: 'Promotion stays gated by the HF Space contract validator and the phase 9 readiness checks.',
    },
  ];

  const callouts = [
    {
      label: 'Launch status',
      detail: 'Positioned as the revenue control surface, ready to roll out once provider wiring is pointed at real services.',
    },
    {
      label: 'Stack fit',
      detail: 'Depends on release metadata, environment defaults, and the lower layers that keep deployments deterministic.',
    },
    {
      label: 'Product focus',
      detail: 'Use the page to explain how billing, plans, and entitlement changes are kept observable and auditable.',
    },
  ];
</script>

<ModuleLanding
  title="Billing Module"
  summary="Revenue and plan-management surface for checkout, entitlements, and the lifecycle of paid product access."
  phase="Phase 6"
  domain="Revenue"
  status="Rolling out"
  launch="Ready behind deployment contracts"
  ctaHref="/"
  ctaLabel="Back to launchpad"
  {layers}
  {callouts}
/>
