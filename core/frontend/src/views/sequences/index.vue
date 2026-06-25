<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div class="bt-title">Sequences</div>
      <n-button type="primary" @click="handleAdd">
        <template #icon><n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></n-icon></template>
        New Sequence
      </n-button>
    </div>

    <!-- Stats Row -->
    <div class="stats-row mb-24px">
      <div class="stat-card" v-for="stat in stats" :key="stat.label">
        <div class="stat-value">{{ stat.value }}</div>
        <div class="stat-label">{{ stat.label }}</div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <n-tabs v-model:value="activeTab" type="line" class="mb-16px" @update:value="fetchData">
      <n-tab name="all">All</n-tab>
      <n-tab name="active">Active</n-tab>
      <n-tab name="paused">Paused</n-tab>
      <n-tab name="draft">Draft</n-tab>
    </n-tabs>

    <!-- Sequence Cards -->
    <div v-if="loading" class="flex justify-center py-60px">
      <n-spin size="large" />
    </div>

    <div v-else-if="sequences.length === 0" class="empty-state">
      <div class="empty-icon">📧</div>
      <div class="empty-title">No sequences yet</div>
      <div class="empty-desc">Create your first cold outreach sequence to start sending</div>
      <n-button type="primary" class="mt-16px" @click="handleAdd">Create Sequence</n-button>
    </div>

    <div v-else class="sequence-grid">
      <div v-for="seq in sequences" :key="seq.id" class="sequence-card">
        <div class="seq-header">
          <div class="seq-info">
            <div class="seq-name">{{ seq.name }}</div>
            <div class="seq-desc">{{ seq.description || 'No description' }}</div>
          </div>
          <n-tag :type="statusType(seq.status)" size="small">{{ statusLabel(seq.status) }}</n-tag>
        </div>

        <div class="seq-steps">
          <div
            v-for="(step, i) in seq.steps"
            :key="i"
            class="step-chip"
            :class="`step-${step.type}`">
            <span class="step-num">{{ i + 1 }}</span>
            <span class="step-label">{{ stepLabel(step.type) }}</span>
            <span v-if="step.delay_days > 0" class="step-delay">+{{ step.delay_days }}d</span>
          </div>
          <div class="step-arrow" v-for="n in seq.steps.length - 1" :key="`a${n}`">→</div>
        </div>

        <div class="seq-metrics">
          <div class="metric">
            <span class="metric-val">{{ seq.leads_count || 0 }}</span>
            <span class="metric-lbl">Leads</span>
          </div>
          <div class="metric">
            <span class="metric-val">{{ seq.open_rate || 0 }}%</span>
            <span class="metric-lbl">Open Rate</span>
          </div>
          <div class="metric">
            <span class="metric-val">{{ seq.reply_rate || 0 }}%</span>
            <span class="metric-lbl">Reply Rate</span>
          </div>
        </div>

        <div class="seq-actions">
          <n-button size="small" @click="handleEdit(seq)">Edit</n-button>
          <n-button size="small" type="primary" @click="handleLaunch(seq)">Launch</n-button>
          <n-button
            size="small"
            :type="seq.status === 'active' ? 'warning' : 'success'"
            @click="handleToggle(seq)">
            {{ seq.status === 'active' ? 'Pause' : 'Resume' }}
          </n-button>
          <n-button size="small" type="error" @click="handleDelete(seq)">Delete</n-button>
        </div>
      </div>
    </div>

    <!-- Launch Modal -->
    <n-modal v-model:show="showLaunch" preset="card" title="Launch Sequence" style="width: 480px">
      <div v-if="selectedSeq">
        <p class="mb-16px text-desc">Select which lead groups to include in <b>{{ selectedSeq.name }}</b></p>
        <n-form-item label="Lead Groups">
          <n-select
            v-model:value="launchGroupIds"
            multiple
            :options="groupOptions"
            placeholder="Select lead groups..." />
        </n-form-item>
        <n-form-item label="Start Time">
          <n-radio-group v-model:value="launchNow">
            <n-radio :value="true">Send Now</n-radio>
            <n-radio :value="false">Schedule</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="!launchNow" label="Scheduled Time">
          <n-date-picker v-model:value="launchTime" type="datetime" />
        </n-form-item>
        <div class="flex justify-end gap-12px mt-16px">
          <n-button @click="showLaunch = false">Cancel</n-button>
          <n-button type="primary" :loading="launching" @click="confirmLaunch">Launch 🚀</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { NButton, NIcon, NModal, NFormItem, NSelect, NDatePicker, NRadioGroup, NRadio, NTag, NTabs, NTab, NSpin } from 'naive-ui'
import { confirm, Message } from '@/utils'
import { getSequenceList, deleteSequence, pauseSequence, resumeSequence, launchSequence } from '@/api/modules/sequences'
import type { Sequence } from '@/api/modules/sequences'
import { getGroupList } from '@/api/modules/contacts/group'

const router = useRouter()

const loading = ref(false)
const sequences = ref<Sequence[]>([])
const activeTab = ref('all')

const showLaunch = ref(false)
const selectedSeq = ref<Sequence | null>(null)
const launchGroupIds = ref<number[]>([])
const launchNow = ref(true)
const launchTime = ref<number | null>(null)
const launching = ref(false)
const groupOptions = ref<{ label: string; value: number }[]>([])

const stats = computed(() => [
  { label: 'Total Sequences', value: sequences.value.length },
  { label: 'Active', value: sequences.value.filter(s => s.status === 'active').length },
  { label: 'Total Leads', value: sequences.value.reduce((a, s) => a + (s.leads_count || 0), 0) },
  { label: 'Avg Reply Rate', value: sequences.value.length ? Math.round(sequences.value.reduce((a, s) => a + (s.reply_rate || 0), 0) / sequences.value.length) + '%' : '0%' },
])

const statusType = (status: string) => ({ active: 'success', paused: 'warning', draft: 'default' }[status] || 'default') as any
const statusLabel = (status: string) => ({ active: 'Active', paused: 'Paused', draft: 'Draft' }[status] || status)
const stepLabel = (type: string) => ({ initial: 'Initial', followup1: 'Follow-up 1', followup2: 'Follow-up 2', reply: 'Reply' }[type] || type)

async function fetchData() {
  loading.value = true
  try {
    const res = await getSequenceList({ page: 1, page_size: 100, keyword: '' }) as any
    if (res?.list) sequences.value = res.list
  } catch (_e) {
    sequences.value = []
  } finally {
    loading.value = false
  }
}

async function fetchGroups() {
  try {
    const res = await getGroupList({ page: 1, page_size: 100, keyword: '' }) as any
    if (res?.list) groupOptions.value = res.list.map((g: any) => ({ label: g.name, value: g.id }))
  } catch (e) { void e }
}

function handleAdd() {
  router.push('/sequences/edit')
}

function handleEdit(seq: Sequence) {
  router.push(`/sequences/edit?id=${seq.id}`)
}

function handleLaunch(seq: Sequence) {
  selectedSeq.value = seq
  launchGroupIds.value = []
  launchNow.value = true
  showLaunch.value = true
}

async function confirmLaunch() {
  if (!selectedSeq.value?.id) return
  if (!launchGroupIds.value.length) {
    Message.warning('Select at least one lead group')
    return
  }
  launching.value = true
  try {
    await launchSequence({ sequence_id: selectedSeq.value.id, lead_group_ids: launchGroupIds.value })
    showLaunch.value = false
    fetchData()
  } catch (_e) { /* ignore */ } finally {
    launching.value = false
  }
}

async function handleToggle(seq: Sequence) {
  if (!seq.id) return
  try {
    if (seq.status === 'active') {
      await pauseSequence({ id: seq.id })
    } else {
      await resumeSequence({ id: seq.id })
    }
    fetchData()
  } catch (e) { void e }
}

async function handleDelete(seq: Sequence) {
  if (!seq.id) return
  await confirm(`Delete sequence "${seq.name}"? This cannot be undone.`)
  try {
    await deleteSequence({ id: seq.id })
    fetchData()
  } catch (e) { void e }
}

onMounted(() => {
  fetchData()
  fetchGroups()
})
</script>

<style scoped>
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.stat-card {
  background: var(--n-card-color, #fff);
  border: 1px solid var(--n-border-color, #e5e7eb);
  border-radius: 8px;
  padding: 16px 20px;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--n-text-color);
  line-height: 1.2;
}
.stat-label {
  font-size: 12px;
  color: var(--n-text-color-3);
  margin-top: 4px;
}
.sequence-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(420px, 1fr));
  gap: 16px;
}
.sequence-card {
  background: var(--n-card-color, #fff);
  border: 1px solid var(--n-border-color, #e5e7eb);
  border-radius: 12px;
  padding: 20px;
  transition: box-shadow 0.2s;
}
.sequence-card:hover {
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}
.seq-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}
.seq-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--n-text-color);
}
.seq-desc {
  font-size: 12px;
  color: var(--n-text-color-3);
  margin-top: 4px;
}
.seq-steps {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 16px;
}
.step-chip {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 500;
  background: #f0f0f0;
}
.step-chip.step-initial { background: #e8f4fd; color: #1677ff; }
.step-chip.step-followup1 { background: #fff7e6; color: #fa8c16; }
.step-chip.step-followup2 { background: #fff0f6; color: #eb2f96; }
.step-chip.step-reply { background: #f6ffed; color: #52c41a; }
.step-num { font-weight: 700; }
.step-delay { font-size: 10px; opacity: 0.8; }
.step-arrow { color: var(--n-text-color-3); font-size: 12px; }
.seq-metrics {
  display: flex;
  gap: 20px;
  padding: 12px 0;
  border-top: 1px solid var(--n-border-color, #f0f0f0);
  border-bottom: 1px solid var(--n-border-color, #f0f0f0);
  margin-bottom: 14px;
}
.metric {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.metric-val {
  font-size: 18px;
  font-weight: 600;
  color: var(--n-text-color);
}
.metric-lbl {
  font-size: 11px;
  color: var(--n-text-color-3);
}
.seq-actions {
  display: flex;
  gap: 8px;
}
.empty-state {
  text-align: center;
  padding: 80px 0;
  color: var(--n-text-color-3);
}
.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-title { font-size: 18px; font-weight: 600; color: var(--n-text-color); }
.empty-desc { font-size: 14px; margin-top: 8px; }
</style>
