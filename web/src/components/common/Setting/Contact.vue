<script setup lang='ts'>
import { computed, ref } from 'vue'
import { NCard, NModal, NTabPane, NTabs } from 'naive-ui'


import Admin from './Admin.vue'
import { SvgIcon } from '@/components/common'

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Props {
	visible: boolean
}

interface Emit {
	(e: 'update:visible', visible: boolean): void
}

const active = ref('General')

const show = computed({
	get() {
		return props.visible
	},
	set(visible: boolean) {
		emit('update:visible', visible)
	},
})

// hidden
const isAdminUser = ref<boolean>(false)
</script>

<template>
	<NModal v-model:show="show" :auto-focus="false">
		<NCard role="dialog" aria-modal="true" :bordered="false" style="width: 75%; max-width: 640px">
			<NTabs v-model:value="active" type="line" animated>
				<NTabPane name="General" tab="General">
					<template #tab>
						<SvgIcon class="text-lg" icon="ri:file-user-line" />
						<span class="ml-2">{{ $t('setting.contact') }}</span>
					</template>
					<div class="min-h-[100px]">
							  <div class="p-4 space-y-5 min-h-[200px]">
									<img src="../../../icons/wechat.jpg" alt="404">
									<p>公众号私信，发送邮箱地址获取更多体验次数</p>
								</div>
						<!-- <General /> -->
					</div>
				</NTabPane>
				<NTabPane v-if="isAdminUser" name="Admin" tab="Admin">
					<template #tab>
						<SvgIcon class="text-lg" icon="ri:list-settings-line" />
						<span class="ml-2">{{ $t('setting.admin') }}</span>
					</template>
					<Admin />
				</NTabPane>
			</NTabs>
		</NCard>
	</NModal>
</template>



<!-- <script setup lang='ts'>
import { computed, ref } from 'vue'
import Wechat from './Wechat.vue'
import { SvgIcon } from '@/components/common'

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

// interface Props {
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
</template> -->
