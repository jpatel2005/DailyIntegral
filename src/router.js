import { createWebHistory, createRouter } from 'vue-router'

import CallbackPage from './pages/CallbackPage.vue'
import HomePage from './pages/HomePage.vue'
import ProfilePage from './pages/ProfilePage.vue'
import ProblemsPage from './pages/ProblemsPage.vue'
import LeaderboardPage from './pages/LeaderboardPage.vue'
import NotFoundPage from './pages/NotFoundPage.vue'

const routes = [
  { path: '/', component: HomePage },
  { path: '/profile/:username?', component: ProfilePage },
  { path: '/problem/:pdate?', component: ProblemsPage },
  { path: '/leaderboard', component: LeaderboardPage },
  { path: '/callback', component: CallbackPage },
  { path: '/:pathMatch(.*)*', component: NotFoundPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
