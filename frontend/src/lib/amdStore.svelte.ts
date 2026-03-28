const amd = $state({
  running: false,
});

function setRunning(value: boolean) {
  amd.running = value;
}

export { amd, setRunning };
