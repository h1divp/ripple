<script lang="ts">
  let { isConnected, nearbyCounter } = $props();
  import { IconUsers } from '@tabler/icons-svelte';

</script>

<div class="mb-4 flex flex-row flex-wrap gap-2 justify-start">
    {#if isConnected}
      <div class="rounded-lg flex items-center justify-center border-4 border-solid border-sky-800 bg-sky-100 p-1 text-sky-900">
        <IconUsers size={20} class="mx-1" />
        <span class="pr-2 text-xl font-bold">
          {nearbyCounter}
        </span>
      </div>
    {/if}
      <div class="invisible"></div>
</div>
