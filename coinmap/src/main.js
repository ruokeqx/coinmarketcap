import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import './assets/fonts/iconfont.css'
import axios from 'axios'
import * as echarts from 'echarts'
Vue.prototype.$echarts = echarts

Vue.use(ElementUI)

axios.defaults.baseURL = 'http://192.168.1.137:8080/'
// axios.defaults.baseURL = 'http://127.0.0.1'

// axios请求拦截
// axios.interceptors.request.use( config =>{
//   // 为请求头对象添加Token验证的Authorization字段
//   config.headers.Authorization = window.sessionStorage.getItem('token');
//   return config;
// })
Vue.prototype.$http = axios
Vue.config.productionTip = false
// import VueAxios from 'vue-axios'

// Vue.use(VueAxios, axios)


new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
