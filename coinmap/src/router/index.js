import Vue from 'vue'
import VueRouter from 'vue-router'

const HeaderFooter = () => import('../components/HeaderFooter')
const Cryptocurrency = () => import('../components/Cryptocurrency')
const Charts = () => import('../components/Charts')
const Cryptocurrency1 = () => import('../components/Cryptocurrency1.vue')
// const Charts1 = () => import('../components/Charts1')
const Login = () => import('../components/Login')
const Collection = () => import('../components/Collection')

Vue.use(VueRouter)
const routes = [
  {
    path: '/',
    redirect:'/login',
  },
  {
    path: '/login',
    component:Login
  },
  {
    path: '/cryptocurrency',
    component: HeaderFooter,
    redirect: '/cryptocurrency',
    children: [
      {
        path: '/cryptocurrency',
        component: Cryptocurrency
      },
      {
        path: '/charts',
        component: Charts
      },
      {
        path: '/collection',
        component:Collection
      },
      {
        path: '/cryptocurrency1',
        component:Cryptocurrency1
      }

    ]
  },
]

const router = new VueRouter({
  routes
})

// 挂载路由导航守卫
router.beforeEach((to,from,next) => {
  // to 将要访问的路径
  // from 从哪个路径跳转而来
  // next 放行函数
  //  next() 放行  next('/login') 强制跳转到login
  if(to.path === '/login') return next();
  //获取token
  const tokenStr=window.sessionStorage.getItem('token');
  if(!tokenStr) return next('/login');
  next();
})


export default router
