<template>
  <div class="flex">
    <!-- Navigation -->
    <nav class="main-nav shadow-lg fixed flex flex-col justify-between top-0 h-full w-72 bg-white z-10">
      <div>
        <div class="px-6 pt-8">
          <div class="flex items-center">
            <div class="w-10 mr-2">
              <svg id="Layer_1" x="0px" y="0px" viewBox="0 0 272.1 119.8" style="enable-background:new 0 0 272.1 119.8;"
                   xml:space="preserve">
              <path fill="#1EAAE7" class="st0"
                    d="M210.6,0h-82.1v24.8h29.8v94.9h24.8V24.8h9.3c1.6-7.4,5.4-14.2,11-19.4C205.6,3.4,208,1.5,210.6,0z"/>
                <path fill="#1EAAE7" class="st0"
                      d="M264.9,61.5c-4.8-5-14.1-10.5-28-16.5c-10.1-4.3-15.1-8.6-15.1-12.9c-0.1-2.3,0.9-4.5,2.6-6.1  c1.7-1.6,3.9-2.4,6.2-2.4l0,0c5,0,9.3,3.5,12.7,10.5l21.6-11.5c-7-15-17.8-22.4-32.3-22.5l0,0c-6.9-0.1-13.7,1.8-19.6,5.4  c-2,1.2-3.8,2.7-5.5,4.2c-4.3,4.1-7.4,9.3-8.9,15.1c-0.8,3.1-1.2,6.2-1.2,9.4c0,7.4,2,13.5,5.9,18.3c3.1,3.7,8.4,7.4,16,11.2  c2.7,1.3,5.4,2.6,8,3.9c2.6,1.2,5.3,2.5,7.9,3.8c6.8,4,10.2,8.3,10.2,13c0,8.1-4.4,12.1-13.2,12.1l0,0c-4.2,0.1-8.3-1.5-11.4-4.4  c-1.6-1.5-4-4.9-7.1-10.4l-22.3,10.3c7.9,18.5,21.3,27.8,40.3,27.8l0,0c11.1,0,20.5-3.5,28.4-10.5c7.9-6.9,12.3-17,12.1-27.4  C272.1,73.3,269.7,66.5,264.9,61.5z"/>
                <path fill="#1EAAE7" class="st0"
                      d="M128.5,30.8v42c0,17.1-6.7,25.6-20.2,25.6c-13.1,0-19.8-8-20.2-24.1V46.4c0-9.5-1-16.6-3-21.4  c-2.2-5.6-5.8-10.6-10.5-14.5C66.6,3.5,56.4,0,44,0C31.7,0,21.5,3.5,13.3,10.5C8.8,14.4,5.2,19.4,3.1,25C0.9,31.9-0.2,39.1,0,46.4  v73.3c0.6-0.7,1.3-1.4,2-2.1c6-5.7,13.7-9.4,21.9-10.6v-60c0-17.1,6.7-25.6,20.2-25.6c13.1,0,19.8,8,20.2,24.1v27.9  c0,9.5,1,16.6,3,21.4c2.2,5.6,5.8,10.6,10.4,14.5c8,7,18.2,10.5,30.6,10.5c12.4,0,22.6-3.5,30.7-10.5c4.6-3.9,8.1-8.9,10.3-14.5  c2.2-6.9,3.2-14.1,3.1-21.4V30.8H128.5z"/>
            </svg>
            </div>

            <div class="main-title text-xl">
              <small>/</small> <span class="font-semibold">Registry Admin</span>
            </div>
          </div>
        </div>

        <div class="px-3 mt-6">
          <div class="grid grid-cols-1">
            <router-link
                id="vendor-menu-link"
                :to="{name: 'admin.serviceProvider'}"
                active-class="menu-link-active"
                class="menu-link">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" enable-background="new 0 0 24 24" viewBox="0 0 24 24"
                     fill="currentColor">
                  <g>
                    <rect fill="none" height="24" width="24"/>
                  </g>
                  <g>
                    <g>
                      <path
                          d="M20,6h-4V4c0-1.1-0.9-2-2-2h-4C8.9,2,8,2.9,8,4v2H4C2.9,6,2,6.9,2,8v12c0,1.1,0.9,2,2,2h16c1.1,0,2-0.9,2-2V8 C22,6.9,21.1,6,20,6z M10,4h4v2h-4V4z M20,20H4V8h16V20z"/>
                      <polygon points="13,10 11,10 11,13 8,13 8,15 11,15 11,18 13,18 13,15 16,15 16,13 13,13"/>
                    </g>
                  </g>
                </svg>
              </div>

              You as a Service Provider
            </router-link>
            <router-link
                id="careorganizations-menu-link"
                :to="{name: 'admin.customers'}"
                active-class="menu-link-active"
                class="menu-link">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M0 0h24v24H0V0z" fill="none"/>
                  <path
                      d="M12 7V3H2v18h20V7H12zM6 19H4v-2h2v2zm0-4H4v-2h2v2zm0-4H4V9h2v2zm0-4H4V5h2v2zm4 12H8v-2h2v2zm0-4H8v-2h2v2zm0-4H8V9h2v2zm0-4H8V5h2v2zm10 12h-8v-2h2v-2h-2v-2h2v-2h-2V9h8v10zm-2-8h-2v2h2v-2zm0 4h-2v2h2v-2z"/>
                </svg>
              </div>

              Your care organizations
            </router-link>
            <router-link
                :to="{name: 'admin.manageVCs'}"
                active-class="menu-link-active"
                class="menu-link">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <g fill="none">
                    <path d="M0 0h24v24H0V0z"/>
                    <path d="M0 0h24v24H0V0z" opacity=".87"/>
                  </g>
                  <path
                      d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zM9 6c0-1.66 1.34-3 3-3s3 1.34 3 3v2H9V6zm9 14H6V10h12v10zm-6-3c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2z"/>
                </svg>
              </div>

              Manage Credential Issuers
            </router-link>
            <router-link
                :to="{name: 'admin.organizationRegistry'}"
                active-class="menu-link-active"
                class="menu-link">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" enable-background="new 0 0 24 24" viewBox="0 0 24 24"
                     fill="currentColor">
                  <g>
                    <path d="M0,0h24v24H0V0z" fill="none"/>
                  </g>
                  <g>
                    <path
                        d="M7,9H2V7h5V9z M7,12H2v2h5V12z M20.59,19l-3.83-3.83C15.96,15.69,15.02,16,14,16c-2.76,0-5-2.24-5-5s2.24-5,5-5s5,2.24,5,5 c0,1.02-0.31,1.96-0.83,2.75L22,17.59L20.59,19z M17,11c0-1.65-1.35-3-3-3s-3,1.35-3,3s1.35,3,3,3S17,12.65,17,11z M2,19h10v-2H2 V19z"/>
                  </g>
                </svg>
              </div>

              Care Organization Registry
            </router-link>
            <router-link
                :to="{name: 'admin.issueVC'}"
                active-class="menu-link-active"
                class="menu-link">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                     stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
                </svg>
              </div>

              Issue Verifiable Credential
            </router-link>
          </div>
        </div>
      </div>

      <div class="px-3 mt-6">
        <router-link
            to="/logout"
            active-class="menu-link-active"
            class="menu-link">
          <div class="w-5 h-5 mr-3">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <g>
                <path d="M0,0h24v24H0V0z" fill="none"/>
              </g>
              <g>
                <path
                    d="M17,8l-1.41,1.41L17.17,11H9v2h8.17l-1.58,1.58L17,16l4-4L17,8z M5,5h7V3H5C3.9,3,3,3.9,3,5v14c0,1.1,0.9,2,2,2h7v-2H5V5z"/>
              </g>
            </svg>
          </div>

          Logout
        </router-link>
      </div>
    </nav>

    <main class="ml-72 mb-14 mt-8 px-12 w-full">
      <!-- Main content -->
      <div class="w-full m-auto max-w-screen-2xl">
        <router-view @statusUpdate="updateStatus"></router-view>
      </div>

      <status-bar :statusMessage="eventMessage"></status-bar>
    </main>
  </div>

</template>

<script>
import StatusBar from '../components/StatusBar.vue'

export default {
  components: {StatusBar},
  data() {
    return {
      eventMessage: ''
    }
  },
  methods: {
    updateStatus(status) {
      this.eventMessage = status
    }
  }
}
</script>

<style>
.menu-link {
  @apply rounded mb-2 py-3 px-3 inline-flex items-center text-gray-500 text-sm;
}

.menu-link:hover {
  @apply text-gray-800 bg-gray-100;
}

.menu-link-active, .menu-link-active:hover {
  @apply bg-gray-200 text-black;
}

.main-title {
  color: #1683af;
}
</style>
