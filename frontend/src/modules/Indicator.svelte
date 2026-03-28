<script lang="ts">
  type Status = "green" | "yellow" | "red" | "off";
  let { status = "off" }: { status?: Status } = $props();

  const colorMap: Record<Status, string> = {
    green:  "oklch(50% 0.1 157)",
    yellow: "oklch(60% 0.1 100)",
    red:    "oklch(50% 0.1 25)",
    off:    "oklch(30% 0 157)",
  };
</script>

<div
  class="w-6 h-6 border border-accent"
  style="background: repeating-linear-gradient(-45deg, transparent, transparent 3px, {colorMap[status]} 3px, {colorMap[status]} 6px);"
></div>
