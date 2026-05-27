import service from '@/utils/request'

export const addCardBin = (data) => {
  return service({
    url: '/admin/card/cardBin/add',
    method: 'post',
    data
  })
}
export const blockCardBin = (data) => {
  return service({
    url: '/admin/card/cardBin/block',
    method: 'post',
    data
  })
}

export const listCardBin = (data) => {
  return service({
    url: '/admin/card/cardBin/list',
    method: 'post',
    data
  })
}


export const listCardHolder = (data) => {
  return service({
    url: '/admin/finance/holder/list',
    method: 'post',
    data
  })
}

export const listCard = (data) => {
  return service({
    url: '/admin/card/list',
    method: 'post',
    data
  })
}

export const listCardTransaction = (data) => {
  return service({
    url: '/admin/card/transaction/list',
    method: 'post',
    data
  })
}

export const rechargeCard = (data) => {
  return service({
    url: '/admin/card/recharge',
    method: 'post',
    data
  })
}

export const cancelCard = (data) => {
  return service({
    url: '/admin/card/cancel',
    method: 'post',
    data
  })
}

export const syncCard = (data) => {
  return service({
    url: 'admin/card/sync',
    method: 'post',
    data
  })
}

export const frozenCard = (data) => {
  return service({
    url: '/admin/card/frozen',
    method: 'post',
    data
  })
}


