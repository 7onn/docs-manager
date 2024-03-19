<template>
  <div>
    <nav>
      <!-- <router-link to="/">Home</router-link> | -->
      <router-link to="/documents">Documents</router-link>
    </nav>
    <router-view />
    <div>
      <h2>Documents</h2>
      <ul style="display: flex; flex-wrap: wrap;">
        <li v-for="document in documents" :key="document.id" style="margin-right: 10px;">
          <a :href="document.url" target="_blank">{{ document.name }}</a>
          <embed :src="document.url" width="300px" height="600px" />
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
    const user = JSON.parse(useCookie('jwt').value || false)
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
          headers: { "Content-Type": "application/json" }
        })
        .then(data => {
          console.log(data)
          this.documents = JSON.stringify(data)
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
