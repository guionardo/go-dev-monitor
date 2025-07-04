export default function useAPI() {
  function getFileStatusClass(status) {
    const statusMap = {
      '?': 'untracked',
      A: 'added',
      M: 'modified',
      D: 'deleted',
      R: 'renamed',
      C: 'copied',
      ' ': 'unmodified',
    }
    return statusMap[status] || 'unknown'
  }

  function parseFile(file) {
    let { name, time, status } = file
    return {
      name: name,
      time: new Date(time).toLocaleString(),
      status: getFileStatusClass(status),
    }
  }

  function parseHost(host) {
    host.untracked_files = (host.untracked_files || []).map((f) => parseFile(f))
    host.changed_files = (host.changed_files || []).map((f) => parseFile(f))
    host.last_changed_files = (host.last_changed_files || []).map((f) => parseFile(f))
    return host
  }
  async function getData() {
    const data = []
    const apiData = await fetchData()
    let dirtyOrigins = 0

    Object.keys(apiData.origins).forEach((origin) => {
      const newOrigin = {
        origin: origin,
        hosts: apiData.origins[origin].map((h) => parseHost(h)),
      }
      newOrigin.dirtyHosts = newOrigin.hosts.reduce((p, c) => (p + c.clean ? 0 : 1), 0)
      newOrigin.countHosts = newOrigin.hosts.length
      newOrigin.clean = newOrigin.dirtyHosts == 0
      newOrigin.description = newOrigin.hosts[0].description
      newOrigin.language = newOrigin.hosts[0].language
      newOrigin.language_icon = newOrigin.hosts[0].language_icon

      data.push(newOrigin)
      if (!newOrigin.clean) {
        dirtyOrigins += 1
      }
    })
    return { origins: data, dirtyOrigins: dirtyOrigins, countOrigins: data.length }
  }

  const apiInfo = { url: null }

  async function tryFetch(url) {
    try {
      const resp = await fetch(url)
      console.debug(`tryFetch: ${url}`, resp)
      return resp.ok
    } catch (e) {
      console.error(`tryFetch: ${url}`, e)
    }
    return false
  }

  async function getApiURL() {
    if (apiInfo.url) {
      return apiInfo.url
    }
    const reqDebug = tryFetch('http://localhost:3800/hc');
    const reqProd = tryFetch('/hc');

    try {
      const response = await Promise.race([reqDebug, reqProd]);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      if (response.url.includes('localhost')) {
        apiInfo.url = 'http://localhost:3800'
      } else {
        apiInfo.url = '/'
      }
      console.debug('Response', response)

      return apiInfo.url;
    } catch (error) {
      console.error("Error with one of the requests:", error);
      throw error; // Re-throw the error to be handled by the caller
    }
  }

  async function fetchData() {
    try {
      const baseUrl = await getApiURL()
      const response = await fetch(`${baseUrl}data`)
      if (!response.ok) {
        throw new Error(`Response status: ${response.statusText}`)
      }
      const json = await response.json()
      console.log('Data', json)
      return json
    } catch (e) {
      console.error('Fetch API data error', e)
    }
    return {}
  }

  return { getData, fetchData }
}

