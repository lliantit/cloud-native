const errorHandler = (err, req, res, next) => {
    if (err.message === 'user_not_found') {
        return res.status(404).json({ error: 'not_found', message: 'User was not found' });
    }

    if (err.message === 'invalid_data') {
        return res.status(400).json({ error: 'bad_request', message: 'Invalid input data' });
    }

    if (err.message === 'conflict') {
        return res.status(409).json({ error: 'conflict', message: 'Resource already exists' });
    }

    // Якщо помилка невідома
    console.error(err);
    res.status(500).json({ error: 'internal_error', message: 'Internal Server Error' });
};

module.exports = errorHandler;