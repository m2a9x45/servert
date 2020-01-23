
fetch("http://localhost:8080/products")
    .then(res => res.json())
    .then(data => console.log(data))
    .catch(err => console.log(err))