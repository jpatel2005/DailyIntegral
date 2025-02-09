<template>
  <VueFinalModal class="flex justify-center items-center"
    content-class="flex flex-col max-w-xl mx-4 p-4 bg-gray-900 border border-gray-700 rounded-lg space-y-2"
    content-transition="vfm-fade" overlay-transition="vfm-fade">
    <h1 class="text-center text-white font-extralight text-2xl mb-2">Select Old Daily Problem</h1>
    <div :class="{ 'mx-auto py-4': isLoading }">
      <span v-if="isLoading" class="pi pi-spin pi-cog text-[8em] text-center text-gray-400"></span>
      <DatePicker v-else transparent color=purple is-dark=true :attributes="solvesData"
        :min-date="convertDateToUTC(new Date(startDate))" :max-date="currDate" v-model.string="problemDate"
        :initial-page="initialPage" :masks="masks" />
      <span></span>
    </div>
    <div class="text-white font-light font-serif flex flex-row gap-x-8 items-center justify-center">
      <div class="flex flex-row gap-x-2 items-center">
        <span>Solved as daily problem</span>
        <span class="text-green-400 text-2xl">&bull;</span>
      </div>
      <div class="flex flex-row gap-x-2 items-center">
        <span>Solved</span>
        <span class="text-gray-400 text-2xl">&bull;</span>
      </div>
    </div>
    <div>
      <div class="mx-3 mt-1">
        <button class="text-white px-3 py-[.1em] border rounded-lg float-left hover:bg-gray-100/25"
          @click="emit('confirm')">
          Close
        </button>
        <button class="text-white px-3 py-[.1em] border rounded-lg float-right bg-purple-600 disabled:bg-gray-400"
          @click="router.push({ path: `/problem/${problemDate}` })" :disabled="problemDate === ''">
          Go To Problem
        </button>
      </div>
    </div>
  </VueFinalModal>
</template>

<script setup>
import router from '../router';
import { onMounted, ref } from 'vue';
import { convertDateToUTC } from '../time';
import { VueFinalModal } from 'vue-final-modal';
import { DatePicker } from 'v-calendar';
import { useAuth0 } from '@auth0/auth0-vue';
import 'vue-final-modal/style.css';
import 'v-calendar/style.css';

const emit = defineEmits('confirm');
const startDate = ref('');
const solvesData = ref([]);
const isLoading = ref(true);
const problemDate = ref('');
const masks = ref({
  modelValue: 'YYYY-MM-DD'
});
const initialPage = ref({});
const getUser = async function (username) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/user/${username}`, {
      method: 'GET'
    });
    // Return null is there was an error (user doesn't exist OR if too many requests were sent)
    if (!response.ok) {
      return null;
    }
    const json = await response.json();
    return json;
  } catch {
    return null;
  }
}
const getStartDate = async function () {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/startdate`, {
      method: 'GET'
    });
    if (!response.ok) {
      return '2025-02-01';
    }
    const data = await response.json();
    return data;
  } catch {
    return '2025-02-01';
  }
}
const currDate = ref(null);
onMounted(async () => {
  isLoading.value = true;
  currDate.value = convertDateToUTC(new Date());
  // Set the starting page so that the current date is visible
  initialPage.value = {
    day: currDate.value.getDate(),
    month: currDate.value.getMonth() + 1,
    year: currDate.value.getFullYear()
  }
  const { user } = useAuth0();
  // Retrieve the date for the first problem
  startDate.value = await getStartDate();
  // Reformat date to a date variable that uses UTC
  startDate.value = convertDateToUTC(new Date(startDate.value));
  if (user.value === undefined) {
    isLoading.value = false;
    return;
  }
  const userData = await getUser(user.value.username);
  if (userData === null) {
    isLoading.value = false;
    return;
  }
  // Process all of the solves and identify each solve as a daily problem solve or a normal problem solve (used for rendering colored dots on the calander)
  const dailyProblems = [];
  const normalProblems = [];
  for (const [date, solveVal] of Object.entries(userData.solves)) {
    if (solveVal === 2) {
      dailyProblems.push(convertDateToUTC(new Date(date)));
    } else {
      normalProblems.push(convertDateToUTC(new Date(date)));
    }
  }
  solvesData.value = [
    {
      dot: 'green',
      dates: dailyProblems
    },
    {
      dot: 'gray',
      dates: normalProblems
    }
  ];
  isLoading.value = false;
});
</script>
