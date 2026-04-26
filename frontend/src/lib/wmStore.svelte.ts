const wm = $state({
  running: false,
});

function setWmRunning(value: boolean) {
  wm.running = value;
}

export { wm, setWmRunning };
