<script lang="ts">
  type Status = "green" | "yellow" | "red" | "off";

  let { status = "off" }: { status?: Status } = $props();

  type StatusColors = {
    primary: string;
    secondary: string;
  };

  const colorMap: Record<Status, StatusColors> = {
    green: {
      primary: "var(--color-green)",
      secondary: "var(--color-green-dim)",
    },
    yellow: {
      primary: "var(--color-yellow)",
      secondary: "var(--color-yellow-dim)",
    },
    red: {
      primary: "var(--color-red)",
      secondary: "var(--color-red-dim)",
    },
    off: {
      primary: "var(--color-text-dim)",
      secondary: "var(--color-border-dim)",
    },
  };

  const colors = $derived(colorMap[status]);
</script>

<div
  class="w-6 h-6 border border-border"
  style="background: repeating-linear-gradient(-45deg, {colors.secondary}, {colors.secondary} 3px, {colors.primary} 3px, {colors.primary} 6px);"
></div>
