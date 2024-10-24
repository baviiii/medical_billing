<template>
    <div>
      <h2>{{ title }}</h2>
      <p>{{ description }}</p>
      <ul>
        <li v-for="bill in bills" :key="bill.id">
          {{ bill.name }} - ${{ bill.amount }}
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  export default {
    name: 'BillingDashboard',
    data() {
      return {
        title: 'Billing Dashboard',
        description: 'Manage your medical billing records here.',
        bills: [
          { id: 1, name: 'Hospital Invoice', amount: 500 },
          { id: 2, name: 'Lab Test', amount: 120 },
          { id: 3, name: 'Consultation Fee', amount: 200 }
        ]
      };
    }
  }
  </script>
  