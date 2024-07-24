const amqp = require('amqplib/callback_api');
const pool = require('../config/db');

amqp.connect('amqp://rabbitmq:5672', (error0, connection) => {
    if (error0) {
        console.error('Failed to connect to RabbitMQ:', error0.message);
        process.exit(1); // Exit the process if RabbitMQ connection fails
    }
    connection.createChannel((error1, channel) => {
        if (error1) {
            console.error('Failed to create a channel:', error1.message);
            process.exit(1);
        }
        const queue = 'user_created';

        channel.assertQueue(queue, {
            durable: false
        });

        channel.consume(queue, async (msg) => {
            if (msg !== null) {
                try {
                    const user = JSON.parse(msg.content.toString());
                    const [results] = await pool.query(
                        'INSERT INTO users (username, email, first_name, last_name, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())',
                        [user.username, user.email, user.first_name, user.last_name]
                    );
                    console.log('User created:', results.insertId);
                    channel.ack(msg); // Acknowledge the message after processing
                } catch (err) {
                    console.error('Database insert error:', err.message);
                }
            }
        }, {
            noAck: false // Ensure to acknowledge the message after processing
        });
    });
});
