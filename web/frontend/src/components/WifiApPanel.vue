<script setup>
import axios from 'axios'
import { computed, defineProps, onBeforeUnmount, ref, shallowRef, watch } from 'vue'
import { getWifiConfig, setWifiConfig } from '~/api/index.js'
import { showInfo } from '~/utils/index.js'

const props = defineProps({
  show: Boolean,
})

let lastReq = null

const defaultValue = {
  device_name: '',
  enable: false,
  password: '',
  pin: 25,
  ssid: '',
}
const old = ref({ ...defaultValue })
const formData = ref({ ...defaultValue })
const showEmpty = shallowRef(false)
const loading = shallowRef(false)
const formRef = ref()
function getCfg() {
  loading.value = true
  lastReq = getWifiConfig()
  lastReq.rsp.then((data) => {
    Object.assign(old.value, data)
    Object.assign(formData.value, data)
    showEmpty.value = false
    if (formRef.value) {
      formRef.value.clearValidate()
    }
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
  formRef.value.validate((valid) => {
    if(valid) {
      loading.value = true
      lastReq = setWifiConfig({
        device_name: formData.value.device_name,
        enable: formData.value.enable,
        password: formData.value.password,
        pin: formData.value.pin,
        ssid: formData.value.ssid,
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
  })
}
const need_update = computed(() => {
  return formData.value.device_name !== old.value.device_name
    || formData.value.enable !== old.value.enable
    || formData.value.password !== old.value.password
    || formData.value.pin !== old.value.pin
    || formData.value.ssid !== old.value.ssid
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
    <el-form
      ref="formRef" v-model:model="formData" :rules="{
        ssid: { required: true, message: 'SSID不能为空', trigger: 'blur' },
        password: {trigger: 'blur', validator: (rule, value, callback) => {
          if(value === '' || value.length >= 8) {
            callback()
          } else {
            callback(new Error('密码长度至少为8位'))
          }
        }}
      }"
    >
      <el-form-item label="启用" prop="enabled">
        <el-checkbox v-model="formData.enable" />
      </el-form-item>
      <el-form-item label="SSID" prop="ssid" style="padding-bottom: 10px">
        <el-input v-model="formData.ssid" />
      </el-form-item>
      <el-form-item label="密码" prop="password" style="padding-bottom: 10px">
        <el-input v-model="formData.password" show-password />
      </el-form-item>
      <el-form-item label="引脚" prop="pin">
        <PinSelector v-model:pin="formData.pin" />
      </el-form-item>
    </el-form>
    <el-button type="primary" :disabled="!need_update" @click="submitForm">
      应用
    </el-button>
    <el-button @click="resetForm">
      重置
    </el-button>
  </div>
</template>

<style scoped>

</style>
