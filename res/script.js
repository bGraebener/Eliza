
$().ready(() => {
    // handle click event on submit button
    $("#questionButton").click(() => {
        askQuestion();
    });

    // if user input text field is in focus allow enter key to trigger the ask question function
    $("#userTextInput").keypress((event) => {
        if (event.which === 13) {
            askQuestion();
        }
    });
});

// function that creates a list item with the passed name and message text
function createListItem(name, msg) {

    var chatbox = document.getElementById("chatbox");

    // new list element
    var newList = document.createElement("li");
    newList.className = "w3-bar w3-round-large";
    newList.style.border = "none";

    // default is a list elemtent for an eliza response, left justified, female icon
    var img = document.createElement("img");
    img.style.width = "85px";
    img.className = "w3-bar-item w3-circle";
    img.src = "res/female_avatar.png";

    var div = document.createElement("div");
    div.className = "w3-bar-item";

    // name and message
    var spanOne = document.createElement("span");
    spanOne.className = "w3-large";
    spanOne.innerText = name;
    var spanTwo = document.createElement("span");
    spanTwo.innerText = msg;

    var br = document.createElement("br");

    // add everything to the list
    div.appendChild(spanOne);
    div.appendChild(br);
    div.appendChild(spanTwo);

    newList.appendChild(img);
    newList.appendChild(div);

    // give the user a different icon and put the list item to the right side of the screen
    if (name != "Eliza") {
        img.src = "res/male_avatar.png";
        div.style = "float:right";
        img.style = "float:right; width:85px";
        newList.style.textAlign = "right";
        newList.style.backgroundColor = "rgba(31, 149, 208,.6)";
        newList.className += " w3-text-white"
        newList.style.float = "right"
    }

    // add it to the list
    chatbox.appendChild(newList);
    // keep the new item in view
    window.scrollTo(0, document.body.scrollHeight);
}

// function that gets executed on the button click
function askQuestion() {
    var userText = $("#userTextInput").value;

    // empty the user input text field
    $("#userTextInput").value = "";

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

        // function that executes after a successfull ajax request
    }).done((data, textStatus, jqXHR) => {
        var userName = jqXHR.getResponseHeader("userName");
        if (data.length > 0) {

            // append user question to the window
            createListItem(userName, userText);

            // wait a random amount of time to simulate eliza thinking about a response and add it to the window
            var rand = Math.floor(Math.random() * 2000);
            setTimeout(() => createListItem("Eliza", data), rand)
        }

        $("#userTextInput").focus();

        // disable the submit button and input field if the user quit the program
        $("#questionButton").disabled = jqXHR.getResponseHeader("quit") === "true";
        $("#userTextInput").disabled = jqXHR.getResponseHeader("quit") === "true";

        if ($("#userTextInput").disabled) {
            $("#userTextInput").setAttribute("placeholder", "Please start a new session!");
        }
    })
}