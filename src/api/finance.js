
import service from '@/utils/request'

export const listCardHolder = (data) => {
  return service({
    url: '/card/holder/list',
    method: 'post',
    data
  })
}

export const addCardHolder = (data) => {
  return service({
    url: '/card/holder/add',
    method: 'post',
    data
  })
}

export const fetchCardHolderAddress = (region = 'us') => {
  return service({
    url: '/card/holder/random-address',
    method: 'get',
    params: { region },
    donNotShowLoading: true,
  })
}

export const updateCardHolder = (data) => {
  return service({
    url: '/card/holder/update',
    method: 'post',
    data
  })
}

export const listCardBin = (data) => {
  return service({
    url: '/card/cardbin/list',
    method: 'post',
    data
  })
}

export const cancelCard = (data) => {
  return service({
    url: '/card/cancel',
    method: 'post',
    data
  })
}
export const createCard = (data) => {
  return service({
    url: '/card/add',
    method: 'post',
    data
  })
}

export const preRecharge = (data) => {
  return service({
    url: '/card/preRecharge',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}

export const rechargeCard = (data) => {
  return service({
    url: '/card/recharge',
    method: 'post',
    data
  })
}
export const withdrawCard = (data) => {
  return service({
    url: '/card/withdraw',
    method: 'post',
    data
  })
}
export const syncCard = (data) => {
  return service({
    url: '/card/sync',
    method: 'post',
    data
  })
}
export const remarkCard = (data) => {
  return service({
    url: '/card/remark',
    method: 'post',
    data
  })
}
export const listCard = (data) => {
  return service({
    url: '/card/list',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const staticsCard = (data) => {
  return service({
    url: '/card/statics',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const showCardDetail = (params) => {
  return service({
    url: '/card/detail',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}
export const cardReport = (data) => {
  return service({
    url: '/card/report',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const cardReportByDay = (data) => {
  return service({
    url: '/card/reportByDay',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const listCardTransactionRecord = (data) => {
  return service({
    url: '/card/transaction/list',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const getCardTransactionRecord = (params) => {
  return service({
    url: '/card/transaction',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}



export const getWalletBalance = () => {
  return service({
    url: '/wallet/balance',
    method: 'get',
    donNotShowLoading: true
  })
}

export const getGzyShareBalance = (data = {}) => {
  return service({
    url: '/wallet/gzy/share',
    method: 'post',
    data,
    donNotShowLoading: true
  })
}

export const gzyShareRecharge = (data) => {
  return service({
    url: '/wallet/gzy/recharge',
    method: 'post',
    data
  })
}

export const gzyShareWithdraw = (data) => {
  return service({
    url: '/wallet/gzy/withdraw',
    method: 'post',
    data
  })
}

export const walletRechargeApply = (data) => {
  return service({
    url: '/wallet/recharge/apply',
    method: 'post',
    data
  })
}

export const walletWithdrawApply = (data) => {
  return service({
    url: '/wallet/withdraw/apply',
    method: 'post',
    data
  })
}
export const listWalletWithdraw = (data) => {
  return service({
    url: '/wallet/withdraw/list',
    method: 'post',
    data
  })
}
export const listWalletRecharge = (data) => {
  return service({
    url: '/wallet/recharge/list',
    method: 'post',
    data
  })
}
export const listRechargeRecord = (data) => {
  return service({
    url: '/wallet/recharge/list',
    method: 'post',
    data
  })
}

export const listWalletHistory = (data) => {
  return service({
    url: '/wallet/history',
    method: 'post',
    data
  })
}
export const walletReport = (data) => {
  return service({
    url: '/wallet/report',
    method: 'post',
    data
  })
}

export const addCardGroup = (data) => {
  return service({
    url: '/card/group',
    method: 'post',
    data
  })
}
export const listCardGroup = (data) => {
  return service({
    url: '/card/group/list',
    method: 'post',
    data
  })
}
export const delCardGroup = (params) => {
  return service({
    url: '/card/group',
    method: 'delete',
    params
  })
}
export const setCardGroup = (data) => {
  return service({
    url: '/card/setGroup',
    method: 'post',
    data
  })
}
export const adjustSubCardLimit = (data) => {
  return service({
    url: '/card/adjustLimit',
    method: 'post',
    data
  })
}