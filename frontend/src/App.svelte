<script lang="ts">
  import { onMount } from "svelte";
  import Commands from "./routes/Commands.svelte";
  import Console from "./routes/Console.svelte";
  import Setup from "./routes/Setup.svelte";
  import About from "./routes/About.svelte";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import { appendLog } from "./lib/logStore.svelte";
  import { setRunning } from "./lib/amdStore.svelte";

  let currentRoute = $state("setup");

  if (localStorage.getItem("theme") === "light") {
    document.documentElement.setAttribute("data-theme", "light");
  }

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
    { id: "commands", label: "Command Builder" },
    { id: "logs", label: "Logs" },
    { id: "about", label: "About" },
  ];
</script>

<div class="flex flex-col h-full w-full bg-bg text-text">
  <nav class="flex flex-row divide-x divide-accent bg-bg border-b border-accent">
    {#each tabs as tab}
      <button
        class="px-5 flex-1 py-2.5 text-sm"
        class:bg-bgmuted={currentRoute === tab.id}
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
    {:else if currentRoute === "about"}
      <About />
    {/if}
  </main>
</div>
