<script lang="ts">
  import { featuredLaunchpadRoutes, launchpadRoutes } from '../../../../packages/ui/src';

  const launchpadPhase = launchpadRoutes[0]?.phase ?? 'Phase 6';
  const launchpadDomains = [...new Set(launchpadRoutes.map((route) => route.domain))].join(', ');
</script>

<section class="home-shell">
  <div class="home-shell__hero ui-card ui-card-raised bg-surface rounded-md p-4">
    <p class="home-shell__eyebrow">Shared launchpad</p>
    <h2>Module launchpad</h2>
    <p>{launchpadPhase} tracks {launchpadRoutes.length} modules across {launchpadDomains}.</p>
    <div class="home-shell__actions">
      {#each featuredLaunchpadRoutes as route, index}
        <a
          class={index === 0 ? 'ui-button ui-button-primary' : 'ui-button ui-button-secondary'}
          href={route.href}
        >
          {route.actionLabel}
        </a>
      {/each}
    </div>
  </div>

  <div class="home-shell__grid" aria-label="Module routes">
    {#each launchpadRoutes as route}
      <article class="ui-card bg-surface rounded-md p-4">
        <h3>{route.name}</h3>
        <p><strong>Phase:</strong> {route.phase}</p>
        <p><strong>Domain:</strong> {route.domain}</p>
        <p><strong>Status:</strong> {route.status}</p>
        <p>{route.summary}</p>
        <a href={route.href}>Go to {route.name}</a>
      </article>
    {/each}
  </div>
</section>
