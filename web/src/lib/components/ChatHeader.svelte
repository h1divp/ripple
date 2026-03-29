<script lang="ts">
  let { nearbyCounter, locationError } = $props();
  import { isConnected } from '$lib/stores/chat'
  import { IconUsers, IconMapPinOff } from '@tabler/icons-svelte';
</script>

<div class="mb-4 flex flex-row flex-wrap justify-start gap-2">
  {#if isConnected && !locationError}
    <div
      class="flex items-center justify-center rounded-lg border-4 border-solid border-sky-800 bg-sky-100 p-1 text-sky-900"
    >
      <IconUsers size={20} class="mx-1" />
      <span class="pr-2 text-xl font-bold">
        {nearbyCounter}
      </span>
    </div>
  {/if}
  {#if locationError}
    <div
      class="tooltip flex items-center justify-center rounded-lg border-4 border-solid border-red-800 bg-sky-100 p-1 text-sky-900"
    >
      <IconMapPinOff size={30} class="text-red-800" />
    </div>
  {/if}
</div>
