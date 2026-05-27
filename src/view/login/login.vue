<template>
  <div class="container-fluid">
    <div class="position-absolute top-0 p-3" style="right: 30px;">
      <a href="/" class="text-success me-5" >{{ $t('lang.to_homepage') }}</a>
      <el-dropdown class="me-2" popper-class="px-2">
        <a href="#" class="link-success link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover" style="line-height: normal;">
          {{ langStore.currentLocale ==='zh'?'中文':langStore.currentLocale ==='en'?'English':'Tiếng Việt' }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
        </a>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="handleLangSet('zh')" class="">中文</el-dropdown-item>
            <el-dropdown-item @click="handleLangSet('en')" class="">English</el-dropdown-item>
            <el-dropdown-item @click="handleLangSet('vn')" class="">Tiếng Việt</el-dropdown-item>

          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <div class="row">
      <!-- 左侧宣传区域 -->
      <div class="col col-6 vh100  d-none d-sm-block p-0">
        <div class="left-container">
          <div class="logo text-center">
            <img src="@/assets/logo-black.png" alt="Logo" class="logo">
          </div>
          <div class="slogan text-left">
            <!-- 使用 $t 进行翻译 -->
            <h1 class="mb-2 ">{{ $t('lang.login.slogan.line1') }}</h1>
            <h1 class="mb-2 text-highlight">{{ $t('lang.login.slogan.line2') }}</h1>
            <h1 class="mb-2 ">{{ $t('lang.login.slogan.line3') }}</h1>
            <h1 class="mb-2 text-highlight">{{ $t('lang.login.slogan.line4') }}</h1>
            <div class="left-main-image text-center">
              <img src="@/assets/cards_login_new_1.png" alt="Main" class="main-image">
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧表单区域 -->
      <div class="col-xs-12 col-sm-6 vh100 p-0">
        <div class="right-side">
          <!-- 登录表单 -->
          <el-card class="login-card login-form" shadow="never" v-if="formType == 'login'">
            <!-- 标题 -->
            <h2 class="mb-3 fs-3">{{ isIamLogin ? $t('lang.login.iam_title') : $t('lang.login.title') }}</h2>
            <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="0px">
              <el-form-item prop="email" class="login-form-item">
                <el-input v-model="loginForm.email" 
                          :placeholder="$t('lang.login.email')" 
                          class="login-form-input">
                </el-input>
              </el-form-item>
              <el-form-item prop="password" class="login-form-item">
                <el-input v-model="loginForm.password" 
                          type="password" 
                          :placeholder="$t('lang.login.password')" 
                          class="login-form-input" 
                          show-password>
                </el-input>
              </el-form-item>
              <el-form-item>
                <div class="justify-content-between w-100 d-flex">
                  <div>
                    <!-- <el-checkbox v-if="!verificationFree" 
                                 v-model="loginForm.remember" 
                                 style="font-size: 12px;">
                      {{ $t('lang.login.remember_me') }}
                    </el-checkbox> -->
                  </div>
                  <div>
                    <!-- 链接文本 -->
                    <!-- <el-link @click="clear()" type="danger" style="float: right; margin-left: 10px;font-size: 12px">
                      {{ $t('lang.login.clear_cache') }}
                    </el-link> -->
                    <el-link type="info" style="float: right;font-size: 12px;" @click="formType = 'forgetPassword'">
                      {{ $t('lang.login.forgot_password') }}
                    </el-link>
                  </div>
                </div>
              </el-form-item>
              <el-form-item>
                <el-button class="w-100 login-btn" type="primary" @click="handleLogin">
                  {{ $t('lang.login.login_btn') }}
                </el-button>
                <div class="w-100" style="display: flex;justify-content: space-between; align-items: center;">
                  <!-- IAM 登录切换 -->
                  <el-link type="success" @click="isIamLogin = !isIamLogin">
                    {{ isIamLogin ? $t('lang.login.switch_to_normal') : $t('lang.login.switch_to_iam') }}
                  </el-link>
                  <!-- 注册链接 -->
                  <div v-if="!isIamLogin" style="display: flex; gap:5px;">
                    <span>{{ $t('lang.login.no_account') }}</span> 
                    <el-link type="primary" @click="formType = 'register'">
                    {{ $t('lang.login.register_link') }}
                  </el-link>
                  </div>
                </div>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 注册表单 -->
          <el-card class="login-card" shadow="never" v-if="formType == 'register'">
            <h2 class="mb-3 fs-5">{{ $t('lang.register.title') }}</h2>
            <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" label-width="0px">
              <el-form-item prop="email" class="login-form-item">
                <el-input v-model="registerForm.email" 
                          :placeholder="$t('lang.register.email')"
                          class="login-form-input" >
                </el-input>
              </el-form-item>
              <el-form-item prop="password" class="login-form-item">
                <el-input v-model="registerForm.password" 
                          type="password" 
                          :placeholder="$t('lang.register.password')" 
                          show-password
                          class="login-form-input" >
                </el-input>
              </el-form-item>
              <el-form-item prop="repeatPassword" class="login-form-item">
                <el-input v-model="registerForm.repeatPassword" 
                          type="password" 
                          :placeholder="$t('lang.register.repeat_password')" 
                          show-password
                          class="login-form-input" >
                </el-input>
              </el-form-item>
              <el-form-item prop="inviteCode" class="login-form-item">
                <el-input v-model="registerForm.inviteCode" 
                          :placeholder="$t('lang.register.invite_code')"
                          class="login-form-input" >
                </el-input>
              </el-form-item>
              <el-form-item>
                <el-button class="w-100" type="primary" @click="handleRegister">
                  {{ $t('lang.register.register_btn') }}
                </el-button>
                <div class="w-100">
                  <el-link class="primary" style="color: #01ad5a;"  @click="formType = 'login'">
                    {{ $t('lang.register.login_link') }}
                  </el-link>
                </div>
              </el-form-item>
            </el-form>
          </el-card>

          <el-card class="login-card" v-if="formType == 'forgetPassword'">
            <h2 class="mb-3 fs-5">{{ $t('lang.forgot.title') }}</h2>
            <p class="text-secondary" style="font-size: small;">{{ $t('lang.forgot.description') }}</p>
            <el-form ref="foundPasswordFormRef" :model="foundPasswordForm" :rules="foundPasswordRules" label-width="0px">
              <el-form-item prop="email">
                <el-input v-model="foundPasswordForm.email" 
                          :placeholder="$t('lang.forgot.email')">
                </el-input>
              </el-form-item>
              <el-form-item prop="password">
                <el-input v-model="foundPasswordForm.password" 
                          type="password" 
                          :placeholder="$t('lang.forgot.new_password')" 
                          show-password>
                </el-input>
              </el-form-item>
              <el-form-item prop="newPassword">
                <el-input v-model="foundPasswordForm.newPassword" 
                          type="password" 
                          :placeholder="$t('lang.forgot.repeat_new_password')" 
                          show-password>
                </el-input>
              </el-form-item>
              <el-form-item>
                <el-button class="w-100" type="primary" @click="handleResetPassword">
                  {{ $t('lang.forgot.confirm_btn') }}
                </el-button>
                <div class="w-100">
                  <!-- 登录链接 -->
                  <el-link class="primary"  @click="formType = 'login'">
                    {{ $t('lang.forgot.login_link') }}
                  </el-link>
                </div>
              </el-form-item>
            </el-form>
          </el-card>

        </div>
      </div>
    </div>
  </div>
</template>
  
<script setup>
import { ElMessage } from 'element-plus';
import { reactive, ref, onMounted } from 'vue';
import { login,sendVerifyCode,register,resetPassword} from '@/api/profile';
import { iamLogin } from '@/api/iam';
import { useUserStore } from '@/pinia/modules/user';
import { validateEmail, validatePassword, validatePhone } from '@/utils/validates';
import {formatDateYYYYMMDD,getCookie,setCookie,deleteCookie} from '@/utils/format';
import { useLanguageStore } from '@/pinia/modules/language';
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
const langStore = useLanguageStore()

const verificationFree = ref(false)
const isIamLogin = ref(false)
const loginFormRef = ref()
const registerFormRef = ref()
const foundPasswordFormRef = ref()
const loginRules = reactive({
  email: [
    { 
      required: true, 
      message: t('lang.validation.email_required'), 
      trigger: 'blur' 
    },
    { 
      validator: validateEmail, 
      message: t('lang.validation.email_invalid'), 
      trigger: 'blur' 
    },
  ],
})

const formType = ref('login')
const loginForm = reactive({
    email: '',
    password: '',
    authMethod: 2,
    phone: '',
    code: '',
    verifyCode: '',
    remember: false,
})
const foundPasswordForm = reactive({
  email: '',
  password: '',
  newPassword: '',
})
const registerForm = reactive({
    email: '',
    repeatPassword:'',
    password: '',
    phone: '',
    verifyCode: '',
    inviteCode:'',
})

const handleLangSet = (lang) =>{
  langStore.changeLocale(lang)
}
onMounted(()=>{
  
  const freeVerify = getCookie('verification_free')
  if(freeVerify){
    verificationFree.value = true
  }
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== registerForm.password) {
    // 使用 $t 函数获取国际化错误消息
    callback(new Error(t('lang.validation.password_mismatch')));
  } else {
    callback();
  }
}

const validateConfirmPassword2 = (rule, value, callback) => {
  if (value !== foundPasswordForm.password) {
    // 使用 $t 函数获取国际化错误消息
    callback(new Error(t('lang.validation.password_mismatch')));
  } else {
    callback();
  }
}

// ... (foundPasswordFormRef 定义保持不变)

// 忘记密码表单规则
const foundPasswordRules = reactive({
  email: [
    { validator: validateEmail, trigger: 'blur' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  newPassword: [
    { validator: validateConfirmPassword2, trigger: 'blur' }
  ]
})

// 注册表单规则
const registerRules = reactive({
  name: [
    { 
      required: true, 
      // 国际化必填消息
      message: t('lang.validation.name_required'), 
      trigger: 'blur' 
    },
    { 
      min: 6, 
      max: 20, 
      // 国际化长度消息
      message: t('lang.validation.name_length'), 
      trigger: 'blur' 
    },
    { 
      whitespace: true, 
      // 国际化空格消息
      message: t('lang.validation.name_whitespace'), 
      trigger: 'blur' 
    }
  ],
  email: [
    { validator: validateEmail, trigger: 'blur' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  repeatPassword: [
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  inviteCode: [
    { 
      required: true, 
      // 国际化必填消息
      message: t('lang.validation.invite_code_required'), 
      trigger: 'blur' 
    },
    { 
      len: 6, 
      // 国际化长度消息
      message: t('lang.validation.invite_code_length'), 
      trigger: 'blur' 
    },
  ],
  verifyCode: [
    { 
      required: true, 
      // 国际化必填消息
      message: t('lang.validation.verify_code_required'), 
      trigger: 'blur' 
    },
    { 
      len: 6, 
      // 国际化长度消息
      message: t('lang.validation.verify_code_length'), 
      trigger: 'blur' 
    },
  ]
})
const userStore = useUserStore()
const handleLogin = async() =>{
  const r = await loginFormRef.value.validate()
  if(!r) return
  
  if (isIamLogin.value) {
    // IAM 登录
    userStore.IamLoginIn(loginForm).then((success)=>{
      if(success && loginForm.remember){
        setCookie('verification_free',"yes",30)
      }
    })
  } else {
    // 普通登录
  userStore.LoginIn(loginForm).then((success)=>{
    if(success && loginForm.remember){
      setCookie('verification_free',"yes",30)
    }
  })
  }
}
const handleResetPassword = async() =>{
  const r =await foundPasswordFormRef.value.validate()
    if(!r){
      return
    }
  resetPassword(foundPasswordForm).then(res =>{
    if(res.code === 0){
      ElMessage.success(t('lang.success'))
      formType.value = 'login'
    }
  })
}
const handleRegister =async () =>{
    const r =await registerFormRef.value.validate()
    if(!r){
      return
    }
    register(registerForm).then(res => {
      if (res.code === 0) {
        formType.value = 'login'
        ElMessage.success(t('lang.success'));
      }
    })
}
const handleGetCode =(t) => {
  if(t == 1){
    startSendEmail(loginForm.email,'loginVerify')
  }else {
    startSendEmail(registerForm.email,'registerVerify')
  }
}

const resetForm = reactive({
  email:'',
  password: '',
  newPassword: '',
  verifyCode: '',
  disable: false,
  sendText: 'Send code',
  intervalID:null,
  countdown: 60,
  lastDateTime: Date.now()
})
const startSendEmail = async(toMail,mailType) =>{
    if (resetForm.intervalID !== null) {
      return
    }
    resetForm.intervalID = -1
    sendVerifyCode({
      to: toMail,
      type:mailType
    }).then(res => {
      if(res.code === 0){
        ElMessage.success('Send success')
        resetForm.lastTime = Date.now()
        resetForm.intervalID  = setInterval(() => {
        const now = Date.now()
        const elapsed = now - resetForm.lastTime;
        resetForm.lastTime = Date.now()
        resetForm.countdown -= Math.floor(elapsed / 1000);
        resetForm.sendText =  "Resend code "+resetForm.countdown+""
        if (resetForm.countdown <= 0) {
            clearInterval(resetForm.intervalID);
            resetForm.intervalID = null
            resetForm.countdown = 60
            resetForm.sendText = "Send code"
        }
        }, 1000);
      }else {
        resetForm.intervalID = null
      }
}) 

}
const clear = ()=>{
  deleteCookie('verification_free')
  window.location.reload()
} 
</script>
  <style lang="scss" root>
  @use "@/style/login.scss";
  </style>