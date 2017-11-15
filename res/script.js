
$().ready(function () {
    // handle click event on submit button
    $("#questionButton").click(function () {
        askQuestion();

    });

    // if user input text field is in focus allow enter key to trigger the ask question function
    $("#userTextInput").keypress((event) => {
        if (event.which === 13) {
            askQuestion();
        }
    });

});

function createListItem(name, msg) {

    var chatbox = document.getElementById("chatbox");

    var newList = document.createElement("li");
    newList.className = "w3-bar w3-round-large";
    newList.style.border = "none";

    var img = document.createElement("img");
    img.style.width = "85px";
    img.className = "w3-bar-item w3-circle";
    img.src = "res/female_avatar.png";

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
        newList.style.float = "right"
    }

    chatbox.appendChild(newList);
    window.scrollTo(0, document.body.scrollHeight);
}

// function that gets executed on the button click
function askQuestion() {
    var userText = document.getElementById("userTextInput").value;
    var questionButton = document.getElementById("questionButton");

    // empty the user input text field
    document.getElementById("userTextInput").value = "";

    // don't send request if user didn't enter text
    if (userText.length < 1) {
        return;
    }

    // send a post request to the path /question
    $.ajax("/question", {

        // for debugging, otherwise the browser doesn't resend requests
        // set a custom header value to the value of the text input field
        headers: { "Cache-Control": "no-cache", "user-question": userText },
        method: "POST"
    }).done(function (data, textStatus, jqXHR) {
        var userName = jqXHR.getResponseHeader("userName");
        if (data.length > 0) {

            // append new text to the window and keep new text in view
            createListItem(userName, userText);

            var rand = Math.floor(Math.random() * 2000);
            setTimeout(() => { createListItem("Eliza", data); }, rand)

            // disable the submit button if the user quit the program
            questionButton.disabled = jqXHR.getResponseHeader("quit") === "true";
            document.getElementById("userTextInput").disabled = jqXHR.getResponseHeader("quit") === "true";
            document.getElementById("userTextInput").setAttribute("placeholder", "Please start a new session!");

        }
    })

}