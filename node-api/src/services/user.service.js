class UserService {
    constructor(userRepository) {
        this.repository = userRepository;
    }

    async getAllUsers() {
        return await this.repository.findAll();
    }

    async getUserById(id) {
        const user = await this.repository.findById(id);
        if (!user) throw new Error('user_not_found');
        return user;
    }

    async createUser(userData) {
        if (!userData.name || !userData.email) {
            throw new Error('invalid_data');
        }
        return await this.repository.create(userData);
    }

    async updateUser(id, userData) {
        const updatedUser = await this.repository.update(id, userData);
        if (!updatedUser) throw new Error('user_not_found');
        return updatedUser;
    }

    async deleteUser(id) {
        const deleted = await this.repository.delete(id);
        if (!deleted) throw new Error('user_not_found');
        return deleted;
    }
}

module.exports = UserService;
