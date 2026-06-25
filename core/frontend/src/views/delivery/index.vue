<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div>
        <div class="bt-title">Delivery Analytics</div>
        <div class="text-desc text-13px mt-4px">Real-time inbox placement, bounce tracking, and reputation monitoring</div>
      </div>
      <div class="flex gap-8px">
        <n-date-picker
          v-model:value="dateRange"
          type="daterange"
          clearable
          @update:value="fetchStats"
        />
        <n-button @click="fetchStats" :loading="loading">Refresh</n-button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="metric-grid mb-24px" v-if="summary">
      <div class="metric-card">
        <div class="metric-icon send-icon">📤</div>
        <div class="metric-val">{{ summary.sent.toLocaleString() }}</div>
        <div class="metric-lbl">Total Sent</div>
      </div>
      <div class="metric-card success">
        <div class="metric-icon">✅</div>
        <div class="metric-val">{{ summary.delivery_rate.toFixed(1) }}%</div>
        <div class="metric-lbl">Delivery Rate</div>
        <div class="metric-sub">{{ summary.delivered.toLocaleString() }} delivered</div>
      </div>
      <div class="metric-card info">
        <div class="metric-icon">👁️</div>
        <div class="metric-val">{{ summary.open_rate.toFixed(1) }}%</div>
        <div class="metric-lbl">Open Rate</div>
        <div class="metric-sub">{{ summary.opened.toLocaleString() }} opens</div>
      </div>
      <div class="metric-card" :class="summary.bounce_rate > 5 ? 'danger' : 'warning'">
        <div class="metric-icon">↩️</div>
        <div class="metric-val">{{ summary.bounce_rate.toFixed(2) }}%</div>
        <div class="metric-lbl">Bounce Rate</div>
        <div class="metric-sub">
          {{ summary.bounced_hard }} hard · {{ summary.bounced_soft }} soft
          <n-tag v-if="summary.bounce_rate > 5" type="error" size="tiny" class="ml-4px">High!</n-tag>
          <n-tag v-else-if="summary.bounce_rate > 2" type="warning" size="tiny" class="ml-4px">Watch</n-tag>
          <n-tag v-else type="success" size="tiny" class="ml-4px">Good</n-tag>
        </div>
      </div>
      <div class="metric-card" :class="summary.complained > 10 ? 'danger' : ''">
        <div class="metric-icon">🚨</div>
        <div class="metric-val">{{ summary.complained }}</div>
        <div class="metric-lbl">Spam Complaints</div>
        <div class="metric-sub">
          <span v-if="summary.sent > 0">{{ (summary.complained / summary.sent * 100).toFixed(3) }}% complaint rate</span>
          <n-tag v-if="summary.complained > 0.1 / 100 * summary.sent" type="error" size="tiny" class="ml-4px">Above Threshold!</n-tag>
        </div>
      </div>
      <div class="metric-card">
        <div class="metric-icon">🖱️</div>
        <div class="metric-val">{{ summary.sent > 0 ? (summary.clicked / summary.sent * 100).toFixed(1) : 0 }}%</div>
        <div class="metric-lbl">Click Rate</div>
        <div class="metric-sub">{{ summary.clicked.toLocaleString() }} clicks</div>
      </div>
    </div>

    <!-- Reputation Health -->
    <div class="reputation-card mb-24px">
      <div class="rep-title">🛡️ Reputation Health</div>
      <div class="rep-checks">
        <div class="rep-check" v-for="check in reputationChecks" :key="check.label">
          <div :class="['rep-dot', check.ok ? 'rep-ok' : 'rep-bad']"></div>
          <div class="rep-label">{{ check.label }}</div>
          <div :class="['rep-status', check.ok ? 'text-success' : 'text-error']">{{ check.status }}</div>
        </div>
      </div>
    </div>

    <!-- Chart placeholder + By Domain Table -->
    <div class="charts-row mb-24px">
      <!-- Daily Trend -->
      <div class="chart-card">
        <div class="chart-title">Sending Volume — Last 30 Days</div>
        <div v-if="byDay.length === 0" class="chart-empty">No data yet. Start sending to see trends.</div>
        <div v-else class="mini-chart">
          <div
            v-for="d in byDay.slice(-14)"
            :key="d.date_str"
            class="bar-col"
            :title="`${d.date_str}: ${d.sent} sent, ${d.delivery_rate.toFixed(1)}% delivered`"
          >
            <div class="bar-stack">
              <div class="bar-delivered" :style="{ height: barH(d.delivered, maxSent) + 'px' }"></div>
              <div class="bar-bounced" :style="{ height: barH(d.bounced_hard + d.bounced_soft, maxSent) + 'px' }"></div>
            </div>
            <div class="bar-label">{{ d.date_str.slice(5) }}</div>
          </div>
        </div>
        <div class="chart-legend">
          <span class="legend-dot delivered"></span>Delivered
          <span class="legend-dot bounced ml-12px"></span>Bounced
        </div>
      </div>
    </div>

    <!-- By Domain Table -->
    <div class="domain-table-card mb-24px">
      <div class="table-title mb-16px">Performance by Domain</div>
      <n-data-table
        :columns="domainColumns"
        :data="byDomain"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        size="small"
      />
    </div>

    <!-- Suppression List Quick View -->
    <div class="suppression-card">
      <div class="flex items-center justify-between mb-16px">
        <div class="table-title">Suppression List</div>
        <div class="flex gap-8px">
          <n-input v-model:value="supKeyword" placeholder="Search email..." size="small" clearable @update:value="fetchSuppression" />
          <n-select v-model:value="supReason" :options="reasonOptions" size="small" clearable placeholder="All reasons" style="width:160px" @update:value="fetchSuppression" />
          <n-button size="small" @click="showAddSup = true">Add</n-button>
        </div>
      </div>
      <n-data-table
        :columns="suppressionColumns"
        :data="suppressionList"
        :loading="supLoading"
        :pagination="{ pageSize: 10 }"
        size="small"
      />
    </div>

    <!-- Add Suppression Modal -->
    <n-modal v-model:show="showAddSup" title="Add to Suppression List" preset="card" style="width:400px">
      <n-form label-placement="top">
        <n-form-item label="Email">
          <n-input v-model:value="supEmail" placeholder="bad@example.com" />
        </n-form-item>
        <n-form-item label="Reason">
          <n-select v-model:value="supNewReason" :options="reasonOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showAddSup = false">Cancel</n-button>
          <n-button type="primary" :loading="addingSupp" @click="handleAddSuppression">Add</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { h } from 'vue'
import { NTag, NButton, useMessage } from 'naive-ui'
import { getDeliveryStats, listSuppression, addSuppression, deleteSuppression } from '@/api/modules/esp'

const message = useMessage()
const loading = ref(false)
const dateRange = ref<[number, number] | null>(null)

const summary = ref<any>(null)
const byDay = ref<any[]>([])
const byDomain = ref<any[]>([])

const suppressionList = ref<any[]>([])
const supLoading = ref(false)
const supKeyword = ref('')
const supReason = ref('')
const showAddSup = ref(false)
const supEmail = ref('')
const supNewReason = ref('manual')
const addingSupp = ref(false)

const reasonOptions = [
  { label: 'Hard Bounce', value: 'hard_bounce' },
  { label: 'Soft Bounce', value: 'soft_bounce' },
  { label: 'Spam Complaint', value: 'complaint' },
  { label: 'Unsubscribe', value: 'unsubscribe' },
  { label: 'Manual', value: 'manual' },
]

const reputationChecks = computed(() => {
  const s = summary.value
  return [
    { label: 'Bounce Rate < 5%', ok: !s || s.bounce_rate < 5, status: s ? `${s.bounce_rate.toFixed(2)}%` : 'No data' },
    { label: 'Complaint Rate < 0.1%', ok: !s || (s.sent === 0 || s.complained / s.sent < 0.001), status: s && s.sent > 0 ? `${(s.complained / s.sent * 100).toFixed(3)}%` : 'No data' },
    { label: 'Delivery Rate > 95%', ok: !s || s.delivery_rate > 95, status: s ? `${s.delivery_rate.toFixed(1)}%` : 'No data' },
    { label: 'Open Rate > 15%', ok: !s || s.open_rate > 15, status: s ? `${s.open_rate.toFixed(1)}%` : 'No data' },
  ]
})

const maxSent = computed(() => {
  return Math.max(...byDay.value.map(d => d.sent), 1)
})

function barH(val: number, max: number): number {
  return Math.max(2, Math.round((val / max) * 80))
}

const domainColumns = [
  { title: 'Domain', key: 'domain' },
  { title: 'Sent', key: 'sent', sorter: (a: any, b: any) => a.sent - b.sent },
  { title: 'Delivery Rate', key: 'delivery_rate', render: (row: any) => h(NTag, { type: row.delivery_rate > 95 ? 'success' : row.delivery_rate > 85 ? 'warning' : 'error', size: 'small' }, { default: () => `${row.delivery_rate.toFixed(1)}%` }) },
  { title: 'Open Rate', key: 'open_rate', render: (row: any) => `${row.open_rate.toFixed(1)}%` },
  { title: 'Bounce Rate', key: 'bounce_rate', render: (row: any) => h(NTag, { type: row.bounce_rate < 2 ? 'success' : row.bounce_rate < 5 ? 'warning' : 'error', size: 'small' }, { default: () => `${row.bounce_rate.toFixed(2)}%` }) },
  { title: 'Complaints', key: 'complained' },
]

const suppressionColumns = [
  { title: 'Email', key: 'email' },
  { title: 'Reason', key: 'reason', render: (row: any) => h(NTag, { type: reasonTagType(row.reason), size: 'small' }, { default: () => row.reason }) },
  { title: 'Domain', key: 'domain' },
  { title: 'Added', key: 'created_at', render: (row: any) => new Date(row.created_at * 1000).toLocaleDateString() },
  { title: 'Action', key: 'actions', render: (row: any) => h(NButton, { size: 'tiny', type: 'error', onClick: () => handleDeleteSuppression(row.id) }, { default: () => 'Remove' }) },
]

function reasonTagType(reason: string) {
  if (reason === 'hard_bounce') return 'error'
  if (reason === 'complaint') return 'error'
  if (reason === 'soft_bounce') return 'warning'
  return 'default'
}

async function fetchStats() {
  loading.value = true
  try {
    const params: any = {}
    if (dateRange.value) {
      params.date_from = new Date(dateRange.value[0]).toISOString().slice(0, 10)
      params.date_to = new Date(dateRange.value[1]).toISOString().slice(0, 10)
    }
    const res = await getDeliveryStats(params) as any
    summary.value = res?.data?.summary || null
    byDay.value = res?.data?.by_day || []
    byDomain.value = res?.data?.by_domain || []
  } catch (_e) { /* ignore */ } finally {
    loading.value = false
  }
}

async function fetchSuppression() {
  supLoading.value = true
  try {
    const res = await listSuppression({ page: 1, page_size: 50, keyword: supKeyword.value, reason: supReason.value }) as any
    suppressionList.value = res?.data?.list || []
  } catch (_e) { /* ignore */ } finally {
    supLoading.value = false
  }
}

async function handleAddSuppression() {
  if (!supEmail.value) { message.error('Email required'); return }
  addingSupp.value = true
  try {
    await addSuppression({ email: supEmail.value, reason: supNewReason.value })
    message.success('Added to suppression list')
    showAddSup.value = false
    supEmail.value = ''
    fetchSuppression()
  } catch (_e) { message.error('Failed') } finally {
    addingSupp.value = false
  }
}

async function handleDeleteSuppression(id: number) {
  try {
    await deleteSuppression({ id })
    message.success('Removed')
    fetchSuppression()
  } catch (_e) { message.error('Failed') }
}

onMounted(() => {
  fetchStats()
  fetchSuppression()
})
</script>

<style lang="scss" scoped>
.metric-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;

  @media (max-width: 900px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.metric-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  padding: 20px;
  position: relative;

  &.success { border-color: #52c41a33; background: #f6ffed; }
  &.info { border-color: #1890ff33; background: #e6f7ff; }
  &.warning { border-color: #faad1433; background: #fffbe6; }
  &.danger { border-color: #ff4d4f44; background: #fff2f0; }

  .metric-icon { font-size: 24px; margin-bottom: 8px; }
  .metric-val { font-size: 32px; font-weight: 700; line-height: 1; }
  .metric-lbl { font-size: 13px; color: var(--color-text-3); margin-top: 4px; font-weight: 600; }
  .metric-sub { font-size: 12px; color: var(--color-text-3); margin-top: 4px; }
}

.reputation-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  padding: 20px;

  .rep-title { font-weight: 700; font-size: 15px; margin-bottom: 16px; }

  .rep-checks {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
  }

  .rep-check {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 12px;
    background: var(--color-fill-1);
    border-radius: 8px;
    text-align: center;

    .rep-dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
      margin-bottom: 6px;
      &.rep-ok { background: #52c41a; }
      &.rep-bad { background: #ff4d4f; }
    }

    .rep-label { font-size: 12px; color: var(--color-text-3); margin-bottom: 4px; }
    .rep-status { font-size: 13px; font-weight: 600; }
    .text-success { color: #52c41a; }
    .text-error { color: #ff4d4f; }
  }
}

.charts-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.chart-card, .domain-table-card, .suppression-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  padding: 20px;

  .chart-title, .table-title { font-weight: 700; font-size: 15px; margin-bottom: 16px; }
  .chart-empty { text-align: center; padding: 40px; color: var(--color-text-3); font-size: 14px; }
}

.mini-chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 100px;
  padding: 0 4px;

  .bar-col {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    cursor: pointer;

    .bar-stack {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: flex-end;
      width: 100%;
      height: 84px;
      gap: 1px;
    }

    .bar-delivered { background: #52c41a; width: 80%; border-radius: 2px 2px 0 0; min-height: 2px; }
    .bar-bounced { background: #ff4d4f; width: 80%; border-radius: 2px 2px 0 0; min-height: 2px; }
    .bar-label { font-size: 9px; color: var(--color-text-3); white-space: nowrap; }
  }
}

.chart-legend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 8px;

  .legend-dot {
    width: 10px;
    height: 10px;
    border-radius: 2px;
    display: inline-block;
    &.delivered { background: #52c41a; }
    &.bounced { background: #ff4d4f; }
  }
}
</style>
