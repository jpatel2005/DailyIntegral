<template>
  <div v-if="!isLoading" class="h-full">
    <div v-if="user" class="w-full inline-block">
      <!-- Search bar with button (normal version) -->
      <div class="max-w-md float-right w-80 mr-10 mt-4">
        <label for="default-search" class="mb-2 text-sm font-medium text-gray-900 sr-only">Search</label>
        <div class="relative">
          <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
            <svg class="w-4 h-4 text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
              viewBox="0 0 20 20">
              <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z" />
            </svg>
          </div>
          <input type="search" id="default-search" @input='event => searchUsername = event.target.value'
            class="block w-full p-4 ps-10 text-sm rounded-lg bg-gray-700 border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500 focus-visible:outline-none"
            placeholder="Enter username" required autocomplete="off" />
          <button @click="searchProfile"
            class="text-white absolute end-2.5 bottom-2.5 focus:ring-4 focus:outline-none font-medium rounded-lg text-sm px-4 py-2 bg-blue-600 hover:bg-blue-700 focus:ring-blue-800">Search</button>
        </div>
      </div>
    </div>
    <h3 v-if="user"
      class="text-white text-center font-semibold text-3xl md:text-4xl lg:text-5xl tracking-wide mt-4 lg:mt-[-.25em]">
      {{
        username }}'s
      Profile</h3>
    <div class="flex flex-col gap-y-24 justify-center align-items mt-[10em]" v-else>
      <h3 class="text-white text-center font-semibold text-5xl md:text-7xl tracking-wide mt-[-.25em]">No user data
        found!</h3>
      <h4 class="text-white text-center text-3xl tracking-wide">{{ username === undefined ? `Please create an account
        first!` : "Invalid user" }}</h4>
      <!-- Search bar with button (no user version) -->
      <div class="flex">
        <div class="flex mx-auto items-center justify-between w-[80%] md:w-[30%]">
          <label for="default-search" class="mb-2 text-sm font-medium text-gray-900 sr-only">Search</label>
          <div class="relative w-full">
            <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
              <svg class="w-4 h-4 text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 20 20">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z" />
              </svg>
            </div>
            <input type="search" id="default-search" @input='event => searchUsername = event.target.value'
              class="block w-full p-4 ps-10 text-sm border rounded-lg bg-gray-700 border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500 focus-visible:outline-none"
              placeholder="Enter username" required autocomplete="off" />
            <button @click="searchProfile"
              class="text-white absolute end-2.5 bottom-2.5 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 bg-blue-600 hover:bg-blue-700 focus:ring-blue-800">Search</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="user" class="flex flex-col justify-center items-center">
      <div class="flex flex-col justify-center items-center">
        <div
          class="flex flex-wrap justify-center items-center w-[65%] mx-auto py-6 my-12 border border-blue-200 border-4 rounded-md text-white">
          <div class="flex flex-col gap-y-4 w-1/2 justify-center items-center py-4">
            <span class="font-semibold text-2xl md:text-4xl lg:text-5xl text-center text-blue-300">Daily Problems
              Solved</span>
            <span class="text-left font-light text-5xl lg-text-6xl font-serif">{{ user.daily_problems }}</span>
          </div>
          <div class="flex flex-col gap-y-4 w-1/2 justify-center items-center py-4">
            <span class="font-semibold text-2xl md:text-4xl lg:text-5xl text-center text-blue-300">Total Problems
              Solved</span>
            <span class="text-left font-light text-5xl lg-text-6xl font-serif font-serif">{{ user.total_problems
              }}</span>
          </div>
          <div class="flex flex-col gap-y-4 w-1/2 justify-center items-center py-4">
            <span class="font-semibold text-2xl md:text-4xl lg:text-5xl text-center text-blue-300">Longest Streak</span>
            <span class="text-left font-light text-6xl font-serif">{{ user.longestStreak }}</span>
            <span v-if="user.longestStreak !== 0" class="text-xs md:text-base">{{ user.longestStreakRange }}</span>
          </div>
          <div class="flex flex-col gap-y-4 w-1/2 justify-center items-center py-4">
            <span class="font-semibold text-2xl md:text-4xl lg:text-5xl text-center text-blue-300">Current Streak</span>
            <span class="text-left font-light text-6xl font-serif">{{ user.currentStreak }}</span>
            <span v-if="user.currentStreak !== 0" class="text-xs md:text-base">{{ user.currentStreakRange }}</span>
          </div>
        </div>
        <div class="!w-full !text-white">
          <calendar-heatmap :round="2" :tooltip="true"
            :tooltip-formatter="(v) => v.date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }) + ' - ' + (v.count === 2 ? '<b>Solved as Daily</b>' : '<b>Solved</b>')"
            tooltip-unit="solve" :range-color="['#1f1f22', '#1f1f22', '#d1d5db', '#4ade80']" :max=3 :values="user.data"
            :end-date="endDate" />
        </div>
        <div class="text-white font-light font-serif flex flex-row gap-x-20 items-center justify-center mt-10">
          <div class="flex flex-row gap-x-8 items-center">
            <span class="text-xl md:text-4xl lg:text-6xl">Solved as daily problem</span>
            <span class="inline-block w-4 h-4 bg-green-400 p-4 md:p-6 lg:p-8 rounded-lg"></span>
          </div>
          <div class="flex flex-row gap-x-4 lg:gap-x-8 items-center justify-center">
            <span class="text-xl md:text-4xl lg:text-6xl">Solved</span>
            <span class="inline-block w-4 h-4 bg-gray-300 p-4 md:p-6 lg:p-8 rounded-lg"></span>
          </div>
        </div>
      </div>
      <!-- absolute right-20 bottom-28 -->
      <div class="hidden md:flex flex-row gap-x-4 ml-auto mr-8 items-baseline justify-center">
        <span class="text-center text-gray-400 text-md lg:text-2xl font-semibold">Account Creation Date:</span>
        <span class=" text-center text-gray-200 text-md lg:text-2xl font-mono">{{ user.creation_date }}</span>
      </div>
    </div>
  </div>
  <LoadingCogs v-else />
</template>

<style>
.vch__legend {
  display: none !important;
}

.vch__day__label {
  margin-right: 40px;
}

.vch__month__label,
.vch__day__label {
  fill: white !important;
}
</style>

<script>
import { onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { CalendarHeatmap } from 'vue3-calendar-heatmap';
import 'vue3-calendar-heatmap/dist/style.css';
import { useAuth0 } from '@auth0/auth0-vue';
import { convertDateToUTC, formatDateISO, observeDateRange, reformatDate } from '../time';
import router from '../router';
import LoadingCogs from '../components/LoadingCogs.vue';

const daySeconds = 86400000;

const getUser = async function (username) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/user/${username}`, {
      method: 'GET'
    });
    // Return null is there was an error (user doesn't exist)
    if (!response.ok) {
      return null;
    }
    const json = await response.json();
    return json;
  } catch {
    return null;
  }
}

export default {
  setup() {
    const { user } = useAuth0();
    const route = useRoute();
    const username = ref('');
    const searchUsername = ref('');
    const userData = ref(null);
    const endDate = ref(null);
    const isLoading = ref(true);
    const searchProfile = async function () {
      const username = searchUsername.value.trim();
      if (username === '') {
        return;
      }
      await router.push({ path: `/profile/${username}` });
    }
    const fetchUserData = async function () {
      isLoading.value = true;
      // exit early is no username is available
      if (username.value === undefined) {
        isLoading.value = false;
        userData.value = null;
        return;
      }
      userData.value = await getUser(username.value);
      // exit early is user could not be retrieved
      if (userData.value === null) {
        isLoading.value = false;
        userData.value = null;
        return;
      }
      // Process all user data from the database
      // First process the solves (determine the streaks)
      const solves = userData.value.solves ?? {};
      const keys = Object.keys(solves).map(date => convertDateToUTC(new Date(date))).sort((a, b) => a - b);
      // Store the problem solves data
      const data = [];
      // TODO Consider modifying this so that it doesn't depending on sorting
      // Use a sliding window to determine the longest streak of solved problems (in the event of a tie, the most recent streak is labeled as the longest streak)
      const len = keys.length;
      let streakStart = keys[0];
      let maxStreak = 0;
      for (let l = 0; l < len;) {
        let r = l + 1;
        while (r < len) {
          const curr = keys[r - 1].getTime();
          const next = keys[r].getTime();
          const currFormat = formatDateISO(keys[r - 1]);
          // If the dates are not continuous (or if they were not solved as daily problems), then break out of the loop
          if (next - curr !== daySeconds || solves[currFormat] !== 2 || solves[formatDateISO(keys[r])] !== 2) {
            break;
          }
          data.push({
            date: keys[r],
            count: solves[currFormat]
          })
          r++;
        }
        if (maxStreak <= r - l) {
          maxStreak = r - l;
          streakStart = keys[l];
        }
        data.push({
          date: keys[l],
          count: solves[formatDateISO(keys[l])]
        })
        l = r;
      }
      if (keys.length === 0) {
        userData.value.currentStreak = 0;
        userData.value.longestStreak = 0;
      } else {
        const lastDay = formatDateISO(keys[keys.length - 1]);
        data.push({
          date: keys[keys.length - 1],
          count: solves[lastDay]
        });
        // Get current streak (don't set to zero if problem for the day hasn't been solved)
        const today = new Date();
        // Check if today's problem has been solved
        const isTodaySolved = solves[formatDateISO(today)] === 2;
        // Convert to Date object for manipulation
        const curr = new Date(today);
        // Start from yesterday (solving today acts as a "bonus")
        curr.setDate(curr.getDate() - 1);
        const yesterday = new Date(curr);
        // Start from yesterday and move backwards day by day while the current date has been solved as a daily
        while (solves[formatDateISO(curr)] === 2) {
          curr.setDate(curr.getDate() - 1);
        }
        // Go back one day to get to the last valid day
        curr.setDate(curr.getDate() + 1);
        // Configure streak data for the longest streak
        const maxStreakEnd = new Date(streakStart);
        maxStreakEnd.setDate(streakStart.getDate() + maxStreak - 1);
        [userData.value.longestStreak, userData.value.longestStreakRange] = observeDateRange(streakStart, maxStreakEnd);
        // Configure streak data for the current streak
        [userData.value.currentStreak, userData.value.currentStreakRange] = observeDateRange(curr, isTodaySolved ? today : yesterday);
        // Set creation date with correct formatting
        userData.value.creation_date = reformatDate(userData.value.creation_date);
        // Set end date for heatmap calander
        const tomorrow = convertDateToUTC(new Date());
        tomorrow.setDate(tomorrow.getDate() + 1);
        endDate.value = formatDateISO(tomorrow);
      }
      userData.value.data = data;
      isLoading.value = false;
    }
    onMounted(async () => {
      username.value = route.params.username || user.value?.username;
      await fetchUserData();
    });
    watch(
      () => route.fullPath,
      async (newPath, oldPath) => {
        if (newPath !== oldPath) {
          // Re-render the page if the path changes but we are still on the profile page (to maintain SPA functionality)
          username.value = route.params.username || user.value?.username;
          await fetchUserData();
        }
      }
    );
    return { username, router, searchUsername, user: userData, isLoading, endDate, searchProfile };
  },
  components: {
    CalendarHeatmap, LoadingCogs
  }
}
</script>
