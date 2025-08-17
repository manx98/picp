<script setup>
import axios from 'axios'
import {computed, onBeforeUnmount, onMounted, ref, shallowRef, watch} from 'vue'
import { getDisplayConfig, setDisplayConfig } from '~/api/index.js'
import { showInfo } from '~/utils/index.js'

const defaultValue = {
  addr: '3c',
  bus: 1,
  enable: false,
  height: 64,
  vcc_state: 0,
  width: 128,
  screen_size: '128x64',
  status_interval: 1,
}
const old = ref({ ...defaultValue })
const data = ref({ ...defaultValue })
const screen_size = {
  '128x64': [128, 64],
  '64x48': [64, 48],
  '128x32': [128, 32],
  '96x16': [96, 16],
}

let lastReq = null

const loading = shallowRef(false)
const showEmpty = shallowRef(false)
const formRef = ref()

function getCfg() {
  loading.value = true
  lastReq = getDisplayConfig()
  lastReq.rsp.then((rsp) => {
    const value = {
      addr: rsp.addr.toString(16),
      bus: rsp.bus,
      enable: rsp.enable,
      vcc_state: rsp.vcc_state,
      screen_size: `${rsp.width}x${rsp.height}`,
      status_interval: rsp.status_interval,
    }
    Object.assign(old.value, value)
    Object.assign(data.value, value)
    if (formRef.value) {
      formRef.value.clearValidate()
    }
    showEmpty.value = false
  }).catch((err) => {
    if (axios.isCancel(err)) {
      return
    }
    showEmpty.value = true
    showInfo(true, err.message)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(getCfg)

function cancelRequest() {
  if (lastReq) {
    lastReq.cancel()
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}
onBeforeUnmount(cancelRequest)
const need_update = computed(() => {
  return data.value.addr !== old.value.addr
    || data.value.bus !== old.value.bus
    || data.value.enable !== old.value.enable
    || data.value.vcc_state !== old.value.vcc_state
    || data.value.screen_size !== old.value.screen_size
    || data.value.status_interval !== old.value.status_interval
})
function checkAddr(rule, value, callback) {
  if (/^[0-9a-f]+$/i.test(value)) {
    return callback()
  }
  else {
    callback(new Error('不是合法的16进制地址'))
  }
}
const formRules = shallowRef({
  addr: { trigger: 'blur', validator: checkAddr },
  screen_size: { tigger: 'blur', validator: (rule, value, callback) => {
    if (screen_size[value]) {
      callback()
    }
    else {
      callback(new Error('无效尺寸'))
    }
  } },
})
function resetForm() {
  formData.value = { ...old.value }
}

function submitForm(formEl) {
  formEl.validate((valid) => {
    if (valid) {
      const size = screen_size[data.value.screen_size]
      lastReq = setDisplayConfig({
        addr: Number.parseInt(data.value.addr, 16),
        bus: data.value.bus,
        enable: data.value.enable,
        vcc_state: data.value.vcc_state,
        width: size[0],
        height: size[1],
        status_interval: data.value.status_interval,
      })
      loading.value = true
      lastReq.rsp.then(() => {
        old.value = { ...data.value }
      }).catch((err) => {
        if (axios.isCancel(err)) {
          return
        }
        showInfo(true, err.message)
      }).finally(() => {
        loading.value = false
      })
    }
  })
}
</script>

<template>
  <el-empty v-show="showEmpty" v-loading="loading" description="加载失败">
    <template #default>
      <el-button type="primary" @click="getCfg">
        刷新
      </el-button>
    </template>
  </el-empty>
  <div v-show="!showEmpty" v-loading="loading" style="max-width: 300px; text-align: right">
    <el-form ref="formRef" :model="data" :rules="formRules" label-width="auto">
      <el-form-item label="启用" prop="enable">
        <el-checkbox v-model="data.enable" />
      </el-form-item>
      <el-form-item label="地址" prop="addr">
        <el-input v-model.trim="data.addr" :maxlength="2">
          <template #prefix>
            0x
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="BUS" prop="bus">
        <el-select v-model="data.bus">
          <el-option :value="0" label="BUS-0" />
          <el-option :value="1" label="BUS-1" />
        </el-select>
      </el-form-item>
      <el-form-item label="宽度" prop="screen_size">
        <el-select v-model="data.screen_size">
          <el-option value="64x48" />
          <el-option value="96x16" />
          <el-option value="128x64" />
          <el-option value="128x32" />
        </el-select>
      </el-form-item>
      <el-form-item label="VCC" prop="vcc_state">
        <el-select v-model="data.vcc_state">
          <el-option :value="0" label="External" />
          <el-option :value="1" label="SwitchCAP" />
        </el-select>
      </el-form-item>
      <el-form-item label="刷新间隔" prop="status_interval">
        <el-input-number v-model="data.status_interval" :min="1">
          <template #prefix>
            秒
          </template>
        </el-input-number>
      </el-form-item>
    </el-form>
    <el-button type="primary" :disabled="!need_update" @click="submitForm(formRef)">
      应用
    </el-button>
    <el-button :disabled="!need_update" @click="resetForm">
      重置
    </el-button>
  </div>
</template>

<style scoped>

</style>
