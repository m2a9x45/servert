const custName = document.getElementById('custName');
const custEmail = document.getElementById('custEmail');
const custPass = document.getElementById('custPass');
const cust18 = document.getElementById('cust18');

const submitButton = document.getElementById('submitButton');

const URL_API = "http://127.0.0.1:8000";

submitButton.addEventListener("click", () => {
    event.preventDefault()
    if (custName.value != "") {
        if (custEmail.value != "" && validateEmail(custEmail.value)) {
            if (custPass.value != "") {
                if (cust18.checked == true) {
                    let formData = {
                        "name" : custName.value,
                        "email" : custEmail.value,
                        "password" : custPass.value
                    }
                    console.log(formData);

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

                } else {
                    console.log("You must be a least 18 to use our services");  
                }
            } else {
                console.log("You must choose a password"); 
            }
        } else {
            console.log("You must enter your email");  
        }
    } else {
        console.log("You must enter your name");
    }
});

function showError(err){

}

function validateEmail(email) 
{
    var re = /\S+@\S+\.\S+/;
    return re.test(email);
}