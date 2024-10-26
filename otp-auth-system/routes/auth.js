const express = require('express');
const bcrypt = require('bcryptjs');
const crypto = require('crypto');
const User = require('../models/User');
const OTP = require('../models/OTP');
const router = express.Router();
const nodemailer = require('nodemailer');

// Email transporter setup (use your email config here)
const transporter = nodemailer.createTransport({
  service: 'Gmail',
  auth: {
    user: 'your-email@gmail.com',
    pass: 'your-email-password'
  }
});

// Generate OTP
const generateOTP = () => {
  return crypto.randomInt(100000, 999999).toString();
};

// Register route
router.post('/register', async (req, res) => {
  const { username, email, password } = req.body;
  try {
    let user = await User.findOne({ email });
    if (user) return res.status(400).json({ msg: 'User already exists' });

    user = new User({ username, email, password });
    await user.save();

    // Generate OTP
    const otp = generateOTP();
    const otpEntry = new OTP({ userId: user._id, otp });
    await otpEntry.save();

    // Send OTP via email
    await transporter.sendMail({
      to: user.email,
      subject: 'Your OTP Code',
      text: `Your OTP code is ${otp}. It will expire in 5 minutes.`
    });

    res.status(200).json({ msg: 'OTP sent to your email' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
  }
});

// Verify OTP
router.post('/verify-otp', async (req, res) => {
  const { email, otp } = req.body;
  try {
    const user = await User.findOne({ email });
    if (!user) return res.status(400).json({ msg: 'User not found' });

    const otpEntry = await OTP.findOne({ userId: user._id, otp });
    if (!otpEntry) return res.status(400).json({ msg: 'Invalid or expired OTP' });

    user.isVerified = true;
    await user.save();
    await OTP.deleteOne({ _id: otpEntry._id });

    res.status(200).json({ msg: 'Account verified successfully' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
  }
});

// Resend OTP
router.post('/resend-otp', async (req, res) => {
  const { email } = req.body;
  try {
    const user = await User.findOne({ email });
    if (!user) return res.status(400).json({ msg: 'User not found' });

    const otp = generateOTP();
    const otpEntry = new OTP({ userId: user._id, otp });
    await otpEntry.save();

    await transporter.sendMail({
      to: user.email,
      subject: 'Your OTP Code',
      text: `Your OTP code is ${otp}. It will expire in 5 minutes.`
    });

    res.status(200).json({ msg: 'OTP resent to your email' });
  } catch (err) {
    console.error(err.message);
    res.status(500).send('Server error');
  }
});

module.exports = router;
