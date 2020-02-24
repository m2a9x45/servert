const submitButton = document.getElementById("submitButton");
const custName = document.getElementById("custName");
const custEmail = document.getElementById("custEmail");
const popup = document.getElementById("id01");
const heyName = document.getElementById("heyName");
const closeButton = document.getElementById("closeButton");
const errorMessage = document.querySelector("#errorMessage");

const URL_API = "http://localhost:8000";

closeButton.addEventListener("click", () => {
    popup.style.display = "none";
});

submitButton.addEventListener("click", () => {
    event.preventDefault();

    const name = custName.value;
    const email = custEmail.value;

    console.log(name, email);

    if (name != "") {
        if (custEmail.value != "" && validateEmail(custEmail.value)) {
            let formData = {
                name,
                email
            }
            console.log(formData);

            fetch(`${URL_API}/account/intrest`, {
                    method: 'post',
                    credentials: 'include',
                    headers: {
                        "Content-type": "application/json"
                    },
                    body: JSON.stringify(formData)
                })
                .then(response => response.json())
                .then((data) => {
                    console.log(data)
                    if (data.success != true) {
                        showError(data.message)
                    } else {
                        errorMessage.innerText = "";
                        heyName.innerText = `Hey ${name}`
                        popup.style.display = "block";
                    }
                })
                .catch(function (error) {
                    console.log('Request failed', error);
                });

        } else {
            console.log("You must enter your email"); 
            showError("You must enter your email");
        }
    } else {
        console.log("You must enter your name");
        showError("You must enter your name");
    }



});

function validateEmail(email) {
    var re = /\S+@\S+\.\S+/;
    return re.test(email);
}

function showError(err){
    errorMessage.style.color = "red";
    errorMessage.innerText = err;
}