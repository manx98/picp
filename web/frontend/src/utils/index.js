import { ElMessage } from 'element-plus'
import { ref } from 'vue'

const windows_size = ref({
  width: window.innerWidth,
  height: window.innerHeight,
})

window.addEventListener('resize', () => {
  windows_size.value = {
    width: window.innerWidth,
    height: window.innerHeight,
  }
})

function showInfo(is_error, msg) {
  if (is_error) {
    ElMessage({
      type: 'error',
      message: msg,
      duration: 1500,
      showClose: true,
      offset: 50,
    })
  }
  else {
    ElMessage({
      type: 'success',
      message: msg,
      duration: 1500,
      showClose: true,
      offset: 50,
    })
  }
}

export function showWarning(msg) {
  ElMessage({
    type: 'warning',
    message: msg,
    duration: 1500,
    showClose: true,
    offset: 50,
  })
}

export { showInfo, windows_size }
