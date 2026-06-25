import { instance } from '@/api'

export const listBackends = () => instance.get('/delivery/backends/list')
export const createBackend = (data: any) => instance.post('/delivery/backends/create', data)
export const updateBackend = (data: any) => instance.post('/delivery/backends/update', data)
export const deleteBackend = (data: any) => instance.post('/delivery/backends/delete', data)
export const testBackend = (data: any) => instance.post('/delivery/backends/test', data)
export const setDefaultBackend = (data: any) => instance.post('/delivery/backends/set_default', data)
