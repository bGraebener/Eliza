
// function that gets executed on the button click
function askQuestion() {
    var textarea = document.getElementById("chatBoxArea");
    var userText = document.getElementById("userTextInput").value;

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

    // handler for all changes off the request ready state
    // only do something if response is available 
    xhr.onreadystatechange = function () {

        if (xhr.readyState === 4 && xhr.status === 200) {
            var userName = xhr.getResponseHeader("userName");
            if (xhr.responseText.length > 0) {

                // append new text to the textarea and keep new text in view
                textarea.value += "\n" + userName + ": " + xhr.responseText + "\nEliza: Why don't you tell me more about that?";
                textarea.scrollTop = textarea.scrollHeight;
            }
        }
    }
}

