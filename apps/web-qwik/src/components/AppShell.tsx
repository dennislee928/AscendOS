import { component$, Slot } from '@builder.io/qwik';

export const AppShell = component$(() => {
  return (
    <div class="app-shell">
      <header class="app-shell__header">
        <div>
          <p class="app-shell__eyebrow">Phase 6 scaffold</p>
          <h1>web-qwik</h1>
        </div>
        <p class="app-shell__subtitle">Route-first Qwik shell for module work.</p>
      </header>
      <main class="app-shell__content">
        <Slot />
      </main>
    </div>
  );
});
