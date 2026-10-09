const UserService = require('../src/services/user.service');
const InMemoryUserRepository = require('../src/repositories/user.repository');

describe('UserService', () => {
    let service;

    beforeEach(() => {
        // Кожен тест отримує чисте сховище
        const repository = new InMemoryUserRepository();
        service = new UserService(repository);
    });

    test('повинен створювати користувача', async () => {
        const user = await service.createUser({ name: 'John', email: 'john@test.com' });
        expect(user.id).toBe(1);
        expect(user.name).toBe('John');
    });

    test('повинен викидати помилку при відсутності email (бізнес-правило)', async () => {
        await expect(service.createUser({ name: 'John' })).rejects.toThrow('invalid_data');
    });

    test('повинен знаходити створеного користувача', async () => {
        const created = await service.createUser({ name: 'Anna', email: 'anna@test.com' });
        const found = await service.getUserById(created.id);
        expect(found.email).toBe('anna@test.com');
    });

    test('повинен оновлювати дані користувача', async () => {
        const user = await service.createUser({ name: 'Max', email: 'max@test.com' });
        const updated = await service.updateUser(user.id, { name: 'Maxim' });
        expect(updated.name).toBe('Maxim');
    });

    test('повинен видаляти користувача', async () => {
        const user = await service.createUser({ name: 'Tom', email: 'tom@test.com' });
        await service.deleteUser(user.id);
        await expect(service.getUserById(user.id)).rejects.toThrow('user_not_found');
    });
});