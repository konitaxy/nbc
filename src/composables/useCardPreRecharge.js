import { ref, watch, onUnmounted } from 'vue'
import { preRecharge as preRechargeApi } from '@/api/finance'

const DEBOUNCE_MS = 400

function parseAmount(value) {
  if (value === null || value === undefined || value === '') return null
  const n = parseFloat(String(value).replace(/,/g, '').trim())
  if (Number.isNaN(n) || n <= 0) return null
  return n.toFixed(2)
}

/**
 * @param {Object} options
 * @param {() => string|undefined} options.getCardId
 * @param {() => string|number|undefined} options.getRechargeAmount 用户输入的到账金额
 * @param {() => boolean|undefined} [options.enabled] dialog visible
 * @param {() => string} [options.getMemberId]
 * @param {() => string} [options.getAccountId]
 */
export function useCardPreRecharge(options) {
  const summary = ref(null)
  const loading = ref(false)
  let debounceTimer = null
  let seq = 0

  async function runFetch() {
    const cardId = options.getCardId?.()
    const arrivalAmount = parseAmount(options.getRechargeAmount?.())
    const memberId = options.getMemberId?.() ?? ''
    const accountId = options.getAccountId?.() ?? ''

    if (!cardId || !arrivalAmount) {
      summary.value = null
      loading.value = false
      return
    }

    const current = ++seq
    loading.value = true
    try {
      const res = await preRechargeApi({
        memberId,
        accountId,
        cardId,
        rechargeAmount: null,
        arrivalAmount,
      })
      if (current !== seq) return
      summary.value = res.code === 0 ? res.data : null
    } catch {
      if (current === seq) summary.value = null
    } finally {
      if (current === seq) loading.value = false
    }
  }

  function scheduleFetch() {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(runFetch, DEBOUNCE_MS)
  }

  function reset() {
    seq++
    clearTimeout(debounceTimer)
    summary.value = null
    loading.value = false
  }

  watch(
    () => {
      const enabled = options.enabled?.()
      if (enabled === false) return ['__closed__']
      return [
        enabled ?? true,
        options.getCardId?.(),
        options.getRechargeAmount?.(),
      ]
    },
    (vals) => {
      if (vals[0] === '__closed__') {
        reset()
        return
      }
      scheduleFetch()
    },
    { immediate: true }
  )

  onUnmounted(reset)

  return { summary, loading, reset }
}
