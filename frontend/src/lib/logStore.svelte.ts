const log = $state({ output: "" });

function appendLog(msg: string) {
    const time = new Date().toLocaleTimeString();
    log.output += `[${time}] ${msg}\n`;
}

function clearLog() {
  log.output = "--- console cleared! ---\n";
}

export { log, appendLog, clearLog };
