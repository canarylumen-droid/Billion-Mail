import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/ip-pools',
  component: Layout,
  meta: { sort: 18, hidden: false, key: 'ip-pools', title: 'IP Pools' },
  children: [
    {
      path: '/ip-pools',
      name: 'IPPools',
      component: () => import('@/views/ip-pools/index.vue'),
    },
  ],
}

export default route
