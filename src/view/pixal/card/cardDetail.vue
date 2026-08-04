<template>
    <div>
    <el-tabs v-model="activeTab">
      <el-tab-pane :label="$t('lang.card_details')" name="cardDetails" class="pane-cardDetails">
          <div class="row mb-4">
            <div class="col-sm-6">
              <div class="card-container p-8-3">
                <div>
                  <div class="d-flex flex-row gap-3">
                    <h3 class="me-2" style="font-size: clamp(1.2rem, 18px, 2rem);">{{ card.cardNo }}<br></h3>
                    <el-button v-if="$userStore.hasRole(4)" class="button-dark" type="secondary" size="small" @click="handleShowCardDetail"><i class="bi bi-eye me-1"></i>{{ $t('lang.show_details') }}</el-button>
                  </div>
                  <span class="text-white">{{ $t('lang.card_number') }}</span>
                </div>

                <div class="d-flex flex-row gap-5 mt-3 mb-4">
                  <div class="fw-medium fw-5">
                    <div style="margin-left: 0px;width: 68px; height: 41px;" :class="`card-bg card-${card.cardBrand}-icon`"></div>
                  </div>
                  <!-- <div class="fw-medium fw-5">{{ $t('lang.brand') }}<br><span>{{ card.cardBrand }}</span></div> -->

                  <div class="fw-medium fw-5">{{ $t('lang.cvv') }}<br><span>{{ card.cvv }}</span></div>
                  <div class="fw-medium fw-5">{{ $t('lang.expiration_date') }}<br><span>{{ card.expirey}}</span></div>
                </div>
              </div>
            </div>
            <div class="col-sm-6 col-xs-12 card-actions-container">
              <el-button v-if="card.cardLevel !== 'SubCard'" type="primary" @click="handleRechargeCard()">{{ $t('lang.top_up') }}</el-button>
              
                <el-button type="danger" @click="handleCancelCard()">{{ $t('lang.terminate_card') }}</el-button>
                          <div class="mt-3">
                          <div class="mt-3 d-flex flex-wrap gap-2">
                <el-button type="primary" @click="handleCopyDetails">{{ $t('lang.copy_details') }}</el-button>
                <el-button type="success" :loading="loading.randomAddressLoading" @click="handleRandomAddress">{{ $t('lang.random_address') }}</el-button>
          </div>
              </div>
             
            </div>
          </div>

          <div class="row gy-3 fs-4  mb-5">
            <span class="fs-4 fw-medium">{{ $t('lang.card_statics') }}</span>
            <div class="col-4">
              <div class="data-item">
                <span class="fs-6 fw-medium">{{ $t('lang.total_card_recharge') }}</span><br>
                <span>${{ filterAmount('Card_Recharge') }}</span>
              </div>
            </div>
            <div class="col-4">
              <div class="data-item">
                <span class="fs-6 fw-medium">{{ $t('lang.total_authorization') }}</span><br>
                <span>${{ filterAmount('Authorization') }}</span>
              </div>
            </div>
            <div class="col-4">
              <div class="data-item">
                <span class="fs-6 fw-medium">{{ $t('lang.total_card_withdraw') }}</span><br>
                <span>${{ filterAmount('Card_Withdraw') }}</span>
              </div>
            </div>
            
          </div>

          <div class="row gy-3 fs-4  mb-5">
            <span class="fs-4 fw-medium">{{ $t('lang.basic_information') }}</span>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.card_id') }}</span><br>
              <span>{{ card.cardId }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.card_status') }}</span><br>
              <el-tag>{{ card.cardStatus }}</el-tag>
            </div>
            
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.available_balance') }}</span><br>
              <span>${{ card.balance }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.card_currency') }}</span><br>
              <span>{{ card.currency }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.card_schema') }}</span><br>
              <span>{{ card.cardBrand }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.valid_from') }}</span><br>
              <span>{{ addYear(card.activeDate) }}</span>
            </div>
            <div v-if="card.primaryCardNo" class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.belonging_primary_card') }}</span><br>
              <span>{{ card.primaryCardNo }}</span>
            </div>
            
          </div>
          <div class="row gy-3 fs-4 " v-if="card.Holder">
            <span class="fs-4 fw-medium">Billing Information</span>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.cardholder_first_name') }}</span><br>
              <span>{{ card.Holder.firstName }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.cardholder_last_name') }}</span><br>
              <span>{{ card.Holder.lastName }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.country') }}</span><br>
              <span>{{ card.Holder.countryCode }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.province_state') }}</span><br>
              <span>{{ card.Holder.state }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.city') }}</span><br>
              <span>{{ card.Holder.city }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.address') }}</span><br>
              <span>{{ card.Holder.address }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.postcode') }}</span><br>
              <span>{{ card.Holder.postcode }}</span>
            </div>
            <div class="col-6">
              <span class="fs-6 fw-medium">{{ $t('lang.billing_address') }}</span><br>
              <span>{{ `${card.Holder.countryCode}, ${card.Holder.state}, ${card.Holder.city}, ${card.Holder.address}`}}</span>
            </div>
            
          </div>

      </el-tab-pane>
      <el-tab-pane :label="$t('lang.transaction_details')" name="transactionDetails">
        <CardTransactionDetail :cardID="card.cardId"></CardTransactionDetail>
      </el-tab-pane>
    </el-tabs>
    <el-dialog
        :title="$t('lang.card_recharge')"
        width="50%"
        v-model="dialogs.rechargeCardDialogVisible"
        style="max-width: 500px;"
        align-center
        >
        <el-form  label-position="top" class="p-4">
          <el-form-item :label="$t('lang.card_number')">
            <el-input v-model="card.cardNo" disabled>
              <template #prepend>
                <span style="color: #1e73be;padding:0 10px; font-weight: bold;">{{ card.cardBrand }}</span>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item :label="$t('lang.card_currency')">
            <el-input v-model="card.currency" disabled></el-input>
          </el-form-item>

          <el-form-item :label="$t('lang.card_available_balance')">
            <el-input v-model="userStore.userInfo.wallet.balance" disabled>
            </el-input>
          </el-form-item>

          <el-form-item :label="$t('lang.recharge_amount')" required :error="$t('lang.please_enter_amount_min',(card.bin.minBalance,'USD'))" :validate-status="card.rechargeAmount>= card.bin.minBalance?'success':'error'">
            <el-input v-model="card.rechargeAmount" :placeholder="$t('lang.please_enter_amount')">
              <template #append>USD</template>
            </el-input>
            <PreRechargeSummary :summary="preRechargeSummary" :loading="preRechargeLoading" />
          </el-form-item>

        </el-form>
          <template #footer>
            <div class="d-flex justify-content-end">
              <el-button :loading="loading.rechargeCardLoading" type="info" @click="dialogs.rechargeCardDialogVisible = false">Cancel</el-button>
              <el-button :loading="loading.rechargeCardLoading" :disabled="preRechargeLoading" type="primary" @click="handleRechargeCardConfirm">Confirm</el-button>
            </div>
          </template>
      </el-dialog>
    </div>
  </template>

<script>
export default {
  name: 'CardDetail',
}
</script>
<script setup>
import { ref, reactive, defineProps, onMounted } from 'vue'
import { staticsCard, cancelCard, rechargeCard, showCardDetail, listCardBin, fetchCardHolderAddress } from '@/api/finance';
import { ElMessage, ElMessageBox } from 'element-plus';
import { buildCancelListPayload, applyCancelCardResult } from '@/utils/cancelCard'
import { useCardPreRecharge } from '@/composables/useCardPreRecharge'
import PreRechargeSummary from '@/components/card/PreRechargeSummary.vue'

import { formatDate, addYear } from '@/utils/format';
import CardTransactionDetail from './transactionDetails.vue'
import { useUserStore } from '@/pinia/modules/user'
import { useI18n } from 'vue-i18n';
import * as clipboard from 'clipboard-polyfill';
const userStore = useUserStore()
const { t } = useI18n();
const props = defineProps({
  card: {
    type: Object,
    default: () => {}
  }
})
const activeTab = ref('cardDetails')
const dialogs =reactive({
    rechargeCardDialogVisible:false,
  })
const loading = reactive({
  rechargeCardLoading:false,
  randomAddressLoading:false,
})
const card = ref(props.card)

const { summary: preRechargeSummary, loading: preRechargeLoading, reset: resetPreRecharge } = useCardPreRecharge({
  enabled: () => dialogs.rechargeCardDialogVisible,
  getCardId: () => card.value?.cardId,
  getRechargeAmount: () => card.value?.rechargeAmount,
})

const handleShowCardDetail = ()=>{
  showCardDetail({id: card.value.ID}).then(res =>{
      if(res.code == 0){
        card.value = res.data
      }
  })
}
const filterAmount = (type) =>{
  let s =statics.value.find(x =>x.transactionType == type)
  console.log(1111,s)
  return s==undefined?0:(s.amount *1).toFixed(2)
}
const statics = ref([])
onMounted(()=>{
  staticsCard({
    cardId:card.value.cardId
  }).then(res =>{
    if (res.code === 0){
      statics.value =res.data.list
    }
  })
  handleListCardBin()
})
const cardBins = ref([])
const handleListCardBin = () => {
    listCardBin({
    page: 1,
    pageSize: 50,
  }).then((res) => {
      cardBins.value = res.data.list
     
    })
  }
  const handleRechargeCard = ()=>{
    if (card.value.cardLevel === 'SubCard') {
      ElMessage.warning(t('lang.sub_card_cannot_recharge'))
      return
    }
    card.value.bin = cardBins.value.find(item => item.cardBinId === card.value.cardBinId)
    resetPreRecharge()
    dialogs.rechargeCardDialogVisible = true
  }
const handleRechargeCardConfirm = ()=>{
    if (!preRechargeSummary.value?.quotationRequestId) {
      ElMessage.warning(t('lang.pre_recharge_quote_required'))
      return
    }
    loading.rechargeCardLoading = true
    rechargeCard({
      id: card.value.ID,
      currency: userStore.userInfo.wallet.currency,
      amount: card.value.rechargeAmount,
      quotationRequestId: preRechargeSummary.value.quotationRequestId,
    }).then(res => {
      loading.rechargeCardLoading = false
      if(res.code === 0){
        ElMessage.success('Rechage success')
        dialogs.rechargeCardDialogVisible = false
        resetPreRecharge()
      }
      userStore.GetUserInfo()
    })
  }
  const handleCancelCard = () => {
    ElMessageBox.confirm(
    t('lang.cancel_card_warning'),
    t('lang.warning'),
    {
      confirmButtonText: t('lang.terminate_card'),
      cancelButtonText: t('lang.cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      handleCancelConfirm()
    })
    .catch(() => {})
  }
  const handleCancelConfirm = () => {
    const row = card.value
    if (!row?.ID || row.cardId == null || row.cardId === '') {
      ElMessage.warning(t('lang.cancel_list_invalid'))
      return
    }
    let payload
    try {
      payload = buildCancelListPayload(row)
    } catch (e) {
      if (e.message === 'cancel_list_too_many') ElMessage.warning(t('lang.cancel_list_too_many'))
      else ElMessage.warning(t('lang.cancel_list_invalid'))
      return
    }
    cancelCard(payload).then((res) => {
      applyCancelCardResult(res, { t, ElMessage })
    })
  }
  const copyDetailsWithHolder = (holder) => {
    if (!holder) {
      ElMessage.warning(t('lang.no_cardholder'))
      return
    }
    const cardNumber = card.value.cardNo
    const cvv = card.value.cvv
    const expiry = card.value.expirey
    const cardholder = `${holder.firstName || ''} ${holder.lastName || ''}`.trim()
    const billingParts = [holder.countryCode, holder.state, holder.city, holder.address, holder.postcode]
      .filter((part) => part != null && String(part).trim() !== '')
    const billingAddress = billingParts.join(', ')
    const details = `${cardNumber}\nCVV ${cvv} Expiration date ${expiry}\nHolder Name: ${cardholder}\nAddress: ${billingAddress}`
    clipboard.writeText(details).then(() => {
      ElMessage.success(t('lang.copy_success'))
    }).catch(() => {
      ElMessage.error(t('lang.copy_failed'))
    })
  }
  const handleCopyDetails = () => {
    copyDetailsWithHolder(card.value.Holder)
  }
  const resolveCardDzRegion = () => {
    const bin = cardBins.value.find(item => item.cardBinId === card.value.cardBinId)
      || card.value.bin
      || card.value.belongCardbin
    const r = String(bin?.region || card.value.Holder?.region || '').toUpperCase()
    if (r === 'HK' || r === 'HKG') return 'hk'
    if (r === 'CN' || r === 'CHN') return 'cn'
    return 'us'
  }
  const handleRandomAddress = () => {
    if (!card.value.Holder) {
      ElMessage.warning(t('lang.no_cardholder'))
      return
    }
    const region = resolveCardDzRegion()
    loading.randomAddressLoading = true
    fetchCardHolderAddress(region).then(res => {
      loading.randomAddressLoading = false
      if (res.code !== 0 || !res.data) {
        ElMessage.error(t('lang.random_address_failed'))
        return
      }
      const data = res.data
      const updatedHolder = {
        ...card.value.Holder,
        countryCode: data.countryCode || card.value.Holder.countryCode,
        state: data.state,
        city: data.city,
        address: data.address,
        postcode: data.postcode,
      }
      card.value.Holder = updatedHolder
      // 用刚生成的地址直接复制详情
      copyDetailsWithHolder(updatedHolder)
    }).catch(() => {
      loading.randomAddressLoading = false
      ElMessage.error(t('lang.random_address_failed'))
    })
  }
</script>
<style scoped>
.card-container {
  width: 100%;
  max-width: 354px;
  aspect-ratio: 1.586 / 1; /* 信用卡比例 */
  padding: clamp(1rem, 5%, 2rem); /* 响应式内边距 */
  background: linear-gradient(342deg, #385DEE 0%, #21EADF 100%);
  box-shadow: 0 1px 20px rgba(133, 137, 151, 0.2);
  border-radius: 16px;
  color: white;
  font-family: sans-serif;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}
.pane-cardDetails {
  padding: 1rem;
}
.data-item {
  background-color: #000;
  border-radius: 12px;
  padding: 1rem;
  color: #fff;
}

@media (max-width: 768px) {
  .pane-cardDetails {
    padding: 0;
  }
  .card-actions-container {
    padding: 1rem;
  }
}

</style>
