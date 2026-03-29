<script lang="ts">
  type Status = "green" | "yellow" | "red" | "off";

  let { status = "off" }: { status?: Status } = $props();

  type StatusColors = {
    primary: string;
    secondary: string;
  };

  const colorMap: Record<Status, StatusColors> = {
    green: {
      primary: "oklch(60% 0.2 175)",
      secondary: "oklch(15% 0.15 175)",
    },
    yellow: {
      primary: "oklch(60% 0.1 100)",
      secondary: "oklch(15% 0.08 100)",
    },
    red: {
      primary: "oklch(50% 0.15 25)",
      secondary: "oklch(15% 0.12 25)",
    },
    off: {
      primary: "oklch(30% 0 157)",
      secondary: "oklch(15% 0 157)",
    },
  };

  const colors = $derived(colorMap[status]);
</script>

<div
  class="w-6 h-6 border border-accent"
  style="background: repeating-linear-gradient(-45deg, {colors.secondary}, {colors.secondary} 3px, {colors.primary} 3px, {colors.primary} 6px);"
></div>
