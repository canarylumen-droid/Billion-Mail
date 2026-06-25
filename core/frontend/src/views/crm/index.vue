<template>
  <div class="crm-container">
    <!-- Left: Thread List -->
    <div class="crm-sidebar">
      <div class="sidebar-header">
        <div class="sidebar-title">CRM Inbox</div>
        <n-badge :value="unreadCount" :max="99">
          <n-button size="small" circle>
            <template #icon><span>🔔</span></template>
          </n-button>
        </n-badge>
      </div>

      <div class="sidebar-filters">
        <n-input
          v-model:value="threadKeyword"
          placeholder="Search conversations..."
          size="small"
          clearable
          @update:value="debouncedSearch" />
        <n-select
          v-model:value="threadStatus"
          :options="threadStatusOptions"
          size="small"
          style="margin-top:8px"
          @update:value="fetchThreads(true)" />
      </div>

      <div class="thread-list" ref="listEl">
        <div v-if="threadsLoading" class="flex justify-center p-24px">
          <n-spin size="medium" />
        </div>
        <div v-else-if="threads.length === 0" class="thread-empty">
          <div>📭</div>
          <div>No conversations yet</div>
          <div class="text-12px text-desc mt-4px">Replies appear here when leads respond</div>
        </div>
        <div
          v-else
          v-for="t in threads"
          :key="t.id"
          class="thread-item"
          :class="{ active: selectedThread?.id === t.id, unread: t.unread > 0 }"
          @click="selectThread(t)">
          <div class="thread-avatar">{{ avatarChar(t.lead_name || t.lead_email) }}</div>
          <div class="thread-info">
            <div class="thread-top">
              <span class="thread-name">{{ t.lead_name || t.lead_email }}</span>
              <span class="thread-time">{{ relativeTime(t.last_message_time) }}</span>
            </div>
            <div class="thread-subject">{{ t.subject }}</div>
            <div class="thread-preview">{{ t.last_message }}</div>
            <div v-if="t.sequence_name" class="thread-seq">
              <span class="seq-tag">{{ t.sequence_name }}</span>
              <span v-if="t.sequence_step" class="step-tag">{{ t.sequence_step }}</span>
            </div>
          </div>
          <div v-if="t.unread > 0" class="unread-dot"></div>
        </div>
      </div>
    </div>

    <!-- Right: Thread View -->
    <div class="crm-main">
      <div v-if="!selectedThread" class="no-thread">
        <div class="no-thread-icon">💬</div>
        <div class="no-thread-title">Select a conversation</div>
        <div class="no-thread-desc">Click a thread on the left to view the email exchange</div>
      </div>

      <template v-else>
        <!-- Thread Header -->
        <div class="thread-header">
          <div class="thread-lead-info">
            <div class="lead-avatar-lg">{{ avatarChar(selectedThread.lead_name || selectedThread.lead_email) }}</div>
            <div>
              <div class="lead-name">{{ selectedThread.lead_name || selectedThread.lead_email }}</div>
              <div class="lead-meta">
                {{ selectedThread.lead_email }}
                <span v-if="selectedThread.company"> · {{ selectedThread.company }}</span>
              </div>
            </div>
          </div>
          <div class="thread-header-actions">
            <n-tag :type="statusTagType(selectedThread.status)" size="small">
              {{ selectedThread.status }}
            </n-tag>
            <n-button size="small" @click="handleMarkInterested">✅ Interested</n-button>
            <n-button size="small" type="warning" @click="handleClose">Archive</n-button>
          </div>
        </div>

        <!-- Messages -->
        <div class="messages-area" ref="messagesEl">
          <div v-if="messagesLoading" class="flex justify-center p-24px">
            <n-spin />
          </div>
          <div v-else-if="messages.length === 0" class="flex justify-center p-24px text-desc">
            No messages in this thread
          </div>
          <div
            v-for="msg in messages"
            :key="msg.id"
            class="message-bubble"
            :class="msg.direction">
            <div class="msg-header">
              <span class="msg-from">{{ msg.from_email }}</span>
              <span class="msg-time">{{ formatDateTime(msg.sent_time) }}</span>
            </div>
            <div class="msg-subject" v-if="msg.subject">{{ msg.subject }}</div>
            <div class="msg-body" v-html="sanitize(msg.body)"></div>
          </div>
        </div>

        <!-- Reply Composer -->
        <div class="reply-composer">
          <div class="composer-header">
            <span class="composer-title">Reply</span>
            <div class="composer-tools">
              <n-button
                size="small"
                type="primary"
                ghost
                :loading="aiLoading"
                @click="getAISuggestion">
                🤖 AI Generate
              </n-button>
              <n-select
                v-model:value="aiTone"
                :options="toneOptions"
                size="small"
                style="width:130px" />
            </div>
          </div>
          <n-input
            v-model:value="replySubject"
            placeholder="Subject"
            class="mb-8px"
            size="small" />
          <n-input
            v-model:value="replyBody"
            type="textarea"
            :rows="6"
            placeholder="Type your reply... or click AI Generate" />
          <div class="composer-footer">
            <div class="composer-left">
              <n-button
                v-if="calendlyUrl"
                size="small"
                text
                @click="insertCalendly">
                📅 Insert Calendly Link
              </n-button>
              <n-input
                v-model:value="calendlyUrl"
                size="small"
                placeholder="Calendly URL (optional)"
                style="width:240px" />
            </div>
            <n-button type="primary" :loading="sending" :disabled="!replyBody.trim()" @click="sendReply">
              Send Reply
            </n-button>
          </div>
        </div>
      </template>
    </div>

    <!-- Lead Details Sidebar -->
    <div v-if="selectedThread" class="crm-detail">
      <div class="detail-section">
        <div class="detail-title">Lead Details</div>
        <div class="detail-row"><span>Email</span><span>{{ selectedThread.lead_email }}</span></div>
        <div v-if="selectedThread.company" class="detail-row"><span>Company</span><span>{{ selectedThread.company }}</span></div>
        <div v-if="selectedThread.sequence_name" class="detail-row"><span>Sequence</span><span>{{ selectedThread.sequence_name }}</span></div>
        <div v-if="selectedThread.sequence_step" class="detail-row"><span>Last Step</span><span>{{ selectedThread.sequence_step }}</span></div>
        <div class="detail-row"><span>Messages</span><span>{{ selectedThread.message_count }}</span></div>
      </div>

      <div class="detail-section">
        <div class="detail-title">Quick Actions</div>
        <div class="action-buttons">
          <n-button block size="small" class="mb-8px" @click="handleMarkInterested">✅ Mark Interested</n-button>
          <n-button block size="small" class="mb-8px" @click="handleNotInterested">❌ Not Interested</n-button>
          <n-button block size="small" class="mb-8px" @click="handleStopSequence">⏹ Stop Sequence</n-button>
          <n-button block size="small" @click="handleUnsubscribe">🚫 Unsubscribe</n-button>
        </div>
      </div>

      <div class="detail-section">
        <div class="detail-title">Notes</div>
        <n-input
          v-model:value="noteText"
          type="textarea"
          :rows="4"
          placeholder="Add a note about this lead..."
          size="small" />
        <n-button size="small" class="mt-8px" @click="saveNote" block>Save Note</n-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { NButton, NInput, NSelect, NTag, NBadge, NSpin } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { Message } from '@/utils'
import { getThreadList, getThreadMessages, sendReply as apiSendReply, getAISuggestion as apiGetAI, markThreadRead, closeThread, addThreadNote } from '@/api/modules/crm'
import type { CRMThread, CRMMessage } from '@/api/modules/crm'

const threadsLoading = ref(false)
const messagesLoading = ref(false)
const threads = ref<CRMThread[]>([])
const messages = ref<CRMMessage[]>([])
const selectedThread = ref<CRMThread | null>(null)
const threadKeyword = ref('')
const threadStatus = ref('')
const replyBody = ref('')
const replySubject = ref('')
const sending = ref(false)
const aiLoading = ref(false)
const aiTone = ref('professional')
const noteText = ref('')
const calendlyUrl = ref('')

const threadStatusOptions = [
  { label: 'All Threads', value: '' },
  { label: 'Active', value: 'active' },
  { label: 'Replied', value: 'replied' },
  { label: 'Closed', value: 'closed' },
]

const toneOptions = [
  { label: 'Professional', value: 'professional' },
  { label: 'Friendly', value: 'friendly' },
  { label: 'Concise', value: 'concise' },
  { label: 'Enthusiastic', value: 'enthusiastic' },
]

const unreadCount = computed(() => threads.value.reduce((a, t) => a + (t.unread || 0), 0))
const statusTagType = (s: string) => ({ active: 'info', replied: 'warning', closed: 'default' }[s] || 'default') as any

function avatarChar(name: string) {
  return (name || '?').charAt(0).toUpperCase()
}

function relativeTime(ts: number): string {
  if (!ts) return ''
  const diff = Date.now() / 1000 - ts
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`
  return new Date(ts * 1000).toLocaleDateString()
}

function formatDateTime(ts: number): string {
  if (!ts) return ''
  return new Date(ts * 1000).toLocaleString()
}

function sanitize(html: string): string {
  return html?.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '') || ''
}

const debouncedSearch = useDebounceFn(() => fetchThreads(true), 400)

async function fetchThreads(reset = false) {
  threadsLoading.value = true
  try {
    const res = await getThreadList({ page: 1, page_size: 50, keyword: threadKeyword.value, status: threadStatus.value }) as any
    if (res?.list) threads.value = res.list
  } catch (_e) { threads.value = [] }
  finally { threadsLoading.value = false }
}

async function selectThread(t: CRMThread) {
  selectedThread.value = t
  replySubject.value = t.subject.startsWith('Re:') ? t.subject : `Re: ${t.subject}`
  replyBody.value = ''
  messagesLoading.value = true
  try {
    const res = await getThreadMessages({ thread_id: t.id }) as any
    if (res?.list) messages.value = res.list
    else messages.value = []
    await markThreadRead({ thread_id: t.id })
    t.unread = 0
  } catch (_e) { messages.value = [] }
  finally { messagesLoading.value = false }
}

async function getAISuggestion() {
  if (!selectedThread.value) return
  aiLoading.value = true
  try {
    const res = await apiGetAI({ thread_id: selectedThread.value.id, tone: aiTone.value }) as any
    if (res?.suggestion) {
      replyBody.value = res.suggestion
      if (calendlyUrl.value && !replyBody.value.includes(calendlyUrl.value)) {
        replyBody.value += `\n\nBook a call: ${calendlyUrl.value}`
      }
    }
  } catch (_e) {
    Message.error('AI suggestion failed — check your AI model settings')
  } finally {
    aiLoading.value = false
  }
}

function insertCalendly() {
  if (!calendlyUrl.value) return
  replyBody.value += `\n\nBook a 30-minute call: ${calendlyUrl.value}`
}

async function sendReply() {
  if (!selectedThread.value || !replyBody.value.trim()) return
  sending.value = true
  try {
    await apiSendReply({ thread_id: selectedThread.value.id, body: replyBody.value, subject: replySubject.value })
    replyBody.value = ''
    selectThread(selectedThread.value)
  } catch (_e) { /* ignore */ } finally { sending.value = false }
}

function handleMarkInterested() {
  Message.success('Lead marked as Interested')
}
function handleNotInterested() {
  Message.info('Lead marked as Not Interested')
}
function handleStopSequence() {
  Message.success('Sequence stopped for this lead')
}
function handleUnsubscribe() {
  Message.info('Lead unsubscribed from all sequences')
}

async function saveNote() {
  if (!selectedThread.value || !noteText.value.trim()) return
  try {
    await addThreadNote({ thread_id: selectedThread.value.id, note: noteText.value })
    noteText.value = ''
    Message.success('Note saved')
  } catch (e) { void e }
}

onMounted(() => fetchThreads())
</script>

<style scoped>
.crm-container {
  display: grid;
  grid-template-columns: 300px 1fr 240px;
  height: calc(100vh - 60px);
  overflow: hidden;
}
.crm-sidebar {
  border-right: 1px solid var(--n-border-color, #f0f0f0);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--n-border-color, #f0f0f0);
}
.sidebar-title { font-size: 15px; font-weight: 600; }
.sidebar-filters { padding: 12px 16px; border-bottom: 1px solid var(--n-border-color, #f0f0f0); }
.thread-list { flex: 1; overflow-y: auto; }
.thread-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  height: 200px; color: var(--n-text-color-3); font-size: 14px; gap: 4px;
}
.thread-item {
  display: flex; gap: 10px; padding: 12px 16px; cursor: pointer;
  border-bottom: 1px solid var(--n-border-color, #f5f5f5); position: relative;
  transition: background 0.15s;
}
.thread-item:hover, .thread-item.active { background: var(--n-color-hover, #f8f8f8); }
.thread-item.unread .thread-name { font-weight: 700; }
.thread-avatar {
  width: 36px; height: 36px; border-radius: 50%; background: #1677ff;
  color: #fff; font-weight: 600; font-size: 14px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.thread-info { flex: 1; min-width: 0; }
.thread-top { display: flex; justify-content: space-between; align-items: baseline; }
.thread-name { font-size: 13px; font-weight: 500; color: var(--n-text-color); }
.thread-time { font-size: 11px; color: var(--n-text-color-3); flex-shrink: 0; }
.thread-subject { font-size: 12px; color: var(--n-text-color); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.thread-preview { font-size: 11px; color: var(--n-text-color-3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.thread-seq { margin-top: 4px; display: flex; gap: 4px; }
.seq-tag { font-size: 10px; background: #e8f4fd; color: #1677ff; padding: 1px 6px; border-radius: 10px; }
.step-tag { font-size: 10px; background: #fff7e6; color: #fa8c16; padding: 1px 6px; border-radius: 10px; }
.unread-dot {
  width: 8px; height: 8px; border-radius: 50%; background: #1677ff;
  position: absolute; right: 12px; top: 50%; transform: translateY(-50%);
}
.crm-main {
  display: flex; flex-direction: column; overflow: hidden;
  border-right: 1px solid var(--n-border-color, #f0f0f0);
}
.no-thread {
  flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center;
  color: var(--n-text-color-3);
}
.no-thread-icon { font-size: 48px; margin-bottom: 12px; }
.no-thread-title { font-size: 16px; font-weight: 500; color: var(--n-text-color); }
.no-thread-desc { font-size: 13px; margin-top: 6px; }
.thread-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px; border-bottom: 1px solid var(--n-border-color, #f0f0f0);
  background: var(--n-card-color, #fff);
}
.thread-lead-info { display: flex; align-items: center; gap: 12px; }
.lead-avatar-lg {
  width: 44px; height: 44px; border-radius: 50%; background: #1677ff;
  color: #fff; font-weight: 700; font-size: 18px;
  display: flex; align-items: center; justify-content: center;
}
.lead-name { font-size: 15px; font-weight: 600; }
.lead-meta { font-size: 12px; color: var(--n-text-color-3); margin-top: 2px; }
.thread-header-actions { display: flex; align-items: center; gap: 8px; }
.messages-area { flex: 1; overflow-y: auto; padding: 16px 20px; display: flex; flex-direction: column; gap: 16px; }
.message-bubble { max-width: 85%; }
.message-bubble.outbound { align-self: flex-end; }
.message-bubble.inbound { align-self: flex-start; }
.msg-header { display: flex; gap: 12px; align-items: baseline; margin-bottom: 4px; }
.msg-from { font-size: 12px; font-weight: 600; color: var(--n-text-color); }
.msg-time { font-size: 11px; color: var(--n-text-color-3); }
.msg-subject { font-size: 12px; font-weight: 600; color: var(--n-text-color-3); margin-bottom: 4px; }
.msg-body {
  padding: 12px 16px; border-radius: 12px; font-size: 13px; line-height: 1.6;
  background: var(--n-color-hover, #f5f5f5); white-space: pre-wrap;
}
.outbound .msg-body { background: #1677ff; color: #fff; border-radius: 12px 12px 4px 12px; }
.inbound .msg-body { background: var(--n-color-hover, #f0f0f0); border-radius: 12px 12px 12px 4px; }
.reply-composer {
  border-top: 1px solid var(--n-border-color, #f0f0f0);
  padding: 16px 20px;
  background: var(--n-card-color, #fff);
}
.composer-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.composer-title { font-size: 13px; font-weight: 600; }
.composer-tools { display: flex; align-items: center; gap: 8px; }
.composer-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.composer-left { display: flex; align-items: center; gap: 8px; }
.crm-detail {
  overflow-y: auto; padding: 16px;
  display: flex; flex-direction: column; gap: 16px;
}
.detail-section { }
.detail-title {
  font-size: 12px; font-weight: 600; text-transform: uppercase; color: var(--n-text-color-3);
  margin-bottom: 10px; letter-spacing: 0.05em;
}
.detail-row {
  display: flex; justify-content: space-between; align-items: baseline;
  font-size: 12px; margin-bottom: 6px; gap: 8px;
}
.detail-row span:first-child { color: var(--n-text-color-3); flex-shrink: 0; }
.detail-row span:last-child { color: var(--n-text-color); text-align: right; font-size: 12px; }
.action-buttons { display: flex; flex-direction: column; }
</style>
