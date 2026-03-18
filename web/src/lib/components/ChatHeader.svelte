<script lang="ts">
  let { isConnected, nearbyCounter } = $props();
  import { IconUsers } from '@tabler/icons-svelte';
</script>

<div class="mb-4 flex flex-row flex-wrap justify-start gap-2">
  {#if isConnected}
    <div
      class="flex items-center justify-center rounded-lg border-4 border-solid border-sky-800 bg-sky-100 p-1 text-sky-900"
    >
      <IconUsers size={20} class="mx-1" />
      <span class="pr-2 text-xl font-bold">
        {nearbyCounter}
      </span>
    </div>
  {/if}
</div>
