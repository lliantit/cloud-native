class InMemoryUserRepository {
    constructor() {
        this.users = [];
        this.currentId = 1;
    }

    async findAll() {
        return this.users;
    }

    async findById(id) {
        return this.users.find(user => user.id === id);
    }

    async create(userData) {
        const newUser = { id: this.currentId++, ...userData };
        this.users.push(newUser);
        return newUser;
    }

    async update(id, userData) {
        const index = this.users.findIndex(user => user.id === id);
        if (index === -1) return null;
        this.users[index] = { ...this.users[index], ...userData };
        return this.users[index];
    }

    async delete(id) {
        const index = this.users.findIndex(user => user.id === id);
        if (index === -1) return false;
        this.users.splice(index, 1);
        return true;
    }
}

module.exports = InMemoryUserRepository;