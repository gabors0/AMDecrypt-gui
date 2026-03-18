<script lang="ts">
    import { RunCmd } from "../../wailsjs/go/app/App.js";
    import { appendLog } from "../lib/logStore.svelte.ts";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";

    function openLink() {
        BrowserOpenURL("https://example.com");
    }

    let isWmStopped = $state(true);
    let isAmdStopped = $state(true);
</script>

<div class="grid max-w-2xl mx-auto grid-cols-2 p-4 gap-4">
    <div class="col-span-2 box flex flex-col">
        <h2 class="p-2 text-xl text-center">Status</h2>
        <hr class="w-full border-accent" />
        <div class="grid grid-cols-2">
            <div class="box p-2 flex justify-between">
                <span>wrapper-manager</span>
                <span class="{isWmStopped ? "text-red-600" : "text-green-500"}">{isWmStopped ? 'Stopped' : 'Running'}</span>
            </div>
            <div class="box p-2 flex justify-between">
                <span>AppleMusicDecrypt</span>
                <span class="{isAmdStopped ? "text-red-600" : "text-green-500"}">{isAmdStopped ? 'Stopped' : 'Running'}</span>

            </div>
        </div>
        <div class="p-2">
            <button class="box w-full py-2 font-bold diagonal-stripes">Check dependencies</button>
        </div>
    </div>
    <div class="box flex flex-col">
        <h2 class="p-2 text-xl text-center">wrapper-manager</h2>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-y-2">
            <div class="w-full flex items-center justify-center">
                <button class="box w-1/2" disabled>Start</button>
                <button class="box w-1/2" disabled>Stop</button>
            </div>
            <button class="box">Install</button>
            <button class="box" disabled>Update</button>
            <button class="box" disabled>Remove</button>
            <button
                class="box"
                onclick={() =>
                    BrowserOpenURL(
                        "https://github.com/WorldObservationLog/wrapper-manager",
                    )}
            >
                Github
            </button>
        </div>
    </div>
    <div class="box flex flex-col">
        <h2 class="p-2 text-xl text-center">AppleMusicDecrypt</h2>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-y-2">
            <div class="w-full flex items-center justify-center">
                <button class="box w-1/2" disabled>Start</button>
                <button class="box w-1/2" disabled>Stop</button>
            </div>
            <button class="box">Install</button>
            <button class="box" disabled>Update</button>
            <button class="box" disabled>Remove</button>
            <button
                class="box"
                onclick={() =>
                    BrowserOpenURL(
                        "https://github.com/WorldObservationLog/AppleMusicDecrypt",
                    )}
            >
                Github
            </button>
        </div>
    </div>
</div>
