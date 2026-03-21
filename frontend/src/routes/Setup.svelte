<script lang="ts">
  import { RunCmd, WhichCmd, GetAppDataDir, OpenAppDataDir } from "../../wailsjs/go/app/App.js";
  import { SetupAmd, RemoveAmd, StartAmd, StopAmd, KillAmd } from "../../wailsjs/go/app/App.js";
  import { appendLog } from "../lib/logStore.svelte.ts";
  import { amd } from "../lib/amdStore.svelte.ts";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
  import Popup from "../modules/Popup.svelte";

  let isWmStopped = $state(true);
  let isWmInstalled = $state(false);
  let isAmdStopped = $derived(!amd.running);
  let isAmdInstalling = $state(false);

  let usePublicInstance = $state(true);

  type DepStatus = null | { installed: boolean; version: string };
  const _cached = JSON.parse(localStorage.getItem("depStatus") ?? "null");
  let isAmdInstalled = $state(_cached?.amdInstalled ?? false);
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

  async function checkStatus() {
    appendLog("[INFO] (Re-)Checking status...");

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

    const appDataDir = await GetAppDataDir();
    const amdOut = await RunCmd(
      `test -f "${appDataDir}/amd/venv/bin/python" && echo ok`,
    );
    isAmdInstalled = amdOut.trim() === "ok";
    appendLog(
      isAmdInstalled
        ? "[INFO] AppleMusicDecrypt: installed"
        : "[WARN] AppleMusicDecrypt: not installed",
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
        amdInstalled: isAmdInstalled,
        lastChecked,
      }),
    );
    appendLog("[INFO] Status check complete!");
  }

  function persistAmdInstalled(value: boolean) {
    isAmdInstalled = value;
    const current =
      JSON.parse(localStorage.getItem("depStatus") ?? "null") ?? {};
    localStorage.setItem(
      "depStatus",
      JSON.stringify({ ...current, amdInstalled: value }),
    );
  }

  async function removeAmd() {
    isAmdInstalling = true;
    await RemoveAmd();
    persistAmdInstalled(false);
    isAmdInstalling = false;
  }

  async function installAmd() {
    isAmdInstalling = true;
    await SetupAmd();
    const appDataDir = await GetAppDataDir();
    const result = await RunCmd(
      `test -f "${appDataDir}/amd/venv/bin/python" && echo ok`,
    );
    persistAmdInstalled(result.trim() === "ok");
    isAmdInstalling = false;
  }
</script>

<div class="grid max-w-2xl mx-auto grid-cols-2 p-4 gap-4 mt-4">
  <h2
    class="col-span-2 box flex flex-col {isReady
      ? 'diagonal-stripes'
      : 'diagonal-stripes-red'} p-2 text-xl text-center"
  >
    Status: <span class={isReady ? "text-green-500" : "text-red-600"}
      >{isReady ? "READY" : "NOT READY"}</span
    >
  </h2>
  <div class="flex items-center col-span-2 gap-x-2">
    <button class="box flex-1 py-2" onclick={checkStatus}>Run check</button>
    <button class="box py-2 px-3" onclick={() => OpenAppDataDir()}>Open folder</button>
    <Popup
      text="Dependencies are required for AMDecrypt to work, please install them and make sure they are on your system PATH!"
      position="left"
      long={true}><button class="box w-10 h-10 cursor-help">?</button></Popup
    >
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
      <div class="grid grid-cols-1 text-sm">
        <div class="box p-2 flex items-center justify-between gap-x-2">
          <div class="flex flex-col">
            <span class="">Docker</span>
            <span class="text-xs text-textmuted">wrapper-manager</span>
          </div>
          <div class="flex flex-col items-end min-w-0">
            {#if dockerStatus === null}
              <span class="text-textmuted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if dockerStatus.installed}
              <span class="text-green-500">Installed</span>
              <span
                class="text-xs text-textmuted truncate"
                title={dockerStatus.version}>{dockerStatus.version}</span
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
            <span class="text-xs text-textmuted">AppleMusicDecrypt</span>
          </div>
          <div class="flex flex-col items-end min-w-0">
            {#if pythonStatus === null}
              <span class="text-textmuted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if pythonStatus.installed}
              <span class="text-green-500">Installed</span>
              <span
                class="text-xs text-textmuted truncate"
                title={pythonStatus.version}>{pythonStatus.version}</span
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
            <span class="text-xs text-textmuted">AppleMusicDecrypt</span>
          </div>
          <div class="flex flex-col items-end min-w-0">
            {#if ffmpegStatus === null}
              <span class="text-textmuted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if ffmpegStatus.installed}
              <span class="text-green-500">Installed</span>
              <span
                class="text-xs text-textmuted truncate"
                title={ffmpegStatus.version}>{ffmpegStatus.version}</span
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
            <span class="text-xs text-textmuted">AppleMusicDecrypt</span>
          </div>
          <div class="flex flex-col items-end min-w-0">
            {#if gpacStatus === null}
              <span class="text-textmuted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if gpacStatus.installed}
              <span class="text-green-500">Installed</span>
              <span
                class="text-xs text-textmuted truncate"
                title={gpacStatus.version}>{gpacStatus.version}</span
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
            <span class="text-xs text-textmuted">AppleMusicDecrypt</span>
          </div>
          <div class="flex flex-col items-end min-w-0">
            {#if bento4Status === null}
              <span class="text-textmuted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if bento4Status.installed}
              <span class="text-green-500">Installed</span>
              <span
                class="text-xs text-textmuted truncate"
                title={bento4Status.version}>{bento4Status.version}</span
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
      class="p-2 text-xl text-center {!usePublicInstance &&
      (!isWmInstalled || isWmStopped)
        ? 'diagonal-stripes-red'
        : 'diagonal-stripes'}"
    >
      wrapper-manager
    </h2>
    <hr class="w-full border-accent" />
    <div class="p-2 flex flex-col h-full gap-y-2">
      <label class="flex items-center justify-between cursor-pointer">
        <span>use public instance</span>
        <input
          type="checkbox"
          bind:checked={usePublicInstance}
          class="sr-only peer"
        />
        <div
          class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"
        ></div>
      </label>
      {#if usePublicInstance}
        <div class="flex flex-col justify-between h-full">
          <input
            type="text"
            class="box p-2 text-sm w-full"
            placeholder="wm.wol.moe"
          />
          <button class="box p-2" disabled>Validate</button>
        </div>
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
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled || !isAmdStopped || isAmdInstalling}
          onclick={() => StartAmd()}
          >Start</button
        >
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled || isAmdStopped || isAmdInstalling}
          onclick={() => StopAmd()}
          >Stop</button
        >
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled || isAmdStopped || isAmdInstalling}
          onclick={() => KillAmd()}
          >Kill</button
        >
      </div>
      <button
        class="box"
        onclick={installAmd}
        disabled={isAmdInstalling || isAmdInstalled}
        >{isAmdInstalling ? "Installing..." : "Install"}</button
      >
      <button class="box" disabled={!isAmdInstalled || isAmdInstalling}
        >Update</button
      >
      <button
        class="box"
        onclick={removeAmd}
        disabled={!isAmdInstalled || isAmdInstalling}>Remove</button
      >
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
  <!-- <div class="box flex flex-col">
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
    </div> -->
</div>
