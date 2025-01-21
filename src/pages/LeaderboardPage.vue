<template>
  <LoadingCogs v-if="isLoading" />
  <div v-else class="shadow-md rounded-md mx-auto mt-12 w-full">
    <div class="bg-blue-500 py-4 px-4">
      <h2 class="text-3xl lg:text-5xl font-bold text-white text-center">Top Users</h2>
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
        <span class="text-white font-bold text-md lg:text-3xl">{{ index <= 2 ? ["🥇", "🥈", "🥉"][index] : (index + 1)
          + '.' }}</span>
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

export default {
  setup() {
    const getLeaderboardUsers = async function () {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/leaderboard`, {
          method: 'GET'
        });
        const json = await response.json();
        //  return this is there was an error
        if (json.error !== undefined) {
          return null;
        }
        return json;
      } catch {
        return null;
      }
    }
    const isLoading = ref(true);
    const users = ref(null);
    onMounted(async () => {
      isLoading.value = true;
      users.value = await getLeaderboardUsers();
      isLoading.value = false;
    })
    return { users, router, reformatDate, isLoading };
  },
  components: {
    LoadingCogs
  }
};
</script>
