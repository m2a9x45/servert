const custName = document.getElementById('custName');
const custEmail = document.getElementById('custEmail');
const custPass = document.getElementById('custPass');
const cust18 = document.getElementById('cust18');

const submitButton = document.getElementById('submitButton');

const URL_API = "http://127.0.0.1:8000";

submitButton.addEventListener("click", () => {
    event.preventDefault()
    console.log("clciekd");
    
    if (cust18.checked == true && custName.value != "" && custEmail.value != "" && custPass.value != "") {
        console.log("data ready");  
    } else {
        console.log("data not ready");  
        // show error
    }

    let formData = {
        "name" : custName.value,
        "email" : custEmail.value,
        "password" : custPass.value
    }

    fetch(`${URL_API}/signup`, {
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

});