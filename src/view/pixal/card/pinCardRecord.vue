<template>
    <div class="p-2">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-input v-model="search.cardNumberLastFour" :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="search.userId" placeholder="请输入使用者ID"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">{{$t('lang.search')}}</el-button>
        </el-form-item>
      </el-form>
  
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="ownerId" label="所有者ID"></el-table-column>
        <el-table-column prop="userId" label="使用者ID"></el-table-column>
        <el-table-column prop="cardNumber" label="卡号"></el-table-column>
        <el-table-column prop="cardBin" label="卡Bin"></el-table-column>
        <el-table-column prop="currency" label="卡币种"></el-table-column>
        <el-table-column prop="cardOrganization" label="卡组织"></el-table-column>
        <el-table-column prop="balance" label="卡片余额"></el-table-column>
        <el-table-column prop="cardStatus" label="卡状态"></el-table-column>
        <el-table-column prop="createTime" label="创建时间"></el-table-column>
        <el-table-column prop="remark" :label="$t('lang.notes')"></el-table-column>
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
      <el-dialog v-model="dialogVisible" title="查看详情">
        <p>所有者ID: {{ selectedRow.ownerId }}</p>
        <p>使用者ID: {{ selectedRow.userId }}</p>
        <p>卡号: {{ selectedRow.cardNumber }}</p>
        <p>卡Bin: {{ selectedRow.cardBin }}</p>
        <p>卡币种: {{ selectedRow.currency }}</p>
        <p>卡组织: {{ selectedRow.cardOrganization }}</p>
        <p>卡片余额: {{ selectedRow.balance }}</p>
        <p>卡状态: {{ selectedRow.cardStatus }}</p>
        <p>创建时间: {{ selectedRow.createTime }}</p>
        <p>备注: {{ selectedRow.remark }}</p>
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
    total: 0,
    pageSize: 10,
    page:1
  });
  const handleSizeChange = (val) => {
    search.pageSize = val;
    // 更新数据
  };
  
  const handleCurrentChange = (val) => {
    search.page = val;
    // 更新数据
  };
  
  const tableData = ref([
    {
      ownerId: '1332',
      userId: '1332',
      cardNumber: '534786******4224',
      cardBin: '534786',
      currency: 'USD',
      cardOrganization: 'MasterCard',
      balance: '0.00',
      cardStatus: '销卡',
      createTime: '2025-05-06 12:24:30',
      remark: ''
    },
    {
      ownerId: '1332',
      userId: '1332',
      cardNumber: '531993******0243',
      cardBin: '531993',
      currency: 'USD',
      cardOrganization: 'MasterCard',
      balance: '0.00',
      cardStatus: '销卡',
      createTime: '2025-05-05 03:22:53',
      remark: ''
    },
    {
      ownerId: '1332',
      userId: '1332',
      cardNumber: '404038******7693',
      cardBin: '404038',
      currency: 'USD',
      cardOrganization: 'VISA',
      balance: '0.00',
      cardStatus: '销卡',
      createTime: '2025-05-05 02:48:45',
      remark: ''
    }
  ]);
  
  const onSearch = () => {
    // 处理搜索逻辑，可以根据 search.cardNumberLastFour 和 search.userId 进行筛选
    console.log('搜索条件:', search.value);
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