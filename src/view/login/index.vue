<template>
    <div class="w-100 h-100">
        <el-header class="bg-white" style="min-height: 80px;">
            <Header />
        </el-header>
        <el-row class="w-100 h-100 justify-content-center align-items-center">
            <el-col :xs="24" :sm="24" :md="24" class="justify-content-center align-items-center w-100 " style="max-width: 504px;padding:0 20px">
                <div class="bg-white rounded-2 login-form ">
                    <h1 class="fs-1 d-flex justify-content-center align-items-center login-title"><span>Login</span></h1>
                <div>
                    <div class="mb-3" style="height: 64px;">
                        <div class="mb-1">Email Address</div>
                        <el-input v-model="form.email" required></el-input>
                    </div>
                    <div style="height: 64px;">
                        <div  class="mb-1">Password</div>
                        <el-input v-model="form.password" type="password" required clearable autocomplete="new-password"></el-input>
                    </div>
                    <div class="verify">
                        <Vcode ref="vcode" :show="isShow" type="inside" :imgs="imgs" :sliderSize="40" :canvasWidth="canvasWidth" successText="success" failText="failed" sliderText="drag to completed" @success="vertifySuccess" />
                    </div>
                    <div>
                        <button :disabled="!hasVerify" @click="login" class="w-100 text-center mt-3 btn-danger btn">Login</button>
                        <div class="w-100 justify-content-end d-flex">
                            <button  @click="forgetPassword" class="text-danger btn btn-link text-decoration-none">forget your password?</button>
                        </div>
                    </div>
                </div>
            </div>
            </el-col>
        </el-row>
    </div>
</template>

<script setup>
import {ref,reactive,computed,onMounted} from 'vue'
import { useRouter } from 'vue-router'
import router from '@/router/index'
import { useUserStore } from '@/pinia/modules/user'
import Header from '@/view/layout/header/index.vue'
import { emitter } from '@/utils/bus.js'

// import Vcode from "./verify_code.vue";
import Vcode from "vue3-puzzle-vcode";

import img1 from '@/assets/imgs/11.jpeg'
import img2 from '@/assets/imgs/12.jpg'
import { ElMessage } from 'element-plus'
import {captcha} from '@/api/profile'
const form = reactive({
    email : '',
    password : '',
    captchaId: '',
    captcha:0
})
const hasVerify = ref(false)
const canvasWidth = computed(()=>{
    return document.body.clientWidth <504? document.body.clientWidth-100:404
})
const sliderVConf = reactive({
    imgUrl:'https://i.pinimg.com/originals/e7/e7/c3/e7e7c37d49d40b043029a147e1158843.jpg',
    isShowSelf: true,
        width: 300,
        height: 180,
        sText: 'sText',
        eText: 'eText',
        isBorder: true,
        isCloseBtn: true,
        isReloadBtn: true,
        isParentNode: false,
        isShowTip: true,
})
const imgs=[img1,img2]
const isShow = ref(true)
const vertifySuccess = (value,obj) =>{
    captcha().then(res =>{
        if(res.code === 0 ){
            form.captcha = res.data.captcha
            form.captchaId = res.data.id
            form.captcha = value*form.captcha
        }
    })
    hasVerify.value = true
    document.getElementsByClassName('auth-body_')[0].style.display = 'none'
}
const sleep = (ms) =>{
    return new Promise(resolve => setTimeout(resolve, ms))
}
const userStore = useUserStore()
const vcode = ref()
const login = async() =>{
    const success = await userStore.LoginIn(form)
    if(!success){
        vcode.value.reset()
        hasVerify.value = false
        document.getElementsByClassName('auth-body_')[0].style.display = ''
    }
}
const forgetPassword = () =>{
    router.push({name:'forgetPassword'})
}
const isMobile = ref(false)

emitter.on('mobile', (item) => {
  isMobile.value = item
})
</script>
<style scoped>
.login-form {

    /* top: 226px; */
    /* left: 589px; */
    padding: 0px 30px;
    gap: 7px;
    border-radius: 15px 0px 0px 0px;
    opacity: 0px;
}
    .login-title {
        
        font-size: 40px;
        font-weight: 700;
        line-height: 48.41px;
        text-align: center;
        gap: 0px;
        opacity: 0px;

    }


</style>
<style>
.auth-body_ {
    display: none;
}
.auth-body_:hover {
    display: none;
    z-index: 2;
    
}
.verify:hover .auth-body_{
    display: block;
    position: absolute;
    top: -410%;
    z-index: 2;
}
/* .range-box:hover + .auth-body_
.range-box:hover ~ .auth-body_{
    display: block;
} */
</style>