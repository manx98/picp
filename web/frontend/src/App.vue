<script setup>
import { Connection, Expand, Setting } from '@element-plus/icons-vue'
import { onMounted, ref, shallowRef, watch } from 'vue'
import { useRouter } from 'vue-router'
import { windows_size } from '~/utils/index.js'

const router = useRouter()

const current = ref('/')

onMounted(() => {
  if (router.currentRoute.value.name === '/setting') {
    current.value = router.currentRoute.value.name
  }
})

const is_mobile = shallowRef(true)
watch(windows_size, (newVal) => {
  is_mobile.value = newVal.width < 600
}, {
  immediate: true,
})
const showDrawer = shallowRef(false)
</script>

<template>
  <el-config-provider namespace="ep">
    <el-container>
      <el-affix v-if="is_mobile" offset="200">
        <el-button style="opacity: 0.7" size="small" type="primary" @click="showDrawer = true">
          <el-icon><Expand /></el-icon>
        </el-button>
      </el-affix>
      <el-drawer
        v-if="is_mobile"
        v-model:model-value="showDrawer"
        :with-header="false"
        direction="ltr"
        size="150px"
      >
        <el-menu router :default-active="current" @click="showDrawer = false">
          <el-menu-item index="/">
            <el-icon><Connection /></el-icon>
            <template #title>
              设备
            </template>
          </el-menu-item>
          <el-menu-item index="/setting">
            <el-icon><Setting /></el-icon>
            <template #title>
              设置
            </template>
          </el-menu-item>
        </el-menu>
      </el-drawer>
      <el-aside v-else width="150px">
        <el-scrollbar>
          <el-menu router style="height:100vh" :default-active="current">
            <el-menu-item index="/">
              <el-icon><Connection /></el-icon>
              <template #title>
                设备
              </template>
            </el-menu-item>
            <el-menu-item index="/setting">
              <el-icon><Setting /></el-icon>
              <template #title>
                设置
              </template>
            </el-menu-item>
          </el-menu>
        </el-scrollbar>
      </el-aside>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-config-provider>
</template>
