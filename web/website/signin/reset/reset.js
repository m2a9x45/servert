const custPass = document.getElementById('custPass');
const custPass1 = document.getElementById('custPass1');

const submitButton = document.getElementById('submitButton');

const URL_API = "http://127.0.0.1:8000";

submitButton.addEventListener("click", () => {
    event.preventDefault()
    console.log("clicked");

    let url_string = window.location.href;
    let url = new URL(url_string);
    let c = url.searchParams.get("t");
    console.log(c);


    if (custPass.value == custPass1.value && (custPass1.value || custPass.value != "")) {
        let formData = {
            "password" : custPass.value,
            "token" : c
        }

        console.log(formData);

        fetch(`${URL_API}/auth/restpassword`, {
            method: 'PATCH',
            credentials: 'include',
            headers: {
                "Content-type": "application/json",
            },
            body: JSON.stringify(formData)
        })
        .then(response => response.json())
        .then((data) => {
            console.log(data);
        })
        .catch(function (error) {
            console.log('Request failed', error);
        });
    } else {
        console.log("Something went wrong");
    }

    


    


});