<template>
    <div class="container form-white">
      <el-form :inline="true">
        <el-form-item>
          <el-input v-model="search.cardNumberLastFour" :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="search.userId" placeholder="使用者ID"></el-input>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="search.dateRange"
            type="daterange"
            :range-separator="$t('lang.to')"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-select v-model="search.transactionType" :placeholder="$t('lang.please_select_transaction_type')" clearable>
            <el-option :label="$t('lang.recharge')" :value="1"></el-option>
            <el-option label="退款" :value="2"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">{{$t('lang.search')}}</el-button>
          <el-button v-if="$userStore.hasRole(6)" type="primary" @click="onExport">{{$t('lang.export')}}</el-button>
        </el-form-item>
      </el-form>
  
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="ownerId" label="所有者ID"></el-table-column>
        <el-table-column prop="userId" label="使用者ID"></el-table-column>
        <el-table-column prop="cardNumber" label="卡号"></el-table-column>
        <el-table-column prop="amount" label="充值/退回金额"></el-table-column>
        <el-table-column prop="fee" label="手续费"></el-table-column>
        <el-table-column prop="arrivalAmount" label="到账金额"></el-table-column>
        <el-table-column prop="balance" label="卡余额"></el-table-column>
        <el-table-column prop="operationType" label="操作类型"></el-table-column>
        <el-table-column prop="operationTime" label="操作时间"></el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="mini" @click="handleView(scope.$index, scope.row)">查看</el-button>
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
        <el-dialog v-model="dialogVisible" title="查看">
            <div class="dialog-content">
                <p>卡号 {{ selectedRow.cardNumber }}</p>
                <p>退款金额 {{ selectedRow.amount }}</p>
                <p>手续费 {{ selectedRow.fee }}</p>
                <p>到账金额 {{ selectedRow.arrivalAmount }}</p>
                <p>卡余额 {{ selectedRow.balance }}</p>
                <p>操作类型 {{ selectedRow.operationType }}</p>
                <p>操作时间 {{ selectedRow.operationTime }}</p>
            </div>
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
  
  const search = ref({
    cardNumberLastFour: '',
    userId: '',
    dateRange: [],
    transactionType: '',
    total: 0,
    pageSize: 10,
    page:1
  });
  const handleSizeChange = (val) => {
    search.pageSize = val;
  };
  
  const handleCurrentChange = (val) => {
    search.page = val;
  };
  
  const tableData = ref([
    {
      ownerId: '1332',
      userId: '1332',
      cardNumber: '553370******6232',
      amount: '47.63',
      fee: '0.00',
      arrivalAmount: '47.63',
      balance: '0.00',
      operationType: '退款',
      operationTime: '2025-07-01 01:58:02'
    },
    // 其他示例数据...
  ]);
  
  const onSearch = () => {
    console.log('搜索条件:', search.value);
  };
  
  const onExport = () => {
    console.log('导出数据');
  };
  
  const dialogVisible = ref(false);
  const selectedRow = ref({});
  
  const handleView = (index, row) => {
    selectedRow.value = row;
    dialogVisible.value = true;
  };
  </script>
  
  <style scoped>
  .container {
    padding: 20px;
  }
  .demo-form-inline {
    display: flex;
    align-items: center;
  }
  </style>