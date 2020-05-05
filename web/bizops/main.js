// make call to get tasks.
const tasks = document.querySelectorAll(".task");
const tasksList = document.querySelector(".tasksList");
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
            // console.log(e);
            if (e.target.parentElement.id) {
                console.log(e.target.parentElement.id);  
            }
    
            if (e.target.parentElement.parentElement.id) {
                console.log(e.target.parentElement.parentElement.id);
            }
    
        });

        divTitle.appendChild(name);
        divTitle.appendChild(made);

        divTask.appendChild(divTitle);
        divTask.appendChild(status);

        tasksList.appendChild(divTask);



    });

}