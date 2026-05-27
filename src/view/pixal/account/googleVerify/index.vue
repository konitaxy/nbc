<template>
    <div class="p-2">
      <div class="row mt-2">
        <div class="col">
          <h4>{{ $t('lang.google_authentication') }}</h4>
        </div>
      </div>
      <div class="mt-4 w-100">
      <el-alert v-if="userStore.userInfo.bind2FA" type="success" :closable="false">{{ $t('lang.bind_2fa_success_message') }}</el-alert>
      <el-alert v-else type="error" :closable="false">{{ $t('lang.bind_2fa_required_message') }}</el-alert>
    </div>
      <div class="row mt-5">
        <div class="col">
          
          <el-form label-width="100px">
            <el-form-item v-if="!userStore.userInfo.bind2FA" :label="$t('lang.secret_key')">
              <el-input :type="startToBind?'text':'password'" disabled class="w-50" v-model="data.secret" ></el-input>
            </el-form-item>
            <el-form-item v-if="!userStore.userInfo.bind2FA" :label="$t('lang.scan_bind')" class="align-items-start">
              <div v-if="!startToBind" style="width: 200px; height: 200px;background-color: white;"></div>
              <img v-else :src="data.qrCode" alt="QR Code" style="width: 200px; height: 200px;">
              <div class="mt-2 help-text">
                <ul>
                  <li>{{ $t('lang.start_binding') }}:</li>
                  <li>1.{{ $t('lang.open_your_authenticator_app') }}</li>
                  <li>2.{{ $t('lang.add_otp_in_the_app') }}</li>
                  <li>3.{{ $t('lang.scan_qr_code') }}</li>
                </ul>
              </div>
            </el-form-item>
            <el-form-item>  
              <el-button v-if="userStore.userInfo.bind2FA" type="danger" @click="handleDisableTOCP">{{ $t('lang.disable_secret_key') }}</el-button>
              <el-button v-if="!userStore.userInfo.bind2FA&&!startToBind" type="primary" style="width: auto;" @click="startToBind = true">{{ $t('lang.start_bind') }}</el-button>
              <el-button v-if="!userStore.userInfo.bind2FA" type="primary" style="width: auto;" @click="handleConfirmBind">{{ $t('lang.next_step') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { reactive, ref , onMounted} from 'vue';
  import {genTOCPSecret,confirmTOCPBind,disableTOCPBind} from '@/api/client';
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useUserStore } from '@/pinia/modules/user'
  const userStore = useUserStore()
  const data = ref({
    secret:'',
    qrCode:''
  });
  const startToBind = ref(false)
  
  const getSecretKey = () => {
      genTOCPSecret().then(res =>{
        if(res.code === 0){
          data.value = res.data
        }
      })
  };

  const handleConfirmBind = () =>{
    ElMessageBox.prompt('Please input 2fa code to completed bind', '2FA Code', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputErrorMessage: 'Invalid Code',
  }).then(({ value }) => {
      confirmTOCPBind({verifyCode:value}).then(res => {
        if(res.code === 0){
          ElMessage.success('Success')
          userStore.GetUserInfo()
        }
      })
    })
    .catch(() => {
    })
  }
  onMounted(() => {
    if(!userStore.userInfo.bind2FA) {
      getSecretKey()
    }
  })
  const handleDisableTOCP = ()=>{
    disableTOCPBind().then(res => {
      if(res.code === 0){
        ElMessage.success('Success')
        userStore.GetUserInfo()
      }
    })
  }
  </script>
  
  <style scoped lang="scss">
  .container {
    padding: 20px;
  }
  
  .el-form-item__label {
    font-size: 16px;
  }
  
  .el-button {
    width: 100px;
  }
  .help-text {
    padding-left: 20px;
    li {
      margin-bottom: 10px;
    }
  }
  @media (max-width: 768px) {
    .help-text {
      padding-left: 0;
    }
  }
  </style>