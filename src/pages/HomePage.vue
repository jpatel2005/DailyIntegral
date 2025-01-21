<template>
  <div class="flex justify-center items-center w-full">
    <div class="flex flex-col items-center w-1/2 md:w-2/5 mx-12 max-sm:w-fit">
      <div class="p-4 w-full text-center">
        <p class="text-white font-medium text-6xl md:text-6xl lg:text-9xl text-left tracking-wide">DailyIntegral</p>
        <p class="text-white font-light text-3xl text-left tracking-wide mt-4 ml-1">
          Improve your integration skills, day-by-day!
        </p>
        <button type="button"
          class="text-white float-left mt-10 focus:outline-none font-medium text-lg rounded-lg text-xl px-4 py-2 text-center bg-blue-600 hover:bg-blue-700"
          @click="router.push('problem')">Go To Daily Problem</button>
        <button type="button"
          class="text-white float-left ml-4 mt-10 focus:outline-none font-medium text-lg rounded-lg text-xl px-4 py-2 text-center bg-transparent border border-blue-400 hover:bg-blue-500"
          @click="open">See Past Problems</button>
      </div>
    </div>
    <div class="max-sm:hidden flex flex-col items-center w-1/2 md:w-2/5 mx-12 mt-12">
      <img :class="{ hidden: gifLoading }" :src="gifSrc" @load="() => gifLoading = true" class="gif" alt="Math Gif">
      <div :class="{ hidden: !gifLoading }" class="img-skeleton"></div>
    </div>
  </div>
  <p
    class="text-white font-bold text-4xl md:text-8xl text-center font-display indent lowercase mt-10 lg:mt-[-.6em] mx-2">
    New
    Problems
    Daily
    @ {{
      problemTime
    }}</p>
  <ModalsContainer />
</template>


<style scoped>
img.gif {
  width: 120em;
  max-width: 150%;
}

.img-skeleton {
  width: 120em;
  max-width: 150%;
  aspect-ratio: 1.77918882451;
}
</style>

<script>
import { getRefreshTime, sleep } from '../time';
import router from '../router';
import { ref, onMounted } from 'vue';
import { ModalsContainer, useModal } from 'vue-final-modal';
import CalanderModal from '../components/CalanderModal.vue';
export default {
  setup() {
    const gifSrc = ref('');;
    const problemTime = ref('');
    const gifLoading = ref(true);
    onMounted(async () => {
      gifLoading.value = true;
      gifSrc.value = '';
      problemTime.value = getRefreshTime();
      // Use timestamp to force refresh on gif
      const timestamp = new Date().getTime();
      gifSrc.value = `https://i.imgur.com/wR1JWd5.gif?v=${timestamp}`;
      // sleep for a bit -- allow the gif to render before showing it
      await sleep(500);
      gifLoading.value = false;
    });
    // Setup past problems modal
    const { open, close } = useModal({
      component: CalanderModal,
      attrs: {
        onConfirm() {
          close()
        },
      },
    });
    return {
      gifSrc,
      router,
      problemTime,
      gifLoading,
      open,
    };
  },
  components: {
    ModalsContainer
  }
};
</script>
