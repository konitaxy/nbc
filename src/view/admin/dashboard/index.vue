<template>
    <div class="p-2">
      <!-- Balance Section -->
      <div class="row mt-4">
        <div class="col-md-6  w-100">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span class="fw-bold">Balance</span>
            </div>
            <el-row :gutter="20">
              <el-col :span="8"><div class="grid-content bg-purple h-100">Channel Total Balance <span class="fw-bold fs-5">${{balances.totalChannelBalance}}</span></div></el-col>
              <el-col :span="8"><div class="grid-content bg-purple-light">Customer Total Wallet Balance <span class="fw-bold fs-5">${{balances.totalWalletBalance}}</span></div></el-col>
              <el-col :span="8"><div class="grid-content bg-purple-light">Customer Total Card Balance <span class="fw-bold fs-5">${{balances.totalCardBalance}}</span></div></el-col>

            </el-row>
          </el-card>
        </div>
      </div>
  
      <!-- Statistics in the time range Section -->
      <div class="row mt-4">
        <div class="col-md-12">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span class="fw-bold">Statistics in the time range</span><br>
              <el-date-picker @change="handleStatDateRangeChange" v-model="search.dateRange" type="daterange" range-separator="To" start-placeholder="Start date" end-placeholder="End date" value-format="YYYY-MM-DD"></el-date-picker>
            </div>
            <el-row :gutter="20">
              <el-col :span="6"><div class="grid-content bg-purple">Customer Recharge Amount <br/><span>${{ summary.walletRechargeAmount }}</span></div></el-col>
              <el-col :span="6"><div class="grid-content bg-purple-light">Card Authorization Amount<br/><span>${{ summary.authorizationAmount }}</span></div></el-col>
              <el-col :span="6"><div class="grid-content bg-purple-light">Customer Fee Amount<br/><span>${{ summary.feeAmount }}</span></div></el-col>

              <el-col :span="6"><div class="grid-content bg-purple">Amount of Withdraw <br/><span>${{ summary.walletWithdrawAmount }}</span></div></el-col>
            </el-row>
            <el-row :gutter="20" class="mt-4">
              <el-col :span="12"><div class="grid-content bg-purple">Counts of Card Creation<br/><span>{{ summary.cardCreateCount }}</span></div></el-col>
              <el-col :span="12"><div class="grid-content bg-purple-light">Counts of Card Cancel<br/><span>{{ summary.cardCancelCount}}</span></div></el-col>
            </el-row>
          </el-card>
        </div>
      </div>
  
      <!-- Last 7 Day Statistics Section -->
      <div class="row mt-4">
        <div class="col-md-12">
          <el-card class="box-card" header-class="card-header">
            <div slot="header" class="clearfix">
              <span class="fw-bold">Last 7 Day Statistics</span>
            </div>
            <el-table :data="tableData" style="width: 100%">
              <el-table-column prop="reportDay" label="Date" width="120"></el-table-column>
              <el-table-column prop="walletRechargeAmount" label="Customer Recharge Amount($)" min-width="180"></el-table-column>
              <el-table-column prop="authorizationAmount" label="Card Authorization Amount($)" min-width="180"></el-table-column>
              <el-table-column prop="feeAmount" label="Fee Amount($)" min-width="120"></el-table-column>

              <el-table-column prop="walletWithdrawAmount" label="Amount of Withdraw($)" min-width="150"></el-table-column>
              <el-table-column prop="cardCreateCount" label="Counts of Card Creation" min-width="180"></el-table-column>
            </el-table>
          </el-card>
        </div>
      </div>
    </div>
  </template>
  <script setup>
  import { ref,onMounted, reactive } from 'vue';
  import { formatDateYYYYMMDD } from '@/utils/format';
  import {statBalance,statReport,listStatReport} from '@/api/finance'
  const search = reactive({
    dateRange: [],
    rangeType: 1,
  });
  const tableData = ref([]);
  const summary = ref({});
  const balances = ref({})
  const handleStatDateRangeChange = ()=>{
    var form = {}
    if (search.dateRange.length > 0) {
      form = {
        startTime: search.dateRange[0],
        endTime: search.dateRange[1],
      }
    }
    statReport(form).then(res => {
      summary.value = res.data;
    })
  }
  onMounted(() => {
    statBalance().then((res) => {
      balances.value = res.data;
    })
    var form = {}
    if (search.dateRange.length > 0) {
      form = {
        startTime: search.dateRange[0],
        endTime: search.dateRange[1],
      }
    }
    statReport(form).then(res => {
      summary.value = res.data;
    })
    const now = new Date();
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(now.getDate() - 7);
    var tableSearch = {
      startTime: formatDateYYYYMMDD(sevenDaysAgo),
      endTime: formatDateYYYYMMDD(now),
    }
    listStatReport(tableSearch).then(res => {
      tableData.value = res.data;
    })
  });
  </script>
  <style scoped>
  .box-card {
    width: 100%;
  }
  .grid-content {
    border-radius: 4px;
    min-height: 36px;
    padding: 20px 10px;
    text-align: start;
  }
  .grid-content span {
    font-weight: 600;
  }
  .bg-purple {
    background: #d3dce6;
  }
  .bg-purple-light {
    background: #e5e9f2;
  }
  </style>