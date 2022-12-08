import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import routes from 'virtual:generated-pages'
import App from './App.vue'

import Toast, { PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";
import { anu } from 'anu-vue'

import '@unocss/reset/tailwind.css'
import 'video.js/dist/video-js.css'
import './styles/main.css'
import 'uno.css'

// anu styles
import 'anu-vue/dist/style.css'

// default theme styles
import '@anu-vue/preset-theme-default/dist/styles.scss'

const app = createApp(App)
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

app.use(router)

const options:PluginOptions = {
  timeout: 8000,
  // You can set your default options here
};

app.use(Toast, options);

app.use(anu)

app.mount('#app')
