const express = require('express');
const redis = require('redis');
const bodyParser = require('body-parser');

const app = express();
const client = redis.createClient({ url: 'redis://redis:6379' });

app.use(bodyParser.json());

app.post('/signup', async (req, res) => {
    const { email, password } = req.body;
    if (!email || !password) return res.status(400).send("Invalid input");

    const key = `user:${email}`;
    const exists = await client.exists(key);
    if (exists) return res.status(409).send("User exists");

    await client.hSet(key, 'email', email, 'password', password);
    res.status(201).send("User created");
});

client.connect().then(() => {
    app.listen(8081, () => console.log("Node.js server on :8081"));
});