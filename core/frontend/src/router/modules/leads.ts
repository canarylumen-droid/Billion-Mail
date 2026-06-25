import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/leads',
  name: 'LeadsLayout',
  component: Layout,
  meta: {
    key: 'leads',
    title: 'Leads',
    titleKey: 'layout.menu.leads',
  },
  children: [
    {
      path: '',
      name: 'LeadsList',
      meta: { title: 'Leads', titleKey: 'layout.menu.leads' },
      component: () => import('@/views/leads/index.vue'),
    },
  ],
}

export default route
