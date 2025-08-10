<script setup>
import { computed, defineProps, ref, shallowRef, watch } from 'vue'
import ConfirmDialog from './ConfirmDialog.vue'

const props = defineProps({
  info: Object,
  show: Boolean,
})

const emits = defineEmits(['ok', 'cancel', 'update:show'])
const config = ref({
  device: '',
  ssid: '',
  auto_connect: false,
  password: '',
  bssid: '',
})
const internalShow = computed({
  set(v) {
    emits('update:show', v)
  },
  get() {
    return props.show
  },
})

const passwordError = shallowRef()
watch(() => props.show, (newVal) => {
  if (newVal) {
    config.value = props.info
    config.value.password = ''
    passwordError.value = ''
  }
}, {
  immediate: true,
})
const formRef = ref()
function onOk(cb) {
  passwordError.value = ''
  const done = (msg) => {
    if (msg) {
      if (msg === 'passwords required') {
        passwordError.value = '密码错误'
      }
      else {
        passwordError.value = msg
      }
    }
    cb()
  }
  formRef.value.validate((valid) => {
    if (valid) {
      emits('ok', {
        config: config.value,
        done,
      })
    }
    else {
      done()
    }
  })
}
</script>

<template>
  <ConfirmDialog v-model:show="internalShow" confirm-text="连接" title="连接" @ok="onOk">
    <template #default>
      <el-form
        ref="formRef"
        :model="config" :rules="{
          password: [
            { required: true, message: '请输入密码', trigger: 'blur' },
            { min: 8, message: '密码最少为8位', trigger: 'blur' },
          ],
        }"
      >
        <el-form-item label="设备">
          <el-input disabled :model-value="config.device" />
        </el-form-item>
        <el-form-item label="SSID">
          <el-input :model-value="config.ssid" disabled />
        </el-form-item>
        <el-form-item label="BSSID">
          <el-input :model-value="config.bssid" disabled />
        </el-form-item>
        <el-form-item v-if="config.need_password" :error="passwordError" label="密码" prop="password">
          <el-input
            v-model:model-value="config.password" type="password" show-password placeholder="密码"
            maxlength="32"
          />
        </el-form-item>
        <el-form-item label="自动连接">
          <el-switch v-model:model-value="config.auto_connect" />
        </el-form-item>
      </el-form>
    </template>
  </ConfirmDialog>
</template>

<style scoped>

</style>
