const db = require('../models');

exports.createUser = async (username, email, first_name, last_name, unique_id) => {
    const user = await db.UserModel.create({
        username,
        email,
        first_name,
        last_name,
        unique_id
    });
    return user.id;
};

exports.getUsers = async () => {
    return await db.UserModel.findAll();
};

exports.getUserById = async (id) => {
    return await db.UserModel.findByPk(id);
}

exports.updateUser = async (id, first_name, last_name) => {
    const user = await db.UserModel.findByPk(id);
    user.first_name = first_name;
    user.last_name = last_name;
    await user.save();
    return user.id;
}

exports.deleteUser = async (id) => {
    return await db.UserModel.destroy({
        where: {
            id
        }
    });
}


