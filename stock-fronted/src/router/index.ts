import { createRouter, createWebHistory } from 'vue-router'
import StocksComponent from '../components/stocks/StocksComponent.vue'
import RecommendationsComponente from '../components/stocks/RecommendationsComponente.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/stocks'
    },
    {
      path: '/stocks',
      name: 'stocks',
      component: StocksComponent
    },
    {
      path: '/recommendations',
      name: 'recommendations', 
      component: RecommendationsComponente
    }
  ]
})

export default router 