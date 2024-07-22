// user/index.js
const express = require('express');
const app = express();
app.use(express.json());

let users = [];

app.post('/users', (req, res) => {
    const user = req.body;
    users.push(user);
    res.status(201).send(user);
});

app.get('/users', (req, res) => {
    res.send(users);
});

app.listen(5001, () => {
    console.log('User service listening on port 5001');
});
