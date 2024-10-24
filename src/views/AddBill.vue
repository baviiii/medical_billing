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
      
      try {
        const response = await axios.post('http://localhost:8080/bills', this.bill);
        console.log("Bill Submitted:", response.data);
        
        // Optionally reset the form after submission
        this.bill.name = '';
        this.bill.amount = '';

        // Optionally, provide feedback to the user
        alert("Bill added successfully!");
      } catch (error) {
        console.error("Error submitting bill:", error);
        alert("There was an error adding the bill.");
      }
    }
  }
};
</script>
