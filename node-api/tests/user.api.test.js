const request = require('supertest');
const app = require('../src/app');

describe('User API Endpoints', () => {
    test('POST /api/users - успішне створення', async () => {
        const res = await request(app)
            .post('/api/users')
            .send({ name: 'Jane Doe', email: 'jane@example.com' });

        expect(res.statusCode).toBe(201);
        expect(res.body).toHaveProperty('id');
        expect(res.body.name).toBe('Jane Doe');
    });

    test('POST /api/users - некоректні дані (помилка валідації)', async () => {
        const res = await request(app)
            .post('/api/users')
            .send({ name: 'Jane Doe' }); // Немає email

        expect(res.statusCode).toBe(400);
        expect(res.body.error).toBe('bad_request');
    });

    test('GET /api/users/:id - отримання неіснуючого ресурсу', async () => {
        const res = await request(app).get('/api/users/999');
        expect(res.statusCode).toBe(404);
        expect(res.body.error).toBe('not_found');
    });

    test('GET /api/users - успішний запит списку', async () => {
        const res = await request(app).get('/api/users');
        expect(res.statusCode).toBe(200);
        expect(Array.isArray(res.body)).toBeTruthy();
    });

    test('DELETE /api/users/:id - успішне видалення', async () => {
        // Спочатку створюємо
        const createRes = await request(app)
            .post('/api/users')
            .send({ name: 'Mark', email: 'mark@test.com' });

        const userId = createRes.body.id;

        // Потім видаляємо
        const deleteRes = await request(app).delete(`/api/users/${userId}`);
        expect(deleteRes.statusCode).toBe(204);
    });
});