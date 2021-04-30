export default {
  install: (app, options = {}) => {

    const authHeader = () => {
      const sessionToken = localStorage.getItem("session")
      if (sessionToken) {
        return {'Authorization': `Bearer ${sessionToken}`}
      }
      return {}
    }

    const {
      headers = {},
    } = options

    app.config.globalProperties.$api = {

      get: (url, reqOptions = {}) => {

        const options = {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            ...authHeader()
          },
          ...reqOptions,
        }

        return fetch(url, options).then((response) => {
          if (!response.ok) {
            throw response
          }
          return response.json()
        })
      },
      post: (url, data, reqOptions = {}) => {

        const options = {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...authHeader()
          },
          ...reqOptions,
          body: JSON.stringify(data)
        }

        return fetch(url, options).then((response) => {
          if (!response.ok) {
            throw response
          }
          return response.json()
        })
      }
    }
  }
}