<script lang="ts">
    console.log("ok");
    import { RunCmd } from "../wailsjs/go/app/App.js";
    
    let outputElement: HTMLTextAreaElement;
    let output = "";
    let cmd = "";

    async function handleClick() {
      const result = await RunCmd(cmd);
      output += `> ${cmd}\n${result}\n\n`;
      setTimeout(() => {
        outputElement.scrollTop = outputElement.scrollHeight;
      }, 0);
    }
    
</script>
<main class="w-full h-full text-neutral-50 bg-slate-950 p-4">
    <div>
        <button class="bg-slate-600 text-neutral-50 p-2 rounded-md" on:click={handleClick}>Run</button>
        <input
            class="bg-slate-800 rounded-md p-2"
            type="text"
            name="cmd"
            id="cmd"
            bind:value={cmd}
        />
    </div>
    <textarea
        id="output"
        bind:this={outputElement}
        bind:value={output}
        readonly
        class="w-full h-64 mt-4 bg-slate-800 text-neutral-50 p-2 rounded-md font-mono text-sm resize-none"
    ></textarea>
</main>