<script lang="ts" setup>
import { computed, ref } from 'vue'
import { NButton, useMessage } from 'naive-ui'
import { useUserStore } from '@/store'
import type { UserInfo } from '@/store/modules/user/helper'
import { t } from '@/locales'

const userStore = useUserStore()

const ms = useMessage()

const userInfo = computed(() => userStore.userInfo)

const name = ref(userInfo.value.name ?? '')

const description = ref(userInfo.value.description ?? '')

function updateUserInfo(options: Partial<UserInfo>) {
  userStore.updateUserInfo(options)
  ms.success(t('common.success'))
}
</script>

<template>
  <div class="p-4 space-y-5 min-h-[200px]">
    <div class="space-y-6">
      <div class="p-4 space-y-5 min-h-[200px]">
        <img src="../../../icons/wechat.jpg" alt="404" class="w-64 h-64">
        <p>关注公众号，发送邮箱地址获取体验天数</p>
        <p>或微信搜索《掘金高效率工具计划》私信留言</p>
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]" />
        <div class="w-[200px]" />
        <NButton size="tiny" text type="primary" @click="updateUserInfo({ name })" />
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]" />
        <div class="flex-1" />
        <NButton size="tiny" text type="primary" @click="updateUserInfo({ description })" />
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]" />
        <div class="flex flex-wrap items-center gap-4" />
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]" />
        <div class="w-[200px]" />
      </div>
    </div>
  </div>
</template>

<!-- <template>
  <div class="p-4 space-y-5 min-h-[200px]">
		<img src="../../../icons/wechat.jpg" alt="404">
		<p>关注公众号，发送邮箱地址获取更多体验次数</p>
  </div>
</template> -->
