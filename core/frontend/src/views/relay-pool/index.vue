<template>
  <div class="p-24px">
    <!-- Header -->
    <div class="flex items-center justify-between mb-20px">
      <div>
        <div class="bt-title">Relay Providers</div>
        <div class="text-desc text-13px mt-4px">15 providers × 2 accounts = 30 API slots. Auto-routes sends through available providers.</div>
      </div>
      <div class="flex gap-8px">
        <n-button :loading="testingPool" @click="showTestPool = true">
          Test Pool Auto-Route
        </n-button>
        <n-button @click="handleResetAll">Reset Daily Counters</n-button>
      </div>
    </div>

    <!-- Aggregate Stats -->
    <div v-if="stats" class="stats-bar mb-24px">
      <div class="stat-item">
        <div class="stat-num">{{ stats.active_slots }}</div>
        <div class="stat-label">Active Slots</div>
      </div>
      <div class="stat-item">
        <div class="stat-num">{{ stats.total_monthly_capacity.toLocaleString() }}</div>
        <div class="stat-label">Monthly Capacity</div>
      </div>
      <div class="stat-item">
        <div class="stat-num">{{ stats.total_monthly_sent.toLocaleString() }}</div>
        <div class="stat-label">Sent This Month</div>
      </div>
      <div class="stat-item">
        <div class="stat-num">{{ (stats.total_monthly_capacity - stats.total_monthly_sent).toLocaleString() }}</div>
        <div class="stat-label">Remaining</div>
      </div>
      <div class="stat-item">
        <div class="stat-num">{{ stats.total_daily_sent }} / {{ stats.total_daily_capacity }}</div>
        <div class="stat-label">Today</div>
      </div>
      <div class="stat-item warn" v-if="stats.exhausted_slots > 0">
        <div class="stat-num">{{ stats.exhausted_slots }}</div>
        <div class="stat-label">Exhausted</div>
      </div>
      <div class="stat-item error" v-if="stats.error_slots > 0">
        <div class="stat-num">{{ stats.error_slots }}</div>
        <div class="stat-label">Errors</div>
      </div>
    </div>

    <!-- Monthly progress bar -->
    <div v-if="stats && stats.total_monthly_capacity > 0" class="mb-24px">
      <div class="flex justify-between text-12px text-desc mb-6px">
        <span>Monthly Usage</span>
        <span>{{ Math.round(stats.total_monthly_sent / stats.total_monthly_capacity * 100) }}% used</span>
      </div>
      <n-progress
        type="line"
        :percentage="Math.min(100, Math.round(stats.total_monthly_sent / stats.total_monthly_capacity * 100))"
        :status="stats.total_monthly_sent / stats.total_monthly_capacity > 0.9 ? 'error' : 'default'"
      />
    </div>

    <div v-if="loading" class="flex justify-center py-60px"><n-spin size="large" /></div>

    <!-- Provider Groups -->
    <div v-else>
      <div v-for="group in providerGroups" :key="group.type" class="provider-group mb-24px">
        <!-- Group Header -->
        <div class="group-header">
          <div class="group-left">
            <div class="group-name">{{ group.meta.label }}</div>
            <div class="group-limits">
              {{ group.meta.daily_free }}/day · {{ group.meta.monthly_free.toLocaleString() }}/month free ·
              <a :href="group.meta.signup_url" target="_blank" class="signup-link">Sign up →</a>
            </div>
          </div>
          <div class="group-combined">
            <span class="combined-label">Combined capacity:</span>
            <span class="combined-val">{{ (group.meta.monthly_free * 2).toLocaleString() }}/month</span>
          </div>
        </div>

        <!-- Two Slot Cards -->
        <div class="slot-grid">
          <div v-for="slot in group.slots" :key="slot.id" class="slot-card" :class="slotClass(slot)">
            <div class="slot-header">
              <div class="slot-title">Account {{ slot.slot }}</div>
              <n-tag :type="statusTagType(slot)" size="small">{{ slotStatusLabel(slot) }}</n-tag>
            </div>

            <!-- API Key Input -->
            <div class="slot-key-row">
              <n-input
                v-if="editingId === slot.id"
                v-model:value="editForm.apiKey"
                :placeholder="group.meta.api_key_label || 'API Key'"
                type="password"
                show-password-on="click"
                size="small"
                class="flex-1"
              />
              <div v-else class="key-display">
                <code>{{ slot.api_key || 'No API key — click Edit to add' }}</code>
              </div>
            </div>

            <!-- API Key 2 (if needed) -->
            <div v-if="group.meta.api_key_2_label" class="slot-key-row mt-8px">
              <n-input
                v-if="editingId === slot.id"
                v-model:value="editForm.apiKey2"
                :placeholder="group.meta.api_key_2_label"
                type="password"
                show-password-on="click"
                size="small"
                class="flex-1"
              />
              <div v-else class="key-display">
                <code class="text-desc">{{ slot.api_key_2 || group.meta.api_key_2_label + ': not set' }}</code>
              </div>
            </div>

            <!-- Usage Bars -->
            <div class="slot-usage" v-if="slot.is_active">
              <div class="usage-row">
                <span class="usage-label">Daily</span>
                <span class="usage-nums">{{ slot.daily_sent }} / {{ slot.daily_limit }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="Math.min(100, Math.round(slot.daily_sent / slot.daily_limit * 100))"
                :status="slot.daily_sent >= slot.daily_limit ? 'error' : 'default'"
                :height="6"
              />
              <div class="usage-row mt-6px">
                <span class="usage-label">Monthly</span>
                <span class="usage-nums">{{ slot.monthly_sent.toLocaleString() }} / {{ slot.monthly_limit.toLocaleString() }}</span>
              </div>
              <n-progress
                type="line"
                :percentage="Math.min(100, Math.round(slot.monthly_sent / slot.monthly_limit * 100))"
                :status="slot.monthly_sent >= slot.monthly_limit ? 'error' : 'default'"
                :height="6"
              />
            </div>

            <!-- Last error -->
            <div v-if="slot.last_error && slot.status === 'error'" class="slot-error">
              ⚠️ {{ slot.last_error.slice(0, 120) }}
            </div>

            <!-- Test result -->
            <div v-if="testResults[slot.id]" :class="['slot-test-result', testResults[slot.id].ok ? 'ok' : 'fail']">
              {{ testResults[slot.id].ok ? '✅' : '❌' }} {{ testResults[slot.id].message }}
            </div>

            <!-- Actions -->
            <div class="slot-actions">
              <template v-if="editingId === slot.id">
                <n-button size="tiny" type="primary" :loading="saving" @click="handleSave(slot)">Save</n-button>
                <n-button size="tiny" @click="editingId = null">Cancel</n-button>
              </template>
              <template v-else>
                <n-button size="tiny" @click="startEdit(slot, group.meta)">Edit Key</n-button>
                <n-button
                  size="tiny"
                  :loading="testingId === slot.id"
                  :disabled="!slot.api_key"
                  @click="handleTest(slot)"
                >Test</n-button>
                <n-switch
                  :value="slot.is_active"
                  size="small"
                  :disabled="!slot.api_key"
                  @update:value="(v: boolean) => handleToggle(slot, v)"
                />
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Pool Modal -->
    <n-modal v-model:show="showTestPool" title="Test Pool Auto-Route" preset="card" style="width:440px">
      <div class="text-desc text-13px mb-16px">
        Sends a test email via the pool. BillionMail picks the highest-priority available provider automatically.
      </div>
      <n-form label-placement="top">
        <n-form-item label="Send test to email">
          <n-input v-model:value="testToEmail" placeholder="your@email.com" />
        </n-form-item>
      </n-form>
      <div v-if="poolTestResult" :class="['slot-test-result mt-12px', poolTestResult.ok ? 'ok' : 'fail']">
        <div>{{ poolTestResult.ok ? '✅ Sent via' : '❌' }} {{ poolTestResult.provider_name }}</div>
        <div v-if="poolTestResult.latency_ms" class="text-12px text-desc mt-4px">{{ poolTestResult.latency_ms }}ms</div>
        <div v-if="!poolTestResult.ok" class="text-12px mt-4px">{{ poolTestResult.error }}</div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showTestPool = false">Close</n-button>
          <n-button type="primary" :loading="testingPool" @click="handleTestPool">Send Test</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import { listProviders, updateProvider, testProvider, testPool, resetCounters } from '@/api/modules/relay-pool'

const message = useMessage()

const loading = ref(false)
const saving = ref(false)
const testingId = ref<number | null>(null)
const testingPool = ref(false)
const showTestPool = ref(false)
const testToEmail = ref('')
const poolTestResult = ref<any>(null)
const editingId = ref<number | null>(null)
const editForm = ref({ apiKey: '', apiKey2: '' })
const testResults = ref<Record<number, any>>({})
const stats = ref<any>(null)

interface ProviderMeta {
  type: string
  label: string
  daily_free: number
  monthly_free: number
  signup_url: string
  api_key_label: string
  api_key_2_label?: string
  docs_url: string
}

interface ProviderSlot {
  id: number
  name: string
  provider_type: string
  slot: number
  api_key: string
  api_key_2: string
  daily_limit: number
  monthly_limit: number
  daily_sent: number
  monthly_sent: number
  status: string
  is_active: boolean
  last_error: string
}

interface ProviderGroup {
  type: string
  meta: ProviderMeta
  slots: ProviderSlot[]
}

const providerGroups = ref<ProviderGroup[]>([])

async function fetchProviders() {
  loading.value = true
  try {
    const res = await listProviders() as any
    const data = res?.data || {}
    stats.value = data.stats || null

    const metaMap: Record<string, ProviderMeta> = {}
    for (const m of (data.meta || [])) {
      metaMap[m.type] = m
    }

    const groupMap: Record<string, ProviderGroup> = {}
    for (const p of (data.providers || [])) {
      if (!groupMap[p.provider_type]) {
        groupMap[p.provider_type] = {
          type: p.provider_type,
          meta: metaMap[p.provider_type] || { type: p.provider_type, label: p.provider_type, daily_free: p.daily_limit, monthly_free: p.monthly_limit, signup_url: '#', api_key_label: 'API Key', docs_url: '#' },
          slots: [],
        }
      }
      groupMap[p.provider_type].slots.push(p)
    }

    // Sort groups by monthly capacity descending
    const groups = Object.values(groupMap)
    groups.sort((a, b) => (b.meta.monthly_free * 2) - (a.meta.monthly_free * 2))
    providerGroups.value = groups
  } catch (_e) { /* ignore */ } finally {
    loading.value = false
  }
}

function slotClass(slot: ProviderSlot) {
  if (slot.status === 'error') return 'slot-error-state'
  if (slot.status === 'exhausted') return 'slot-exhausted'
  if (slot.is_active) return 'slot-active'
  return 'slot-inactive'
}

function statusTagType(slot: ProviderSlot) {
  if (slot.status === 'active') return 'success'
  if (slot.status === 'error') return 'error'
  if (slot.status === 'exhausted') return 'warning'
  return 'default'
}

function slotStatusLabel(slot: ProviderSlot) {
  if (!slot.api_key) return 'No key'
  if (slot.status === 'exhausted') return 'Exhausted'
  if (slot.status === 'error') return 'Error'
  if (slot.is_active) return 'Active'
  return 'Disabled'
}

function startEdit(slot: ProviderSlot, _meta: ProviderMeta) {
  editingId.value = slot.id
  editForm.value = { apiKey: '', apiKey2: '' }
  delete testResults.value[slot.id]
}

async function handleSave(slot: ProviderSlot) {
  if (!editForm.value.apiKey && !editForm.value.apiKey2) {
    message.error('Enter an API key')
    return
  }
  saving.value = true
  try {
    const payload: any = { id: slot.id }
    if (editForm.value.apiKey) payload.api_key = editForm.value.apiKey
    if (editForm.value.apiKey2) payload.api_key_2 = editForm.value.apiKey2
    const active = true
    payload.is_active = active
    await updateProvider(payload)
    message.success('API key saved — slot activated')
    editingId.value = null
    fetchProviders()
  } catch (_e) {
    message.error('Failed to save')
  } finally {
    saving.value = false
  }
}

async function handleToggle(slot: ProviderSlot, value: boolean) {
  try {
    await updateProvider({ id: slot.id, is_active: value })
    message.success(value ? 'Slot enabled' : 'Slot disabled')
    fetchProviders()
  } catch (_e) { message.error('Failed') }
}

async function handleTest(slot: ProviderSlot) {
  const email = prompt('Send test to email address:')
  if (!email) return
  testingId.value = slot.id
  delete testResults.value[slot.id]
  try {
    const res = await testProvider({ id: slot.id, test_to_email: email }) as any
    const d = res?.data || {}
    testResults.value[slot.id] = {
      ok: d.ok,
      message: d.ok ? `Sent! ${d.latency_ms}ms${d.provider_message_id ? ' · ID: ' + d.provider_message_id : ''}` : (d.error || 'Failed'),
    }
  } catch (_e) {
    testResults.value[slot.id] = { ok: false, message: 'Request failed' }
  } finally {
    testingId.value = null
  }
}

async function handleTestPool() {
  if (!testToEmail.value) { message.error('Enter an email address'); return }
  testingPool.value = true
  poolTestResult.value = null
  try {
    const res = await testPool({ test_to_email: testToEmail.value }) as any
    const d = res?.data || {}
    poolTestResult.value = {
      ok: d.ok,
      provider_name: d.provider_name,
      latency_ms: d.latency_ms,
      error: d.error,
    }
  } catch (_e) {
    poolTestResult.value = { ok: false, error: 'Request failed' }
  } finally {
    testingPool.value = false
  }
}

async function handleResetAll() {
  try {
    await resetCounters({ reset_type: 'daily' })
    message.success('Daily counters reset')
    fetchProviders()
  } catch (_e) { message.error('Failed') }
}

onMounted(fetchProviders)
</script>

<style lang="scss" scoped>
.stats-bar {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;

  .stat-item {
    background: var(--color-fill-1);
    border: 1px solid var(--color-border-1);
    border-radius: 8px;
    padding: 12px 20px;
    text-align: center;
    min-width: 100px;

    &.warn { border-color: #faad14; background: #fffbe6; }
    &.error { border-color: #ff4d4f; background: #fff2f0; }

    .stat-num { font-size: 22px; font-weight: 700; }
    .stat-label { font-size: 11px; color: var(--color-text-3); margin-top: 2px; }
  }
}

.provider-group {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  overflow: hidden;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  background: var(--color-fill-1);
  border-bottom: 1px solid var(--color-border-1);

  .group-name { font-weight: 700; font-size: 15px; }
  .group-limits { font-size: 12px; color: var(--color-text-3); margin-top: 2px; }
  .signup-link { color: var(--color-primary); }
  .combined-label { font-size: 12px; color: var(--color-text-3); margin-right: 6px; }
  .combined-val { font-weight: 600; font-size: 14px; color: var(--color-primary); }
}

.slot-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;

  > :first-child { border-right: 1px solid var(--color-border-1); }
}

.slot-card {
  padding: 16px 20px;
  transition: background 0.2s;

  &.slot-active { background: var(--color-bg-2); }
  &.slot-inactive { background: var(--color-fill-1); opacity: 0.8; }
  &.slot-exhausted { background: #fffbe6; }
  &.slot-error-state { background: #fff2f0; }

  .slot-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }

  .slot-title { font-weight: 600; font-size: 13px; }

  .slot-key-row {
    display: flex;
    align-items: center;
    gap: 8px;

    .key-display {
      code {
        font-size: 12px;
        background: var(--color-fill-2);
        padding: 2px 8px;
        border-radius: 4px;
        word-break: break-all;
        color: var(--color-text-2);
      }
    }
  }

  .slot-usage {
    margin-top: 12px;
    .usage-row {
      display: flex;
      justify-content: space-between;
      font-size: 11px;
      margin-bottom: 3px;
    }
    .usage-label { color: var(--color-text-3); }
    .usage-nums { color: var(--color-text-2); font-weight: 500; }
  }

  .slot-error {
    margin-top: 8px;
    font-size: 11px;
    color: #a8071a;
    background: #fff2f0;
    padding: 4px 8px;
    border-radius: 4px;
    border: 1px solid #ffa39e;
  }

  .slot-test-result {
    margin-top: 8px;
    font-size: 12px;
    padding: 6px 10px;
    border-radius: 6px;
    &.ok { background: #f6ffed; color: #237804; border: 1px solid #b7eb8f; }
    &.fail { background: #fff2f0; color: #a8071a; border: 1px solid #ffa39e; }
  }

  .slot-actions {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 12px;
    flex-wrap: wrap;
  }
}

@media (max-width: 640px) {
  .slot-grid {
    grid-template-columns: 1fr;
    > :first-child { border-right: none; border-bottom: 1px solid var(--color-border-1); }
  }
}
</style>
