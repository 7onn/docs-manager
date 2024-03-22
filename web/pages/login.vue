<template>
  <div>
    <nav>
      <ul>
        <li :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">Login</li>
        <li :class="{ active: activeTab === 'signup' }" @click="activeTab = 'signup'">Sign up</li>
      </ul>
    </nav>

    <div v-if="activeTab === 'login'">
      <h2>Login</h2>
      <form @submit.prevent="login">
        <input type="text" v-model="loginEmail" placeholder="Email">
        <input type="password" v-model="loginPassword" placeholder="Password">
        <button type="submit">Login</button>
      </form>
    </div>

    <div v-if="activeTab === 'signup'">
      <h2>Sign up</h2>
      <form @submit.prevent="signup">
        <input type="text" v-model="signupEmail" placeholder="Email">
        <input type="password" v-model="signupPassword" placeholder="Password">
        <input type="password" v-model="signupConfirmPassword" placeholder="Confirm password"
          @input="validatePasswordMatch">
        <p v-if="signupPasswordMismatch" style="color: red;">Passwords do not match</p>
        <button type="submit">Sign up</button>
      </form>
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
      activeTab: 'login',
      loginEmail: '',
      loginPassword: '',
      signupEmail: '',
      signupPassword: '',
      signupConfirmPassword: '',
      signupPasswordMismatch: false
    };
  },
  methods: {
    validatePasswordMatch() {
      this.signupPasswordMismatch = this.signupPassword !== this.signupConfirmPassword;
    },
    login() {
      $fetch
        (useRuntimeConfig().public.apiUrl + "login", {
          method: 'POST',
          body: JSON.stringify({
            "email": this.loginEmail,
            "password": this.loginPassword
          }),
          headers: { "Content-Type": "application/json" }
        })
        .then(data => {
          useCookie('jwt').value = JSON.stringify(data)
          this.$router.push('/documents');
        })
        .catch(error => {
          console.error('Error signin up:', error);
        })
    },
    signup() {
      $fetch
        (useRuntimeConfig().public.apiUrl + "signup", {
          method: 'POST',
          body: JSON.stringify({
            "email": this.signupEmail,
            "password": this.signupPassword
          }),
          headers: { "Content-Type": "application/json" }
        })
        .then(data => {
          useCookie('jwt').value = JSON.stringify(data)
          this.$router.push('/documents');
        })
        .catch(error => {
          console.error('Error signin up:', error);
        })
    }
  },
};
</script>


<style scoped>
nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

nav ul li {
  display: inline-block;
  padding: 10px;
  cursor: pointer;
}

nav ul li.active {
  background-color: #ccc;
}
</style>
