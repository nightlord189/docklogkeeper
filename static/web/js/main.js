
function isAuthorized () {
    return document.cookie.startsWith('default=') && document.cookie !== 'default='
}