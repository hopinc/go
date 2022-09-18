'use strict';

// This is some incredibly simple JS that effectively just acts as glue.

(function () {
    // Make the websocket connection.
    const ws = new WebSocket(`ws://${window.location.host}/ws`);
    ws.binaryType = 'arraybuffer';

    // Add the events.
    ws.onerror = () => {
        document.body.innerHTML = '<h1>WebSocket error. Please refresh.</h1>';
    };
    ws.onmessage = ev => {
        const view = new DataView(ev.data);
        const type = view.getInt8(0);
        switch (type) {
            case 0: {
                // Ping.
                ws.send(new Uint8Array([0]));
                break;
            }
            case 1: {
                // JSON data.
                const json = new TextDecoder().decode(ev.data.slice(1));
                document.getElementById('socket-data').value += json + '\r\n';
            }
        }
    };

    // Handle form submits.
    document.getElementById('channel-form').addEventListener('submit', async ev => {
        ev.preventDefault();

        const params = new URLSearchParams();
        params.append('channel', document.getElementById('channel-id').value);

        await fetch(`/subscribe?${params.toString()}`, {method: 'POST'}).then(res => res.text()).then(body => {
            document.getElementById('join-result').innerText = body;
        });
    });
})();
