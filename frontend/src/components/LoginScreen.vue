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
                    <input type="password" v-model="signupConfirmPassword" placeholder="Confirm password" @input="validatePassword">
                    <p v-if="signupPasswordMismatch" style="color: red;">Passwords do not match</p>
                    <button type="submit">Sign up</button>
                </form>
            </div>
        </div>
    </template>


    <!-- <template>
         <div>
         <h2>login</h2>
         <form @submit.prevent="login">
         <input type="text" v-model="loginemail" placeholder="email">
         <input type="password" v-model="loginpassword" placeholder="password">
         <button type="submit">login</button>
         </form>

         <h2>sign up</h2>
         <form @submit.prevent="login">
         <input type="text" v-model="signupemail" placeholder="email">
         <input type="password" v-model="signuppassword" placeholder="password">
         <input type="password" v-model="signupconfirmpassword" placeholder="confirm password" @input="validatepassword">
         <p v-if="signuppasswordmismatch" style="color: red;">passwords do not match</p>
         <button type="submit">sign up</button>
         </form>
         </div>
         </template> -->

    <script>
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
             validatePassword() {
                 this.signupPasswordMismatch = this.signupPassword !== this.signupConfirmPassword;
             },
             login() {
                 if (this.loginEmail && this.loginPassword) {
                     localStorage.setItem('user', JSON.stringify({ isAuthenticated: true }));
                     this.$router.push('/documents');
                 } else {
                     alert('Invalid username or password');
                 }
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
