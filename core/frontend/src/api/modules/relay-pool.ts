import { instance } from '@/api'

export const listProviders = () => instance.get('/relay_pool/list')
export const updateProvider = (data: any) => instance.post('/relay_pool/update', data)
export const testProvider = (data: any) => instance.post('/relay_pool/test', data)
export const testPool = (data: any) => instance.post('/relay_pool/test_pool', data)
export const resetCounters = (data: any) => instance.post('/relay_pool/reset_counters', data)
