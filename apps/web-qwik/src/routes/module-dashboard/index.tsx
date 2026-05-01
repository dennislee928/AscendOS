import { component$ } from '@builder.io/qwik';
import { launchpadRoutes } from '../../../../packages/ui/src';

export default component$(() => {
  const currentRoute = launchpadRoutes.find((route) => route.href === '/module-dashboard') ?? launchpadRoutes[0];
  const supportingRoutes = launchpadRoutes.filter((route) => route.href !== currentRoute.href);

  const statusItems = [
    { label: 'Freshness', value: 'Live operational snapshots' },
    { label: 'Signal', value: 'Ready for provider and telemetry wiring' },
    { label: 'Navigation', value: 'Launchpad, gateway, and module pages linked' },
  ];

  return (
    <section class="module-shell">
      <header class="ui-card ui-card-raised bg-surface rounded-md p-4">
        <p class="module-shell__eyebrow">Operations console</p>
        <h2>{currentRoute.name} module</h2>
        <p>{currentRoute.summary} This page is the entry point for live status, rollups, and module navigation.</p>
        <div class="home-shell__actions">
          <a class="ui-button ui-button-primary" href="/">
            Open shared launchpad
          </a>
          <a class="ui-button ui-button-secondary" href={currentRoute.href}>
            Stay on dashboard
          </a>
        </div>
      </header>

      <div class="home-shell__grid" aria-label="Dashboard status summary">
        {statusItems.map((item) => (
          <article key={item.label} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">{item.label}</p>
            <h3>{item.value}</h3>
          </article>
        ))}
      </div>

      <div class="home-shell__grid" aria-label="Module navigation">
        {launchpadRoutes.map((route) => (
          <article key={route.href} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">
              {route.domain} · {route.phase}
            </p>
            <h3>{route.name}</h3>
            <p>Status: {route.status}</p>
            <p>{route.summary}</p>
            <a href={route.href}>{route.actionLabel}</a>
          </article>
        ))}
      </div>

      <aside class="ui-card ui-card-subtle bg-surface rounded-md p-4">
        <h3>Gateway and adjacent areas</h3>
        <ul>
          <li>
            <a href="/">Gateway launchpad</a>
          </li>
          {supportingRoutes.map((route) => (
            <li key={route.href}>
              <a href={route.href}>{route.name} module</a>
            </li>
          ))}
        </ul>
      </aside>
    </section>
  );
});
