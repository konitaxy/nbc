<template>
    <div class="recharge-records">
      <!-- 操作栏 -->
      <div class="operation-bar">
        <el-select v-model="search.transactionType" :placeholder="$t('lang.please_select_transaction_type')" @change="handleSearch" style="width: 150px;" clearable>
          <!-- <el-option label="开卡" value="Pending"></el-option> -->
          <el-option :label="$t('Wallet_Recharge')" value="Wallet_Recharge"></el-option>
          <el-option :label="$t('Wallet_Withdraw')" value="Wallet_Withdraw"></el-option>
          <el-option :label="$t('Card_Recharge')" value="Card_Recharge"></el-option>
          <el-option :label="$t('Card_Withdraw')" value="Card_Withdraw"></el-option>
          <el-option :label="$t('Fee')" value="Fee"></el-option>


        </el-select>
        <el-button class="ms-2" type="primary" icon="search" @click="handleSearch">{{$t('lang.search')}}</el-button>
        <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="onExport">{{$t('lang.export')}}</el-button>
      </div>
  
      <!-- 充值记录表格 -->
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="orderId" :label="$t('lang.order_number')" width="180">
          <template #default="{row}">
            <span v-if="row.transactionRecord.ID == 0">{{ row.orderId }}</span>
            <el-link class="fs-6 fw-normal" v-else @click="handleViewOrderDetail(row)">{{ row.orderId }}<i class="bi bi-search"></i></el-link>
          </template>
        </el-table-column>
        <el-table-column prop="cardNo" :label="$t('lang.card_number')" width="200">
          
        </el-table-column>
        <el-table-column :label="$t('lang.type')" min-width="140">
          <template #default="scope">
            <span>{{ $t(scope.row.transactionType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('lang.amount')" width="150">
          <template #default="scope">
            <span>{{ `${scope.row.amount} ${scope.row.amountCurrency}` }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" :label="$t('lang.account_balance')"  width="150">
          <template #default="scope">
            <span>{{ `${scope.row.balance} ${scope.row.currency}` }}</span>
          </template>
        </el-table-column>
        <el-table-column  :label="$t('lang.update_time')"  min-width="150">
          <template #default="scope">
            {{ formatDate(scope.row.UpdatedAt) }}
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
        destroy-on-close
      >
        <RechargeForm/>
      </el-dialog>
      <el-dialog
        title="Detail"
        v-model="dialogs.viewOrderDetailDialogVisible"
      >
      <el-form label-position="right" label-width="150px">
        <!-- <el-form-item :label="$t('lang.card_number')+':'">
          {{ transactionDetail.cardId }}
        </el-form-item> -->
        <el-form-item :label="$t('lang.transaction_id')+':'">
          {{ transactionDetail.transactionId }}
        </el-form-item>
        <el-form-item :label="$t('lang.transaction_amount')+':'" class="text-nowrap">
          {{ `${transactionDetail.amount} ${transactionDetail.amountCurrency == null?'USD':transactionDetail.amountCurrency}` }}
        </el-form-item>
        <el-form-item :label="$t('lang.transaction_type')+':'">
          {{ $t(transactionDetail.transactionType) }}
        </el-form-item>
        <el-form-item :label="$t('lang.merchant_name')+':'">
          {{ transactionDetail.merchantName }}
        </el-form-item>
        <el-form-item :label="$t('lang.transaction_time')+':'">
          {{ formatDate(transactionDetail.transactionTime) }}
        </el-form-item>
      </el-form>

      </el-dialog>
    </div>
  </template>

<script>
export default {
  name: 'WalletHistory',
}
</script>
<script setup>
import {reactive, ref,onMounted} from 'vue'
import RechargeForm from '@/view/pixal/common/recharge_form.vue'
import {listWalletHistory,getCardTransactionRecord} from '@/api/finance'
import {formatDate} from '@/utils/format'
import {buildExcel} from '@/utils/excel'
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
const dialogs = reactive({
  rechargeDetailDialog:false,
  rechargeDialogVisible:false,
  viewOrderDetailDialogVisible:false,
})
const search = reactive({
  timeRange: [],
  transactionType: null,
  total: 0,
  pageSize: 10,
  page:1
});

const tableData = ref([])

const getTableData = ()=>{
  listWalletHistory(search).then(res => {
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
  // 处理导出逻辑
  console.log('点击了导出按钮');
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
const submitRecharge =() =>{

}
const handleFileChange = (file,fl) => {
  rechargeForm.value.voucher = file.raw
  console.log(rechargeForm.value.voucher)
  };

  const beforeUpload = (file) => {
  const isJPG = file.type === 'image/jpeg';
  const isPNG = file.type === 'image/png';
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isJPG && !isPNG) {
      ElMessage.error('上传图片只能是 JPG 或 PNG 格式!');
      return false;
  }
  if (!isLt2M) {
      ElMessage.error('上传图片大小不能超过 2MB!');
      return false;
  }
  return true;
  };
  const onExport = () => {
    var s = {
      ...search,
      pageSize:99999,
      page:1
    }
    listWalletHistory(s).then(res=>{ 
      if(res.code === 0){
        let list = res.data.list.map(x =>{
          return {
            
            "Order No":x.orderId,
            "Card No":x.cardNo,
            "Type":t(x.transactionType),
            "Transfer Amount":`${x.amount} ${x.amountCurrency}`,
            "Wallet Balance":`${x.balance} ${x.currency}`,
            "Update Time":formatDate(x.CreatedAt)

          }
        })
        buildExcel(list,"wallet_history")
      }
    })
  }
const transactionDetail = ref({})
const handleViewOrderDetail = (row)=>{
    transactionDetail.value = row.transactionRecord
    dialogs.viewOrderDetailDialogVisible = true
}
</script>

<style scoped>
.recharge-records {
  padding: 20px;
  background-color: #fff;
  border-radius: 24px;

}

.operation-bar {
  margin-bottom: 20px;
}

</style>
