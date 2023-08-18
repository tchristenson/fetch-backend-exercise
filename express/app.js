const express = require('express');
const app = express();

app.get('/', (req, res) => {
    res.status(200).send('Home Page')
})

app.all('*', (req, res) => {
    res.status(404).send('Resource not found')
})








const port = 5000
app.listen(port, () => {
    console.log(`Server is running on ${port}...`)
})
