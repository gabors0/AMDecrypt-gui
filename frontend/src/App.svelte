<script lang="ts">
  import { onMount } from "svelte";
  import Commands from "./routes/Commands.svelte";
  import Console from "./routes/Console.svelte";
  import Overview from "./routes/Overview.svelte";
  import Setup from "./routes/Setup.svelte";
  import About from "./routes/About.svelte";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import { appendLog, log } from "./lib/logStore.svelte";
  import { setRunning } from "./lib/amdStore.svelte";
  import { setWmRunning } from "./lib/wmStore.svelte";

  let currentRoute = $state("overview");

  if (localStorage.getItem("theme") === "light") {
    document.documentElement.setAttribute("data-theme", "light");
  }

  onMount(() => {
    const unsub1 = EventsOn("log", (msg: string) => appendLog(msg));
    const unsub2 = EventsOn("amd:started", () => setRunning(true));
    const unsub3 = EventsOn("amd:stopped", () => setRunning(false));
    const unsub4 = EventsOn("wm:started", () => setWmRunning(true));
    const unsub5 = EventsOn("wm:stopped", () => setWmRunning(false));

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
    { id: "overview", label: "Overview" },
    { id: "setup", label: "Setup" },
    { id: "commands", label: "Commands" },
    { id: "logs", label: "Logs" },
    { id: "about", label: "About" },
  ];
</script>

<div class="flex flex-col h-full w-full bg-bg text-text">
  <nav class="flex min-w-0 flex-row divide-x divide-border bg-bg">
    {#each tabs as tab}
      <button
        class="min-w-0 flex-1 truncate px-2 py-2.5 text-sm border-b border-border"
        class:border-b-bg={currentRoute === tab.id}
        class:hover:border-b-bg-active={currentRoute === tab.id}
        class:text-text={currentRoute === tab.id}
        class:text-text-muted={currentRoute !== tab.id}
        class:hover:text-text={currentRoute !== tab.id}
        class:hover:bg-bg={currentRoute !== tab.id}
        class:error={tab.id === "logs" && log.hasError}
        onclick={() => {
          currentRoute = tab.id;
          if (tab.id === "logs") log.hasError = false;
        }}
      >
        {tab.label}
      </button>
    {/each}
  </nav>
  <main class="flex-1 overflow-auto">
    {#if currentRoute === "overview"}
      <Overview />
    {:else if currentRoute === "setup"}
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
