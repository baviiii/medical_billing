const express = require('express');
const bodyParser = require('body-parser');
const connectDB = require('./config/db');
const authRoutes = require('./routes/auth');

const app = express();
connectDB();

app.use(bodyParser.json());
app.use('/api/auth', authRoutes);

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => console.log(`Server running on port ${PORT}`));
