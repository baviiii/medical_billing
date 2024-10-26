<template>
  <div>
    <h1>Register</h1>
    <form @submit.prevent="register">
      <div>
        <label for="username">Username:</label>
        <input type="text" v-model="username" required />
      </div>
      <div>
        <label for="phone">Phone Number:</label>
        <input type="tel" v-model="phone" required />
      </div>
      <button type="submit">Send OTP</button>
    </form>

    <div v-if="otpSent">
      <h2>Verify OTP</h2>
      <form @submit.prevent="verifyOtp">
        <div>
          <label for="otp">OTP:</label>
          <input type="text" v-model="otp" required />
        </div>
        <button type="submit">Verify OTP</button>
      </form>
    </div>

    <p v-if="error">{{ error }}</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      username: '',
      phone: '',
      otp: '',
      otpSent: false,
      error: null,
    };
  },
  methods: {
    async register() {
      try {
        await axios.post('http://localhost:8080/register', {
          username: this.username,
          phone: this.phone,
        });
        this.otpSent = true;
      } catch (error) {
        this.error = 'Error sending OTP';
      }
    },
    async verifyOtp() {
      try {
        await axios.post('http://localhost:8080/verify-otp', {
          phone: this.phone,
          otp: this.otp,
        });
        // Redirect or show success message after successful verification
      } catch (error) {
        this.error = 'Invalid OTP';
      }
    },
  },
};
</script>

<style scoped>
/* Add your styles here */
</style>
