import { createI18n } from 'vue-i18n'
import enLocale from './en'
import zhLocale from './zh'
import vnLocale from './vn'

const messages = {
  en: enLocale,
  zh: zhLocale,
  vn: vnLocale,
}

const i18n = createI18n({
  legacy: false,
  locale: 'en', // 默认语言 
  fallbackLocale: 'zh', // 备用语言
  messages,
})

export default i18n