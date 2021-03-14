"use strict";

// Closes the edit panel.
function closeEditPanel() {
    document.getElementById("user-edit-modal").classList.remove("is-active");
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

// Deletes a invoice job.
function deleteJob(jobId) {
    submitForm("/", {"job_id": jobId});
}

// Validates the creation form.
function validateCreationForm() {
    var email = document.getElementById("insert-email").value;
    var name = document.getElementById("insert-name").value;
    var markdown = document.getElementById("insert-markdown").value;
    var emailValid = email.includes("@") && email.includes(".") && email.length >= 6;
    document.getElementById("insert-button").disabled = !(emailValid && name !== "" && markdown !== "");
}

// Validate the markdown is not blank.
function validateNotBlank() {
    document.getElementById("edit-job-submit").disabled = document.getElementById("edit-job-markdown").value === "";
}

// Update the e-mail contents.
function editJob(jobId) {
    // Get the applicable text for the e-mails.
    var emailContent = document.getElementById(jobId + "-content").innerText;
    var subject = document.getElementById(jobId + "-subject").innerText;
    var markdown = document.getElementById(jobId + "-markdown").innerText;

    // Set the edit job ID.
    document.getElementById("edit-job-id").value = jobId;

    // Set the e-mail content.
    document.getElementById("edit-job-email").value = emailContent;

    // Set the subject.
    document.getElementById("edit-job-subject").value = subject;

    // Set the markdown.
    document.getElementById("edit-job-markdown").value = markdown;

    // Display the form.
    document.getElementById("user-edit-modal").classList.add("is-active");
}
