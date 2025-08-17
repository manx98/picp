<script setup>
import { ref } from 'vue'
import { EpKey, EpUserFilled } from 'vue-icons-plus/ep'
import { login } from '~/api'
import { showInfo } from '~/utils/index.js'

const loginData = ref({
  user: '',
  password: '',
})

const formRef = ref()
const loading = ref(false)
function submitForm() {
  formRef.value.validate((valid) => {
    if (valid) {
      loading.value = true
      const req = login(loginData.value)
      req.rsp.then(() => {
        showInfo(false, '登陆成功，即将跳转...')
        window.location.href = '/'
      }).catch((e) => {
        showInfo(true, e.message)
      }).finally(() => {
        loading.value = false
      })
    }
  })
}
</script>

<template>
  <div style="display: flex; justify-self: center;align-items: center;height: 100vh">
    <el-card v-loading="loading" style="width: 300px;">
      <template #header>
        <div style="text-align: center">
          <h3>登陆</h3>
        </div>
      </template>
      <el-form
        ref="formRef"
        v-model:model="loginData" :rules="{
          user: [{ required: true, message: '用户名不能为空', trigger: 'blur' }],
          password: [{ required: true, message: '密码不能为空', trigger: 'blur' }],
        }"
      >
        <el-form-item label-width="auto" prop="user">
          <el-input v-model="loginData.user" placeholder="用户名">
            <template #prefix>
              <EpUserFilled />
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label-width="auto" prop="password">
          <el-input v-model="loginData.password" @keydown.enter="submitForm(formRef)" type="password" placeholder="密码" show-password>
            <template #prefix>
              <EpKey />
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <div style="text-align: center">
        <el-button type="primary" @click="submitForm">
          登陆
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<style scoped>

</style>
