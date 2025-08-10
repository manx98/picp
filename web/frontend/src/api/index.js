import axios from 'axios'

function handlerResponse(response) {
  if (response.status === 200) {
    if (response.data.code === 200) {
      return response.data.data
    }
    else if (response.status === 401) {
      location.reload()
      return
    }
    throw new Error(response.data.msg)
  }
  else {
    throw new Error(`server error ${response.status}`)
  }
}

export function getDevices() {
  return axios.get('/api/devices').then(handlerResponse)
}

export function listWifiAp(device) {
  return axios.get(`/api/wifi`, {
    params: {
      device,
    },
  }).then(handlerResponse)
}

export function createWifiAp(device) {
  return axios.post('/api/wifi', device).then(handlerResponse)
}

export function deleteWifiAp(connection_uuid) {
  return axios.delete(`/api/wifi`, {
    params: {
      connection_uuid,
    },
  }).then(handlerResponse)
}
