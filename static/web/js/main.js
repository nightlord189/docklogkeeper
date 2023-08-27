
function isAuthorized () {
    return document.cookie.startsWith('default=') && document.cookie !== 'default='
}

function logout () {
    console.log('logout');
    document.cookie = 'default='
    window.location.href = "/";
}