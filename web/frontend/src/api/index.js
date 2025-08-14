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

function sendGet(url, params) {
  const token = axios.CancelToken.source()
  if (params) {
    params.cancelToken = token.token
  }
  else {
    params = {
        cancelToken: token.token
    }
  }
  return {
      rsp: axios.get(url, params).then(handlerResponse),
      cancel: token.cancel
  }
}

export function getDevices() {
  return sendGet('/api/devices')
}

export function listWifiAp(device) {
  return sendGet(`/api/wifi`, {
    params: {
      device,
    },
  })
}

function sendPost(url, data) {
  const token = axios.CancelToken.source()
  return {
      rsp: axios.post(url, data, {
          cancelToken: token.token,
      }).then(handlerResponse),
      cancel: token.cancel
  }
}

export function createWifiAp(device) {
  return sendPost('/api/wifi', device)
}

function sendDelete(url, params) {
  const token = axios.CancelToken.source()
  if (params) {
    params.cancelToken = token.token
  }
  else {
    params = {
      cacheToken: token.token,
    }
  }
  return {
      rsp: axios.delete(url, params).then(handlerResponse),
      cancel: token.cancel
  }
}

export function deleteWifiAp(connection_uuid) {
  return sendDelete(`/api/wifi`, {
    params: {
      connection_uuid,
    },
  })
}
