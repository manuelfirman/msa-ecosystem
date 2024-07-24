const userRepository = require('../repositories/userRepository');

exports.createUser = async (username, email, first_name, last_name, unique_id) => {
    return userRepository.createUser(username, email, first_name, last_name, unique_id);
};

exports.getUsers = async () => {
    return userRepository.getUsers();
};

exports.getUserById = async (id) => {
    return userRepository.getUserById(id);
}

exports.updateUser = async (id, first_name, last_name) => {
    return userRepository.updateUser(id, first_name, last_name);
}

exports.deleteUser = async (id) => {
    return userRepository.deleteUser(id);
}

