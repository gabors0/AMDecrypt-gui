<script lang="ts">
  type Status = "green" | "yellow" | "red" | "off";

  let { status = "off" }: { status?: Status } = $props();

  type StatusColors = {
    primary: string;
    secondary: string;
  };

  const colorMap: Record<Status, StatusColors> = {
    green: {
      primary: "var(--color-themegreen)",
      secondary: "var(--color-themegreenmuted)",
    },
    yellow: {
      primary: "var(--color-themeyellow)",
      secondary: "var(--color-themeyellowmuted)",
    },
    red: {
      primary: "var(--color-themered)",
      secondary: "var(--color-themeredmuted)",
    },
    off: {
      primary: "var(--color-themegray)",
      secondary: "var(--color-themegraymuted)",
    },
  };

  const colors = $derived(colorMap[status]);
</script>

<div
  class="w-6 h-6 border border-accent"
  style="background: repeating-linear-gradient(-45deg, {colors.secondary}, {colors.secondary} 3px, {colors.primary} 3px, {colors.primary} 6px);"
></div>
