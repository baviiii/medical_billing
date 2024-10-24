<template>
  <div id="app">
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
export default {
  name: 'App',
}
</script>
