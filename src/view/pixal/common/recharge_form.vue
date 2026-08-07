<template>
    <div class="wallet-recharge-flow d-flex justify-content-center align-items-center">
        <div class="wallet-recharge-panel w-100 mt-5">
          <h3 class="recharge-section-title">{{$t('lang.recharge_type')}}</h3>
          <el-radio-group v-model="rechargeForm.rechargeType" class="recharge-type-tabs w-100 mt-2">
              
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
      
          <el-card class="box-card recharge-chain-card" v-if="rechargeForm.rechargeType === 'BLOCKCHAIN'">
              <el-alert
                v-if="hasRechargeApplyResult"
                :title="$t('lang.please_transfer_to_the_following_address')"
                type="info"
                :description="$t('lang.please_contact_customer_manager_for_more_information')"
                show-icon
                :closable="false"
                effect="light"
                class="recharge-alert p-4"
              >
            </el-alert>
            <div class="chain-info">
              <div class="chain-info-item chain-info-select">
                <span>{{ $t('lang.chain_type') }}</span
                ><br />
                <el-radio-group
                  v-model="rechargeForm.chain"
                  :disabled="hasRechargeApplyResult"
                  class="chain-network-tabs mt-1"
                  size="small"
                >
                  <el-radio-button value="TRC20">TRC20</el-radio-button>
                  <el-radio-button value="ERC20">ERC20</el-radio-button>
                </el-radio-group>
              </div>
              <div class="chain-info-item">
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
        <el-button class="button-dark" @click="goBack">{{ $t('lang.back') }}</el-button>
        <el-button
          :loading="loading.rechargeApplyLoading"
          type="primary"
          @click="handleRechargeApply"
          >{{ $t('lang.next_step') }}</el-button
        >
      </div>

      <!-- 充值结果卡片 (生成结果后) -->
      <el-card class="box-card recharge-result-card" v-if="hasRechargeApplyResult">
        <div class="d-flex flex-column justify-content-center align-items-center">
          <h3 class="result-title">{{ $t('lang.transfer_amount') }}</h3>
          <p class="fs-4 balance-view result-amount">
            <span class="currency-token-icon">{{ rechargeApplyResult.currency === 'USDT' ? '₮' : rechargeApplyResult.currency?.slice(0, 1) }}</span>
            <span class="result-amount-value">{{ rechargeApplyResult.remmitAmount }}</span>
            <span class="result-currency">{{ rechargeApplyResult.currency }}</span>
            <span class="result-network">{{ rechargeApplyResult.chain }}</span>
          </p>
          <div class="position-relative vstack ms-2 exact-amount-note">
            <p
              class="exact-amount-arrow"
              >&nbsp;</p
            >
            <div
              class="exact-amount-copy px-2 py-1 rounded-1 "
            >
              {{ $t('lang.please_transfer_exact_amount',{amount:rechargeApplyResult.remmitAmount,currency:rechargeApplyResult.currency,network:rechargeApplyResult.chain,minutes:240}) }}
              
            </div>
            <div class="countdown-wrap mt-2 d-flex justify-content-center align-items-center">
              <el-countdown title="Valid time:" :value="expireTime" />
            </div>
          </div>
          <img
            :src="rechargeApplyResult.qrCode"
            class="recharge-qr img-fluid mt-2"
          />
          <div class="address-card row rounded-3 p-2 mt-2">
            <div class="col-1">
              <i class="bi bi-exclamation-circle-fill fs-5"></i>
            </div>
            <div class="col-10">
              <p class="text-nowrap">{{ $t('lang.receiving_address') }} ({{ displayChainLabel }}):</p>
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
          <el-button class="button-dark" @click="reset">{{ $t('lang.reset_amount') }}</el-button>
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
  currency: 'USDT',
  rechargeType: 'BLOCKCHAIN',
  chain: 'TRC20',
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
const displayChainLabel = computed(() => {
  const chain = String(rechargeApplyResult.value?.chain || rechargeForm.chain || 'TRC20').toUpperCase()
  if (chain === 'ETHEREUM' || chain === 'ETH' || chain === 'ERC20') return 'ERC20'
  if (chain === 'TRON' || chain === 'TRC20') return 'TRC20'
  return rechargeForm.chain || 'TRC20'
})
const handleRechargeApply = () => {
  rechargeApplyResult.value = {}
  loading.rechargeApplyLoading = true
  walletRechargeApply({
    ...rechargeForm,
    chain: rechargeForm.chain || 'TRC20',
  }).then(res =>{
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
.wallet-recharge-flow {
  --console-bg: #020812;
  --console-panel: rgba(5, 16, 29, 0.86);
  --console-line: rgba(139, 214, 255, 0.16);
  --console-line-strong: rgba(139, 214, 255, 0.3);
  --console-text: #f4fbff;
  --console-muted: rgba(232, 247, 255, 0.84);
  --console-dim: rgba(232, 247, 255, 0.64);
  --console-cyan: #44d5ff;
  --console-blue: #2f7dff;
  --console-green: #7dffcc;
  --console-home-cta: linear-gradient(135deg, var(--console-green), var(--console-cyan) 58%, #83a8ff);
  color: var(--console-text);
}

.wallet-recharge-panel {
  max-width: 550px;
}

.recharge-section-title {
  margin: 0 0 8px;
  color: var(--console-text);
  font-size: 18px;
  font-weight: 800;
}

.box-card {
  margin-top: 20px;
  color: var(--console-text);
  border: 1px solid rgba(68, 213, 255, 0.28);
  border-radius: 24px;
  background:
    linear-gradient(145deg, rgba(68, 213, 255, 0.1), rgba(47, 125, 255, 0.05)),
    rgba(5, 16, 29, 0.78);
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.22);
  backdrop-filter: blur(14px);
}

.box-card :deep(.el-card__body) {
  color: var(--console-text);
  background: transparent;
}

.button-group :deep(.el-button) {
  min-height: 42px;
  padding: 11px 22px;
  color: var(--console-text);
  border-color: var(--console-line);
  border-radius: 14px;
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.06), rgba(255, 255, 255, 0.02)),
    rgba(2, 8, 18, 0.7);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.18);
  font-weight: 800;
}

.button-group :deep(.el-button:hover) {
  color: var(--console-text);
  border-color: var(--console-line-strong);
  background:
    linear-gradient(145deg, rgba(68, 213, 255, 0.12), rgba(47, 125, 255, 0.06)),
    rgba(2, 8, 18, 0.76);
}

.button-group :deep(.el-button--primary) {
  color: #061221;
  border-color: transparent;
  background: var(--console-home-cta);
  box-shadow: 0 18px 36px rgba(68, 213, 255, 0.18);
}

.chain-info {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-top: 22px;
}

.chain-info-item {
  min-width: 0;
  color: var(--console-text);
}

.chain-info span:first-child {
  color: var(--console-muted);
  font-size: 14px;
  font-weight: 700;
}

.chain-info .fw-bold {
  display: inline-block;
  margin-top: 4px;
  color: var(--console-text);
  font-size: 16px;
}

.chain-info-select {
  flex: 1;
}

.chain-network-tabs {
  display: inline-flex;
}

:deep(.chain-network-tabs .el-radio-button__inner) {
  min-width: 72px;
  color: var(--console-muted);
  border-color: rgba(139, 214, 255, 0.26);
  background: rgba(2, 8, 18, 0.52);
  font-weight: 800;
}

:deep(.chain-network-tabs .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  color: #061221;
  border-color: transparent;
  background: var(--console-home-cta);
  box-shadow: 0 10px 22px rgba(68, 213, 255, 0.16);
}

:deep(.chain-network-tabs .el-radio-button.is-disabled .el-radio-button__inner) {
  color: rgba(232, 247, 255, 0.5);
  border-color: rgba(139, 214, 255, 0.16);
  background: rgba(255, 255, 255, 0.045);
}

.button-group {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 22px;
}

.recharge-result-card {
  text-align: center;
}

.result-title {
  color: var(--console-text);
  font-size: 18px;
  font-weight: 800;
}

.result-amount {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #7dffcc;
  font-weight: 900;
}

.currency-token-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  color: #022216;
  border-radius: 50%;
  background: linear-gradient(135deg, #7dffcc, #44d5ff);
  box-shadow: 0 10px 24px rgba(125, 255, 204, 0.22);
  font-size: 18px;
  font-weight: 900;
  line-height: 1;
}

.result-amount-value,
.result-currency {
  color: #7dffcc;
}

.result-network {
  padding: 4px 9px;
  color: rgba(232, 247, 255, 0.86);
  border: 1px solid rgba(125, 255, 204, 0.2);
  border-radius: 999px;
  background: rgba(125, 255, 204, 0.08);
  font-size: 14px;
  font-weight: 800;
}

.exact-amount-note {
  display: flex;
  align-items: center;
}

.exact-amount-arrow {
  position: relative;
  bottom: -5px;
  width: 10px;
  height: 10px;
  margin: 0;
  transform: rotate(45deg);
  background: rgba(255, 111, 145, 0.92);
}

.exact-amount-copy {
  position: relative;
  left: 4px;
  color: #fff;
  background: linear-gradient(135deg, rgba(255, 111, 145, 0.96), rgba(255, 77, 109, 0.86));
  box-shadow: 0 12px 26px rgba(255, 77, 109, 0.18);
}

.address-card {
  color: var(--console-text);
  border: 1px solid rgba(139, 214, 255, 0.18);
  background: rgba(2, 8, 18, 0.34);
}

.countdown-wrap {
  color: var(--console-text);
  line-height: 1;
}

.countdown-wrap :deep(.el-countdown__title) {
  margin: 0 8px 0 0;
  color: var(--console-text) !important;
  font-size: 13px;
  font-weight: 800;
  line-height: 1;
}

.countdown-wrap :deep(.el-countdown) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.countdown-wrap :deep(.el-countdown__value) {
  color: #7dffcc !important;
  font-size: 20px;
  font-weight: 900;
  letter-spacing: 0;
  text-shadow: 0 0 18px rgba(125, 255, 204, 0.28);
}

.address-card {
  width: 100%;
  text-align: left;
}

.address-card i {
  color: var(--console-cyan);
}

.address-card p {
  margin-bottom: 4px;
  color: var(--console-muted);
}

.address-card p:last-child {
  overflow: hidden;
  color: var(--console-text);
  font-weight: 700;
  text-overflow: ellipsis;
}

.recharge-qr {
  width: 150px;
  padding: 10px;
  border: 1px solid rgba(139, 214, 255, 0.24);
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 18px 42px rgba(0, 0, 0, 0.22);
}

:deep(.recharge-type-tabs .el-radio-button) {
  width: 50%;
}

:deep(.recharge-type-tabs .el-radio-button__inner) {
  width: 100%;
  color: var(--console-muted);
  border-color: rgba(139, 214, 255, 0.26);
  background: rgba(2, 8, 18, 0.52);
  box-shadow: none;
}

:deep(.recharge-type-tabs .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  color: #061221;
  border-color: transparent;
  background: var(--console-home-cta);
  box-shadow: 0 14px 30px rgba(68, 213, 255, 0.18);
}

:deep(.recharge-type-tabs .el-radio-button.is-disabled .el-radio-button__inner) {
  color: rgba(232, 247, 255, 0.44);
  border-color: rgba(139, 214, 255, 0.16);
  background: rgba(255, 255, 255, 0.045);
}

:deep(.recharge-alert.el-alert) {
  border: 1px solid rgba(68, 213, 255, 0.28);
  border-radius: 18px;
  background:
    linear-gradient(135deg, rgba(68, 213, 255, 0.16), rgba(47, 125, 255, 0.08)),
    rgba(2, 8, 18, 0.42);
}

:deep(.recharge-alert .el-alert__title) {
  color: var(--console-cyan);
  font-weight: 800;
}

:deep(.recharge-alert .el-alert__description),
:deep(.el-form-item__label) {
  color: var(--console-muted);
}

@media (max-width: 768px) {
  .wallet-recharge-panel {
    max-width: 100%;
  }

  .chain-info {
    flex-direction: column;
    gap: 12px;
  }

  .button-group {
    flex-direction: column-reverse;
  }

  .button-group :deep(.el-button) {
    width: 100%;
    margin-left: 0;
  }
}
</style>
