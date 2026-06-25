<template>
  <n-card title="Schedule Settings" class="mb-24px">
    <div class="schedule-grid">
      <!-- Working Days -->
      <div class="setting-group">
        <div class="setting-label">Working Days</div>
        <div class="days-row">
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="day-toggle"
            :class="{ active: schedule.working_days.includes(day.value) }"
            @click="toggleDay(day.value)">
            {{ day.short }}
          </div>
        </div>
        <div class="setting-hint">Click to select/deselect. Emails won't send on deselected days.</div>
      </div>

      <!-- Send Window -->
      <div class="setting-group">
        <div class="setting-label">Send Window</div>
        <div class="flex items-center gap-12px">
          <n-time-picker
            v-model:value="startTimeMs"
            format="HH:mm"
            :minutes-step="15"
            style="width: 110px"
            @update:value="updateStartTime" />
          <span class="text-desc">to</span>
          <n-time-picker
            v-model:value="endTimeMs"
            format="HH:mm"
            :minutes-step="15"
            style="width: 110px"
            @update:value="updateEndTime" />
        </div>
        <div class="setting-hint">Emails send only within this window (local server time)</div>
      </div>

      <!-- Timezone -->
      <div class="setting-group">
        <div class="setting-label">Timezone</div>
        <n-select
          v-model:value="schedule.timezone"
          :options="timezoneOptions"
          filterable
          style="max-width: 280px" />
      </div>

      <!-- Daily Limit -->
      <div class="setting-group">
        <div class="setting-label">Daily Limit Per Mailbox</div>
        <n-input-number
          v-model:value="schedule.daily_limit"
          :min="1"
          :max="10000"
          style="width: 160px" />
        <div class="setting-hint">Max emails sent per mailbox per day</div>
      </div>

      <!-- Randomize Window -->
      <div class="setting-group">
        <div class="setting-label">Randomize Send Window</div>
        <div class="flex items-center gap-12px">
          <n-input-number
            v-model:value="schedule.randomize_minutes"
            :min="0"
            :max="60"
            style="width: 120px" />
          <span class="text-desc">minutes variance</span>
        </div>
        <div class="setting-hint">Adds random delay (0 = disabled) to avoid spam detection</div>
      </div>

      <!-- Stop on Reply -->
      <div class="setting-group">
        <div class="setting-label">Stop on Reply</div>
        <div class="flex items-center gap-12px">
          <n-switch
            v-model:value="stopOnReplyBool"
            @update:value="v => schedule.stop_on_reply = v ? 1 : 0" />
          <span class="text-desc">
            {{ schedule.stop_on_reply ? 'Sequence stops when lead replies' : 'Sequence continues after reply' }}
          </span>
        </div>
      </div>

      <!-- Tracking -->
      <div class="setting-group">
        <div class="setting-label">Tracking</div>
        <div class="flex items-center gap-24px">
          <div class="flex items-center gap-8px">
            <n-switch
              v-model:value="trackOpenBool"
              @update:value="v => schedule.track_open = v ? 1 : 0" />
            <span class="text-desc">Track Opens</span>
          </div>
          <div class="flex items-center gap-8px">
            <n-switch
              v-model:value="trackClickBool"
              @update:value="v => schedule.track_click = v ? 1 : 0" />
            <span class="text-desc">Track Clicks</span>
          </div>
        </div>
      </div>
    </div>
  </n-card>
</template>

<script lang="ts" setup>
import { NCard, NSwitch, NSelect, NInputNumber, NTimePicker } from 'naive-ui'
import type { SequenceSchedule } from '@/api/modules/sequences'

const props = defineProps<{ schedule: SequenceSchedule }>()

const weekDays = [
  { value: 'mon', short: 'Mon', full: 'Monday' },
  { value: 'tue', short: 'Tue', full: 'Tuesday' },
  { value: 'wed', short: 'Wed', full: 'Wednesday' },
  { value: 'thu', short: 'Thu', full: 'Thursday' },
  { value: 'fri', short: 'Fri', full: 'Friday' },
  { value: 'sat', short: 'Sat', full: 'Saturday' },
  { value: 'sun', short: 'Sun', full: 'Sunday' },
]

const stopOnReplyBool = computed({
  get: () => props.schedule.stop_on_reply === 1,
  set: (v: boolean) => { props.schedule.stop_on_reply = v ? 1 : 0 },
})
const trackOpenBool = computed({
  get: () => props.schedule.track_open === 1,
  set: (v: boolean) => { props.schedule.track_open = v ? 1 : 0 },
})
const trackClickBool = computed({
  get: () => props.schedule.track_click === 1,
  set: (v: boolean) => { props.schedule.track_click = v ? 1 : 0 },
})

function timeStringToMs(t: string): number {
  const [h, m] = t.split(':').map(Number)
  return ((h || 0) * 3600 + (m || 0) * 60) * 1000
}

function msToTimeString(ms: number): string {
  const totalSec = ms / 1000
  const h = Math.floor(totalSec / 3600)
  const m = Math.floor((totalSec % 3600) / 60)
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
}

const startTimeMs = ref(timeStringToMs(props.schedule.send_window_start || '08:00'))
const endTimeMs = ref(timeStringToMs(props.schedule.send_window_end || '18:00'))

function updateStartTime(v: number | null) {
  if (v !== null) props.schedule.send_window_start = msToTimeString(v)
}
function updateEndTime(v: number | null) {
  if (v !== null) props.schedule.send_window_end = msToTimeString(v)
}

function toggleDay(day: string) {
  const idx = props.schedule.working_days.indexOf(day)
  if (idx >= 0) {
    props.schedule.working_days.splice(idx, 1)
  } else {
    props.schedule.working_days.push(day)
  }
}

const timezoneOptions = [
  { label: 'UTC', value: 'UTC' },
  { label: 'US/Eastern (EST)', value: 'America/New_York' },
  { label: 'US/Central (CST)', value: 'America/Chicago' },
  { label: 'US/Mountain (MST)', value: 'America/Denver' },
  { label: 'US/Pacific (PST)', value: 'America/Los_Angeles' },
  { label: 'Europe/London (GMT)', value: 'Europe/London' },
  { label: 'Europe/Paris (CET)', value: 'Europe/Paris' },
  { label: 'Europe/Berlin (CET)', value: 'Europe/Berlin' },
  { label: 'Asia/Dubai (GST)', value: 'Asia/Dubai' },
  { label: 'Asia/Kolkata (IST)', value: 'Asia/Kolkata' },
  { label: 'Asia/Singapore (SGT)', value: 'Asia/Singapore' },
  { label: 'Asia/Tokyo (JST)', value: 'Asia/Tokyo' },
  { label: 'Australia/Sydney (AEST)', value: 'Australia/Sydney' },
]
</script>

<style scoped>
.schedule-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px 40px;
}
.setting-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.setting-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--n-text-color);
}
.setting-hint {
  font-size: 11px;
  color: var(--n-text-color-3);
}
.days-row {
  display: flex;
  gap: 6px;
}
.day-toggle {
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid var(--n-border-color, #e5e7eb);
  color: var(--n-text-color-3);
  transition: all 0.15s;
  user-select: none;
}
.day-toggle.active {
  background: #1677ff;
  border-color: #1677ff;
  color: #fff;
}
</style>
