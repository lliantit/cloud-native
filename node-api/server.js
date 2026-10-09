const express = require('express');
const app = express();

app.use(express.json());

let movies = [];
let currentId = 1;

// Health check
app.get('/health', (req, res) => {
    res.status(200).json({ status: 'ok', uptime: process.uptime() });
});

// Create
app.post('/movies', (req, res) => {
    const movie = { id: currentId++, ...req.body };
    movies.push(movie);
    res.status(201).json(movie);
});

// Read All
app.get('/movies', (req, res) => {
    res.status(200).json(movies);
});

// Read One
app.get('/movies/:id', (req, res) => {
    const movie = movies.find(m => m.id === parseInt(req.params.id));
    if (!movie) {
        return res.status(404).json({ error: 'Фільм не знайдено' });
    }
    res.status(200).json(movie);
});

// Update
app.put('/movies/:id', (req, res) => {
    const movie = movies.find(m => m.id === parseInt(req.params.id));
    if (!movie) {
        return res.status(404).json({ error: 'Фільм не знайдено' });
    }
    Object.assign(movie, req.body);
    res.status(200).json(movie);
});

// Delete
app.delete('/movies/:id', (req, res) => {
    const index = movies.findIndex(m => m.id === parseInt(req.params.id));
    if (index === -1) {
        return res.status(404).json({ error: 'Фільм не знайдено' });
    }
    movies.splice(index, 1);
    res.status(204).send();
});

app.get('/io', (req, res) => {
    setTimeout(() => {
        res.status(200).json({ message: 'I/O operation completed' });
    }, 2000);
});

app.get('/cpu', (req, res) => {
    let count = 0;
    for (let i = 0; i < 500_000_000; i++) {
        count++;
    }
    res.status(200).json({ message: 'CPU operation completed', count });
});




const PORT = 5000;
app.listen(PORT, () => {
    console.log(`Node.js сервер запущено на http://localhost:${PORT}`);
});