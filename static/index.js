const terminal = document.getElementById('console');
const body = document.getElementsByTagName('body')[0];
const commandLine = document.getElementById('command-line');
const saveConfigurationBtn = document.getElementById('save-configuration');

const websocket = new WebSocket(`ws://localhost:8080/ws-serial-debugger`);

websocket.onmessage = (event) => {
	terminal.value += `${event.data}`;
	terminal.scrollTop = terminal.scrollHeight;
};

body.addEventListener('keypress', function (event) {
	if (event.key === 'Enter') {
		if (commandLine.value.length > 0) {
			terminal.value += `${commandLine.value}\n`;
			websocket.send(JSON.stringify({ type: 'command', message: commandLine.value }));
			commandLine.value = '';
		}
	}
});
