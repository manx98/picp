<script>
import {windows_size} from '~/utils'
export default {
  name: 'ConfirmDialog',
  props: {
    show: Boolean,
    title: String,
    disableConfirm: {
      type: Boolean,
      default: false,
    },
    confirmText: {
      type: String,
      default: '确认',
    },
    width: {
      type: Number,
      default: 400,
    },
    cancelText: {
      type: String,
      default: '取消',
    },
    cancelMethod: {
      type: Function,
      default: null,
    },
  },
  emits: ['ok', 'update:show'],
  data() {
    return {
      value: '',
      loading: false,
    }
  },
  computed: {
    windows_size() {
      return windows_size
    },
    internalShow: {
      get() {
        return this.show
      },
      set(v) {
        this.$emit('update:show', v)
      },
    },
  },
  watch: {
    show(v) {
      if (v) {
        this.value = this.defaultValue
      }
    },
  },
  methods: {
    confirm() {
      this.loading = true
      this.$emit('ok', () => this.loading = false)
    },
    cancel() {
      if (this.cancelMethod) {
        this.loading = true
        this.cancelMethod(() => this.loading = false)
      }
      else {
        this.internalShow = false
      }
    },
  },
}
</script>

<template>
  <el-dialog
      :style="{maxWidth: windows_size.value.width + 'px'}"
    v-model="internalShow" :width="width" :modal="false" :close-on-click-modal="!loading" append-to-body draggable
    destroy-on-close :show-close="!loading"
  >
    <template #title>
      <slot name="title">
        {{ title }}
      </slot>
    </template>
    <template #default>
      <div v-loading="loading">
        <slot />
      </div>
    </template>
    <template #footer>
      <div v-show="!loading" style="text-align: right">
        <el-button :disabled="disableConfirm" type="primary" @click="confirm">
          {{ confirmText }}
        </el-button>
        <el-button :disabled="loading" @click="cancel">
          {{ cancelText }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>

</style>
