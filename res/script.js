
// function that gets executed on the button click
function askQuestion() {
    var textarea = document.getElementById("chatBoxArea");
    var userText = document.getElementById("userTextInput").value;
    var questionButton = document.getElementById("questionButton");

    // empty the user input text field
    document.getElementById("userTextInput").value = "";

    // don't send request if user didn't enter text
    if (userText.length < 1) {
        return;
    }

    // create a new AJAX object
    var xhr = new XMLHttpRequest();

    // send a post request to the path /question
    xhr.open("POST", "/question");

    // for debugging, otherwise the browser doesn't resend requests
    xhr.setRequestHeader("Cache-Control", "no-cache");

    // set a custom header value to the value of the text input field
    xhr.setRequestHeader("user-question", userText);

    // send the request
    xhr.send();

    // handler for all changes of the request ready state
    // only do something if response is available 
    xhr.onreadystatechange = () => {

        if (xhr.readyState === 4 && xhr.status === 200) {
            // retrieve the user name
            var userName = xhr.getResponseHeader("userName");
            if (xhr.responseText.length > 0) {

                // append new text to the window and keep new text in view
                createListItem(userName, userText);

                var rand = Math.floor(Math.random() * 2000);
                setTimeout(() => { createListItem("Eliza", xhr.responseText); }, rand)

                // disable the submit button if the user quit the program
                questionButton.disabled = xhr.getResponseHeader("quit") === "true";
                document.getElementById("userTextInput").disabled = xhr.getResponseHeader("quit") === "true";

            }
        }
    }
}

function createListItem(name, msg) {

    var chatbox = document.getElementById("chatbox");

    var newList = document.createElement("li");
    newList.className = "w3-bar w3-round-large";
    newList.style.border = "none";

    var img = document.createElement("img");
    img.style.width = "85px";
    img.className = "w3-bar-item w3-circle";
    img.src = "female_avatar.png";

    var div = document.createElement("div");
    div.className = "w3-bar-item";

    var spanOne = document.createElement("span");
    spanOne.className = "w3-large";
    spanOne.innerText = name;
    var spanTwo = document.createElement("span");
    spanTwo.innerText = msg;

    var br = document.createElement("br");

    div.appendChild(spanOne);
    div.appendChild(br);
    div.appendChild(spanTwo);

    newList.appendChild(img);
    newList.appendChild(div);

    if (name != "Eliza") {
        img.src = "male_avatar.png";
        div.style = "float:right";
        img.style = "float:right; width:85px";
        newList.style.textAlign = "right";
        newList.style.backgroundColor = "rgba(31, 149, 208,.6)";
        newList.className += " w3-text-white"
    }

    chatbox.appendChild(newList);
    window.scrollTo(0, document.body.scrollHeight);
}
// if user input text field is in focus allow enter key to trigger the ask question function
document.getElementById("userTextInput").addEventListener("keypress", (event) => {
    if (event.which === 13) {
        askQuestion();
    }
});

