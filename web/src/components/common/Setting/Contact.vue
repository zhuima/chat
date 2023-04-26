<script setup lang='ts'>
import { computed, ref } from 'vue'
import Wechat from './Wechat.vue'
import { SvgIcon } from '@/components/common'

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Props {
  visible: boolean
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
}

const active = ref('Wechat')

const show = computed({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emit('update:visible', visible)
  },
})

</script>

<template>
  <NModal v-model:show="show" :auto-focus="false">
    <NCard role="dialog" aria-modal="true" :bordered="false" style="width: 95%; max-width: 640px">
      <NTabs v-model:value="active" type="line" animated>
        <NTabPane name="General" tab="General">
          <template #tab>
            <SvgIcon class="text-lg" icon="ri:file-user-line" />
            <span class="ml-2">{{ $t('setting.contact') }}</span>
          </template>
          <div class="min-h-[100px]">
            <Wechat />
          </div>
        </NTabPane>
      </NTabs>
    </NCard>
  </NModal>
</template>
