import { component$ } from '@builder.io/qwik';
import { launchpadRoutes } from '../../../../packages/ui/src';

export default component$(() => {
  const currentRoute = launchpadRoutes.find((route) => route.href === '/module-auth') ?? launchpadRoutes[0];
  const metrics = [
    { label: 'Entry points', value: 'Login, sign-up, session recovery' },
    { label: 'Policy', value: 'Auth flows ready for provider wiring' },
    { label: 'Output', value: 'Preserved redirects and module navigation' },
  ];

  return (
    <section class="module-shell">
      <header class="ui-card ui-card-raised bg-surface rounded-md p-4">
        <p class="module-shell__eyebrow">Identity module</p>
        <h2>{currentRoute.name} module</h2>
        <p>{currentRoute.summary} It keeps the sign-in surface small, direct, and ready for provider integration.</p>
        <div class="home-shell__actions">
          <a class="ui-button ui-button-primary" href="/">
            Open shared launchpad
          </a>
          <a class="ui-button ui-button-secondary" href="/module-billing">
            Go to billing
          </a>
        </div>
      </header>

      <div class="home-shell__grid" aria-label="Auth status summary">
        {metrics.map((metric) => (
          <article key={metric.label} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">{metric.label}</p>
            <h3>{metric.value}</h3>
          </article>
        ))}
      </div>

      <div class="home-shell__grid" aria-label="Auth module cards">
        {launchpadRoutes.map((route) => (
          <article key={route.href} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">
              {route.domain} · {route.status}
            </p>
            <h3>{route.name}</h3>
            <p>{route.summary}</p>
            <a href={route.href}>{route.actionLabel}</a>
          </article>
        ))}
      </div>

      <aside class="ui-card ui-card-subtle bg-surface rounded-md p-4">
        <h3>Useful navigation</h3>
        <ul>
          <li>
            <a href="/">Gateway launchpad</a>
          </li>
          <li>
            <a href="/module-dashboard">Dashboard module</a>
          </li>
          <li>
            <a href="/module-billing">Billing module</a>
          </li>
        </ul>
      </aside>
    </section>
  );
});
