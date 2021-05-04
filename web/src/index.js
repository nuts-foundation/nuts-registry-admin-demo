import {createApp, h} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import Landing from './Landing.vue'
import Login from './Login.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Customers from './admin/Customers.vue'
import ServiceProvider from './admin/ServiceProvider.vue'
import NewCustomer from './admin/NewCustomer.vue'
import EditCustomer from './admin/EditCustomer.vue'
import Modal from './components/Modal.vue'
import Api from './plugins/api'

const routes = [
  {path: '/', component: Landing},
  {path: '/login', component: Login},
  {path: '/logout', component: Logout},
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
          {
            path: ':id/edit',
            name: 'admin.editCustomer',
            components: {
              modal: EditCustomer
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
  {path: '/:pathMatch*', name: 'NotFound', component: NotFound}
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
app.use(Api)
app.mount('#app')
app.component('nrad-modal', Modal)