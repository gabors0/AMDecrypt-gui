<script lang="ts">
    import { RunCmd } from "../../wailsjs/go/app/App.js";
    import { appendLog } from "../lib/logStore.svelte.ts";

    async function runTest() {
        const result = await RunCmd("echo test");
        appendLog(result);
    }
</script>

<button onclick={runTest}> Test </button>
<div class="flex flex-row">
    <div class="w-full"></div>
    <div class="w-full"></div>
</div>
