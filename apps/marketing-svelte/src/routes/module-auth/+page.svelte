<script lang="ts">
  import ModuleLanding from '../../lib/components/ModuleLanding.svelte';

  const layers = [
    {
      label: 'Edge gateway',
      detail: 'Auth entry points stay aligned with request-id propagation and the gateway contract surface.',
    },
    {
      label: 'Identity flow',
      detail: 'Sign-in, onboarding, and session handoff are framed as reusable product flows instead of stubs.',
    },
    {
      label: 'Persistence',
      detail: 'Session and profile state are designed to sit on top of the same durable service patterns used elsewhere.',
    },
    {
      label: 'Quality gates',
      detail: 'Release promotion depends on the deployment metadata checks and the launch checklist remaining green.',
    },
  ];

  const callouts = [
    {
      label: 'Launch status',
      detail: 'Integrated into the module catalog and positioned as the identity surface for the product.',
    },
    {
      label: 'Stack fit',
      detail: 'Pairs with the core gateway and the persisted service tier instead of carrying its own ad hoc state.',
    },
    {
      label: 'Product focus',
      detail: 'The page should sell the trust boundary: login, onboarding, and controlled access to the rest of the stack.',
    },
  ];
</script>

<ModuleLanding
  title="Auth Module"
  summary="Identity entry point for onboarding, access control, and the secure handoff into the rest of the product."
  phase="Phase 6"
  domain="Identity"
  status="Integrated"
  launch="Wired to gateway and shared module catalog"
  ctaHref="/"
  ctaLabel="Back to launchpad"
  {layers}
  {callouts}
/>
