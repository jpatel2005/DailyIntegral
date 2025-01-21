<template>
  <div class="app-container">
    <header>
      <Navbar />
    </header>
    <main class="flex-grow">
      <RouterView />
    </main>
    <footer>
      <Footer />
    </footer>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

main {
  flex-grow: 1;
}
</style>

<script setup>
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';
</script>
