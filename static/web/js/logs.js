const limit = 1000;

let currentContainer = '';
let currentContainers = [];

let lastLogID = 0; // for next
let firstLogID = 0; // for past

let updateLogsInterval;

async function setCurrentContainer (newVal) {
    if (newVal === currentContainer) {
        return
    }
    currentContainer = newVal;
    lastLogID = 0;
    firstLogID = 0;
    renderLogs([], appendValues.NONE, 'Loading...');

    const resp = await getLogs('future', 0);

    firstLogID = resp.firstCursor;
    lastLogID = resp.lastCursor;
    renderLogs(resp.records, appendValues.NONE, 'There is no logs');

    startLogsUpdate();
}

window.onload = function() {
    if (!isAuthorized()) {
        console.log('user is not authorized');
        window.location.href = '/';
    }
};

function renderContainers(containers) {
    if (arraysOfObjectsAreEqual(currentContainers, containers)) {
        console.log(`containers didn't change`)
        return
    }

    console.log(`rendering containers`);

    currentContainers = containers;

    const contRoot = document.getElementById('containersRoot');

    // 1. Remove all children of containersRoot
    while (contRoot.firstChild) {
        contRoot.removeChild(contRoot.firstChild);
    }

    // 2. Iterate over each container item and append to containersRoot
    containers.forEach(item => {
        const containerLink = document.createElement('a');
        containerLink.href = "#"; // You can modify the href as required
        containerLink.setAttribute('data-itemId', item.shortName);
        containerLink.className = 'text-dark container';

        // Create the status bulb and assign the appropriate class
        const statusBulb = document.createElement('span');
        statusBulb.className = 'bulb ' + (item.isAlive ? 'alive' : 'dead');

        // Append the status bulb and container name
        containerLink.appendChild(statusBulb);
        containerLink.appendChild(document.createTextNode(item.shortName));

        // Click event to update the selected container
        containerLink.addEventListener('click', async function(event) {
            event.preventDefault();
            // Remove the selected-container class from any other container
            const selected = document.querySelector('a.container.selected');
            if (selected) {
                selected.classList.remove('selected');
            }

            // Add the selected-container class to the clicked container
            containerLink.classList.add('selected');

            await setCurrentContainer(event.target.getAttribute('data-itemId'));
        });

        contRoot.appendChild(containerLink);
    });
}

async function updateContainers () {
    document.getElementById('updateContainers').disabled = true;
    try {
        const resp = await fetch('/api/container', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });
        const respJson = await resp.json();

        if (resp.status <= 299) {
            console.log('success get containers');
            renderContainers(respJson.containers);
        } else {
            const errorMessage = respJson.message || 'Something went wrong';
            console.log(errorMessage);
        }
    } catch (err) {
        console.log(err);
    }
    document.getElementById('updateContainers').disabled = false;
}

const appendValues = { NONE: 'none', UP: 'up', DOWN: 'down' };

function renderLogs(logs, append, messageForEmpty, highlightTerm = '') {
    const logsRoot = document.getElementById('logsRoot');

    console.log(`rendering logs, length ${logs.length}, append ${append}`);

    // Store current scroll position
    const initialScrollTop = logsRoot.scrollTop;

    // Check if the current scroll position is at the top
    const isAtTop = initialScrollTop === 0;

    if (highlightTerm) {
        const regex = new RegExp(highlightTerm, 'gi');
        logs = logs.map(log => log.replace(regex, match => `<span class="highlight">${match}</span>`));
    }

    switch (append) {
        case appendValues.NONE:
            if (!logs || logs.length === 0) {
                logsRoot.innerHTML = messageForEmpty;
            } else {
                logsRoot.innerHTML = logs.join('<br>');
            }
            break;
        case appendValues.DOWN:
            if (logs.length > 0) {
                logsRoot.innerHTML += '<br>' + logs.join('<br>');
            }
            break;
        case appendValues.UP:
            if (logs.length > 0) {
                // Create a dummy div to measure the height of the new content
                const dummyDiv = document.createElement('div');
                dummyDiv.innerHTML = logs.join('<br>') + '<br>';
                document.body.appendChild(dummyDiv);
                const addedHeight = dummyDiv.getBoundingClientRect().height;
                document.body.removeChild(dummyDiv);

                logsRoot.innerHTML = logs.join('<br>') + '<br>' + logsRoot.innerHTML;

                if (!isAtTop) {
                    // Adjust the scroll position by the height of the newly added content
                    logsRoot.scrollTop = initialScrollTop + addedHeight;
                }
                // If scroll position was at the top, it will remain at the top naturally
            }
            break;
        default:
            console.error(`unknown append: ${append}`);
    }

    if (append !== appendValues.UP) {
        // Restore the original scroll position
        logsRoot.scrollTop = initialScrollTop;
    }
}


async function getLogs (dir, cursor) {
    try {
        const resp = await fetch(`/api/container/${currentContainer}/log?direction=${dir}&cursor=${cursor}&limit=${limit}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });
        const respJson = await resp.json();

        if (resp.status <= 299) {
            console.log(`got logs for container ${currentContainer} cursor ${cursor}`);

            return respJson;
        } else {
            const errorMessage = respJson.message || 'Something went wrong';
            console.log(errorMessage);
            return {logs:[], chunkNumber:0, offset:0}
        }
    } catch (err) {
        console.log(err);
    }
    return {logs:[], chunkNumber:0, offset:0};
}

async function search (contains) {
    try {
        const resp = await fetch(`/api/container/${currentContainer}/log/search?contains=${contains}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });
        const respJson = await resp.json();

        if (resp.status <= 299) {
            console.log('success search');
            renderLogs(respJson.records, appendValues.NONE, 'Not found', contains);
        } else {
            const errorMessage = respJson.message || 'Something went wrong';
            console.log(errorMessage);
        }
    } catch (err) {
        console.log(err);
    }
}

async function updateLogs () {
    if (currentContainer === '') {
        return
    }

    const updLogsBtn = document.getElementById('updateLogsButton');
    updLogsBtn.disabled = true;

    const resp = await getLogs('future', lastLogID);
    if (resp.records.length === 0) {
        updLogsBtn.disabled = false;
        return
    }

    lastLogID = resp.lastCursor;

    renderLogs(resp.records, appendValues.UP, '');

    updLogsBtn.disabled = false;
}

async function nextLogs () {
    const nextLogsButton = document.getElementById('nextLogsButton');
    console.log('next logs');

    const resp = await getLogs('past', firstLogID);
    if (resp.records.length === 0) {
        nextLogsButton.disabled = false;
        return
    }

    firstLogID = resp.firstCursor;

    renderLogs(resp.records, appendValues.DOWN, '');

    nextLogsButton.disabled = false;
}

function startLogsUpdate () {
    console.log('startLogsUpdate');
    if (!updateLogsInterval) {
        updateLogsInterval = setInterval(updateLogs, 5000);
    }
    document.getElementById('autoRefreshCheckbox').checked = updateLogsInterval !== null;
    console.log(`checkbox checked: ${updateLogsInterval !== null}`);
}

function stopLogsUpdate () {
    console.log('stopLogsUpdate');
    if (updateLogsInterval) {
        clearInterval(updateLogsInterval);
        updateLogsInterval = null;
    }
    document.getElementById('autoRefreshCheckbox').checked = updateLogsInterval !== null;
    console.log(`checkbox checked: ${updateLogsInterval !== null}`);
}

document.addEventListener('DOMContentLoaded', async function() {
    document.getElementById('logoutLink').addEventListener('click', function () {
        logout();
    });

    const searchInput = document.getElementById('searchInput');
    const findButton = document.getElementById('findButton');
    const searchForm = document.getElementById('searchForm');

    const updContainersBtn = document.getElementById('updateContainers');
    const updLogsBtn = document.getElementById('updateLogsButton');
    const nextLogsButton = document.getElementById('nextLogsButton');

    async function searchListener (event) {
        event.preventDefault();
        console.log('search');

        if (!currentContainer) {
            return
        }

        findButton.disabled = true;

        const searchText = searchInput.value;
        if (!searchText) {
            findButton.disabled = false;
            return
        }

        stopLogsUpdate();
        await search(searchText);

        findButton.disabled = false;
    }

    searchForm.addEventListener('submit', searchListener);
    findButton.addEventListener('click', searchListener);

    updContainersBtn.addEventListener('click', updateContainers);

    updLogsBtn.addEventListener('click', updateLogs);

    nextLogsButton.addEventListener('click', nextLogs)

    await updateContainers();

    setInterval(updateContainers, 10000);

    startLogsUpdate();

    document.getElementById('autoRefreshCheckbox').addEventListener('change', function(event) {
        console.log('change checkmark');

        if (event.target.checked) {
            startLogsUpdate();
        } else {
            stopLogsUpdate();
        }
    });
});