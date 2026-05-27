<template>
    <div class="h-100">
        <el-header  class="bg-white">
            <Header />
        </el-header>
        <el-row v-if="visible" class="w-100 h-100 justify-content-center align-items-center">
            <el-col :xs="24" :sm="24" :md="12" class="d-flex justify-content-center align-items-center" >
                <div class="bg-white rounded-2 login-form">
                <!-- <h1 class="text-center fs-2">think you for join us</h1> -->
                <h3 class="text-center fs-2 login-title  d-flex justify-content-center align-items-center ">Set your account password</h3>
                <el-form class="bg-white container rounded-4 p-4">
                    <div class="mb-1">
                        <div class="mb-1">Email Address</div>
                        <el-input :disabled="form.disable" v-model="form.email" required></el-input>
                    </div>
                    <div>
                        <div class=" w-100 align-items-end justify-content-between d-flex"><span  class="mb-1">Received Code</span><button :class="form.intervalID === null?'text-danger btn btn-link':'text-secondary btn btn-link'" @click.prevent="startSendEmail"><span>{{ form.sendText }}</span></button></div>
                        <el-input v-model="form.verifyCode"></el-input>
                    </div>
                    <div class="mb-1">
                        <div class="mb-1">Password</div>
                        <el-form-item label="Set New Password">
                            <el-tooltip trigger="click" effect="light" placement="bottom-start">
                                <template #content>
                                    <div>
                                        <span class="fw-medium">Password must include:</span>
                                        <div class="d-flex text-nowrap"><i :class="form.password.length>7&&form.password.length<20?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>8-20 characters</p></div>
                                        <div class="d-flex text-nowrap"><i :class="containsLetterRegex.test(form.password)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>At least one capital letter</p></div>
                                        <div class="d-flex text-nowrap"><i :class="containsNumberRegex.test(form.password)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>At least one number</p></div>
                                        <div class="d-flex text-nowrap"><i :class="!containsSpaceRegex.test(form.password)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>No spaces</p></div>
                                    </div>
                                </template>
                        <el-input type="password" v-model="form.password" :minlength="6" show-password clearable autocomplete="new-password"></el-input>
                        </el-tooltip>
                        </el-form-item>
                    </div>
                    <div class="mb-4">
                        <div  class="mb-1">Repeat Password</div>
                        <el-form-item label="Set New Password">
                            <el-tooltip trigger="click" effect="light" placement="bottom-start">
                                <template #content>
                                    <div>
                                        <span class="fw-medium">Password must include:</span>
                                        <div class="d-flex text-nowrap"><i :class="form.newPassword.length>7&&form.newPassword.length<20?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>8-20 characters</p></div>
                                        <div class="d-flex text-nowrap"><i :class="containsLetterRegex.test(form.newPassword)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>At least one capital letter</p></div>
                                        <div class="d-flex text-nowrap"><i :class="containsNumberRegex.test(form.newPassword)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>At least one number</p></div>
                                        <div class="d-flex text-nowrap"><i :class="!containsSpaceRegex.test(form.newPassword)?'bi bi-check text-success fs-4 me-1':'bi bi-x text-danger fs-4 me-1'"></i><p>No spaces</p></div>
                                    </div>
                                </template>
                        <el-input type="password" v-model="form.newPassword" :minlength="6" show-password clearable autocomplete="new-password"></el-input>
                            </el-tooltip>
                            </el-form-item>
                    </div>
                    <div>
                        <el-button @click="onSubmitRegister" class="w-100 text-center mt-1 btn btn-danger">Confirm</el-button>
                        <button class="btn btn-link text-danger text-decoration-none text-end w-100" @click="toTarget">Login</button>
                    </div>

                </el-form>
            </div>
            </el-col>
        </el-row>
        <el-row v-else class="w-100  justify-content-center">
            <div class="bg-white rounded-2 justify-content-center text-center p-5 mt-5">
                <h1 class="text-center fs-2">Congratulations！</h1>
                <h3 class="text-center fs-6 w-100 d-flex justify-content-center mt-3"><span class="w-50">Password reset successful, please start your show</span></h3>
                <el-button @click="toTarget" class="w-50 text-center mt-3 btn btn-danger">Login Now</el-button>

            </div>
        </el-row>
    </div>
</template>
<script setup>
import {ref,reactive, onMounted} from 'vue'
import { useRoute } from 'vue-router'
import { sendVerifyCode,resetPassword} from '@/api/profile'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import Header from '@/view/layout/header/index.vue'


const route = useRoute()
const visible = ref(true)
const verificationFree = ref(false)
const containsLetterRegex = /[a-zA-Z]/
const containsNumberRegex = /[0-9]/
const containsSpaceRegex = /\s/
const form = reactive({
    email:'',
    password: '',
    newPassword: '',
    verifyCode: '',
    disable: false,
    sendText: 'Send verification code',
    intervalID:null,
    countdown: 60,
    lastDateTime: Date.now()
})
const registerVerofyCodeForm =reactive({
    to : '',
    type: 'resetPassword'
})
const startSendEmail = async() =>{
    if (form.intervalID !== null) {
      return
    }
    form.intervalID = -1
    registerVerofyCodeForm.to = form.email

    sendVerifyCode(registerVerofyCodeForm).then(res => {
      if(res.code === 0){
        ElMessage.success('Send success')
        form.lastTime = Date.now()
        form.intervalID  = setInterval(() => {
        const now = Date.now()
        const elapsed = now - form.lastTime;
        form.lastTime = Date.now()
        form.countdown -= Math.floor(elapsed / 1000);
        form.sendText =  "Resend verification code "+form.countdown+""
        if (form.countdown <= 0) {
            clearInterval(form.intervalID);
            form.intervalID = null
            form.countdown = 60
            form.sendText = "Send verification code"
        }
        }, 1000);
      }else {
        form.intervalID = null
      }
    })
    
}
const onSubmitRegister = () =>{
    if(form.password !== form.newPassword){
        ElMessage.error('Password is not same')
        return
    }else {
        resetPassword(form).then(res => {
            if (res.code === 0) {
                visible.value = false
            }
        })
    }
    
}

onMounted(() =>{
    const feeVerify = window.localStorage.getItem('verificationFree')
    if(feeVerify == null){
        verificationFree.value = false
    }else {
        var date = new Date().setDate(new Date().getDate() + 30)
        if (date > new Date(feeVerify)){
            verificationFree.value = false
            window.localStorage.removeItem('verificationFree')
        }else {
            verificationFree.value = true
        }
    }
})

const toTarget = () =>{
    router.push({name:'Login',replace:true})
}

</script>
<style scoped>
.login-form {
    width: Fixed (604px)px;
    height: Fixed (449px)px;
    padding: 0px 50px 0px 50px;
    gap: 7px;
    border-radius: 15px 0px 0px 0px;
    opacity: 0px;
}
    .login-title {
        
        font-size: 40px;
        font-weight: 700;
        line-height: 48.41px;
        text-align: center;
        width: 504px;
        height: 108px;
        gap: 0px;
        opacity: 0px;

    }


</style>