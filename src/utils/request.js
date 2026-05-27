import axios from 'axios'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'
import { emitter } from '@/utils/bus.js'
import { sendVerifyCode } from '@/api/profile'
import {useLanguageStore} from '@/pinia/modules/language'
import router from '@/router/index'
const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 99999
})
let acitveAxios = 0
let timer
const showLoading = () => {
  acitveAxios++
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    if (acitveAxios > 0) {
      emitter.emit('showLoading')
    }
  }, 400)
}

const closeLoading = () => {
  acitveAxios--
  if (acitveAxios <= 0) {
    clearTimeout(timer)
    emitter.emit('closeLoading')
  }
}
// http request 拦截器
service.interceptors.request.use(
  config => {
    if (!config.donNotShowLoading) {
      showLoading()
    }
    const userStore = useUserStore()
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': userStore.token,
      'x-user-id': userStore.userInfo.ID,
      ...config.headers
    }
    return config
  },
  error => {
    closeLoading()
    ElMessage({
      showClose: true,
      message: error,
      type: 'error'
    })
    return error
  }
)

// http response 拦截器
service.interceptors.response.use(
  async (response) => {
    const userStore = useUserStore()
    closeLoading()
    if (response.headers['new-token']) {
      userStore.setToken(response.headers['new-token'])
    }

    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else if(response.data.code === 14){
      const language = useLanguageStore()

      const originalRequest = response.config;
        var t = response.headers['x-auth-type']
        if (!response.headers['x-auth-required']){
          ElNotification({
            type: 'error',
            message: `Try to login again`,
            title: 'Request Authorization First',
          })
          return
        }
        if (t == 'otcp'){
         return ElMessageBox.prompt('', language.t('lang.mfa_verify_code'), {
            roundButton:true,
            center:true,
            inputPlaceholder: language.t('lang.verify_code_placeholder'),
            confirmButtonText: language.t('lang.confirm'),
            cancelButtonText: language.t('lang.cancel'),
            inputErrorMessage: 'Invalid code',
            beforeClose: async(action,instance,done) =>{
              if(action === 'confirm'){
                originalRequest.headers['X-Auth-Code'] = instance.inputValue
                var res = await service(originalRequest)
                if(res.code === 0){
                  instance.inputValue= res
                  done()
                }
              }else {
                done()
              }
           }
         }).then(async({value}) => {
            if(originalRequest.url == '/client/login'){
              const userStore = useUserStore()
              userStore.LoginInRes(value)
            }else {
             return value
            }
         })
            // .then(async({ value }) => {
            //   originalRequest.headers['X-Auth-Code'] = value
            //   if(originalRequest.url == '/client/login'){
            //     const userStore = useUserStore()
            //     service(originalRequest).then(res =>{
            //       userStore.LoginInRes(res)
            //     })
            //   }else {
            //     return await service(originalRequest)
            //   }
            // })
            .catch(() => {
            })
        }else if (t == 'pin'){
          return ElMessageBox.prompt('', language.t('lang.pin_verify_code'), {
             roundButton:true,
             center:true,
             inputPlaceholder: language.t('lang.verify_code_placeholder'),
             confirmButtonText: language.t('lang.confirm'),
             cancelButtonText: language.t('lang.cancel'),
             inputErrorMessage: 'Invalid code',
             beforeClose: async(action,instance,done) =>{
                if(action === 'confirm'){
                  originalRequest.headers['X-Auth-Code'] = instance.inputValue
                  var res = await service(originalRequest)
                  if(res.code === 0){
                    instance.inputValue= res
                    done()
                  }
                }else {
                  done()
                }
             }
           }).then(async({value}) => {
                return value
            })
             .catch(() => {
             })
         }else {
          
            return await sendVerifyCode({
              to: originalRequest.headers['X-Auth-Mail'],
              path: `${originalRequest.method.toUpperCase()}:${originalRequest.url}`
            }).then(async(sres)=>{
              if (sres.code === 0){
                return ElMessageBox.prompt(language.t('lang.verify_code_sent'), language.t('lang.email_verify_code'), {
                  inputPlaceholder: language.t('lang.verify_code_placeholder'),
                  confirmButtonText: language.t('lang.confirm'),
                  cancelButtonText: language.t('lang.cancel'),
                  inputErrorMessage: 'Invalid code',
                  beforeClose: async(action,instance,done) =>{
                    if(action === 'confirm'){
                      originalRequest.headers['X-Auth-Code'] = instance.inputValue
                      var res = await service(originalRequest)
                      if(res.code === 0){
                        instance.inputValue= res
                        done()
                      }
                    }else {
                      done()
                    }
                 }
               }).then(async({value}) => {
                  if(originalRequest.url == '/client/login'){
                    const userStore = useUserStore()
                    userStore.LoginInRes(value)
                  }else {
                   return value
                  }
               }).catch(() => {
                  })
              }
            })
        }
    } else if (response.data.code === 15){
        const language = useLanguageStore()
        ElMessageBox.confirm(
          language.t('lang.need_completed_kyc_required'),
          language.t('lang.kyc_required_warning'),
          {
            confirmButtonText: language.t('lang.kyc_now'),
            type: 'warning',
          }
        ).then(() => {
          router.push({name: 'identityVerify',replace: true})
        })
    } else if (response.data.code === 16){
      const language = useLanguageStore()
      ElMessageBox.confirm(
        language.t('lang.kyc_pending_review'),
        language.t('lang.kyc_required_warning'),
        {
          confirmButtonText: language.t('lang.confirm'),
          type: 'warning',
        }
      ).then(() => {
      })
    } else{
      if (response.data.msg) {
        const lang = useLanguageStore()
        response.data.msg = lang.t(response.data.msg)
        response.data.msg = response.data.msg.charAt(0).toUpperCase() + response.data.msg.slice(1);
      }
      ElMessage({
        showClose: true,
        message: response.data.msg || decodeURI(response.headers.msg),
        type: 'error',
        duration: 5000,
        customClass:'el-message-box-custom'
      })
      if (response.data.data && response.data.data.reload) {
        userStore.token = ''
        router.push({ name: 'Login', replace: true })
      }
      return response.data.msg ? response.data : response
    }
  },
  error => {
    closeLoading()
    switch (error.response.status) {
      case 500:
        ElMessageBox.confirm(`
        <p>Interface error detected${error}</p>
        <p>Oops! Something went wrong on our end. Please try again later or contact support if the issue continues.</p>
        `, 'Interface error', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          cancelButtonText: 'Cancel'
        })
          .then(() => {
            const userStore = useUserStore()
            userStore.token = ''
            localStorage.clear()
            router.push({ name: 'Login', replace: true })
          })
        break
      case 404:
        ElMessageBox.confirm(`
          <p>Interface error detected${error}</p>
          <p>Error code<span style="color:red"> 404 </span>：This type of error is mostly caused by the interface not being registered (or not restarted) or the request path (method) not matching the API path (method) - if it is automated code, please check whether there are spaces.</p>
          `, 'Interface error', {
          dangerouslyUseHTMLString: true,
          distinguishCancelAndClose: true,
          confirmButtonText: 'I see',
          cancelButtonText: 'Cancel'
        })
        break
      case 413:
          ElNotification({
            type: 'error',
            message: `Upload file exceeds the maximum allowed size of 64MB. Please reduce the file size or try uploading a smaller file`,
            title: 'Request Entity Too Large',
          })
          break
      case 401:
        const originalRequest = error.config;
        var t = error.response.headers['x-auth-type']
        if (!error.response.headers['x-auth-required']){
          ElNotification({
            type: 'error',
            message: `Try to login again`,
            title: 'Request Authorization First',
          })
          return
        }
        if (t == 'otcp'){
          ElMessageBox.prompt('', 'Two-factor authentication code', {
            roundButton:true,
            center:true,
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            inputErrorMessage: 'Invalid code',
          })
            .then(({ value }) => {
              originalRequest.headers['X-Auth-Code'] = value
              if(originalRequest.url == '/client/login'){
                const userStore = useUserStore()
                service(originalRequest).then(res =>{
                  userStore.LoginInRes(res)
                })
              }else {
                service(originalRequest).then(res =>{ 
                  error.response=res
                  if(res.code === 0){
                    ElMessage.success(res.msg)
                  }
                })
              }
            })
            .catch(() => {
            })
        }else {
          ElMessageBox.confirm('Ok to Send a code to your email', 'Send Email Code', {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            inputErrorMessage: 'Invalid code',
          }).then(()=>{
            const data = {
              to: originalRequest.headers['X-Auth-Mail'],
              path: `${originalRequest.method.toUpperCase()}:${originalRequest.url}`
            }
            sendVerifyCode(data).then((sres)=>{
              if (sres.code === 0){
                ElMessageBox.prompt('', 'Email authentication code', {
                  confirmButtonText: 'OK',
                  cancelButtonText: 'Cancel',
                  inputErrorMessage: 'Invalid code',
                }).then(({ value }) => {
                    originalRequest.headers['X-Auth-Code'] = value
                    if(originalRequest.url == '/client/login'){
                      const userStore = useUserStore()
                      service(originalRequest).then(res =>{
                        userStore.LoginInRes(res)
                      })
                      
                    }else {
                      service(originalRequest)
                    }
                  })
                  .catch(() => {
                  })
              }
            })
          })
            .catch(() => {
            })
          
        }
        break
    }

    return error
  }
)
export default service
