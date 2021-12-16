export function parseQuery() {
  let query = window.location.search.substring(1)
  let result = new Map<string, string>()
  query.split("&").forEach(item => {
    let [key, value] = item.split("=")
    result.set(key, value)
  })
  return result
}

export function queryGet(m: Map<string, string>, key: string) {
  let ret = m.get(key)
  if (!ret) {
    return ""
  }
  return ret
}

function asString(m: Map<string, string>) {
  let ret = ""
  m.forEach((value, key) => {
    ret += `${key}=${value}&`
  })
  return ret
}

export function setQuery(query: Map<string, string>, k: string, v: string) {
  query.set(k, v)
  window.history.pushState(null, "", "?" + asString(query))
}

export function mustSplit3(joint: string) {
  let splited = joint.split("/")
  if (splited.length < 3) {
    return ["", "", ""]
  }
  return splited
}