<script lang="ts">
  export type ModuleLayer = {
    label: string;
    detail: string;
  };

  export type ModuleCallout = {
    label: string;
    detail: string;
  };

  export let title: string;
  export let summary: string;
  export let phase: string;
  export let domain: string;
  export let status: string;
  export let launch: string;
  export let ctaHref: string;
  export let ctaLabel: string;
  export let layers: ModuleLayer[] = [];
  export let callouts: ModuleCallout[] = [];
</script>

<section class="module-landing">
  <div class="module-landing__hero ui-card ui-card-raised bg-surface rounded-md p-4">
    <p class="module-landing__eyebrow">Marketing module</p>
    <div class="module-landing__title-row">
      <div>
        <h2>{title}</h2>
        <p class="module-landing__summary">{summary}</p>
      </div>
      <a class="ui-button ui-button-primary module-landing__cta" href={ctaHref}>{ctaLabel}</a>
    </div>

    <div class="module-landing__meta" aria-label={`${title} launch metadata`}>
      <span><strong>Phase:</strong> {phase}</span>
      <span><strong>Domain:</strong> {domain}</span>
      <span><strong>Status:</strong> {status}</span>
      <span><strong>Launch:</strong> {launch}</span>
    </div>
  </div>

  <div class="module-landing__grid">
    <article class="ui-card bg-surface rounded-md p-4">
      <h3>Stack layers</h3>
      <ul class="module-landing__list">
        {#each layers as layer}
          <li>
            <strong>{layer.label}</strong>
            <span>{layer.detail}</span>
          </li>
        {/each}
      </ul>
    </article>

    <article class="ui-card bg-surface rounded-md p-4">
      <h3>Operational callouts</h3>
      <ul class="module-landing__list">
        {#each callouts as callout}
          <li>
            <strong>{callout.label}</strong>
            <span>{callout.detail}</span>
          </li>
        {/each}
      </ul>
    </article>
  </div>
</section>

<style>
  .module-landing {
    display: grid;
    gap: 1rem;
  }

  .module-landing__hero,
  .module-landing__grid {
    display: grid;
    gap: 1rem;
  }

  .module-landing__eyebrow {
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.18em;
    font-size: 0.75rem;
    color: var(--token-muted, #6b7280);
  }

  .module-landing__title-row {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    align-items: flex-start;
    justify-content: space-between;
  }

  .module-landing__summary {
    margin: 0.5rem 0 0;
    max-width: 68ch;
    color: var(--token-muted, #4b5563);
  }

  .module-landing__cta {
    white-space: nowrap;
  }

  .module-landing__meta {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem 1.25rem;
    color: var(--token-muted, #4b5563);
    font-size: 0.95rem;
  }

  .module-landing__list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    gap: 0.875rem;
  }

  .module-landing__list li {
    display: grid;
    gap: 0.2rem;
  }

  .module-landing__list strong {
    color: var(--token-foreground, #111827);
  }

  .module-landing__list span {
    color: var(--token-muted, #4b5563);
  }

  @media (min-width: 48rem) {
    .module-landing__grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
</style>
