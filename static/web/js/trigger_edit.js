document.addEventListener('DOMContentLoaded', async function() {
    let mode = 'create';
    let triggerID = '';

    const currentURL = window.location.pathname; // Gets the path without the query string

    if (currentURL.endsWith('/edit')) {
        mode = 'update';
        const pathParts = window.location.pathname.split('/'); // Split the path into segments
        triggerID = pathParts[pathParts.length - 2];
        if (!triggerID || triggerID ==='') {
            throw Error('invalid trigger id');
        }
    }
    console.log(`mode ${mode}, trigger_id ${triggerID}`);

    const triggerForm = document.getElementById('triggerForm');
    const submitButton = document.getElementById('submitBtn');
    const nameInput = document.getElementById('name');
    const containerNameInput = document.getElementById('containerName');
    const containsInput = document.getElementById('contains');
    const notContainsInput = document.getElementById('notContains');
    const regexpInput = document.getElementById('regexp');
    const webhookUrlInput = document.getElementById('webhookUrl');
    const webhookHeadersInput = document.getElementById('webhookHeaders');
    const webhookBodyInput = document.getElementById('webhookBody');
    const errorText = document.getElementById('errorText');

    if (mode === 'update') {
        submitButton.textContent = 'Update trigger';
        document.getElementById('formHeader').textContent = 'Update existing trigger';

        try {
            const resp =  await fetch(`/api/trigger?trigger_id=${triggerID}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            const jsonBody = await resp.json();
            if (resp.status <= 299) {
                const currentTrigger = jsonBody.records[0]
                nameInput.value = currentTrigger.name;
                containerNameInput.value = currentTrigger.containerName;
                containsInput.value = currentTrigger['contains'];
                notContainsInput.value = currentTrigger.notContains;
                regexpInput.value = currentTrigger.regexp;
                webhookUrlInput.value = currentTrigger.webhookURL;
                webhookHeadersInput.value = currentTrigger.webhookHeaders;
                webhookBodyInput.value = currentTrigger.webhookBody;
            } else {
                alert (resp.body.message);
            }
        } catch (err) {
            alert (err);
        }
    }

    triggerForm.addEventListener('submit', async function(event) {
        console.log('submit call');

        event.preventDefault();
        submitButton.disabled = true;

        const body = {
            name: nameInput.value,
            containerName: containerNameInput.value,
            contains: containsInput.value,
            notContains: notContainsInput.value,
            regexp: regexpInput.value,
            method: 'webhook',
            webhookURL: webhookUrlInput.value,
            webhookHeaders: webhookHeadersInput.value,
            webhookBody: webhookBodyInput.value,
        };

        console.log(body);

        let url = '/api/trigger';
        let method = 'POST';

        if (mode === 'update') {
            url = `/api/trigger/${triggerID}`;
            method = 'PUT';
        }

        try {
            const resp = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(body)
            });
            const jsonBody = await resp.json();
            if (resp.status <= 299) {
                if (confirm(mode === 'create' ? 'Trigger created' : 'Trigger updated')) {
                    window.location.href = "/triggers";
                }
            }
            else {
                alert (jsonBody.message);
                submitButton.disabled = false;
            }
        } catch (err) {
            alert (err);
            submitButton.disabled = false;
        }
    });
});