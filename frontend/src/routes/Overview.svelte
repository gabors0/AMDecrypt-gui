<script lang="ts">
  import { onMount } from "svelte";
  import {
    IsAmdInstalled,
    IsWmInstalled,
    OpenAppDataDir,
    OpenDownloadsDir,
  } from "../../wailsjs/go/app/App.js";
  import {
    SetupAmd,
    RemoveAmd,
    StartAmd,
    StopAmd,
    KillAmd,
    LoginAmd,
    LogoutAmd,
    SetupWm,
    RemoveWm,
    StartWm,
    StopWm,
    KillWm,
  } from "../../wailsjs/go/app/App.js";
  import { GetInstanceConfig } from "../../wailsjs/go/app/App.js";
  import { amd } from "../lib/amdStore.svelte.ts";
  import { wm } from "../lib/wmStore.svelte.ts";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
  import Indicator from "../modules/Indicator.svelte";

  let isWmStopped = $derived(!wm.running);
  let isAmdStopped = $derived(!amd.running);
  let isAmdInstalling = $state(false);
  let isWmInstalling = $state(false);

  let useCustomInstance = $state(true);

  async function loadInstanceConfig() {
    try {
      const cfg = await GetInstanceConfig();
      useCustomInstance = cfg.url !== "127.0.0.1:8080";
    } catch {
      // config.toml missing — AMD not installed yet, keep defaults
    }
  }

  type DepStatus = null | { installed: boolean; version: string };
  const _cached = JSON.parse(localStorage.getItem("depStatus") ?? "null");
  let isAmdInstalled = $state(_cached?.amdInstalled ?? false);
  let isWmInstalled = $state(_cached?.wmInstalled ?? false);
  let isQuickStarting = $state(false);
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
    const installed = await IsAmdInstalled();
    persistAmdInstalled(installed);
    if (installed) await loadInstanceConfig();
    isAmdInstalling = false;
  }

  function persistWmInstalled(value: boolean) {
    isWmInstalled = value;
    const current =
      JSON.parse(localStorage.getItem("depStatus") ?? "null") ?? {};
    localStorage.setItem(
      "depStatus",
      JSON.stringify({ ...current, wmInstalled: value }),
    );
  }

  async function installWm() {
    isWmInstalling = true;
    await SetupWm();
    persistWmInstalled(await IsWmInstalled());
    isWmInstalling = false;
  }

  async function removeWm() {
    isWmInstalling = true;
    await RemoveWm();
    persistWmInstalled(false);
    isWmInstalling = false;
  }

  function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  async function quickStart() {
    isQuickStarting = true;
    try {
      if (!isWmInstalled) {
        await installWm();
      }
      if (!isAmdInstalled) {
        await installAmd();
      }

      if (!wm.running) {
        await StartWm(false);
        await delay(2500);
      } else {
        await delay(500);
      }

      if (!amd.running) {
        await StartAmd();
      }
    } finally {
      isQuickStarting = false;
    }
  }

  onMount(() => {
    if (isAmdInstalled) void loadInstanceConfig();
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
    <div class="flex items-center justify-end">
      <button class="box py-2 px-3" onclick={() => OpenAppDataDir()}
        >App folder</button
      >
      <button class="box py-2 px-3" onclick={() => OpenDownloadsDir()}
        >Downloads folder</button
      >
    </div>
  </div>
  <!-- modules -->
  <!-- wrapper-manager -->
  <div class="box flex flex-col">
    <h2 class="p-2 text-xl flex items-center justify-between">
      wrapper-manager
      <Indicator
        status={!useCustomInstance && (!isWmInstalled || isWmStopped)
          ? "red"
          : "green"}
      />
    </h2>
    <hr class="w-full border-border" />
    <div class="p-2 flex flex-col h-full gap-y-2">
      {#if useCustomInstance}
        <span class="text-sm text-text-muted text-center py-4"
          >Using selfhosted wrapper-manager</span
        >
      {:else}
        <div class="w-full flex items-center justify-center">
          <button
            class="box w-1/3"
            disabled={!isWmInstalled ||
              !isWmStopped ||
              isWmInstalling ||
              isQuickStarting}
            onclick={() => StartWm(false)}>Start</button
          >
          <button
            class="box w-1/3"
            title="start but log into the app's logs"
            disabled={!isWmInstalled ||
              !isWmStopped ||
              isWmInstalling ||
              isQuickStarting}
            onclick={() => StartWm(true)}>Verbose</button
          >
          <button
            class="box w-1/3"
            disabled={!isWmInstalled ||
              isWmStopped ||
              isWmInstalling ||
              isQuickStarting}
            onclick={() => StopWm()}>Stop</button
          >
          <button
            class="box w-1/3"
            disabled={!isWmInstalled ||
              isWmStopped ||
              isWmInstalling ||
              isQuickStarting}
            onclick={() => KillWm()}>Kill</button
          >
        </div>

        <button
          class="box"
          onclick={installWm}
          disabled={isWmInstalling || isWmInstalled || isQuickStarting}
          >{isWmInstalling ? "Working..." : "Install"}</button
        >
        <button
          class="box"
          disabled
          title="not implemented (in the meantime, reinstalling will get the latest version)"
          >Update</button
        >
        <button
          class="box"
          onclick={removeWm}
          disabled={!isWmInstalled || isWmInstalling || isQuickStarting}
          >Remove</button
        >
        <button
          class="inline-flex items-center justify-center gap-0.5 text-text-muted"
          onclick={() =>
            BrowserOpenURL(
              "https://github.com/WorldObservationLog/wrapper-manager",
            )}
        >
          <span>Github</span>
          <span class="relative top-[-1px] text-md" aria-hidden="true">↗</span>
        </button>
      {/if}
    </div>
  </div>
  <!-- AppleMusicDecrypt -->
  <div class="box flex flex-col">
    <h2 class="p-2 text-xl flex items-center justify-between">
      AppleMusicDecrypt
      <Indicator status={!isAmdInstalled || isAmdStopped ? "red" : "green"} />
    </h2>
    <hr class="w-full border-border" />
    <div class="p-2 flex flex-col gap-y-2">
      <div class="w-full flex items-center justify-center">
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled ||
            !isAmdStopped ||
            isAmdInstalling ||
            isQuickStarting}
          onclick={() => StartAmd()}>Start</button
        >
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled ||
            isAmdStopped ||
            isAmdInstalling ||
            isQuickStarting}
          onclick={() => StopAmd()}>Stop</button
        >
        <button
          class="box w-1/3"
          disabled={!isAmdInstalled ||
            isAmdStopped ||
            isAmdInstalling ||
            isQuickStarting}
          onclick={() => KillAmd()}>Kill</button
        >
      </div>
      <button
        class="box"
        onclick={installAmd}
        disabled={isAmdInstalling || isAmdInstalled || isQuickStarting}
        >{isAmdInstalling ? "Working..." : "Install"}</button
      >
      <button
        class="box"
        disabled
        title="not implemented (in the meantime, reinstalling will get the latest version)"
        >Update</button
      >
      <button
        class="box"
        onclick={removeAmd}
        disabled={!isAmdInstalled || isAmdInstalling || isQuickStarting}
        >Remove</button
      >
      <button
        class="inline-flex items-center justify-center gap-0.5 text-text-muted"
        onclick={() =>
          BrowserOpenURL(
            "https://github.com/WorldObservationLog/AppleMusicDecrypt",
          )}
      >
        <span>Github</span>
        <span class="relative top-[-1px] text-md" aria-hidden="true">↗</span>
      </button>
    </div>
  </div>
  <div
    class="box p-2 grid grid-cols-2 col-span-2 gap-2 w-full mx-auto mt-8 relative"
  >
    <svg
      class="absolute -top-12 left-0 h-12 w-full pointer-events-none text-border"
    >
      <line
        x1="25%"
        y1="0"
        x2="33.3333%"
        y2="100%"
        stroke="currentColor"
        stroke-width="2"
      />
      <line
        x1="75%"
        y1="0"
        x2="66.6667%"
        y2="100%"
        stroke="currentColor"
        stroke-width="2"
      />
    </svg>
    <button
      class="box p-2"
      disabled={!isAmdInstalled ||
        !isWmInstalled ||
        isAmdInstalling ||
        isWmInstalling ||
        isQuickStarting}
      onclick={() => LoginAmd()}>Login</button
    >
    <button
      class="box p-2"
      disabled={!isAmdInstalled ||
        !isWmInstalled ||
        isAmdInstalling ||
        isWmInstalling ||
        isQuickStarting}
      onclick={() => LogoutAmd()}>Logout</button
    >
    <button
      class="box p-2 col-span-2 w-full mx-auto"
      disabled={isQuickStarting || isAmdInstalling || isWmInstalling}
      onclick={quickStart}
      >{isQuickStarting ? "Working..." : "Quick Start"}</button
    >
  </div>
  <div class="col-span-2 box p-2 flex items-center justify-between gap-x-2">
    <details class="w-full">
      <summary class="text-xl">Help</summary>
      <section class="mt-1 ml-2 [&_p]:ml-3">
        <h1 class="text-md">Modules</h1>
        <p class="text-sm text-text-muted">
          There are 2 modules needed for this to work, these are wrapper-manager
          which decrypts Apple Music songs, and AppleMusicDecrypt which uses a
          wrapper instance and provides a tui to make downloading possible.
        </p>
        <h1 class="text-md mt-1">Logging in</h1>
        <p class="text-sm text-text-muted">
          To use your Apple Music subscription (which is needed), you need to
          log in. Make sure both modules are installed, wrapper-manager is
          running and press the login button. Your username is your Apple ID,
          most likely the email you use. <br /> If you get a login failed error,
          try these: <br /> - Reinstall both modules <br /> - Try logging out,
          then log back in <br /> - Wait <br /> If you get a python error, make sure
          wrapper-manager is running
        </p>
        <h1 class="text-md mt-1">Quick Start</h1>
        <p class="text-sm text-text-muted">
          The quick start button will install both modules if they aren't
          already installed and start them in the correct order
        </p>
      </section>
    </details>
  </div>
</div>
