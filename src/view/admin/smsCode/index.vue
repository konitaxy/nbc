<template>
  <div class="p-2">
    <!-- Filter Form -->
    <el-form :inline="true"  label-position="top" :model="search" class="demo-form-inline">
      
          <el-form-item label="Email">
            <el-input v-model="search.to" placeholder="Please enter email"></el-input>
          </el-form-item>
        
          <el-form-item label="Type">
            <el-select v-model="search.eventType" placeholder="Please select client status" clearable>
              <el-option v-for="item of eventTypes" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
          <!-- <el-form-item label="Client Risk Level">
            <el-select v-model="search.riskLevel" placeholder="Please select client risk level">
              <el-option label="Low" value="low"></el-option>
              <el-option label="Medium" value="medium"></el-option>
              <el-option label="High" value="high"></el-option>
            </el-select>
          </el-form-item> -->
          <el-form-item label=" ">
            <el-button type="primary" @click="getTableData">Search</el-button>
            <el-button @click="reset">Reset</el-button>          
          </el-form-item>
      <!-- Add more rows as needed -->
          
    </el-form>
        
       

    <!-- Table -->
    <div  style="overflow-x: auto;">
      <el-table :data="tableData" show-overflow-tooltip	>
        <el-table-column prop="to" label="to" min-width="120"></el-table-column>
        <!-- <el-table-column prop="enName" label="English Name" width="120"></el-table-column> -->
        <el-table-column prop="code" label="code" min-width="120"></el-table-column>
        <el-table-column prop="createdAt" label="Send time" min-width="120">
          <template #default="{row}">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="eventType" label="Event type" min-width="120">
          <template #default="{row}">
            {{ formatFeeType(row.eventType) }}
          </template>

        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from '@/utils/format';
import { listSmsCode} from '@/api/finance';
import { ElMessage, ElMessageBox } from 'element-plus';
import {writeText} from 'clipboard-polyfill'

const dialogs=reactive({
  clientDDDetailDialogVisible:false
})
const search = ref({
  to:'',
  eventType:'',
});
const eventTypes = [
{
  label: 'Card Create',
  value: 'cardAdd'
},{
  label: 'Card cancel',
  value: 'cardCancel'
},{
  label: 'Card withdraw',
  value: 'cardWithdraw'
},{
  label: 'Card recharge',
  value: 'cardRecharge'
},{
  label: 'Wallet withdraw',
  value: 'walletWithdraw'
},{
  label: 'Verify setting',
  value: 'verifySetting'
},{
  label: 'Bind tocp',
  value: 'tocp'
},{
  label: 'Login',
  value: 'login'
},{
  label: 'Change password',
  value: 'changePassword'
},{
  label: 'Reset password',
  value: 'resetPassword'
},{
  label: 'Register',
  value: 'register'
},{
  label: 'IAM Create',
  value: 'iamAccountCreate'
},{
  label: 'IAM Login',
  value: 'iamLogin'
}]
const tableData = ref([])
const formatFeeType =(str) =>{
  const v =eventTypes.find(item => item.value === str)
  if(v) return v.label
  return "Unknown"
}
const getTableData = () => {
  listSmsCode(search.value).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.value.total = res.data.total
    }
    
  })
};
const handleCopy = (text) => {
  writeText(text)
  ElMessage.success('Copy Success')
};
const reset = () => {
  search.value = {
    to:'',
    eventType:'',
  };
}
onMounted(()=>{
  getTableData()
})
</script>
<style scoped>
.container {
  padding: 10px;
}
:deep(.el-select__wrapper) {
    min-width: 150px;
    padding: 9px 16px;
}
:deep(.el-dialog .el-dialog__body){
  padding: 0px;
}
</style>