<script lang="ts">
    import { RunCmd, WhichCmd } from "../../wailsjs/go/app/App.js";
    import { appendLog } from "../lib/logStore.svelte.ts";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";

    let isWmStopped = $state(true);
    let isWmInstalled = $state(false);
    let isAmdStopped = $state(true);
    let isAmdInstalled = $state(false);

    let usePublicInstance = $state(true);

    type DepStatus = null | { installed: boolean; version: string };
    const _cached = JSON.parse(localStorage.getItem("depStatus") ?? "null");
    let dockerStatus: DepStatus = $state(_cached?.docker ?? null);
    let pythonStatus: DepStatus = $state(_cached?.python ?? null);
    let ffmpegStatus: DepStatus = $state(_cached?.ffmpeg ?? null);
    let gpacStatus: DepStatus = $state(_cached?.gpac ?? null);
    let bento4Status: DepStatus = $state(_cached?.bento4 ?? null);
    let lastChecked: string | null = $state(_cached?.lastChecked ?? null);

    let isReady = $derived(
        dockerStatus?.installed &&
            pythonStatus?.installed &&
            ffmpegStatus?.installed &&
            gpacStatus?.installed &&
            bento4Status?.installed &&
            (usePublicInstance || (isWmInstalled && !isWmStopped)) &&
            isAmdInstalled &&
            !isAmdStopped,
    );

    async function checkDependencies() {
        appendLog("[INFO] (Re-)Checking dependencies...");

        const dockerOut = await RunCmd("docker --version");
        dockerStatus = dockerOut.startsWith("Error:")
            ? { installed: false, version: "" }
            : { installed: true, version: dockerOut };
        appendLog(
            dockerStatus.installed
                ? `[INFO] Docker: ${dockerOut}`
                : "[WARN] Docker: not found",
        );

        let pythonOut = await RunCmd("python3 --version");
        if (pythonOut.startsWith("Error:")) {
            pythonOut = await RunCmd("python --version");
        }
        pythonStatus = pythonOut.startsWith("Error:")
            ? { installed: false, version: "" }
            : { installed: true, version: pythonOut };
        appendLog(
            pythonStatus.installed
                ? `[INFO] Python: ${pythonOut}`
                : "[WARN] Python: not found",
        );

        const ffmpegOut = await RunCmd("ffmpeg -version");
        ffmpegStatus = ffmpegOut.startsWith("Error:")
            ? { installed: false, version: "" }
            : { installed: true, version: ffmpegOut.split("\n")[0] };
        appendLog(
            ffmpegStatus.installed
                ? `[INFO] ffmpeg: ${ffmpegOut.split("\n")[0]}`
                : "[WARN] ffmpeg: not found",
        );

        const gpacOut = await RunCmd("gpac -version");
        gpacStatus = gpacOut.startsWith("Error:")
            ? { installed: false, version: "" }
            : { installed: true, version: gpacOut.split("\n")[0] };
        appendLog(
            gpacStatus.installed
                ? `[INFO] gpac: ${gpacOut.split("\n")[0]}`
                : "[WARN] gpac: not found",
        );

        const bento4Path = await WhichCmd("mp4decrypt");
        bento4Status = bento4Path.startsWith("Error:")
            ? { installed: false, version: "" }
            : { installed: true, version: bento4Path };
        appendLog(
            bento4Status.installed
                ? `[INFO] mp4decrypt: ${bento4Path}`
                : "[WARN] mp4decrypt (Bento4): not found",
        );

        lastChecked = new Date().toLocaleString();
        localStorage.setItem(
            "depStatus",
            JSON.stringify({
                docker: dockerStatus,
                python: pythonStatus,
                ffmpeg: ffmpegStatus,
                gpac: gpacStatus,
                bento4: bento4Status,
                lastChecked,
            }),
        );
        appendLog("[INFO] Dependency check complete.");
    }
</script>

<div class="grid max-w-2xl mx-auto grid-cols-2 p-4 gap-4 mt-4">
    <div
        class="col-span-2 box flex flex-col {isReady
            ? 'diagonal-stripes'
            : 'diagonal-stripes-red'}"
    >
        <h2 class="p-2 text-xl text-center">
            Status: <span class={isReady ? "text-green-500" : "text-red-600"}
                >{isReady ? "READY" : "NOT READY"}</span
            >
        </h2>
    </div>
    <div class="box flex flex-col col-span-2">
        <h2
            class="p-2 text-xl text-center {dockerStatus?.installed &&
            pythonStatus?.installed &&
            ffmpegStatus?.installed &&
            gpacStatus?.installed &&
            bento4Status?.installed
                ? 'diagonal-stripes'
                : 'diagonal-stripes-red'}"
        >
            Dependencies
        </h2>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-y-2">
            <div class="flex items-center gap-x-2">
                <button class="box flex-1 py-2" onclick={checkDependencies}
                    >Run check</button
                >
            </div>
            <div class="grid grid-cols-1 text-sm">
                <div class="box p-2 flex items-center justify-between gap-x-2">
                    <div class="flex flex-col">
                        <span class="">Docker</span>
                        <span class="text-xs text-textmuted"
                            >wrapper-manager</span
                        >
                    </div>
                    <div class="flex flex-col items-end min-w-0">
                        {#if dockerStatus === null}
                            <span class="text-textmuted">Not checked</span>
                            <span class="text-xs invisible">_</span>
                        {:else if dockerStatus.installed}
                            <span class="text-green-500">Installed</span>
                            <span
                                class="text-xs text-textmuted truncate"
                                title={dockerStatus.version}
                                >{dockerStatus.version}</span
                            >
                        {:else}
                            <span class="text-red-600">Not found</span>
                            <span class="text-xs invisible">_</span>
                        {/if}
                    </div>
                </div>
                <div class="box p-2 flex items-center justify-between gap-x-2">
                    <div class="flex flex-col">
                        <span class="">Python</span>
                        <span class="text-xs text-textmuted"
                            >AppleMusicDecrypt</span
                        >
                    </div>
                    <div class="flex flex-col items-end min-w-0">
                        {#if pythonStatus === null}
                            <span class="text-textmuted">Not checked</span>
                            <span class="text-xs invisible">_</span>
                        {:else if pythonStatus.installed}
                            <span class="text-green-500">Installed</span>
                            <span
                                class="text-xs text-textmuted truncate"
                                title={pythonStatus.version}
                                >{pythonStatus.version}</span
                            >
                        {:else}
                            <span class="text-red-600">Not found</span>
                            <span class="text-xs invisible">_</span>
                        {/if}
                    </div>
                </div>
                <div class="box p-2 flex items-center justify-between gap-x-2">
                    <div class="flex flex-col">
                        <span class="">ffmpeg</span>
                        <span class="text-xs text-textmuted"
                            >AppleMusicDecrypt</span
                        >
                    </div>
                    <div class="flex flex-col items-end min-w-0">
                        {#if ffmpegStatus === null}
                            <span class="text-textmuted">Not checked</span>
                            <span class="text-xs invisible">_</span>
                        {:else if ffmpegStatus.installed}
                            <span class="text-green-500">Installed</span>
                            <span
                                class="text-xs text-textmuted truncate"
                                title={ffmpegStatus.version}
                                >{ffmpegStatus.version}</span
                            >
                        {:else}
                            <span class="text-red-600">Not found</span>
                            <span class="text-xs invisible">_</span>
                        {/if}
                    </div>
                </div>
                <div class="box p-2 flex items-center justify-between gap-x-2">
                    <div class="flex flex-col">
                        <span class="">gpac (MP4Box)</span>
                        <span class="text-xs text-textmuted"
                            >AppleMusicDecrypt</span
                        >
                    </div>
                    <div class="flex flex-col items-end min-w-0">
                        {#if gpacStatus === null}
                            <span class="text-textmuted">Not checked</span>
                            <span class="text-xs invisible">_</span>
                        {:else if gpacStatus.installed}
                            <span class="text-green-500">Installed</span>
                            <span
                                class="text-xs text-textmuted truncate"
                                title={gpacStatus.version}
                                >{gpacStatus.version}</span
                            >
                        {:else}
                            <span class="text-red-600">Not found</span>
                            <span class="text-xs invisible">_</span>
                        {/if}
                    </div>
                </div>
                <div class="box p-2 flex items-center justify-between gap-x-2">
                    <div class="flex flex-col">
                        <span class="">Bento4 (mp4decrypt)</span>
                        <span class="text-xs text-textmuted"
                            >AppleMusicDecrypt</span
                        >
                    </div>
                    <div class="flex flex-col items-end min-w-0">
                        {#if bento4Status === null}
                            <span class="text-textmuted">Not checked</span>
                            <span class="text-xs invisible">_</span>
                        {:else if bento4Status.installed}
                            <span class="text-green-500">Installed</span>
                            <span
                                class="text-xs text-textmuted truncate"
                                title={bento4Status.version}
                                >{bento4Status.version}</span
                            >
                        {:else}
                            <span class="text-red-600">Not found</span>
                            <span class="text-xs invisible">_</span>
                        {/if}
                    </div>
                </div>
                {#if lastChecked}
                    <span class="text-xs text-center mt-2 text-textmuted"
                        >Last checked: {lastChecked}</span
                    >
                {/if}
            </div>
        </div>
    </div>
    <div class="box flex flex-col">
        <h2
            class="p-2 text-xl text-center {!usePublicInstance && (!isWmInstalled || isWmStopped)
                ? 'diagonal-stripes-red'
                : 'diagonal-stripes'}"
        >
            wrapper-manager
        </h2>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-y-2">
            <label class="flex items-center justify-between cursor-pointer">
                <span>use public instance</span>
                <input type="checkbox" bind:checked={usePublicInstance} class="sr-only peer" />
                <div class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"></div>
            </label>
            {#if usePublicInstance}
                <input type="text" class="box p-2 text-sm w-full" placeholder="wm.wol.moe" />
                <button class="box p-2">Validate</button>
            {:else}
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
            {/if}
        </div>
    </div>
    <div class="box flex flex-col">
        <h2
            class="p-2 text-xl text-center {!isAmdInstalled || isAmdStopped
                ? 'diagonal-stripes-red'
                : 'diagonal-stripes'}"
        >
            AppleMusicDecrypt
        </h2>
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
    <div class="box flex flex-col">
        <h2
            class="p-2 text-xl text-center {bento4Status?.installed
                ? 'diagonal-stripes'
                : 'diagonal-stripes-red'}"
        >
            Bento4
        </h2>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-y-2">
            <button
                class="box"
                onclick={() =>
                    BrowserOpenURL("https://www.bento4.com/downloads/")}
            >
                Download
            </button>
            <button
                class="box"
                onclick={() =>
                    BrowserOpenURL(
                        "https://github.com/axiomatic-systems/Bento4",
                    )}
            >
                Github
            </button>
        </div>
    </div>
</div>
