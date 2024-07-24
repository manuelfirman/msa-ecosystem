// user/index.js
const express = require('express');
const mysql = require('mysql2/promise');
require('dotenv').config();
// const amqp = require('amqplib/callback_api');
const app = express();
app.use(express.json());

// const PORT = process.env.PORT || 5001;

// const pool = mysql.createPool({
//     host: process.env.DB_HOST,
//     user: process.env.DB_USER,
//     password: process.env.DB_PASSWORD,
//     database: process.env.DB_NAME,
//     waitForConnections: true,
//     connectionLimit: 10,
//     queueLimit: 0
//   });

app.use((req, res, next) => {
req.requestId = req.headers['x-request-id'];
if (!req.requestId) {
    res.status(400).json({ error: 'Missing X-Request-ID header' });
    return;
}
res.setHeader('X-Request-ID', req.requestId);
next();
});

app.post('/users', async (req, res) => {
    const { username, email, first_name, last_name, unique_id } = req.body;
    try {
        const [results] = await pool.query(
            'INSERT INTO users (username, email, first_name, last_name, unique_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())',
            [username, email, first_name, last_name, unique_id]
        );
        console.log('User created:', results.insertId);
        res.status(201).json({ id: results.insertId });
    } catch (err) {
        console.error('Database insert error:', err.message);
        res.status(500).json({ error: err.message });
    }
});
  
app.get('/users', async (req, res) => {
    try {
        const [rows] = await pool.query('SELECT * FROM users');
        res.json(rows);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.listen(PORT, () => {
    console.log(`User service listening on port ${PORT}`);
});


// codigo para consumir mensagens da fila
// amqp.connect('amqp://rabbitmq:5672', (error0, connection) => {
//     if (error0) {
//         console.error('Failed to connect to RabbitMQ:', error0.message);
//         process.exit(1); // Exit the process if RabbitMQ connection fails
//     }
//     connection.createChannel((error1, channel) => {
//         if (error1) {
//             console.error('Failed to create a channel:', error1.message);
//             process.exit(1);
//         }
//         const queue = 'user_created';

//         channel.assertQueue(queue, {
//             durable: false
//         });

//         channel.consume(queue, async (msg) => {
//             if (msg !== null) {
//                 try {
//                     const user = JSON.parse(msg.content.toString());
//                     const [results] = await pool.query(
//                         'INSERT INTO users (username, email, first_name, last_name, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())',
//                         [user.username, user.email, user.first_name, user.last_name]
//                     );
//                     console.log('User created:', results.insertId);
//                     channel.ack(msg); // Acknowledge the message after processing
//                 } catch (err) {
//                     console.error('Database insert error:', err.message);
//                 }
//             }
//         }, {
//             noAck: false // Ensure to acknowledge the message after processing
//         });
//     });
// });