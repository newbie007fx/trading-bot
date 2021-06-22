$(document).ready(function(){
    $(document.body).append('<div class="container-loader">' +
        '<div class="overlay-loader"><img src="/assets/img/ajax-loader.gif"/> Loading...</div>' +
        '</div>')
});

function showLoader() {
    $(".overlay-loader").show();
}

function hideLoader() {
    $(".overlay-loader").hide();
}