// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.

// Element
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import Prism from 'vue-prismjs'
import 'prismjs/themes/prism.css'
import Vue from 'vue'
import App from './App'
import router from './router'

Vue.config.productionTip = false
Vue.use(ElementUI)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App, Prism },
  template: '<App/>'
})
