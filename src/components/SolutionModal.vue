<template>
  <VueFinalModal class="flex flex-col justify-center items-center w-full h-full"
    content-class="flex flex-col w-fit mx-2 p-4 bg-gray-900 border border-gray-700 rounded-lg space-y-4"
    content-transition="vfm-fade" overlay-transition="vfm-fade">
    <h1 class="text-center text-blue-300 font-serif font-light text-3xl sm:text-5xl my-2 sm:my-4">
      Solution{{ props.steps === '' ? '' : ' With Steps' }}
    </h1>
    <div v-if="validAccount"
      class="flex flex-col sm:flex-row justify-center items-center w-full h-full gap-y-4 sm:gap-x-4">
      <vue-latex :expression="props.problem + '='" :fontsize="screenSize" display-mode
        class="text-center text-white text-xl sm:text-4xl select-none" />
      <vue-latex :expression="props.solution" :fontsize="screenSize" display-mode
        class="text-center text-white text-xl sm:text-4xl select-none" />
    </div>
    <div v-if="validAccount && props.steps !== ''">
      <h3 class="text-center text-blue-300 font-serif font-light text-3xl sm:text-5xl mt-4">
        Steps to Solve
      </h3>
      <h3 class="text-center text-blue-200 font-light text-sm sm:text-md mb-4">(horizontal scroll enabled)</h3>
      <div
        class="my-8 w-fit max-h-[30vh] max-w-[90vw] overflow-x-scroll scrollbar-thin scrollbar-thumb-gray-500 scrollbar-track-gray-700">
        <vue-latex :expression="props.steps" :fontsize="18" display-mode class="text-center text-gray-200" />
      </div>
    </div>
    <h3 v-if="!validAccount" class="text-center text-white font-serif font-light text-4xl mx-2 my-4">
      Please sign-in before viewing solution!
    </h3>
    <div class="flex justify-center items-center mx-6 sm:mx-12 mt-4">
      <button class="text-white px-3 py-2 border rounded-lg w-28 sm:w-40 bg-gray-100/50 hover:bg-gray-100/25"
        @click="emit('confirm')">
        Close
      </button>
    </div>
  </VueFinalModal>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { useAuth0 } from '@auth0/auth0-vue';
import { VueLatex } from 'vatex';
import 'vue-final-modal/style.css';
import 'v-calendar/style.css';

const props = defineProps(['problem', 'solution', 'steps']);
const emit = defineEmits('confirm');
const validAccount = ref(false);
onMounted(async () => {
  const { user } = useAuth0();
  validAccount.value = user.value !== undefined;
});
// Manage resizing of the latex expression
const screenSize = ref('default');
const updateScreenSize = function () {
  const width = window.innerWidth;
  if (width >= 1024) {
    screenSize.value = 72;
  } else if (width >= 768) {
    screenSize.value = 40;
  } else {
    screenSize.value = 18;
  }
  screenSize.value *= (7 / 20);
};
onMounted(() => {
  updateScreenSize();
  window.addEventListener('resize', updateScreenSize);
});
onMounted(() => {
  updateScreenSize();
  window.addEventListener('resize', updateScreenSize);
});
</script>
