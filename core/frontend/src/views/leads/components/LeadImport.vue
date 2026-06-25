<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="Import Leads"
    style="width: 680px"
    :mask-closable="false">
    <div class="import-body">
      <n-steps :current="currentStep" size="small" class="mb-24px">
        <n-step title="Upload CSV" />
        <n-step title="Map Columns" />
        <n-step title="Review & Import" />
      </n-steps>

      <!-- Step 1: Upload -->
      <div v-if="currentStep === 1">
        <div
          class="drop-zone"
          :class="{ dragging: isDragging }"
          @dragover.prevent="isDragging = true"
          @dragleave="isDragging = false"
          @drop.prevent="handleDrop">
          <div class="drop-icon">📂</div>
          <div class="drop-title">Drop your CSV file here</div>
          <div class="drop-subtitle">or</div>
          <input ref="fileInput" type="file" accept=".csv" style="display:none" @change="handleFileChange" />
          <n-button @click="fileInput?.click()">Browse File</n-button>
          <div class="drop-hint">Supports CSV files up to 100MB · {{ leadCount.toLocaleString() }} max leads</div>
        </div>

        <div v-if="fileName" class="file-selected">
          <span class="file-icon">📄</span>
          <span class="file-name">{{ fileName }}</span>
          <span class="file-rows">{{ previewRows.length }} rows detected</span>
          <n-button text type="error" size="small" @click="clearFile">Remove</n-button>
        </div>

        <n-divider />
        <div class="paste-section">
          <div class="paste-label">Or paste CSV data directly:</div>
          <n-input
            v-model:value="pastedData"
            type="textarea"
            :rows="6"
            placeholder="email,first_name,last_name,company,title&#10;john@example.com,John,Smith,Acme Corp,CEO"
            @update:value="handlePaste" />
        </div>
      </div>

      <!-- Step 2: Map columns -->
      <div v-if="currentStep === 2">
        <div class="preview-header">
          <b>{{ previewRows.length }}</b> leads detected ·
          <span class="text-desc">{{ headers.length }} columns found</span>
        </div>

        <n-alert type="info" class="mb-16px" :show-icon="false">
          Map your CSV columns to lead fields. The <b>email</b> column is required.
        </n-alert>

        <div class="column-mapper">
          <div class="mapper-header">
            <span>CSV Column</span>
            <span>Sample Value</span>
            <span>Maps To</span>
          </div>
          <div v-for="(header, i) in headers" :key="i" class="mapper-row">
            <span class="col-name">{{ header }}</span>
            <span class="col-sample">{{ previewRows[0]?.[i] || '—' }}</span>
            <n-select
              v-model:value="columnMap[header]"
              :options="fieldOptions"
              size="small"
              placeholder="Skip" />
          </div>
        </div>

        <div class="variable-preview mt-16px">
          <div class="var-preview-title">Your custom variables (usable in templates):</div>
          <div class="var-chips">
            <span v-for="v in customVariables" :key="v" class="var-chip">{{ v }}</span>
          </div>
        </div>
      </div>

      <!-- Step 3: Review -->
      <div v-if="currentStep === 3">
        <div class="review-summary">
          <div class="summary-item">
            <span class="summary-icon">👥</span>
            <span class="summary-val">{{ previewRows.length.toLocaleString() }}</span>
            <span class="summary-lbl">Total Leads</span>
          </div>
          <div class="summary-item">
            <span class="summary-icon">✅</span>
            <span class="summary-val">{{ validCount.toLocaleString() }}</span>
            <span class="summary-lbl">Valid Emails</span>
          </div>
          <div class="summary-item">
            <span class="summary-icon">⚠️</span>
            <span class="summary-val">{{ previewRows.length - validCount }}</span>
            <span class="summary-lbl">Skipped</span>
          </div>
        </div>

        <n-form-item label="Add to Group" class="mt-16px">
          <n-input v-model:value="groupName" placeholder="e.g. SaaS Founders 2025" />
          <template #feedback>A new lead group will be created with this name</template>
        </n-form-item>

        <n-form-item label="Assign to Sequence (optional)">
          <n-select
            v-model:value="selectedSequenceId"
            :options="sequenceOptions"
            placeholder="Select a sequence..."
            clearable />
        </n-form-item>

        <n-form-item label="Overwrite Duplicates">
          <n-switch v-model:value="overwrite" />
          <span class="ml-12px text-desc text-12px">Update existing leads with same email</span>
        </n-form-item>

        <!-- Preview Table -->
        <n-data-table
          :columns="previewColumns"
          :data="previewRows.slice(0, 5).map(rowToLead)"
          size="small"
          class="mt-16px"
          :pagination="false" />
        <div class="text-desc text-12px mt-8px" v-if="previewRows.length > 5">
          Showing first 5 of {{ previewRows.length }} leads
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <n-button v-if="currentStep > 1" @click="currentStep--">Back</n-button>
        <div v-else></div>
        <div class="flex gap-8px">
          <n-button @click="$emit('update:show', false)">Cancel</n-button>
          <n-button
            v-if="currentStep < 3"
            type="primary"
            :disabled="!canProceed"
            @click="nextStep">
            Next →
          </n-button>
          <n-button
            v-else
            type="primary"
            :loading="importing"
            @click="confirmImport">
            Import {{ previewRows.length.toLocaleString() }} Leads
          </n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
import { NModal, NSteps, NStep, NButton, NInput, NSelect, NSwitch, NAlert, NDataTable, NFormItem, NDivider } from 'naive-ui'
import { Message } from '@/utils'
import { getSequenceList } from '@/api/modules/sequences'
import { importLeads } from '@/api/modules/leads'

const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ 'update:show': [v: boolean]; imported: [] }>()

const show = computed({
  get: () => props.show,
  set: (v) => emit('update:show', v),
})

const fileInput = ref<HTMLInputElement>()
const currentStep = ref(1)
const isDragging = ref(false)
const fileName = ref('')
const fileData = ref('')
const pastedData = ref('')
const headers = ref<string[]>([])
const previewRows = ref<string[][]>([])
const columnMap = ref<Record<string, string>>({})
const groupName = ref('')
const selectedSequenceId = ref<number | null>(null)
const overwrite = ref(false)
const importing = ref(false)
const sequenceOptions = ref<{ label: string; value: number }[]>([])
const leadCount = 50000

const fieldOptions = [
  { label: 'Email', value: 'email' },
  { label: 'First Name', value: 'first_name' },
  { label: 'Last Name', value: 'last_name' },
  { label: 'Company', value: 'company' },
  { label: 'Title/Role', value: 'title' },
  { label: 'Phone', value: 'phone' },
  { label: 'LinkedIn', value: 'linkedin' },
  { label: 'Website', value: 'website' },
  { label: 'Custom Variable', value: 'custom' },
  { label: 'Skip', value: '' },
]

const customVariables = computed(() =>
  headers.value
    .filter(h => columnMap.value[h] === 'custom' || (!columnMap.value[h] && h !== 'email'))
    .map(h => `{{${h.toLowerCase().replace(/\s+/g, '_')}}}`)
)

const validCount = computed(() => previewRows.value.filter(r => {
  const emailIdx = headers.value.findIndex(h => columnMap.value[h] === 'email' || h.toLowerCase() === 'email')
  return emailIdx >= 0 && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(r[emailIdx] || '')
}).length)

const canProceed = computed(() => {
  if (currentStep.value === 1) return previewRows.value.length > 0
  if (currentStep.value === 2) return headers.value.some(h => columnMap.value[h] === 'email' || h.toLowerCase() === 'email')
  return groupName.value.trim().length > 0
})

function parseCSV(text: string) {
  const lines = text.trim().split('\n').filter(Boolean)
  if (lines.length < 2) return
  headers.value = lines[0].split(',').map(h => h.trim().replace(/^"|"$/g, ''))
  previewRows.value = lines.slice(1).map(line =>
    line.split(',').map(c => c.trim().replace(/^"|"$/g, ''))
  )
  autoMapColumns()
}

function autoMapColumns() {
  const map: Record<string, string> = {}
  const knownFields: Record<string, string> = {
    email: 'email', 'e-mail': 'email',
    first_name: 'first_name', firstname: 'first_name', first: 'first_name',
    last_name: 'last_name', lastname: 'last_name', last: 'last_name',
    company: 'company', organization: 'company', org: 'company',
    title: 'title', role: 'title', position: 'title', jobtitle: 'title',
    phone: 'phone', telephone: 'phone', mobile: 'phone',
    linkedin: 'linkedin', linkedin_url: 'linkedin',
    website: 'website', url: 'website',
  }
  for (const h of headers.value) {
    const key = h.toLowerCase().replace(/[^a-z_]/g, '')
    map[h] = knownFields[key] || ''
  }
  columnMap.value = map
}

function handleFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  fileName.value = file.name
  const reader = new FileReader()
  reader.onload = (ev) => {
    fileData.value = ev.target?.result as string
    parseCSV(fileData.value)
  }
  reader.readAsText(file)
}

function handleDrop(e: DragEvent) {
  isDragging.value = false
  const file = e.dataTransfer?.files[0]
  if (!file || !file.name.endsWith('.csv')) {
    Message.error('Please drop a CSV file')
    return
  }
  fileName.value = file.name
  const reader = new FileReader()
  reader.onload = (ev) => {
    fileData.value = ev.target?.result as string
    parseCSV(fileData.value)
  }
  reader.readAsText(file)
}

function handlePaste(v: string) {
  if (v.includes(',') && v.includes('\n')) {
    fileData.value = v
    parseCSV(v)
  }
}

function clearFile() {
  fileName.value = ''
  fileData.value = ''
  previewRows.value = []
  headers.value = []
  if (fileInput.value) fileInput.value.value = ''
}

function nextStep() {
  currentStep.value++
  if (currentStep.value === 3) loadSequences()
}

async function loadSequences() {
  try {
    const res = await getSequenceList({ page: 1, page_size: 100 }) as any
    if (res?.list) sequenceOptions.value = res.list.map((s: any) => ({ label: s.name, value: s.id }))
  } catch (e) { void e }
}

function rowToLead(row: string[]) {
  const obj: Record<string, string> = {}
  headers.value.forEach((h, i) => { obj[columnMap.value[h] || h] = row[i] || '' })
  return obj
}

const previewColumns = computed(() =>
  headers.value.slice(0, 5).map(h => ({
    key: columnMap.value[h] || h,
    title: columnMap.value[h] || h,
    ellipsis: true,
  }))
)

async function confirmImport() {
  if (!groupName.value.trim()) {
    Message.error('Please enter a group name')
    return
  }
  importing.value = true
  try {
    await importLeads({
      file_data: fileData.value,
      file_type: 'csv',
      group_id: 0,
      overwrite: overwrite.value ? 1 : 0,
    })
    emit('imported')
    emit('update:show', false)
  } catch (_e) { /* ignore */ } finally {
    importing.value = false
    currentStep.value = 1
  }
}
</script>

<style scoped>
.import-body { min-height: 300px; }
.drop-zone {
  border: 2px dashed var(--n-border-color, #d9d9d9);
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  transition: all 0.2s;
  cursor: pointer;
}
.drop-zone.dragging { border-color: #1677ff; background: #e8f4fd; }
.drop-icon { font-size: 40px; margin-bottom: 12px; }
.drop-title { font-size: 16px; font-weight: 600; color: var(--n-text-color); margin-bottom: 8px; }
.drop-subtitle { color: var(--n-text-color-3); margin-bottom: 12px; }
.drop-hint { font-size: 12px; color: var(--n-text-color-3); margin-top: 12px; }
.file-selected {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 16px; background: #f6ffed; border: 1px solid #b7eb8f; border-radius: 8px;
  margin-top: 12px;
}
.file-icon { font-size: 18px; }
.file-name { font-weight: 500; flex: 1; }
.file-rows { font-size: 12px; color: #52c41a; }
.paste-section { margin-top: 8px; }
.paste-label { font-size: 13px; font-weight: 500; margin-bottom: 8px; color: var(--n-text-color-3); }
.preview-header { font-size: 13px; margin-bottom: 12px; color: var(--n-text-color); }
.column-mapper { border: 1px solid var(--n-border-color, #f0f0f0); border-radius: 8px; overflow: hidden; }
.mapper-header {
  display: grid; grid-template-columns: 1fr 1fr 1fr;
  padding: 8px 12px; background: var(--n-color-hover, #fafafa);
  font-size: 11px; font-weight: 600; color: var(--n-text-color-3); text-transform: uppercase;
}
.mapper-row {
  display: grid; grid-template-columns: 1fr 1fr 1fr;
  align-items: center; padding: 8px 12px; gap: 8px;
  border-top: 1px solid var(--n-border-color, #f0f0f0);
}
.col-name { font-weight: 500; font-size: 13px; }
.col-sample { font-size: 12px; color: var(--n-text-color-3); font-family: monospace; }
.variable-preview { }
.var-preview-title { font-size: 12px; font-weight: 600; margin-bottom: 8px; color: var(--n-text-color-3); }
.var-chips { display: flex; flex-wrap: wrap; gap: 6px; }
.var-chip { background: #e8f4fd; color: #1677ff; font-family: monospace; font-size: 11px; padding: 2px 8px; border-radius: 4px; }
.review-summary {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px;
  margin-bottom: 20px;
}
.summary-item {
  display: flex; flex-direction: column; align-items: center;
  padding: 16px; border: 1px solid var(--n-border-color, #f0f0f0); border-radius: 8px;
}
.summary-icon { font-size: 24px; margin-bottom: 6px; }
.summary-val { font-size: 24px; font-weight: 700; color: var(--n-text-color); }
.summary-lbl { font-size: 12px; color: var(--n-text-color-3); }
</style>
