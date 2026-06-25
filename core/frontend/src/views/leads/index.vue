<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div class="bt-title">Leads</div>
      <div class="flex gap-8px">
        <n-button @click="handleExport">Export CSV</n-button>
        <n-button type="primary" @click="showImport = true">
          + Import Leads
        </n-button>
      </div>
    </div>

    <!-- Stats -->
    <div class="leads-stats mb-24px">
      <div v-for="s in statusStats" :key="s.key" class="stat-chip" :class="`chip-${s.key}`" @click="filterByStatus(s.key)">
        <span class="chip-val">{{ s.count.toLocaleString() }}</span>
        <span class="chip-label">{{ s.label }}</span>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="flex items-center gap-12px mb-16px">
      <n-input
        v-model:value="keyword"
        placeholder="Search by email, name, company..."
        style="max-width: 320px"
        clearable
        @update:value="debouncedFetch" />
      <n-select
        v-model:value="statusFilter"
        :options="statusOptions"
        style="width: 160px"
        @update:value="fetchLeads(true)" />
      <n-select
        v-model:value="sequenceFilter"
        :options="[{ label: 'All Sequences', value: '' }, ...sequenceOptions]"
        style="width: 200px"
        @update:value="fetchLeads(true)" />
      <span v-if="selectedIds.length" class="text-desc text-13px ml-8px">
        {{ selectedIds.length }} selected
      </span>
      <n-button v-if="selectedIds.length" size="small" @click="handleBatchAssign">
        Assign to Sequence
      </n-button>
      <n-button v-if="selectedIds.length" size="small" type="error" @click="handleBatchDelete">
        Delete Selected
      </n-button>
    </div>

    <!-- Table -->
    <n-data-table
      v-model:checked-row-keys="selectedIds"
      :loading="loading"
      :columns="columns"
      :data="leads"
      :row-key="row => row.id"
      :pagination="pagination"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange">
      <template #empty>
        <div class="empty-leads">
          <div class="empty-icon">📋</div>
          <div class="empty-title">No leads yet</div>
          <div class="empty-desc">Import a CSV file to get started</div>
          <n-button type="primary" class="mt-12px" @click="showImport = true">Import Leads</n-button>
        </div>
      </template>
    </n-data-table>

    <!-- Import Modal -->
    <lead-import v-model:show="showImport" @imported="fetchLeads(true)" />

    <!-- Assign Sequence Modal -->
    <n-modal v-model:show="showAssign" preset="card" title="Assign to Sequence" style="width: 400px">
      <n-form-item label="Select Sequence">
        <n-select v-model:value="assignSequenceId" :options="sequenceOptions" placeholder="Choose a sequence..." />
      </n-form-item>
      <div class="flex justify-end gap-8px mt-16px">
        <n-button @click="showAssign = false">Cancel</n-button>
        <n-button type="primary" :loading="assigning" @click="confirmAssign">Assign</n-button>
      </div>
    </n-modal>
  </div>
</template>

<script lang="tsx" setup>
import { NButton, NInput, NSelect, NDataTable, NModal, NFormItem, NTag } from 'naive-ui'
import { useDebounceFn } from '@vueuse/core'
import { confirm, formatTime, Message } from '@/utils'
import { getLeadList, deleteLead, batchDeleteLeads, assignLeadsToSequence, type Lead } from '@/api/modules/leads'
import { getSequenceList } from '@/api/modules/sequences'
import LeadImport from './components/LeadImport.vue'

const loading = ref(false)
const leads = ref<Lead[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const statusFilter = ref('')
const sequenceFilter = ref('')
const selectedIds = ref<number[]>([])
const showImport = ref(false)
const showAssign = ref(false)
const assigning = ref(false)
const assignSequenceId = ref<number | null>(null)
const sequenceOptions = ref<{ label: string; value: number }[]>([])

const pagination = computed(() => ({
  page: page.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
}))

const statusStats = ref([
  { key: '', label: 'All Leads', count: 0 },
  { key: 'not_started', label: 'Not Started', count: 0 },
  { key: 'in_sequence', label: 'In Sequence', count: 0 },
  { key: 'replied', label: 'Replied', count: 0 },
  { key: 'interested', label: 'Interested', count: 0 },
  { key: 'bounced', label: 'Bounced', count: 0 },
])

const statusOptions = [
  { label: 'All Status', value: '' },
  { label: 'Not Started', value: 'not_started' },
  { label: 'In Sequence', value: 'in_sequence' },
  { label: 'Replied', value: 'replied' },
  { label: 'Interested', value: 'interested' },
  { label: 'Not Interested', value: 'not_interested' },
  { label: 'Bounced', value: 'bounced' },
  { label: 'Unsubscribed', value: 'unsubscribed' },
]

const statusTagType = (s: string) => ({
  not_started: 'default',
  in_sequence: 'info',
  replied: 'warning',
  interested: 'success',
  not_interested: 'default',
  bounced: 'error',
  unsubscribed: 'default',
}[s] || 'default') as any

const columns = [
  { type: 'selection' as const },
  {
    key: 'email', title: 'Email', minWidth: 200,
    render: (row: Lead) => (
      <div>
        <div style="font-weight:500">{row.email}</div>
        {row.first_name && <div style="font-size:12px;color:var(--n-text-color-3)">{row.first_name} {row.last_name}</div>}
      </div>
    ),
  },
  { key: 'company', title: 'Company', minWidth: 130, ellipsis: { tooltip: true } },
  { key: 'title', title: 'Title', minWidth: 120, ellipsis: { tooltip: true } },
  {
    key: 'status', title: 'Status', width: 120,
    render: (row: Lead) => (
      <NTag type={statusTagType(row.status)} size="small">
        {statusOptions.find(s => s.value === row.status)?.label || row.status}
      </NTag>
    ),
  },
  {
    key: 'sequence_name', title: 'Sequence', minWidth: 140, ellipsis: { tooltip: true },
    render: (row: Lead) => row.sequence_name || <span style="color:var(--n-text-color-3)">—</span>,
  },
  {
    key: 'sequence_step', title: 'Step', width: 110,
    render: (row: Lead) => row.sequence_step || <span style="color:var(--n-text-color-3)">—</span>,
  },
  {
    key: 'last_contacted', title: 'Last Contacted', width: 140,
    render: (row: Lead) => row.last_contacted ? formatTime(row.last_contacted) : <span style="color:var(--n-text-color-3)">—</span>,
  },
  {
    key: 'actions', title: '', width: 80, fixed: 'right' as const,
    render: (row: Lead) => (
      <NButton text type="error" size="small" onClick={() => handleDelete(row)}>Remove</NButton>
    ),
  },
]

const debouncedFetch = useDebounceFn(() => fetchLeads(true), 400)

async function fetchLeads(reset = false) {
  if (reset) page.value = 1
  loading.value = true
  try {
    const res = await getLeadList({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      status: statusFilter.value || undefined,
      sequence_id: sequenceFilter.value ? Number(sequenceFilter.value) : undefined,
    }) as any
    if (res?.list) {
      leads.value = res.list
      total.value = res.total || 0
    }
  } catch (_e) {
    leads.value = []
  } finally {
    loading.value = false
  }
}

function filterByStatus(status: string) {
  statusFilter.value = status
  fetchLeads(true)
}

async function loadSequences() {
  try {
    const res = await getSequenceList({ page: 1, page_size: 100 }) as any
    if (res?.list) sequenceOptions.value = res.list.map((s: any) => ({ label: s.name, value: s.id }))
  } catch (e) { void e }
}

async function handleDelete(row: Lead) {
  await confirm(`Remove ${row.email} from leads?`)
  try {
    if (row.id) await deleteLead({ id: row.id })
    fetchLeads()
  } catch (e) { void e }
}

async function handleBatchDelete() {
  if (!selectedIds.value.length) return
  await confirm(`Delete ${selectedIds.value.length} leads?`)
  try {
    await batchDeleteLeads({ ids: selectedIds.value as number[] })
    selectedIds.value = []
    fetchLeads()
  } catch (e) { void e }
}

function handleBatchAssign() {
  showAssign.value = true
}

async function confirmAssign() {
  if (!assignSequenceId.value) {
    Message.warning('Select a sequence')
    return
  }
  assigning.value = true
  try {
    await assignLeadsToSequence({ ids: selectedIds.value as number[], sequence_id: assignSequenceId.value })
    showAssign.value = false
    selectedIds.value = []
    fetchLeads()
  } catch (_e) { /* ignore */ } finally {
    assigning.value = false
  }
}

function handleExport() {
  Message.info('Export feature requires backend deployment')
}

function handlePageChange(p: number) { page.value = p; fetchLeads() }
function handlePageSizeChange(ps: number) { pageSize.value = ps; fetchLeads(true) }

onMounted(() => {
  fetchLeads()
  loadSequences()
})
</script>

<style scoped>
.leads-stats {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.stat-chip {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 20px;
  border: 1px solid var(--n-border-color, #e5e7eb);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  min-width: 90px;
}
.stat-chip:hover { border-color: #1677ff; background: #e8f4fd; }
.chip-val { font-size: 20px; font-weight: 700; color: var(--n-text-color); }
.chip-label { font-size: 11px; color: var(--n-text-color-3); margin-top: 2px; }
.empty-leads { text-align: center; padding: 60px 0; }
.empty-icon { font-size: 40px; margin-bottom: 12px; }
.empty-title { font-size: 16px; font-weight: 600; color: var(--n-text-color); }
.empty-desc { font-size: 13px; color: var(--n-text-color-3); margin-top: 6px; }
</style>
