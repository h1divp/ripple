<script lang="ts">
  import { isConnected, nearbyCount } from '$lib/stores/chat';
  import { IconUsers, IconMapPinOff, IconSettings, IconInfoCircle } from '@tabler/icons-svelte';
  // import DistanceButton from '$lib/components/DistanceButton.svelte';
  
  let { locationError }: { locationError: boolean } = $props();

  function handleSettingsClick() {
    console.log('Settings clicked');
  }
  
  function handleInfoClick() {
    console.log('Info clicked');
  }
</script>

<div class="mb-4 lg:px-4 flex flex-row items-center gap-2">
  <div class="flex flex-row gap-1 flex-1 justify-start">
    {#if locationError}
      <div
        class="tooltip flex items-center justify-center rounded-lg border-4 border-solid border-red-800 bg-sky-100 p-1 text-sky-800"
      >
        <IconMapPinOff size={30} class="text-red-800" />
      </div>
    {/if}
    {#if $isConnected && !locationError}
      <div
        class="flex items-center justify-center rounded-lg border-4 border-solid border-sky-800 bg-sky-100 p-1 text-sky-800"
      >
        <IconUsers size={20} class="mx-1" />
        <span class="pr-2 text-xl font-bold select-none">
          {$nearbyCount}
        </span>
      </div>
      <!-- <DistanceButton /> -->
    {/if}
  </div>

  <!-- <div class="flex justify-center">
    <img src="/logo.svg" alt="Echo Chat" class="h-8 w-auto" />
  </div> -->

  <div class="flex flex-row gap-1 flex-1 justify-end">
    <button
      class="flex items-center justify-center rounded-lg bg-sky-100 text-sky-900 hover:[&>svg]:fill-sky-50 active:text-sky-700 cursor-pointer"
      onclick={handleInfoClick}
      title="Information"
    >
      <IconInfoCircle size={40} />
    </button>
    <button
      class="flex items-center justify-center rounded-lg bg-sky-100 text-sky-900 hover:[&>svg]:fill-sky-50 active:text-sky-700 cursor-pointer"
      onclick={handleSettingsClick}
      title="Settings"
    >
      <IconSettings size={40} />
    </button>
  </div>
</div>
