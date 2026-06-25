import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/crm',
  name: 'CRMLayout',
  component: Layout,
  meta: {
    key: 'crm',
    title: 'Enrichment',
    titleKey: 'layout.menu.crm',
  },
  children: [
    {
      path: '',
      name: 'CRMInbox',
      meta: { title: 'Enrichment', titleKey: 'layout.menu.crm' },
      component: () => import('@/views/crm/index.vue'),
    },
  ],
}

export default route
