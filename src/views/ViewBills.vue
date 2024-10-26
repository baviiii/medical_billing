<template>
  <div>
    <h1>Bills</h1>
    <ul v-if="bills.length">
      <li v-for="bill in bills" :key="bill.id">
        {{ bill.name }} - ${{ bill.amount }}
      </li>
    </ul>
    <p v-else-if="error">{{ error }}</p>
    <p v-else>No bills found.</p>
  </div>
</template>

<script>
import axios from 'axios';
const token = localStorage.getItem('token');
console.log("Token from localStorage:", token);  // Check if token is retrieved correctly

export default {
  name: 'ViewBills',
  data() {
    return {
      bills: [],
      error: null,
    };
  },
  async created() {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get('http://localhost:8080/bills', {
        headers: {
          Authorization: `Bearer ${token}`, // Token passed in the header with Bearer scheme
        },
      });
      this.bills = response.data;
    } catch (error) {
      console.error("Error fetching bills:", error.response ? error.response.data : error.message);
      this.error = 'Error fetching bills';
    }
  },
};
</script>

<style scoped>
ul {
  list-style: none;
  padding: 0;
}
li {
  background-color: #f4f4f4;
  padding: 10px;
  margin-bottom: 5px;
}
</style>
