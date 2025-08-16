<script setup>
import axios from 'axios'
import { computed, defineProps, onBeforeUnmount, ref, shallowRef, watch } from 'vue'
import { getFanConfig, setFanConfig } from '~/api/index.js'
import { showInfo } from '~/utils/index.js'

const props = defineProps({
  show: Boolean,
})

let lastReq = null

const defaultValue = {
  enable: false,
  temp: [45, 50],
  pin: 14,
  speed: 60,
}
const old = ref({ ...defaultValue })
const formData = ref({ ...defaultValue })
const showEmpty = shallowRef(false)
const loading = shallowRef(false)
const formRef = ref()
function getCfg() {
  loading.value = true
  lastReq = getFanConfig()
  lastReq.rsp.then((data) => {
    data.temp = [data.min_temp, data.max_temp]
    Object.assign(old.value, data)
    Object.assign(formData.value, data)
    showEmpty.value = false
    if (formRef.value) {
      formRef.value.clearValidate()
    }
    showEmpty.value = false
  }).catch((err) => {
    if (axios.isCancel(err)) {
      return
    }
    showInfo(true, err.message)
    showEmpty.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.show, (value) => {
  if (value) {
    getCfg()
  }
  else {
    cancelRequest()
  }
}, {
  immediate: true,
})

function cancelRequest() {
  if (lastReq) {
    lastReq.cancel()
  }
  if(formRef.value) {
    formRef.value.clearValidate()
  }
}
onBeforeUnmount(cancelRequest)
function resetForm() {
  formData.value = { ...old.value }
}
function submitForm() {
  loading.value = true
  lastReq = setFanConfig({
    enable: formData.value.enable,
    min_temp: formData.value.temp[0],
    max_temp: formData.value.temp[1],
    pin: formData.value.pin,
    speed: formData.value.speed,
  })
  lastReq.rsp.then(() => {
    old.value = { ...formData.value }
  }).catch((err) => {
    if (axios.isCancel(err)) {
      return
    }
    showInfo(true, err.message)
  }).finally(() => {
    loading.value = false
  })
}
const can_update = computed(() => {
  return formData.value.enable !== old.value.enable
    || formData.value.temp[0] !== old.value.temp[0]
    || formData.value.temp[1] !== old.value.temp[1]
    || formData.value.pin !== old.value.pin
    || formData.value.speed !== old.value.speed
})
</script>

<template>
  <el-empty v-show="showEmpty" v-loading="loading" description="加载失败">
    <template #default>
      <el-button type="primary" @click="getCfg">
        刷新
      </el-button>
    </template>
  </el-empty>
  <div v-show="!showEmpty" v-loading="loading" style="max-width: 300px;text-align: center">
    <el-form ref="formRef" v-model:model="formData">
      <el-form-item label="启用" prop="enabled">
        <el-checkbox v-model="formData.enable" />
      </el-form-item>
      <el-form-item label="温度" prop="temp" style="padding-bottom: 10px">
        <el-slider
          v-model:model-value="formData.temp" :step="1" :min="30" :max="80" range :marks="{
            1: '1°C',
            30: '30°C',
            60: '60°C',
            80: '80°C',
          }"
        />
      </el-form-item>
      <el-form-item label="引脚" prop="pin">
        <el-select v-model="formData.pin">
          <el-option :value="12" label="GPIO 12" />
          <el-option :value="13" label="GPIO 13" />
          <el-option :value="40" label="GPIO 40" />
          <el-option :value="41" label="GPIO 41" />
          <el-option :value="45" label="GPIO 45" />
          <el-option :value="18" label="GPIO 18" />
          <el-option :value="19" label="GPIO 19" />
        </el-select>
      </el-form-item>
      <el-form-item label="速度" prop="speed">
        <el-slider v-model="formData.speed" :max="100" :min="60" />
      </el-form-item>
    </el-form>
    <el-button type="primary" :disabled="!can_update" @click="submitForm">
      应用
    </el-button>
    <el-button @click="resetForm">
      重置
    </el-button>
  </div>
</template>

<style scoped>

</style>
