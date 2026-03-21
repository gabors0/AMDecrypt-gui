<script lang="ts">
  import { onMount } from "svelte";
  import Amd from "./routes/Amd.svelte";
  import Console from "./routes/Console.svelte";
  import Setup from "./routes/Setup.svelte";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import { appendLog } from "./lib/logStore.svelte";
  import { addLine, setRunning } from "./lib/amdStore.svelte";

  let currentRoute = $state("setup");

  onMount(() => {
    const unsub1 = EventsOn("log", (msg: string) => appendLog(msg));
    const unsub2 = EventsOn("amd:stdout", (text: string) => addLine(text, "stdout"));
    const unsub3 = EventsOn("amd:stderr", (text: string) => addLine(text, "stderr"));
    const unsub4 = EventsOn("amd:started", () => setRunning(true));
    const unsub5 = EventsOn("amd:stopped", () => setRunning(false));

    return () => {
      unsub1();
      unsub2();
      unsub3();
      unsub4();
      unsub5();
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
    { id: "amd", label: "AMD" },
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
    {:else if currentRoute === "amd"}
      <Amd />
    {:else if currentRoute === "logs"}
      <Console />
    {/if}
  </main>
</div>
