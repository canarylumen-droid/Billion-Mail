import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/smtp-relay',
  component: Layout,
  meta: { sort: 17, hidden: false, key: 'smtp-relay', title: 'SMTP Relay' },
  children: [
    {
      path: '/smtp-relay',
      name: 'SMTPRelay',
      component: () => import('@/views/smtp-relay/index.vue'),
    },
  ],
}

export default route
