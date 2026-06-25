import { instance } from '@/api'

export interface Lead {
  id?: number
  email: string
  first_name?: string
  last_name?: string
  company?: string
  title?: string
  phone?: string
  linkedin?: string
  website?: string
  attributes?: Record<string, string>
  status: 'not_started' | 'in_sequence' | 'replied' | 'bounced' | 'unsubscribed' | 'interested' | 'not_interested'
  sequence_id?: number
  sequence_name?: string
  sequence_step?: string
  last_contacted?: number
  next_send?: number
  group_id?: number
  create_time?: number
}

export interface LeadListParams {
  page: number
  page_size: number
  keyword?: string
  status?: string
  group_id?: number
  sequence_id?: number
}

export function getLeadList(params: LeadListParams) {
  return instance.get('/leads/list', { params })
}

export function getLeadStats(params?: { group_id?: number }) {
  return instance.get('/leads/stats', { params })
}

export function importLeads(params: {
  file_data: string
  file_type: string
  group_id: number
  overwrite: number
}) {
  return instance.post('/leads/import', params, {
    fetchOptions: { successMessage: true },
  })
}

export function deleteLead(params: { id: number }) {
  return instance.post('/leads/delete', params, {
    fetchOptions: { successMessage: true },
  })
}

export function batchDeleteLeads(params: { ids: number[] }) {
  return instance.post('/leads/batch_delete', params, {
    fetchOptions: { successMessage: true },
  })
}

export function assignLeadsToSequence(params: { ids: number[]; sequence_id: number }) {
  return instance.post('/leads/assign_sequence', params, {
    fetchOptions: { successMessage: true },
  })
}

export function unsubscribeLead(params: { id: number }) {
  return instance.post('/leads/unsubscribe', params, {
    fetchOptions: { successMessage: true },
  })
}

export function verifyLeadEmail(params: { id: number }) {
  return instance.post('/leads/verify_email', params)
}

export function exportLeads(params: { group_id?: number; status?: string }) {
  return instance.get('/leads/export', { params })
}
