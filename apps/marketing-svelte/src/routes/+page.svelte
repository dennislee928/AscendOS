<script lang="ts">
  import { featuredLaunchpadRoutes, launchpadRoutes } from '../../../../packages/ui/src';
</script>

<section class="home-shell">
  <div class="home-shell__hero ui-card ui-card-raised bg-surface rounded-md p-4">
    <p class="home-shell__eyebrow">Phase 6 frontend scaffold</p>
    <h2>Module launchpad</h2>
    <p>Route placeholders are ready. Use the cards below to jump into auth, billing, or dashboard work.</p>
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
        <p>{route.summary}</p>
        <a href={route.href}>Go to {route.name}</a>
      </article>
    {/each}
  </div>
</section>
