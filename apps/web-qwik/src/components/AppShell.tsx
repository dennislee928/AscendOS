import { component$, Slot } from '@builder.io/qwik';

export const AppShell = component$(() => {
  return (
    <div class="app-shell">
      <header>
        <h1>web-qwik</h1>
      </header>
      <main>
        <Slot />
      </main>
    </div>
  );
});
