document.addEventListener("htmx:afterRequest", function(evt) {
    var xhr = evt.detail.xhr;
    var response = JSON.parse(xhr.responseText);
    if (response.redirect) {
        window.location.href = response.redirect;
    }
});