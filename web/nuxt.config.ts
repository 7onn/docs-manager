// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  runtimeConfig: {
      apiSecret : '123',
      public: {
        apiUrl: 'http://localhost:7777/'
      }
    },
    routeRules: {
      '/**': {
          cors: true
       },
  },
})
