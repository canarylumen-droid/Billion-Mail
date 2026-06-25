<template>
  <div class="step-card" :class="`step-${step.type}`">
    <div class="step-header">
      <div class="step-badge">
        <span class="step-num">{{ index + 1 }}</span>
      </div>
      <div class="step-title-area">
        <div class="step-type-label">{{ typeLabel }}</div>
        <div class="step-desc">{{ typeDesc }}</div>
      </div>
      <div v-if="showDelay" class="step-delay-badge">
        After {{ step.delay_days }} day{{ step.delay_days !== 1 ? 's' : '' }}
      </div>
    </div>

    <div class="step-body">
      <n-grid :cols="2" :x-gap="16">
        <n-form-item-gi :span="1" label="From Address">
          <n-input v-model:value="step.addresser" placeholder="sender@yourdomain.com" />
        </n-form-item-gi>
        <n-form-item-gi :span="1" label="Display Name">
          <n-input v-model:value="step.full_name" placeholder="John Smith" />
        </n-form-item-gi>
      </n-grid>

      <n-form-item label="Subject Line">
        <div class="flex-1 flex gap-8px">
          <n-input
            v-model:value="step.subject"
            :placeholder="subjectPlaceholder"
            class="flex-1" />
          <n-popover trigger="click" placement="bottom-end">
            <template #trigger>
              <n-button size="small">{'{ }'} Variables</n-button>
            </template>
            <div class="variable-list">
              <div class="var-title">Available Variables</div>
              <div
                v-for="v in availableVars"
                :key="v"
                class="var-item"
                @click="insertVar('subject', v)">
                {{ v }}
              </div>
            </div>
          </n-popover>
        </div>
      </n-form-item>

      <n-form-item label="Email Template">
        <div class="flex-1 flex gap-8px items-center">
          <n-select
            v-model:value="step.template_id"
            :options="templateOptions"
            placeholder="Select template..."
            class="flex-1" />
          <n-button size="small" @click="$emit('createTemplate', step.type)">
            + New
          </n-button>
        </div>
      </n-form-item>

      <n-form-item v-if="showDelay" label="Delay (days after previous step)">
        <div class="flex items-center gap-12px">
          <n-input-number
            v-model:value="step.delay_days"
            :min="1"
            :max="90"
            style="width: 120px" />
          <span class="text-desc text-12px">
            days after {{ index === 0 ? 'sequence start' : 'previous step' }}
          </span>
        </div>
      </n-form-item>

      <div v-if="step.type === 'reply'" class="ai-settings">
        <n-form-item label="AI Auto-Reply">
          <div class="flex items-center gap-12px flex-1">
            <n-switch
              v-model:value="aiEnabledBool"
              @update:value="v => step.ai_enabled = v ? 1 : 0" />
            <span class="text-desc text-12px">
              {{ step.ai_enabled ? 'AI will auto-reply when lead responds' : 'Manual reply only' }}
            </span>
          </div>
        </n-form-item>
        <n-form-item v-if="step.ai_enabled" label="AI Tone">
          <n-select
            v-model:value="step.ai_tone"
            :options="toneOptions"
            style="max-width: 200px" />
        </n-form-item>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { NGrid, NFormItemGi, NInput, NSelect, NInputNumber, NSwitch, NButton, NPopover, NFormItem } from 'naive-ui'
import type { SequenceStep } from '@/api/modules/sequences'

const props = defineProps<{
  step: SequenceStep
  index: number
  templateOptions: { label: string; value: number }[]
  availableVars?: string[]
}>()

defineEmits<{ createTemplate: [type: string] }>()

const typeMap = {
  initial: { label: 'Initial Email', desc: 'First email sent to the lead', subject: 'Hey {{first_name}}, quick question about {{company}}' },
  followup1: { label: 'Follow-up 1', desc: 'Sent if no reply to initial', subject: 'Re: Hey {{first_name}}' },
  followup2: { label: 'Follow-up 2', desc: 'Final follow-up before closing', subject: 'Last touch — {{first_name}}' },
  reply: { label: 'Reply Handler', desc: 'Triggered when lead replies — AI responds', subject: 'Re: {{subject}}' },
}

const typeLabel = computed(() => typeMap[props.step.type]?.label || props.step.type)
const typeDesc = computed(() => typeMap[props.step.type]?.desc || '')
const subjectPlaceholder = computed(() => typeMap[props.step.type]?.subject || 'Subject line...')
const showDelay = computed(() => props.step.type !== 'initial')

const aiEnabledBool = computed({
  get: () => props.step.ai_enabled === 1,
  set: (v: boolean) => { props.step.ai_enabled = v ? 1 : 0 },
})

const toneOptions = [
  { label: 'Professional', value: 'professional' },
  { label: 'Friendly', value: 'friendly' },
  { label: 'Concise', value: 'concise' },
  { label: 'Enthusiastic', value: 'enthusiastic' },
  { label: 'Formal', value: 'formal' },
]

function insertVar(field: string, varName: string) {
  props.step.subject = (props.step.subject || '') + varName
}
</script>

<style scoped>
.step-card {
  border: 1px solid var(--n-border-color, #e5e7eb);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 8px;
}
.step-initial { border-top: 3px solid #1677ff; }
.step-followup1 { border-top: 3px solid #fa8c16; }
.step-followup2 { border-top: 3px solid #eb2f96; }
.step-reply { border-top: 3px solid #52c41a; }

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: var(--n-color-hover, #fafafa);
  border-bottom: 1px solid var(--n-border-color, #f0f0f0);
}
.step-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.step-num { color: #fff; font-weight: 700; font-size: 14px; }
.step-title-area { flex: 1; }
.step-type-label { font-weight: 600; font-size: 14px; color: var(--n-text-color); }
.step-desc { font-size: 12px; color: var(--n-text-color-3); margin-top: 2px; }
.step-delay-badge {
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 20px;
  padding: 4px 12px;
  font-size: 12px;
  color: #fa8c16;
  font-weight: 500;
}
.step-body {
  padding: 20px;
}
.ai-settings {
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 8px;
  padding: 16px;
  margin-top: 8px;
}
.variable-list {
  min-width: 200px;
}
.var-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--n-text-color-3);
  text-transform: uppercase;
  margin-bottom: 8px;
}
.var-item {
  padding: 6px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
  cursor: pointer;
  color: #1677ff;
}
.var-item:hover { background: #e8f4fd; }
</style>
