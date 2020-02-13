const orderInfo = document.querySelector(".orderInfo");
var form = document.getElementById('payment-form');
let sercret = "";
const URL_API = "http://127.0.0.1:8000";

var url_string = window.location.href;
var url = new URL(url_string);
var c = url.searchParams.get("id");
console.log(c);

let productID = "";

if (c == "" || c == null) {
    window.location = "../"
}

window.addEventListener("load", () => {

    // is user logged in

    fetch(`${URL_API}/loggedIn`, {
        method: 'get',
        credentials: 'include',
        headers: {
            "Content-type": "application/json",
        }
    })
    .then(response => {
        console.log(response.status); // Will show you the status
            if (!response.ok) {
                if (confirm("You need to login")) {
                    window.location = "../signin";
                } else {
                    window.location = "../signin";
                }
            }
            return response.json();
    })
    .then (data => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
        console.log(response.status); 
    });

    refershToken()

    fetch(`http://localhost:8000/products/${c}`)
    .then(res => {
        console.log(res.status); // Will show you the status
        if (!res.ok) {
            window.location = "../"
        }
        return res.json();
    })
    .then(data => {
        console.log(data);
        productID = data[0].uuid;
        addProdToPage(data[0]);
    })
    .catch(err => console.log(err))
});

function addProdToPage(product) {
    let name = document.createElement("p");
    name.innerText = product.name;

    let spec = document.createElement("h3");
    spec.innerText = "Specs"

    let cpu = document.createElement("p");
    cpu.innerText = `Number of CPU cores : ${product.cpu} cores`;

    let ram = document.createElement("p");
    ram.innerText = `Amount of RAM : ${product.ram} GB`;

    let disk = document.createElement("p");
    disk.innerText = `Disk space : ${product.disk} GB`;

    let price = document.createElement("p");
    price.innerText = `Price : £ ${product.price}`;

    orderInfo.appendChild(name);
    orderInfo.appendChild(spec);
    orderInfo.appendChild(cpu);
    orderInfo.appendChild(ram);
    orderInfo.appendChild(disk);
    orderInfo.appendChild(price);


}


    var stripe = Stripe('pk_test_ERYWSEs8exlFbm3glnzDeiga00VmESFxNg');
    var elements = stripe.elements();
    
    var response = fetch(`http://localhost:8000/create-payment-intent/${c}`).then(function(response) {
        return response.json();
    }).then(function(responseJson) {
        var clientSecret = responseJson.clientecret;
        sercret = clientSecret;
    });
    
    var style = {
        base: {
        color: "#32325d",
        fontFamily: '"Helvetica Neue", Helvetica, sans-serif',
        fontSmoothing: "antialiased",
        fontSize: "16px",
        "::placeholder": {
        color: "#aab7c4"
        }
    },
    invalid: {
        color: "#fa755a",
        iconColor: "#fa755a"
    }};
    
    var card = elements.create("card", { style: style });
    card.mount("#card-element");
    
    card.addEventListener('change', ({error}) => {
    const displayError = document.getElementById('card-errors');
    if (error) {
        displayError.textContent = error.message;
    } else {
        displayError.textContent = '';
    }});


form.addEventListener('submit', function(ev) {
    ev.preventDefault();
    stripe.confirmCardPayment(sercret, {
        payment_method: {
        card: card,
        billing_details: {
            name: 'Jenny Rosen'
        }}
    }).then(function(result) {
        if (result.error) {
            console.log(result.error.message);
        } else {
        if (result.paymentIntent.status === 'succeeded') {
            console.log("payment made yaaaa");
            console.log(result);

            // add order to DB
            createOrder(result.paymentIntent.id, productID);
        }}
});
});

setInterval(function() {
    refershToken();
}, 60000); // Every 1 minitue


function refershToken() {
    fetch(`${URL_API}/refresh`, {
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
        console.log(response.status); 
    });
}

function createOrder(payID, prodID) {

    let formData = {
        "PaymentID" : payID,
        "ProductID" : prodID
    }

    console.log(formData);
    

    fetch(`${URL_API}/order`, {
        method: 'post',
        credentials: 'include',
        headers: {
            "Content-type": "application/json",
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then((data) => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
    });
}