<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div>
        <div class="bt-title">IP Pools</div>
        <div class="text-desc text-13px mt-4px">Manage shared and dedicated IP pools for email delivery</div>
      </div>
      <n-button type="primary" @click="showCreate = true">
        <template #icon><n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></n-icon></template>
        New IP Pool
      </n-button>
    </div>

    <!-- Warming Guide -->
    <div class="warming-guide mb-24px">
      <div class="guide-title">📈 IP Warming Schedule</div>
      <div class="warming-days">
        <div v-for="w in warmingSchedule" :key="w.day" class="warm-step">
          <div class="warm-day">Day {{ w.day }}</div>
          <div class="warm-limit">{{ w.limit.toLocaleString() }}/day</div>
        </div>
      </div>
      <div class="guide-note">New IPs start at 100 emails/day and ramp up over 4 weeks. This signals to Gmail/Outlook that you're a legitimate sender, not a spam bot.</div>
    </div>

    <!-- Stats Row -->
    <div class="stats-row mb-24px">
      <div class="stat-card" v-for="stat in summaryStats" :key="stat.label">
        <div class="stat-value">{{ stat.value }}</div>
        <div class="stat-label">{{ stat.label }}</div>
      </div>
    </div>

    <!-- Pool List -->
    <div v-if="loading" class="flex justify-center py-60px"><n-spin size="large" /></div>
    <div v-else-if="pools.length === 0" class="empty-state">
      <div class="empty-icon">🔗</div>
      <div class="empty-title">No IP pools yet</div>
      <div class="empty-desc">Create a shared pool for free-plan users or dedicated pools for premium senders</div>
      <n-button type="primary" class="mt-16px" @click="showCreate = true">Create IP Pool</n-button>
    </div>
    <div v-else class="pool-grid">
      <div v-for="pool in pools" :key="pool.id" class="pool-card">
        <div class="pool-header">
          <div class="pool-name">{{ pool.name }}</div>
          <div class="pool-badges">
            <n-tag :type="pool.pool_type === 'dedicated' ? 'success' : 'info'" size="small" class="mr-6px">
              {{ pool.pool_type === 'dedicated' ? '⭐ Dedicated' : '🔄 Shared' }}
            </n-tag>
            <n-tag :type="statusColor(pool.status)" size="small">{{ pool.status }}</n-tag>
          </div>
        </div>

        <!-- Warming Progress -->
        <div v-if="pool.status === 'warming'" class="warming-progress mb-12px">
          <div class="flex justify-between text-12px mb-4px">
            <span class="text-desc">Warming Day {{ pool.warm_day }}/28</span>
            <span class="text-desc">{{ currentLimit(pool.warm_day) }}/day limit</span>
          </div>
          <n-progress type="line" :percentage="Math.round(pool.warm_day / 28 * 100)" status="warning" />
        </div>

        <!-- Daily Usage -->
        <div class="usage mb-12px">
          <div class="flex justify-between text-12px mb-4px">
            <span class="text-desc">Today's usage</span>
            <span class="text-desc">{{ pool.sent_today }} / {{ pool.daily_limit }}</span>
          </div>
          <n-progress
            type="line"
            :percentage="Math.min(100, Math.round(pool.sent_today / pool.daily_limit * 100))"
            :status="usageStatus(pool)"
          />
        </div>

        <!-- IPs -->
        <div class="pool-ips mb-12px">
          <div class="text-12px text-desc mb-6px">IP Addresses ({{ pool.ips.length }})</div>
          <div class="ip-tags">
            <n-tag v-for="ip in pool.ips" :key="ip" size="small" class="mr-4px mb-4px">{{ ip }}</n-tag>
            <span v-if="pool.ips.length === 0" class="text-desc text-12px">No IPs added yet</span>
          </div>
        </div>

        <div class="pool-actions">
          <n-button size="small" @click="handleAdvance(pool)" v-if="pool.status === 'warming'">
            Advance Warm Day
          </n-button>
          <n-button size="small" type="primary" v-if="pool.status === 'warming'" @click="handleActivate(pool)">
            Mark Active
          </n-button>
          <n-popconfirm @positive-click="handleDelete(pool.id)">
            <template #trigger>
              <n-button size="small" type="error">Delete</n-button>
            </template>
            Delete this IP pool?
          </n-popconfirm>
        </div>
      </div>
    </div>

    <!-- Create Pool Modal -->
    <n-modal v-model:show="showCreate" title="New IP Pool" preset="card" style="width:500px">
      <n-form :model="createForm" label-placement="top">
        <n-form-item label="Name">
          <n-input v-model:value="createForm.name" placeholder="e.g. Shared Pool - EU, Dedicated - Client A" />
        </n-form-item>
        <n-form-item label="Pool Type">
          <n-radio-group v-model:value="createForm.pool_type">
            <n-space>
              <n-radio value="shared">
                <div>
                  <div class="font-medium">Shared</div>
                  <div class="text-desc text-12px">Multiple senders share IPs (free plan)</div>
                </div>
              </n-radio>
              <n-radio value="dedicated">
                <div>
                  <div class="font-medium">Dedicated</div>
                  <div class="text-desc text-12px">Single sender owns IPs (paid plan)</div>
                </div>
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="IP Addresses (one per line)">
          <n-input
            v-model:value="ipInput"
            type="textarea"
            placeholder="192.168.1.1&#10;10.0.0.1"
            :rows="4"
          />
        </n-form-item>
        <n-form-item label="Daily Send Limit">
          <n-input-number v-model:value="createForm.daily_limit" :min="100" :max="1000000" />
          <template #feedback>Start at 100-500 for new IPs. Increase after warming.</template>
        </n-form-item>
        <n-form-item label="Description">
          <n-input v-model:value="createForm.description" placeholder="Optional" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showCreate = false">Cancel</n-button>
          <n-button type="primary" :loading="creating" @click="handleCreate">Create Pool</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import { listIpPools, createIpPool, deleteIpPool } from '@/api/modules/esp'

const message = useMessage()
const pools = ref<any[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const ipInput = ref('')
const createForm = ref({ name: '', pool_type: 'shared', daily_limit: 500, description: '', ips: [] as string[] })

const warmingSchedule = [
  { day: 1, limit: 100 },
  { day: 3, limit: 300 },
  { day: 7, limit: 1000 },
  { day: 14, limit: 5000 },
  { day: 21, limit: 20000 },
  { day: 28, limit: 100000 },
]

function currentLimit(warmDay: number): number {
  for (let i = warmingSchedule.length - 1; i >= 0; i--) {
    if (warmDay >= warmingSchedule[i].day) return warmingSchedule[i].limit
  }
  return 100
}

const summaryStats = computed(() => [
  { label: 'Total Pools', value: pools.value.length },
  { label: 'Dedicated', value: pools.value.filter(p => p.pool_type === 'dedicated').length },
  { label: 'Active', value: pools.value.filter(p => p.status === 'active').length },
  { label: 'Total IPs', value: pools.value.reduce((s, p) => s + (p.ips?.length || 0), 0) },
])

function statusColor(status: string) {
  if (status === 'active') return 'success'
  if (status === 'warming') return 'warning'
  return 'error'
}

function usageStatus(pool: any) {
  const pct = pool.sent_today / pool.daily_limit
  if (pct >= 0.9) return 'error'
  if (pct >= 0.7) return 'warning'
  return 'default'
}

async function fetchPools() {
  loading.value = true
  try {
    const res = await listIpPools() as any
    pools.value = res?.data || []
  } catch (_e) { /* ignore */ } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!createForm.value.name) { message.error('Name required'); return }
  creating.value = true
  const ips = ipInput.value.split('\n').map(s => s.trim()).filter(Boolean)
  try {
    await createIpPool({ ...createForm.value, ips })
    message.success('IP pool created')
    showCreate.value = false
    createForm.value = { name: '', pool_type: 'shared', daily_limit: 500, description: '', ips: [] }
    ipInput.value = ''
    fetchPools()
  } catch (_e) { message.error('Failed') } finally {
    creating.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteIpPool({ id })
    message.success('Deleted')
    fetchPools()
  } catch (_e) { message.error('Failed') }
}

function handleAdvance(pool: any) {
  message.info(`Warming advancement requires MTA server integration. Pool: ${pool.name}`)
}

function handleActivate(pool: any) {
  message.info(`Activation requires MTA server integration. Pool: ${pool.name}`)
}

onMounted(fetchPools)
</script>

<style lang="scss" scoped>
.warming-guide {
  background: linear-gradient(135deg, #f0f9ff 0%, #e8f4ff 100%);
  border: 1px solid #bfd7f5;
  border-radius: 8px;
  padding: 16px;

  .guide-title { font-weight: 600; margin-bottom: 12px; font-size: 14px; }

  .warming-days {
    display: flex;
    gap: 8px;
    margin-bottom: 10px;
    flex-wrap: wrap;
  }

  .warm-step {
    background: white;
    border-radius: 6px;
    padding: 8px 12px;
    text-align: center;
    min-width: 80px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.08);
    .warm-day { font-size: 11px; color: #666; }
    .warm-limit { font-size: 13px; font-weight: 600; color: #1a7fd4; }
  }

  .guide-note { font-size: 12px; color: #666; line-height: 1.5; }
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.stat-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  .stat-value { font-size: 28px; font-weight: 700; color: var(--color-primary); }
  .stat-label { font-size: 12px; color: var(--color-text-3); margin-top: 4px; }
}

.pool-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.pool-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;

  .pool-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 14px;
  }

  .pool-name { font-weight: 600; font-size: 15px; }
  .pool-badges { display: flex; align-items: center; }

  .pool-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }
}

.ip-tags { display: flex; flex-wrap: wrap; }

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  text-align: center;
  .empty-icon { font-size: 48px; margin-bottom: 16px; }
  .empty-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
  .empty-desc { font-size: 14px; color: var(--color-text-3); }
}
</style>
