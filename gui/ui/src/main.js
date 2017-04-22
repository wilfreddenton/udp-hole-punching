// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import VueSocketio from 'vue-socket.io'

import io from 'socket.io-client'
import store from './store'

Vue.config.productionTip = false
Vue.use(VueSocketio, io('http://localhost:8000/socket.io/'), store)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  template: '<App/>',
  components: { App }
})
