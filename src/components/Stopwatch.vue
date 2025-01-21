<template>
  <div class="text-center font-sans bg-gray-900 px-6 py-4 lg:px-8 lg:py-6 rounded-lg">
    <h2 class="text-2xl lg:text-4xl font-bold mb-4 text-white uppercase">Problem Timer</h2>
    <div class="text-3xl lg:text-5xl my-2 lg:mt-3 text-orange-400/80 font-serif">{{ formattedTime }}</div>
    <div class="flex flex-row gap-x-4 justify-center align-items">
      <button @click="startTimer"
        class="py-2 text-gray-400 hover:text-white transition disabled:text-gray-300 text-xl lg:text-2xl font-serif">
        Start
      </button>
      <button @click="pauseTimer"
        class="py-2 text-gray-400 hover:text-white transition disabled:text-gray-300 text-xl lg:text-2xl font-serif">
        Pause
      </button>
      <button @click="resetTimer" class="py-2 text-gray-400 hover:text-white transition text-xl lg:text-2xl font-serif">
        Reset
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  setup() {
    const time = ref(0);
    const isRunning = ref(false);
    const interval = ref(null);
    // retrieve the time in the following format <hours>:<minutes>:<seconds>
    const formattedTime = computed(() => {
      const hours = Math.floor(time.value / 3600);
      const minutes = Math.floor((time.value % 3600) / 60);
      const seconds = time.value % 60;
      return `${hours.toString().padStart(2, '0')} : ${minutes.toString().padStart(2, '0')} : ${seconds.toString().padStart(2, '0')}`;
    });
    onMounted(() => isRunning.value = false);
    // function for starting timer -- dont do anything if the timer is already running
    const startTimer = function () {
      if (isRunning.value) {
        return;
      }
      isRunning.value = true
      interval.value = setInterval(() => {
        time.value++
      }, 1000)
    }
    // function for pausing the timer -- dont pause unless the timer is running
    const pauseTimer = function () {
      clearInterval(interval.value)
      isRunning.value = false
    }
    // function for resetting timer -- pause the timer as well
    const resetTimer = function () {
      pauseTimer()
      time.value = 0;
    }
    return {
      formattedTime,
      isRunning,
      startTimer,
      pauseTimer,
      resetTimer
    }
  }
}
</script>
