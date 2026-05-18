<template>
  <LoadingCogs v-if="isLoading" />
  <div v-else class="shadow-md rounded-md mx-auto mt-12 w-full">
    <div class="bg-blue-500 py-4 px-4">
      <div class="relative flex items-center justify-between w-full px-4 md:px-10">
        <div class="hidden md:block w-24"></div>
        <h2 class="text-white text-3xl md:text-4xl lg:text-5xl font-semibold tracking-wide text-center flex-grow">
          Top Users
        </h2>
        <div class="flex items-center justify-end gap-2 w-auto md:w-24">
          <button @click="goBack" :disabled="!prevCursor"
            class="flex items-center justify-center w-10 h-10 rounded-full border-2 border-blue-600 bg-transparent text-gray-300 transition-all hover:border-blue-400 hover:text-blue-400 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:border-gray-600 disabled:hover:text-gray-300">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5"
              stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
            </svg>
          </button>
          <button @click="goForward" :disabled="!nextCursor"
            class="flex items-center justify-center w-10 h-10 rounded-full border-2 border-blue-600 bg-transparent text-gray-300 transition-all hover:border-blue-400 hover:text-blue-400 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:border-gray-600 disabled:hover:text-gray-300">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5"
              stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
          </button>
        </div>
      </div>
    </div>
    <div class="grid grid-cols-6 bg-blue-400 text-white font-semibold text-xs lg:text-[1.7em] py-3 lg:py-6 px-6">
      <span></span>
      <span class="text-left ml-[-1rem] lg:ml-0">Username</span>
      <span class="text-left ml-[1rem] lg:ml-0">Daily Problems Solved</span>
      <span class="text-left ml-[1.5rem] lg:ml-0">Total Problems Solved</span>
      <span class="text-left ml-[1.75rem] lg:ml-0">Profile Page</span>
      <span class="text-left ml-[1rem] lg:ml-0">Account Creation Date</span>
    </div>
    <ul class="divide-y divide-gray-500">
      <li v-for="(user, index) in users" :key="index" class="grid grid-cols-6 items-center py-4 px-6">
        <span class="text-white font-bold text-md lg:text-3xl">{{ startingRank + index <= 3 ? ["🥇", "🥈", "🥉"
        ][startingRank + index - 1] : (startingRank + index) + '.' }}</span>
            <h3 class="text-gray-300 text-sm lg:text-xl font-medium ml-[-1rem] lg:ml-0">{{ user.username }}</h3>
            <p class="text-white text-sm lg:text-[2em] font-serif ml-12 lg:ml-0">{{ user.daily_problems }}</p>
            <p class="text-white text-sm lg:text-[2em] font-serif ml-8 lg:ml-0">{{ user.total_problems }}</p>
            <span class="text-blue-400 text-sm lg:text-4xl pi pi-user hover:cursor-pointer w-fit ml-8 lg:ml-0"
              @click="router.push({ path: `/profile/${user.username}` })"></span>
            <p class="text-white text-xs lg:text-3xl font-light ml-6 lg:ml-0">{{ reformatDate(user.creation_date) }}</p>
      </li>
    </ul>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue';
import router from '../router';
import { reformatDate } from '../time';
import LoadingCogs from '../components/LoadingCogs.vue';

const getLeaderboardUsers = async function (cursor = null) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/leaderboard${cursor ? `?cursor=${cursor}` : ''}`, {
      method: 'GET'
    });
    if (!response.ok) {
      return null;
    }
    const json = await response.json();
    //  return this is there was an error
    if (json.error !== undefined) {
      return null;
    }
    return { users: json.data, prevCursor: json.prev_cursor, nextCursor: json.next_cursor, startingRank: json.starting_rank };
  } catch {
    return null;
  }
}

export default {
  setup() {
    const isLoading = ref(true);
    const users = ref(null);
    const prevCursor = ref(null);
    const nextCursor = ref(null);
    const startingRank = ref(1);
    // go back a page
    const goBack = async function () {
      if (!prevCursor.value) {
        return;
      }
      isLoading.value = true;
      const result = await getLeaderboardUsers(prevCursor.value);
      if (result !== null) {
        users.value = result.users;
        prevCursor.value = result.prevCursor;
        nextCursor.value = result.nextCursor;
        startingRank.value = result.startingRank;
      }
      isLoading.value = false;
    }
    // go forward a page
    const goForward = async function () {
      if (!nextCursor.value) {
        return;
      }
      isLoading.value = true;
      const result = await getLeaderboardUsers(nextCursor.value);
      if (result !== null) {
        users.value = result.users;
        prevCursor.value = result.prevCursor;
        nextCursor.value = result.nextCursor;
        startingRank.value = result.startingRank;
      }
      isLoading.value = false;
    }
    onMounted(async () => {
      isLoading.value = true;
      const result = await getLeaderboardUsers();
      if (result !== null) {
        users.value = result.users;
        prevCursor.value = result.prevCursor;
        nextCursor.value = result.nextCursor;
        startingRank.value = result.startingRank;
      }
      isLoading.value = false;
    })
    return { users, router, reformatDate, isLoading, prevCursor, nextCursor, startingRank, goBack, goForward };
  },
  components: {
    LoadingCogs
  }
};
</script>
