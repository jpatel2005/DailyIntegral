<template>
  <div class="flex flex-col justify-center items-center h-full px-4">
    <div class="text-center">
      <h1 class="text-white text-4xl font-bold mb-12">Oops! The page you're looking for doesn't exist.</h1>
      <router-link to="/"
        class="inline-block px-6 py-3 bg-blue-600 text-white text-sm font-medium rounded-lg shadow-md hover:bg-blue-700 transition-colors">
        Go Back Home</router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: "NotFoundPage",
};
</script>

<style scoped>
/* Add any additional scoped styles here if needed */
</style>
