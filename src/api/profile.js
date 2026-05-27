
import service from '@/utils/request'


export const login = (data) => {
  return service({
    url: '/client/login',
    method: 'post',
    headers: { 'X-Auth-Mail': data.email },
    data
  })
}
export const applyArtist = (data) => {
  return service({
    url: '/client/applyArtist',
    method: 'post',
    data
  })
}

export const setAvatar = (data) => {
    return service({
      url: '/client/avatar',
      method: 'post',
      headers: { 'Content-Type': 'multipart/form-data' },
      data
    })
}

export const setBgImg = (data) => {
  return service({
    url: '/client/bgImg',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data
  })
}

export const setProfile = (data) => {
  return service({
    url: '/client/profile',
    method: 'post',
    data
  })
}

export const changeCouponCode = (params) => {
  return service({
    url: '/client/changeCouponCode',
    method: 'get',
    params
  })
}

export const addCouponCode = (params) => {
  return service({
    url: '/client/addCouponCode',
    method: 'get',
    params
  })
}

export const setPayment = (data) => {
  return service({
    url: '/client/payment',
    method: 'post',
    data
  })
}
export const getPayment = (params) => {
  return service({
    url: '/client/payment',
    method: 'get',
    params
  })
}

export const getInviteInfo = (params) => {
  return service({
    url: '/client/inviteInfo',
    method: 'get',
    params
  })
}

export const getRecentTags = (params) => {
  return service({
    url: '/client/recentTags',
    method: 'get',
    params
  })
}

export const adminLogin = (params) => {
  return service({
    url: '/client/rologin',
    method: 'get',
    params
  })
}

export const sendVerifyCode = (data) => {
  return service({
    url: '/client/verifyCode',
    method: 'post',
    donNotShowLoading:true,
    data
  })
}

export const changePassword = (data) => {
  return service({
    url: '/client/changePassword',
    method: 'post',
    data
  })
}

export const setPIN = (data) => {
  return service({
    url: '/client/pin',
    method: 'post',
    data
  })
}


export const register = (data) => {
  return service({
    url: '/client/register',
    method: 'post',
    headers: { 'X-Auth-Mail': data.email },
    data
  })
}

export const resetPassword = (data) => {
  
  return service({
    url: '/client/resetPassword',
    method: 'post',
    headers: { 'X-Auth-Mail': data.email },
    data
  })
}

export const  getShopifyProfile = async(params) => {
  return service({
    url: '/admin/artwork/getShopifyProfile',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

export const  captcha = async(params) => {
  return service({
    url: '/client/captcha',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

export const  getDueDiligence = async(params) => {
  return service({
    url: '/client/dueDiligence',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

export const  setDueDiligence = async(data) => {
  return service({
    url: '/client/dueDiligence',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const listSessionLog = (params) => {
  return service({
    url: '/client/sessionLog',
    method: 'get',
    params
  })
}







