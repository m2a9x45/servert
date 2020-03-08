const custEmail = document.getElementById('custEmail');

const submitButton = document.getElementById('submitButton');

const URL_API = "http://127.0.0.1:8000";

submitButton.addEventListener("click", () => {
    event.preventDefault()
    console.log("clciekd");
    
    let formData = {
        "email" : custEmail.value
    }

    console.log(formData);
    

    fetch(`${URL_API}/auth/reset`, {
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
        if (data.success == true) {
            window.location = "../"
        } else {
            alert("Login didn't work")
        }
    })
    .catch(function (error) {
        console.log('Request failed', error);
    });

});