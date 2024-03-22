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
          <!-- <p>{{ document.uuid  }}</p> -->
          <p><a :href="document.url" target="_blank">{{ document.name }}</a></p>
          <embed :src="document.url" width="200px" height="300px" />
          <!-- <embed src="http://localhost:7777/doc?uuid=5df63112-4050-419c-9913-ff13a6e4b80e" width="300px" height="600px" /> -->
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
          console.log(data)
          this.documents = JSON.parse(data)
        })
        .catch(error => {
          console.error('Error pulling docs:', error);
        })

      // this.documents = [
      //   { id: "1", name: "resume1", url: "https://www.7onn.dev./resume.pdf" },
      //   { id: "2", name: "resume2", url: "https://www.7onn.dev./resume.pdf" },
      // ]
    }
  }
};
</script>

<style scoped>
/* Add your protected page styles here */
</style>
