<script setup>
 import VuePdfEmbed from 'vue-pdf-embed'

 import 'vue-pdf-embed/dist/style/index.css'

 import 'vue-pdf-embed/dist/style/annotationLayer.css'
 import 'vue-pdf-embed/dist/style/textLayer.css'
</script>

<template>
    <div>
        <nav>
            <!-- <router-link to="/">Home</router-link> | -->
            <router-link to="/documents">Documents</router-link>
        </nav>
        <router-view/>
        <div>
            <h2>Documents</h2>
            <ul style="display: flex; flex-wrap: wrap;">
                <li v-for="document in documents" :key="document.id" style="margin-right: 10px;">
                    <a :href="document.url" target="_blank">{{ document.name }}</a>
                        <VuePdfEmbed annotation-layer text-layer :source="document.url" />
                </li>
            </ul>
        </div>
    </div>
</template>
<script>
 export default {
     data() {
         return {
             documents: []
         };
     },
     created() {
         const user = JSON.parse(localStorage.getItem('user'));
         if (!user || !user.isAuthenticated) {
             this.$router.push('/login');
         } else {
             this.fetchDocuments();
         }
     },
     methods: {
         fetchDocuments() {
             this.documents = [
                 {id: "1", name: "resume1", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "2", name: "resume2", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "3", name: "resume3", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "4", name: "resume4", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "5", name: "resume5", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "1", name: "resume1", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "2", name: "resume2", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "3", name: "resume3", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "4", name: "resume4", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "5", name: "resume5", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "1", name: "resume1", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "2", name: "resume2", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "3", name: "resume3", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "4", name: "resume4", url: "https://www.7onn.dev./resume.pdf"},
                 {id: "5", name: "resume5", url: "https://www.7onn.dev./resume.pdf"},
             ]
             /* fetch('/docs')
              *     .then(response => response.json())
              *     .then(data => {
              *         this.documents = data;
              *     })
              *     .catch(error => {
              *         console.error('Error fetching documents:', error);
              *     }); */
         }
     }
 };
</script>

<style scoped>
 /* Add your protected page styles here */
</style>

<!-- 

     <template>
     <nav>
     <router-link to="/">Home</router-link> |
     <router-link to="/documents">Documents</router-link>
     </nav>
     <router-view/>
     <div>
     <h2>Documents</h2>
     <p>List of documents</p>
     </div>
     </template>

     <script>
     export default {
     created() {
     // todo: use jwt with some backend
     const user = JSON.parse(localStorage.getItem('user'));
     if (!user || !user.isAuthenticated) {
     this.$router.push('/login');
     }
     }
     };
     </script>

     <style scoped>
     /* Add your protected page styles here */
     </style> -->
