<template>
  <div class="p-2">
    <!-- Search Form -->
    <el-form :inline="true" label-position="top" :model="search">
      <el-row :gutter="20" class="w-100">
        <el-col :span="6">
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client no"></el-input>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="6">
          <el-form-item label="Name">
            <el-input v-model="search.name" placeholder="Please enter name"></el-input>
          </el-form-item>
        </el-col> -->
        <el-col :span="6">
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="6">
          <el-form-item label="Transaction Date">
            <el-date-picker
              v-model="search.timeRange"
              type="daterange"
              range-separator="To"
              value-format="yyyy-MM-dd"
              start-placeholder="Start time"
              end-placeholder="End time"
            ></el-date-picker>
          </el-form-item>
        </el-col>
        <el-col :span="6"> -->
        <el-col :span="6">
          <el-form-item label="Order ID">
            <el-input v-model="search.orderId" placeholder="Please enter order id"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Status">
            <el-select v-model="search.status" placeholder="Please select status" clearable> 
              <el-option label="Pending" value="Pending"></el-option>
              <el-option label="Proceed" value="Success"></el-option>
              <el-option label="Decline" value="Decline"></el-option>
              <el-option label="Failure" value="Failure"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
        <el-form-item label=" ">
          <el-button type="primary" @click="onSearch">Search</el-button>
          <el-button @click="onReset">Reset</el-button>
        </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        
      </el-row>
    </el-form>
    <!-- <div style="text-align: right; padding-bottom:5px;">
      <el-button type="primary" @click="channelDataCompensation">Channel Data Compensation</el-button>
      <el-button type="primary" @click="exportData">Export</el-button>
    </div> -->
    <!-- Table -->
    <el-table :data="tableData" style="width: 100%" show-overflow-tooltip>
      <el-table-column prop="client.clientNo" label="Client No" width="120"></el-table-column>
      <el-table-column prop="client.email" label="Email" width="160"></el-table-column>
      <el-table-column label="Order Id" width="200">
        <template #default="{row}">
            <div style="font-size: 14px;  font-size: 12px;line-height: 1.2;">{{ row.thirdOrderId }}</div>
           <div style="color: #909399; font-size: 10px; line-height: 1.2;">{{ row.orderId }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="accountNumber" label="Receive Account" width="150"></el-table-column>
      <el-table-column prop="originAmount" label="Origin Remit Amount" width="180">
        <template #default="{row}">
          {{ row.originAmount }}<span style="font-size: small;" class="ms-1 px-1 rounded-1 bg-secondary text-light">{{ row.currency }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remitAmount" label="Actual Amount" width="150">
        <template #default="{row}">
          {{ row.remitAmount }}<span style="font-size: small;" class="ms-1 px-1 rounded-1 bg-secondary text-light">{{ row.currency }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="Final Amount" width="150"></el-table-column>
      <el-table-column prop="status" label="Status" width="100" fixed="right">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'Pending'" type="info">{{row.status}}</el-tag>
          <el-tag v-else-if="row.status === 'Success'" type="success">{{row.status}}</el-tag>
          <el-tag v-else-if="row.status === 'Decline'" type="danger">{{row.status}}</el-tag>
          <el-tag v-else-if="row.status === 'Failure'" type="danger">{{row.status}}</el-tag>

        </template>
      </el-table-column>
      <el-table-column prop="operator" label="operator" width="120"></el-table-column>
      <el-table-column label="Created Date" width="160">
        <template #default="{ row }">
          {{formatDate(row.CreatedAt)}}
        </template>
      </el-table-column>

      <el-table-column label="Action" width="100" fixed="right">
        <template #default="{ row }">
          <el-dropdown>
            <span class="el-dropdown-link">
              Operate
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :disabled="loading.editRechargeRecordLoading || row.status != 'Pending'" @click="handleEdit(row)">Edit</el-dropdown-item>
                <el-dropdown-item :disabled="loading.editRechargeRecordLoading || row.status != 'Pending'" @click="handleReviewRechargeConfirm(row,'Success')"> Inbound</el-dropdown-item>
                <el-dropdown-item :disabled="loading.editRechargeRecordLoading || row.status != 'Pending'" @click="handleReviewRechargeConfirm(row,'Decline')">Decline</el-dropdown-item>

              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
         
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="dialogs.editRechargeRecordDialog" title="Edit Recharge Record" width="50%" style="max-width: 500px;">
      <el-form label-position="top" label-width="auto" class="p-3">
          <el-form-item label="Origin Remit Amount">
            <el-input readonly v-model="currRow.originAmount">
              <template #append>
                {{ currRow.currency }}
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Actual Received Amount"  error="请输入金额" :validate-status="currRow.form.remitAmount >0?'success':'error'">
            <el-input v-model="currRow.form.remitAmount">
              <template #append>
                {{ currRow.currency }}
              </template>
            </el-input>
          </el-form-item>
          <!-- <el-form-item label="Remark">
              <el-input v-model="currRow.remark"></el-input>
            </el-form-item> -->
           
      </el-form>
      <template #footer>
              <div>
                <el-button type="info" @click="dialogs.editRechargeRecordDialog = false">cancel</el-button>
                <el-button :loading="loading.editRechargeRecordLoading" type="success" @click.prevent="handleEditRechargeRecordConfirm">confirm</el-button>
              </div>
            </template>
    </el-dialog>

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
</template>
<script setup>
import { ref,onMounted,reactive } from 'vue';
import { ElMessage } from 'element-plus';
import {formatDate} from '@/utils/format'
import {reviewRechargeRecord,listRechargeRecord,editRechargeRecord } from '@/api/finance'
const dialogs = reactive({
  editRechargeRecordDialog:false
})
const loading = reactive({
  editRechargeRecordLoading:false
})
const search = reactive({
  clientId: '',
  name: '',
  email: '',
  cardNo: '',
  status: '',
  orderId: '',
  page:1,
  pageSize:10,
  total:0
});

const tableData = ref([]);


const getTableData = () =>{
  listRechargeRecord(search).then(res => {
    if(res.code ===0){
      tableData.value = res.data.list
      search.total = res.data.total 
    }
  })
}

onMounted(() => {
  getTableData()
});
const onSearch = () => {
  getTableData()
};

const onReset = () => {
  
  search.clientId = ''
  search.name=''
  search.email=''
  search.cardNo=''
  search.status=''
  search.orderId=''
}

const handleOperate = (row) => {
  // Implement operate logic here
  console.log('Operating:', row);
};
const handleReviewRechargeConfirm = (row,status) =>{
  loading.editRechargeRecordLoading = true

  reviewRechargeRecord({
    id:row.ID,
    status:status
  }).then(res =>{
    loading.editRechargeRecordLoading = false

    if (res.code === 0){
      getTableData()
      ElMessage.success('Success')
    }
  })
}
const channelDataCompensation = () => {
  // Implement channel data compensation logic here
  ElMessage.success('Channel Data Compensation button clicked');
};

const exportData = () => {
  // Implement export logic here
  ElMessage.success('Export button clicked');
};
const currRow = ref(null);
const handleEdit = (row) => {
    currRow.value = row
    currRow.value.form = {
      remitAmount: row.remitAmount
    }
    dialogs.editRechargeRecordDialog = true
}
const handleEditRechargeRecordConfirm = () => {
  loading.editRechargeRecordLoading = true
  editRechargeRecord({
    id: currRow.value.ID,
    amount: currRow.value.form.remitAmount
  }).then(res =>{
    loading.editRechargeRecordLoading = false
    if(res.code === 0){
      dialogs.editRechargeRecordDialog = false
      getTableData()
      ElMessage.success('Edit Success')
    }
  })
}
const handlePageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.page = val
  getTableData()
}
</script>
<style scoped>
.container {
  padding: 20px;
}
.demo-form-inline {
  margin-bottom: 20px;
}
.el-pagination {
  margin-top: 20px;
  text-align: center;
}
:deep(.el-row){
  width: 100%;
  padding:0px
}
.el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
</style>