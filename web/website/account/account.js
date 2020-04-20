const orderInfo = document.querySelector(".orderInfo");
const userData = document.querySelector(".userData");

const URL_API = "http://127.0.0.1:8000";

refershToken();

fetch(`${URL_API}/auth/isloggedin`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => {
    console.log(response.status); // Will show you the status
        if (!response.ok) {
                window.location = "../signin";
        }
        return response.json();
})
.then (data => console.log(data))
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
});

fetch(`${URL_API}/account/accountinfo`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => response.json())
.then (data => {
    console.log(data);
    let name = document.createElement("p")
    name.innerText = "Name : " + data[0].Name;

    let email = document.createElement("p")
    email.innerText = "Email : " + data[0].Email;

    userData.appendChild(name);
    userData.appendChild(email);
})
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
});

fetch(`${URL_API}/getorders`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => response.json())
.then (data => {
    console.log(data);

    for (let i = 0; i < data.length; i++) {
        // console.log(data[i].OrderID);
        getProducts(data[i].ProdID, data[i].OrderID);
    }
})
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
});

setInterval(function() {
    refershToken();
}, 60000); // Every 1 minitue

function refershToken() {
    fetch(`${URL_API}/auth/refresh`, {
        method: 'get',
        credentials: 'include',
        headers: {
            "Content-type": "application/json",
        }
    })
    .then(response => response.json())
    .then (data => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
    });
}

function getProducts(productID, OrderID) {
    fetch(`${URL_API}/products/${productID}`)
    .then(res => {
        // console.log(res.status); // Will show you the status
        if (!res.ok) {
            console.log(res.status + " product not found");  
        }
        return res.json();
    })
    .then(data => {
        // console.log(data);
        productID = data[0].uuid;
        addProdToPage(data[0], OrderID);
    })
    .catch(err => console.log(err))
}

// http://127.0.0.1:8080/web/website/receipts/index.html?id=order_1ZMD18zw0NUulLXt2ftFpvT9QTl

function addProdToPage(product, OrderID) {
    let orderidtitle = document.createElement("p");
    orderidtitle.innerText = "Order ID : ";

    let orderid = document.createElement("a");
    orderid.innerText = `${OrderID}`;
    orderid.href = `../receipts/index.html?id=${OrderID}`

    let name = document.createElement("p");
    name.innerText = `Name : ${product.name}`;

    let spec = document.createElement("h3");
    spec.innerText = "Specs"

    let cpu = document.createElement("p");
    cpu.innerText = `Number of CPU cores : ${product.cpu} cores`;

    let ram = document.createElement("p");
    ram.innerText = `Amount of RAM : ${product.ram} GB`;

    let disk = document.createElement("p");
    disk.innerText = `Disk space : ${product.disk} GB`;

    let price = document.createElement("p");
    var formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'GBP',
      });
    let priceFM = formatter.format(product.price);
    price.innerText = `Price : ${priceFM} per a month`;

    const orderdiv = document.createElement("div");
    orderdiv.className = "orderdiv";


    orderidtitle.appendChild(orderid);

    orderdiv.appendChild(name);
    orderdiv.appendChild(orderidtitle);
    orderdiv.appendChild(spec);
    orderdiv.appendChild(cpu);
    orderdiv.appendChild(ram);
    orderdiv.appendChild(disk);
    orderdiv.appendChild(price);

    orderInfo.appendChild(orderdiv);


}