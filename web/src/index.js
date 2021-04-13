import {createApp, h} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import AdminMenu from './admin/AdminMenu.vue'
import Landing from './Landing.vue'
import PublicMenu from './layout/PublicMenu.vue'
import Login from './Login.vue'
import Customers from './admin/Customers.vue'

const routes = [
  {path: '/', components: {default: Landing, menu: PublicMenu}},
  {path: '/login', components: {default: Login, menu: PublicMenu}},
  {
    path: '/admin',
    components: {default: AdminApp, menu: AdminMenu},
    children: [
      {
        path: '',
        component: Customers
      }
    ]
  }
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

// const app = createApp({})
const app = createApp(App)

app.use(router)
app.mount('#app')
