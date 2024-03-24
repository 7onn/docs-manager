<template>
  <div>
    <nav>
      <ul>
        <li>
          <NuxtLink to="/">Home</NuxtLink>
        </li>
        <li>
          <NuxtLink to="/documents">Documents</NuxtLink>
        </li>
      </ul>
    </nav>
    <router-view />
    <div>
      <h2>Documents</h2>
      <ul style="display: flex; flex-wrap: wrap;">
        <li v-for="document, index in documents" :key="index" style="margin-right: 10px;">
          <p><a :href="document.url" target="_blank">{{ document.name }}</a></p>
          <embed :src="document.url" width="200px" height="300px" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
definePageMeta({
  layout: 'web'
})

export default {
  data() {
    return {
      documents: []
    };
  },
  created() {
    const jwt = useCookie('jwt').value
    if (!jwt) {
      this.$router.push('/login');
      return
    } else {
      $fetch
        (useRuntimeConfig().public.apiUrl, {
          method: 'GET',
          headers: {
            "Authorization": jwt || ""
          }
        })
        .then(() => {
          this.fetchDocuments();
        })
        .catch(err => {
          if (err.response.status === 401) {
            useCookie('jwt').value = ""
            this.$router.push('/login');
          }
        })
    }
  },
  methods: {
    fetchDocuments() {
      $fetch
        (useRuntimeConfig().public.apiUrl + "docs", {
          method: 'GET',
          headers: {
            "Content-Type": "application/json",
            "Authorization": useCookie('jwt').value || ""
          },
        })
        .then(data => {
          if (data) {
            for (const doc of data) {
              doc.url = useRuntimeConfig().public.apiUrl + "doc?uuid=" + doc.uuid
              this.documents.push(doc)
            }
          }
        })
        .catch(err => {
          console.error('Error pulling docs:', err);
          if (err.response.status === 401) {
            useCookie('jwt').value = ""
            this.$router.push('/login');
          }
        })
    }
  }
};
</script>

<style scoped>
/* Add your protected page styles here */
</style>
