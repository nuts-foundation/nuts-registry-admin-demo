import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import Login from './Login.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Customers from './admin/Customers.vue'
import ServiceProvider from './admin/ServiceProvider.vue'
import NewEndpoint from './admin/NewEndpoint.vue'
import NewCustomer from './admin/NewCustomer.vue'
import EditCustomer from './admin/EditCustomer.vue'
import OrganizationRegistry from './admin/OrganizationRegistry.vue'
import ManageVCs from './admin/ManageVCs.vue'
import IssueVC from './admin/IssueVC.vue'
import Api from './plugins/api'
import NewCompoundService from './admin/NewCompoundService.vue'
import EditCompoundService from './admin/EditCompoundService.vue'

const routes = [
  { path: '/', component: Login },
  {
    name: 'login',
    path: '/login',
    component: Login
  },
  {
    name: 'logout',
    path: '/logout',
    component: Logout
  },
  {
    path: '/admin',
    components: {
      default: AdminApp
    },
    children: [
      {
        path: '',
        name: 'admin.home',
        redirect: '/admin/service-provider'
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
          }
        ]
      },
      {
        path: 'service-provider',
        name: 'admin.serviceProvider',
        component: ServiceProvider,
        children: [
          {
            path: 'endpoints/new',
            name: 'admin.newEndpoint',
            components: {
              modal: NewEndpoint
            }
          },
          {
            path: 'services/new',
            name: 'admin.newCompoundService',
            components: {
              modal: NewCompoundService
            }
          },
          {
            path: 'services/:serviceID/edit',
            name: 'admin.editCompoundService',
            components: {
              modal: EditCompoundService
            }
          }
        ]
      },
      {
        path: 'manage-vcs',
        name: 'admin.manageVCs',
        component: ManageVCs
      },
      {
        path: 'organization-registry',
        name: 'admin.organizationRegistry',
        component: OrganizationRegistry
      },
      {
        path: 'issue-vc',
        name: 'admin.issueVC',
        component: IssueVC
      },
    ],
    meta: { requiresAuth: true }
  },
  { path: '/:pathMatch*', name: 'NotFound', component: NotFound }
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

router.beforeEach((to, from) => {
  if (to.meta.requiresAuth) {
    if (localStorage.getItem('session')) {
      return true
    }
    return '/login'
  }
})

const app = createApp(App)

app.use(router)
app.use(Api, { forbiddenRoute: { name: 'logout' } })
app.mount('#app')
