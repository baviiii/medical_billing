<template>
  <div id="app">
    <NavBarComponent /> <!-- Use NavBarComponent instead of Navbar -->
    <header>
      <h1>Medical Billing Portal</h1>
    </header>
    <nav>
      <ul>
        <li><router-link to="/">Home</router-link></li>
        <li><router-link to="/view-bills">View Bills</router-link></li>
        <li><router-link to="/add-bill">Add Bill</router-link></li>
        <li><router-link to="/manage-bills">Manage Bills</router-link></li>
      </ul>
    </nav>
    <router-view /> <!-- This will render the current route's component -->
  </div>
</template>

<script>
import NavBarComponent from './components/NavBarComponent.vue';
export default {
  name: 'App',
  components: {
    NavBarComponent // Keep this as NavBarComponent
  }
};
</script>
