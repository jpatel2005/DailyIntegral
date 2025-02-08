<template>
  <nav class="border-gray-200 bg-gray-900 rounded-2xl mt-4 mx-3.5 mb-2">
    <div class="flex flex-wrap items-center justify-between lg:mx-16 p-2 py-4 lg:p-4 lg:my-0">
      <RouterLink to="/" class="flex items-center space-x-3 rtl:space-x-reverse">
        <img :class="{ hidden: logoLoading }" src="../assets/logo.png" @load="() => logoLoading = false" class="h-20"
          alt="DailyIntegral Logo" />
        <div :class="{ hidden: !logoLoading }" src="../assets/logo.png" class="img-skeleton !mx-0"
          alt="DailyIntegral Logo">
        </div>
        <span
          class="self-center text-3xl md:text-4xl lg:text-5xl font-semibold whitespace-nowrap text-white">DailyIntegral</span>
      </RouterLink>
      <div class="flex md:order-2 space-x-3 md:space-x-0 rtl:space-x-reverse">
        <div class="w-[180px] h-[48px] flex items-center justify-center">
          <button v-if="!isLoading && !isAuthenticated" type="button"
            class="text-white focus:outline-none font-medium rounded-lg text-xl px-4 py-2 text-center bg-blue-600 hover:bg-blue-700"
            @click="handleLogin">Log In / Sign Up</button>
          <button v-else-if="!isLoading" type="button"
            class="text-white focus:outline-none font-medium rounded-lg text-xl px-4 py-2 text-center bg-blue-600 hover:bg-blue-700"
            @click="handleLogout">Sign Out</button>
        </div>
        <button data-collapse-toggle="navbar-cta" type="button"
          class="inline-flex items-center p-2 w-10 h-10 justify-center text-sm rounded-lg md:hidden focus:outline-none text-gray-400 hover:bg-gray-700"
          aria-controls="navbar-cta" aria-expanded="false" ref="navbarMenu">
          <span class="sr-only">Open main menu</span>
          <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 14">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M1 1h15M1 7h15M1 13h15" />
          </svg>
        </button>
      </div>
      <div class="items-center justify-between hidden w-full md:flex md:w-auto md:order-1" id="navbar-cta">
        <ul
          class="flex flex-col font-medium p-4 md:p-0 mt-4 border rounded-lg md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0 md:bg-gray-900 border-gray-700">
          <li>
            <RouterLink to="/" class="block py-2 px-3 md:p-0 text-white bg-blue-700 rounded md:bg-transparent text-xl"
              active-class="md:text-blue-500" aria-current="page" @click="hideNavbarMenu">Home
            </RouterLink>
          </li>
          <li>
            <RouterLink to="/profile"
              class="block py-2 px-3 md:p-0 rounded md:hover:text-blue-500 text-white hover:bg-gray-700 hover:text-white md:hover:bg-transparent border-gray-700 text-xl"
              active-class="md:text-blue-500" :class="{ 'md:text-blue-500': $route.path.startsWith('/profile') }"
              @click="hideNavbarMenu">
              Profile
            </RouterLink>
          </li>
          <li>
            <RouterLink to="/leaderboard"
              class="block py-2 px-3 md:p-0 rounded md:hover:text-blue-500 text-white hover:bg-gray-700 hover:text-white md:hover:bg-transparent border-gray-700 text-xl"
              active-class="md:text-blue-500" @click="hideNavbarMenu">
              Leaderboard</RouterLink>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.img-skeleton {
  width: 5rem;
  height: 5rem;
}
</style>

<script>
import { onMounted, ref } from 'vue';
import { initFlowbite } from 'flowbite';
import { useAuth0 } from '@auth0/auth0-vue';
const logoLoading = ref(true);
const navbarMenu = ref(true);
export default {
  setup() {
    logoLoading.value = true;
    const { loginWithRedirect, logout, isLoading, isAuthenticated } = useAuth0();
    const handleLogin = () =>
      loginWithRedirect({
        appState: {
          target: '/problem',
        },
        authorizationParams: {
          prompt: 'login',
          scope: 'openid profile email',
        }
      });
    const handleLogout = () =>
      logout({
        logoutParams: {
          returnTo: window.location.origin
        }
      });
    // Function for hiding navbar menu after option is selected
    const hideNavbarMenu = function () {
      if (!navbarMenu.value) {
        return;
      }
      navbarMenu.value.click();
    }
    onMounted(() => {
      initFlowbite();
    })
    return { handleLogin, handleLogout, isLoading, isAuthenticated, logoLoading, navbarMenu, hideNavbarMenu }
  }
}
</script>
