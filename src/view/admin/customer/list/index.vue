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
          <el-form-item label="Manager">
            <el-input v-model="search.accountManager" placeholder="Please enter account manager"></el-input>
          </el-form-item>
        
          <el-form-item label="Client Status">
            <el-select v-model="search.clientStatus" placeholder="Please select client status">
              <el-option label="Review" :value="1"></el-option>
              <el-option label="Active" :value="2"></el-option>
              <el-option label="Suspend" :value="3"></el-option>
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
        <el-table-column prop="clientNo" label="Client No" min-width="120">
          <template #default="{row}">
            <span style="cursor: pointer;font-weight: 500;" >{{ row.clientNo }}&nbsp;&nbsp;<i class="bi bi-copy" @click="handleCopy(row.clientNo)"></i></span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="enName" label="English Name" width="120"></el-table-column> -->
        <!-- <el-table-column prop="ema" label="name" min-width="120"></el-table-column> -->
        <el-table-column prop="email" label="email" min-width="120"></el-table-column>
        <el-table-column prop="accountManager" label="Manager"  min-width="120">
          <template #default="{row}">
            <i @click="handleSetManager(row)" class="bi bi-pencil-square"></i>{{ row.accountManager }}
          </template>
        </el-table-column>
        <el-table-column prop="inviter" label="Inviter"  min-width="120">
          <template #default="{row}">
            {{ row.inviteUser.userName }}
          </template>

        </el-table-column>

        <el-table-column prop="clientType" label="Client Type" min-width="120"></el-table-column>
        <el-table-column prop="location" label="Location" min-width="90"></el-table-column>
        <el-table-column prop="clientStatus" label="Client Status" min-width="120">
          <template #default="{row}">
            <el-tag :type="row.clientStatus == 1?'info':row.clientStatus == 2?'success':'danger'">{{row.clientStatus == 1?'Review':row.clientStatus == 2?'Active':'Suspend'}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="Register Time" min-width="120">
          <template #default="{row}">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="Remark" min-width="120">
          <template #default="{row}">
            <i @click="handleEditRemark(row)" class="bi bi-pencil-square"></i>{{ row.remark }}
          </template>

        </el-table-column>
        <el-table-column label="Action" fixed="right" min-width="200">
          <template #default="{row}"> 
            <div>
              <el-button @click="handleClientStatusChange(row,3)" :disabled="row.clientStatus == 3" size="small" type="text">Suspence</el-button>
              <el-button @click="handleClientStatusChange(row,2)" :disabled="row.clientStatus == 2" size="small" type="text">Active</el-button>
              <el-button :disabled="row.clientStatus != 2" size="small" type="text" @click="handleROLogin(row)">RO Login</el-button>
          </div>
          </template>
          
        </el-table-column>
      </el-table>
      <div class="mt-4">
      <el-pagination
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="search.total"
        :page-sizes="[10, 50, 100]"
        :page-size="search.pageSize"
        :current-page="search.page"
        @size-change="handlePageSizeChange"
        @current-change="handlePageChange"
      ></el-pagination>
    </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from '@/utils/format';
import { listClient,setClientName,setClientManager,remarkClient,ddClient,reviewClient,changeClientStatus,adminLogin} from '@/api/client';
import { ElMessage, ElMessageBox } from 'element-plus';
import {writeText} from 'clipboard-polyfill'

const dialogs=reactive({
  clientDDDetailDialogVisible:false
})
const activeNames = ref(['1']);
const search = ref({
  clientNo: '',
  location: '',
  enName: '',
  accountManager: '',
  clientStatus: '',
  reviewStatus: '',
  riskLevel: '',
  email:'',
  page:1,
  pageSize:20,
  total:0
});

const tableData = ref([])

const getTableData = () => {
  listClient(search.value).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.value.total = res.data.total
    }
    
  })
};
const handlePageSizeChange =(val) =>{
  search.value.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.value.page = val
  getTableData()
}
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
  search.value = {
    clientId: '',
    clientNo: '',
    locationName: '',
    englishName: '',
    accountManager: '',
    clientStatus: '',
    reviewStatus: '',
    riskLevel: ''
  };
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