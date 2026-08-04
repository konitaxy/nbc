<template>
    <div class="tx-detail-panel">
        <el-form :inline="true" class="demo-form-inline">
            <el-form-item>
                <el-input v-model="search.cardNoSuffix" :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')" clearable></el-input>
            </el-form-item>
            <el-form-item>
                <el-input v-model="search.transactionId" :placeholder="$t('lang.transaction_id')" clearable></el-input>
            </el-form-item>
            <el-form-item >
                <el-date-picker v-model="search.timeRange" type="datetimerange" style="width: auto;"
                    :start-placeholder="$t('lang.start_time')"
                    value-format="YYYY-MM-DD"
                    :end-placeholder="$t('lang.end_time')"></el-date-picker>
            </el-form-item>
            <el-form-item>
                <el-select v-model="search.transactionType" :placeholder="$t('lang.transaction_type')" clearable style="min-width: 180px;">
                    <el-option :label="$t('Recharge')" value="Card_Recharge"></el-option>
                    <el-option :label="$t('Withdraw')" value="Card_Withdraw"></el-option>
                    <el-option :label="$t('Authorization')" value="Authorization"></el-option>
                    <el-option :label="$t('Settlement')" value="Settlement"></el-option>

                </el-select>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="onSearch">{{$t('lang.search')}}</el-button>
                <el-button v-if="$userStore.hasRole(6)" type="primary" @click="onExport">{{$t('lang.export')}}</el-button>
            </el-form-item>
            </el-form>
  
            <el-table :data="tableData" style="width: 100%" flexible @row-click="handleRowClick">
  <!-- 卡号 -->
  <el-table-column 
    min-width="145"
    prop="card.cardNo" 
    :label="$t('lang.card_number')">
  </el-table-column>

  <!-- 商户名称 -->
  <el-table-column show-overflow-tooltip prop="merchantName" 
  :label="$t('lang.merchant_name')">
  </el-table-column>

  <!-- 交易金额 -->
  <el-table-column :label="$t('lang.transaction_amount')">
    <template #default="{row}">
      {{ row.amount }}<span class="ms-1">{{row.currency}}</span>
    </template>
  </el-table-column>

  <el-table-column :label="$t('lang.transaction_type')">
    <template #default="{row}">
      <div>
        <p style="line-height: 18px;">{{$t(row.transactionType)}}</p>
        <p class="hstack"><span :class="row.status=='Success'?'bg-success':'bg-danger'" style="width: 6px;height: 6px;border-radius: 50%;margin-right: 3px;"></span><span :class="row.status=='Success'?'text-success':'text-danger'" style="font-size: 11px;line-height: 8px;">{{ row.status }}</span></p>
      </div>
    </template>
  </el-table-column>

  <!-- <el-table-column 
    prop="status" 
    :label="$t('lang.transaction_status')">
    <template #default="{ row }">
      <el-tag :type="row.status=='Success'?'success':'danger'">{{ row.status }}</el-tag>

    </template>
  </el-table-column> -->

  <el-table-column :label="$t('lang.transaction_id')" width="220" >
      <template #default="{row}">
      <div class="vstack justify-content-center text-end text-nowrap" style="display: inline-block;">
        <p class="" style="line-height: 18px;"><span>{{$t(row.transactionId)}}</span></p>
        <p style="line-height: 11px;"><span style="font-size: 11px;line-height: 8px;">{{ formatDate(row.transactionTime)  }}</span></p>
      </div>
    </template>
  </el-table-column>

  <!-- 操作 (如果需要，可以取消注释) -->
  <!--
  <el-table-column :label="$t('lang.operation')">
    <template #default="scope">
      <el-button @click="handleView(scope.$index, scope.row)">{{ $t('lang.view') }}</el-button>
    </template>
  </el-table-column>
  -->
</el-table>
  
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
  
      <el-drawer v-model="drawer" :title="$t('lang.detail')" :size="400">
        <div class="col-12 ps-4">
              <div class="card-container p-8-3">
                <div>
                  <div class="d-flex flex-row gap-3">
                    <h3 class="me-2" style="font-size: clamp(1.2rem, 18px, 2rem);">{{ drawerRow.card.cardNo }}<br></h3>
                    <el-button v-if="$userStore.hasRole(4)" class="button-dark" type="secondary" size="small" @click="handleShowCardDetail"><i class="bi bi-eye me-1"></i>{{ $t('lang.show_details') }}</el-button>
                  </div>
                  <span class="text-white">{{ $t('lang.card_number') }}</span>
                </div>

                <div class="d-flex flex-row gap-5 mt-3 mb-4">
                  <div class="fw-medium fw-5">
                    <div style="margin-left: 0px;width: 68px; height: 41px;" :class="`card-bg card-${drawerRow.card.cardBrand}-icon`"></div>
                  </div>
                  <!-- <div class="fw-medium fw-5">{{ $t('lang.brand') }}<br><span>{{ card.cardBrand }}</span></div> -->

                  <div class="fw-medium fw-5">{{ $t('lang.cvv') }}<br><span>{{ drawerRow.card.cvv }}</span></div>
                  <div class="fw-medium fw-5">{{ $t('lang.expiration_date') }}<br><span>{{ drawerRow.card.expirey.substring(2, 4) }}/{{ drawerRow.card.expirey.substring(0, 2) }}</span></div>
                </div>
              </div>
            </div>
            <div class="p-3 mt-5">
              <p class="fs-4 mb-2">{{ $t(drawerRow.transactionType) }} <el-tag size="small" :type="drawerRow.status =='Success'?'success':'danger'">{{ drawerRow.status }}</el-tag></p>
              <el-form label-position="top">
                <el-form-item :label="$t('lang.merchant_name')">
                  <p class="text-secondary">{{ drawerRow.merchantName?drawerRow.merchantName:'-' }}</p>
                </el-form-item>
                <el-form-item label="Transaction/Billing Amount">
                  <p class="text-secondary">{{ `${drawerRow.originAmount} ${drawerRow.originCurrency} ` }}/{{ ` ${drawerRow.amount} ${drawerRow.currency}` }}</p>
                </el-form-item>
              </el-form>
              <el-card body-class="bg-body-secondary">
                <el-form label-position="top" inline>
                <el-form-item :label="$t('lang.card_id')">
                  <p>{{ drawerRow.card.cardId }}</p>
                </el-form-item>
                <el-form-item :label="$t('Fee')">
                  {{ drawerRow.fee }}
                </el-form-item>
                <el-form-item :label="$t('lang.transaction_id')">
                  {{ drawerRow.transactionId }}
                </el-form-item>
                <el-form-item v-if="drawerRow.failReason" :label="$t('lang.fail_reason')">
                  {{ drawerRow.failReason }}
                </el-form-item>
                <el-form-item v-if="drawerRow.authCode" :label="$t('lang.auth_code')">
                  {{ drawerRow.authCode }}
                </el-form-item>
              </el-form>
              </el-card>
            </div>
    </el-drawer>
    </div>
  </template>
  
  <script>
 export default {
  name: 'CardTransactionDetail',
}
  </script>
  
  <script setup>
  import { ref,reactive,onMounted,defineProps } from 'vue';
  import {listCardTransactionRecord,showCardDetail} from '@/api/finance';
  import {formatDate} from '@/utils/format'
  import {buildExcel} from '@/utils/excel'
  import { useI18n } from 'vue-i18n';
  const { t } = useI18n();
  const props = defineProps({
    cardID: {
      type: String,
      default: ""
    }
  })
  const drawer = ref(false);
  const search = reactive({
    cardId: props.cardID,
    cardNumber: '',
    timeRange: [],
    transactionType: '',
    transactionStatus: '',
    total: 0,
    pageSize: 20,
    page:1
  });

  const getTableData = async()=>{
    
    listCardTransactionRecord(search).then(res=>{ 
      if(res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }
  onMounted(()=>{
    getTableData();
  })
  
  const tableData = ref([]);
  
  
  const onSearch = () => {
     getTableData();
  };
  
  const onExport = () => {
    var s = {
      ...search,
      pageSize:99999,
      page:1
    }
    listCardTransactionRecord(s).then(res=>{ 
      if(res.code === 0){
        let list = res.data.list.map(x =>{
          return {
            
            "Card No":x.card.cardNo,
            "Merchant Name":x.merchantName,
            "Transaction Amount":parseFloat(x.amount),
            "Currency":x.currency,
            "Transaction Type":t(x.transactionType),
            "Transaction ID":t(x.transactionId),
            "Trabsaction Status":x.status,
            "Transaction Time":formatDate(x.transactionTime),
          }
        })
        buildExcel(list,"transaction_details")

      }
    })
  }
  const handleSizeChange = (val) => {
    search.pageSize = val
    getTableData()
  }
  const drawerRow = ref()
  const handleRowClick = (row) =>{
    drawerRow.value = row
    drawer.value = true
  }
  const handleShowCardDetail = ()=>{
  showCardDetail({id: drawerRow.value.card.ID}).then(res =>{
      if(res.code == 0){
        drawerRow.value.card = res.data
      }
  })
}
  
const handleCurrentChange = (val) => {
  search.page = val
  getTableData()
};

  </script>
  
  <style scoped>
  .tx-detail-panel {
    padding: 0;
    background: transparent;
    border-radius: 0;
    color: inherit;
  }
  .tx-detail-panel :deep(.el-table),
  .tx-detail-panel :deep(.el-table__expanded-cell),
  .tx-detail-panel :deep(.el-table__body-wrapper),
  .tx-detail-panel :deep(.el-table__header-wrapper),
  .tx-detail-panel :deep(.el-table__inner-wrapper),
  .tx-detail-panel :deep(.el-table__body),
  .tx-detail-panel :deep(.el-table__header),
  .tx-detail-panel :deep(.el-table__empty-block),
  .tx-detail-panel :deep(.el-table th),
  .tx-detail-panel :deep(.el-table tr),
  .tx-detail-panel :deep(.el-table .el-table__row),
  .tx-detail-panel :deep(.el-table td.el-table__cell),
  .tx-detail-panel :deep(.el-table th.el-table__cell) {
    background-color: transparent !important;
  }
  .tx-detail-panel :deep(.el-table) {
    --el-table-bg-color: transparent;
    --el-table-tr-bg-color: transparent;
    --el-table-header-bg-color: rgba(8, 24, 43, 0.74);
    --el-table-row-hover-bg-color: rgba(68, 213, 255, 0.08);
    --el-table-border-color: rgba(139, 214, 255, 0.16);
    --el-table-text-color: rgba(232, 247, 255, 0.84);
    --el-table-header-text-color: #f4fbff;
    color: rgba(232, 247, 255, 0.84);
  }
  .tx-detail-panel :deep(.el-pagination),
  .tx-detail-panel :deep(.el-pagination button),
  .tx-detail-panel :deep(.el-pager li),
  .tx-detail-panel :deep(.el-select .el-select__wrapper),
  .tx-detail-panel :deep(.el-pagination .el-select .el-select__wrapper) {
    background-color: transparent !important;
    color: rgba(232, 247, 255, 0.84);
  }
  .container {
    padding: 20px;
  }
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
  </style>