import service from '@/utils/request'
// @Tags finance

export const addFeeGlobalCfg = (data) => {
  return service({
    url: '/admin/finance/fee/addGlobalCfg',
    method: 'post',
    data
  })
}
export const addFeeUserCfg = (data) => {
  return service({
    url: '/admin/finance/fee/addUserCfg',
    method: 'post',
    data
  })
}
export const addFeeMonthCfg = (data) => {
  return service({
    url: '/admin/finance/fee/addMonthCfg',
    method: 'post',
    data
  })
}

export const setUserCfgGlobal = (data) => {
  return service({
    url: '/admin/finance/fee/setUserCfgGlobal',
    method: 'post',
    data
  })
}

export const listFeeUserCfg = (data) => {
  return service({
    url: '/admin/finance/fee/list/user',
    method: 'post',
    data
  })
}
export const listFeeGlobalCfg = (data) => {
  return service({
    url: '/admin/finance/fee/list/global',
    method: 'post',
    data
  })
}

export const listLogs = (data) => {
  return service({
    url: '/admin/common/log/list',
    method: 'post',
    data
  })
}

export const listSmsCode = (data) => {
  return service({
    url: '/admin/common/smscode/list',
    method: 'post',
    data
  })
}



export const listWithdrawRecord = (data) => {
  return service({
    url: '/admin/finance/wallet/withdraw/list',
    method: 'post',
    data
  })
}

export const reviewWithdrawRecord = (data) => {
  return service({
    url: '/admin/finance/wallet/withdraw/review',
    method: 'post',
    data
  })
}


export const listRechargeRecord = (data) => {
  return service({
    url: '/admin/finance/wallet/recharge/list',
    method: 'post',
    data
  })
}

export const reviewRechargeRecord = (data) => {
  return service({
    url: '/admin/finance/wallet/recharge/review',
    method: 'post',
    data
  })
}

export const editRechargeRecord = (data) => {
  return service({
    url: '/admin/finance/wallet/recharge/edit',
    method: 'post',
    data
  })
}

export const statBalance = (params) => {
  return service({
    url: '/admin/finance/report/balance',
    method: 'get',
    params
  })
}
export const statReport = (data) => {
  return service({
    url: '/admin/finance/report/all',
    method: 'post',
    data
  })
}
export const listStatReport = (data) => {
  return service({
    url: '/admin/finance/report/list',
    method: 'post',
    data
  })
}

export const listStatReportByClient = (data) => {
  return service({
    url: '/admin/finance/report/listByClient',
    method: 'post',
    data
  })
}

export const addChainWatchAddress = (data) => {
  return service({
    url: '/admin/finance/chain/address/add',
    method: 'post',
    data
  })
}

export const deleteChainWatchAddress = (data) => {
  return service({
    url: '/admin/finance/chain/address/delete',
    method: 'post',
    data
  })
}

export const listChainWatchAddress = (data) => {
  return service({
    url: '/admin/finance/chain/address/list',
    method: 'post',
    data
  })
}

export const listChainInboundTransaction = (data) => {
  return service({
    url: '/admin/finance/chain/transaction/list',
    method: 'post',
    data
  })
}






