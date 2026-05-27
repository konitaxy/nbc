// src/stores/languageStore.js
import { defineStore } from 'pinia';
import i18n from '@/i18n'
export const useLanguageStore = defineStore('language', {
  state: () => ({
    currentLocale: 'zh' 
  }),
  actions: {
    changeLocale(locale) {
      this.currentLocale = locale;
      i18n.global.locale.value = locale
    },
    t(code) {
      return i18n.global.t(code)
    }
  },
  persist: true
});