import { instance } from '@/api'

// ── SMTP Credentials ──────────────────────────────────────────────────────────
export const listCredentials = () => instance.get('/esp/credentials/list')
export const createCredential = (data: any) => instance.post('/esp/credentials/create', data)
export const deleteCredential = (data: any) => instance.post('/esp/credentials/delete', data)
export const regeneratePassword = (data: any) => instance.post('/esp/credentials/regenerate', data)

// ── Sending Domains ───────────────────────────────────────────────────────────
export const listSendingDomains = () => instance.get('/esp/domains/list')
export const addSendingDomain = (data: any) => instance.post('/esp/domains/add', data)
export const verifyDomain = (data: any) => instance.post('/esp/domains/verify', data)
export const deleteSendingDomain = (data: any) => instance.post('/esp/domains/delete', data)

// ── IP Pools ──────────────────────────────────────────────────────────────────
export const listIpPools = () => instance.get('/esp/pools/list')
export const createIpPool = (data: any) => instance.post('/esp/pools/create', data)
export const deleteIpPool = (data: any) => instance.post('/esp/pools/delete', data)

// ── Analytics ─────────────────────────────────────────────────────────────────
export const getDeliveryStats = (params?: any) => instance.get('/esp/analytics/stats', { params })

// ── Suppression ───────────────────────────────────────────────────────────────
export const listSuppression = (params?: any) => instance.get('/esp/suppression/list', { params })
export const addSuppression = (data: any) => instance.post('/esp/suppression/add', data)
export const deleteSuppression = (data: any) => instance.post('/esp/suppression/delete', data)

// ── MTA Servers ───────────────────────────────────────────────────────────────
export const listMtaServers = () => instance.get('/esp/mta/list')
export const addMtaServer = (data: any) => instance.post('/esp/mta/add', data)
export const deleteMtaServer = (data: any) => instance.post('/esp/mta/delete', data)
