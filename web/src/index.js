import {createApp, h} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import AdminMenu from './admin/AdminMenu.vue'
import Landing from './Landing.vue'
import PublicMenu from './layout/PublicMenu.vue'
import Login from './Login.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Customers from './admin/Customers.vue'
import ServiceProvider from './admin/ServiceProvider.vue'
import NewCustomer from './admin/NewCustomer.vue'
import Modal from './layout/Modal.vue'

const routes = [
  {path: '/', components: {default: Landing, menu: PublicMenu}},
  {path: '/login', components: {default: Login, menu: PublicMenu}},
  {path: '/logout', components: {default: Logout, menu: PublicMenu}},
  {
    path: '/admin',
    components: {
      default: AdminApp,
    },
    children: [
      {
        path: '',
        name: 'admin.home',
        redirect: '/admin/customers'
      },
      {
        path: 'customers',
        name: 'admin.customers',
        component: Customers,
        children: [
          {
            path: 'new',
            name: 'admin.newCustomer',
            components: {
              modal: NewCustomer
            }
          },
        ]
      },
      {
        path: 'service-provider',
        name: 'admin.serviceProvider',
        component: ServiceProvider
      }
    ],
    meta: {requiresAuth: true}
  },
  {path: '/:pathMatch*', name: 'NotFound', components: {default: NotFound, menu: PublicMenu}},
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

router.beforeEach((to, from) => {
  if (to.meta.requiresAuth) {
    if (localStorage.getItem("session")) {
      return true
    }
    return '/login'
  }
})

const app = createApp(App)

app.use(router)
app.mount('#app')
app.component('nrad-modal', Modal)