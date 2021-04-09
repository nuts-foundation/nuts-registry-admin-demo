import { createApp } from 'vue'
import ElementPlus from 'element-plus';
import 'element-plus/lib/theme-chalk/index.css';
import App from './admin/App.vue'
import Customers from './admin/Customers.vue'
import Login from './Login.vue'
import Landing from './Landing.vue'
import * as VueRouter from 'vue-router'

const About = { template: '<div>Provided by the Nuts community for demo purposes.</div>' }

const routes = [
  { path: '/', component: Landing },
  { path: '/about', component: About },
  { path: '/login', component: Login },
  { path: '/admin',
    component: App,
    children: [
      {
        path: 'customers',
        component: Customers
      }
    ]
  }
]

const router = VueRouter.createRouter({
  // We are using the hash history for simplicity here.
  history: VueRouter.createWebHashHistory(),
  routes // short for `routes: routes`
})

const app = createApp({})
app.use(ElementPlus)
app.use(router)
app.mount('#app')
