<template>
    <div class="h-100">
        <el-row class="h-100">
          <el-col :span="8">
            <div class="position-relative d-flex" style="height: 60px;top:6px;left:10px">
            <i class="bi bi-list h1 me-2 align-self-center d-lg-none" @click="toggleSideBar"></i>
            <img :src="logo" class="logo"></img>
          </div>
          </el-col>
              <el-col  :span="6" :offset="10" class=" align-items-center justify-content-end d-flex">
                <i class="bi bi-envelope fs-5 me-4"></i>
                <!-- <i class="bi bi-bell fs-3 me-2"></i> -->
                <!-- <i class="bi bi-question-circle fs-5 me-3"></i> -->
                <el-dropdown class="me-2">
                  <span class="el-dropdown-link">
                    <el-avatar shape="square" :size="35" :src="form.headerImg"></el-avatar><el-icon class="fs-4 el-icon--right"><arrow-down /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="visibleDialog = true" class="text-center d-flex justify-content-center">详情</el-dropdown-item>
                      <el-dropdown-item ><button class="btn btn-link  text-nowrap" @click="loginout">退出</button></el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-col>
                
        </el-row>
        <el-dialog
          title="个人信息"
          width="65%"
          v-model="visibleDialog"

        >
        <div>
          <div class="row">
            <div class="col-4">
              <div>
                <img class="object-fit-scale w-100 rounded-1" :src="form.headerImg">
                <div class="d-flex justify-content-center">
                    <p class="fs-4 fw-bold">{{ form.userName }}</p>
                </div>
                <div class="d-flex justify-content-center">
                    <div class="fs-4">{{ form.email }}</div>
                </div>
                
              </div>
            </div>
            <div class="col-8">
              <el-form  label-position="right" label-width="90px">
                <el-form-item label="头像" style="max-width: 200px;">
                  <el-select v-model="form.headerImg">
                    <template #label>
                      <img style="scale: 0.42;" class="object-fit-scale w-100 h-100" :src="form.headerImg"></img>
                    </template>
                    <el-option v-for="item in headimgs" :value="item">
                      
                        <img class="object-fit-scale w-100 h-100"  :src="item"></img>
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="昵称">
                  <el-input v-model="form.nickName"></el-input>
                </el-form-item>
                <el-form-item label="邀请码">
                  <el-input disabled v-model="form.inviteCode"></el-input>
                </el-form-item>
                <el-form-item label="密码">
                  <el-input v-model="form.password"></el-input>
                </el-form-item>
                <el-form-item label="新密码">
                  <el-input v-model="form.newPassword"></el-input>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>
        <template #footer>
          <button class="btn btn-danger me-2" @click="handleSetSelfInfo">保存信息</button>
          <button class="btn btn-primary" @click="handleChangePassword">修改密码</button>
        </template>
        </el-dialog>
      </div>
  </template>
  
  <script>
export default {
  name: 'Header',
}
</script>
  <script setup>
  import {ref,computed, reactive,onMounted} from 'vue'
  import {ElMessage} from 'element-plus'
  import logopng from '@/assets/logo.png'
  import { useUserStore } from '@/pinia/modules/user'
  import { setSelfInfo, changePassword } from '@/api/user.js'

  const userStore = useUserStore()
  const visibleDialog = ref(false)
  const form = reactive({
    ...userStore.userInfo
  })
  const loginout = ()=>{
    userStore.LoginOut()
}
const handleSetSelfInfo = ()=>{
  setSelfInfo(form).then(res=>{
    if(res.code==0){
      ElMessage.success("修改成功")
      userStore.GetUserInfo()
    }
  })
}
const handleChangePassword =()=>{
  var formData = {
      password:form.password,
      newPassword:form.newPassword,
      username:form.userName
  }
  changePassword(formData).then((res) => {
        if (res.code === 0) {
          ElMessage.success('修改密码成功！')
        }
        showPassword.value = false
      })
}
onMounted(()=>{
})
const headimgs = [
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_1.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_2.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_3.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_4.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_5.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_6.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_7.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_8.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_9.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_10.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_11.jpeg",
  "https://metalposterpro.s3.us-east-1.amazonaws.com/static/headimg/border_12.jpeg",
] 
const emit = defineEmits(['expland']);
const isMobile = computed(()=>{
  return window.innerWidth < 768
})
const toggleSideBar = () => {
  emit('expland')
}
  // No script needed in this example.
  const logo = ref(logopng)
  </script>
  
  <style scoped>
  .header {
    padding: 0;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  }
  
  .right-icons {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }
  .logo {
      align-self: center;
      height: 60px; padding-bottom:20px
    }
  @media (max-width: 1000px) {
    .logo {
      align-self: center;
      height: 40px; padding-bottom:10px
    }
  }
  
  </style>