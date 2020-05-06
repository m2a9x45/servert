// make call to get tasks.
const tasks = document.querySelectorAll(".task");
const tasksList = document.querySelector(".tasksList");
const content = document.querySelector(".content");
const URL_API = "http://127.0.0.1:8000";

fetch(`${URL_API}/internal/tasks`)
    .then(res => res.json())
    .then(data => {
        console.log(data);
        addTaskToPage(data);
    })
    .catch(err => console.log(err))

tasks.forEach(task => {
    task.addEventListener("click", (e) => {
        // console.log(e);
        if (e.target.parentElement.id) {
            console.log(e.target.parentElement.id);  
        }

        if (e.target.parentElement.parentElement.id) {
            console.log(e.target.parentElement.parentElement.id);
        }

    })
    
});

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

        const d = new Date(task.CreatedAt * 1000);
        const strDate = d.toLocaleString('en-GB', {dateStyle: "short", timeStyle : "short"});

        made.innerText = strDate;

        let status = document.createElement("p");
        status.innerText = task.Status;

        divTask.addEventListener("click", (e) => {
            selectTask(e, task.UUID);
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

function selectTask(e, taskID) {
    console.log(taskID);
    console.log(Date.now());
    

    let url = window.history.pushState( {} , '', `?task_id=${taskID}` );

    content.style.visibility = "visible";
    

}