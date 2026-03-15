<script lang="ts">
  let { isConnected, nearbyCounter } = $props();
  import { IconUsers } from '@tabler/icons-svelte';

</script>

<div class="mb-4 flex flex-row flex-wrap items-center justify-center gap-2 sm:justify-start">
  <div
    class="flex w-fit flex-row gap-3 rounded-lg border-4 border-solid border-sky-800 bg-sky-100 p-1 text-sky-900"
  >
    {#if isConnected}
      <div class="rounded-lgp-2 flex items-center justify-center">
        <IconUsers size={20} class="mx-1" />
        <span class="pr-2 text-xl font-bold">
          {nearbyCounter}
        </span>
      </div>
    {:else}
      <div class="rounded-lgp-2 flex items-center justify-center">
        <span class="px-2 text-xl font-bold">Loading</span>
      </div>
    {/if}
  </div>
</div>
