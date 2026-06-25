import { instance } from '@/api'

export interface SequenceStep {
  id?: number
  type: 'initial' | 'followup1' | 'followup2' | 'reply'
  subject: string
  template_id: number
  addresser: string
  full_name: string
  delay_days: number
  ai_enabled: number
  ai_tone: string
}

export interface SequenceSchedule {
  working_days: string[]
  send_window_start: string
  send_window_end: string
  timezone: string
  daily_limit: number
  randomize_minutes: number
  stop_on_reply: number
  track_open: number
  track_click: number
}

export interface Sequence {
  id?: number
  name: string
  description: string
  status: 'draft' | 'active' | 'paused'
  steps: SequenceStep[]
  schedule: SequenceSchedule
  calendly_url: string
  leads_count?: number
  open_rate?: number
  reply_rate?: number
  create_time?: number
  update_time?: number
}

export function getSequenceList(params: { page: number; page_size: number; keyword?: string }) {
  return instance.get('/sequence/list', { params })
}

export function getSequenceById(params: { id: number }) {
  return instance.get('/sequence/find', { params })
}

export function createSequence(params: Sequence) {
  return instance.post('/sequence/create', params, {
    fetchOptions: { successMessage: true },
  })
}

export function updateSequence(params: Sequence & { id: number }) {
  return instance.post('/sequence/update', params, {
    fetchOptions: { successMessage: true },
  })
}

export function deleteSequence(params: { id: number }) {
  return instance.post('/sequence/delete', params, {
    fetchOptions: { successMessage: true },
  })
}

export function launchSequence(params: { sequence_id: number; lead_group_ids: number[] }) {
  return instance.post('/sequence/launch', params, {
    fetchOptions: { successMessage: true },
  })
}

export function pauseSequence(params: { id: number }) {
  return instance.post('/sequence/pause', params, {
    fetchOptions: { successMessage: true },
  })
}

export function resumeSequence(params: { id: number }) {
  return instance.post('/sequence/resume', params, {
    fetchOptions: { successMessage: true },
  })
}
