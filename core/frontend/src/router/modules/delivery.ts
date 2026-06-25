import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/delivery',
  component: Layout,
  meta: { sort: 19, hidden: false, key: 'delivery', title: 'Delivery' },
  children: [
    {
      path: '/delivery',
      name: 'Delivery',
      component: () => import('@/views/delivery/index.vue'),
    },
  ],
}

export default route
