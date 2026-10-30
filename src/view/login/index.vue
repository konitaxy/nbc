<template>
  <div id="userLayout">
    <div class="login_panle">
      <div class="login_panle_form">
        <div class="login_panle_form_title">
          <img
            class="login_panle_form_title_logo"
            :src="$GIN_VUE_ADMIN.appLogo"
            alt
          >
          <p class="login_panle_form_title_p">{{ $GIN_VUE_ADMIN.appName }}</p>
        </div>
        <el-form
          v-if="step === 'password'"
          ref="loginForm"
          :model="loginFormData"
          :rules="rules"
          @keyup.enter="submitForm"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginFormData.username"
              placeholder="请输入用户名"
            >
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <user />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginFormData.password"
              :type="lock === 'lock' ? 'password' : 'text'"
              placeholder="请输入密码"
            >
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <component
                      :is="lock"
                      @click="changeLock"
                    />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="vPicBox">
              <el-input
                v-model="loginFormData.captcha"
                placeholder="请输入验证码"
                style="width: 60%"
              />
              <div class="vPic">
                <img
                  v-if="picPath"
                  :src="picPath"
                  alt="请输入验证码"
                  @click="loginVerify()"
                >
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <!-- <el-button
              type="primary"
              style="width: 46%"
              size="large"
              @click="checkInit"
            >前往初始化</el-button> -->
            <el-button
              type="primary"
              size="large"
              style="width: 46%;margin-top: 12px;"
              @click="submitForm"
            >登 录</el-button>
          </el-form-item>
        </el-form>
        <div v-else class="twofa-box">
          <p class="twofa-title">{{ step === 'bind' ? '绑定两步验证' : '两步验证' }}</p>
          <p class="twofa-desc">
            {{ step === 'bind' ? '使用验证器扫描二维码，然后输入 6 位动态码完成绑定。' : '请输入验证器中的 6 位动态码。' }}
          </p>
          <img v-if="step === 'bind'" :src="qrCode" alt="2FA QR code">
          <p v-if="step === 'bind'" class="secret">{{ secret }}</p>
          <el-input v-model="otpCode" maxlength="6" placeholder="6 位动态码" @keyup.enter="submitOtp" />
          <div class="twofa-actions">
            <el-button size="large" @click="backToPassword">返回</el-button>
            <el-button type="primary" size="large" @click="submitOtp">确认</el-button>
          </div>
        </div>
      </div>
      <div class="login_panle_right" />
      <div class="login_panle_foot">
        <div class="copyright">
          <bootomInfo />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
}
</script>

<script setup>
import { bindAdmin2FA, captcha, verifyAdmin2FA } from '@/api/user'
import { checkDB } from '@/api/initdb'
import bootomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
const router = useRouter()
// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}

// 获取验证码
const loginVerify = () => {
  captcha({}).then((ele) => {
    rules.captcha[1].max = ele.data.captchaLength
    rules.captcha[1].min = ele.data.captchaLength
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
  })
}
loginVerify()

// 登录相关操作
const lock = ref('lock')
const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}

const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: '',
  password: '',
  captcha: '',
  captchaId: '',
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const step = ref('password')
const bindTicket = ref('')
const loginTicket = ref('')
const qrCode = ref('')
const secret = ref('')
const otpCode = ref('')

const submitForm = () => {
  loginForm.value.validate(async(v) => {
    if (!v) {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      loginVerify()
      return false
    }
    const result = await userStore.LoginIn(loginFormData)
    if (result?.needBind2FA) {
      step.value = 'bind'
      bindTicket.value = result.bindTicket
      qrCode.value = result.qrCode
      secret.value = result.secret
      otpCode.value = ''
      return
    }
    if (result?.need2FA) {
      step.value = 'otp'
      loginTicket.value = result.loginTicket
      otpCode.value = ''
      return
    }
    if (!result) {
      loginVerify()
    }
  })
}

const backToPassword = () => {
  step.value = 'password'
  otpCode.value = ''
  loginVerify()
}

const submitOtp = async() => {
  const code = otpCode.value.trim()
  if (!/^\d{6}$/.test(code)) {
    ElMessage({ type: 'error', message: '请输入 6 位动态码', showClose: true })
    return
  }
  const res = step.value === 'bind'
    ? await bindAdmin2FA({ bindTicket: bindTicket.value, code })
    : await verifyAdmin2FA({ loginTicket: loginTicket.value, code })
  if (res.code === 0) {
    await userStore.completeLogin(res.data)
  }
}

// 跳转初始化
const checkInit = async() => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: '已配置数据库信息，无法初始化',
      })
    }
  }
}

</script>

<style lang="scss" scoped>
@import "@/style/newLogin.scss";
</style>
