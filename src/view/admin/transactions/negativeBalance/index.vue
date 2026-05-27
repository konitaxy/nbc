<template>
  <div class="p-2">
    <!-- Search Form -->
    <el-form :inline="true" :model="search" label-position="top" class="demo-form-inline">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client no"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="6">
          <el-form-item label="Name">
            <el-input v-model="search.name" placeholder="Please enter unit nickname"></el-input>
          </el-form-item>
        </el-col> -->
        <el-col :span="6">
          <el-form-item label="Card ID">
            <el-input v-model="search.cardId" placeholder="Please enter card id"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="Card No.">
            <el-input v-model="search.cardNo" placeholder="Please enter card no."></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Account Type">
            <el-select v-model="search.accountType" placeholder="Please select account type">
              <el-option label="Option 1" value="1"></el-option>
              <el-option label="Option 2" value="2"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="6">
          <el-form-item label="Available Balance Quick Sort">
            <el-select v-model="search.availableBalanceSort" placeholder="Please select available balance sort">
              <el-option label="Ascending" value="asc"></el-option>
              <el-option label="Descending" value="desc"></el-option>
            </el-select>
          </el-form-item>
        </el-col> -->
        <el-col :span="6" class="d-flex align-items-center justify-content-center">
          <el-button type="primary" @click="onSearch">Search</el-button>
          <el-button @click="onReset">Reset</el-button>
        </el-col>
      </el-row>
    </el-form>

    <!-- Table -->
    <el-table :data="tableData" style="width: 100%" border>
      <el-table-column prop="client.clientNo" label="Client No" width="150"></el-table-column>
      <el-table-column prop="client.email" label="Email" width="180"></el-table-column>
      <el-table-column prop="cardNo" label="Card No." ></el-table-column>
      <el-table-column prop="cardBrand" label="Card Model" width="120"></el-table-column>
      <el-table-column prop="cardType" label="Card Type" width="100"></el-table-column>
      <el-table-column prop="client.accountType" label="Account Type" width="150"></el-table-column>
      <el-table-column prop="balance" label="Balance" width="150"></el-table-column>
      <el-table-column prop="currency" label="Currency" width="100"></el-table-column>
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
import { reactive, ref,onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { listCard} from '@/api/card'

const search = reactive({
  negative: true,
  clientNo: '',
  email: '',
  unitNickname: '',
  cardId: '',
  cardNo: '',
  accountType: '',
  availableBalanceSort: '',
  page:1,
  pageSize:10,
  total:0
});

const tableData = ref([]);

const getTableData = () => {
  listCard(search).then(res =>{
    if (res.code === 0) {
      tableData.value = res.data.list;
      search.total = res.data.total;
    }
  })
  
};
onMounted(() => {
  getTableData();
});
const onSearch = () => {
  getTableData()
};
const handlePageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.page = val
  getTableData()
}
const onReset = () => {
  search.page = 1;
  search.pageSize = 10;
  search.negative = true;
};

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