import { createApp } from 'vue'
import 'element-plus/dist/index.css'
import './style/element_visiable.scss'
import ElementPlus from 'element-plus'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap-icons/font/bootstrap-icons.css'

import en from 'element-plus/dist/locale/en.mjs'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

// import vn from 'element-plus/dist/locale/vi.mjs'

// import gin-vue-admin font install
import './core/metalposter-admin'
// import router
import router from '@/router/index'
import '@/permission'
import run from '@/core/metalposter-admin.js'
import auth from '@/directive/auth'
import { store } from '@/pinia'
import { useLanguageStore } from '@/pinia/modules/language'
import { useUserStore } from '@/pinia/modules/user'
import App from './App.vue'
import { mainLogo } from '@/core/config'
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import i18n from '@/i18n'
const app = createApp(App)
app.config.productionTip = false
store.use(piniaPluginPersistedstate); // 使用插件
mainLogo(window)
app 
  .use(run)
  .use(store)
  .use(auth)
  .use(router)
  .use(ElementPlus, { locale: en })
  .use(i18n)
  const languageStore = useLanguageStore()
  const userStore = useUserStore()
  
  // 设置 userStore 为全局变量，模板中可用 $userStore 访问
  app.config.globalProperties.$userStore = userStore
  
  i18n.global.locale.value = languageStore.currentLocale;
  const rootSelector = '#app'

app.mount(rootSelector)

export default app
