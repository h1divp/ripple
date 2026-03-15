<script lang="ts">
let { message } = $props();
</script>

<div class="my-4 flex items-center gap-4">
  <div class="h-[1px] flex-1 bg-gray-500"></div>

  <div class="flex flex-col items-center text-center">
    <span class="text-xs font-bold uppercase tracking-widest text-sky-900/80">
      {message}
    </span>
  </div>

  <div class="h-[1px] flex-1 bg-gray-500"></div>
</div>

