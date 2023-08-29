
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

function isAuthorized () {
    return document.cookie.startsWith('default=') && document.cookie !== 'default='
}

function logout () {
    console.log('logout');
    document.cookie = 'default='
    window.location.href = "/";
}