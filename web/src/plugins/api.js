export default {
  install: (app, apiOptions = {}) => {

    let { defaultOptions } = apiOptions

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
              if (apiOptions.forbiddenRoute && response.status === 401) {
                app.config.globalProperties.$router.push(apiOptions.forbiddenRoute)
              } else {
                throw response
              }
            }
            return response.json()
          })
      }
    })

    app.config.globalProperties.$api = api
  }
}