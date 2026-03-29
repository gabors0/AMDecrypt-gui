<script lang="ts">
    import { ClipboardSetText } from "../../wailsjs/runtime/runtime";
    import Popup from "../modules/Popup.svelte";
    import Indicator from "../modules/Indicator.svelte";

    let selectedMode = $state<"download" | "quality">("download");
    let downloadCodec = $state("default");
    let qualityUrl = $state("");
    let downloadUrls = $state("");
    let downloadOverwrite = $state(false);
    let downloadLang = $state("");
    let downloadIncludeParticipate = $state(false);
    let copyFeedback = $state(false);

    const downloadCommand = $derived.by(() => {
        const parts: string[] = ["dl"];
        if (downloadCodec !== "default") parts.push("-c", downloadCodec);
        if (downloadOverwrite) parts.push("-f");
        if (downloadLang.trim()) parts.push("-l", downloadLang.trim());
        if (downloadIncludeParticipate) parts.push("--include-participate-songs");
        const urls = downloadUrls.trim();
        if (urls) parts.push(...urls.split(/\s+/));
        return parts.join(" ");
    });

    const qualityCommand = $derived.by(() => {
        const parts: string[] = ["qa"];
        const url = qualityUrl.trim();
        if (url) parts.push(url);
        return parts.join(" ");
    });

    const currentCommand = $derived(selectedMode === "download" ? downloadCommand : qualityCommand);

    async function copyCommand() {
        await ClipboardSetText(currentCommand);
        copyFeedback = true;
        setTimeout(() => { copyFeedback = false; }, 1500);
    }

    const codecs = [
        { value: "default", label: "ALAC" },
        { value: "ec3", label: "EC3" },
        { value: "ac3", label: "AC3" },
        { value: "aac", label: "AAC" },
        { value: "aac-binaural", label: "AAC Binaural" },
        { value: "aac-downmix", label: "AAC Downmix" },
        { value: "aac-legacy", label: "AAC Legacy" },
    ];
</script>

<div class="grid max-w-2xl mx-auto p-4 gap-4 mt-4">
    <div class="box flex flex-col w-full">
        <h1 class="p-2 text-xl">Command output</h1>
        <hr class="w-full border-accent" />
        <div class="p-2 flex flex-col gap-2">
            <input type="text" value={currentCommand} readonly class="text-xl box p-2 w-full cursor-text focus:!bg-bgmuted" />
            <button class="w-full box p-2" onclick={copyCommand}>{copyFeedback ? "Copied!" : "Copy to clipboard"}</button>
        </div>
    </div>
    <div class="flex flex-col gap-4">
        <div class="box flex flex-col w-full" onpointerdown={() => selectedMode = "download"}>
            <label
                class="p-2 text-xl cursor-pointer flex items-center justify-between"
               
            >
                <input type="radio" name="commandMode" value="download" bind:group={selectedMode} class="sr-only" />
                Download
                <Indicator status={selectedMode === 'download' ? 'green' : 'off'} />
            </label>
            <hr class="w-full border-accent" />
            <div class="p-2 flex flex-col items-center gap-2">
                <!-- download link -->
                <input type="text" class="box p-2 text-sm w-full" placeholder="Artist/album/playlist/song link(s) // Separate with space" bind:value={downloadUrls} />
                <!-- overwrite (-f) -->
                <label class="flex w-full items-center justify-between cursor-pointer">
                    <span>Overwrite existing files</span>
                    <input type="checkbox" class="sr-only peer" bind:checked={downloadOverwrite} />
                    <div class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"></div>
                </label>
                <!-- codec (-c) -->
                <div class="flex w-full items-center justify-between gap-2">
                    <span>Codec</span>
                    <div class="flex flex-wrap justify-end">
                        {#each codecs as opt}
                            <label class="box px-2 py-1 text-sm cursor-pointer hover:bg-bgactive! {downloadCodec === opt.value ? 'bg-bgactive!' : ''}">
                                <input type="radio" name="downloadCodec" value={opt.value} bind:group={downloadCodec} class="sr-only" />
                                {opt.label}
                            </label>
                        {/each}
                    </div>
                </div>
                <!-- metadata language (-l) -->
                <div class="flex w-full items-center justify-between">
                    <span>Metadata language</span>
                    <input
                      type="text"
                      class="box p-2 text-sm h-5"
                      placeholder="en-US"
                      bind:value={downloadLang}
                    />
                </div>
                <!-- include participate songs (--include-participate-songs) -->
                <label class="flex w-full items-center justify-between cursor-pointer">
                    <span>Include songs as featured artist</span>
                    <input type="checkbox" class="sr-only peer" bind:checked={downloadIncludeParticipate} />
                    <div class="w-5 h-5 box flex items-center justify-center text-text text-sm leading-none peer-checked:after:content-['✕'] after:content-['']"></div>
                </label>
            </div>
        </div>
        <div class="box flex flex-col w-full" onpointerdown={() => selectedMode = "quality"}>
            <label
                class="p-2 text-xl cursor-pointer flex items-center justify-between"
               
            >
                <input type="radio" name="commandMode" value="quality" bind:group={selectedMode} class="sr-only" />
                <span>Quality <span class="underline cursor-help"><Popup long text="Checks which qualities are available for the given link." position="bottom">[?]</Popup></span></span>
                <Indicator status={selectedMode === 'quality' ? 'green' : 'off'} />
            </label>
            <hr class="w-full border-accent" />
            <div class="p-2 flex flex-col items-center gap-2">
                <input type="text" class="box p-2 text-sm w-full" placeholder="Song/album/playlist URL" bind:value={qualityUrl} />
            </div>
        </div>
    </div>
</div>
