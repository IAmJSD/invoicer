"use strict";

// For context here is how POST /emails treats params:
// name, email: Delete
// name, email and username: Change username
// name, email and password: Change password
// name, email, hostname, username and password: Create e-mail

// Gets the information about a record by its index.
function getInformation(index) {
    return {
        "name": document.getElementById(index + "-name").innerText,
        "email": document.getElementById(index + "-email").innerText,
    };
}

// Creates and submits a form.
function submitForm(path, params) {
    var form = document.createElement("form");
    form.method = "POST";
    form.action = path;
    Object.keys(params).forEach(function (key) {
        var input = document.createElement("input");
        input.type = "hidden";
        input.name = key;
        input.value = params[key];
        form.appendChild(input);
    });
    document.body.appendChild(form);
    form.submit();
}

// Deletes the record.
function deleteRecord(index) {
    submitForm("/emails", getInformation(index));
}

// Used to open and load the username panel.
function openUsernamePanel(index) {
    var information = getInformation(index);
    document.getElementById("username-form-username").value = document.getElementById(index + "-username").innerText;
    document.getElementById("username-form-name").value = information.name;
    document.getElementById("username-form-email").value = information.email;
    document.getElementById("username-form-modal").classList.add("is-active");
}

// Closes the username panel.
function closeUsernamePanel() {
    document.getElementById("username-form-modal").classList.remove("is-active");
    return false;
}

// Validates the username.
function validateUsername() {
    document.getElementById("username-form-submit").disabled = document.getElementById("username-form-username").value === "";
}

// Used to open and load the password panel.
function openPasswordPanel(index) {
    var information = getInformation(index);
    document.getElementById("password-form-name").value = information.name;
    document.getElementById("password-form-email").value = information.email;
    document.getElementById("password-form-modal").classList.add("is-active");
}

// Closes the password panel.
function closePasswordPanel() {
    document.getElementById("password-form-modal").classList.remove("is-active");
    return false;
}

// Validates the form.
function validateForm() {
    var email = document.getElementById("insert-email").value;
    var hostname = document.getElementById("insert-hostname").value;
    var emailValid = email.includes("@") && email.includes(".") && email.length >= 6;
    document.getElementById("insert-button").disabled = !(emailValid && hostname !== "");
    return false;
}
