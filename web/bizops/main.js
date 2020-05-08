// make call to get tasks.
const tasks = document.querySelectorAll(".task");
const tasksList = document.querySelector(".tasksList");
const content = document.querySelector(".content");
const serverTask = document.querySelector(".serverTask")
const loginFrom = document.querySelector(".loginFrom");
const username = document.querySelector("#username");
const password = document.querySelector("#password");
const contaniner = document.querySelector(".contaniner");
const URL_API = "http://127.0.0.1:8000";

window.addEventListener("load", () => {
    fetch(`${URL_API}/auth/isstaffloggedin`, {
        method: 'get',
        credentials: 'include',
        headers: {
            "Content-type": "application/json",
        }
    })
    .then(response => {
        if (response.status == 401) {
            console.log(response.status);
            return {"success" : false};
        } else {
            return response.json();
        }    
    })
    .then (data => {
        console.log(data);

        if (data.success == true) {
            loginFrom.style.display = "none";
            contaniner.style.display = "grid";
            getTasks();
            refershToken();
        }

    })
    .catch(function (error) {
        console.log('Request failed', error);
    });
});

loginFrom.addEventListener("submit", (e) => {
    e.preventDefault();

    let formData = {
        "email" : username.value,
        "password" : password.value,
    }

    fetch(`${URL_API}/auth/stafflogin`, {
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
            loginFrom.style.display = "none";
            contaniner.style.display = "grid";
            getTasks();
            refershToken();
        }

    })
    .catch(function (error) {
        console.log('Request failed', error);
    });
    
})

function getTasks() {
    fetch(`${URL_API}/internal/tasks`, {
        method: 'get',
        credentials: 'include',
        headers: {
            "Content-type": "application/json",
        }
    })
    .then(res => res.json())
    .then(data => {
        console.log(data);
        addTaskToPage(data);
    })
    .catch(err => console.log(err))
}

function addTaskToPage(tasks) {

    tasks.forEach(task => {
        
        let divTask = document.createElement("div");
        divTask.setAttribute("class", "task");
        divTask.setAttribute("id", task.LinkID);

        let divTitle = document.createElement("div");
        divTitle.setAttribute("class", "taskTitle");

        let name = document.createElement("p");
        name.innerText = task.ID;

        let made = document.createElement("p");

        let now = new Date().getTime();
        let created = new Date(task.CreatedAt * 1000).getTime();

        // console.log(now);
        // console.log(created);

        let diff = now - created;
        let secs = diff / 1000;
        let mins = secs / 60;
        let hours = mins / 60;
        let days = hours / 24;

        if (secs <= 60) {
            // display in seconds as it's less than an minits
            // console.log("seconds :", Math.round(secs));
            made.innerText = `${Math.round(secs)} Seconds ago`
        } else if (mins <= 60) {
            // display in minitues as it's less than an hour
            //console.log("minitues :", Math.round(mins));      
            made.innerText = `${Math.round(mins)} minutes ago`
        } else if (hours <= 24) {
            // display in hours as it's less than an days
            //console.log("hours :", Math.round(hours));     
            made.innerText = `${Math.round(hours)} hours ago`
        } else {
            //console.log("days :", Math.round(days));
            made.innerText = `${Math.round(days)} days ago`
            
        }

        let status = document.createElement("p");
        status.innerText = task.Status;

        divTask.addEventListener("click", (e) => {
            serverTask.innerHTML = "";
            selectTask(e, task.UUID, task.LinkID, task.UserID);
        });

        divTitle.appendChild(name);
        divTitle.appendChild(made);

        divTask.appendChild(divTitle);
        divTask.appendChild(status);

        tasksList.appendChild(divTask);



    });

}

window.onpopstate = function(event) {
    alert("location: " + document.location + ", state: " + JSON.stringify(event.state));
  };

function selectTask(e, taskID, LinkID, User_ID) {
    // console.log(taskID, LinkID, User_ID);
    let url = window.history.pushState( {} , '', `?task_id=${taskID}` );
    content.style.visibility = "visible";

    let OrderID = document.createElement("p");
    OrderID.innerText = LinkID;

    let UserID = document.createElement("p");
    UserID.innerText = User_ID;

    let closeTask = document.createElement("button");
    closeTask.innerText = "Close Task";
    closeTask.addEventListener("click", () => {
        console.log(taskID, "closed");
        
    });

    serverTask.appendChild(OrderID);
    serverTask.appendChild(UserID);
    serverTask.appendChild(closeTask);

}

setInterval(function() {
    refershToken();
}, 60000); // Every 1 minitue


function refershToken() {
    fetch(`${URL_API}/internal/refresh`, {
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