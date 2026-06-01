<template>
    <div class="d-flex justify-content-center align-items-center">
        <div class="w-100 mt-5" style="max-width: 550px;">
          <h3>{{$t('lang.recharge_type')}}</h3>
          <el-radio-group v-model="rechargeForm.rechargeType" class="w-100 mt-2">
              
            <el-radio-button
          size="large"
          style="width: 50%;"
          :label="$t('lang.chain_transfer')"
          value="BLOCKCHAIN"
        />
        <el-radio-button
          disabled
          size="large"
          style="width: 50%;"
          :label="$t('lang.through_va_settlement')"
          value="CARD"
        />
            
          </el-radio-group>
      
          <el-card class="box-card" v-if="rechargeForm.rechargeType === 'BLOCKCHAIN'">
              <el-alert
                v-if="hasRechargeApplyResult"
                :title="$t('lang.please_transfer_to_the_following_address')"
                type="info"
                :description="$t('lang.please_contact_customer_manager_for_more_information')"
                show-icon
                :closable="false"
                effect="light"
                class="p-4 text-primary"
              >
            </el-alert>
            <div class="chain-info">
              <div>
                <span>{{ $t('lang.chain_type') }}</span
                ><br />
                <span class="fw-bold">TRC20</span>
              </div>
              <div>
                <span>{{ $t('lang.currency') }}</span
                ><br />
                <span class="fw-bold">USDT</span>
              </div>
            </div>
          </el-card>
      
          <el-form-item
        required
        inline-message
        :validate-status="rechargeForm.amount > 0 ? 'success' : 'error'"
        :error="$t('lang.this_field_is_required')"
        :label="$t('lang.amount')"
        label-position="top"
        label-class="fw-bold"
      >
        <el-input
          :readonly="hasRechargeApplyResult"
          v-model="rechargeForm.amount"
          :placeholder="$t('lang.please_enter_amount')"
          clearable
        ></el-input>
      </el-form-item>

      <!-- 操作按钮组 (未生成结果时) -->
      <div class="button-group" v-if="!hasRechargeApplyResult">
        <el-button @click="goBack">{{ $t('lang.back') }}</el-button>
        <el-button
          :loading="loading.rechargeApplyLoading"
          type="primary"
          @click="handleRechargeApply"
          >{{ $t('lang.next_step') }}</el-button
        >
      </div>

      <!-- 充值结果卡片 (生成结果后) -->
      <el-card class="box-card" v-if="hasRechargeApplyResult">
        <div class="d-flex flex-column justify-content-center align-items-center">
          <h3>{{ $t('lang.transfer_amount') }}</h3>
          <p class="fs-4 text-danger balance-view">
            {{ rechargeApplyResult.remmitAmount }} {{ rechargeApplyResult.currency }} {{ rechargeApplyResult.chain }}
          </p>
          <div style="display: contents;" class="position-relative vstack ms-2">
            <p
              style="position: relative; bottom: -5px; height: 10px; width: 10px; transform: rotate(45deg);"
              class="bg-danger"
              >&nbsp;</p
            >
            <div
              style="left: 4px"
              class="bg-danger px-2 py-1 text-light rounded-1 "
            >
              {{ $t('lang.please_transfer_exact_amount',{amount:rechargeApplyResult.remmitAmount,currency:rechargeApplyResult.currency,network:rechargeApplyResult.chain,minutes:240}) }}
              
            </div>
            <div class="mt-2 p-2 d-flex flex-column justify-content-center align-items-center">
              <el-countdown title="Valid time:" :value="expireTime" />
            </div>
          </div>
          <img
            :src="rechargeApplyResult.qrCode"
            style="width: 150px;"
            class="img-fluid img-thumbnail mt-2"
          />
          <div class="row bg-body-secondary rounded-3 p-2 mt-2">
            <div class="col-1">
              <i class="bi bi-exclamation-circle-fill fs-5 text-primary"></i>
            </div>
            <div class="col-10">
              <p class="text-nowrap">{{ $t('lang.receiving_address') }} (TRC20):</p>
              <p
                @click="handleCopy(rechargeApplyResult.accountNumber)"
                class="text-nowrap"
                style="cursor: pointer;"
              >
                {{ rechargeApplyResult.accountNumber }}
              </p>
            </div>
          </div>
        </div>
        <div class="button-group">
          <el-button type="primary" @click="reset">{{ $t('lang.reset_amount') }}</el-button>
          <el-button type="primary" @click="goBack">{{ $t('lang.i_have_completed_deposit') }}</el-button>
        </div>
      </el-card>
      </div>
    </div>
  </template>
  
<script>
  export default {
  name: 'RechargeForm',
}
</script>
<script setup>
import { reactive, ref, computed } from 'vue';
import {walletRechargeApply} from '@/api/finance'
import { ElMessage } from 'element-plus'
import {writeText} from 'clipboard-polyfill'
const rechargeType = ref('链转账');
const amount = ref('');
const rechargeForm = reactive({
  amount: '0',
  currency:'USDT',
  rechargeType:'BLOCKCHAIN'
})
const loading = reactive({
  rechargeApplyLoading:false
})
const handleCopy = (v) => {
  writeText(v)
  ElMessage.success('复制成功')
};
const expireTime = ref()
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
/** 接口返回 orderId；旧版可能为 traceId */
const hasRechargeApplyResult = computed(() => {
  const r = rechargeApplyResult.value
  return !!(r?.orderId || r?.traceId)
})
const handleRechargeApply = () => {
  rechargeApplyResult.value = {}
  loading.rechargeApplyLoading = true
  walletRechargeApply(rechargeForm).then(res =>{
    if (res.code === 0){
      rechargeApplyResult.value = res.data
      const expireAt = res.data.expireAtUnix
        ? new Date(res.data.expireAtUnix * 1000)
        : new Date(String(res.data.expireTime || '').replace(' ', 'T'))
      expireTime.value = expireAt
    }
    loading.rechargeApplyLoading = false
  })
};
const reset = () => {
  rechargeApplyResult.value = {}
};
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

.button-group {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>