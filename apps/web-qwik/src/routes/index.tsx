import { component$ } from '@builder.io/qwik';
import { featuredLaunchpadRoutes, launchpadRoutes } from '../../../../packages/ui/src';

export default component$(() => {
  return (
    <section class="home-shell">
      <div class="home-shell__hero ui-card ui-card-raised bg-surface rounded-md p-4">
        <p class="home-shell__eyebrow">Phase 6 frontend scaffold</p>
        <h2>Module launchpad</h2>
        <p>
          Route placeholders are ready. Use the cards below to jump into auth, billing, or dashboard work.
        </p>
        <div class="home-shell__actions">
          {featuredLaunchpadRoutes.map((route, index) => (
            <a
              key={route.href}
              class={index === 0 ? 'ui-button ui-button-primary' : 'ui-button ui-button-secondary'}
              href={route.href}
            >
              {route.actionLabel}
            </a>
          ))}
        </div>
      </div>

      <div class="home-shell__grid" aria-label="Module routes">
        {launchpadRoutes.map((route) => (
          <article key={route.href} class="ui-card bg-surface rounded-md p-4">
            <h3>{route.name}</h3>
            <p>{route.summary}</p>
            <a href={route.href}>Go to {route.name}</a>
          </article>
        ))}
      </div>
    </section>
  );
});
