import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import getPageTitle from '@/utils/page'
import router from '@/router'

let asyncRouterFlag = 0

const whiteList = ['PixelCard','Login', 'Init','Invite','forgetPassword','adminLogin']
const getRouter = async(userStore) => {
  
  const routerStore = useRouterStore()
  await routerStore.SetAsyncRouter()
  await userStore.GetUserInfo()
  const asyncRouters = routerStore.asyncRouters
  asyncRouters.forEach(asyncRouter => {
    router.addRoute(asyncRouter)
  })
}

async function handleKeepAlive(to) {
  if (to.matched && to.matched.length > 2) {
    for (let i = 1; i < to.matched.length; i++) {
      const element = to.matched[i - 1]
      if (element.name === 'layout') {
        to.matched.splice(i, 1)
        await handleKeepAlive(to)
      }
      // after loaded
      if (typeof element.components.default === 'function') {
        await element.components.default()
        await handleKeepAlive(to)
      }
    }
  }
}

router.beforeEach(async(to, from, next) => {
  const userStore = useUserStore()
  to.meta.matched = [...to.matched]
  handleKeepAlive(to)
  const token = userStore.token
  // check in whitelist
  document.title = getPageTitle(to.meta.title,userStore.userInfo)
  if (whiteList.indexOf(to.name) > -1) {
    if (token) {
      if(to.name === 'adminLogin' && from.name !== 'adminLogin'){
        next()
      }

      if (!asyncRouterFlag && whiteList.indexOf(from.name) < 0) {
        asyncRouterFlag++
        await getRouter(userStore)
      }
      next({ name: userStore.userInfo.authority.defaultRouter })
    } else {
      next()
    }
  } else {
    if (token) {
      // add flag to prevent overflow
      if (!asyncRouterFlag && whiteList.indexOf(from.name) < 0) {
        asyncRouterFlag++
        await getRouter(userStore)
        next({ ...to, replace: true })
      } else {
        if (to.matched.length) {
          next()
        } else {
          next({ path: '/layout/404' })
        }
      }
    }
    // when not in whitelist
    if (!token) {
      next({
        name: 'PixelCard',
        query: {
          // redirect: document.location.hash
        }
      })
    }
  }
})
