import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/relay-pool',
  component: Layout,
  meta: { sort: 21, hidden: false, key: 'relay-pool', title: 'Relay Providers' },
  children: [
    {
      path: '/relay-pool',
      name: 'RelayPool',
      component: () => import('@/views/relay-pool/index.vue'),
    },
  ],
}

export default route
