<template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="login">
      <div>
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" required />
      </div>
      <div>
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
    <p v-if="error">{{ error }}</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'UserLogin',
  data() {
    return {
      username: '',
      password: '',
      error: null,
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('http://localhost:8080/login', {
          username: this.username,
          password: this.password,
        });
        const token = response.data.token;
        localStorage.setItem('token', token); // Save token to local storage
        this.$router.push('/'); // Redirect to the home page after login
      } catch (err) {
        this.error = 'Invalid username or password';
      }
    },
  },
};
</script>

<style scoped>
/* Add basic styles here */
</style>
