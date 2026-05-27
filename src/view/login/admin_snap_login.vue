<template>
    <div>
        <el-header>
            <Header />
        </el-header>
        <h1>wait for logining</h1>
    </div>
</template>

<script setup>
import {ref,reactive, onMounted} from 'vue'
import { useRoute } from 'vue-router'
import { adminLogin} from '@/api/profile'
import router from '@/router/index'
import { useUserStore } from '@/pinia/modules/user'
import { ElLoading, ElMessage } from 'element-plus'
import Header from '@/view/layout/header/index.vue'

const userStore = useUserStore()
const route = useRoute()

onMounted(() =>{

    if(route.params.code){
        adminLogin({code:route.params.code}).then(res => {
            if(res.code == 0){
                const token = window.localStorage.getItem('token')
                if (token === res.data) {
                    window.location.replace("/")
                    return
                }
                window.localStorage.setItem('token',res.data)
                userStore.setToken(res.data)
                router.push({ name: 'Login',replace: true })
                window.location.reload()
            }
        })
    }else {
        ElMessage.error('Invalid to Login')
    }
})

</script>
