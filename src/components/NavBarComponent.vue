<template>
  <nav class="navbar">
    <ul>
      <li><router-link to="/">Home</router-link></li>
      <li><router-link to="/login">Login</router-link></li>
      <li><router-link to="/register">Register</router-link></li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'NavBarComponent'
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

.navbar a {
  color: white;
  text-decoration: none;
}
</style>
