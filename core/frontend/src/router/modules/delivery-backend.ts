import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/delivery-backend',
  component: Layout,
  meta: { sort: 20, hidden: false, key: 'delivery-backend', title: 'Delivery Backend' },
  children: [
    {
      path: '/delivery-backend',
      name: 'DeliveryBackend',
      component: () => import('@/views/delivery-backend/index.vue'),
    },
  ],
}

export default route
