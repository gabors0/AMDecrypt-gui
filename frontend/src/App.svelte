<script lang="ts">
  import { onMount } from "svelte";
  import Console from "./routes/Console.svelte";
  import Setup from "./routes/Setup.svelte";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import { appendLog } from "./lib/logStore.svelte";

  let currentRoute = $state("setup");

  onMount(() => {
    return EventsOn("log", (msg: string) => appendLog(msg));
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
</script>

<div class="flex flex-col h-full w-full bg-bg text-text">
  <nav class="flex flex-row bg-bgmuted border-b border-bgmuted">
    <button
      class="px-5 py-2.5 text-sm transition-colors border-b-2 border-transparent"
      class:border-b-text={currentRoute === "setup"}
      class:text-text={currentRoute === "setup"}
      class:text-textmuted={currentRoute !== "setup"}
      class:hover:text-text={currentRoute !== "setup"}
      class:hover:bg-bg={currentRoute !== "setup"}
      onclick={() => (currentRoute = "setup")}
    >
      Setup
    </button>
    <button
      class="px-5 py-2.5 text-sm transition-colors border-b-2 border-transparent"
      class:border-b-text={currentRoute === "console"}
      class:text-text={currentRoute === "console"}
      class:text-textmuted={currentRoute !== "console"}
      class:hover:text-text={currentRoute !== "console"}
      class:hover:bg-bg={currentRoute !== "console"}
      onclick={() => (currentRoute = "console")}
    >
      Console
    </button>
  </nav>
  <main class="flex-1 overflow-auto">
    {#if currentRoute === "console"}
      <Console />
    {:else if currentRoute === "setup"}
      <Setup />
    {/if}
  </main>
</div>
