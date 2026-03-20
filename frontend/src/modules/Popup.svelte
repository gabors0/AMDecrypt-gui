<script lang="ts">
  export let text: string;
  export let position: "top" | "bottom" | "left" | "right" = "top";
  export let long: boolean = false;
</script>

<div class="popup-wrapper relative inline-block">
  <slot />
  <div
    class="popup pointer-events-none absolute z-99 px-2 py-1 text-sm {long
      ? 'min-w-64'
      : ''}
               bg-bgmuted border border-accent text-text
               whitespace-pre-line opacity-0 hover:opacity-100 transition-opacity duration-200
               {position === 'top'
      ? 'bottom-full mb-2 left-1/2 -translate-x-1/2'
      : ''}
               {position === 'bottom'
      ? 'top-full mt-2 left-1/2 -translate-x-1/2'
      : ''}
               {position === 'left'
      ? 'right-full mr-2 top-1/2 -translate-y-1/2'
      : ''}
               {position === 'right'
      ? 'left-full ml-2 top-1/2 -translate-y-1/2'
      : ''}"
    role="tooltip"
  >
    {text}
  </div>
</div>

<style>
  .popup-wrapper:hover .popup {
    opacity: 1;
  }
</style>
