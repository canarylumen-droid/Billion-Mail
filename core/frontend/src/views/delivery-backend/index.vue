<template>
  <div class="p-24px">
    <div class="flex items-center justify-between mb-24px">
      <div>
        <div class="bt-title">Delivery Backend</div>
        <div class="text-desc text-13px mt-4px">
          Configure how BillionMail delivers email — zero-cost options with day-1 inbox placement
        </div>
      </div>
      <n-button type="primary" @click="showAdd = true">
        <template #icon><n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></n-icon></template>
        Add Backend
      </n-button>
    </div>

    <!-- Free Options Banner -->
    <div class="options-banner mb-24px">
      <div class="option-card cf" @click="selectedType = 'cloudflare_worker'; showAdd = true">
        <div class="opt-icon">☁️</div>
        <div class="opt-body">
          <div class="opt-title">Cloudflare Workers</div>
          <div class="opt-desc">Free · 100k sends/day · Cloudflare's trusted IPs · No port 25 needed</div>
          <div class="opt-badge free">FREE FOREVER</div>
        </div>
        <div class="opt-arrow">→</div>
      </div>
      <div class="option-card oracle" @click="selectedType = 'oracle_vm'; showAdd = true">
        <div class="opt-icon">🖥️</div>
        <div class="opt-body">
          <div class="opt-title">Oracle Cloud Free VM</div>
          <div class="opt-desc">Free · Unlimited sends · Your own IPs · Port 25 via request · Haraka MTA</div>
          <div class="opt-badge free">FREE FOREVER</div>
        </div>
        <div class="opt-arrow">→</div>
      </div>
      <div class="option-card hetzner" @click="selectedType = 'haraka'; showAdd = true">
        <div class="opt-icon">⚡</div>
        <div class="opt-body">
          <div class="opt-title">Hetzner / OVH VPS</div>
          <div class="opt-desc">€4.51/mo · Unlimited sends · Dedicated IPs · Port 25 open by default</div>
          <div class="opt-badge paid">€4.51/mo</div>
        </div>
        <div class="opt-arrow">→</div>
      </div>
    </div>

    <!-- Comparison Table -->
    <div class="compare-card mb-24px">
      <div class="compare-title mb-16px">📊 Which Backend Should You Use?</div>
      <table class="compare-table">
        <thead>
          <tr>
            <th>Feature</th>
            <th>Cloudflare Workers</th>
            <th>Oracle Free VM</th>
            <th>Hetzner VPS</th>
          </tr>
        </thead>
        <tbody>
          <tr><td>Cost</td><td class="good">Free forever</td><td class="good">Free forever</td><td class="ok">€4.51/mo</td></tr>
          <tr><td>Daily send limit</td><td class="ok">100k requests/day</td><td class="good">Unlimited</td><td class="good">Unlimited</td></tr>
          <tr><td>IP reputation (day 1)</td><td class="good">Cloudflare's trusted IPs ✅</td><td class="ok">New IP (needs warming)</td><td class="ok">New IP (needs warming)</td></tr>
          <tr><td>Port 25</td><td class="good">Not needed (Cloudflare handles)</td><td class="ok">Request to enable</td><td class="good">Open by default ✅</td></tr>
          <tr><td>Setup time</td><td class="good">5 minutes</td><td class="ok">30 minutes</td><td class="ok">30 minutes</td></tr>
          <tr><td>Your own infrastructure</td><td class="ok">Cloudflare network</td><td class="good">100% yours ✅</td><td class="good">100% yours ✅</td></tr>
          <tr><td>Best for</td><td>Start immediately</td><td>Scale, own everything</td><td>Serious volume</td></tr>
        </tbody>
      </table>
    </div>

    <!-- Configured Backends -->
    <div class="backends-section">
      <div class="section-title mb-16px">Configured Backends</div>
      <div v-if="loading" class="flex justify-center py-40px"><n-spin size="large" /></div>
      <div v-else-if="backends.length === 0" class="empty-state">
        <div class="empty-icon">🚀</div>
        <div class="empty-title">No delivery backends yet</div>
        <div class="empty-desc">Choose one of the free options above to get started</div>
      </div>
      <div v-else class="backend-list">
        <div v-for="b in backends" :key="b.id" class="backend-card" :class="{ 'is-default': b.is_default }">
          <div class="backend-header">
            <div class="backend-left">
              <div class="backend-icon">{{ backendIcon(b.backend_type) }}</div>
              <div>
                <div class="backend-name">{{ b.name }}</div>
                <div class="backend-type">{{ backendLabel(b.backend_type) }}</div>
              </div>
            </div>
            <div class="backend-right">
              <n-tag v-if="b.is_default" type="success" size="small" class="mr-8px">DEFAULT</n-tag>
              <n-tag :type="b.status === 'active' ? 'success' : 'error'" size="small">{{ b.status }}</n-tag>
            </div>
          </div>

          <!-- Config Preview -->
          <div class="backend-config">
            <div v-if="b.backend_type === 'cloudflare_worker'" class="config-row">
              <span class="config-label">Worker URL</span>
              <code>{{ b.config.worker_url || '—' }}</code>
            </div>
            <div v-if="b.backend_type === 'cloudflare_worker'" class="config-row">
              <span class="config-label">Secret</span>
              <code>{{ b.config.secret || '—' }}</code>
            </div>
            <div v-if="b.backend_type !== 'cloudflare_worker'" class="config-row">
              <span class="config-label">Host</span>
              <code>{{ b.config.host || '—' }}</code>
            </div>
            <div class="config-row">
              <span class="config-label">Daily Limit</span>
              <span>{{ b.daily_sent.toLocaleString() }} / {{ b.daily_limit.toLocaleString() }}</span>
            </div>
          </div>
          <n-progress
            type="line"
            :percentage="Math.min(100, Math.round(b.daily_sent / b.daily_limit * 100))"
            class="mb-12px"
          />

          <div class="backend-actions">
            <n-button size="small" :loading="testingId === b.id" @click="handleTest(b)">
              Test Connection
            </n-button>
            <n-button size="small" type="primary" v-if="!b.is_default" @click="handleSetDefault(b.id)">
              Set as Default
            </n-button>
            <n-popconfirm @positive-click="handleDelete(b.id)">
              <template #trigger>
                <n-button size="small" type="error">Remove</n-button>
              </template>
              Remove this delivery backend?
            </n-popconfirm>
          </div>

          <!-- Last error -->
          <div v-if="b.last_error" class="backend-error">⚠️ {{ b.last_error }}</div>

          <!-- Test result -->
          <div v-if="testResults[b.id]" :class="['test-result', testResults[b.id].ok ? 'test-ok' : 'test-fail']">
            {{ testResults[b.id].ok ? '✅' : '❌' }} {{ testResults[b.id].message }}
            <span v-if="testResults[b.id].latency_ms" class="text-desc ml-8px text-12px">{{ testResults[b.id].latency_ms }}ms</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Backend Modal -->
    <n-modal v-model:show="showAdd" title="Add Delivery Backend" preset="card" style="width:640px">
      <n-tabs v-model:value="selectedType" type="card" class="mb-20px">
        <n-tab name="cloudflare_worker">☁️ Cloudflare Workers</n-tab>
        <n-tab name="oracle_vm">🖥️ Oracle Free VM</n-tab>
        <n-tab name="haraka">⚡ Haraka / VPS</n-tab>
      </n-tabs>

      <!-- Cloudflare Workers Setup -->
      <div v-if="selectedType === 'cloudflare_worker'">
        <div class="setup-guide">
          <div class="guide-step">
            <div class="step-num">1</div>
            <div class="step-body">
              <div class="step-title">Enable Email Routing on your domain</div>
              <div class="step-desc">Go to <a href="https://dash.cloudflare.com" target="_blank" class="link">dash.cloudflare.com</a> → your domain → <strong>Email → Email Routing</strong> → Enable it</div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">2</div>
            <div class="step-body">
              <div class="step-title">Create a Cloudflare Worker</div>
              <div class="step-desc">Workers & Pages → Create Worker → paste the script below → Deploy</div>
              <div class="code-block">
                <div class="code-header">
                  <span>email-relay.js — paste this into Cloudflare Workers</span>
                  <n-button size="tiny" @click="copyWorkerScript">Copy Script</n-button>
                </div>
                <pre class="code-pre">{{ workerScriptPreview }}</pre>
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">3</div>
            <div class="step-body">
              <div class="step-title">Add "Send Email" binding + secret variable</div>
              <div class="step-desc">
                Worker → Settings → Bindings → Add → <strong>Send Email</strong> → Name it <code>EMAIL</code><br>
                Then: Settings → Environment Variables → Add <code>BILLIONMAIL_SECRET</code> = any secret string
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">4</div>
            <div class="step-body">
              <div class="step-title">Enter your Worker details below</div>
            </div>
          </div>
        </div>

        <n-form label-placement="top" class="mt-16px">
          <n-form-item label="Backend Name">
            <n-input v-model:value="form.name" placeholder="e.g. Cloudflare Worker EU" />
          </n-form-item>
          <n-form-item label="Worker URL">
            <n-input v-model:value="form.config.worker_url" placeholder="https://your-worker.your-subdomain.workers.dev" />
          </n-form-item>
          <n-form-item label="BILLIONMAIL_SECRET (same value you set in Worker env vars)">
            <n-input v-model:value="form.config.secret" type="password" show-password-on="click" placeholder="your-secret-string" />
          </n-form-item>
          <n-form-item label="Set as default delivery backend">
            <n-switch v-model:value="form.is_default" />
          </n-form-item>
        </n-form>
      </div>

      <!-- Oracle Free VM Setup -->
      <div v-if="selectedType === 'oracle_vm'">
        <div class="setup-guide">
          <div class="guide-step">
            <div class="step-num">1</div>
            <div class="step-body">
              <div class="step-title">Create Oracle Cloud Free Account</div>
              <div class="step-desc">
                Go to <a href="https://cloud.oracle.com/free" target="_blank" class="link">cloud.oracle.com/free</a> →
                Sign up (needs a credit card for verification, but Always Free VMs are never charged) →
                Create 2 AMD E2.1.Micro VMs (free forever)
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">2</div>
            <div class="step-body">
              <div class="step-title">Request Port 25 to be enabled</div>
              <div class="step-desc">
                Oracle Cloud Console → Support → Create Service Request → "Remove outbound SMTP restrictions" →
                Explain you're building a legitimate transactional email platform (usually approved in 24h)
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">3</div>
            <div class="step-body">
              <div class="step-title">Install Haraka MTA (copy-paste into your VM)</div>
              <div class="code-block">
                <div class="code-header">
                  <span>SSH into your Oracle VM, then run:</span>
                  <n-button size="tiny" @click="copyHarakaScript">Copy Commands</n-button>
                </div>
                <pre class="code-pre">{{ harakaInstallScript }}</pre>
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">4</div>
            <div class="step-body">
              <div class="step-title">Set PTR (Reverse DNS) record</div>
              <div class="step-desc">
                In Oracle Cloud → your VM → Network → Reserved IPs → Edit → Set reverse DNS to <code>mail.yourdomain.com</code><br>
                This is critical for inbox placement — it proves your IP belongs to your domain.
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">5</div>
            <div class="step-body">
              <div class="step-title">Enter your Oracle VM details below</div>
            </div>
          </div>
        </div>

        <n-form label-placement="top" class="mt-16px">
          <n-form-item label="Backend Name">
            <n-input v-model:value="form.name" placeholder="e.g. Oracle Free VM - US East" />
          </n-form-item>
          <n-form-item label="VM IP Address / Hostname">
            <n-input v-model:value="form.config.host" placeholder="e.g. 140.238.1.100 or mail.yourdomain.com" />
          </n-form-item>
          <n-form-item label="SMTP Submission Port">
            <n-input-number v-model:value="form.config.port" :default-value="587" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="Set as default delivery backend">
            <n-switch v-model:value="form.is_default" />
          </n-form-item>
        </n-form>
      </div>

      <!-- Haraka / VPS Setup -->
      <div v-if="selectedType === 'haraka'">
        <div class="setup-guide">
          <div class="guide-step">
            <div class="step-num">1</div>
            <div class="step-body">
              <div class="step-title">Get a VPS with port 25 open</div>
              <div class="step-desc">
                <strong>Hetzner CX22</strong> — €4.51/mo — port 25 open by default →
                <a href="https://www.hetzner.com/cloud" target="_blank" class="link">hetzner.com/cloud</a><br>
                <strong>OVH VPS SSD</strong> — €3.99/mo →
                <a href="https://www.ovhcloud.com/en/vps/" target="_blank" class="link">ovhcloud.com</a><br>
                <strong>Contabo VPS S</strong> — €4.99/mo →
                <a href="https://contabo.com" target="_blank" class="link">contabo.com</a>
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">2</div>
            <div class="step-body">
              <div class="step-title">Install Haraka MTA</div>
              <div class="code-block">
                <div class="code-header">
                  <span>Run on your VPS:</span>
                  <n-button size="tiny" @click="copyHarakaScript">Copy</n-button>
                </div>
                <pre class="code-pre">{{ harakaInstallScript }}</pre>
              </div>
            </div>
          </div>
          <div class="guide-step">
            <div class="step-num">3</div>
            <div class="step-body">
              <div class="step-title">Set PTR (reverse DNS) in your VPS control panel</div>
              <div class="step-desc">Set your IP's reverse DNS to <code>mail.yourdomain.com</code> — required for inbox placement</div>
            </div>
          </div>
        </div>

        <n-form label-placement="top" class="mt-16px">
          <n-form-item label="Backend Name">
            <n-input v-model:value="form.name" placeholder="e.g. Hetzner EU MTA" />
          </n-form-item>
          <n-form-item label="VPS IP / Hostname">
            <n-input v-model:value="form.config.host" placeholder="e.g. 94.130.1.100 or mail.yourdomain.com" />
          </n-form-item>
          <n-form-item label="SMTP Port">
            <n-input-number v-model:value="form.config.port" :default-value="587" :min="1" :max="65535" />
          </n-form-item>
          <n-form-item label="Authentication Username">
            <n-input v-model:value="form.config.username" placeholder="SMTP username" />
          </n-form-item>
          <n-form-item label="Authentication Password">
            <n-input v-model:value="form.config.password" type="password" show-password-on="click" />
          </n-form-item>
          <n-form-item label="Daily Send Limit">
            <n-input-number v-model:value="form.daily_limit" :default-value="100000" :min="1" />
          </n-form-item>
          <n-form-item label="Set as default delivery backend">
            <n-switch v-model:value="form.is_default" />
          </n-form-item>
        </n-form>
      </div>

      <template #footer>
        <div class="flex items-center justify-between">
          <n-button :loading="testingNew" @click="handleTestNew">
            Test Connection First
          </n-button>
          <div class="flex gap-8px">
            <n-button @click="showAdd = false">Cancel</n-button>
            <n-button type="primary" :loading="saving" @click="handleSave">Save Backend</n-button>
          </div>
        </div>
        <div v-if="newTestResult" :class="['test-result mt-12px', newTestResult.ok ? 'test-ok' : 'test-fail']">
          {{ newTestResult.ok ? '✅' : '❌' }} {{ newTestResult.message }}
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import {
  listBackends, createBackend, deleteBackend,
  testBackend, setDefaultBackend,
} from '@/api/modules/delivery-backend'

const message = useMessage()
const backends = ref<any[]>([])
const loading = ref(false)
const showAdd = ref(false)
const saving = ref(false)
const testingId = ref<number | null>(null)
const testResults = ref<Record<number, any>>({})
const testingNew = ref(false)
const newTestResult = ref<any>(null)

const selectedType = ref('cloudflare_worker')
const form = ref({
  name: '',
  is_default: false,
  daily_limit: 100000,
  config: {
    worker_url: '',
    secret: '',
    host: '',
    port: 587,
    username: '',
    password: '',
  },
})

watch(showAdd, (v) => {
  if (!v) {
    newTestResult.value = null
    form.value = {
      name: '', is_default: false, daily_limit: 100000,
      config: { worker_url: '', secret: '', host: '', port: 587, username: '', password: '' },
    }
  }
})

watch(selectedType, () => { newTestResult.value = null })

const workerScriptPreview = `// Paste into Cloudflare Workers → edit code
export default {
  async fetch(request, env) {
    if (request.method !== 'POST') return new Response('Method not allowed', { status: 405 })
    const auth = request.headers.get('X-BillionMail-Secret')
    if (auth !== env.BILLIONMAIL_SECRET) return new Response('Unauthorized', { status: 401 })
    const { to, from, subject, html, text } = await request.json()
    const msg = new EmailMessage(from, to, buildRaw(from, to, subject, html || text))
    await env.EMAIL.send(msg)
    return Response.json({ ok: true })
  }
}
// Full script: BillionMail → Delivery Backend → Copy Full Script`

const harakaInstallScript = `# Install Node.js (LTS)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo bash -
sudo apt-get install -y nodejs

# Install Haraka globally
sudo npm install -g Haraka

# Initialize Haraka config in /etc/haraka
sudo haraka -i /etc/haraka

# Configure outbound hostname (replace with your domain)
echo "mail.yourdomain.com" | sudo tee /etc/haraka/config/me

# Start Haraka (it listens on port 25 + 587)
sudo haraka -c /etc/haraka

# Enable on boot (systemd)
sudo systemctl enable haraka 2>/dev/null || true`

function copyWorkerScript() {
  const fullScript = `// BillionMail Cloudflare Email Relay Worker
// Full script available at: core/cloudflare-worker/email-relay.js in BillionMail repo
export default {
  async fetch(request, env) {
    if (request.method !== 'POST') return new Response('Method not allowed', { status: 405 })
    const auth = request.headers.get('X-BillionMail-Secret')
    if (!env.BILLIONMAIL_SECRET || auth !== env.BILLIONMAIL_SECRET) return new Response('Unauthorized', { status: 401 })
    const payload = await request.json()
    const { to, from, from_name, subject, html, text, message_id } = payload
    const msgId = message_id || \`<\${Date.now()}.\${Math.random().toString(36).slice(2)}@billionmail.worker>\`
    const hdrs = [\`From: \${from_name ? \`"\${from_name}" <\${from}>\` : from}\`, \`To: \${to}\`, \`Subject: \${subject}\`, \`Message-ID: \${msgId}\`, \`Date: \${new Date().toUTCString()}\`, 'MIME-Version: 1.0', 'Content-Type: text/html; charset=UTF-8']
    const raw = [...hdrs, '', html || text || ''].join('\\r\\n')
    if (!env.EMAIL) return Response.json({ ok: false, error: 'Add EMAIL send_email binding in Worker settings' }, { status: 503 })
    await env.EMAIL.send(new EmailMessage(from, Array.isArray(to) ? to[0] : to, raw))
    return Response.json({ ok: true, message_id: msgId })
  }
}`
  navigator.clipboard.writeText(fullScript)
  message.success('Worker script copied!')
}

function copyHarakaScript() {
  navigator.clipboard.writeText(harakaInstallScript)
  message.success('Install commands copied!')
}

function backendIcon(type: string) {
  const icons: Record<string, string> = {
    cloudflare_worker: '☁️',
    oracle_vm: '🖥️',
    haraka: '⚡',
    smtp_relay: '📧',
  }
  return icons[type] || '📧'
}

function backendLabel(type: string) {
  const labels: Record<string, string> = {
    cloudflare_worker: 'Cloudflare Workers — Free, Cloudflare trusted IPs',
    oracle_vm: 'Oracle Cloud Free VM — Haraka MTA',
    haraka: 'Haraka MTA on VPS',
    smtp_relay: 'Custom SMTP Relay',
  }
  return labels[type] || type
}

async function fetchBackends() {
  loading.value = true
  try {
    const res = await listBackends() as any
    backends.value = res?.data || []
  } catch (_e) { /* ignore */ } finally {
    loading.value = false
  }
}

function buildTestPayload() {
  const cfg: any = { ...form.value.config }
  if (selectedType.value === 'cloudflare_worker') {
    return { backend_type: 'cloudflare_worker', config: { worker_url: cfg.worker_url, secret: cfg.secret } }
  }
  return { backend_type: selectedType.value, config: { host: cfg.host, port: cfg.port } }
}

async function handleTestNew() {
  testingNew.value = true
  newTestResult.value = null
  try {
    const res = await testBackend(buildTestPayload()) as any
    newTestResult.value = res?.data || { ok: false, message: 'No response' }
  } catch (_e) {
    newTestResult.value = { ok: false, message: 'Test request failed' }
  } finally {
    testingNew.value = false
  }
}

async function handleTest(b: any) {
  testingId.value = b.id
  delete testResults.value[b.id]
  try {
    const res = await testBackend({ backend_type: b.backend_type, config: b.config }) as any
    testResults.value[b.id] = res?.data || { ok: false, message: 'No response' }
  } catch (_e) {
    testResults.value[b.id] = { ok: false, message: 'Test request failed' }
  } finally {
    testingId.value = null
  }
}

async function handleSave() {
  if (!form.value.name) { message.error('Name is required'); return }
  if (selectedType.value === 'cloudflare_worker' && !form.value.config.worker_url) {
    message.error('Worker URL is required'); return
  }
  if (selectedType.value !== 'cloudflare_worker' && !form.value.config.host) {
    message.error('Host is required'); return
  }

  saving.value = true
  try {
    const config: any = {}
    if (selectedType.value === 'cloudflare_worker') {
      config.worker_url = form.value.config.worker_url
      config.secret = form.value.config.secret
    } else {
      config.host = form.value.config.host
      config.port = form.value.config.port || 587
      if (form.value.config.username) config.username = form.value.config.username
      if (form.value.config.password) config.password = form.value.config.password
    }

    await createBackend({
      name: form.value.name,
      backend_type: selectedType.value,
      config,
      daily_limit: form.value.daily_limit || 100000,
      is_default: form.value.is_default,
    })
    message.success('Delivery backend saved')
    showAdd.value = false
    fetchBackends()
  } catch (_e) {
    message.error('Failed to save')
  } finally {
    saving.value = false
  }
}

async function handleSetDefault(id: number) {
  try {
    await setDefaultBackend({ id })
    message.success('Default backend updated')
    fetchBackends()
  } catch (_e) { message.error('Failed') }
}

async function handleDelete(id: number) {
  try {
    await deleteBackend({ id })
    message.success('Backend removed')
    fetchBackends()
  } catch (_e) { message.error('Failed') }
}

onMounted(fetchBackends)
</script>

<style lang="scss" scoped>
.options-banner {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.option-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px;
  border-radius: 10px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;

  &:hover { transform: translateY(-2px); box-shadow: 0 4px 16px rgba(0,0,0,0.1); }
  &.cf { background: #f0f7ff; border-color: #bfd7f5; &:hover { border-color: #1890ff; } }
  &.oracle { background: #fff7e6; border-color: #ffd591; &:hover { border-color: #fa8c16; } }
  &.hetzner { background: #f6ffed; border-color: #b7eb8f; &:hover { border-color: #52c41a; } }

  .opt-icon { font-size: 28px; flex-shrink: 0; }
  .opt-body { flex: 1; }
  .opt-title { font-weight: 700; font-size: 15px; margin-bottom: 4px; }
  .opt-desc { font-size: 12px; color: #666; line-height: 1.4; margin-bottom: 6px; }
  .opt-badge {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 20px;
    font-size: 10px;
    font-weight: 700;
    &.free { background: #52c41a; color: white; }
    &.paid { background: #1890ff; color: white; }
  }
  .opt-arrow { font-size: 18px; color: #999; flex-shrink: 0; }
}

.compare-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  padding: 20px;
  overflow-x: auto;

  .compare-title { font-weight: 700; font-size: 15px; }
}

.compare-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;

  th {
    text-align: left;
    padding: 10px 12px;
    background: var(--color-fill-1);
    border-bottom: 2px solid var(--color-border-1);
    font-weight: 600;
  }

  td {
    padding: 10px 12px;
    border-bottom: 1px solid var(--color-border-1);
  }

  tr:last-child td { border-bottom: none; }

  .good { color: #52c41a; font-weight: 500; }
  .ok { color: #fa8c16; }
  .bad { color: #ff4d4f; }
}

.section-title { font-weight: 700; font-size: 16px; }

.backend-list { display: flex; flex-direction: column; gap: 16px; }

.backend-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 10px;
  padding: 20px;
  transition: box-shadow 0.2s;

  &.is-default {
    border-color: #52c41a;
    box-shadow: 0 0 0 1px #52c41a33;
  }

  .backend-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 14px;
  }

  .backend-left { display: flex; align-items: center; gap: 12px; }
  .backend-icon { font-size: 24px; }
  .backend-name { font-weight: 700; font-size: 15px; }
  .backend-type { font-size: 12px; color: var(--color-text-3); margin-top: 2px; }
  .backend-right { display: flex; align-items: center; }

  .backend-config {
    margin-bottom: 12px;
    .config-row {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 6px;
      font-size: 13px;
      .config-label { color: var(--color-text-3); min-width: 100px; }
      code { background: var(--color-fill-2); padding: 1px 6px; border-radius: 4px; font-size: 12px; word-break: break-all; }
    }
  }

  .backend-actions { display: flex; gap: 8px; margin-top: 12px; }
  .backend-error { margin-top: 8px; font-size: 12px; color: #ff4d4f; background: #fff2f0; padding: 6px 10px; border-radius: 6px; }
}

.test-result {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  &.test-ok { background: #f6ffed; color: #237804; border: 1px solid #b7eb8f; }
  &.test-fail { background: #fff2f0; color: #a8071a; border: 1px solid #ffa39e; }
}

.setup-guide { display: flex; flex-direction: column; gap: 16px; }

.guide-step {
  display: flex;
  gap: 14px;

  .step-num {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--color-primary);
    color: white;
    font-weight: 700;
    font-size: 13px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    margin-top: 2px;
  }

  .step-title { font-weight: 600; font-size: 14px; margin-bottom: 4px; }
  .step-desc { font-size: 13px; color: var(--color-text-3); line-height: 1.6; }
  code { background: var(--color-fill-2); padding: 1px 5px; border-radius: 4px; font-size: 12px; }
  .link { color: var(--color-primary); text-decoration: underline; }
}

.code-block {
  margin-top: 8px;
  background: #1a1a2e;
  border-radius: 8px;
  overflow: hidden;

  .code-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #16213e;
    font-size: 11px;
    color: #aaa;
  }

  .code-pre {
    margin: 0;
    padding: 12px;
    font-size: 11px;
    color: #e2e8f0;
    overflow-x: auto;
    font-family: 'Courier New', monospace;
    line-height: 1.6;
    white-space: pre;
    max-height: 200px;
    overflow-y: auto;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  text-align: center;
  background: var(--color-bg-2);
  border: 1px dashed var(--color-border-1);
  border-radius: 10px;
  .empty-icon { font-size: 48px; margin-bottom: 12px; }
  .empty-title { font-size: 18px; font-weight: 600; margin-bottom: 6px; }
  .empty-desc { font-size: 14px; color: var(--color-text-3); }
}
</style>
