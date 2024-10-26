<template>
  <nav class="navbar">
    <ul>
      <li><router-link to="/">Home</router-link></li>
      <!-- Conditionally render Login/Register if not authenticated -->
      <li v-if="!isAuthenticated">
        <router-link to="/login">Login</router-link>
      </li>
      <li v-if="!isAuthenticated">
        <router-link to="/register">Register</router-link>
      </li>
      <!-- Show Logout button if authenticated -->
      <li v-if="isAuthenticated">
        <button @click="logout" class="logout-button">Logout</button>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'NavBarComponent',
  computed: {
    // Check if user is authenticated
    isAuthenticated() {
      const token = localStorage.getItem('token');
      console.log('Is Authenticated:', !!token); // Debugging line
      return !!token;
    },
  },
  methods: {
    // Logout method to clear token and redirect to login page
    logout() {
      localStorage.removeItem('token');
      this.$router.push('/login');
    },
  },
};
</script>

<style scoped>
.navbar {
  background-color: #333;
  color: white;
  padding: 10px;
}

.navbar ul {
  list-style: none;
  display: flex;
  gap: 15px;
}

.navbar a,
.logout-button {
  color: white;
  text-decoration: none;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
}

.navbar a:hover,
.logout-button:hover {
  color: #ddd;
}

.navbar a:active,
.logout-button:active {
  color: #bbb;
}
</style>
