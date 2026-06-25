import { instance } from '@/api'

export interface CRMThread {
  id: number
  lead_email: string
  lead_name: string
  company?: string
  subject: string
  last_message: string
  last_message_time: number
  unread: number
  status: 'active' | 'replied' | 'closed'
  sequence_name?: string
  sequence_step?: string
  message_count: number
}

export interface CRMMessage {
  id: number
  thread_id: number
  from_email: string
  to_email: string
  subject: string
  body: string
  sent_time: number
  direction: 'outbound' | 'inbound'
  is_read: number
}

export function getThreadList(params: {
  page: number
  page_size: number
  keyword?: string
  status?: string
}) {
  return instance.get('/crm/threads', { params })
}

export function getThreadMessages(params: { thread_id: number }) {
  return instance.get('/crm/thread/messages', { params })
}

export function sendReply(params: {
  thread_id: number
  body: string
  subject: string
}) {
  return instance.post('/crm/reply', params, {
    fetchOptions: { successMessage: true },
  })
}

export function getAISuggestion(params: {
  thread_id: number
  tone?: string
  context?: string
}) {
  return instance.post('/crm/ai_suggest', params)
}

export function markThreadRead(params: { thread_id: number }) {
  return instance.post('/crm/mark_read', params)
}

export function closeThread(params: { thread_id: number }) {
  return instance.post('/crm/close', params, {
    fetchOptions: { successMessage: true },
  })
}

export function addThreadNote(params: { thread_id: number; note: string }) {
  return instance.post('/crm/note', params, {
    fetchOptions: { successMessage: true },
  })
}
