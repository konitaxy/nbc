<template>
  <div class="p-4">
    <!-- Search Form -->
    <el-form :inline="true" :model="search" class="demo-form-inline">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client No."></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
        </el-col>
        
        <el-col :span="6">
          <el-form-item label="Status">
            <el-select v-model="search.status" placeholder="Please select status" clearable>
              <el-option label="Pending" value="Pending"></el-option>
              <el-option label="Proceed" value="Proceed"></el-option>
              <el-option label="Decline" value="Decline"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="getTableData">Search</el-button>
          <el-button @click="onReset">Reset</el-button>
        </el-col>
      </el-row>
      <el-row :gutter="20">
       
      </el-row>
    </el-form>

    <!-- Table -->
    <el-table :data="tableData" style="width: 100%" border show-overflow-tooltip>
      <el-table-column prop="orderId" label="Order ID" width="200">
        <template #default="{ row }">
          <a href="#" @click.prevent="viewDetails(row.orderId)">{{ row.orderId }}</a>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="Withdraw Amount" width="150">
        <!-- <template #default="{ row }">
          <a href="#" @click.prevent="viewDetails(row.orderId)">{{ row.orderId }}</a>
        </template> -->
      </el-table-column>
      <el-table-column prop="CreatedAt" label="Created Time" min-width="170">
        <template #default="{ row }">
          {{ formatDate(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column prop="finishTime" label="Finish Time" min-width="170">
        <template #default="{ row }">
          {{ formatDate(row.finishTime) }}
        </template>
      </el-table-column>
      <el-table-column prop="operator" label="Operator" width="100"></el-table-column>
      <el-table-column prop="status" label="Status" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'Pending'" type="info">Pending</el-tag>
          <el-tag v-else-if="row.status === 'Proceed'" type="success">Proceed</el-tag>
          <el-tag v-else-if="row.status === 'Decline'" type="danger">Decline</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="client.clientNo" label="Client No" width="100"></el-table-column>
      <el-table-column prop="client.email" label="Email" width="100"></el-table-column>
      <el-table-column prop="memo" label="Memo" width="100"></el-table-column>
      <el-table-column label="Action" min-width="200" fixed="right">
        <template #default="{ row }">
          <el-popconfirm
            width="220"
            icon-color="#626AEF"
            title="Sure to Decline this?"
            @confirm="handleReviewWithdrawConfirm(row,'Decline')"
          >
            <template #reference>
              <el-button :disabled="row.status != 'Pending'" type="text">Decline</el-button>
            </template>
            <template #actions="{ confirm, cancel }">
              <el-button size="small" @click="cancel">No!</el-button>
              <el-button
                type="danger"
                size="small"
                @click="confirm"
              >
                Confirm?
              </el-button>
            </template>
            </el-popconfirm>
            <el-popconfirm
                width="220"
                icon-color="#626AEF"
                title="Sure to Proceed this?"
                @confirm="handleReviewWithdrawConfirm(row,'Proceed')"
              >
            <template #reference>
              <el-button :disabled="row.status != 'Pending'" type="text">Proceed</el-button>
            </template>
            <template #actions="{ confirm, cancel }">
              <el-button size="small" @click="cancel">No!</el-button>
              <el-button
                type="danger"
                size="small"
                @click="confirm"
              >
                Confirm?
              </el-button>
            </template>
            </el-popconfirm>
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
</template>
<script setup>
import { ref,onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import {reviewWithdrawRecord,listWithdrawRecord } from '@/api/finance'
import { formatDate } from '@/utils/format';

const search = ref({
  clientNo: '',
  email: '',
  clientNo: '',
  status: '',
  page:1,
  pageSize:10,
  total:0
});

const tableData = ref([]);

const onSearch = () => {
  // Implement search logic here
  ElMessage.success('Search button clicked');
};

const onReset = () => {
  // Implement reset logic here
  Object.assign(search.value, {
    clientNo: '',
    email: '',
    clientNo: '',
    status: '',
    page:1,
    pageSize:10,
    total:0
  });
};
const getTableData = ()=>{
  listWithdrawRecord(search.value).then(res=>{
    tableData.value = res.data.list
    search.total = res.data.total
  })
}

const viewDetails = (orderId) => {
  // Implement view details logic here
  console.log('Viewing details for order:', orderId);
};
onMounted(()=>{
  getTableData()
})
const handleReviewWithdrawConfirm = (row,status) => {
  reviewWithdrawRecord({
    id:row.ID,
    status:status
  }).then(res =>{
    if (res.code === 0){
      getTableData()
      ElMessage.success('Success')
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

/* Custom styles for table */
.el-table .cell a {
  color: #67C23A; /* Green color for links */
  text-decoration: none;
}
:deep(.el-select__wrapper){
  min-width: 150px;
}
</style>