import { createAuth0 } from '@auth0/auth0-vue'
import { createApp } from 'vue'
import { createVfm } from 'vue-final-modal'
import { setupCalendar } from 'v-calendar'
import App from './App.vue'
import './assets/main.css'
import router from './router'

const app = createApp(App)

const vfm = createVfm()

app
  .use(router)
  .use(vfm)
  .use(setupCalendar)
  .use(
    createAuth0({
      domain: import.meta.env.VITE_AUTH0_DOMAIN,
      clientId: import.meta.env.VITE_AUTH0_CLIENT_ID,
      authorizationParams: {
        redirect_uri: window.location.origin,
        audience: import.meta.env.API_IDENTIFIER,
        scope: 'openid profile email',
      },
      cacheLocation: 'localstorage',
    }),
  )
  .mount('#app')
