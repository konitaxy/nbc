<template>
  <div class="">
    <!-- 中心钱包 -->
    <div class="row ">
      <div class="col col-12 page-title"><h5>My Wallet</h5></div>
      <div class="col col-md-7 col-sm-12">
        <el-card class="center-wallet" shadow="never">
          <div slot="header" class="d-flex justify-content-between">
            <span class="fw-medium fs-5 mb-2">{{ $t('lang.central_wallet') }}</span>
          </div>
          <div>
            <div class="row g-0">
                <div class="col col-8 p-0">
                  <div class="card-container p-8-3">
                  <div class="card-brand-mark">
                    <img class="card-logo" src="@/assets/NEWBEECARD-logo-black.png" alt="NEWBEECARD">
                    <div class="card-logo-text">NEWBEECARD</div>
                  </div>
                  <div class="currency-container d-flex a-center col-3" >
                    <img src="@/assets/flag-us.png" class="currency-icon">
                    <div class="currency-text" >{{userStore.userInfo.wallet.currency == ''?'USD':userStore.userInfo.wallet.currency }}</div>
                  </div>
                  <div class="balance-amount" >${{userStore.userInfo.wallet.balance}}  <div class="balance-text" >{{ $t('lang.available_balance') }}</div></div>
                  <div class="fw-bold" style="">
                    <p class="balance-text">${{ report.walletRechargeAmount }}</p>
                    <p class="balance-text">{{ $t('lang.total_recharge_amount') }}</p>
                  </div>
                </div>
              </div>
              <div class="col col-4 p-0">
  
                <div class="card-actions" v-if="$userStore.hasRole(7)">
                  <el-button type="primary" @click="dialogs.rechargeDialogVisible = true">{{ $t('lang.top_up') }}</el-button>
                </div>
                  <div class="card-actions" v-if="$userStore.hasRole(7)">
                  <el-button type="secondary button-dark" @click="dialogs.withdrawDialogVisible = true">{{ $t('lang.withdraw') }}</el-button>
                </div>
     
              </div>
              <el-card class="star-chain-card col-12 pd-0" body-style="padding:20px 0px;" shadow="never">
                <div slot="header" class="clearfix virtual-card-title">
                  <span>{{ $t('lang.virtual_card') }}</span>
                </div>
                <div class="metrics">
                    <el-progress  :width="80"type="circle" :percentage="getVisibleRiskRate(report.rechargeBackRate)" :color="riskColor" :format="() => formatRiskRate(report.rechargeBackRate)"></el-progress>
                    <div class=""><span>{{ $t('lang.chargeback_rate') }}<br/>{{ report.authorizationFailureCount }} / {{ report.authorizationCount }}</span></div>
                    <el-progress  :width="80"type="circle" :percentage="getVisibleRiskRate(report.refoundRate)" :color="riskColor" :format="() => formatRiskRate(report.refoundRate)"></el-progress>
                    <div><span>{{ $t('lang.refund_rate') }}<br/>{{ report.refundCount }} / {{ report.authorizationCount }}</span></div>
                    <el-progress  :width="80"type="circle" :percentage="getVisibleRiskRate(report.refoundAmountRate)" :color="riskColor" :format="() => formatRiskRate(report.refoundAmountRate)"></el-progress>
                    <div><span>{{ $t('lang.refund_amount_rate') }}<br/>{{ report.cardWithdrawAmount }} / {{ report.cardRechargeAmount }}</span></div>
                </div>
              </el-card>
            </div>
          </div>
        </el-card>
      </div>
      <div class="col col-md-5 col-sm-12">
        <el-card class="card-transaction-stats" shadow="never">
          <div slot="header" class="clearfix mb-3 font-weight-bold">
            <span>{{ $t('lang.virtual_card_transaction_summary') }}</span>
          </div>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-card shadow="never" class="summary-item">
                <div class="stat-item">
                  <div class="summary-item-title">
                    <el-icon><money /></el-icon>
                    <span>{{ $t('lang.total_balance') }}</span>
                  </div>
                    <span class="money">${{ report.cardTotalBalance}}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" class="summary-item">
                <div class="stat-item">
                  <div class="summary-item-title">
                    <el-icon><money /></el-icon>
                    <span>{{ $t('lang.recharge') }}</span>
                  </div>
                    <span class="money">${{ report.cardRechargeAmount}}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" class="summary-item">
                <div class="stat-item">
                  <div class="summary-item-title">
                  <el-icon><shopping-bag /></el-icon>
                  <span class="text-nowrap">{{ $t('lang.consume') }}</span>
                </div>
                  <span class="money">${{ report.authorizationAmount}}</span>
                </div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" class="summary-item">
                <div class="stat-item">
                  <div class="summary-item-title">
                    <el-icon><refresh-right /></el-icon>
                    <span>{{ $t('lang.refund') }}</span>
                  </div>
                    <span class="money">${{ report.cardWithdrawAmount}}</span>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
        <!-- 最新公告 -->
        <el-card class="latest-announcements" shadow="never">
          <div slot="header" class="clearfix">
            <span>{{ $t('lang.lastest_notice') }}</span>
            <el-button type="text"  @click="viewAllAnnouncements">View all</el-button>
          </div>
          <div class="announcement">
            <h3>...</h3>
            <p>none</p>
          </div>
        </el-card>
        <!-- 转账记录 -->
        <!-- <el-card class="transfer-record" shadow="never">
          <div slot="header" class="clearfix">
            <span>转账记录</span>
          </div>
          <el-form :model="transferForm" label-width="80px">
            <el-form-item label="收款码">
              <el-input v-model="transferForm.receiverCode"></el-input>
            </el-form-item>
            <el-form-item :label="$t('lang.amount')">
              <el-input v-model="transferForm.amount"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitTransfer">转账</el-button>
            </el-form-item>
          </el-form>
        </el-card> -->
        
      </div>
      <div class="col col-12 page-title"><h5>{{ $t('lang.transaction_history') }}</h5></div>
      <div class="col-12">
        <WalletHistory />
      </div>
</div>
<el-dialog
  :title="$t('lang.wallet_recharge')"
  v-model="dialogs.rechargeDialogVisible"
  class="recharge-form-dialog"
  fullscreen
  destroy-on-close
>
  <RechargeForm/>
</el-dialog>
<el-dialog
  :title="$t('lang.wallet_withdrawal')"
  v-model="dialogs.withdrawDialogVisible"
  destroy-on-close
  fullscreen
  class="recharge-form-dialog"
>
  <WithdrawForm/>
</el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref,onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import {writeText} from 'clipboard-polyfill'
import {useUserStore} from '@/pinia/modules/user'
import RechargeForm from '@/view/pixal/common/recharge_form.vue';
import WithdrawForm from '@/view/pixal/common/withdraw_form.vue';
import WalletHistory from '@/view/pixal/wallet/history.vue';
import {cardReport} from '@/api/finance'
const dialogs = reactive({
  rechargeDialogVisible: false,
  withdrawDialogVisible:false
})
const walletSum = ref({});
const userStore = useUserStore();
const transferForm = ref({
  receiverCode: '',
  amount: ''
});
const riskColor = (percentage) => {
  const value = Number(percentage)
  if (!Number.isFinite(value) || value <= 30) return '#7dffcc'
  if (value < 70) return '#ffd166'
  return '#ff5f7d'
}

const getRiskRate = (numerator, denominator) => {
  const total = Number(denominator)
  if (!Number.isFinite(total) || total <= 0) return 0
  const value = Number(numerator)
  if (!Number.isFinite(value) || value <= 0) return 0
  return Math.min(100, Number(((100 * value) / total).toFixed(2)))
}

const getVisibleRiskRate = (percentage) => {
  const value = Number(percentage)
  if (!Number.isFinite(value) || value <= 0) return 3
  return Math.min(100, value)
}

const formatRiskRate = (percentage) => {
  const value = Number(percentage)
  if (!Number.isFinite(value) || value <= 0) return '0%'
  return `${value}%`
}
const handleCopy = (v) => {
  writeText(v)
  Ellang.success('复制成功')
}

const report = ref({})
onMounted(()=>{
  // walletReport().then(res=>{
  //   if (res.code === 0){
  //     walletSum.value = res.data
  //   }
  // })
  cardReport({
  }).then(res=>{
    if (res.code === 0){
      report.value = res.data
      report.value.refoundRate = getRiskRate(report.value.refundCount, report.value.authorizationCount)
      report.value.refoundAmountRate = getRiskRate(report.value.cardWithdrawAmount, report.value.cardRechargeAmount)
      report.value.rechargeBackRate = getRiskRate(report.value.authorizationFailureCount, report.value.authorizationCount)
      // if(report.value.Card_Recharge_Success == null){
      //   report.value.Card_Recharge_Success = 0
      // }
      // if(report.value.Card_Recharge_count_Success == null){
      //   report.value.Card_Recharge_count_Success = 0
      // }
      // if(report.value.Card_Withdraw_count_Success == null){
      //   report.value.Card_Withdraw_count_Success = 0
      // }
      // report.value.totalCount = report.value.Card_Recharge_count_Success + report.value.Card_Withdraw_count_Success
      // if(report.value.Card_Withdraw_Success == null){
      //   report.value.Card_Withdraw_Success = 0
      // }
      
      // if (report.value.totalCount > 0){
      //   report.value.Card_Recharge_percent = (report.value.Card_Recharge_count_Success*100 / report.value.totalCount).toFixed(1)
      //   report.value.Card_Withdraw_percent = (report.value.Card_Withdraw_count_Success*100 / report.value.totalCount).toFixed(1)
      // }
      // if(report.value.Card_Recharge_Success > 0){
      //   report.value.Card_Withdraw_in_percent = (report.value.Card_Withdraw_Success*100/report.value.Card_Recharge_Success).toFixed(1)
      // }
      // if(report.value.Authorization_Success == null) {
      //   report.value.Authorization_Success = 0
      // }
      // if(report.value.Authorization_Failure == null) {
      //   report.value.Authorization_Failure = 0
      // }
      // if(report.value.Authorization_count_Success == null) {
      //   report.value.Authorization_count_Success = 0
      // }
      // if(report.value.Authorization_count_Failure == null) {
      //   report.value.Authorization_count_Failure = 0
      // }
      // report.value.Authorization_count_Total = report.value.Authorization_count_Success*1 + report.value.Authorization_count_Failure*1
      // if (report.value.Authorization_count_Total > 0){
      //   report.value.Authorization_count_Failure_percent = (report.value.Authorization_count_Failure*100/report.value.Authorization_count_Total).toFixed(1)
      // }
    }
  })
})
const viewAllAnnouncements = () => {
  // 处理查看所有公告的逻辑
};

const submitTransfer = () => {
  if (transferForm.value.receiverCode && transferForm.value.amount) {
    Ellang.success('转账成功');
    // 这里可以添加API请求代码
  } else {
    Ellang.error('请填写完整的转账信息');
  }
};
</script>

<style scoped>
@import '@/style/my-wallet.scss';
</style>
