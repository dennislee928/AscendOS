import { component$ } from '@builder.io/qwik';

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
          <a class="ui-button ui-button-primary" href="/module-auth">
            Open auth
          </a>
          <a class="ui-button ui-button-secondary" href="/module-billing">
            Open billing
          </a>
        </div>
      </div>

      <div class="home-shell__grid" aria-label="Module routes">
        <article class="ui-card bg-surface rounded-md p-4">
          <h3>Auth</h3>
          <p>Authentication and onboarding entry points.</p>
          <a href="/module-auth">Go to /module-auth</a>
        </article>
        <article class="ui-card bg-surface rounded-md p-4">
          <h3>Billing</h3>
          <p>Checkout, plan management, and account billing flows.</p>
          <a href="/module-billing">Go to /module-billing</a>
        </article>
        <article class="ui-card bg-surface rounded-md p-4">
          <h3>Dashboard</h3>
          <p>Summary views and operational dashboards.</p>
          <a href="/module-dashboard">Go to /module-dashboard</a>
        </article>
      </div>
    </section>
  );
});
