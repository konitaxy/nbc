<template>
  <div class="p-2">
    <!-- Filter Form -->
    <el-form :inline="true"  label-position="top" :model="search" class="demo-form-inline">
      
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client no"></el-input>
          </el-form-item>
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
          <el-form-item label="Month">
            <el-date-picker type="month" format="YYYY-MM" value-format="YYYY-MM" v-model="search.month"></el-date-picker>
          </el-form-item>
          <el-form-item label=" ">
            <el-button type="primary" @click="getTableData">Search</el-button>
            <el-button @click="reset">Reset</el-button>          
          </el-form-item>
      <!-- Add more rows as needed -->
          
    </el-form>
        
       

    <!-- Table -->
    <div  style="overflow-x: auto;">
      <el-table :data="tableData" show-overflow-tooltip	>
        <el-table-column prop="clientNo" label="Client No" width="90">

        </el-table-column>
        <!-- <el-table-column prop="ema" label="name" min-width="120"></el-table-column> -->
        <el-table-column prop="email" label="email" width="170"></el-table-column>
        <el-table-column label="Month" width="120">
          <template #default>
            {{ search.month || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="totalCardBalance" label="Card Balance" width="170"></el-table-column>
        <el-table-column prop="walletBalance" label="Wallet Balance" width="170"></el-table-column>

        <el-table-column prop="walletRechargeAmount" label="Total recharge" width="120"></el-table-column>
        <el-table-column prop="authorizationAmount" label="Authorization" width="120"></el-table-column>
        <el-table-column prop="clearingAmount" label="Clearing Amount" width="170"></el-table-column>
        <el-table-column prop="clearingCrossBoardAmount" label="Cross Clearing Amount" width="170"></el-table-column>
      </el-table>
      <div class="mt-1">
            <el-pagination
              style="padding-top: 0px;float:right"
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="search.total"
              :page-sizes="[10, 50, 100]"
              :page-size="search.pageSize"
              :current-page="search.page"
              @size-change="handleLogPageSizeChange"
              @current-change="handleLogPageChange"
            ></el-pagination>
          </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from '@/utils/format';
import { listClient,setClientName,setClientManager,remarkClient,ddClient,reviewClient,changeClientStatus,adminLogin} from '@/api/client';
import { listStatReportByClient} from '@/api/finance';
import { ElMessage, ElMessageBox } from 'element-plus';
import {writeText} from 'clipboard-polyfill'

const dialogs=reactive({
  clientDDDetailDialogVisible:false
})
// 获取当前月份（格式：YYYY-MM）
const getCurrentMonth = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  return `${year}-${month}`
}

const activeNames = ref(['1']);
const search = reactive({
  clientNo: '',
  location: '',
  enName: '',
  accountManager: '',
  clientStatus: '',
  reviewStatus: '',
  riskLevel: '',
  email:'',
  month: getCurrentMonth(),
  dateRange: [],
  page:1,
  pageSize:10,
  total:0
});

const tableData = ref([])

// 将月份转换为日期范围
const convertMonthToDateRange = (month) => {
  if (!month) {
    return []
  }
  // 月份格式: YYYY-MM
  const [year, monthNum] = month.split('-')
  // 获取该月的第一天
  const startDate = `${year}-${monthNum}-01`
  // 获取该月的最后一天
  const lastDay = new Date(year, monthNum, 0).getDate()
  const endDate = `${year}-${monthNum}-${String(lastDay).padStart(2, '0')}`
  return [startDate, endDate]
}

const getTableData = () => {
  // 将月份转换为日期范围
  const params = { ...search }
  if (params.month) {
    params.dateRange = convertMonthToDateRange(params.month)
  } else {
    params.dateRange = []
  }
  // 删除 month 字段，只传递 dateRange
  delete params.month
  
  listStatReportByClient(params).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.total = res.data.total
    }
    
  })
};
const handleCopy = (text) => {
  writeText(text)
  ElMessage.success('Copy Success')
};
const handleClientStatusChange = (row,status)=>{
  changeClientStatus({
    id: row.ID,
    clientStatus: status
  }).then(res =>{
    if (res.code === 0){
      getTableData()
      ElMessage.success('Success')
    }
  })
}
const handleROLogin = (row)=>{
  adminLogin({
    id: row.ID,
  }).then(res =>{
    if (res.code === 0){
      if(res.data != null){
       let a = document.createElement('a')
       a.target = '_blank'
       a.href = res.data
       a.click()
       a.remove()
      }
    }
  })
}
const reset = () => {
  search.clientId = ''
  search.clientNo =''
  search.email =''
  search.month = ''
  search.dateRange = []
}
onMounted(()=>{
  getTableData()
})
const handleEditRemark = (row) => {
  ElMessageBox.prompt('Please enter your remark', 'Edit Remark', {
    confirmButtonText: 'Save',
    cancelButtonText: 'Cancel',
    inputValue: row.remark,
  }
  ).then((val)=>{
    remarkClient({
      id: row.ID,
      remark: val.value
    }).then(res =>{
      if (res.code === 0){
        ElMessage({
          type: 'success',
          message: 'Edit Success'
        });
      }
    })
  })
};
const handleLogPageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handleLogPageChange =(val) =>{
  search.page = val
  getTableData()
}


const handleSetManager = (row) => {
  ElMessageBox.prompt('Please enter manager name', 'Set manager', {
    confirmButtonText: 'Save',
    cancelButtonText: 'Cancel',
    inputValue: row.accountManager,
  }
  ).then((val)=>{
    setClientManager({
      id: row.ID,
      accountManager: val.value
    }).then(res =>{
      if (res.code === 0){
        ElMessage({
          type: 'success',
          message: 'Edit Success'
        });
      }
    })
  })
};
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