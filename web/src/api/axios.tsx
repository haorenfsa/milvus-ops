import axios, { AxiosResponse } from "axios";
import { notification, message } from "antd";

export async function axiosGet(
  url: string,
): Promise<any|false> {
  let msg = "sending get request"
  let config = {
    validateStatus: validateStatus
  }
  let ret: void | AxiosResponse
  ret = await axios.get(`${url}`, config)
    .then((res: AxiosResponse) => {
      if (res.status >= 400) {
        message.error(msg + " failed: " + res.data)
        return false
      }
      if (res.status !== 200) {
        message.error(msg + ` unexpected code :${res.status} data:` + res.data)
        return false
      }
      return res.data
    })
  return ret
}

export async function axiosPost(
  url: string,
  data: any,
  returnData: boolean = false,
  msg?: string,
  // config?: object,
): Promise<any | false> {
  let loadDone: any
  if (msg) {
    msg = "sending post request: " + msg
    loadDone = message.loading(msg)
  }
  console.debug(data)
  
  let config = {
    validateStatus: validateStatus
  }
  const ret = await axios.post(`${url}`, data, config)
    .then((res: AxiosResponse) => {
      if (loadDone) {
        loadDone()
      }
      if (res.status >= 400) {
        message.error(msg + " failed: " + res.data)
        return false
      }
      if (res.status !== 200) {
        message.error(msg + ` unexpected code :${res.status} data:` + res.data)
        return false
      }
      if (loadDone) {
        message.success(msg + " success")
      }
      return returnData ? res.data : true
    })
    .catch((reason: string) => {
      notification.error({
        message: msg + " failed",
        description:  reason,
      })
      return false
    })
  return ret;
}


export async function axiosPut(
  url: string,
  data: any,
  returnData: boolean = false,
  msg?: string,
  // config?: object,
): Promise<any | false> {
  let loadDone: any
  if (msg) {
    msg = "sending post request: " + msg
    loadDone = message.loading(msg)
  }
  console.debug(data)
  
  let config = {
    validateStatus: validateStatus
  }
  const ret = await axios.put(`${url}`, data, config)
    .then((res: AxiosResponse) => {
      if (loadDone) {
        loadDone()
      }
      if (res.status >= 400) {
        message.error(msg + " failed: " + res.data)
        return false
      }
      if (res.status !== 200) {
        message.error(msg + ` unexpected code :${res.status} data:` + res.data)
        return false
      }
      if (loadDone) {
        message.success(msg + " success")
      }
      return returnData ? res.data : true
    })
    .catch((reason: string) => {
      notification.error({
        message: msg + " failed",
        description:  reason,
      })
      return false
    })
  return ret;
}

export async function axiosPatch(
  url: string,
  data: any,
  msg?: string,
  returnData: boolean = false,
  // config?: object,
): Promise<any | false> {
  let loadDone: any
  if (msg) {
    msg = "sending patch request: " + msg
    loadDone = message.loading(msg)
  }
  console.debug(data)
  
  let config = {
    validateStatus: validateStatus
  }
  const ret = await axios.patch(`${url}`, data, config)
    .then((res: AxiosResponse) => {
      if (loadDone) {
        loadDone()
      }
      if (res.status >= 400) {
        message.error(msg + " failed: " + res.data)
        return false
      }
      if (res.status !== 200) {
        message.error(msg + ` unexpected code :${res.status} data:` + res.data)
        return false
      }
      if (loadDone) {
        message.success(msg + " success")
      }
      return returnData ? res.data : true
    })
    .catch((reason: string) => {
      notification.error({
        message: msg + " failed",
        description:  reason,
      })
      return false
    })
  return ret;
}


export async function axiosDelete(
  url: string,
  msg?: string,
  returnData: boolean = false,
  // config?: object,
): Promise<any | false> {
  let loadDone: any
  if (msg) {
    msg = "sending delete request: " + msg
    loadDone = message.loading(msg)
  }
  
  let config = {
    validateStatus: validateStatus
  }
  const ret = await axios.delete(`${url}`, config)
    .then((res: AxiosResponse) => {
      if (loadDone) {
        loadDone()
      }
      if (res.status >= 400) {
        message.error(msg + " failed: " + res.data)
        return false
      }
      if (res.status !== 200) {
        message.error(msg + ` unexpected code :${res.status} data:` + res.data)
        return false
      }
      if (loadDone) {
        message.success(msg + " success")
      }
      return returnData ? res.data : true
    })
    .catch((reason: string) => {
      notification.error({
        message: msg + " failed",
        description:  reason,
      })
      return false
    })
  return ret;
}

function validateStatus () {
  return true; // Reject only if the status code is greater than or equal to 500
}