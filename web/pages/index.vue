<template>
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
  <div>
    <h1>Upload file</h1>
    <div>
      <input type="file" @change="handleFileUpload">
      <button @click="uploadFile">Upload</button>
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
      file: null
    };
  },
  created() {
    const jwt = useCookie('jwt').value
    if (!jwt) {
      this.$router.push('/login');
      return
    }
    $fetch
      (useRuntimeConfig().public.apiUrl, {
        method: 'GET',
        headers: {
          "Authorization": jwt || ""
        }
      })
      .catch(err => {
        if (err.response.status === 401) {
          useCookie('jwt').value = ""
          this.$router.push('/login');
        }
      })
  },
  methods: {
    handleFileUpload(event) {
      this.file = event.target.files[0];
    },
    uploadFile() {
      if (this.file) {
        let formData = new FormData();
        console.log(this.file)
        formData.append('file', this.file);

        $fetch
          (useRuntimeConfig().public.apiUrl + "upload", {
            method: 'POST',
            body: formData,
            headers: {
              "Authorization": useCookie('jwt').value || ""
            }
          })
          .then(data => {
            console.log(data)
            console.log('File uploaded:', this.file);
          })
          .catch(err => {
            console.error('Error uploading file:', err);
            if (err.response.status === 401) {
              useCookie('jwt').value = ""
              this.$router.push('/login');
            }
          })
      }
    }
  }
}
</script>
