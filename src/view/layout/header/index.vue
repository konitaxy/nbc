<template>
    <div class="console-header h-100 w-100">
      <el-row justify="start" class="h-100 w-100 align-items-center">
        <el-col :xs="4" :sm="4" :md="0">
          <button class="console-menu-toggle d-block d-lg-none" type="button" aria-label="Toggle sidebar" @click="handleExpland">
            <i class="bi bi-list"></i>
          </button>
        </el-col>
        <el-col :span="8" :xs="10" :sm="8" :md="10">
          <div class="logo-container">
            <img src="@/assets/NEWBEECARD-logo.png" class="logo-img" alt="NEWBEECARD" />
            <div class="console-brand-copy d-none d-md-flex">
              <strong>NEWBEECARD</strong>
              <span>Global subscription cards</span>
            </div>
            <span class="console-brand-tag d-none d-xl-inline-flex">Control center</span>
          </div>
          <!-- <img :src="logo" class="logo" style="max-width: 186px; height:75px;object-fit:cover;"></img> -->
        </el-col>
        <el-col :xs="10" :sm="12" :md="14">
          <el-col  class="console-actions h-100 align-items-center justify-content-end d-flex me-2 gap-3">
                <div class="console-client-pill d-none d-sm-flex">
                  <span>{{ $t('lang.clientNo') }}</span>
                  <strong>{{ userStore.userInfo.clientNo }}</strong>
                </div>
                <div v-if="loginEmail" class="console-email-pill d-none d-lg-flex">
                  <i class="bi bi-envelope"></i>
                  <strong>{{ loginEmail }}</strong>
                </div>
                <!-- <i class="bi bi-envelope fs-5 me-4"></i> -->
                <!-- <i class="bi bi-envelope fs-5 me-4"></i> -->
                <el-dropdown class="console-dropdown me-2" popper-class="px-2">
                  <a href="#" class="console-language-link text-nowrap" style="line-height: normal;">
                    {{ langStore.currentLocale ==='zh'?'中文':langStore.currentLocale ==='en'?'English':'Tiếng Việt' }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
                  </a>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="handleLangSet('zh')" class="justify-content-center">中文</el-dropdown-item>
                      <el-dropdown-item @click="handleLangSet('en')" class="justify-content-center">English</el-dropdown-item>
                      <el-dropdown-item @click="handleLangSet('vn')" class="justify-content-center">Tiếng Việt</el-dropdown-item>

                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <!-- <i class="bi bi-bell fs-3 me-2"></i> -->
                <!-- <i class="bi bi-question-circle fs-5 me-3"></i> -->
                <el-dropdown class="console-avatar-dropdown me-2" popper-class="px-2">
                  <span class="console-avatar-trigger el-dropdown-link">
                    <el-avatar shape="circle" :size="38" :src="avatar"></el-avatar>
                    <!-- <el-icon class="fs-4 el-icon--right"><arrow-down /></el-icon> -->
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item v-if="loginEmail" disabled class="console-email-dropdown-item justify-content-start">
                        <span>Email</span>
                        <strong>{{ loginEmail }}</strong>
                      </el-dropdown-item>
                      <el-dropdown-item @click="handleGoto('accountInfo')" class="justify-content-start">{{ $t('lang.detail') }}</el-dropdown-item>
                      <el-dropdown-item @click="handleGoto('googleVerify')" class="justify-content-start">{{ $t('lang.auth') }}</el-dropdown-item>
                      <el-dropdown-item ><el-button type="text" class="text-nowrap w-100 justify-content-start" @click="loginout">{{ $t('lang.loginout') }}</el-button></el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-col>
       
        </el-col>
      </el-row>
      </div>
      <el-dialog
        v-model="earningGuideVisible"
        title="Earning Guide"
        width="60%"
        append-to="html"
      >
        <iframe v-if="earningGuideVisible" width="100%" height="450px" src="https://www.youtube.com/embed/Ei-7a6LmbU0?si=bqeqw2-0sHIQIIYc" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
      </el-dialog>

      <el-dialog
        v-model="contactHelpVisible"
        title="Need help? Get in touch anytime!"
        max-width="60%"
        append-to="html"
        align-center
      >
        <div class="row gap-3">
          <div class="col-12 fs-5 text-secondary">
            We're here for you — choose your preferred way to contact us:
          </div>
          <div class="col-9 d-flex gap-2 ms-4 fs-5 text-nowrap">
            <i class="bi bi-envelope"></i><p>Email: <span class="link-primary">support@metalposter.com</span></p>
          </div>
          <div class="col-9 d-flex gap-2 ms-4 fs-5 text-nowrap">
            <i class="bi bi-whatsapp text-success"></i><p>Whatsapp: <a href="https://wa.me/+6596659699" target="_blank" class="link-primary">+65 96659699</a></p>
          </div>
          <div class="col-9 d-flex gap-2 ms-4 fs-5 text-nowrap">
            <i class="bi bi-telegram text-primary"></i><p>Telegram: <a href="https://t.me/metalpostercom" target="_blank" class="link-primary">@metalpostercom</a></p>
          </div>
        </div>
        <template #footer>
          <button class="btn btn-secondary mt-1" @click="contactHelpVisible = false">Close</button>
        </template>
      </el-dialog>
  </template>
  
  <script>
export default {
  name: 'Header',
}
</script>
<script setup>
import {computed, ref,onMounted} from 'vue'
import { emitter } from '@/utils/bus.js'
import logoArt from '@/assets/logo.png'
import avatar from '@/assets/avatar-default.png'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import { on } from 'screenfull'
import i18n from '@/i18n'
import { useLanguageStore } from '@/pinia/modules/language';
const langStore = useLanguageStore()
const userStore = useUserStore()
const router = useRouter()
const loginEmail = computed(() => {
  return userStore.userInfo.email || userStore.userInfo.userName || userStore.userInfo.username || ''
})
const loginout = ()=>{
    userStore.LoginOut()
}
  // No script needed in this example.

const emit = defineEmits(['expland']);
const isMobile = ref(false)
const earningGuideVisible = ref(false)
const contactHelpVisible = ref(false)
const logo = logoArt

onMounted(() => {
  
})
const handleLangSet = (lang) =>{
  langStore.changeLocale(lang)
}
const handleExpland = () => {
  emit('expland')
}
emitter.on('mobile', (item) => {
  isMobile.value = item
  })
  const handleGoto = (name) => {
    router.push({name:name})
  }
  // No script needed in this example.
</script>
  
  <style scoped>
  .header {
    padding: 0;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  }
  .right-icons-container {
    width: Fixed (205.22px)px;
    height: Hug (25px)px;
    top: 20px;
    left: 1366.78px;
    gap: 32px;
    opacity: 0px;

  }
  .right-icons {
    width: 8.13px;
    height: 13.46px;
    top:5.23px;
    left:7.88px;
    margin-bottom: 10px;
    color: var(--bs-secondary-color);
    cursor: pointer;
  }
  .right-icons-mobile {
    width: 8.13px;
    height: 13.46px;
    top:5.23px;
    left:7.88px;
    gap:32px;
    margin-bottom: 10px;
    color: var(--bs-secondary-color)
  }
  .logo {
    padding:10px 10px 20px 20px;
    margin-left:10px
  }
  @media (max-width: 991px) {
    .logo {
      padding:10px 10px 20px 10px;
    }
   
  }


  </style>
