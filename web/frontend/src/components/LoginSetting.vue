<script setup>
import axios from 'axios'
import { computed, onBeforeUnmount, onMounted, ref, shallowRef } from 'vue'
import { getLoginSetting, setLoginSetting } from '~/api/index.js'
import { showInfo } from '~/utils/index.js'

let lastReq = null

const defaultValue = {
  enabled: false,
  user: '',
  password: '',
  max_age: 1,
}
const old = ref({ ...defaultValue })
const formData = ref({ ...defaultValue })
const showEmpty = shallowRef(false)
const loading = shallowRef(false)
const formRef = ref()
function getCfg() {
  loading.value = true
  lastReq = getLoginSetting()
  lastReq.rsp.then((data) => {
    data.enabled = !!data.user
    data.password = ''
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
function resetForm() {
  formData.value = { ...old.value }
}
function submitForm(formEl) {
  formEl.validate((valid) => {
    if (valid) {
      loading.value = true
      const value = {
        max_age: formData.value.max_age,
      }
      if (formData.value.enabled) {
        value.user = formData.value.user
        value.password = formData.value.password
      }
      else {
        value.user = ''
        value.password = ''
      }
      lastReq = setLoginSetting(value)
      lastReq.rsp.then(() => {
        formData.value.user = value.user
        formData.value.password = value.password
        formData.value.max_age = value.max_age
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
const can_update = computed(() => {
  if (formData.value.enabled) {
    return formData.value.user !== old.value.user
      || formData.value.password !== old.value.password
      || formData.value.max_age !== old.value.max_age
  }
  return formData.value.enabled !== old.value.enabled
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
      ref="formRef" v-model:model="formData" label-width="auto" :rules="{
        user: [{ required: true, message: '用户名不能为空' }],
        password: [{ required: true, message: '密码不能为空' }],
      }"
    >
      <el-form-item label="启用" prop="enabled">
        <el-checkbox v-model="formData.enabled" />
      </el-form-item>
      <el-form-item v-if="formData.enabled" label="用户名" prop="user">
        <el-input v-model="formData.user" />
      </el-form-item>
      <el-form-item v-if="formData.enabled" label="密码" prop="password">
        <el-input v-model="formData.password" type="password" show-password />
      </el-form-item>
      <el-form-item v-if="formData.enabled" label="Cookie有效期" prop="max_age">
        <el-input-number v-model="formData.max_age" :min="1">
          <template #suffix>
            小时
          </template>
        </el-input-number>
      </el-form-item>
    </el-form>
    <el-button type="primary" :disabled="!can_update" @click="submitForm(formRef)">
      应用
    </el-button>
    <el-button @click="resetForm">
      重置
    </el-button>
  </div>
</template>

<style scoped>

</style>
