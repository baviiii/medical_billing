<template>
  <div>
    <h1>Add a New Bill</h1>
    <form @submit.prevent="submitBill">
      <label for="name">Bill Name:</label>
      <input type="text" v-model="bill.name" required />

      <label for="amount">Amount:</label>
      <input type="number" v-model="bill.amount" required />

      <button type="submit">Add Bill</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      bill: {
        name: '',
        amount: ''
      }
    };
  },
  methods: {
    async submitBill() {
      // Validate inputs
      if (!this.bill.name || !this.bill.amount) {
        alert("Please fill out all fields.");
        return;
      }

      const token = localStorage.getItem('token');
      try {
        const response = await axios.post('http://localhost:8080/bills', this.bill, {
          headers: {
            Authorization: `Bearer ${token}`, // Include the token in the request
          },
        });
        console.log("Bill Submitted:", response.data);
        
        // Reset the form after submission
        this.bill.name = '';
        this.bill.amount = '';

        // Provide feedback to the user
        alert("Bill added successfully!");
      } catch (error) {
        console.error("Error submitting bill:", error.response ? error.response.data : error.message);
        alert("There was an error adding the bill.");
      }
    }
  }
};
</script>

<style scoped>
/* Add basic styles here */
</style>
