export default {
  install: (app, defaultOptions = {}) => {

    const authHeader = () => {
      const sessionToken = localStorage.getItem("session")
      if (sessionToken) {
        return {'Authorization': `Bearer ${sessionToken}`}
      }
      return {}
    }
    let api = {}

    let httpMethods = ['get', 'post', 'put', 'delete']
    httpMethods.forEach((method) => {
      api[method] = (url, data = null, requestOptions = {}) => {
        const options = {
          ...defaultOptions,
          method: method.toUpperCase(),
          headers: {
            'Content-Type': 'application/json',
            ...authHeader()
          },
          ...requestOptions,
        }
        if (data) {
          options.body = JSON.stringify(data)
        }

        return fetch(url, options)
          .then((response) => {
            if (!response.ok) {
              throw response
            }
            return response.json()
          })
      }
    })

    app.config.globalProperties.$api = api
  }
}