const submitButton = document.getElementById("submitButton");
const custName = document.getElementById("custName");
const custEmail= document.getElementById("custEmail");

const URL_API = "http://localhost:8000";

submitButton.addEventListener("click", () => {
    event.preventDefault();
    
    const name = custName.value;
    const email = custEmail.value;

    console.log(name, email);
    
    let formData = {
        name,
        email
    }

    console.log(formData);
    
    fetch(`${URL_API}/intrest`, {
        method: 'post',
        credentials: 'include',
        headers: {
            "Content-type": "application/json"
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then((data) => console.log(data))
    .catch(function (error) {
        console.log('Request failed', error);
    });
    
        
    
});

