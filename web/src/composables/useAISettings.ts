import { ref, watch } from 'vue'

const STORAGE_KEY = 'ai_max_concurrent_threads'
const DEFAULT_THREADS = 4
const MIN_THREADS = 1
const MAX_THREADS = 16

/**
 * Composable for managing global AI generation settings.
 * Currently stores max concurrent threads for batch image/video generation.
 * Settings persist via localStorage.
 */
export function useAISettings() {
  const maxConcurrentThreads = ref(loadSetting())

  function loadSetting(): number {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored !== null) {
        const parsed = parseInt(stored, 10)
        if (!isNaN(parsed) && parsed >= MIN_THREADS && parsed <= MAX_THREADS) {
          return parsed
        }
      }
    } catch {
      // localStorage unavailable
    }
    return DEFAULT_THREADS
  }

  function saveSetting(value: number) {
    const clamped = Math.max(MIN_THREADS, Math.min(MAX_THREADS, value))
    localStorage.setItem(STORAGE_KEY, String(clamped))
  }

  // Auto-persist changes
  watch(maxConcurrentThreads, (newVal) => {
    saveSetting(newVal)
  })

  return {
    maxConcurrentThreads,
    MIN_THREADS,
    MAX_THREADS,
    DEFAULT_THREADS
  }
}
