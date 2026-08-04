<template>
    <div class="card-management form-white">
      
        <el-form inline label-position="top">
            <el-form-item  prop="email">
                <el-input v-model="search.email" :placeholder="$t('lang.email')" clearable></el-input>
            </el-form-item>
            <el-form-item  prop="email">
                <el-input v-model="search.mobile" :placeholder="$t('lang.phone_number')" clearable></el-input>
            </el-form-item>
            <el-form-item label="">
                <el-button type="primary" icon="search" @click="getTableData">{{ $t('lang.search') }}</el-button>
            </el-form-item>
            <el-form-item label="">
                <el-button type="primary" @click="openAddCardHolderDialog" icon="plus" >{{ $t('lang.add_new_cardholder') }}</el-button>
            </el-form-item>
        </el-form>
        
<el-table :data="tableData" style="width: 100%" border show-overflow-tooltip flexible>
  
<el-table-column prop="firstName" :label="$t('lang.cardholder_first_name')" >
<template #default="{row}">
    {{ `${row.firstName} ${row.lastName}` }}
</template>
</el-table-column>

  <el-table-column prop="mobile" :label="$t('lang.phone_number')" >
    <template #default="{row}">
        {{ `${row.mobilePrefix} ${row.mobile}` }}
    </template>
  </el-table-column>
  <el-table-column prop="birthDate" :label="$t('lang.date_of_birth')" >
  </el-table-column>
  
  <!-- 表格列 - 卡状态 -->
  <el-table-column prop="email" :label="$t('lang.email')" >
  </el-table-column>
  <el-table-column prop="region" :label="$t('lang.country_or_region')" >
  </el-table-column>
  <el-table-column prop="cardCount" :label="$t('lang.card_count')" >
  </el-table-column>
  <el-table-column :label="$t('lang.actions')" width="120" fixed="right">
    <template #default="{ row }">
      <el-button type="primary" size="small" @click="handleOpenEditCardHolder(row)">{{ $t('lang.edit') }}</el-button>
    </template>
  </el-table-column>

</el-table>
      <el-pagination
        background
        layout="total, sizes, prev, pager, next"
        :total="search.total"
        :page-size="search.pageSize"
        :page-sizes="[10, 50, 100]"
        :pager-count="5"
        :current-page="search.page"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
    ></el-pagination>
   
  <el-dialog 
    v-model="dialogs.addCardHolderDialogVisible" 
    :title="$t('lang.add_cardholder')" 
    width="50%" 
    align-center>
    <el-form  label-width="auto">
      <el-form-item :label="$t('lang.nationality') + ' *'">
        <el-select v-model="cardHolder.region" :placeholder="$t('lang.please_select_nationality')" clearable style="min-width: 180px;" @change="onCardHolderRegionChange">
            <el-option :label="$t('lang.united_states')" value="USA"></el-option>
            <el-option :label="$t('lang.hong_kong')" value="HK"></el-option>
        </el-select>
      </el-form-item>

    <el-form-item :label="$t('lang.date_of_birth') + ' *'">
      <div>
        <el-date-picker 
          v-model="cardHolder.birthDate" 
          type="date" 
          value-format="YYYY-MM-DD" 
          style="width: 150px;" >
        </el-date-picker>
        <el-button @click="randomBirthDate" type="primary" size="small" class="ms-1">{{ $t('lang.random_birthdate') }}</el-button>
      </div>
    </el-form-item>
    
    <el-form-item :label="$t('lang.cardholder_last_name') + ' *'">
      <el-input v-model="cardHolder.lastName" :placeholder="$t('lang.please_enter_cardholder_last_name')"></el-input>
    </el-form-item>

      <el-form-item :label="$t('lang.cardholder_first_name') + ' *'">
        <el-input v-model="cardHolder.firstName" :placeholder="$t('lang.please_enter_cardholder_first_name')"></el-input>
      </el-form-item>

    <!-- 账单地国家 -->
    <el-form-item :label="$t('lang.billing_country') + ' *'">
      <el-select v-model="cardHolder.countryCode" :placeholder="$t('lang.please_select_country')" clearable style="min-width: 200px;">
          <el-option :label="$t('lang.united_states')" value="USA"></el-option>
          <el-option :label="$t('lang.hong_kong')" value="HK"></el-option>
      </el-select>
    </el-form-item>

      <!-- 账单地州省 -->
      <el-form-item :label="$t('lang.billing_state_province') + ' *'">
        <el-input v-model="cardHolder.state" :placeholder="$t('lang.please_enter_billing_state_province')"></el-input>
      </el-form-item>

      <!-- 账单城市 -->
      <el-form-item :label="$t('lang.billing_city') + ' *'">
        <el-input v-model="cardHolder.city" :placeholder="$t('lang.please_enter_billing_city')"></el-input>
      </el-form-item>
      
      <!-- 账单详细地址 -->
      <el-form-item :label="$t('lang.billing_address') + ' *'">
        <el-input v-model="cardHolder.address" :placeholder="$t('lang.please_enter_billing_address')"></el-input>
      </el-form-item>

      <!-- 账单地邮编 -->
      <el-form-item :label="$t('lang.billing_postcode') + ' *'">
        <el-input v-model="cardHolder.postcode" :placeholder="$t('lang.please_enter_billing_postcode')"></el-input>
      </el-form-item>

    <!-- 手机号码 -->
    <el-form-item :label="$t('lang.phone_number') + ' *'">
      <el-input v-model="cardHolder.mobile" :placeholder="$t('lang.please_enter_phone_number')">
        <template #prefix>
            <span class="pe-2">{{ cardHolder.mobilePrefix }}</span>
          </template>
      </el-input>
    </el-form-item>
    
    <el-form-item :label="$t('lang.email') + ' *'">
      <el-input v-model="cardHolder.email" :placeholder="$t('lang.please_enter_email')"></el-input>
    </el-form-item>
  </el-form>
  <template #footer>
    <div>
      <el-button type="success"  @click="randomCardHolder">{{ $t('lang.random_holder') }}</el-button>
      <el-button @click="dialogs.addCardHolderDialogVisible = false">{{ $t('lang.cancel') }}</el-button>
      <el-button 
        :loading="loading.addCardHolderLoading" 
        type="primary" 
        @click="handleAddCardHoderConfirm">
        {{ $t('lang.confirm') }}
      </el-button>
    </div>
  </template>
</el-dialog>

  <el-dialog
    v-model="dialogs.editCardHolderDialogVisible"
    :title="$t('lang.edit_cardholder')"
    width="50%"
    align-center>
    <el-form label-width="auto">
      <el-form-item :label="$t('lang.cardholder_last_name')">
        <el-input v-model="editCardHolderForm.lastName" disabled />
      </el-form-item>
      <el-form-item :label="$t('lang.cardholder_first_name')">
        <el-input v-model="editCardHolderForm.firstName" disabled />
        <div class="text-muted small mt-1">{{ $t('lang.cardholder_name_readonly_hint') }}</div>
      </el-form-item>

      <el-form-item :label="$t('lang.billing_country')">
        <el-select v-model="editCardHolderForm.countryCode" :placeholder="$t('lang.please_select_country')" clearable style="min-width: 200px;">
          <el-option :label="$t('lang.united_states')" value="USA"></el-option>
          <el-option :label="$t('lang.hong_kong')" value="HK"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('lang.billing_state_province')">
        <el-input v-model="editCardHolderForm.state" :placeholder="$t('lang.please_enter_billing_state_province')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('lang.billing_city')">
        <el-input v-model="editCardHolderForm.city" :placeholder="$t('lang.please_enter_billing_city')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('lang.billing_address')">
        <el-input v-model="editCardHolderForm.address" :placeholder="$t('lang.please_enter_billing_address')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('lang.billing_postcode')">
        <el-input v-model="editCardHolderForm.postcode" :placeholder="$t('lang.please_enter_billing_postcode')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('lang.phone_number')">
        <el-input v-model="editCardHolderForm.mobile" :placeholder="$t('lang.please_enter_phone_number')">
          <template #prefix>
            <span class="pe-2">{{ editCardHolderForm.mobilePrefix }}</span>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item :label="$t('lang.email')">
        <el-input v-model="editCardHolderForm.email" :placeholder="$t('lang.please_enter_email')"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <div>
        <el-button @click="dialogs.editCardHolderDialogVisible = false">{{ $t('lang.cancel') }}</el-button>
        <el-button
          :loading="loading.updateCardHolderLoading"
          type="primary"
          @click="handleUpdateCardHolderConfirm">
          {{ $t('lang.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

</div>
  </template>
  
  <script setup>
  import { reactive, ref,onMounted } from 'vue';
  import { ElMessage,ElMessageBox } from 'element-plus';
  import { formatDate,addYear} from '@/utils/format';
  import { listCardHolder,syncCard,addCardHolder,updateCardHolder,fetchCardHolderAddress,listCardBin,createCard,cancelCard,listCard,rechargeCard,withdrawCard } from '@/api/finance';
  import { useUserStore } from '@/pinia/modules/user'
  import CardDetail from './cardDetail.vue'
  import {buildExcel} from '@/utils/excel'
  import { buildCancelListPayload, applyCancelCardResult } from '@/utils/cancelCard'
  import { randomBirth, randomHKMobile } from '@/utils/random'
  import { useI18n } from 'vue-i18n'
  const { t } = useI18n()

  const userStore = useUserStore()
  const dialogs =reactive({
    activeCardDialogVisible:false,
    addCardHolderDialogVisible:false,
    editCardHolderDialogVisible:false,
    rechargeCardDialogVisible:false,
    withdrawCardDialogVisible:false,
    cardDetailDialogVisible:false,
  })
  const loading = reactive({
    addCardHolderLoading:false,
    updateCardHolderLoading:false,
    addCardLoading:false,
    rechargeCardLoading:false,
    withdrawCardLoading:false,
  })
  const activeTab = ref('normal')
  const filters = reactive({
    currency:'',
    brand:'',
    cardModel:'',
  })
  const filterParams = reactive({
    currency:[],
    brand:[],
    cardModel:[],
  })
  const cardList = ref([])
  const search = reactive({
    email:'',
    mobile:'',
    total: 0,
    pageSize: 10,
    page:1
  })
  
  const tableData = ref([])
  const holders = ref([])
  const holderSearch = reactive({
    page: 1,
    pageSize: 50,
  })
  
  const dialogTitle = ref('')
  const dialogType = ref('')
  const cardForm = ref({
    card: null,
    cardHolderId: '',
    cardBinId: '',
    cardBin:''
  })
  const getTableData = ()=>{
    listCardHolder(search).then(res => { 
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }
  onMounted(()=>{
    getTableData()
  })
  const handleRefresh = (row) => {
    row.loading = true
    syncCard({id:row.ID}).then(res => { 
      if (res.code ===0){
        handleListCard()
      } 
      row.loading = false
    }) 
  }
  const cardFormRef = ref(null)
  const cardBins = ref([])
  const handleListCardBin = () => {
    listCardBin(holderSearch).then((res) => {
      cardBins.value = res.data.list
      if(res.data.total > 0){
        filterParams.brand = [...new Set(res.data.list.map(item => item.cardBrand))]
        filterParams.currency = [...new Set(res.data.list.map(item => item.currency))]
        filterParams.cardModel = [...new Set(res.data.list.map(item => item.cardModel))]
        
      }
    })
  }
  const filterResult = (list) =>{
    cardList.value = cardBins.value.filter(card => {
      let match = true;

      if (filters.cardModel && card.cardModel !== filters.cardModel) match = false;
      if (filters.brand && card.cardBrand !== filters.brand) match = false;
      if (filters.currency && card.currency !== filters.currency) match = false;
      return match;
    })
    console.log(11123123,cardList.value)
  }
  
  const activeCardForm = ref({
    
  })
  const cardHolder = ref({
    region: 'USA',
    countryCode: 'USA',
    mobilePrefix: '+1',
  })
  const editCardHolderForm = ref({})
  const handleOpenEditCardHolder = (row) => {
    editCardHolderForm.value = {
      cardHolderId: row.cardHolderId,
      firstName: row.firstName,
      lastName: row.lastName,
      email: row.email,
      mobile: row.mobile,
      mobilePrefix: row.mobilePrefix || '+1',
      countryCode: row.countryCode === 'HKG' ? 'HK' : row.countryCode,
      state: row.state,
      city: row.city,
      postcode: row.postcode,
      address: row.address,
      region: row.region,
    }
    dialogs.editCardHolderDialogVisible = true
  }
  const handleUpdateCardHolderConfirm = () => {
    const form = editCardHolderForm.value
    const payload = {
      cardHolderId: form.cardHolderId,
      email: form.email,
      mobile: form.mobile,
      mobilePrefix: String(form.mobilePrefix || '').replace(/^\+/, ''),
      countryCode: form.countryCode,
      state: form.state,
      city: form.city,
      postcode: form.postcode,
      address: form.address,
      region: form.region,
    }
    loading.updateCardHolderLoading = true
    updateCardHolder(payload).then(res => {
      loading.updateCardHolderLoading = false
      if (res.code === 0) {
        ElMessage.success('Success')
        dialogs.editCardHolderDialogVisible = false
        getTableData()
      }
    })
  }
  const openAddCardHolderDialog = () => {
    cardHolder.value = {
      region: 'USA',
      countryCode: 'USA',
      mobilePrefix: '+1',
    }
    dialogs.addCardHolderDialogVisible = true
  }
  const onCardHolderRegionChange = (region) => {
    if (region === 'HK') {
      cardHolder.value.countryCode = 'HK'
      cardHolder.value.mobilePrefix = '+852'
    } else {
      cardHolder.value.region = 'USA'
      cardHolder.value.countryCode = 'USA'
      cardHolder.value.mobilePrefix = '+1'
    }
    // 切换身份后重新随机生成姓名与账单地址
    randomCardHolder()
  }
  const randomCardHolder = () => {
    const isHK = cardHolder.value.region === 'HK'
    const region = isHK ? 'hk' : 'us'
    fetchCardHolderAddress(region).then(res => {
      if (res.code !== 0 || !res.data) {
        return
      }
      const data = res.data
      cardHolder.value = {
        region: isHK ? 'HK' : 'USA',
        countryCode: isHK ? (data.countryCode === 'HKG' ? 'HK' : (data.countryCode || 'HK')) : (data.countryCode || 'USA'),
        firstName: data.firstName,
        lastName: data.lastName,
        email: data.email,
        mobilePrefix: isHK ? '+852' : (data.mobilePrefix || '+1'),
        mobile: isHK ? randomHKMobile() : data.mobile,
        birthDate: data.birthDate,
        state: data.state,
        city: data.city,
        postcode: data.postcode,
        address: data.address,
      }
    })
  }
  
  const handleAddCardHoderConfirm = ()=>{
    loading.addCardHolderLoading = true
    addCardHolder(cardHolder.value).then(res=>{
      loading.addCardHolderLoading = false
      if(res.code === 0){
        ElMessage.success('Success')
        dialogs.addCardHolderDialogVisible = false
        getTableData()
      }
    })
  }
  const randomBirthDate = ()=>{
    
    cardHolder.value.birthDate =  randomBirth();
  }
  const handleActiveCardConfirm = () => {
    loading.addCardLoading = true;
    const form = {
      cardBinId: cardForm.value.card.cardBinId,
      cardBin:cardForm.value.card.cardBin,
      cardHolderId: cardForm.value.cardHolderId,
      amount: cardForm.value.amount,
    }
    createCard(form).then(res =>{
      if(res.code === 0){
        ElMessage.success('Create Card Success')
        dialogs.activeCardDialogVisible = false
      }
      loading.addCardLoading = false;
    })
  }
  const cardView = ref({})
  const handleViewDetail = (row) =>{
    cardView.value = row
    dialogs.cardDetailDialogVisible = true
  }
  const quickOpenCards = () => {
    // 处理快速开卡逻辑
    console.log('点击了快速开卡按钮')
  };
  
  const handleActiveCard = () => {
    dialogs.activeCardDialogVisible = true
  }
  const batchFreezeCards = () => {
    // 处理批量冻结逻辑
    console.log('点击了批量冻结按钮')
  };
  const handleRechargeCardConfirm = ()=>{
    loading.rechargeCardLoading = true
    rechargeCard({
      id: currCard.value.ID,
      currency: userStore.userInfo.wallet.currency,
      amount: currCard.value.rechargeAmount
    }).then(res => {
      loading.rechargeCardLoading = false
      if(res.code === 0){
        ElMessage.success('Rechage success')
      }
      userStore.GetUserInfo()
    })
  }
  const handleWithdrawCardConfirm = ()=>{
    loading.withdrawCardLoading = true
    withdrawCard({
      id: currCard.value.ID,
      currency: userStore.userInfo.wallet.currency,
      amount: currCard.value.withdrawAmount
    }).then(res => {
      loading.withdrawCardLoading = false
      if(res.code === 0){
        ElMessage.success('Success')
      }
      userStore.GetUserInfo()
    })
  }
  
  const handleCancelCard = (row)=>{
    ElMessageBox.confirm(
    'The card will be cancelled immediately and cannot be restored!',
    'Warning',
    {
      confirmButtonText: 'Terminate',
      cancelButtonText: 'Quit',
      type: 'warning',
    }
  )
    .then(() => {
      handleCancelConfirm(row)
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'canceled',
      })
    })
  }
  const handleCancelConfirm = (row) => {
    let payload
    try {
      payload = buildCancelListPayload(row)
    } catch (e) {
      if (e.message === 'cancel_list_too_many') ElMessage.warning(t('lang.cancel_list_too_many'))
      else ElMessage.warning(t('lang.cancel_list_invalid'))
      return
    }
    cancelCard(payload).then((res) => {
      applyCancelCardResult(res, { t, ElMessage, onAfter: (ok) => { if (ok) getTableData() } })
    })
  }
  const batchUnfreezeCards = () => {
    // 处理批量解冻逻辑
    console.log('点击了批量解冻按钮')
  };
  
  const exportData = () => {
    // 处理导出逻辑
    console.log('点击了导出按钮')
  };
  
  const handleEdit = (index, row) => {
    dialogTitle.value = '编辑卡片';
    dialogType.value = 'edit';
    cardForm.value = { ...row };
    dialogVisible.value = true;
  };
  
  const handleDelete = (index, row) => {
    // 处理删除逻辑
    console.log('删除卡片:', row)
  };
  
  const handleSizeChange = (val) => {
    search.pageSize = val
    getTableData()
  }
  
const handleCurrentChange = (val) => {
  search.page = val
  getTableData()
};

const cardDetail = ref()
const currCard = ref()
const handleRechargeCard = (row) => {
  row.belongCardbin = cardBins.value.find(item => item.cardBinId === row.cardBinId)
  currCard.value = row
  dialogs.rechargeCardDialogVisible = true
};

const handleWithdrawCard = (row) => {
  currCard.value = row
  dialogs.withdrawCardDialogVisible = true

};


const handleCancel = (index, row) => {
  // 处理销卡逻辑
  console.log('销卡卡片:', row)
}
const onExport = () => {
    var s = {
      ...search,
      pageSize:99999,
      page:1
    }
    listCard(s).then(res=>{ 
      if(res.code === 0){
        let list = res.data.list.map(x =>{
          return {
            
            "Card No":x.cardNo,
            "Expiry Date":addYear(x.activeDate),
            "Status":x.cardStatus,
            "Balance":parseFloat(x.balance),
            "Creation Time":x.activeDate,
            "Currency":x.currency,
            "Brand":x.cardBrand,
          }
        })
        buildExcel(list,"card_list")
      }
    })
  }


</script>

<style lang="scss" scoped>
:deep(.el-radio-button__inner) {
  display: none !important;
}
.card-management {
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px 5px;
}
.form-actions {
  display: flex;
  flex-wrap: nowrap;
  overflow: auto;
  gap: 5px;
  .el-button {
    margin: 0px
  }
}

.balance-range {
  margin-bottom: 20px;
}
.el-select {
  max-width: 150px;
  // margin: 0 5px;
}
:deep(.el-input-group__prepend) {
  padding: 0px;
}

.card-bg {
  // transform: scale(0.);
  padding: 0px;
  height: 65px;
  width: 109px;
  border-radius: 8px;
  background-size: cover;
  margin-left: 10px;
}

.card-bg.is-active {
  border: 2px solid #01ad5a;
}

.card-MasterCard {
  background-image: url(/src/assets/mastercard.png);
}
.card-Visa {
  background-image: url(/src/assets/visacard.png);
}

.card-Discover {
  background-image: url(/src/assets/discovercard.png);
}

.el-input-number .el-input__inner {
  text-align: left ;
}

.note-bg  {
  border-radius: 16px;
}
</style>