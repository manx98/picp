<script setup>
import { computed, shallowRef, watch } from 'vue'
import { MdMoreHoriz, MdSettings, MdSettingsEthernet } from 'vue-icons-plus/md'
import { useRouter } from 'vue-router'
import { windows_size } from '~/utils/index.js'

const router = useRouter()

const showNavPath = new Set(['/', '/setting', '/wifi'])

const current = computed(() => {
  if (router.currentRoute) {
    if (router.currentRoute.value.name === '/setting') {
      return router.currentRoute.value.name
    }
  }
  return '/'
})

const hiddenNav = computed(() => {
  if (router.currentRoute) {
    return !showNavPath.has(router.currentRoute.value.name)
  }
  else {
    return false
  }
})

const is_mobile = shallowRef(true)
watch(windows_size, (newVal) => {
  is_mobile.value = newVal.width < 600
}, {
  immediate: true,
})
const showDrawer = shallowRef(false)
function logout() {
  window.location.href = '/logout'
}
</script>

<template>
  <el-config-provider namespace="ep">
    <el-container>
      <el-affix v-if="is_mobile && !hiddenNav" offset="200">
        <el-button style="opacity: 0.7" size="small" type="primary" @click="showDrawer = true">
          <MdMoreHoriz />
        </el-button>
      </el-affix>
      <el-drawer
        v-if="is_mobile && !hiddenNav"
        v-model:model-value="showDrawer"
        :with-header="false"
        direction="ltr"
        size="150px"
      >
        <el-menu router :default-active="current" @click="showDrawer = false">
          <el-menu-item index="/">
            <MdSettingsEthernet style="margin-right: 10px" />
            <template #title>
              设备
            </template>
          </el-menu-item>
          <el-menu-item index="/setting">
            <MdSettings style="margin-right: 10px" />
            <template #title>
              设置
            </template>
          </el-menu-item>
          <div style="text-align: center;margin-top: 50px">
            <el-button :disabled="false" type="primary" size="small" @click="logout">
              退出登陆
            </el-button>
          </div>
        </el-menu>
      </el-drawer>
      <el-aside v-if="!is_mobile && !hiddenNav" width="110px">
        <el-scrollbar>
          <el-menu router style="height:100vh" :default-active="current">
            <el-menu-item index="/">
              <MdSettingsEthernet style="margin-right: 10px" />
              <template #title>
                设备
              </template>
            </el-menu-item>
            <el-menu-item index="/setting">
              <MdSettings style="margin-right: 10px" />
              <template #title>
                设置
              </template>
            </el-menu-item>
            <div style="text-align: center;margin-top: 50px">
              <el-button :disabled="false" type="primary" size="small" @click="logout">
                退出登陆
              </el-button>
            </div>
          </el-menu>
        </el-scrollbar>
      </el-aside>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-config-provider>
</template>

<style scoped>
:deep(.ep-main) {
  padding: 0;
}
</style>
