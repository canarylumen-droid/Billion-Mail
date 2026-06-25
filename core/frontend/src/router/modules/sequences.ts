import { RouteRecordRaw } from 'vue-router'
import { Layout } from '@/router/constant'

const route: RouteRecordRaw = {
  path: '/sequences',
  name: 'SequencesLayout',
  component: Layout,
  meta: {
    key: 'sequences',
    title: 'Sequences',
    titleKey: 'layout.menu.sequences',
  },
  children: [
    {
      path: '',
      name: 'SequencesList',
      meta: { title: 'Sequences', titleKey: 'layout.menu.sequences' },
      component: () => import('@/views/sequences/index.vue'),
    },
    {
      path: 'edit',
      name: 'SequencesEdit',
      meta: { title: 'Sequences', hidden: true },
      component: () => import('@/views/sequences/edit.vue'),
    },
  ],
}

export default route
