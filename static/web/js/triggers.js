window.onload = function() {
    if (!isAuthorized()) {
        console.log('user is not authorized');
        window.location.href = '/';
    }
};

document.addEventListener('DOMContentLoaded', async function() {
    document.getElementById('logoutLink').addEventListener('click', function () {
        logout();
    });

    await getTriggers();

    async function getTriggers() {
        try {
            const resp = await fetch('/api/trigger', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            const jsonBody = await resp.json();
            if (resp.status <= 299) {
                updateTable(jsonBody.records);
            } else {
                alert (jsonBody.message);
            }
        } catch (err) {
            alert (err);
        }
    }

    function updateTable(records) {
        const table = document.getElementById("data-table");
        while (table.rows.length > 1) {
            table.deleteRow(1);
        }

        records.forEach(record => {
            const tr = document.createElement("tr");

            const tdID = document.createElement("td");
            tdID.textContent = record.id;
            tr.appendChild(tdID);

            const tdName = document.createElement("td");
            tdName.textContent = record.name;
            tr.appendChild(tdName);

            const tdContainerName = document.createElement("td");
            tdContainerName.textContent = record.containerName;
            tr.appendChild(tdContainerName);

            const tdActions = document.createElement("td");

            const editBtn = document.createElement("a");  // Change this from button to anchor element.
            editBtn.textContent = "Edit";
            editBtn.className = "table-button";
            editBtn.setAttribute('data-itemId', record.id);
            editBtn.href = `/trigger/${record.id}/edit`; // Set the href attribute based on the record's id.
            tdActions.appendChild(editBtn);

            const deleteBtn = document.createElement("button");
            deleteBtn.textContent = "Delete";
            deleteBtn.className = "table-button";
            deleteBtn.setAttribute('data-itemId', record.id);
            deleteBtn.addEventListener('click', onDeleteClick);
            tdActions.appendChild(deleteBtn);

            tr.appendChild(tdActions);

            table.appendChild(tr);
        });
    }

    async function onDeleteClick (event) {
        const itemId = event.target.getAttribute('data-itemId');

        if (confirm("Are you sure you want to delete this trigger?")) {
            try {
                const resp = await fetch(`/api/trigger/${itemId}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                })
                const jsonBody = await resp.json();
                if (resp.status <= 299) {
                    await getTriggers();
                } else {
                    alert (jsonBody.message);
                }
            } catch (err) {
                alert(err);
            }
        }
    }
});