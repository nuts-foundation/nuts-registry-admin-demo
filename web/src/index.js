import * as Vue from 'vue'
import App from './App.vue'
import * as VueRouter from 'vue-router'

const About = { template: '<div>Provided by the Nuts comunity for demo purposes.</div>' }

const routes = [
  { path: '/', component: App },
  { path: '/about', component: About }
]

const router = VueRouter.createRouter({
  // We are using the hash history for simplicity here.
  history: VueRouter.createWebHashHistory(),
  routes // short for `routes: routes`
})

const app = Vue.createApp({})
app.use(router)
app.mount('#app')
