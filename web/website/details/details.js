const form = document.getElementById('payment-form');
const orderInfo = document.querySelector(".orderInfo");
const rentdurButtons = document.querySelectorAll(".rentrange");
const rent = document.getElementsByName("rent");
const cardInfo = document.querySelector(".cardInfo");
const newCard = document.querySelector("#newCard");
const payment = document.querySelector(".payment");
const buybutton = document.querySelector("#buybutton");
const radios = document.getElementsByName('card');
const selectCardText = document.querySelector("#selectCardText");
const terms = document.querySelector("#terms");
const recurringPayments = document.querySelector("#recurringPayments");
const recurringPaymentsDiv = document.querySelector(".recurringPaymentsDiv");
const errorDiv = document.querySelector(".errorDiv");
const errorMessage = document.querySelector("#errorMessage");

const URL_API = "http://127.0.0.1:8000";


let duration = "";
let custName = "";

var url_string = window.location.href;
var url = new URL(url_string);
var c = url.searchParams.get("id");
console.log(c);

let productID = "";

if (c == "" || c == null) {
    window.location = "../"
}



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
    .then(data => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
        console.log(response.status);
    });

refershToken()

fetch(`${URL_API}/products/${c}`)
    .then(res => {
        console.log(res.status); // Will show you the status
        if (!res.ok) {
            window.location = "../"
        }
        return res.json();
    })
    .then(data => {
        console.log(data);
        productID = data[0];
        addProdToPage(data[0]);
    })
    .catch(err => console.log(err))

rentdurButtons.forEach(button => {
    button.addEventListener("click", (e) => {
        // console.log(e);
        if (e.target.className == "rentrange") {

            rentdurButtons.forEach(element => {
                element.style.borderColor = "#0e1c31";
                element.style.backgroundColor = "#fafafa";
                // e.target.style.opacity = 1;
            });

            e.srcElement.children[0].checked = true;
            e.target.style.borderColor = "#2a4f87";
            e.target.style.backgroundColor = "#e6e6e6";
            // e.target.style.opacity = 0.5;

        }

        if (e.target.tagName == "LABEL") {

            rentdurButtons.forEach(element => {
                element.style.borderColor = "#0e1c31";
                element.style.backgroundColor = "#fafafa";
            });

            e.srcElement.previousElementSibling.checked = true;

            e.srcElement.parentElement.style.borderColor = "#2a4f87";
            e.srcElement.parentElement.style.backgroundColor = "#e6e6e6";
        }

        if (e.target.tagName == "INPUT") {

            rentdurButtons.forEach(element => {
                element.style.borderColor = "#0e1c31";
                element.style.backgroundColor = "#fafafa";
            });

            e.target.checked = true;

            e.target.parentElement.style.borderColor = "#2a4f87";
            e.target.parentElement.style.backgroundColor = "#e6e6e6";
        }

        rent.forEach(button => {
            if (button.checked) {
                // console.log(button.value);
                duration = button.value;
                // console.log(duration);

            }
        })

        orderInfo.innerHTML = "";
        addProdToPage(productID);

    });
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
    custName = data[0].Name;
})
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
});

fetch(`${URL_API}/account/customercards`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => response.json())
.then (data => {
    console.log(data);

    if (data.length > 0) {

        form.style.display = "none";
        recurringPaymentsDiv.style.display = "none";

        for (let i = 0; i < data.length; i++) {
            addcard(data[i])
        }
    } else {
        buybutton.style.display = "none";
        newCard.style.display = "none";
        selectCardText.style.display = "none";
        

    }  
})
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
});


newCard.addEventListener("click", () => {
    event.preventDefault();

    cardInfo.style.display = "none";
    form.style.display = "block";
    payment.style.display = "none";
    newCard.style.display = "none";

})

var stripe = Stripe('pk_test_ERYWSEs8exlFbm3glnzDeiga00VmESFxNg');
var elements = stripe.elements();

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
    }
};

var card = elements.create("card", {
    style: style
});
card.mount("#card-element");

card.addEventListener('change', ({
    error
}) => {
    const displayError = document.getElementById('card-errors');
    if (error) {
        displayError.textContent = error.message;
    } else {
        displayError.textContent = '';
    }
});

function addProdToPage(product) {
    let name = document.createElement("p");
    name.innerText = product.name;

    // let spec = document.createElement("h3");
    // spec.innerText = "Specs"

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
    let pricef = formatter.format(product.price * duration);
    price.innerText = `Price : ${pricef}`;

    orderInfo.appendChild(name);
    // orderInfo.appendChild(spec);
    orderInfo.appendChild(cpu);
    orderInfo.appendChild(ram);
    orderInfo.appendChild(disk);
    orderInfo.appendChild(price);
}

form.addEventListener('submit', function (ev) {
    ev.preventDefault();
    const dur = duration;
    canSubmit = checkData(dur, terms.checked); 
    console.log(canSubmit);
   
    if (canSubmit == true) {
        var response = fetch(`${URL_API}/order/create-payment-intent/${c}/${dur}/`, {
            method: "GET",
            credentials: "include",
        }).then(function (response) {
            return response.json();
        }).then(function (responseJson) {
            stripe.confirmCardPayment(responseJson.clientecret, {
                payment_method: {
                    card: card,
                    billing_details: {
                        name: custName
                    }
                },
                setup_future_usage: recurringPayments.checked ? 'off_session' : ''
            }).then(function (result) {
                if (result.error) {
                    console.log(result.error.message);
                } else {
                    if (result.paymentIntent.status === 'succeeded') {
                        console.log("payment made yaaaa");
                        console.log(result);
    
                        // add order to DB
                        console.log(productID);
    
                        createOrder(result.paymentIntent.id, productID.uuid, dur);
                    }
                }
            });
        });  
    }

});

buybutton.addEventListener('click', () => {
    event.preventDefault(); 

    let cardID = "";

    for (var i = 0, length = radios.length; i < length; i++) {
        if (radios[i].checked) {
            cardID = radios[i].value;
            break;
        }
    }

    const dur = duration;

    canSubmit = checkData(dur, terms.checked); 
    console.log(canSubmit);
   
    if (canSubmit == true) {
        var response = fetch(`${URL_API}/order/create-payment-intent/${c}/${dur}/${cardID}`, {
            method: "GET",
            credentials: "include",
        }).then(function (response) {
            return response.json();
        }).then(function (responseJson) {
            console.log(responseJson);
            stripe
                .confirmCardPayment(responseJson.clientecret, {
                    payment_method: cardID,
                })
                .then(function(result) {
                    // Handle result.error or result.paymentIntent
                    if (result.error) {
                        console.log(result.error.message);
                    } else {
                        if (result.paymentIntent.status === 'succeeded') {
                            console.log("payment made yaaaa");
                            console.log(result);
        
                            // add order to DB
                            console.log(productID);
        
                            createOrder(result.paymentIntent.id, productID.uuid, dur);  
                        }
                    }
                });
        });
    }
})

setInterval(function () {
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
        .then(data => console.log(data))
        .catch(function (error) {
            console.log('Request failed', error);
            console.log(response.status);
        });
}

function createOrder(payID, prodID, dur) {

    let formData = {
        "PaymentID": payID,
        "ProductID": prodID,
        "Dur": dur,
    }

    console.log(formData);


    fetch(`${URL_API}/order/makeorder`, {
            method: 'post',
            credentials: 'include',
            headers: {
                "Content-type": "application/json",
            },
            body: JSON.stringify(formData)
        })
        .then(response => response.json())
        .then((data) => {
            console.log(data);
            window.location = `../receipts/index.html?id=${data.message}`;
        })
        .catch(function (error) {
            console.log('Request failed', error);
        });
}

function addcard(data) {

    let input = document.createElement("input");
    input.setAttribute("type", "radio");
    input.setAttribute("name", "card");
    input.setAttribute("value", data.ID);

    let label = document.createElement("label");
    label.innerText = `${data.Brand} ${data.Last4} Expires : ${data.Exp_month}/${data.Exp_year}`

    let linebreak = document.createElement("br");

    cardInfo.appendChild(input);
    cardInfo.appendChild(label);
    cardInfo.appendChild(linebreak);


}

function checkData(duration, terms) {
    console.log(duration, terms);

    if (duration > 0) {
        if (terms == true) {
            errorMessage.textContent = "";
            return true;
        } else {
            console.log("You must agree to out terms and condisons");
            showError("You must agree to out terms and condisons")
            return false;
        }
    } else {
        console.log("You must select a duration"); 
        showError("You must select a duration")
        return false;
    }
}

function showError(message){
    errorDiv.style.display = "block";
    errorMessage.textContent = message;
}