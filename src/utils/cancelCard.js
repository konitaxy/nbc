/**
 * 批量销卡 POST /card/cancel：{ list: [{ id, cardId }], ... }
 * row 使用列表接口的 ID、cardId 字段
 */
export function buildCancelListPayload (rows) {
  const arr = (Array.isArray(rows) ? rows : [rows]).filter(Boolean)
  const seen = new Set()
  const list = []
  for (const r of arr) {
    const id = r.ID
    if (id == null) continue
    const cardId = r.cardId != null ? String(r.cardId) : ''
    if (!cardId) continue
    if (seen.has(id)) continue
    seen.add(id)
    list.push({ id, cardId })
  }
  if (list.length < 1) {
    throw new Error('cancel_list_empty')
  }
  if (list.length > 100) {
    throw new Error('cancel_list_too_many')
  }
  return { list }
}

/**
 * 处理 /card/cancel 返回：code 0 可能部分失败；code 7 为整批失败
 */
export function applyCancelCardResult (res, { t, ElMessage, onAfter }) {
  const tr = t || ((k) => k)
  if (res.code === 7) {
    ElMessage.error(res.msg || 'Error')
    const f = res.data?.failed
    if (f?.length) {
      const detail = f.map((x) => `${x.id}: ${x.reason || ''}`).join('; ')
      ElMessage({ type: 'warning', message: detail, duration: 8000, showClose: true })
    }
    onAfter && onAfter(false, res)
    return
  }
  if (res.code !== 0) {
    ElMessage.error(res.msg || 'Error')
    onAfter && onAfter(false, res)
    return
  }
  const d = res.data || {}
  const success = d.success
  const failed = d.failed || []
  if (failed.length) {
    const detail = failed.map((x) => `${x.id}: ${x.reason || ''}`).join('; ')
    ElMessage({
      type: 'warning',
      message: tr('lang.cancel_batch_partial', { success: success ?? 0, fail: failed.length }) + (detail ? ` — ${detail}` : ''),
      duration: 10000,
      showClose: true
    })
  } else {
    ElMessage.success(tr('lang.success'))
  }
  onAfter && onAfter(true, res)
}
