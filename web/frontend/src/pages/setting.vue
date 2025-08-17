<script setup>
import { shallowRef } from 'vue'
import DisplayPanel from '~/components/DisplayPanel.vue'
import FanPanel from '~/components/FanPanel.vue'
import LoginSetting from '~/components/LoginSetting.vue'
import WifiApPanel from '~/components/WifiApPanel.vue'

const current_tab = shallowRef('fan')
</script>

<template>
  <el-tabs v-model="current_tab" style="padding: 20px" type="card">
    <el-tab-pane label="扇热" name="fan">
      <FanPanel v-if="current_tab === 'fan'" />
    </el-tab-pane>
    <el-tab-pane label="网络" name="wifi">
      <WifiApPanel v-if="current_tab === 'wifi'" />
    </el-tab-pane>
    <el-tab-pane label="显示" name="display">
      <DisplayPanel v-if="current_tab === 'display'" />
    </el-tab-pane>
    <el-tab-pane label="登陆" name="login">
      <LoginSetting v-if="current_tab === 'login'" />
    </el-tab-pane>
  </el-tabs>
</template>

<style scoped>

</style>
