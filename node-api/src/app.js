const express = require('express');
const userRoutes = require('./routes/user.routes');
const logger = require('./middleware/logger');
const errorHandler = require('./middleware/error-handler');

const app = express();

// Мідлвари рівня застосунку
app.use(express.json());
app.use(logger);

// Ендпойнт для перевірки здоров'я
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// Підключення маршрутів
app.use('/api/users', userRoutes);

// Глобальний обробник помилок завжди має бути останнім
app.use(errorHandler);

// Запуск сервера, якщо файл викликається напряму (не через тести)
if (require.main === module) {
    const PORT = process.env.PORT || 3000;
    app.listen(PORT, () => {
        console.log(`Node.js server is running on port ${PORT}`);
    });
}

module.exports = app;