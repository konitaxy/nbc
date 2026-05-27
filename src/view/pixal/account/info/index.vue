<template>
    <div class="container form-white">
      <el-tabs v-model="activeTab" @tab-click="handleClick">
        <el-tab-pane :label="$t('lang.update_information')" name="1">
          <el-form ref="resetFormRef" :model="resetForm" :rules="formRules" label-width="140px">
            <el-form-item :label="$t('lang.email')">
              <el-input v-model="userInfo.email" disabled></el-input>
            </el-form-item>
            <!-- <el-form-item label="Email:">
              <el-input class="w-50" v-model="userInfo.email"></el-input>
            </el-form-item> -->
            <!-- <el-form-item label="收款码:">
              <el-input v-model="userInfo.paymentCode"></el-input>
            </el-form-item> -->
            <el-form-item prop="password" :label="$t('lang.old_password')">
              <el-input v-model="resetForm.password" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item prop="newPassword" :label="$t('lang.new_password')">
              <el-input v-model="resetForm.newPassword" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item prop="repeatPassword" :label="$t('lang.confirm_new_password')">
              <el-input v-model="resetForm.repeatPassword" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleConfirm">{{ $t('lang.submit_changes') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane :label="$t('lang.security_settings')" name="2">
          <el-form ref="pinFormRef" :model="resetForm" label-width="140px" :rules="pinRules">
            <el-form-item prop="pin" :label="$t('lang.pin')"> 
              <el-input v-model="resetForm.pin" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item prop="repeatPin" :label="$t('lang.confirm_pin')"> 
              <el-input v-model="resetForm.repeatPin" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item> 
              <el-button type="primary" @click="handleSetPin">{{ userStore.userInfo.bindPin?$t('lang.update_pin'):$t('lang.confirm') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane :label="$t('lang.session_active_info')" name="3">
            <el-table :data="sessionActions">
              <el-table-column prop="opSystem" :label="$t('lang.op_system')"></el-table-column>
              <el-table-column prop="application" :label="$t('lang.application')"></el-table-column>
              <el-table-column prop="address" :label="$t('lang.address')"></el-table-column>
              <el-table-column prop="ipAddress" :label="$t('lang.ip_address')"></el-table-column>
              <el-table-column :label="$t('lang.last_active_time')">
                <template #default="scope">
                  {{ formatTimeDifference(scope.row.lastActiveTime) }}
                </template>
                
              </el-table-column>
              <el-table-column prop="status" :label="$t('lang.status')">
                <template #default="scope">
                  <el-tag type="primary" v-if="scope.row.status">{{ $t('lang.current_session') }}</el-tag>
                  <el-tag type="info" v-else>{{ $t('lang.loginout') }}</el-tag>
                </template>
              </el-table-column>
            </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </template>
  
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import {formatTimeDifference} from '@/utils/format'
import { validatePassword} from '@/utils/validates'
import {changePassword,setPIN,listSessionLog} from '@/api/profile'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
const userStore = useUserStore()
const userInfo = ref({})
const sessionActions = ref([])
onMounted(() => {
  userInfo.value = userStore.userInfo
  if(userStore.userInfo.bindPin){
    resetForm.pin = "******"
  }
  listSessionLog().then(res =>{
    if(res.code ===0){
      sessionActions.value = res.data
    }
  })
})
const activeTab = ref('1')
const resetFormRef = ref()
const pinFormRef = ref()
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== resetForm.newPassword) {
    callback(new Error(t('lang.validation.password_mismatch')));
  } else {
    callback();
  }
}

const validateConfirmPin = (rule, value, callback) => {
  if (value !== resetForm.pin) {
    callback(new Error(t('lang.validation.pin_mismatch')));
  } else {
    callback();
  }
}
const handleConfirm = async ()=>{
  let r = await resetFormRef.value.validate()
  if(!r){
    return
  }
  resetForm.email = userStore.userInfo.email
  changePassword(resetForm).then(res =>{
    if(res.code === 0){
      ElMessage.success('change success')
    }
  })
}
const resetForm = reactive({
  password: '',
  newPassword: '',
  repeatPassword: ''
})
const formRules = reactive({
    password: [
      { validator: validatePassword, trigger: 'blur' }
    ],
    newPassword: [
      { validator: validatePassword, trigger: 'blur' }
    ],
    repeatPassword: [
      { validator: validateConfirmPassword, trigger: 'blur' }
    ]
  })

  const pinRules = reactive({
    pin: [
          {
            pattern: /^\d{6}$/,
            message: t('lang.validation.pin_format'),
            trigger: 'blur'
          }
        ],
    repeatPin: [
      { validator: validateConfirmPin, trigger: 'blur' }
    ]
  })
  const handleSetPin = async ()=>{
    let r = await pinFormRef.value.validate()
    if(!r){
      return
    }
    setPIN(resetForm).then(res =>{
      if(res.code ===0){
        ElMessage.success('Success')
      }
    })
  }
  </script>
  
  <style scoped>
  .container {
    padding: 20px;
  }
  :deep(.el-input) {
    max-width: 400px;
}
  </style>