<script lang="ts">
  import { isConnected, nearbyCount } from '$lib/stores/chat';
  import { IconUsers, IconMapPinOff, IconSettings, IconInfoCircle, IconMinus, IconPlus } from '@tabler/icons-svelte';
  
  let { locationError }: { locationError: boolean } = $props();

  let messageDistance = $state(50); // Default 50m
  const distances = [25, 50, 100, 200, 500]; // Available distances in meters
  
  function handleSettingsClick() {
    console.log('Settings clicked');
  }
  
  function handleInfoClick() {
    console.log('Info clicked');
  }

  function decreaseDistance() {
    const currentIndex = distances.indexOf(messageDistance);
    if (currentIndex > 0) {
      messageDistance = distances[currentIndex - 1];
      console.log('Distance decreased to:', messageDistance + 'm');
    }
  }
  
  function increaseDistance() {
    const currentIndex = distances.indexOf(messageDistance);
    if (currentIndex < distances.length - 1) {
      messageDistance = distances[currentIndex + 1];
      console.log('Distance increased to:', messageDistance + 'm');
    }
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
    <div class="flex items-center bg-sky-100 overflow-hidden border-4 border-sky-800 rounded-lg p-1">
      <button
        class="flex items-center justify-center text-sky-900 text-xl font-bold transition-colors hover:text-sky-700 isabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        onclick={decreaseDistance}
        disabled={distances.indexOf(messageDistance) === 0}
        title="Decrease distance"
      >
        <IconMinus size={20} />
      </button>
    
      <div class="text-sky-900 text-xl font-bold w-16 text-center select-none">
        {messageDistance}m
      </div>
    
      <button
        class="flex items-center justify-center text-sky-900 font-bold text-xl transition-colors hover:text-sky-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        onclick={increaseDistance}
        disabled={distances.indexOf(messageDistance) === distances.length - 1}
        title="Increase distance"
      >
        <IconPlus size={20} />
      </button>
    </div>
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
