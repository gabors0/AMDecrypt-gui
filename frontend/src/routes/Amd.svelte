<script lang="ts">
  import { amd, clearLines } from "../lib/amdStore.svelte.ts";
  import { parseAmdLine, levelColor } from "../lib/amdParser.ts";
  import { StartAmd, SendInput } from "../../wailsjs/go/app/App.js";

  let scrollContainer: HTMLDivElement;
  let inputValue = $state("");
  let history: string[] = $state([]);
  let historyIndex = $state(-1);

  let autoScroll = $state(true);

  $effect(() => {
    amd.lines.length;
    if (autoScroll && scrollContainer) {
      requestAnimationFrame(() => {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      });
    }
  });

  function handleScroll() {
    if (!scrollContainer) return;
    const { scrollTop, scrollHeight, clientHeight } = scrollContainer;
    autoScroll = scrollHeight - scrollTop - clientHeight < 40;
  }

  async function send() {
    const text = inputValue.trim();
    if (!text) return;
    inputValue = "";
    history = [...history, text];
    historyIndex = -1;
    await SendInput(text);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      send();
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      if (history.length === 0) return;
      if (historyIndex === -1) {
        historyIndex = history.length - 1;
      } else if (historyIndex > 0) {
        historyIndex--;
      }
      inputValue = history[historyIndex];
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      if (historyIndex === -1) return;
      if (historyIndex < history.length - 1) {
        historyIndex++;
        inputValue = history[historyIndex];
      } else {
        historyIndex = -1;
        inputValue = "";
      }
    }
  }

  async function handleStart() {
    try {
      await StartAmd();
    } catch (err) {
      console.error("Failed to start AMD:", err);
    }
  }
</script>

<div class="w-full h-full flex flex-col p-4 gap-y-2">
  {#if !amd.running && amd.lines.length === 0}
    <div class="flex-1 flex flex-col items-center justify-center gap-y-4">
      <span class="text-textmuted text-lg">AMD is not running</span>
      <button class="box px-6 py-2" onclick={handleStart}>Start</button>
    </div>
  {:else}
    <div
      class="flex-1 min-h-0 overflow-y-auto box p-2 font-mono text-sm bg-bgmuted"
      bind:this={scrollContainer}
      onscroll={handleScroll}
    >
      {#each amd.lines as line}
        {@const parsed = parseAmdLine(line.text, line.source)}
        {#if parsed.type === "table"}
          <pre class="text-text whitespace-pre">{parsed.raw}</pre>
        {:else if parsed.type === "loguru"}
          <div class={levelColor(parsed.level)}>
            <span class="text-textmuted">{parsed.timestamp}</span>
            <span class="font-bold">[{parsed.level}]</span>
            {parsed.message}
          </div>
        {:else}
          <div class={line.source === "stderr" ? "text-red-500" : "text-text"}>
            {parsed.raw}
          </div>
        {/if}
      {/each}
    </div>
    <div class="flex gap-x-2">
      <input
        type="text"
        class="flex-1 box p-2 text-sm font-mono bg-bgmuted"
        placeholder={amd.running ? "Type a command..." : "AMD is not running"}
        bind:value={inputValue}
        onkeydown={handleKeydown}
        disabled={!amd.running}
      />
      <button class="box px-4" onclick={send} disabled={!amd.running}>Send</button>
      <button class="box px-4" onclick={clearLines}>Clear</button>
    </div>
  {/if}
</div>
