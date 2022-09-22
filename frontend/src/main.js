import { createApp } from 'vue'
import PrimeVue from 'primevue/config'

import App from './App.vue'
import 'primevue/resources/themes/luna-amber/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
import './main.css'

const app = createApp(App)

app.use(PrimeVue)

app.mount('#app')
