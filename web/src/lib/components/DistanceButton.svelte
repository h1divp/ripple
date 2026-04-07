<script lang="ts">
import { messageDistance } from '$lib/stores/chat';
import { IconMinus, IconPlus } from '@tabler/icons-svelte';

function decreaseDistance() {
  const currentIndex = distances.indexOf($messageDistance);
  if (currentIndex > 0) {
    messageDistance.set(distances[currentIndex - 1]);
    console.log('Distance decreased to:', distances[currentIndex - 1] + 'm');
  }
}

function increaseDistance() {
  const currentIndex = distances.indexOf($messageDistance);
  if (currentIndex < distances.length - 1) {
    messageDistance.set(distances[currentIndex + 1]);
    console.log('Distance increased to:', distances[currentIndex + 1] + 'm');
  }
}

const distances = [25, 50, 100, 200, 500];
</script>

<div class="flex items-center bg-sky-100 overflow-hidden border-4 border-sky-800 rounded-lg p-1">
  <button
    class="flex items-center justify-center text-sky-900 text-xl font-bold transition-colors hover:text-sky-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
    onclick={decreaseDistance}
    disabled={distances.indexOf($messageDistance) === 0}
    title="Decrease distance"
  >
    <IconMinus size={20} />
  </button>
  <div class="text-sky-900 text-xl font-bold w-16 text-center select-none">
    {$messageDistance}m
  </div>
  <button
    class="flex items-center justify-center text-sky-900 font-bold text-xl transition-colors hover:text-sky-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
    onclick={increaseDistance}
    disabled={distances.indexOf($messageDistance) === distances.length - 1}
    title="Increase distance"
  >
    <IconPlus size={20} />
  </button>
</div>
