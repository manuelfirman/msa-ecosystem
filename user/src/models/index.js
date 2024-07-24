const { Sequelize } = require('sequelize');
const config = require('../config/config');

// Configuración de Sequelize
const sequelize = new Sequelize(
    config.db.database, 
    config.db.user, 
    config.db.password, 
    {
        host: config.db.host,
        dialect: 'mysql',
    }
);

const UserModel = require('./userModel')(sequelize, Sequelize.DataTypes);

const db = {
    sequelize,
    UserModel
};

module.exports = db;
