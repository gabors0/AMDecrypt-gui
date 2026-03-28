<script lang="ts">
  import { onMount } from "svelte";
  import Commands from "./routes/Commands.svelte";
  import Console from "./routes/Console.svelte";
  import Setup from "./routes/Setup.svelte";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import { appendLog } from "./lib/logStore.svelte";
  import { setRunning } from "./lib/amdStore.svelte";

  let currentRoute = $state("setup");

  onMount(() => {
    const unsub1 = EventsOn("log", (msg: string) => appendLog(msg));
    const unsub2 = EventsOn("amd:started", () => setRunning(true));
    const unsub3 = EventsOn("amd:stopped", () => setRunning(false));

    return () => {
      unsub1();
      unsub2();
      unsub3();
    };
  });

  // Disable zoom (Ctrl+scroll and Ctrl+/-/0)
  window.addEventListener(
    "wheel",
    (e) => {
      if (e.ctrlKey) e.preventDefault();
    },
    { passive: false },
  );
  window.addEventListener("keydown", (e) => {
    if (e.ctrlKey && ["=", "-", "0", "+"].includes(e.key)) e.preventDefault();
  });

  const tabs = [
    { id: "setup", label: "Setup" },
    { id: "commands", label: "Command Creator" },
    { id: "logs", label: "Logs" },
  ];
</script>

<div class="flex flex-col h-full w-full bg-bg text-text">
  <nav class="flex flex-row bg-bgmuted border-b border-bgmuted">
    {#each tabs as tab}
      <button
        class="px-5 py-2.5 text-sm transition-colors border-b-2 border-transparent"
        class:border-b-text={currentRoute === tab.id}
        class:text-text={currentRoute === tab.id}
        class:text-textmuted={currentRoute !== tab.id}
        class:hover:text-text={currentRoute !== tab.id}
        class:hover:bg-bg={currentRoute !== tab.id}
        onclick={() => (currentRoute = tab.id)}
      >
        {tab.label}
      </button>
    {/each}
  </nav>
  <main class="flex-1 overflow-auto">
    {#if currentRoute === "setup"}
      <Setup />
    {:else if currentRoute === "commands"}
      <Commands />
    {:else if currentRoute === "logs"}
      <Console />
    {/if}
  </main>
</div>
