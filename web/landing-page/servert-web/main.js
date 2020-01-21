const submitButton = document.getElementById("submitButton");
const custName = document.getElementById("custName");
const custEmail= document.getElementById("custEmail");

const URL_API = "";

submitButton.addEventListener("click", (e) => {
    event.preventDefault(e);
    
    let name = custName.value;
    let email = custEmail.value;

    if (name.trim() || email.trim() == "") {
        console.log("no values");
        // show error
    } else {
        let data = {
            name,
            email
        }
    
        console.log(data);
        sendData(data).then(data => console.log(data));
        
    }
});

async function sendData(data){
    let response = await fetch(URL_API);
    let data = await response.json();
    return data; 
}