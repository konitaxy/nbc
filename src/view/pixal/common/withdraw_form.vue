<template>
  <div class="d-flex justify-content-center align-items-center">
    <!-- 侧边步骤条 -->
    <div style="height: 150px; max-width: 100px; position: absolute; left: 50px; top: 100px" class="d-none d-sm-block">
      <el-steps direction="vertical" :active="step" finish-status="success">
        <el-step :title="$t('lang.withdrawal_info')" />
        <el-step :title="$t('lang.information_verification')" />
        <el-step :title="$t('lang.completed')" />
      </el-steps>
    </div>

    <el-card class="box-card w-100 mt-5" style="max-width: 550px;">
      <!-- 提现方式选择 -->
      <el-radio-group v-model="withdrawForm.accountType" class="w-100 mt-2">
        <el-radio-button disabled size="large" style="width: 50%;" :label="$t('lang.bank_account')" value="CARD" />
        <el-radio-button size="large" style="width: 50%;" :label="$t('lang.tron_address')" value="BLOCKCHAIN" />
      </el-radio-group>

      <!-- 区块链提现表单 (Step 1) -->
      <div v-if="withdrawForm.accountType === 'BLOCKCHAIN'" class="mt-4">
        <div v-if="step == 1">
          <!-- Tron 地址输入 -->
          <el-form-item
            :label="$t('lang.tron_address')"
            :validate-status="withdrawForm.accountNumber != null ? 'success' : 'error'"
            label-position="top"
            required
            inline-message
            :error="$t('lang.this_field_is_required')"
          >
            <el-input v-model="withdrawForm.accountNumber" :placeholder="$t('lang.please_enter_tron_address')" />
          </el-form-item>

          <!-- 提现金额与计算 -->
          <div class="p-2 bg-body-secondary rounded-3 mt-5">
            <el-form-item
              :label="$t('lang.withdrawal_amount')"
              :validate-status="withdrawForm.remitAmount > 0 ? 'success' : 'error'"
              label-position="top"
              required
              inline-message
              :error="$t('lang.this_field_is_required')"
            >
              <el-input @change="handleRemitAmountChange" v-model="withdrawForm.remitAmount" :placeholder="$t('lang.please_enter_amount')">
                <template #append>
                  <el-select v-model="withdrawForm.currency" style="width: auto;">
                    <el-option
                      v-for="item in currencyOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    >
                    </el-option>
                  </el-select>
                </template>
              </el-input>
              <div class="w-100">
                <p class="float-end">
                  <span style="font-size: small" class="text-secondary">
                    {{ $t('lang.available_balance') }} {{ userStore.userInfo.wallet.balance | 0 }} {{ withdrawForm.currency }}
                  </span>
                  <a
                    class="link-offset-2 link-underline link-underline-opacity-0 ms-2 text-primary"
                    href="#"
                    @click="withdrawForm.remitAmount = userStore.userInfo.wallet.balance; handleRemitAmountChange()"
                  >
                    {{ $t('lang.withdraw_all') }}
                  </a>
                </p>
              </div>
            </el-form-item>
          </div>

          <!-- 费率提示 -->
          <div class="position-relative">
            <el-divider direction="vertical" style="height: 100px;" class="ms-3 p-2"></el-divider>
            <div class="position-absolute hstack" style="top: 40%; left: 5px">
              <div
                class="text-dark bg-body-secondary fs-5 text-center"
                style="width: 22px; height: 22px; border-radius: 50%; transform: rotate(45deg);"
              >
                ➗
              </div>
              &nbsp;&nbsp;
              <span class="text-danger">{{ $t('lang.fee_rate') }} {{ withdrawForm.feeRate * 100 }}%</span>
            </div>
          </div>

          <!-- 您将收到金额 -->
          <div class="p-2 bg-body-secondary rounded-3">
            <el-form-item :label="$t('lang.you_will_receive')" label-position="top">
              <el-input readonly v-model="withdrawForm.amount" :placeholder="$t('lang.please_enter_tron_address')">
                <template #append>
                  <el-select v-model="withdrawForm.currency" style="width: auto;">
                    <el-option
                      v-for="item in currencyOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    >
                    </el-option>
                  </el-select>
                </template>
              </el-input>
            </el-form-item>
          </div>

          <!-- 按钮组 -->
          <div class="button-group">
            <el-button type="info" @click="goBack">{{ $t('lang.back') }}</el-button>
            <el-button
              type="primary"
              :disabled="withdrawForm.amount <= 0 || withdrawForm.accountNumber == ''"
              @click="step = 2"
            >
              {{ $t('lang.next_step') }}
            </el-button>
          </div>
        </div>

        <!-- 确认信息 (Step 2) -->
        <div v-if="step == 2">
          <h2 class="h5">{{ $t('lang.confirm_withdrawal_info') }}</h2>
          <el-card class="tron-address-card">
            {{ $t('lang.tron_address') }}: <span>{{ withdrawForm.accountNumber }}</span>
          </el-card>

          <div class="row">
            <!-- 信息时间线 -->
            <div class="col-sm-4 col-xs-12">
              <div class="p-4">
                <el-timeline style="max-width: 600px">
                  <el-timeline-item>
                    <div>
                      <p>
                        <span class="text-body-secondary">{{ $t('lang.you_withdraw') }}</span><br />
                        <span class="fw-medium fs-5">{{ withdrawForm.remitAmount }} {{ withdrawForm.currency }}</span><br />
                        <span style="font-size: small;" class="text-nowrap">
                          {{ $t('lang.wallet_balance') }} {{ userStore.userInfo.wallet.balance | 0 }} {{ withdrawForm.currency }}
                        </span>
                      </p>
                    </div>
                  </el-timeline-item>
                  <el-timeline-item>
                    <div>
                      <p>
                        <span class="text-body-secondary">{{ $t('lang.fee_rate') }}</span><br />
                        <span class="fw-medium fs-5">{{ withdrawForm.feeRate * 100 }}%</span>/<span class="text-bold">{{ $t('lang.per_transaction') }}</span>
                      </p>
                    </div>
                  </el-timeline-item>
                  <el-timeline-item>
                    <div>
                      <p>
                        <span class="text-body-secondary">{{ $t('lang.receiver_gets') }}</span><br />
                        <span class="fw-medium fs-5 text-nowrap">{{ withdrawForm.amount }} {{ withdrawForm.currency }}</span>
                      </p>
                    </div>
                  </el-timeline-item>
                </el-timeline>
              </div>
            </div>

            <!-- 交易附言输入 -->
            <div class="col-sm-7 pt-4 col-xs-12">
              <el-form :model="withdrawForm" label-position="top" :rules="rules" ref="formRef" class="ms-3">
                <el-form-item :label="$t('lang.transaction_memo')" prop="memo">
                  <el-input v-model="withdrawForm.memo" :placeholder="$t('lang.please_enter_transaction_memo')" />
                </el-form-item>

                <div class="button-group">
                  <el-button :loading="loading.withdrawApplyLoading" @click="step = 1">{{ $t('lang.previous_step') }}</el-button>
                  <el-button :loading="loading.withdrawApplyLoading" type="primary" @click="handleWithdrawConfirm">
                    {{ $t('lang.submit') }}
                  </el-button>
                </div>
              </el-form>
            </div>
          </div>
        </div>

        <!-- 提交成功 (Step 3) -->
        <div v-if="step == 3">
          <div>
            <el-alert type="success" center :closable="false">{{ $t('lang.withdrawal_request_submitted_successfully') }}</el-alert>
          </div>
          <el-button class="w-100 mt-4" @click="goBack">{{ $t('lang.return') }}</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>
  
<script>
  export default {
  name: 'WithdrawForm',
}
</script>
<script setup>
import { reactive, ref } from 'vue';
import {walletWithdrawApply} from '@/api/finance'
import { ElMessage } from 'element-plus'
import {writeText} from 'clipboard-polyfill'
import {useUserStore} from '@/pinia/modules/user'
const userStore = useUserStore()
const accountType = ref('链转账');
const amount = ref('');
const step = ref(1);
const withdrawForm = reactive({
  remitAmount: '0',
  feeRate: 0.00,
  amount: '0',
  currency:'USD',
  accountType:'BLOCKCHAIN',
  verifyCode:''
})
const resetForm = ()=>{
  withdrawForm.remitAmount = '0'
  withdrawForm.feeRate = 0.00
  withdrawForm.amount = '0'
  withdrawForm.currency='USD'
  withdrawForm.accountType='BLOCKCHAIN'
  withdrawForm.verifyCode=''
}
const handleRemitAmountChange = ()=>{
  withdrawForm.amount = (withdrawForm.remitAmount * (1-withdrawForm.feeRate)).toFixed(2)
}
const loading = reactive({
  withdrawApplyLoading:false
})
const handleCopy = (v) => {
  writeText(v)
  ElMessage.success('复制成功')
};
const goBack = () => {
  const escEvent = new KeyboardEvent('keydown', {
      key: 'Escape',
      code: 'Escape',
      keyCode: 27,
      which: 27,
      bubbles: true
    });
    document.dispatchEvent(escEvent)
};
const rechargeApplyResult = ref({
})
const handleWithdrawConfirm = () => {
  rechargeApplyResult.value = {}
  loading.withdrawApplyLoading = true
  if (withdrawForm.memo === '' || withdrawForm.accountNumber == '' ){
   ElMessage.error('Please fill the memo or account number')
   return
  }
  walletWithdrawApply(withdrawForm).then(res =>{
    if (res.code === 0){
      rechargeApplyResult.value = res.data
      ElMessage.success('申请成功')
      resetForm()
      step.value = 3
      userStore.GetUserInfo()
      // goBack()
    }
    loading.withdrawApplyLoading = false
  })
};
const reset = () => {
  rechargeApplyResult.value = {}
};

const currencyOptions = [{label:'USD',value:'USD'}]

</script>

<style scoped>
.container {
  padding: 20px;
}

.box-card {
  margin-top: 20px;
}

.chain-info {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.chain-info p {
  display: flex;
  align-items: center;
}

.chain-info span {
  margin-right: 10px;
}

</style>