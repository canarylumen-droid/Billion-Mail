<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div>
        <div class="bt-title">SMTP Relay</div>
        <div class="text-desc text-13px mt-4px">Your own SMTP provider — give users host:587 credentials to send through you</div>
      </div>
      <n-button type="primary" @click="showCreate = true">
        <template #icon><n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></n-icon></template>
        New Credential
      </n-button>
    </div>

    <!-- Info Banner -->
    <div class="info-banner mb-24px">
      <div class="info-icon">📡</div>
      <div class="info-body">
        <div class="info-title">Your SMTP Relay Server</div>
        <div class="info-desc">
          Users connect to <code>smtp.yourdomain.com:587</code> with their credentials. Your MTA server handles delivery to Gmail, Outlook, Yahoo via port 25.
          <a href="https://www.hetzner.com/cloud" target="_blank" class="text-primary ml-4px">Deploy MTA on Hetzner →</a>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <n-tabs v-model:value="activeTab" type="line" class="mb-20px">
      <n-tab name="credentials">SMTP Credentials</n-tab>
      <n-tab name="domains">Sending Domains</n-tab>
      <n-tab name="mta">MTA Servers</n-tab>
    </n-tabs>

    <!-- SMTP Credentials Tab -->
    <div v-if="activeTab === 'credentials'">
      <div v-if="loading" class="flex justify-center py-60px"><n-spin size="large" /></div>
      <div v-else-if="credentials.length === 0" class="empty-state">
        <div class="empty-icon">🔑</div>
        <div class="empty-title">No credentials yet</div>
        <div class="empty-desc">Create SMTP credentials to give users relay access</div>
        <n-button type="primary" class="mt-16px" @click="showCreate = true">Create First Credential</n-button>
      </div>
      <div v-else class="cred-grid">
        <div v-for="cred in credentials" :key="cred.id" class="cred-card">
          <div class="cred-header">
            <div class="cred-name">{{ cred.name }}</div>
            <n-tag :type="cred.status === 1 ? 'success' : 'error'" size="small">
              {{ cred.status === 1 ? 'Active' : 'Disabled' }}
            </n-tag>
          </div>
          <div class="cred-detail">
            <div class="detail-row">
              <span class="detail-label">Host</span>
              <code class="detail-val">smtp.yourdomain.com:587</code>
            </div>
            <div class="detail-row">
              <span class="detail-label">Username</span>
              <code class="detail-val">{{ cred.username }}</code>
              <n-button size="tiny" text @click="copy(cred.username)">Copy</n-button>
            </div>
            <div class="detail-row">
              <span class="detail-label">Plan</span>
              <n-tag size="small" :type="planType(cred.plan)">{{ cred.plan }}</n-tag>
            </div>
            <div class="detail-row">
              <span class="detail-label">Daily Limit</span>
              <span>{{ cred.sent_today }} / {{ cred.daily_limit }} sent today</span>
            </div>
          </div>
          <n-progress
            type="line"
            :percentage="Math.min(100, Math.round(cred.sent_today / cred.daily_limit * 100))"
            :status="cred.sent_today >= cred.daily_limit ? 'error' : 'default'"
            class="mb-12px"
          />
          <div class="cred-actions">
            <n-button size="small" @click="handleRegenerate(cred)">Regenerate Password</n-button>
            <n-popconfirm @positive-click="handleDelete(cred.id)">
              <template #trigger>
                <n-button size="small" type="error">Delete</n-button>
              </template>
              Delete this credential?
            </n-popconfirm>
          </div>
        </div>
      </div>
    </div>

    <!-- Sending Domains Tab -->
    <div v-if="activeTab === 'domains'">
      <div class="flex justify-end mb-16px">
        <n-button type="primary" size="small" @click="showAddDomain = true">Add Domain</n-button>
      </div>
      <div v-if="domainsLoading" class="flex justify-center py-40px"><n-spin /></div>
      <div v-else-if="domains.length === 0" class="empty-state">
        <div class="empty-icon">🌐</div>
        <div class="empty-title">No sending domains</div>
        <div class="empty-desc">Add a domain to get SPF, DKIM, and DMARC DNS records</div>
        <n-button type="primary" class="mt-16px" @click="showAddDomain = true">Add Domain</n-button>
      </div>
      <div v-else>
        <div v-for="d in domains" :key="d.id" class="domain-card mb-16px">
          <div class="domain-header">
            <div class="domain-name">{{ d.domain }}</div>
            <div class="domain-badges">
              <n-tag :type="verifyBadge(d.spf_verified)" size="small" class="mr-6px">SPF</n-tag>
              <n-tag :type="verifyBadge(d.dkim_verified)" size="small" class="mr-6px">DKIM</n-tag>
              <n-tag :type="verifyBadge(d.dmarc_verified)" size="small" class="mr-6px">DMARC</n-tag>
              <n-tag :type="statusBadge(d.status)" size="small">{{ d.status }}</n-tag>
            </div>
          </div>
          <div class="domain-actions">
            <n-button size="small" @click="handleShowDns(d)">View DNS Records</n-button>
            <n-button size="small" type="primary" :loading="verifyingId === d.id" @click="handleVerify(d)">Verify DNS</n-button>
            <n-popconfirm @positive-click="handleDeleteDomain(d.id)">
              <template #trigger>
                <n-button size="small" type="error">Remove</n-button>
              </template>
              Remove this domain?
            </n-popconfirm>
          </div>
        </div>
      </div>
    </div>

    <!-- MTA Servers Tab -->
    <div v-if="activeTab === 'mta'">
      <div class="flex justify-end mb-16px">
        <n-button type="primary" size="small" @click="showAddMta = true">Register MTA Server</n-button>
      </div>
      <div class="mta-info-card mb-20px">
        <div class="mta-info-title">🖥️ What is an MTA Server?</div>
        <div class="mta-info-body">
          The MTA (Mail Transfer Agent) is the server that actually delivers emails to Gmail, Outlook, Yahoo via port 25.
          Deploy on <strong>Hetzner, OVH, or Contabo</strong> (€4-20/mo) — port 25 is open by default.
          Install <a href="https://haraka.github.io" target="_blank" class="text-primary">Haraka</a> or
          <a href="https://postal.atech.media" target="_blank" class="text-primary">Postal</a> on it.
        </div>
      </div>
      <div v-if="mtaLoading" class="flex justify-center py-40px"><n-spin /></div>
      <div v-else-if="mtaServers.length === 0" class="empty-state">
        <div class="empty-icon">🖥️</div>
        <div class="empty-title">No MTA servers registered</div>
        <div class="empty-desc">Register your Hetzner/OVH server with Haraka installed</div>
        <n-button type="primary" class="mt-16px" @click="showAddMta = true">Register Server</n-button>
      </div>
      <div v-else class="mta-grid">
        <div v-for="srv in mtaServers" :key="srv.id" class="mta-card">
          <div class="mta-header">
            <div class="mta-name">{{ srv.name }}</div>
            <n-tag :type="srv.status === 'active' ? 'success' : 'warning'" size="small">{{ srv.status }}</n-tag>
          </div>
          <div class="mta-details">
            <div class="detail-row"><span class="detail-label">Host</span><code>{{ srv.host }}</code></div>
            <div class="detail-row"><span class="detail-label">Region</span><span>{{ srv.region || '—' }}</span></div>
            <div class="detail-row"><span class="detail-label">IPs</span><span>{{ srv.ip_count }}</span></div>
          </div>
          <n-popconfirm @positive-click="handleDeleteMta(srv.id)">
            <template #trigger>
              <n-button size="small" type="error" class="mt-12px">Remove</n-button>
            </template>
            Remove this MTA server?
          </n-popconfirm>
        </div>
      </div>
    </div>

    <!-- Create Credential Modal -->
    <n-modal v-model:show="showCreate" title="New SMTP Credential" preset="card" style="width:480px">
      <n-form :model="createForm" label-placement="top">
        <n-form-item label="Name">
          <n-input v-model:value="createForm.name" placeholder="e.g. Audnix AI Outreach" />
        </n-form-item>
        <n-form-item label="Description">
          <n-input v-model:value="createForm.description" placeholder="Optional description" />
        </n-form-item>
        <n-form-item label="Plan">
          <n-select v-model:value="createForm.plan" :options="planOptions" />
        </n-form-item>
        <n-form-item label="Daily Send Limit">
          <n-input-number v-model:value="createForm.daily_limit" :min="1" :max="100000" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showCreate = false">Cancel</n-button>
          <n-button type="primary" :loading="creating" @click="handleCreate">Create</n-button>
        </div>
      </template>
    </n-modal>

    <!-- Add Domain Modal -->
    <n-modal v-model:show="showAddDomain" title="Add Sending Domain" preset="card" style="width:480px">
      <n-form label-placement="top">
        <n-form-item label="Domain">
          <n-input v-model:value="newDomain" placeholder="e.g. yourdomain.com" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showAddDomain = false">Cancel</n-button>
          <n-button type="primary" :loading="addingDomain" @click="handleAddDomain">Add & Get DNS Records</n-button>
        </div>
      </template>
    </n-modal>

    <!-- DNS Records Modal -->
    <n-modal v-model:show="showDns" title="DNS Records" preset="card" style="width:680px">
      <div class="mb-12px text-desc text-13px">Add these DNS records to your domain to enable sending. After adding, click "Verify DNS" on the domain.</div>
      <div v-for="rec in dnsRecords" :key="rec.host" class="dns-record">
        <div class="dns-row">
          <span class="dns-label">Type</span><code class="dns-val">{{ rec.type }}</code>
        </div>
        <div class="dns-row">
          <span class="dns-label">Host</span>
          <code class="dns-val">{{ rec.host }}</code>
          <n-button size="tiny" text @click="copy(rec.host)">Copy</n-button>
        </div>
        <div class="dns-row">
          <span class="dns-label">Value</span>
          <code class="dns-val dns-val-long">{{ rec.value }}</code>
          <n-button size="tiny" text @click="copy(rec.value)">Copy</n-button>
        </div>
        <div class="dns-row"><span class="dns-label">TTL</span><span>{{ rec.ttl }}</span></div>
        <div class="dns-note">{{ rec.note }}</div>
      </div>
      <template #footer>
        <n-button @click="showDns = false">Close</n-button>
      </template>
    </n-modal>

    <!-- Add MTA Modal -->
    <n-modal v-model:show="showAddMta" title="Register MTA Server" preset="card" style="width:480px">
      <n-form :model="mtaForm" label-placement="top">
        <n-form-item label="Name">
          <n-input v-model:value="mtaForm.name" placeholder="e.g. Hetzner EU MTA 1" />
        </n-form-item>
        <n-form-item label="Host / IP">
          <n-input v-model:value="mtaForm.host" placeholder="e.g. 94.130.1.100 or mta.yourdomain.com" />
        </n-form-item>
        <n-form-item label="Region">
          <n-select v-model:value="mtaForm.region" :options="regionOptions" />
        </n-form-item>
        <n-form-item label="SSH Port">
          <n-input-number v-model:value="mtaForm.ssh_port" :min="1" :max="65535" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-8px">
          <n-button @click="showAddMta = false">Cancel</n-button>
          <n-button type="primary" :loading="addingMta" @click="handleAddMta">Register</n-button>
        </div>
      </template>
    </n-modal>

    <!-- Password Reveal Modal -->
    <n-modal v-model:show="showPassword" title="New Password Generated" preset="card" style="width:400px">
      <div class="text-desc mb-8px text-13px">Save this password — it won't be shown again.</div>
      <div class="password-box">
        <code>{{ newPassword }}</code>
        <n-button size="small" @click="copy(newPassword)" class="ml-8px">Copy</n-button>
      </div>
      <template #footer>
        <n-button type="primary" @click="showPassword = false">Done</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import {
  listCredentials, createCredential, deleteCredential, regeneratePassword,
  listSendingDomains, addSendingDomain, verifyDomain, deleteSendingDomain,
  listMtaServers, addMtaServer, deleteMtaServer,
} from '@/api/modules/esp'

const message = useMessage()
const activeTab = ref('credentials')

// ── Credentials ──────────────────────────────────────────────────────────────
const credentials = ref<any[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const showPassword = ref(false)
const newPassword = ref('')

const createForm = ref({ name: '', description: '', plan: 'free', daily_limit: 1000, ip_pool_id: 0 })

const planOptions = [
  { label: 'Free (Shared IP)', value: 'free' },
  { label: 'Starter (Shared IP, higher limits)', value: 'starter' },
  { label: 'Pro (Dedicated IP)', value: 'pro' },
]

async function fetchCredentials() {
  loading.value = true
  try {
    const res = await listCredentials() as any
    credentials.value = res?.data || []
  } catch (_e) { /* ignore */ } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!createForm.value.name) { message.error('Name required'); return }
  creating.value = true
  try {
    const res = await createCredential(createForm.value) as any
    message.success(res?.message || 'Created')
    if (res?.message?.includes('password:')) {
      newPassword.value = res.message.split('password: ')[1] || ''
      showPassword.value = true
    }
    showCreate.value = false
    createForm.value = { name: '', description: '', plan: 'free', daily_limit: 1000, ip_pool_id: 0 }
    fetchCredentials()
  } catch (_e) { message.error('Failed to create') } finally {
    creating.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteCredential({ id })
    message.success('Deleted')
    fetchCredentials()
  } catch (_e) { message.error('Delete failed') }
}

async function handleRegenerate(cred: any) {
  try {
    const res = await regeneratePassword({ id: cred.id }) as any
    newPassword.value = res?.data?.password || ''
    showPassword.value = true
  } catch (_e) { message.error('Failed') }
}

// ── Domains ───────────────────────────────────────────────────────────────────
const domains = ref<any[]>([])
const domainsLoading = ref(false)
const showAddDomain = ref(false)
const addingDomain = ref(false)
const newDomain = ref('')
const showDns = ref(false)
const dnsRecords = ref<any[]>([])
const verifyingId = ref<number | null>(null)

async function fetchDomains() {
  domainsLoading.value = true
  try {
    const res = await listSendingDomains() as any
    domains.value = res?.data || []
  } catch (_e) { /* ignore */ } finally {
    domainsLoading.value = false
  }
}

async function handleAddDomain() {
  if (!newDomain.value) { message.error('Domain required'); return }
  addingDomain.value = true
  try {
    const res = await addSendingDomain({ domain: newDomain.value }) as any
    dnsRecords.value = res?.data?.dns_records || []
    showAddDomain.value = false
    showDns.value = true
    newDomain.value = ''
    fetchDomains()
  } catch (_e) { message.error('Failed to add domain') } finally {
    addingDomain.value = false
  }
}

async function handleVerify(d: any) {
  verifyingId.value = d.id
  try {
    const res = await verifyDomain({ id: d.id }) as any
    const dom = res?.data?.domain
    if (dom?.status === 'verified') {
      message.success('All DNS records verified! ✅')
    } else {
      message.warning('Some records not yet verified. Check DNS propagation (can take up to 48h).')
    }
    dnsRecords.value = res?.data?.dns_records || []
    fetchDomains()
  } catch (_e) { message.error('Verification failed') } finally {
    verifyingId.value = null
  }
}

function handleShowDns(d: any) {
  const sel = d.dkim_selector || 'bm1'
  const pub = d.dkim_public_key || ''
  const dom = d.domain
  dnsRecords.value = [
    { type: 'TXT', host: dom, value: d.spf_record || `v=spf1 include:relay.yourdomain.com ~all`, ttl: '3600', note: 'SPF record' },
    { type: 'TXT', host: `${sel}._domainkey.${dom}`, value: `v=DKIM1; k=rsa; p=${pub}`, ttl: '3600', note: 'DKIM public key' },
    { type: 'TXT', host: `_dmarc.${dom}`, value: d.dmarc_record || `v=DMARC1; p=none; rua=mailto:dmarc@${dom}`, ttl: '3600', note: 'DMARC policy' },
  ]
  showDns.value = true
}

async function handleDeleteDomain(id: number) {
  try {
    await deleteSendingDomain({ id })
    message.success('Domain removed')
    fetchDomains()
  } catch (_e) { message.error('Failed') }
}

// ── MTA Servers ───────────────────────────────────────────────────────────────
const mtaServers = ref<any[]>([])
const mtaLoading = ref(false)
const showAddMta = ref(false)
const addingMta = ref(false)
const mtaForm = ref({ name: '', host: '', region: 'eu-central', ssh_port: 22 })
const regionOptions = [
  { label: 'EU Central (Germany/France)', value: 'eu-central' },
  { label: 'EU West (UK/Netherlands)', value: 'eu-west' },
  { label: 'US East', value: 'us-east' },
  { label: 'US West', value: 'us-west' },
  { label: 'Asia Pacific', value: 'ap' },
]

async function fetchMta() {
  mtaLoading.value = true
  try {
    const res = await listMtaServers() as any
    mtaServers.value = res?.data || []
  } catch (_e) { /* ignore */ } finally {
    mtaLoading.value = false
  }
}

async function handleAddMta() {
  if (!mtaForm.value.name || !mtaForm.value.host) { message.error('Name and host required'); return }
  addingMta.value = true
  try {
    await addMtaServer(mtaForm.value)
    message.success('MTA server registered')
    showAddMta.value = false
    mtaForm.value = { name: '', host: '', region: 'eu-central', ssh_port: 22 }
    fetchMta()
  } catch (_e) { message.error('Failed') } finally {
    addingMta.value = false
  }
}

async function handleDeleteMta(id: number) {
  try {
    await deleteMtaServer({ id })
    message.success('Removed')
    fetchMta()
  } catch (_e) { message.error('Failed') }
}

// ── Utils ─────────────────────────────────────────────────────────────────────
function copy(text: string) {
  navigator.clipboard.writeText(text)
  message.success('Copied!')
}

function planType(plan: string) {
  return plan === 'pro' ? 'success' : plan === 'starter' ? 'warning' : 'default'
}

function verifyBadge(ok: boolean) {
  return ok ? 'success' : 'error'
}

function statusBadge(status: string) {
  if (status === 'verified') return 'success'
  if (status === 'partial') return 'warning'
  return 'default'
}

watch(activeTab, (tab) => {
  if (tab === 'domains') fetchDomains()
  if (tab === 'mta') fetchMta()
})

onMounted(() => { fetchCredentials() })
</script>

<style lang="scss" scoped>
.info-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  background: var(--color-fill-1);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;

  .info-icon { font-size: 22px; }
  .info-title { font-weight: 600; font-size: 14px; margin-bottom: 4px; }
  .info-desc { font-size: 13px; color: var(--color-text-3); }
  code { background: var(--color-fill-2); padding: 1px 6px; border-radius: 4px; font-size: 12px; }
}

.cred-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.cred-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;

  .cred-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .cred-name { font-weight: 600; font-size: 15px; }

  .cred-detail { margin-bottom: 12px; }

  .cred-actions {
    display: flex;
    gap: 8px;
  }
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-size: 13px;

  .detail-label { color: var(--color-text-3); min-width: 90px; }
  code { background: var(--color-fill-2); padding: 1px 6px; border-radius: 4px; font-size: 12px; }
}

.domain-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;

  .domain-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .domain-name { font-weight: 600; font-size: 15px; }
  .domain-badges { display: flex; align-items: center; }
  .domain-actions { display: flex; gap: 8px; }
}

.mta-info-card {
  background: #f0f7ff;
  border: 1px solid #bfd7f5;
  border-radius: 8px;
  padding: 16px;

  .mta-info-title { font-weight: 600; margin-bottom: 8px; }
  .mta-info-body { font-size: 13px; line-height: 1.6; color: #444; }
}

.mta-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.mta-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 8px;
  padding: 16px;

  .mta-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .mta-name { font-weight: 600; }
  .mta-details { font-size: 13px; }
}

.dns-record {
  background: var(--color-fill-1);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;

  .dns-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 4px;
    font-size: 13px;
    .dns-label { color: var(--color-text-3); min-width: 50px; }
    code { background: var(--color-fill-2); padding: 1px 6px; border-radius: 4px; font-size: 12px; }
    .dns-val-long { word-break: break-all; max-width: 400px; display: inline-block; }
  }

  .dns-note { font-size: 12px; color: var(--color-text-3); margin-top: 4px; font-style: italic; }
}

.password-box {
  display: flex;
  align-items: center;
  background: var(--color-fill-1);
  border-radius: 8px;
  padding: 12px;
  code { font-size: 14px; word-break: break-all; }
}

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
