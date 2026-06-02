<template>
  <div class="login-page">
    <div class="login-aurora login-aurora-one"></div>
    <div class="login-aurora login-aurora-two"></div>

    <header class="login-topbar">
      <a href="/" class="login-brand">
        <img :src="brandLogo" alt="NEWBEECARD" />
        <div>
          <strong>NEWBEECARD</strong>
          <span>Global subscription cards</span>
        </div>
      </a>

      <div class="top-actions">
        <a href="/" class="home-link">{{ $t('lang.to_homepage') }}</a>
        <el-dropdown popper-class="px-2">
          <button class="language-button" type="button">
            {{ langStore.currentLocale === 'zh' ? '中文' : langStore.currentLocale === 'en' ? 'English' : 'Tiếng Việt' }}
            <el-icon><arrow-down /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLangSet('zh')">中文</el-dropdown-item>
              <el-dropdown-item @click="handleLangSet('en')">English</el-dropdown-item>
              <el-dropdown-item @click="handleLangSet('vn')">Tiếng Việt</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <main class="login-shell">
      <section class="login-story" aria-label="NEWBEECARD product overview">
        <p class="login-eyebrow">Secure virtual card console</p>
        <h1>登录虚拟卡控制台</h1>
        <p class="story-copy">
          统一管理虚拟卡、广告预算、AI 工具订阅和团队账单。每一张卡都可以独立设置额度、用途和状态。
        </p>

        <div class="visual-stage">
          <img :src="loginVisual" alt="NEWBEECARD dashboard preview" />
        </div>

        <div class="story-metrics">
          <div>
            <strong>50+</strong>
            <span>卡头资源</span>
          </div>
          <div>
            <strong>1%</strong>
            <span>充值手续费</span>
          </div>
          <div>
            <strong>24/7</strong>
            <span>在线支持</span>
          </div>
        </div>
      </section>

      <section class="auth-panel" aria-label="Authentication forms">
        <div class="panel-header">
          <p class="login-eyebrow">
            {{ formType === 'login' ? 'Access account' : formType === 'register' ? 'Create account' : 'Reset password' }}
          </p>
          <h2 v-if="formType === 'login'">{{ isIamLogin ? $t('lang.login.iam_title') : $t('lang.login.title') }}</h2>
          <h2 v-else-if="formType === 'register'">{{ $t('lang.register.title') }}</h2>
          <h2 v-else>{{ $t('lang.forgot.title') }}</h2>
          <p v-if="formType === 'forgetPassword'" class="panel-description">{{ $t('lang.forgot.description') }}</p>
          <p v-else class="panel-description">使用你的邮箱和密码进入 NEWBEECARD 控制台。</p>
        </div>

        <el-form
          v-if="formType === 'login'"
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-width="0px"
          class="auth-form"
        >
          <el-form-item prop="email" class="auth-form-item">
            <el-input v-model="loginForm.email" :placeholder="$t('lang.login.email')" class="auth-input" />
          </el-form-item>
          <el-form-item prop="password" class="auth-form-item">
            <el-input
              v-model="loginForm.password"
              type="password"
              :placeholder="$t('lang.login.password')"
              class="auth-input"
              show-password
            />
          </el-form-item>
          <div class="form-row">
            <span></span>
            <button class="text-action" type="button" @click="formType = 'forgetPassword'">
              {{ $t('lang.login.forgot_password') }}
            </button>
          </div>
          <el-button class="auth-submit" type="primary" @click="handleLogin">
            {{ $t('lang.login.login_btn') }}
          </el-button>
          <div class="panel-footer">
            <button class="text-action" type="button" @click="isIamLogin = !isIamLogin">
              {{ isIamLogin ? $t('lang.login.switch_to_normal') : $t('lang.login.switch_to_iam') }}
            </button>
            <div v-if="!isIamLogin" class="inline-action">
              <span>{{ $t('lang.login.no_account') }}</span>
              <button class="text-action strong" type="button" @click="formType = 'register'">
                {{ $t('lang.login.register_link') }}
              </button>
            </div>
          </div>
        </el-form>

        <el-form
          v-if="formType === 'register'"
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          label-width="0px"
          class="auth-form"
        >
          <el-form-item prop="email" class="auth-form-item">
            <el-input v-model="registerForm.email" :placeholder="$t('lang.register.email')" class="auth-input" />
          </el-form-item>
          <el-form-item prop="password" class="auth-form-item">
            <el-input v-model="registerForm.password" type="password" :placeholder="$t('lang.register.password')" show-password class="auth-input" />
          </el-form-item>
          <el-form-item prop="repeatPassword" class="auth-form-item">
            <el-input v-model="registerForm.repeatPassword" type="password" :placeholder="$t('lang.register.repeat_password')" show-password class="auth-input" />
          </el-form-item>
          <el-form-item prop="inviteCode" class="auth-form-item">
            <el-input v-model="registerForm.inviteCode" :placeholder="$t('lang.register.invite_code')" class="auth-input" />
          </el-form-item>
          <el-button class="auth-submit" type="primary" @click="handleRegister">
            {{ $t('lang.register.register_btn') }}
          </el-button>
          <div class="panel-footer single">
            <button class="text-action strong" type="button" @click="formType = 'login'">
              {{ $t('lang.register.login_link') }}
            </button>
          </div>
        </el-form>

        <el-form
          v-if="formType === 'forgetPassword'"
          ref="foundPasswordFormRef"
          :model="foundPasswordForm"
          :rules="foundPasswordRules"
          label-width="0px"
          class="auth-form"
        >
          <el-form-item prop="email" class="auth-form-item">
            <el-input v-model="foundPasswordForm.email" :placeholder="$t('lang.forgot.email')" class="auth-input" />
          </el-form-item>
          <el-form-item prop="password" class="auth-form-item">
            <el-input v-model="foundPasswordForm.password" type="password" :placeholder="$t('lang.forgot.new_password')" show-password class="auth-input" />
          </el-form-item>
          <el-form-item prop="newPassword" class="auth-form-item">
            <el-input v-model="foundPasswordForm.newPassword" type="password" :placeholder="$t('lang.forgot.repeat_new_password')" show-password class="auth-input" />
          </el-form-item>
          <el-button class="auth-submit" type="primary" @click="handleResetPassword">
            {{ $t('lang.forgot.confirm_btn') }}
          </el-button>
          <div class="panel-footer single">
            <button class="text-action strong" type="button" @click="formType = 'login'">
              {{ $t('lang.forgot.login_link') }}
            </button>
          </div>
        </el-form>
      </section>
    </main>
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
import brandLogo from '@/assets/NEWBEECARD-logo.png'
import loginVisual from '@/assets/login-card-transparent.png'
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

<style scoped>
.login-page {
  --bg: #020812;
  --panel: rgba(5, 16, 29, 0.86);
  --panel-strong: rgba(4, 13, 24, 0.96);
  --line: rgba(139, 214, 255, 0.16);
  --line-strong: rgba(139, 214, 255, 0.3);
  --text: #f4fbff;
  --muted: rgba(225, 242, 255, 0.74);
  --dim: rgba(225, 242, 255, 0.56);
  --cyan: #44d5ff;
  --blue: #2f7dff;
  --green: #7dffcc;
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  color: var(--text);
  background:
    radial-gradient(circle at 13% 13%, rgba(68, 213, 255, 0.1), transparent 28%),
    radial-gradient(circle at 86% 22%, rgba(47, 125, 255, 0.12), transparent 26%),
    linear-gradient(180deg, #020812 0%, #030b16 52%, #04101d 100%);
  font-family: "Avenir Next", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.login-aurora {
  position: fixed;
  width: 34vw;
  height: 34vw;
  min-width: 320px;
  min-height: 320px;
  border-radius: 999px;
  filter: blur(76px);
  opacity: 0.22;
  pointer-events: none;
}

.login-aurora-one {
  top: 10%;
  left: -10%;
  background: #1fa6d6;
}

.login-aurora-two {
  right: -12%;
  bottom: 8%;
  background: #1d55c9;
}

.login-topbar,
.login-shell {
  position: relative;
  z-index: 1;
}

.login-topbar {
  width: min(1180px, calc(100% - 36px));
  margin: 0 auto;
  padding: 22px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.login-brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: var(--text);
  text-decoration: none;
}

.login-brand img {
  width: 58px;
  height: 58px;
  object-fit: contain;
}

.login-brand strong {
  display: block;
  font-size: 17px;
  letter-spacing: 0.08em;
}

.login-brand span {
  display: block;
  margin-top: 2px;
  color: var(--dim);
  font-size: 12px;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.home-link,
.language-button {
  min-height: 42px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 0 16px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: rgba(4, 15, 29, 0.62);
  text-decoration: none;
  transition: color 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.language-button {
  cursor: pointer;
}

.home-link:hover,
.language-button:hover {
  color: var(--text);
  border-color: var(--line-strong);
  background: rgba(8, 24, 43, 0.82);
}

.login-shell {
  width: min(1180px, calc(100% - 36px));
  min-height: calc(100vh - 104px);
  margin: 0 auto;
  padding: 36px 0 58px;
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(390px, 0.78fr);
  gap: 52px;
  align-items: stretch;
}

.login-story {
  min-width: 0;
  min-height: 430px;
  display: grid;
  grid-template-rows: auto auto 1fr auto;
  align-self: center;
}

.login-eyebrow {
  display: inline-flex;
  align-items: center;
  justify-self: start;
  width: fit-content;
  gap: 10px;
  margin: 0 0 20px;
  padding: 8px 13px;
  border: 1px solid rgba(125, 255, 204, 0.24);
  border-radius: 999px;
  background: rgba(125, 255, 204, 0.08);
  color: var(--green);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.login-eyebrow::before {
  content: "";
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 18px currentColor;
}

.login-story h1 {
  max-width: 720px;
  margin: 0;
  font-size: clamp(2.2rem, 4vw, 3.8rem);
  line-height: 1.08;
  letter-spacing: -0.045em;
}

.story-copy {
  max-width: 620px;
  margin: 14px 0 0;
  color: var(--muted);
  font-size: 15px;
  line-height: 1.7;
}

.visual-stage {
  position: relative;
  min-height: 210px;
  margin-top: 10px;
  display: grid;
  place-items: center;
}

.visual-stage img {
  width: min(82%, 500px);
  max-height: 230px;
  object-fit: contain;
  filter: drop-shadow(0 34px 58px rgba(0, 0, 0, 0.48));
}

.floating-ticket {
  position: absolute;
  min-width: 150px;
  padding: 16px;
  border: 1px solid rgba(139, 214, 255, 0.18);
  border-radius: 20px;
  background: rgba(5, 17, 31, 0.46);
  backdrop-filter: blur(12px);
}

.floating-ticket span {
  color: var(--dim);
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.floating-ticket strong {
  display: block;
  margin-top: 6px;
  font-size: 28px;
  line-height: 1;
}

.ticket-one {
  left: 2%;
  top: 44px;
}

.ticket-two {
  right: 4%;
  bottom: 42px;
}

.story-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  max-width: 560px;
}

.story-metrics div {
  padding: 12px 14px;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: rgba(8, 24, 43, 0.58);
  backdrop-filter: blur(12px);
}

.story-metrics strong {
  display: block;
  color: #ffffff;
  font-size: 21px;
  line-height: 1;
}

.story-metrics span {
  display: block;
  margin-top: 6px;
  color: var(--dim);
  font-size: 12px;
}

.auth-panel {
  padding: 34px;
  border: 1px solid var(--line);
  border-radius: 34px;
  background:
    linear-gradient(140deg, rgba(68, 213, 255, 0.1), transparent 38%),
    rgba(4, 15, 29, 0.82);
  box-shadow: 0 34px 90px rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(18px);
  align-self: center;
}

.panel-header h2 {
  margin: 0;
  font-size: clamp(1.8rem, 3vw, 2.7rem);
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.panel-description {
  margin: 14px 0 0;
  color: var(--muted);
  line-height: 1.7;
}

.auth-form {
  margin-top: 26px;
}

.auth-form-item {
  margin-bottom: 18px;
}

.auth-input :deep(.el-input__wrapper) {
  min-height: 52px;
  border-radius: 16px;
  border: 1px solid rgba(139, 214, 255, 0.14);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: none;
}

.auth-input :deep(.el-input__wrapper.is-focus) {
  border-color: rgba(68, 213, 255, 0.45);
  box-shadow: 0 0 0 3px rgba(68, 213, 255, 0.1);
}

.auth-input :deep(.el-input__inner) {
  color: var(--text);
}

.auth-input :deep(.el-input__inner::placeholder) {
  color: rgba(225, 242, 255, 0.46);
}

.form-row,
.panel-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 14px;
}

.form-row {
  margin: -2px 0 18px;
}

.text-action {
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--muted);
  cursor: pointer;
  font: inherit;
  transition: color 0.2s ease;
}

.text-action:hover,
.text-action.strong {
  color: var(--cyan);
}

.auth-submit {
  width: 100%;
  min-height: 54px;
  border: 0;
  border-radius: 999px;
  color: #04111e;
  font-weight: 800;
  background: linear-gradient(135deg, var(--green), var(--cyan) 58%, #83a8ff);
  box-shadow: 0 22px 48px rgba(68, 213, 255, 0.24);
}

.auth-submit:hover,
.auth-submit:focus {
  color: #04111e;
  background: linear-gradient(135deg, var(--green), var(--cyan) 58%, #83a8ff);
}

.panel-footer {
  margin-top: 18px;
  color: var(--dim);
  font-size: 14px;
}

.panel-footer.single {
  justify-content: center;
}

.inline-action {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

@media (max-width: 980px) {
  .login-page {
    overflow-y: auto;
  }

  .login-shell {
    grid-template-columns: 1fr;
    padding-top: 20px;
  }

  .login-story {
    min-height: auto;
    display: block;
    text-align: left;
  }

  .visual-stage {
    min-height: 300px;
  }
}

@media (max-width: 640px) {
  .login-topbar {
    align-items: flex-start;
    flex-direction: column;
  }

  .top-actions {
    width: 100%;
  }

  .home-link,
  .language-button {
    flex: 1;
    justify-content: center;
  }

  .login-shell {
    width: min(100% - 28px, 1180px);
    gap: 28px;
  }

  .login-story h1 {
    font-size: clamp(2.25rem, 12vw, 3.6rem);
  }

  .story-metrics {
    grid-template-columns: 1fr;
  }

  .visual-stage {
    min-height: 250px;
  }

  .floating-ticket {
    display: none;
  }

  .auth-panel {
    padding: 24px;
    border-radius: 26px;
  }

  .panel-footer {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
