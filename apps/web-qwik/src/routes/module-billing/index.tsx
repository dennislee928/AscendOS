import { component$ } from '@builder.io/qwik';
import { launchpadRoutes } from '../../../../packages/ui/src';

export default component$(() => {
  const currentRoute = launchpadRoutes.find((route) => route.href === '/module-billing') ?? launchpadRoutes[0];
  const billingPanels = [
    { label: 'Plans', value: 'Starter and growth tiers ready for pricing wiring' },
    { label: 'Invoices', value: 'Account history and payment state placeholders replaced' },
    { label: 'Entitlements', value: 'Module access and usage summary surfaced' },
  ];

  return (
    <section class="module-shell">
      <header class="ui-card ui-card-raised bg-surface rounded-md p-4">
        <p class="module-shell__eyebrow">Revenue module</p>
        <h2>{currentRoute.name} module</h2>
        <p>{currentRoute.summary} It presents account state, plan context, and direct routes back to the gateway.</p>
        <div class="home-shell__actions">
          <a class="ui-button ui-button-primary" href="/">
            Open shared launchpad
          </a>
          <a class="ui-button ui-button-secondary" href="/module-dashboard">
            Go to dashboard
          </a>
        </div>
      </header>

      <div class="home-shell__grid" aria-label="Billing status summary">
        {billingPanels.map((panel) => (
          <article key={panel.label} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">{panel.label}</p>
            <h3>{panel.value}</h3>
          </article>
        ))}
      </div>

      <div class="home-shell__grid" aria-label="Billing module cards">
        {launchpadRoutes.map((route) => (
          <article key={route.href} class="ui-card bg-surface rounded-md p-4">
            <p class="module-shell__eyebrow">
              {route.phase} · {route.domain}
            </p>
            <h3>{route.name}</h3>
            <p>Status: {route.status}</p>
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
            <a href="/module-auth">Auth module</a>
          </li>
          <li>
            <a href="/module-dashboard">Dashboard module</a>
          </li>
        </ul>
      </aside>
    </section>
  );
});
