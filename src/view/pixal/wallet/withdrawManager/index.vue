<template>
    <div class="recharge-records">
      <!-- 操作栏 -->
      <div class="operation-bar">
        <el-input v-model="search.orderNumber" :placeholder="$t('lang.enter_order_number')" style="max-width: 140px;"></el-input>
        <!-- <el-select v-model="search.transcationType" :placeholder="$t('lang.please_select_audit_status')" @change="handleSearch" style="width: 150px;" clearable>
          <el-option label="转账" :value="1"></el-option>
          <el-option label="收款" :value="2"></el-option>
        </el-select> -->
        <el-select v-model="search.auditStatus" :placeholder="$t('lang.please_select_audit_status')" @change="handleSearch" style="width: 150px;" clearable>
          <el-option :label="$t('lang.transfer_pending')" value="Pending"></el-option>
          <el-option :label="$t('lang.proceed')" value="Proceed"></el-option>
          <el-option :label="$t('lang.decline')" value="Decline"></el-option>
        </el-select>
        <el-button class="ms-2" type="primary" icon="search" @click="handleSearch">{{$t('lang.search')}}</el-button>
        <el-button type="primary" icon="plus" @click="dialogs.withdrawDialogVisible = true">{{$t('lang.transfer')}}</el-button>
        <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="handleExport">{{$t('lang.export')}}</el-button>
      </div>
  
      <!-- 充值记录表格 -->
      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="orderId" :label="$t('lang.order_number')" width="200"></el-table-column>
        <el-table-column prop="originAmount" :label="$t('lang.withdrawal_amount')"></el-table-column>
        <!-- <el-table-column prop="amount" :label="$t('lang.arrival_amount')">
          <template #default="scope">
            <span>{{scope.row.status =='Pending'?'':scope.row.amount}}</span>
          </template>
        </el-table-column> -->
        <el-table-column prop="CreatedAt" :label="$t('lang.withdrawal_time')" min-width="100">
          <template #default="scope">
            <span>{{formatDate(scope.row.CreatedAt)}}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('lang.arrival_time')" min-width="100">
          <template #default="scope">
            <span>{{scope.row.finishTime==null?'':formatDate(scope.row.finishTime)}}</span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="finishTime" :label="转出时间" min-width="100">
          <template #default="scope">
            <span>{{scope.row.finishTime?formatDate(scope.row.finishTime):''}}</span>
          </template>
        </el-table-column> -->
        <el-table-column prop="status" :label="$t('lang.audit_status')">
          <template #default="scope">
            <el-tag v-if="scope.row.status =='Pending'" type="info">{{scope.row.status}}</el-tag>
            <el-tag v-else-if="scope.row.status =='Proceed'" type="success">{{scope.row.status}}</el-tag>
            <el-tag v-else type="danger">{{scope.row.status}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('lang.operations')" width="150">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">{{ $t('lang.view') }}</el-button>
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
      <el-dialog title="详情" v-model="dialogs.withdrawDetailDialogVisible" width="30%"  align-center>
        <div class="d-flex flex-column gap-2 withdraw-detail">
            <p><span>单&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;号:</span> {{ currentRecord.orderId }}</p>
            <p><span>{{ $t('lang.audit_status') }}:</span> {{ currentRecord.status }}</p>
            <p><span>转账金额:</span> {{ `${currentRecord.originAmount} ${currentRecord.currency}` }}</p>
            <p><span>到账金额:</span> {{ `${currentRecord.amount} ${currentRecord.currency}` }}</p>
            <p><span>充值时间:</span> {{ formatDate(currentRecord.CreatedAt) }}</p>
            <p><span>到账时间:</span> {{ currentRecord.finishTime?formatDate(currentRecord.finishTime):'' }}</p>
            <p><span>充值方式:</span> {{ currentRecord.accountType }}</p>
            <p><span>提现账号:</span> {{ currentRecord.accountNumber }}</p>
            <p><span>提现备注:</span> {{ currentRecord.memo }}</p>
            <p><span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;审核人:</span> {{ currentRecord.operator }}</p>
            <p><span>审核备注:</span> {{ currentRecord.remark }}</p>

        </div>
        </el-dialog>
        <el-dialog
          :title="$t('lang.wallet_withdrawal')"
          v-model="dialogs.withdrawDialogVisible"
          fullscreen
          append-to-body
          class="recharge-form-dialog"
          destroy-on-close
        >
          <WithdrawForm/>
          
        </el-dialog>
    </div>
  </template>
<script>
export default {
  name: 'WithdrawManager',
}
</script>

<script setup>
import {reactive, ref,onMounted} from 'vue'
import WithdrawForm from '@/view/pixal/common/withdraw_form.vue';
import {listWalletWithdraw} from '@/api/finance'
import {formatDate} from '@/utils/format'
import {buildExcel} from '@/utils/excel'

const dialogs = reactive({
  withdrawDialogVisible:false,
  withdrawDetailDialogVisible:false
})
const search = reactive({
  status: null,
  transcationType: null,
  total: 0,
  pageSize: 10,
  page:1
});

const tableData = ref([])
const getTableData = () => {
  listWalletWithdraw(search).then(res => {
    if(res.code === 0){
      tableData.value = res.data.list
      search.total = res.data.total
    }
    
  })
}
onMounted(()=>{getTableData()})
  const handleSearch =() =>{
    getTableData()
  }
  const handleExport =() =>{
    var s = {
      ...search,
      pageSize:99999,
      page:1
    }
    listWalletWithdraw(s).then(res=>{ 
      if(res.code === 0){
        let list = res.data.list.map(x =>{
          return {
            "Order No":x.orderId,
            "Withdrawal Amount":parseFloat(x.originAmount),
            "Arrival Amount":parseFloat(x.amount),
            "Withdrawal Time":formatDate(x.CreatedAt),
            "Arrival Time":x.finishTime ==null?'':formatDate(x.finishTime),
            "Audit Status": x.status,
          }
        })
        buildExcel(list,"wallet_withdraw_history")
      }
    })
  }
  const currentRecord = ref({})
  const handleView =(val) => {
    currentRecord.value = val
    dialogs.withdrawDetailDialogVisible = true
  }
  const handleSizeChange=(val) => {
    search.pageSize = val
    getTableData()
  }
  const  handleCurrentChange =(val) => {
    search.page = val
    getTableData()
  }
  const submitRecharge =() =>{

  }
  const elUpladRef = ref(null);
  const handleFileChange = (file,fl) => {
    rechargeForm.value.voucher = file.raw
    console.log(rechargeForm.value.voucher)
    };
    const handleFileRemove =(file) =>{
        elUpladRef.value.clearFiles();
    }

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
  </script>
  
  <style scoped>
 .recharge-records {
  padding: 20px;
  background-color: #fff;
  border-radius: 24px;

}

  .operation-bar {
    margin-bottom: 20px;
    display: flex;
    /* display: inline; */
    gap: 10px;
  }
  .withdraw-detail span {
    font-size: 14px;
    width: 300px;
    font-weight: bold;
    margin-right: 10px;
  }
  </style>
