<script setup>
import {Refresh, Search, SwitchButton, TurnOff, Open } from "@element-plus/icons-vue";
import {computed, onMounted, shallowRef} from 'vue'
import {getDevices} from '~/api'
import {showInfo, windows_size} from '~/utils'

const loading = shallowRef(false)
const devicesValue = shallowRef([])
const tableHeight = computed(() => {
  return windows_size.value.height - 75
})
let devicesTemp = []

function loadDevices() {
  loading.value = true
  getDevices().then((devices) => {
    devicesTemp = devices
    doFilterInputKey()
  }).catch((e) => {
    showInfo(true, e.message)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(loadDevices)
const filterInputKey = shallowRef("")
const filterType = shallowRef("0")

function doFilterInputKey() {
  devicesValue.value = devicesTemp.filter((device) => {
    if (filterInputKey.value) {
      switch (filterType.value) {
        case '0':
          return device.device && device.device.indexOf(filterInputKey.value) > -1
        case "1":
          return device.type && device.type.indexOf(filterInputKey.value) > -1
        default:
          return device.state && device.state.indexOf(filterInputKey.value) > -1
      }
    }
    return true;
  })
}
</script>

<template>
  <div style="display:flex; flex-direction: row;align-items: center;">
    <el-input style="max-width: 400px" v-model:model-value="filterInputKey" @keydown.enter="doFilterInputKey" :disabled="loading" clearable @clear="doFilterInputKey">
      <template #prepend>
        <el-select style="width: 80px" v-model:model-value="filterType">
          <el-option value="0" label="名称"/>
          <el-option value="1" label="类型"/>
          <el-option value="2" label="状态"/>
        </el-select>
      </template>
      <template #append>
        <el-button type="text" @click="doFilterInputKey">
          <el-icon>
            <Search/>
          </el-icon>
        </el-button>
      </template>
    </el-input>
    <el-button style="margin-left: 10px; width: 80px" type="primary" :icon="Refresh" :disabled="loading"
               @click="loadDevices">
      刷新
    </el-button>
  </div>
  <el-table v-loading="loading" :data="devicesValue" :max-height="tableHeight">
    <el-table-column prop="device" label="名称" min-width="200">
      <template #default="scope">
        <el-tooltip v-if="scope.row.type === 'wifi'" :placement="'top'" content="扫描WIFI热点">
          <router-link :to="{ query: { device: scope.row.device }, path: '/wifi' }">
            {{scope.row.device}}
          </router-link>
        </el-tooltip>
      </template>
    </el-table-column>
    <el-table-column min-width="100px" prop="type" label="类型"/>
    <el-table-column min-width="160px" prop="hardware_addr" label="MAC"/>
    <el-table-column min-width="150px" prop="ipv4" label="IPV4">
      <template #default="scope">
        <p v-for="ipv4 in scope.row.ipv4" :key="ipv4">
          {{ ipv4 }}
        </p>
      </template>
    </el-table-column>
    <el-table-column min-width="270px" prop="ipv6" label="IPV6">
      <template #default="scope">
        <p v-for="ipv6 in scope.row.ipv6" :key="ipv6">
          {{ ipv6 }}
        </p>
      </template>
    </el-table-column>
    <el-table-column min-width="70px" prop="up" label="电源">
      <template #default="scope">
        <el-button size="small" :icon="SwitchButton" type="danger" v-if="scope.row.up">断开</el-button>
        <el-button size="small" :icon="SwitchButton" type="info" v-else>开启</el-button>
      </template>
    </el-table-column>
    <el-table-column min-width="150px" prop="state" label="状态">
      <template #default="scope">
        <el-button size="small" :icon="TurnOff" v-if="scope.row.state === 'connected'" type="danger">断开</el-button>
        <el-button  size="small" :icon="Open" v-else-if="scope.row.state === 'disconnected'" type="info">连接</el-button>
        <el-tag v-else type="warning">
          {{ scope.row.state }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>
</template>
