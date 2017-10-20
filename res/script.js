
// function that gets executed on the button click
function askQuestion() {
    var textarea = document.getElementById("chatBoxArea");

    // create a new AJAX object
    var xhr = new XMLHttpRequest();

    // send a post request to the path /question
    xhr.open("POST", "/question");

    // for debugging, otherwise the browser doesn't resend requests
    xhr.setRequestHeader("Cache-Control", "no-cache");

    // set a custom header value to the value of the text input field
    xhr.setRequestHeader("user-question", document.getElementById("userTextInput").value);
    
    // send the request
    xhr.send();

    // handler for all changes off the request ready state
    // only do something if response is available 
    xhr.onreadystatechange = function () {

        if (xhr.readyState === 4 && xhr.status === 200) {
            // append new text to the textarea and keep new text in view
            textarea.value += "\nUser: " + xhr.responseText + "\nEliza: Why don't you tell me more about that?\n";
            textarea.scrollTop = textarea.scrollHeight;
        }
    }
}
