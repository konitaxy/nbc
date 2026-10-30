import { login, getUserInfo, setSelfInfo } from '@/api/user'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useRouterStore } from './router'

export const useUserStore = defineStore('user', () => {
  const loadingInstance = ref(null)

  const userInfo = ref({
    uuid: '',
    nickName: '',
    headerImg: '',
    authority: {},
    sideMode: 'dark',
    activeColor: '#4D70FF',
    baseColor: '#fff'
  })
  const token = ref(window.localStorage.getItem('token') || '')
  const setUserInfo = (val) => {
    userInfo.value = val
  }

  const setToken = (val) => {
    token.value = val
  }
  const setDevice = (val) => {
    token.value = val
  }

  const NeedInit = () => {
    token.value = ''
    window.localStorage.removeItem('token')
    localStorage.clear()
    router.push({ name: 'Init', replace: true })
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }
  /* 获取用户信息*/
  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }
  const completeLogin = async(data) => {
    setUserInfo(data.user)
    setToken(data.token)
    const routerStore = useRouterStore()
    await routerStore.SetAsyncRouter()
    routerStore.asyncRouters.forEach(asyncRouter => {
      router.addRoute(asyncRouter)
    })
    await router.push({ name: userInfo.value.authority.defaultRouter })
    return true
  }
  /* 登录。未绑定 2FA 或需要动态码时返回挑战数据，不写入登录态。*/
  const LoginIn = async(loginInfo) => {
    loadingInstance.value = ElLoading.service({
      fullscreen: true,
      text: '登陆中，请稍候...',
    })
    try {
      const res = await login(loginInfo)
      if (res.code === 0) {
        if (res.data?.needBind2FA || res.data?.need2FA) {
          return res.data
        }
        return await completeLogin(res.data)
      }
    } catch (e) {
      console.error(e)
    } finally {
      loadingInstance.value?.close()
    }
  }
  /* 登出*/
  const LoginOut = async() => {
    const res = await jsonInBlacklist()
    if (res.code === 0) {
      token.value = ''
      sessionStorage.clear()
      localStorage.clear()
      router.push({ name: 'Login', replace: true })
      window.location.reload()
    }
  }
  /* 设置侧边栏模式*/
  const changeSideMode = async(data) => {
    const res = await setSelfInfo({ sideMode: data })
    if (res.code === 0) {
      userInfo.value.sideMode = data
      ElMessage({
        type: 'success',
        message: '设置成功'
      })
    }
  }

  const mode = computed(() => userInfo.value.sideMode)
  const sideMode = computed(() => {
    if (userInfo.value.sideMode === 'dark') {
      return '#191a23'
    } else if (userInfo.value.sideMode === 'light') {
      return '#fff'
    } else {
      return userInfo.value.sideMode
    }
  })
  const baseColor = computed(() => {
    if (userInfo.value.sideMode === 'dark') {
      return '#fff'
    } else if (userInfo.value.sideMode === 'light') {
      return '#191a23'
    } else {
      return userInfo.value.baseColor
    }
  })
  const activeColor = computed(() => {
    if (userInfo.value.sideMode === 'dark' || userInfo.value.sideMode === 'light') {
      return '#4D70FF'
    }
    return userInfo.activeColor
  })

  watch(token, () => {
    window.localStorage.setItem('token', token.value)
  })

  return {
    userInfo,
    token,
    NeedInit,
    ResetUserInfo,
    GetUserInfo,
    LoginIn,
    completeLogin,
    LoginOut,
    changeSideMode,
    mode,
    sideMode,
    setToken,
    setDevice,
    baseColor,
    activeColor,
    loadingInstance
  }
})
