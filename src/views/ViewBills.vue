<template>
  <div>
    <h1>View Bills</h1>
    <ul>
      <li v-for="bill in bills" :key="bill.id">{{ bill.name }} - ${{ bill.amount }}</li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      bills: []  // Empty array to hold the bills from the database
    };
  },
  mounted() {
    this.fetchBills();  // Fetch bills when the component is mounted
  },
  methods: {
    async fetchBills() {
      try {
        const response = await axios.get('http://localhost:8080/bills');
        this.bills = response.data;  // Set the bills data from the response
      } catch (error) {
        console.error("Error fetching bills:", error);
        alert("There was an error fetching the bills.");
      }
    }
  }
};
</script>
