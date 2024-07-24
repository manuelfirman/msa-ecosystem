const userService = require('../services/userService');

exports.createUser = async (req, res) => {
    const { username, email, first_name, last_name, unique_id } = req.body;
    try {
        const userId = await userService.createUser(username, email, first_name, last_name, unique_id);
        res.status(201).json({ id: userId });
    } catch (err) {
        console.error('Error creating user:', err.message);
        res.status(500).json({ error: err.message });
    }
};

exports.getUsers = async (req, res) => {
    try {
        const users = await userService.getUsers();
        res.json(users);
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
};

exports.getUserById = async (req, res) => {
    const { id } = req.params;
    try {
        const user = await userService.getUserById(id);
        if (user) {
            res.json(user);
        } else {
            res.status(404).json({ error: 'User not found' });
        }
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
};

exports.updateUser = async (req, res) => {
    const { id } = req.params;
    const { first_name, last_name } = req.body;
    try {
        const userId = await userService.updateUser(id, first_name, last_name);
        if (userId) {
            res.json({ id: userId });
        } else {
            res.status(404).json({ error: 'User not found' });
        }
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
};

exports.deleteUser = async (req, res) => {
    const { id } = req.params;
    try {
        const deleted = await userService.deleteUser(id);
        if (deleted) {
            res.json({ id });
        } else {
            res.status(404).json({ error: 'User not found' });
        }
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
};