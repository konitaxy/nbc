<template>
    <div class="container-base">
      <div class="row mt-1">
        <div class="col">
          <el-form :inline="true" class="demo-form-inline" size="small">
            <el-form-item >
              <el-input v-model="search.cardLastFour" :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')"></el-input>
            </el-form-item>
            <el-form-item>
              <el-input v-model="search.userId" placeholder="请输入使用者ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-date-picker v-model="search.timeRange" type="datetimerange"
              :start-placeholder="$t('lang.start_time')"
              :end-placeholder="$t('lang.end_time')"
              ></el-date-picker>
            </el-form-item>
            <el-form-item>
              <el-select v-model="search.transactionType" :placeholder="$t('lang.please_select_transaction_type')" style="width: 150px;"  clearable>
                <el-option
                  v-for="(item, index) in transactionTypes"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button  type="primary" @click="getData">{{$t('lang.search')}}</el-button>
              <el-button v-if="$userStore.hasRole(6)" type="primary" @click="exportData">{{$t('lang.export')}}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
      <div class="row mt-5">
        <div class="col">
          <el-table :data="transactions" style="width: 100%">
            <el-table-column prop="ownerId" label="所有者ID"></el-table-column>
            <el-table-column prop="userId" label="使用者ID"></el-table-column>
            <el-table-column prop="cardNumber" label="卡号"></el-table-column>
            <el-table-column prop="transactionAmount" label="交易金额"></el-table-column>
            <el-table-column prop="channelType" label="渠道类型"></el-table-column>
            <el-table-column prop="deductionAmount" label="扣款金额"></el-table-column>
            <el-table-column prop="transactionFeeRate" label="交易费率"></el-table-column>
            <el-table-column prop="walletBalance" :label="$t('lang.wallet_balance')"></el-table-column>
            <el-table-column prop="transactionType" label="交易类型"></el-table-column>
            <el-table-column prop="operationTime" label="操作时间"></el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button size="small" @click="viewDetails(scope.row)">查看详情</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            background
            layout="total, sizes, prev, pager, next"
            :total="search.total"
            :page-size="search.pageSize"
            :page-sizes="[10, 50, 100]"
            :current-page="search.page"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          ></el-pagination>
        </div>
      </div>
      <el-dialog v-model="dialogVisible" title="查看详情" width="30%" center>
      <el-descriptions direction="vertical" :column="1" border>
        <el-descriptions-item label="流水号">{{ currentDetail.transactionId }}</el-descriptions-item>
        <el-descriptions-item label="卡号">{{ currentDetail.cardNumber }}</el-descriptions-item>
        <el-descriptions-item label="交易金额">{{ currentDetail.transactionAmount }}</el-descriptions-item>
        <el-descriptions-item label="扣款金额">{{ currentDetail.deductionAmount }}</el-descriptions-item>
        <el-descriptions-item label="交易费率">{{ currentDetail.transactionFeeRate }}%</el-descriptions-item>
        <el-descriptions-item :label="$t('lang.wallet_balance')">{{ currentDetail.walletBalance }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ currentDetail.transactionType }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentDetail.operationTime }}</el-descriptions-item>
        <el-descriptions-item :label="$t('lang.notes')">{{ currentDetail.remark }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
        </span>
      </template>
    </el-dialog>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { ElMessage } from 'element-plus';
  
  // 查询表单，包含分页信息
  const search = ref({
    cardLastFour: '',
    userId: '',
    startDate: '',
    endDate: '',
    transactionType: '',
    total: 0,
    pageSize: 10,
    page:1
  });
  const dialogVisible = ref(false);
  
  // 表格数据
  const transactions = ref([]);
  const currentDetail = ref({})
  // 交易类型列表
  const transactionTypes = [
    { value: 1, label: '开卡' },
    { value: 2, label: '充卡' },
    { value: 3, label: '交易退款' },
    { value: 4, label: '小额交易费(常规卡)' },
    { value: 5, label: '小额交易费(共享卡)' },
    { value: 6, label: '服务费' }
  ];
  
  // 获取数据的方法
  const getData = async () => {
    ElMessage.success("获取数据成功")
  };
  
  // 导出数据
  const exportData = () => {
    ElMessage.success('导出成功');
  };
  
  // 查看详情
  const viewDetails = (row) => {
    // console.log(row);
    currentDetail.value = row;
  };
  
  // 分页处理
  const handleSizeChange = (val) => {
    search.value.pageSize = val;
    search.value.page = 1;
    getData();
  };
  
  const handleCurrentChange = (val) => {
    search.value.page = val;
    getData();
  };
  </script>
  
  <style scoped>
  .container-base {
    padding: 10px;
  }
  
  .el-form-item__label {
    font-size: 16px;
  }
  
  .el-button {
    width: 100px;
  }
  </style>