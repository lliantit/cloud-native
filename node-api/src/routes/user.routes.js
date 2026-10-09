const express = require('express');
const InMemoryUserRepository = require('../repositories/user.repository');
const UserService = require('../services/user.service');
const UserController = require('../controllers/user.controller');
const validateUser = require('../middleware/validation');

const router = express.Router();

const repository = new InMemoryUserRepository();
const service = new UserService(repository);
const controller = new UserController(service);

router.get('/', (req, res, next) => controller.getAll(req, res, next));
router.get('/:id', (req, res, next) => controller.getById(req, res, next));
// Додаємо validateUser перед викликом контролера
router.post('/', validateUser, (req, res, next) => controller.create(req, res, next));
router.put('/:id', validateUser, (req, res, next) => controller.update(req, res, next));
router.delete('/:id', (req, res, next) => controller.delete(req, res, next));

module.exports = router;