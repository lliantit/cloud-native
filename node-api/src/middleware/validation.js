const validateUser = (req, res, next) => {
    const { name, email } = req.body;

    if (!name || typeof name !== 'string' || name.trim() === '') {
        return res.status(400).json({ error: 'bad_request', message: 'Name is required and must be a string' });
    }

    if (!email || typeof email !== 'string' || !email.includes('@')) {
        return res.status(400).json({ error: 'bad_request', message: 'Valid email is required' });
    }

    next();
};

module.exports = validateUser;
