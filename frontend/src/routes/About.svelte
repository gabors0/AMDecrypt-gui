<script>
    import appIcon from "../../../build/no_radius.png";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
      
    let selectedTheme = $state(localStorage.getItem("theme") ?? "dark");

    $effect(() => {
        localStorage.setItem("theme", selectedTheme);
        if (selectedTheme === "light") {
            document.documentElement.setAttribute("data-theme", "light");
        } else {
            document.documentElement.removeAttribute("data-theme");
        }
    });

    const themes = [
        { value: "light", label: "Light" },
        { value: "dark", label: "Dark" },
    ];
</script>

<!-- svelte-ignore a11y_invalid_attribute -->
<div class="grid max-w-2xl mx-auto p-4 gap-4 mt-4">
    <div class="box flex flex-col w-full">
        <div class="p-2 text-xl">About</div>
        <hr class="w-full border-accent" />

        <!-- App header -->
        <div class="grid grid-cols-[1fr_auto]">
            <div class="flex flex-col gap-1 p-2">
                <h2 class="text-lg font-semibold">AMDecrypt-gui</h2>
                <span class="text-sm text-textmuted">Version 0.1.1</span>
                <div class="mt-auto flex justify-between *:text-sm *:text-textmuted">
                    
                    <span>made with &lt;3 by <a onclick={() => BrowserOpenURL("https://gs0.me/")} href="#">@gabors0</a> under the <b>MIT License</b></span>
                    <a onclick={() => BrowserOpenURL("https://github.com/gabors0/AMDecrypt-gui")} href="#">github</a>
                </div>
            </div>
            <div class="border-l border-accent aspect-square flex items-center justify-center">
                <img src={appIcon} alt="logo" class="w-32" />
            </div>
        </div>
        <hr class="w-full border-accent">
        <div class="p-2 text-center text-sm text-textmuted">
            This project uses (but does not modify, bundle, or embed)
            <a onclick={() => BrowserOpenURL("https://github.com/WorldObservationLog/AppleMusicDecrypt")} href="#">AppleMusicDecrypt</a>
            &amp;
            <a onclick={() => BrowserOpenURL("https://github.com/WorldObservationLog/wrapper-manager")} href="#">wrapper-manager</a>
            by <a onclick={() => BrowserOpenURL("https://github.com/WorldObservationLog")} href="#">@WorldObservationLog</a>
            under the <b>AGPL-3.0 License</b>
        </div>
        <hr class="w-full border-accent">
        <div class="p-2 text-center text-sm text-textmuted">
            <div class="flex w-full items-center justify-between gap-2">
                <span>Theme</span>
                <div class="flex flex-wrap justify-end">
                {#each themes as opt}
                        <label class="box border-r-0 px-2 py-1 text-sm cursor-pointer hover:bg-bgactive! {selectedTheme === opt.value ? 'bg-bgactive!' : ''}">
                            <input type="radio" name="downloadCodec" value={opt.value} bind:group={selectedTheme} class="sr-only" />
                            {opt.label}
                        </label>
                {/each}
                </div>
            </div>
        </div>
    </div>
</div>
