const express = require('express');
const userRoutes = require('./routes/userRoutes');
const requestIdMiddleware = require('./middleware/requestIdMiddleware');
const config = require('./config/config');
const db = require('./models');
require('dotenv').config();

const app = express();
app.use(express.json());
app.use(requestIdMiddleware);
app.use(userRoutes);

// Sincronizar los modelos con la base de datos
db.sequelize.sync()
    .then(() => {
        app.listen(config.port, () => {
            console.log(`User service listening on port ${config.port}`);
        });
    })
    .catch(error => {
        console.error('Unable to connect to the database:', error);
    });
