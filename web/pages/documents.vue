<template>
  <div>
    <nav>
      <!-- <router-link to="/">Home</router-link> | -->
      <!-- <router-link to="/documents">Documents</router-link> -->
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
    const user = useCookie('jwt').value
    console.log(user)
    if (!user) {
      this.$router.push('/login');
    } else {
      this.fetchDocuments();
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
        .catch(error => {
          console.error('Error pulling docs:', error);
        })
    }
  }
};
</script>

<style scoped>
/* Add your protected page styles here */
</style>
