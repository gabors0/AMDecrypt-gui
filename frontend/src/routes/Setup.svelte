<script lang="ts">
  import { onMount } from "svelte";
  import {
    RunCmd,
    WhichCmd,
    SetupBento4,
    RemoveBento4,
    GetInstanceConfig,
    SetInstanceConfig,
    GetOS,
    IsAmdInstalled,
    IsWmInstalled,
    GetSettings,
    SaveSettings,
    DetectTerminal,
  } from "../../wailsjs/go/app/App.js";
  import { appendLog } from "../lib/logStore.svelte.ts";
  import Popup from "../modules/Popup.svelte";
  import Indicator from "../modules/Indicator.svelte";

  let currentOS = $state("");
  let terminalBin = $state("");

  GetOS().then((os) => {
    currentOS = os;
  });
  GetSettings().then((s) => {
    if (s.terminal) terminalBin = s.terminal;
  });

  async function onTerminalChange() {
    await SaveSettings({ terminal: terminalBin });
  }

  let useCustomInstance = $state(true);
  let instanceUrl = $state("wm.wol.moe");
  let useSecure = $state(true);

  async function loadInstanceConfig() {
    try {
      const cfg = await GetInstanceConfig();
      useCustomInstance = cfg.url !== "127.0.0.1:8080";
      if (useCustomInstance && cfg.url) {
        instanceUrl = cfg.url;
      }
      useSecure = cfg.secure;
    } catch {
      // config.toml missing — AMD not installed yet, keep defaults
    }
  }

  async function onToggleInstance() {
    useCustomInstance = !useCustomInstance;
    if (useCustomInstance) {
      await SetInstanceConfig(instanceUrl, useSecure);
    } else {
      await SetInstanceConfig("127.0.0.1:8080", false);
    }
  }

  async function onToggleSecure() {
    useSecure = !useSecure;
    await SetInstanceConfig(instanceUrl, useSecure);
  }

  async function onInstanceUrlChange() {
    if (useCustomInstance) {
      await SetInstanceConfig(instanceUrl, useSecure);
    }
  }

  type DepStatus = null | { installed: boolean; version: string };
  const _cached = JSON.parse(localStorage.getItem("depStatus") ?? "null");
  let isAmdInstalled = $state(_cached?.amdInstalled ?? false);
  let isWmInstalled = $state(_cached?.wmInstalled ?? false);
  let dockerStatus: DepStatus = $state(_cached?.docker ?? null);
  let goStatus: DepStatus = $state(_cached?.go ?? null);
  let pythonStatus: DepStatus = $state(_cached?.python ?? null);
  let ffmpegStatus: DepStatus = $state(_cached?.ffmpeg ?? null);
  let gpacStatus: DepStatus = $state(_cached?.gpac ?? null);
  let bento4Status: DepStatus = $state(_cached?.bento4 ?? null);
  let lastChecked: string | null = $state(_cached?.lastChecked ?? null);

  let isReady = $derived(
    (useCustomInstance || (dockerStatus?.installed && goStatus?.installed)) &&
      pythonStatus?.installed &&
      ffmpegStatus?.installed &&
      gpacStatus?.installed &&
      bento4Status?.installed &&
      (useCustomInstance || isWmInstalled) &&
      isAmdInstalled,
  );

  async function checkStatus() {
    appendLog("[INFO] (Re-)Checking status...");
    const os = currentOS || (await GetOS());
    currentOS = os;

    const dockerOut = await RunCmd("docker --version");
    dockerStatus = dockerOut.startsWith("Error:")
      ? { installed: false, version: "" }
      : { installed: true, version: dockerOut };
    appendLog(
      dockerStatus.installed
        ? `[INFO] Docker: ${dockerOut}`
        : "[WARN] Docker: not found",
    );

    const goOut = await RunCmd("go version");
    goStatus = goOut.startsWith("Error:")
      ? { installed: false, version: "" }
      : { installed: true, version: goOut };
    appendLog(
      goStatus.installed ? `[INFO] Go: ${goOut}` : "[WARN] Go: not found",
    );

    let pythonOut = "";
    if (os === "windows") {
      pythonOut = await RunCmd("py -3 --version");
    } else {
      pythonOut = await RunCmd("python3 --version");
    }
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

    const gpacPath = await WhichCmd("gpac");
    if (gpacPath.startsWith("Error:")) {
      gpacStatus = { installed: false, version: "" };
    } else {
      const gpacOut = await RunCmd("gpac -version");
      const ver =
        gpacOut.split("\n").find((l: string) => l.includes("version")) ??
        gpacPath;
      gpacStatus = { installed: true, version: ver };
    }
    appendLog(
      gpacStatus.installed
        ? `[INFO] gpac: ${gpacStatus.version}`
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

    isAmdInstalled = await IsAmdInstalled();
    appendLog(
      isAmdInstalled
        ? "[INFO] AppleMusicDecrypt: installed"
        : "[WARN] AppleMusicDecrypt: not installed",
    );
    if (isAmdInstalled) await loadInstanceConfig();

    isWmInstalled = await IsWmInstalled();
    appendLog(
      isWmInstalled
        ? "[INFO] wrapper-manager: installed"
        : "[WARN] wrapper-manager: not installed",
    );

    if (currentOS === "linux" && !terminalBin) {
      const detected = await DetectTerminal();
      if (detected) {
        terminalBin = detected;
        await SaveSettings({ terminal: terminalBin });
        appendLog("[INFO] Terminal auto-detected: " + terminalBin);
      } else {
        appendLog(
          "[WARN] No terminal emulator detected. Set one manually in the Terminal field.",
        );
      }
    }

    lastChecked = new Date().toLocaleString();
    localStorage.setItem(
      "depStatus",
      JSON.stringify({
        docker: dockerStatus,
        go: goStatus,
        python: pythonStatus,
        ffmpeg: ffmpegStatus,
        gpac: gpacStatus,
        bento4: bento4Status,
        amdInstalled: isAmdInstalled,
        wmInstalled: isWmInstalled,
        lastChecked,
      }),
    );
    appendLog("[INFO] Status check complete!");
  }

  let isBento4Working = $state(false);

  async function installBento4() {
    isBento4Working = true;
    await SetupBento4();
    const bento4Path = await WhichCmd("mp4decrypt");
    bento4Status = bento4Path.startsWith("Error:")
      ? { installed: false, version: "" }
      : { installed: true, version: bento4Path };
    isBento4Working = false;
  }

  async function removeBento4() {
    isBento4Working = true;
    await RemoveBento4();
    const bento4Path = await WhichCmd("mp4decrypt");
    bento4Status = bento4Path.startsWith("Error:")
      ? { installed: false, version: "" }
      : { installed: true, version: bento4Path };
    isBento4Working = false;
  }

  onMount(() => {
    if (isAmdInstalled) void loadInstanceConfig();
    void checkStatus();
  });
</script>

<div class="grid max-w-2xl mx-auto grid-cols-2 p-4 gap-4 mt-4">
  <div
    class="col-span-2 box p-2 flex items-center justify-between gap-x-2 mb-4"
  >
    <h2 class="text-xl flex items-center gap-x-2 shrink-0">
      <span
        class="font-bold p-1 px-2 text-bg leading-none {isReady
          ? 'bg-green'
          : 'bg-red'}">{isReady ? "Ready" : "Not ready"}</span
      >
    </h2>
    <div class="flex items-center">
      <button class="box py-2 px-3" onclick={checkStatus}>Run check</button>
    </div>
  </div>
  <div class="box flex flex-col col-span-2">
    <!-- Dependencies -->
    <h2 class="p-2 text-xl flex items-center justify-between">
      <span
        >Dependencies <span class="underline cursor-help"
          ><Popup
            long
            text="Dependencies are required for AMDecrypt to work, please install them and make sure they are on your system PATH!"
            position="right">[?]</Popup
          ></span
        ></span
      >
      <Indicator
        status={pythonStatus?.installed &&
        ffmpegStatus?.installed &&
        gpacStatus?.installed &&
        bento4Status?.installed &&
        (useCustomInstance || (dockerStatus?.installed && goStatus?.installed))
          ? "green"
          : "red"}
      />
    </h2>
    <hr class="w-full border-border" />
    <div class="p-2 flex flex-col gap-y-2">
      <div class="grid grid-cols-1 text-sm">
        <div class="box p-2 flex items-center justify-between gap-x-2">
          <span class="">Python</span>
          <div class="flex flex-col items-end min-w-0">
            {#if pythonStatus === null}
              <span class="text-text-muted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if pythonStatus.installed}
              <span
                class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Installed</span
              >
            {:else}
              <span
                class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Not found</span
              >
            {/if}
          </div>
        </div>
        <div class="box p-2 flex items-center justify-between gap-x-2">
          <span class="">ffmpeg</span>
          <div class="flex flex-col items-end min-w-0">
            {#if ffmpegStatus === null}
              <span class="text-text-muted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if ffmpegStatus.installed}
              <span
                class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Installed</span
              >
            {:else}
              <span
                class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Not found</span
              >
            {/if}
          </div>
        </div>
        <div class="box p-2 flex items-center justify-between gap-x-2">
          <span class="">gpac (MP4Box)</span>
          <div class="flex flex-col items-end min-w-0">
            {#if gpacStatus === null}
              <span class="text-text-muted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if gpacStatus.installed}
              <span
                class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Installed</span
              >
            {:else}
              <span
                class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Not found</span
              >
            {/if}
          </div>
        </div>
        <div class="box p-2 flex items-center justify-between gap-x-2">
          <span class="">Bento4</span>
          <div class="flex flex-row ml-auto *:p-1 *:px-2">
            <button
              class="box"
              onclick={installBento4}
              disabled={isBento4Working || bento4Status?.installed}
            >
              {#if isBento4Working}
                Working...
              {:else}
                Install <span class="underline cursor-help"
                  ><Popup
                    long
                    text={currentOS === "windows"
                      ? "Installs Bento4 with winget. Requires Windows Package Manager to be installed and available on PATH."
                      : "Builds Bento4 from source using git + cmake + make install. Requires cmake to be installed and available on PATH."}
                    position="top">[?]</Popup
                  ></span
                >
              {/if}
            </button>
            <button
              class="box {isBento4Working ? 'p-0! outline-0!' : ''}"
              onclick={removeBento4}
              disabled={isBento4Working ||
                !bento4Status?.installed ||
                currentOS === "windows"}
              >{isBento4Working ? "" : "Remove"}</button
            >
          </div>
          <div class="flex flex-col items-end">
            {#if bento4Status === null}
              <span class="text-text-muted">Not checked</span>
              <span class="text-xs invisible">_</span>
            {:else if bento4Status.installed}
              <span
                class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Installed</span
              >
            {:else}
              <span
                class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                >Not found</span
              >
            {/if}
          </div>
        </div>
        {#if !useCustomInstance}
          <div class="box p-2 flex items-center justify-between gap-x-2">
            <span class="">Docker</span>
            <div class="flex flex-col items-end min-w-0">
              {#if dockerStatus === null}
                <span class="text-text-muted">Not checked</span>
                <span class="text-xs invisible">_</span>
              {:else if dockerStatus.installed}
                <span
                  class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                  >Installed</span
                >
              {:else}
                <span
                  class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                  >Not found</span
                >
              {/if}
            </div>
          </div>
          <div class="box p-2 flex items-center justify-between gap-x-2">
            <span class="">Go</span>
            <div class="flex flex-col items-end min-w-0">
              {#if goStatus === null}
                <span class="text-text-muted">Not checked</span>
              {:else if goStatus.installed}
                <span
                  class="bg-green py-1.5 text-bg font-bold text-center w-20 leading-none"
                  >Installed</span
                >
              {:else}
                <span
                  class="bg-red py-1.5 text-bg font-bold text-center w-20 leading-none"
                  >Not found</span
                >
              {/if}
            </div>
          </div>
        {/if}
        {#if lastChecked}
          <span class="text-xs text-center mt-2 text-text-muted"
            >Last checked: {lastChecked}</span
          >
        {/if}
      </div>
    </div>
  </div>
  <div class="box flex flex-col col-span-2">
    <h2 class="p-2 text-xl flex items-center justify-between">
      Other options
      <Indicator
        status={useCustomInstance ||
        (isWmInstalled && dockerStatus?.installed && goStatus?.installed)
          ? "green"
          : "red"}
      />
    </h2>
    <hr class="w-full border-border" />
    <div class="p-2 flex flex-col gap-y-2">
      <div class="box p-2 flex flex-col gap-y-2">
        <label class="flex items-center justify-between gap-x-2 cursor-pointer">
          <span
            >Use selfhosted wrapper-manager <Popup
              text="Check if you're running your own wrapper instance or someone else's"
              position="right"
              long>[?]</Popup
            >
          </span>
          <input
            type="checkbox"
            checked={useCustomInstance}
            onchange={onToggleInstance}
            class="sr-only peer"
          />
          <div
            class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"
          ></div>
        </label>
        {#if useCustomInstance}
          <label
            class="flex items-center justify-between gap-x-2 cursor-pointer"
          >
            <span>Secure (https)</span>
            <input
              type="checkbox"
              checked={useSecure}
              onchange={onToggleSecure}
              class="sr-only peer"
            />
            <div
              class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"
            ></div>
          </label>
          <input
            type="text"
            class="box p-2 text-sm w-full"
            placeholder="wm.wol.moe"
            bind:value={instanceUrl}
            onchange={onInstanceUrlChange}
          />
        {/if}
      </div>
      {#if currentOS === "linux"}
        <div class="box p-2 flex flex-col gap-y-2">
          <div class="flex items-center justify-between gap-x-2">
            <span>Terminal emulator</span>
            <button
              class="box px-3 py-1 text-sm"
              onclick={async () => {
                const detected = await DetectTerminal();
                if (detected) {
                  terminalBin = detected;
                  await SaveSettings({ terminal: terminalBin });
                }
              }}>Detect</button
            >
          </div>
          <input
            type="text"
            class="box p-2 text-sm w-full"
            placeholder="e.g. konsole, kitty, xterm"
            bind:value={terminalBin}
            onchange={onTerminalChange}
          />
        </div>
      {/if}
    </div>
  </div>
</div>
