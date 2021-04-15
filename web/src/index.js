import { createApp } from 'vue'
import ElementPlus from 'element-plus';
import 'element-plus/lib/theme-chalk/index.css';
import App from './admin/App.vue'
import Customers from './admin/Customers.vue'
import Vendor from './admin/Vendor.vue'
import Vendors from './admin/Vendors.vue'
import Login from './Login.vue'
import * as VueRouter from 'vue-router'

const About = { template: '<div>Provided by the Nuts community for demo purposes.</div>' }

const routes = [
  { path: '/', component: Login },
  { path: '/about', component: About },
  { path: '/admin',
    component: App,
    children: [
      { path: 'customers', component: Customers },
      { path: 'vendor', component: Vendor },
      { path: 'vendors', component: Vendors },
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
