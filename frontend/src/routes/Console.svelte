<script lang="ts">
    import { log } from "../lib/logStore.svelte.ts";
    import { clearLog } from "../lib/logStore.svelte.ts";
    let outputElement: HTMLTextAreaElement;
</script>

<div class="w-full h-full text-neutral-50 p-4 flex flex-col gap-y-2">
    <textarea
        id="output"
        bind:this={outputElement}
        value={log.output}
        readonly
        disabled
    class="w-full flex-1 min-h-0 p-2 box font-mono text-sm resize-none cursor-text bg-bgmuted text-text outline outline-1 outline-accent focus:bg-bgactive"
    ></textarea>
    <button class="p-2 box" on:click={clearLog}>Clear</button>
</div>
