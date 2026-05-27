<template>
  <div v-if="loading" class="pre-recharge-summary is-loading">
    {{ $t('lang.pre_recharge_loading') }}
  </div>
  <div v-else-if="summary" class="pre-recharge-summary">
    <div class="pre-recharge-final">
      {{ $t('lang.pre_recharge_final_amount') }}:
      <span class="amount">{{ summary.rechargeAmount }} {{ summary.rechargeCurrency || 'USD' }}</span>
    </div>
    <div class="pre-recharge-detail">
      {{ $t('lang.receive_amount') }} {{ summary.arrivalAmount }} {{ summary.arrivalAmountCurrency || 'USD' }}
      <template v-if="hasFee">
        · {{ $t('lang.pre_recharge_fee') }} {{ summary.rechargeFee }} {{ summary.rechargeFeeCurrency || 'USD' }}
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  summary: { type: Object, default: null },
  loading: { type: Boolean, default: false },
})

const hasFee = computed(() => {
  const fee = parseFloat(props.summary?.rechargeFee)
  return !Number.isNaN(fee) && fee > 0
})
</script>

<style scoped>
.pre-recharge-summary {
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.5;
  color: #606266;
}
.pre-recharge-summary.is-loading {
  color: #909399;
}
.pre-recharge-final .amount {
  font-weight: 600;
  color: #303133;
}
.pre-recharge-detail {
  margin-top: 2px;
  color: #909399;
  font-size: 12px;
}
</style>
