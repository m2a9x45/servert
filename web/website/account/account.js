const URL_API = "http://127.0.0.1:8000";

refershToken();

fetch(`${URL_API}/account`, {
    method: 'get',
    credentials: 'include',
    headers: {
        "Content-type": "application/json",
    }
})
.then(response => response.json())
.then (data => {
    console.log(data);
    let userID = document.createElement("p")
    userID.innerText = "Hey " + data.message;

    document.body.append(userID);
})
.catch(function (error) {
    console.log('Request failed', error);
    console.log(response.status); 
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