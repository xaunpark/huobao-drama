/**
 * Worker-based timer that is NOT throttled when browser tab is in background.
 * 
 * Modern browsers aggressively throttle setTimeout/setInterval in background tabs
 * (from ~2s up to 60s+ minimum intervals). This causes batch processing to stall
 * because polling loops stop detecting task completion while the tab is not visible.
 * 
 * This module creates a shared Web Worker that handles timing, so all timers
 * continue to fire at their scheduled intervals regardless of tab visibility.
 */

let worker: Worker | null = null;
let callbackMap = new Map<number, () => void>();
let nextId = 0;

const getWorker = (): Worker => {
  if (worker) return worker;

  const blob = new Blob([`
    // Web Worker timer - immune to background tab throttling
    const timers = new Map();
    
    self.onmessage = function(e) {
      const { type, id, delay } = e.data;
      
      if (type === 'setTimeout') {
        const timer = setTimeout(() => {
          self.postMessage({ id });
          timers.delete(id);
        }, delay);
        timers.set(id, { type: 'timeout', timer });
      } 
      else if (type === 'setInterval') {
        const timer = setInterval(() => {
          self.postMessage({ id });
        }, delay);
        timers.set(id, { type: 'interval', timer });
      }
      else if (type === 'clear') {
        const entry = timers.get(id);
        if (entry) {
          if (entry.type === 'timeout') clearTimeout(entry.timer);
          else clearInterval(entry.timer);
          timers.delete(id);
        }
      }
    };
  `], { type: 'application/javascript' });

  worker = new Worker(URL.createObjectURL(blob));
  worker.onmessage = (e: MessageEvent) => {
    const cb = callbackMap.get(e.data.id);
    if (cb) {
      cb();
      // For setTimeout, clean up after firing
      // For setInterval, keep it alive (caller must clear manually)
    }
  };

  return worker;
};

/**
 * A drop-in replacement for `new Promise(r => setTimeout(r, delay))` 
 * that is NOT throttled in background tabs.
 * 
 * Usage: `await workerDelay(2000)` instead of `await new Promise(r => setTimeout(r, 2000))`
 */
export const workerDelay = (delay: number): Promise<void> => {
  return new Promise((resolve) => {
    const id = nextId++;
    callbackMap.set(id, () => {
      callbackMap.delete(id);
      resolve();
    });
    getWorker().postMessage({ type: 'setTimeout', id, delay });
  });
};

/**
 * Worker-based setInterval that is NOT throttled in background tabs.
 * Returns an ID that can be used to clear the interval.
 */
export const workerSetInterval = (callback: () => void, delay: number): number => {
  const id = nextId++;
  callbackMap.set(id, callback);
  getWorker().postMessage({ type: 'setInterval', id, delay });
  return id;
};

/**
 * Clear a worker-based interval or timeout.
 */
export const workerClearInterval = (id: number): void => {
  callbackMap.delete(id);
  if (worker) {
    worker.postMessage({ type: 'clear', id });
  }
};

/**
 * Cleanup: terminate the shared worker when no longer needed.
 * Call this on app unmount if desired.
 */
export const terminateWorkerTimer = (): void => {
  if (worker) {
    worker.terminate();
    worker = null;
  }
  callbackMap.clear();
  nextId = 0;
};
