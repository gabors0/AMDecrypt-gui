const log = $state({ output: "", hasError: false });
let errorTimeout: ReturnType<typeof setTimeout> | null = null;

function appendLog(msg: string) {
  const time = new Date().toLocaleTimeString();
  log.output += `[${time}] ${msg}\n`;
  if (msg.includes("[ERROR]")) {
    log.hasError = true;
    clearTimeout(errorTimeout);
    errorTimeout = setTimeout(() => {
      log.hasError = false;
    }, 1000);
  }
}

function clearLog() {
  log.output = "--- Console cleared! ---\n";
}

export { log, appendLog, clearLog };
