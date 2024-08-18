window.onload = function() {
    var message = document.getElementById('message');
    if (message) {
        message.classList.add('show');
        setTimeout(function() {
            message.classList.remove('show');
        }, 5000); // 5 секунд
    }
};