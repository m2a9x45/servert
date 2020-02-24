
const productsList = document.querySelector('.products');
const LoginButton = document.querySelector('#Login');
const SignupButton = document.querySelector('#Signup');
const URL_API = "http://127.0.0.1:8000";

refershToken();

fetch(`${URL_API}/products`)
    .then(res => res.json())
    .then(data => {
        console.log(data)
        addProducts(data);
    })
    .catch(err => console.log(err))


function addProducts(products){
    products.forEach(product => {
        console.log(product);
        addProductToPage(product);
    });
}

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
        LoginButton.href = "./account"
        SignupButton.style.display = "none"
    } else {
        LoginButton.innerText = "Login"
        LoginButton.href = "./signin/"
        SignupButton.style.display = "block"
    }
})
.catch(function (error) {
    console.log('Request failed', error);
});


function addProductToPage(product){

    let div = document.createElement('div');
    div.setAttribute("class", "product");

    let h5 = document.createElement('h5');
    h5.setAttribute("id", "productName");
    h5.innerText = product.name;

    let productCore = document.createElement('h5');
    productCore.setAttribute("id", "productCore");
    productCore.innerText = product.cpu + " core";

    let productRam = document.createElement('h5');
    productRam.setAttribute("id", "productRam");
    productRam.innerText = product.ram + " GB";

    let productDrive = document.createElement('h5');
    productDrive.setAttribute("id", "productDrive");
    productDrive.innerText = product.disk + " GB"

    let productprice = document.createElement('h5');
    productprice.setAttribute("id", "productprice");
    var formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'GBP',
      });
    let price = formatter.format(product.price);
    productprice.innerText = price;

    let addProduct = document.createElement('button');
    addProduct.setAttribute("id", "addProduct");
    addProduct.setAttribute("value", product.id)
    addProduct.addEventListener("click", (e) => {
        console.log(e);
        console.log(e.target.value);

        let url = "./details?id=" + product.id;
        console.log(url);
        

        window.location.href = "./details/index.html?id=" + product.uuid;
        
    })
    addProduct.innerText = "Get";


    div.appendChild(h5);
    div.appendChild(productCore);
    div.appendChild(productRam);
    div.appendChild(productDrive);
    div.appendChild(productprice);
    div.appendChild(addProduct);

    productsList.appendChild(div);

}

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
    });
}