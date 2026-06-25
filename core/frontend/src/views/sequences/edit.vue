<template>
  <div class="seq-edit-container">
    <div class="seq-edit-header">
      <n-breadcrumb>
        <n-breadcrumb-item>
          <router-link to="/sequences">Sequences</router-link>
        </n-breadcrumb-item>
        <n-breadcrumb-item>{{ isEdit ? form.name || 'Edit Sequence' : 'New Sequence' }}</n-breadcrumb-item>
      </n-breadcrumb>
      <div class="header-actions">
        <n-button @click="$router.back()">Cancel</n-button>
        <n-button type="default" @click="handleSave('draft')">Save Draft</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave('active')">
          Save & Activate
        </n-button>
      </div>
    </div>

    <div class="seq-edit-body">
      <!-- Left column: steps -->
      <div class="seq-steps-col">
        <!-- Sequence Info -->
        <n-card class="mb-24px" title="Sequence Info">
          <n-grid :cols="2" :x-gap="16">
            <n-form-item-gi :span="2" label="Sequence Name" required>
              <n-input v-model:value="form.name" placeholder="e.g. SaaS Founders Outreach Q1" />
            </n-form-item-gi>
            <n-form-item-gi :span="2" label="Description">
              <n-input v-model:value="form.description" type="textarea" :rows="2" placeholder="What is this sequence for?" />
            </n-form-item-gi>
            <n-form-item-gi :span="2" label="Calendly URL (optional)">
              <n-input v-model:value="form.calendly_url" placeholder="https://calendly.com/yourname/30min" />
              <template #feedback>
                <span class="text-11px text-desc">AI will include booking link in replies when applicable</span>
              </template>
            </n-form-item-gi>
          </n-grid>
        </n-card>

        <!-- Variable Hint -->
        <div v-if="availableVars.length" class="var-hint-bar mb-16px">
          <span class="var-hint-icon">📋</span>
          <span class="var-hint-text">Available variables from your leads:</span>
          <span
            v-for="v in availableVars"
            :key="v"
            class="var-chip"
            @click="copyVar(v)">
            {{ v }}
          </span>
        </div>
        <div v-else class="var-hint-bar mb-16px">
          <span class="var-hint-icon">💡</span>
          <span class="var-hint-text">Standard variables: <code>&#123;&#123;first_name&#125;&#125;</code> <code>&#123;&#123;last_name&#125;&#125;</code> <code>&#123;&#123;company&#125;&#125;</code> <code>&#123;&#123;title&#125;&#125;</code> <code>&#123;&#123;email&#125;&#125;</code></span>
        </div>

        <!-- Steps -->
        <div class="steps-header mb-16px">
          <span class="steps-title">Outreach Steps</span>
          <span class="steps-subtitle">{{ form.steps.length }} steps · sequence stops when lead replies</span>
        </div>

        <div class="steps-timeline">
          <div v-for="(step, i) in form.steps" :key="i" class="timeline-item">
            <div class="timeline-connector" v-if="i > 0">
              <div class="connector-line"></div>
              <div class="connector-label">Wait {{ step.delay_days }} day{{ step.delay_days !== 1 ? 's' : '' }}</div>
            </div>
            <step-card
              :step="step"
              :index="i"
              :template-options="templateOptions"
              :available-vars="availableVars"
              @create-template="handleCreateTemplate" />
          </div>
        </div>
      </div>

      <!-- Right column: schedule -->
      <div class="seq-schedule-col">
        <schedule-settings :schedule="form.schedule" />

        <!-- Quick Stats Preview -->
        <n-card title="Estimated Timeline" class="mb-24px">
          <div class="timeline-preview">
            <div v-for="(step, i) in form.steps" :key="i" class="timeline-row">
              <div class="tl-dot" :class="`dot-${step.type}`"></div>
              <div class="tl-info">
                <div class="tl-name">{{ stepLabels[step.type] }}</div>
                <div class="tl-day">
                  {{ i === 0 ? 'Day 1' : `Day ${form.steps.slice(0, i + 1).reduce((a, s) => a + (s.delay_days || 0), 1)}` }}
                </div>
              </div>
            </div>
          </div>
        </n-card>

        <!-- AI Settings -->
        <n-card title="AI Settings">
          <n-form-item label="AI Provider">
            <n-select v-model:value="aiProvider" :options="aiProviderOptions" />
          </n-form-item>
          <n-form-item label="Auto-Optimize Copy">
            <div class="flex items-center gap-12px">
              <n-switch v-model:value="autoOptimize" />
              <span class="text-desc text-12px">AI adjusts subject/copy based on open/reply analytics</span>
            </div>
          </n-form-item>
          <n-form-item label="Smart Warm-up">
            <div class="flex items-center gap-12px">
              <n-switch v-model:value="warmupEnabled" />
              <span class="text-desc text-12px">Gradually increase send volume to protect deliverability</span>
            </div>
          </n-form-item>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { NBreadcrumb, NBreadcrumbItem, NButton, NCard, NGrid, NFormItemGi, NInput, NFormItem, NSelect, NSwitch } from 'naive-ui'
import { Message } from '@/utils'
import { createSequence, updateSequence, getSequenceById } from '@/api/modules/sequences'
import type { Sequence, SequenceStep } from '@/api/modules/sequences'
import { getTemplateAll } from '@/api/modules/market/template'
import StepCard from './components/StepCard.vue'
import ScheduleSettings from './components/ScheduleSettings.vue'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.query.id)
const saving = ref(false)
const aiProvider = ref('openai')
const autoOptimize = ref(false)
const warmupEnabled = ref(true)
const templateOptions = ref<{ label: string; value: number }[]>([])

const availableVars = ref<string[]>([
  '{{first_name}}', '{{last_name}}', '{{company}}', '{{title}}', '{{email}}', '{{website}}', '{{phone}}'
])

const stepLabels: Record<string, string> = {
  initial: 'Initial Email',
  followup1: 'Follow-up 1',
  followup2: 'Follow-up 2',
  reply: 'Reply Handler',
}

const defaultStep = (type: SequenceStep['type'], delay: number = 0): SequenceStep => ({
  type,
  subject: '',
  template_id: 0,
  addresser: '',
  full_name: '',
  delay_days: delay,
  ai_enabled: type === 'reply' ? 1 : 0,
  ai_tone: 'professional',
})

const form = reactive<Sequence>({
  name: '',
  description: '',
  status: 'draft',
  calendly_url: '',
  steps: [
    defaultStep('initial', 0),
    defaultStep('followup1', 3),
    defaultStep('followup2', 7),
    defaultStep('reply', 1),
  ],
  schedule: {
    working_days: ['mon', 'tue', 'wed', 'thu', 'fri'],
    send_window_start: '08:00',
    send_window_end: '18:00',
    timezone: 'America/New_York',
    daily_limit: 200,
    randomize_minutes: 15,
    stop_on_reply: 1,
    track_open: 1,
    track_click: 1,
  },
})

const aiProviderOptions = [
  { label: 'OpenAI (GPT-4)', value: 'openai' },
  { label: 'Anthropic (Claude)', value: 'anthropic' },
  { label: 'DeepSeek', value: 'deepseek' },
  { label: 'Gemini', value: 'gemini' },
]

function copyVar(v: string) {
  navigator.clipboard?.writeText(v)
  Message.success(`Copied ${v}`)
}

function handleCreateTemplate(type: string) {
  router.push('/template/ai-template/new?type=' + type)
}

async function handleSave(status: 'draft' | 'active') {
  if (!form.name.trim()) {
    Message.error('Please enter a sequence name')
    return
  }
  saving.value = true
  form.status = status
  try {
    if (isEdit.value && route.query.id) {
      await updateSequence({ ...form, id: Number(route.query.id) })
    } else {
      await createSequence(form)
    }
    router.push('/sequences')
  } catch (_e) { /* ignore */ } finally {
    saving.value = false
  }
}

async function loadTemplates() {
  try {
    const res = await getTemplateAll() as any
    if (Array.isArray(res)) {
      templateOptions.value = res.map((t: any) => ({ label: t.temp_name || t.name, value: t.id }))
    }
  } catch (e) { void e }
}

async function loadSequence() {
  if (!route.query.id) return
  try {
    const res = await getSequenceById({ id: Number(route.query.id) }) as any
    if (res) {
      Object.assign(form, res)
    }
  } catch (e) { void e }
}

onMounted(() => {
  loadTemplates()
  loadSequence()
})
</script>

<style scoped>
.seq-edit-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}
.seq-edit-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid var(--n-border-color, #f0f0f0);
  background: var(--n-card-color, #fff);
  flex-shrink: 0;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.seq-edit-body {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 24px;
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}
.var-hint-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  padding: 10px 14px;
  background: var(--n-color-hover, #fafafa);
  border: 1px dashed var(--n-border-color, #d9d9d9);
  border-radius: 8px;
  font-size: 13px;
}
.var-hint-icon { font-size: 16px; }
.var-hint-text { color: var(--n-text-color-3); font-size: 12px; }
.var-chip {
  background: #e8f4fd;
  color: #1677ff;
  font-family: monospace;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  cursor: pointer;
}
.var-chip:hover { background: #bae0ff; }
.steps-header { display: flex; align-items: baseline; gap: 12px; }
.steps-title { font-size: 15px; font-weight: 600; color: var(--n-text-color); }
.steps-subtitle { font-size: 12px; color: var(--n-text-color-3); }
.timeline-item { position: relative; }
.timeline-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 4px 0;
}
.connector-line {
  width: 2px;
  height: 20px;
  background: var(--n-border-color, #e5e7eb);
}
.connector-label {
  font-size: 11px;
  color: var(--n-text-color-3);
  background: var(--n-base-color, #fff);
  padding: 2px 8px;
  border-radius: 10px;
  border: 1px dashed var(--n-border-color, #d9d9d9);
  margin: 2px 0;
}
.timeline-preview { display: flex; flex-direction: column; gap: 12px; }
.timeline-row { display: flex; align-items: center; gap: 10px; }
.tl-dot {
  width: 10px; height: 10px; border-radius: 50%;
  flex-shrink: 0;
}
.dot-initial { background: #1677ff; }
.dot-followup1 { background: #fa8c16; }
.dot-followup2 { background: #eb2f96; }
.dot-reply { background: #52c41a; }
.tl-name { font-size: 13px; font-weight: 500; color: var(--n-text-color); }
.tl-day { font-size: 11px; color: var(--n-text-color-3); }
</style>
