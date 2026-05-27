<template>
    <div class="container form-white">
      <h2>{{ $t('lang.authentication_settings') }}</h2>
      <el-form label-width="150px" :model="settingsForm" class="mt-4">
        <el-row :gutter="20">
          <!-- 第一列 -->
          <el-col :span="8">
            <el-checkbox v-model="settingsForm.cardCancel">{{ $t('lang.cancel_card_application') }}(Pixel)</el-checkbox>
            <el-checkbox v-model="settingsForm.cardAdd">{{ $t('lang.open_card') }}(Pixel)</el-checkbox>
            <el-checkbox v-model="settingsForm.cardDetail">{{ $t('lang.card_detail') }}</el-checkbox>
            <el-checkbox v-model="settingsForm.cardWithdraw">{{ $t('lang.withdrawal_application') }}(Pixel)</el-checkbox>
            <el-checkbox v-model="settingsForm.cardRecharge">{{$t('lang.card_recharge')}}(Pixel)</el-checkbox>
            <el-checkbox v-model="settingsForm.walletWithdraw">{{$t('lang.wallet_withdraw')}}</el-checkbox>
          </el-col>
  
          <!-- 第二列 -->
          <el-col :span="8">
            <el-checkbox disabled v-model="settingsForm.changePassword">{{ $t('lang.update_password') }}({{ $t('lang.required') }})</el-checkbox>
            <el-checkbox disabled v-model="settingsForm.verifySetting">{{ $t('lang.verify_configuration') }}({{ $t('lang.required') }})</el-checkbox>
            <el-checkbox disabled v-model="settingsForm.disableTocp">{{ $t('lang.disable_tocp') }}({{ $t('lang.required') }})</el-checkbox>

          </el-col>
        </el-row>
  
        <!-- 提交按钮 -->
        <el-button type="primary" @click="submitSettings">{{ $t('lang.save_settings') }}</el-button>
      </el-form>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import axios from 'axios';
  import { ElMessage } from 'element-plus';
  import { useUserStore } from '@/pinia/modules/user'
  import {verifySetting} from '@/api/client'
  const userStore = useUserStore()

  const settingsForm = ref(userStore.userInfo.verifySetting);
  
  // 提交设置的方法
  const submitSettings = async () => {
    verifySetting({
      "setting": settingsForm.value,
    }).then(res =>{
      if(res.code === 0){
        userStore.GetUserInfo()
        settingsForm.value = userStore.userInfo.verifySetting
        ElMessage.success('Success')
      }
    })
  };
  </script>
  
  <style scoped>
  .container {
    padding: 20px;
  }
  
  .el-checkbox {
    margin-bottom: 10px;
    display: flex !important;
    padding: 10px;
  }
  
  .el-button {
    margin-top: 20px;
  }
  </style>