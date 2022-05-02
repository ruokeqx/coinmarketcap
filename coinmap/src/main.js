import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import './assets/fonts/iconfont.css'
import axios from 'axios'
import * as echarts from 'echarts'
import baseUrl from './config.js'

Vue.prototype.$echarts = echarts
Vue.use(ElementUI)
axios.defaults.baseURL = baseUrl
// axios请求拦截
// axios.interceptors.request.use( config =>{
//   // 为请求头对象添加Token验证的Authorization字段
//   config.headers.Authorization = window.sessionStorage.getItem('token');
//   return config;
// })
Vue.prototype.$http = axios
Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
