
function arraysOfObjectsAreEqual(arr1, arr2) {
    if (arr1.length !== arr2.length) {
        return false;
    }

    for (let i = 0; i < arr1.length; i++) {
        const obj1 = arr1[i];
        const obj2 = arr2[i];

        for (const key in obj1) {
            if (Object.prototype.hasOwnProperty.call(obj1, key)) {
                if (typeof obj1[key] === 'object' && obj1[key] !== null && obj2[key] !== null) {
                    // Recursive check for nested objects
                    if (!arraysOfObjectsAreEqual([obj1[key]], [obj2[key]])) {
                        return false;
                    }
                } else if (obj1[key] !== obj2[key]) {
                    return false;
                }
            }
        }
    }

    return true;
}

function getCookie(cookieKey) {
    // Split the cookies string into an array of individual cookies
    const cookiesArray = document.cookie.split('; ');

    // Iterate over each cookie to find the one with the specified key
    for (const cookie of cookiesArray) {
        const [key, value] = cookie.split('=');

        // Trim any leading or trailing spaces
        const trimmedKey = key.trim();

        // Check if the current cookie's key matches the specified key
        if (trimmedKey === cookieKey) {
            // Return the corresponding value
            return decodeURIComponent(value);
        }
    }

    // Return null if the cookie with the specified key is not found
    return null;
}

function isAuthorized () {
    //console.log('raw cookie: ', document.cookie)
    const defaultCookie = getCookie('default')
    //console.log('defaultCookie: ', defaultCookie)
    return defaultCookie !== null && defaultCookie.trim() !== ''
}

function logout () {
    console.log('logout');
    document.cookie = 'default='
    window.location.href = "/";
}