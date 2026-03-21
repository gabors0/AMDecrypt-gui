type AmdLine = {
  text: string;
  source: "stdout" | "stderr";
  timestamp: string;
};

const MAX_LINES = 10000;

const amd = $state({
  running: false,
  lines: [] as AmdLine[],
});

function addLine(text: string, source: "stdout" | "stderr") {
  amd.lines.push({ text, source, timestamp: new Date().toLocaleTimeString() });
  if (amd.lines.length > MAX_LINES) {
    amd.lines.splice(0, amd.lines.length - MAX_LINES);
  }
}

function clearLines() {
  amd.lines = [];
}

function setRunning(value: boolean) {
  amd.running = value;
}

export { amd, addLine, clearLines, setRunning };
