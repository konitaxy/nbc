<template>
    <div class="form-list-container">
      <!-- 操作栏 -->
      <div class="operation-bar">
        <el-select v-model="search.status" :placeholder="$t('lang.please_select_audit_status')" @change="handleSearch" style="width: 150px;" clearable>
          <el-option :label="$t('lang.transfer_pending')" value="Pending"></el-option>
          <el-option :label="$t('lang.success')" value="Success"></el-option>
          <el-option :label="$t('lang.decline')" value="Decline"></el-option>
        </el-select>
        <el-button class="ms-2" type="primary" icon="search" @click="handleSearch">{{$t('lang.search')}}</el-button>
        <el-button type="primary" icon="plus" @click="handleRecharge">{{$t('lang.top_up')}}</el-button>
        <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="handleExport">{{$t('lang.export')}}</el-button>
      </div>
  
      <!-- 充值记录表格 -->
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="orderId" :label="$t('lang.order_number')" width="180"></el-table-column>
        <el-table-column :label="$t('lang.account_type')">
          <template #default="scope">
            <span>{{ $t(scope.row.rechargeType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="accountNumber" :label="$t('lang.remit_account')" width="200">
            <template #default="scope">
              <el-tooltip
                :content="formatRemitAccount(scope.row)"
                placement="top"
                effect="dark"
                popper-class="remit-account-tooltip"
                :show-after="120"
                :hide-after="0"
                :offset="8"
              >
                <span class="remit-account-cell">{{ formatRemitAccount(scope.row) }}</span>
              </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column prop="originAmount" :label="$t('lang.remit_amount')"></el-table-column>
        <el-table-column prop="amount" :label="$t('lang.receive_amount')"></el-table-column>
        <el-table-column prop="status"  :label="$t('lang.audit_status')">
          <template #default="scope">
            <span>{{ $t(scope.row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column  :label="$t('lang.created_time')" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
      </el-table>
  
      <!-- 分页 -->
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
      <el-dialog
        :title="$t('lang.wallet_recharge')"
        v-model="dialogs.rechargeDialogVisible"
        fullscreen
        append-to-body
        class="recharge-form-dialog"
      >
        <RechargeForm/>
      </el-dialog>
    </div>
  </template>

<script>
export default {
  name: 'RechargeManager',
}
</script>
<script setup>
import {reactive, ref,onMounted} from 'vue'
import RechargeForm from '@/view/pixal/common/recharge_form.vue'
import {listWalletRecharge} from '@/api/finance'
import {formatDate} from '@/utils/format'
import {buildExcel} from '@/utils/excel'

const dialogs = reactive({
  rechargeDetailDialog:false,
  rechargeDialogVisible:false
})
const search = reactive({
  timeRange: [],
  transactionType: null,
  total: 0,
  pageSize: 10,
  page:1
});

const tableData = ref([])
const formatRemitAccount = (row) => {
  const accountType = row?.accountType ? `(${row.accountType})` : ''
  return `${accountType}${row?.accountNumber || ''}`
}

const getTableData = ()=>{
  listWalletRecharge(search).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.total = res.data.total
    }
   })
}
const handleSearch =() =>{
  getTableData()
}
onMounted(()=>{
  getTableData()
})
const handleRecharge =() =>{
  dialogs.rechargeDialogVisible = true;
}
const handleExport =() =>{
  var search2 = {
    ...search,
    pageSize: 999999,
    page:1
  };
  listWalletRecharge(search2).then(res => {
    if (res.code === 0){
      let list = res.data.list.map(x =>{
          return {
            
            "Order No":x.orderId,
            "Account Type":x.accountType,
            "Remit Account": `(${x.accountType})${x.accountNumber}`,
            "Remit Amount":parseFloat(x.originAmount),
            "Receive Amount":parseFloat(x.amount),
            "Audit Status":x.status,
            "Update Time":formatDate(x.CreatedAt)
          }
        })
        buildExcel(list,"wallet_recharge_history")
    }
   })
}
const currentRecord = ref()
const handleView =(val) => {
  currentRecord.value = val
  dialogs.rechargeDetailDialog = true
}
const handleSizeChange=(val) => {
  search.pageSize = val;
  getTableData()
}
const  handleCurrentChange =(val) => {
  search.page = val;
  getTableData()
}
</script>

<style scoped>


.operation-bar {
  margin-bottom: 20px;
}

.remit-account-cell {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  white-space: nowrap;
}

:global(.remit-account-tooltip) {
  max-width: min(560px, calc(100vw - 32px));
  border: 1px solid rgba(139, 214, 255, 0.28) !important;
  border-radius: 10px !important;
  background: rgba(5, 16, 29, 0.96) !important;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.28);
  color: #f4fbff !important;
  font-weight: 700;
}

:global(.remit-account-tooltip .el-popper__arrow::before) {
  border-color: rgba(139, 214, 255, 0.28) !important;
  background: rgba(5, 16, 29, 0.96) !important;
}

</style>
