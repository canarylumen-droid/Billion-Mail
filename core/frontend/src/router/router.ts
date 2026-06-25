import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { isDev } from '@/utils'

// Routes reflect list — order determines sidebar position
const routesReflectList = [
        'Overview',
        'Email Marketing',
        'template',
        'Send API',
        'Contacts',
        'Sequences',
        'Leads',
        'Enrichment',
        'MailDomain',
        'MailBoxes',
        'SMTP',
        'SMTP Relay',
        'IP Pools',
        'Delivery',
        'Delivery Backend',
        'Logs',
        'Settings',
        'Automation',
        'Video Outreach',
]

// Import routes from modules.
// NOTE: Rspack webpackContext keys include the "./" prefix (e.g. "./api.ts"),
// so the regex must start with "\.\/" to match. The original "/^[^.]+\.ts$/"
// was broken in Rspack because it rejected all keys starting with ".".
// Test files like "./api.test.ts" are excluded since [^.]+ stops at the first ".".
const modules = import.meta.webpackContext('./modules', {
        recursive: false,
        regExp: /^\.\/[^.]+\.ts$/,
})

// Module routes
export let menuList: RouteRecordRaw[] = []

// Iterate through the module list to generate module routes
for (const path of modules.keys()) {
        const mod = modules(path)
        // Support both ES module format (mod.default) and direct CommonJS export
        const route: RouteRecordRaw | undefined =
                (mod?.default?.path ? mod.default : undefined) ??
                ((mod as RouteRecordRaw)?.path ? (mod as RouteRecordRaw) : undefined)
        if (route) {
                menuList.push(route)
        }
}

// Sort module routes into the fixed sidebar order
menuList = menuList.reduce((p: RouteRecordRaw[], v: RouteRecordRaw) => {
        const routeIndex = routesReflectList.findIndex(item => item === v.meta?.title)
        if (routeIndex >= 0) {
                p[routeIndex] = v
        }
        return p
}, [] as RouteRecordRaw[])

// Remove undefined holes from the sparse array
const filteredMenuList = menuList.filter(Boolean)

const otherArray: RouteRecordRaw[] = []

if (isDev) {
        otherArray.push({
                path: '/test',
                name: 'Test',
                component: () => import('@/views/test/index.vue'),
        })
}

export const routes: RouteRecordRaw[] = [
        {
                path: '/login',
                name: 'Login',
                component: () => import('@/views/login/index.vue'),
        },
        {
                path: '/',
                redirect: '/overview',
        },
        ...filteredMenuList,
        ...otherArray,
]

const router = createRouter({
        history: createWebHistory('/'),
        routes,
        strict: false,
        scrollBehavior: () => ({ left: 0, top: 0 }),
})

export default router
