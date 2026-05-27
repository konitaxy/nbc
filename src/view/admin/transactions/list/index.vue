<template>
  <div class="p-2">
    <el-form :inline="true" label-position="top" :model="search">
      <el-row :gutter="20" class="w-100">
        <el-col :span="6">
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client id"></el-input>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="6">
          <el-form-item label="Name">
            <el-input v-model="search.name" placeholder="Please enter unit nickname"></el-input>
          </el-form-item>
        </el-col> -->
        <el-col :span="6">
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Transaction Date">
            <el-date-picker
              v-model="search.timeRange"
              type="daterange"
              range-separator="To"
              value-format="YYYY-MM-DD"
              start-placeholder="Start time"
              end-placeholder="End time"
            ></el-date-picker>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        
        <el-col :span="6">
          <el-form-item label="Card NO.">
            <el-input v-model="search.cardNo" placeholder="Please enter card no."></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Transaction ID">
            <el-input v-model="search.transactionId" placeholder="Please enter transaction ID"></el-input>

          </el-form-item>
        </el-col>
        <el-col :span="6">
        <el-form-item label=" ">
          <el-button type="primary" @click="onSearch">Search</el-button>
          <el-button @click="onReset">Reset</el-button>
        </el-form-item>
        </el-col>
      </el-row>
      
      <!-- <el-row :gutter="20">
        <el-col :span="6">
          <el-button type="primary" @click="onSearch">Search</el-button>
          <el-button @click="onReset">Reset</el-button>
        </el-col>
      </el-row> -->
    </el-form>
    <!-- <div style="text-align: right; padding-bottom:5px;">
      <el-button type="primary" @click="channelDataCompensation">Channel Data Compensation</el-button>
      <el-button type="primary" @click="exportData">Export</el-button>
    </div> -->
    <!-- Table -->
    <el-table :data="tableData" style="width: 100%" border show-overflow-tooltip>
      <el-table-column prop="client.clientNo" label="Client No" width="100"></el-table-column>
      <el-table-column prop="client.email" label="Email" width="200"></el-table-column>
      <!-- <el-table-column prop="cardId" label="Card ID" width="200"></el-table-column> -->
      <el-table-column prop="card.cardNo" label="Card NO." width="180">
        <template #default="{row}">
          <div class="vstack justify-content-center text-end text-nowrap" style="display: inline-block;">
            <p class="" style="line-height: 18px;"><span>{{row.card.cardNo}}</span></p>
            <p style="line-height: 11px;"><span style="font-size: 11px;line-height: 8px;">{{ row.cardId  }}</span></p>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="Amount" width="100">
        <template #default="{row}">
            {{ row.amount }}<span class="ms-1">{{row.currency}}</span>
          </template>
      </el-table-column>
      <el-table-column prop="transactionType" label="Transaction Type" width="150">
        <template #default="{row}">
      <div>
        <p style="line-height: 18px;">{{ row.transactionType =='Card_Recharge'?'Recharge':row.transactionType =='Card_Withdraw'?'Withdraw':row.transactionType }}</p>
        <p class="hstack"><span :class="row.status=='Success'?'bg-success':'bg-danger'" style="width: 6px;height: 6px;border-radius: 50%;margin-right: 3px;"></span><span :class="row.status=='Success'?'text-success':'text-danger'" style="font-size: 11px;line-height: 8px;">{{ row.status }}</span></p>
      </div>
    </template>
      </el-table-column>
      <el-table-column label="Transaction ID" min-width="220" >
          <template #default="{row}">
          <div class="vstack justify-content-center text-end text-nowrap" style="display: inline-block;">
            <p class="" style="line-height: 18px;"><span>{{row.transactionId}}</span></p>
            <p style="line-height: 11px;"><span style="font-size: 11px;line-height: 8px;">{{ formatDate(row.transactionTime)  }}</span></p>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Other Message" min-width="220" >
          <template #default="{row}">
          <div class="vstack justify-content-center text-end text-nowrap" style="display: inline-block;">
            <p class="" style="line-height: 15px;"><span>{{row.failReason}}</span></p>
            <p style="line-height: 15px;"><span style="font-size: 15px;line-height: 8px;">{{ row.authCode  }}</span></p>
          </div>
        </template>
      </el-table-column>
      <!-- <el-table-column label="Action" width="100">
        <template #default="{ row }">
          <el-button type="info" @click="handleOperate(row)">Operate</el-button>
        </template>
      </el-table-column> -->
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
    <!-- Buttons -->
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from "@/utils/format";
import { ElMessage } from 'element-plus';
import { listCardTransaction} from '@/api/card'

const search = reactive({
  email: '',
  timeRange: [],
  clientNo: '',
  cardNo: '',
  name: '',
  email: '',
  page:1,
  pageSize:10,
  total:0
});
const tableData = ref([]);
const getTableData = () => {
  listCardTransaction(search).then((res) => {
    if(res.code === 0){
      tableData.value = res.data.list;
      search.total = res.data.total;
    }
  });
};
onMounted(() => {
  getTableData();
});
const onSearch = () => {
  getTableData()
};

const onReset = () => {
  // Implement reset logic here
  Object.assign(search, {
    userId: '',
    clientId: '',
    clientNo: '',
    unitNickname: '',
    email: '',
    cardNo: '',
    primaryCardId: '',
    primaryCardNo: '',
    ptxId: '',
    channelTransactionId: '',
    bizType: '',
    cardChannel: '',
    messageType: '',
    authResult: '',
    transactionDate: ''
  });
};
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
</style>