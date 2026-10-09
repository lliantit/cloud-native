class UserController {
    constructor(userService) {
        this.userService = userService;
    }

    async getAll(req, res, next) {
        try {
            const users = await this.userService.getAllUsers();
            res.json(users);
        } catch (error) {
            next(error);
        }
    }

    async getById(req, res, next) {
        try {
            const id = parseInt(req.params.id);
            if (isNaN(id)) return res.status(400).json({ error: 'bad_request', message: 'Invalid identifier' });

            const user = await this.userService.getUserById(id);
            res.json(user);
        } catch (error) {
            next(error);
        }
    }

    async create(req, res, next) {
        try {
            const user = await this.userService.createUser(req.body);
            res.status(201).json(user);
        } catch (error) {
            next(error);
        }
    }

    async update(req, res, next) {
        try {
            const id = parseInt(req.params.id);
            if (isNaN(id)) return res.status(400).json({ error: 'bad_request', message: 'Invalid identifier' });

            const user = await this.userService.updateUser(id, req.body);
            res.json(user);
        } catch (error) {
            next(error);
        }
    }

    async delete(req, res, next) {
        try {
            const id = parseInt(req.params.id);
            if (isNaN(id)) return res.status(400).json({ error: 'bad_request', message: 'Invalid identifier' });

            await this.userService.deleteUser(id);
            res.status(204).send();
        } catch (error) {
            next(error);
        }
    }
}

module.exports = UserController;