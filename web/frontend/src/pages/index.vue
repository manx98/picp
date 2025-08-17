<script setup>
import axios from 'axios'
import { computed, onBeforeUnmount, onMounted, shallowRef } from 'vue'
import { MdLink, MdLinkOff, MdPower, MdPowerOff, MdRefresh, MdSearch } from 'vue-icons-plus/md'
import { getDevices } from '~/api'
import { showInfo, windows_size } from '~/utils'

const loading = shallowRef(false)
const devicesValue = shallowRef([])
const tableHeight = computed(() => {
  return windows_size.value.height - 75
})
let devicesTemp = []

let lastLoadDevicesQuery = null

onBeforeUnmount(() => {
  if (lastLoadDevicesQuery) {
    lastLoadDevicesQuery.cancel()
  }
})

function loadDevices() {
  loading.value = true
  lastLoadDevicesQuery = getDevices()
  lastLoadDevicesQuery.rsp.then((devices) => {
    devicesTemp = devices
    doFilterInputKey()
  }).catch((err) => {
    if (axios.isCancel(err)) {
      return
    }
    devicesTemp = []
    doFilterInputKey()
    showInfo(true, err.message)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(loadDevices)
const filterInputKey = shallowRef('')
const filterType = shallowRef('0')

function doFilterInputKey() {
  devicesValue.value = devicesTemp.filter((device) => {
    if (filterInputKey.value) {
      switch (filterType.value) {
        case '0':
          return device.device && device.device.includes(filterInputKey.value)
        case '1':
          return device.type && device.type.includes(filterInputKey.value)
        default:
          return device.state && device.state.includes(filterInputKey.value)
      }
    }
    return true
  })
}
</script>

<template>
  <div style="display:flex; flex-direction: row;align-items: center;padding: 20px 20px 0 20px">
    <el-input v-model:model-value="filterInputKey" style="max-width: 400px" :disabled="loading" clearable @keydown.enter="doFilterInputKey" @clear="doFilterInputKey">
      <template #prepend>
        <el-select v-model:model-value="filterType" style="width: 80px">
          <el-option value="0" label="名称" />
          <el-option value="1" label="类型" />
          <el-option value="2" label="状态" />
        </el-select>
      </template>
      <template #append>
        <el-button :icon="MdSearch" type="text" @click="doFilterInputKey" />
      </template>
    </el-input>
    <el-button
      style="margin-left: 10px; width: 80px" type="primary" :icon="MdRefresh" :disabled="loading"
      @click="loadDevices"
    >
      刷新
    </el-button>
  </div>
  <el-table v-loading="loading" :data="devicesValue" :height="tableHeight" style="padding: 0 20px 20px 20px;">
    <el-table-column prop="device" label="名称" min-width="200">
      <template #default="scope">
        <el-tooltip v-if="scope.row.type === 'wifi'" placement="top" content="扫描WIFI热点">
          <router-link :to="{ query: { device: scope.row.device }, path: '/wifi' }">
            {{ scope.row.device }}
          </router-link>
        </el-tooltip>
      </template>
    </el-table-column>
    <el-table-column min-width="100px" prop="type" label="类型" />
    <el-table-column min-width="160px" prop="hardware_addr" label="MAC" />
    <el-table-column min-width="150px" prop="ipv4" label="IPV4">
      <template #default="scope">
        <p v-for="ipv4 in scope.row.ipv4" :key="ipv4">
          {{ ipv4 }}
        </p>
      </template>
    </el-table-column>
    <el-table-column min-width="300px" prop="ipv6" label="IPV6">
      <template #default="scope">
        <p v-for="ipv6 in scope.row.ipv6" :key="ipv6">
          {{ ipv6 }}
        </p>
      </template>
    </el-table-column>
    <el-table-column min-width="100px" prop="up" label="电源">
      <template #default="scope">
        <el-button v-if="scope.row.up" :icon="MdPowerOff" size="small" type="danger">
          关闭
        </el-button>
        <el-button v-else :icon="MdPower" size="small" type="info">
          开启
        </el-button>
      </template>
    </el-table-column>
    <el-table-column min-width="200px" prop="state" label="状态">
      <template #default="scope">
        <el-button v-if="scope.row.state === 'connected'" :icon="MdLinkOff" size="small" type="danger">
          断开
        </el-button>
        <el-button v-else-if="scope.row.state === 'disconnected'" :icon="MdLink" size="small" type="info">
          连接
        </el-button>
        <el-tag v-else type="warning">
          {{ scope.row.state }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>
</template>
