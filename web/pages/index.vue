<template>
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
    const user = JSON.parse(useCookie('jwt').value || false)
    console.log(user)
    if (!user) {
      this.$router.push('/login');
    }
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
          })
          .then(data => {
            console.log(data)
            console.log('File uploaded:', this.file);
          })
          .catch(error => {
            console.error('Error uploading file:', error);
          })
      }
    }
  }
}
</script>
