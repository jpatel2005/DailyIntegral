<template>
  <PageLayout v-if="error">
    <div class="content-layout">
      <h1 id="page-title" class="content__title">Error</h1>
      <div class="content__body">
        <p id="page-description">
          <span>{{ error.message }}</span>
        </p>
      </div>
    </div>
  </PageLayout>
  <div class="page-layout" v-else>
    <Navbar />
    <div class="page-layout__content" />
  </div>
</template>

<script setup>
import Navbar from '../components/Navbar.vue';
import { useAuth0 } from "@auth0/auth0-vue";

const { error } = useAuth0();
</script>
