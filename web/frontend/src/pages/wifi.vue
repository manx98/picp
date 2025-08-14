<script setup>
import {computed, onBeforeUnmount, onMounted, ref, shallowRef} from 'vue'
import { MdRefresh, MdSearch } from 'vue-icons-plus/md'
import { useRouter } from 'vue-router'
import { createWifiAp, deleteWifiAp, listWifiAp } from '~/api'
import ConfirmDialog from '~/components/ConfirmDialog.vue'
import WifiConnectDialog from '~/components/WifiConnectDialog.vue'
import { showInfo, windows_size } from '~/utils'
import axios from 'axios'

const router = useRouter()

const loading = shallowRef(false)
const wifiList = shallowRef([])
const filterInputKey = shallowRef('')
let oldValue = []
let lastReq = null

onBeforeUnmount(() => {
  if (lastReq) {
    lastReq.cancel()
  }
})

function doFilterInputKey() {
  wifiList.value = oldValue.filter((item) => {
    if(filterInputKey.value && filterInputKey.value.length > 0) {
      if (item.SSID) {
        return item.SSID.indexOf(filterInputKey.value) !== -1
      }
      return false
    }
    return true
  })
}
function loadWifiAp() {
  loading.value = true
  lastReq = listWifiAp(router.currentRoute.value.query.device)
  lastReq.rsp.then((apList) => {
    apList.forEach((ap) => {
      if (ap.active) {
        ap.show_tag = true
        ap.tag_value = '已连接'
        ap.tag_type = 'success'
      }
      else if (ap.connection_uuid) {
        ap.show_tag = true
        ap.tag_value = '已保存'
        ap.tag_type = 'info'
      }
    })
    oldValue = apList
    doFilterInputKey()
  }).catch((err) => {
    if(axios.isCancel(err)) {
      return
    }
    showInfo(true, err.message)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(loadWifiAp)

const tableHeight = computed(() => {
  return windows_size.value.height - 82
})

const show_wifi_connect = shallowRef(false)

function configWifi({ config, done }) {
  const query = {
    device: config.device,
    auto_connect: config.auto_connect,
  }
  if (config.ssid && config.ssid.length > 0) {
    query.ssid = config.ssid
  }
  else {
    query.bssid = config.bssid
  }
  if (config.password && config.password.length > 0) {
    query.password = config.password.trim()
  }
  lastReq = createWifiAp(query)
  lastReq.rsp.then(() => {
    showInfo(false, '添加成功')
    show_wifi_connect.value = false
    loadWifiAp()
    done()
  }).catch((e) => {
    if(axios.isCancel(e)){
      return
    }
    done(e.message)
  })
}

const configInfo = ref({
  device: '',
  ssid: '',
  auto_connect: false,
  need_password: false,
  bssid: '',
  connection_uuid: '',
})

const showConfirmDialog = shallowRef(false)
const confirmDialogTitle = shallowRef('')
const confirmRow = shallowRef({})
const confirmType = shallowRef(0)

function changeConnectionState(connection_uuid, active, done) {
  lastReq = createWifiAp({
    connection_uuid,
    active,
  })
  lastReq.rsp.then(() => {
    showInfo(false, `连接${active ? '已建立' : '已断开'}`)
    showConfirmDialog.value = false
    loadWifiAp()
  }).catch((e) => {
    if(axios.isCancel(e)) {
      return
    }
    showInfo(true, e.message)
  }).finally(done)
}

function showConfigWifi(info) {
  if (info.connection_uuid) {
    showConnectConfirm(info)
    return
  }
  configInfo.value.auto_connect = false
  configInfo.value.device = info.device
  configInfo.value.ssid = info.SSID
  configInfo.value.bssid = info.BSSID
  configInfo.value.connection_uuid = info.connection_uuid
  configInfo.value.need_password = !!(info.security || info.connection_uuid)
  show_wifi_connect.value = true
}

function showDeleteConfirm(row) {
  confirmType.value = 0
  confirmRow.value = row
  showConfirmDialog.value = true
  confirmDialogTitle.value = '删除'
}

function doConfirm(done) {
  if (confirmType.value === 0) {
    lastReq=deleteWifiAp(confirmRow.value.connection_uuid)
    lastReq.rsp.then(() => {
      showConfirmDialog.value = false
      loadWifiAp()
    })
    lastReq.rsp.catch((e) => {
      if(axios.isCancel(e)) {
        return
      }
      showInfo(true, e.message)
    }).finally(done)
  }
  else if (confirmType.value === 1) {
    changeConnectionState(confirmRow.value.connection_uuid, true, done)
  }
  else {
    changeConnectionState(confirmRow.value.connection_uuid, false, done)
  }
}

function showConnectConfirm(row) {
  confirmRow.value = row
  confirmType.value = 1
  showConfirmDialog.value = true
  confirmDialogTitle.value = '连接'
}

function doShowDisconnectConfirm(row) {
  confirmRow.value = row
  confirmType.value = 2
  confirmDialogTitle.value = '断开'
  showConfirmDialog.value = true
}
</script>

<template>
  <ConfirmDialog v-model:show="showConfirmDialog" :title="confirmDialogTitle" @ok="doConfirm">
    <template #default>
      <span v-if="confirmType === 0">
        是否删除
      </span>
      <span v-if="confirmType === 1">
        是否连接
      </span>
      <span v-if="confirmType === 2">
        是否断开
      </span>
      <b style="color: #589ef8">{{ confirmRow.SSID }}({{ confirmRow.BSSID }})</b> ?
    </template>
  </ConfirmDialog>
  <WifiConnectDialog v-model:show="show_wifi_connect" :info="configInfo" @ok="configWifi" />
  <div style="display:flex; flex-direction: row;align-items: center;justify-items: center;padding: 20px 20px 0 20px">
    <el-input v-model:model-value="filterInputKey" style="max-width: 400px" @keydown.enter="doFilterInputKey" :disabled="loading" clearable @clear="doFilterInputKey">
      <template #prefix>
        <span>名称</span>
      </template>
      <template #append>
        <el-button type="text" :icon="MdSearch" @click="doFilterInputKey"/>
      </template>
    </el-input>
    <el-button style="margin-left: 10px;width: 100px" type="primary" :icon="MdRefresh" :disabled="loading" @click="loadWifiAp">
      刷新
    </el-button>
  </div>
  <el-table v-loading="loading" :data="wifiList" :height="tableHeight" style="padding: 0 20px 20px 20px">
    <el-table-column min-width="125px" label="SSID">
      <template #default="scope">
        <el-tag v-if="scope.row.show_tag" style="margin-right: 5px" :type="scope.row.tag_type" size="small">
          {{ scope.row.tag_value }}
        </el-tag>
        <b v-if="scope.row.show_tag">{{ scope.row.SSID }}</b>
        <span v-else>{{ scope.row.SSID }}</span>
      </template>
    </el-table-column>
    <el-table-column min-width="70px" prop="signal" label="信号" />
    <el-table-column min-width="150px" prop="BSSID" label="BSSID" />
    <el-table-column min-width="70px" prop="chan" label="信道" />
    <el-table-column min-width="100px" prop="freq" label="频率" />
    <el-table-column min-width="70px" prop="security" label="加密">
      <template #default="scope">
        <el-tag v-for="security in scope.row.security" :key="security" style="margin: 2px">
          {{ security }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column min-width="150px" label="操作">
      <template #default="scope">
        <el-button
          v-if="scope.row.connection_uuid" size="small" type="danger"
          @click="showDeleteConfirm(scope.row)"
        >
          删除
        </el-button>
        <el-button
          v-if="scope.row.active" size="small" type="warning"
          @click="doShowDisconnectConfirm(scope.row)"
        >
          断开
        </el-button>
        <el-button v-else size="small" type="primary" @click="showConfigWifi(scope.row)">
          连接
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped>

</style>
