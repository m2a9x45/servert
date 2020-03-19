//http://127.0.0.1:8080/web/website/receipts/index.html?id=order_1ZMD18zw0NUulLXt2ftFpvT9QTl

const table = document.getElementById("table");
const orderIDtext = document.getElementById("orderID");
const payID = document.getElementById("payID");
const LoginButton = document.querySelector('#Login');
const SignupButton = document.querySelector('#Signup');
const dateOfPur = document.getElementById("dateOfPur");


const URL_API = "http://127.0.0.1:8000";

let url_string = window.location.href;
let url = new URL(url_string);
let orderID = url.searchParams.get("id");
console.log(orderID);

if (orderID == "" || orderID == null) {
    window.location = "../"
}

// get order ID call DB and get the prod_id, transaction ID, maybe add time at some point

refershToken();

fetch(`${URL_API}/auth/isloggedin`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => {
    console.log(response.status);
    return response.json();
})
.then (data => {
    console.log(data);
    if (data.success == true) {
        LoginButton.innerText = "Account"
        LoginButton.href = "../account"
        SignupButton.style.display = "none"
    } else {
        LoginButton.innerText = "Login"
        LoginButton.href = "../signin/"
        SignupButton.style.display = "block"
    }
})
.catch(function (error) {
    console.log('Request failed', error);
});

fetch(`${URL_API}/account/receipt/${orderID}`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => {
        return response.json();
})
.then (data => {
    console.log(data);
    displayinfo(data);
    displayTable(data);

})
.catch(function (error) {
    console.log('Request failed', error);
});

function displayinfo(data) {
    orderIDtext.innerText = `Order ID : ${data.order.OrderID}`;
    payID.innerText = `Payment ID : ${data.order.PaymentID}`;

    const d = new Date(data.order.Time * 1000);
    const strDate = d.toLocaleString('en-GB', {dateStyle: "short", timeStyle : "short"});


    dateOfPur.innerText = `Date of purchase : ${strDate}`;

}

function displayTable(data) {

    let prodName = document.createElement("td");
    prodName.innerText = data.product.name;

    let rentDur = document.createElement("td");
    rentDur.innerText = "1 month";

    let tr1 = document.createElement("tr");

    let prodCpulabel = document.createElement("td");
    prodCpulabel.innerText = "Cpu";

    let prodCpu = document.createElement("td");
    prodCpu.innerText = data.product.cpu;

    let tr2 = document.createElement("tr");

    let prodRamlabel = document.createElement("td");
    prodRamlabel.innerText = "Ram";

    let prodRam = document.createElement("td");
    prodRam.innerText = data.product.ram;

    let tr3 = document.createElement("tr");

    let prodDisklabel = document.createElement("td");
    prodDisklabel.innerText = "Disk";

    let prodDisk = document.createElement("td");
    prodDisk.innerText = data.product.disk;

    let tr4 = document.createElement("tr");

    let prodTotallabel = document.createElement("td");
    prodTotallabel.innerText = "Total";

    let prodTotal = document.createElement("td");
    var formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'GBP',
      });
    let price = formatter.format(data.product.price);
    prodTotal.innerText = price;

    let tr5 = document.createElement("tr");

    tr1.appendChild(prodName);
    tr1.appendChild(rentDur);

    tr2.appendChild(prodCpulabel);
    tr2.appendChild(prodCpu);

    tr3.appendChild(prodRamlabel);
    tr3.appendChild(prodRam);

    tr4.appendChild(prodDisklabel);
    tr4.appendChild(prodDisk);

    tr5.appendChild(prodTotallabel);
    tr5.appendChild(prodTotal);

    
    table.appendChild(tr1);
    table.appendChild(tr2);
    table.appendChild(tr3);
    table.appendChild(tr4);
    table.appendChild(tr5);

}

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
    .then(response => {
        if (!response.ok) {
            return
        }
        return response.json();
    })
    .then (data => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
    });
}