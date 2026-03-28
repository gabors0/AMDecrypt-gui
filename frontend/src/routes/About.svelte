<script>
    import appIcon from "../../../build/no_radius.png";
    import Indicator from "../modules/Indicator.svelte";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";

    const links = [
        {
            label: "AMDecrypt-gui",
            url: "https://github.com/gabors0/AMDecrypt-gui",
        },
        {
            label: "AppleMusicDecrypt",
            url: "https://github.com/WorldObservationLog/AppleMusicDecrypt",
        },
        {
            label: "wrapper-manager",
            url: "https://github.com/WorldObservationLog/wrapper-manager",
        },
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
                <span class="text-sm text-textmuted">Version 0.1.0</span>
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
    </div>
</div>
