<template>
  <div>
    <h1>Welcome to the homepage</h1>
  </div>
</template>

<script>
definePageMeta({
  layout: 'web'
})

export default {
  data() {
    return {
      a: []
    };
  },
  created() {
    const user = JSON.parse(useCookie('jwt').value || false)
    console.log(user)
    if (!user) {
      this.$router.push('/login');
    }
  },
}
</script>
