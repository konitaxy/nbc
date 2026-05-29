<template>
  <div class="p-2">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="Watch Addresses" name="addresses">
        <el-form :inline="true" label-position="top" :model="addressSearch">
          <el-form-item label="Chain Type">
            <el-select v-model="addressSearch.chainType" clearable placeholder="All" style="width: 140px">
              <el-option label="TRON" value="TRON" />
            </el-select>
          </el-form-item>
          <el-form-item label="Address">
            <el-input v-model="addressSearch.address" placeholder="Watch address" clearable style="width: 280px" />
          </el-form-item>
          <el-form-item label="Enabled">
            <el-select v-model="addressSearch.enabled" clearable placeholder="All" style="width: 120px">
              <el-option label="Enabled" :value="true" />
              <el-option label="Disabled" :value="false" />
            </el-select>
          </el-form-item>
          <el-form-item label=" ">
            <el-button type="primary" @click="searchAddresses">Search</el-button>
            <el-button @click="resetAddressSearch">Reset</el-button>
          </el-form-item>
        </el-form>

        <div class="toolbar">
          <el-button type="primary" icon="Plus" @click="openAddAddressDialog">Add Address</el-button>
        </div>

        <el-table :data="addressTableData" show-overflow-tooltip>
          <el-table-column prop="ID" label="ID" width="70" />
          <el-table-column prop="chainType" label="Chain" width="90" />
          <el-table-column prop="address" label="Address" min-width="280" />
          <el-table-column prop="contractAddress" label="Contract" min-width="200">
            <template #default="{ row }">
              {{ row.contractAddress || '(default USDT)' }}
            </template>
          </el-table-column>
          <el-table-column prop="watchTrx" label="Watch TRX" width="100">
            <template #default="{ row }">
              <el-tag :type="row.watchTrx ? 'success' : 'info'" size="small">
                {{ row.watchTrx ? 'Yes' : 'No' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="enabled" label="Enabled" width="90">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'danger'" size="small">
                {{ row.enabled ? 'Yes' : 'No' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="Remark" min-width="140" />
          <el-table-column label="Created At" width="170">
            <template #default="{ row }">
              {{ formatDate(row.CreatedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="Action" width="100" fixed="right">
            <template #default="{ row }">
              <el-popconfirm title="Delete this watch address?" @confirm="handleDeleteAddress(row)">
                <template #reference>
                  <el-button type="danger" size="small" link>Delete</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <div class="mt-4">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :total="addressSearch.total"
            :page-sizes="[10, 20, 50]"
            :page-size="addressSearch.pageSize"
            :current-page="addressSearch.page"
            @size-change="handleAddressPageSizeChange"
            @current-change="handleAddressPageChange"
          />
        </div>
      </el-tab-pane>

      <el-tab-pane label="Inbound Transactions" name="transactions">
        <el-form :inline="true" label-position="top" :model="txSearch">
          <el-form-item label="Chain Type">
            <el-select v-model="txSearch.chainType" clearable placeholder="All" style="width: 140px">
              <el-option label="TRON" value="TRON" />
            </el-select>
          </el-form-item>
          <el-form-item label="To Address">
            <el-input v-model="txSearch.toAddress" placeholder="Receive address" clearable style="width: 280px" />
          </el-form-item>
          <el-form-item label="Tx Hash">
            <el-input v-model="txSearch.transactionId" placeholder="Transaction hash" clearable style="width: 240px" />
          </el-form-item>
          <el-form-item label="Symbol">
            <el-select v-model="txSearch.symbol" clearable placeholder="All" style="width: 120px">
              <el-option label="USDT" value="USDT" />
              <el-option label="TRX" value="TRX" />
            </el-select>
          </el-form-item>
          <el-form-item label="Transaction Time">
            <el-date-picker
              v-model="txSearch.timeRange"
              type="datetimerange"
              range-separator="To"
              start-placeholder="Start"
              end-placeholder="End"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 360px"
            />
          </el-form-item>
          <el-form-item label=" ">
            <el-button type="primary" @click="searchTransactions">Search</el-button>
            <el-button @click="resetTxSearch">Reset</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="txTableData" show-overflow-tooltip>
          <el-table-column prop="ID" label="ID" width="70" />
          <el-table-column prop="chainType" label="Chain" width="80" />
          <el-table-column prop="symbol" label="Symbol" width="80" />
          <el-table-column prop="amount" label="Amount" width="120" />
          <el-table-column prop="kind" label="Kind" width="80" />
          <el-table-column prop="toAddress" label="To Address" min-width="240" />
          <el-table-column prop="fromAddress" label="From Address" min-width="240" />
          <el-table-column prop="transactionId" label="Tx Hash" min-width="240" />
          <el-table-column prop="watchAddressId" label="Watch ID" width="90" />
          <el-table-column prop="contractAddress" label="Contract" min-width="200" />
          <el-table-column label="Transaction Time" width="170">
            <template #default="{ row }">
              {{ formatDate(row.transactionTime) }}
            </template>
          </el-table-column>
        </el-table>

        <div class="mt-4">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :total="txSearch.total"
            :page-sizes="[10, 20, 50]"
            :page-size="txSearch.pageSize"
            :current-page="txSearch.page"
            @size-change="handleTxPageSizeChange"
            @current-change="handleTxPageChange"
          />
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="addAddressDialogVisible" title="Add Watch Address" width="560px" align-center>
      <el-form :model="addressForm" label-width="140px">
        <el-form-item label="Chain Type">
          <el-select v-model="addressForm.chainType" style="width: 100%">
            <el-option label="TRON" value="TRON" />
          </el-select>
        </el-form-item>
        <el-form-item label="Address *">
          <el-input v-model="addressForm.address" placeholder="TRON Base58 address" />
        </el-form-item>
        <el-form-item label="Contract Address">
          <el-input
            v-model="addressForm.contractAddress"
            placeholder="Leave empty for default USDT contract"
          />
        </el-form-item>
        <el-form-item label="Watch TRX">
          <el-switch v-model="addressForm.watchTrx" />
        </el-form-item>
        <el-form-item label="Enabled">
          <el-switch v-model="addressForm.enabled" />
          <span class="form-hint">Must be enabled to participate in polling</span>
        </el-form-item>
        <el-form-item label="Remark">
          <el-input v-model="addressForm.remark" type="textarea" :rows="2" placeholder="Optional remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addAddressDialogVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="addingAddress" @click="handleAddAddress">Confirm</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'ChainWatch'
}
</script>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  addChainWatchAddress,
  deleteChainWatchAddress,
  listChainWatchAddress,
  listChainInboundTransaction
} from '@/api/finance'

const activeTab = ref('addresses')
const addressTableData = ref([])
const txTableData = ref([])
const addAddressDialogVisible = ref(false)
const addingAddress = ref(false)

const defaultAddressForm = () => ({
  chainType: 'TRON',
  address: '',
  contractAddress: '',
  watchTrx: false,
  enabled: true,
  remark: ''
})

const addressForm = ref(defaultAddressForm())

const addressSearch = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  chainType: '',
  address: '',
  enabled: undefined
})

const txSearch = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
  chainType: '',
  toAddress: '',
  transactionId: '',
  symbol: '',
  timeRange: null
})

const buildAddressListParams = () => {
  const params = {
    page: addressSearch.page,
    pageSize: addressSearch.pageSize
  }
  if (addressSearch.chainType) params.chainType = addressSearch.chainType
  if (addressSearch.address) params.address = addressSearch.address
  if (addressSearch.enabled === true || addressSearch.enabled === false) {
    params.enabled = addressSearch.enabled
  }
  return params
}

const buildTxListParams = () => {
  const params = {
    page: txSearch.page,
    pageSize: txSearch.pageSize
  }
  if (txSearch.chainType) params.chainType = txSearch.chainType
  if (txSearch.toAddress) params.toAddress = txSearch.toAddress
  if (txSearch.transactionId) params.transactionId = txSearch.transactionId
  if (txSearch.symbol) params.symbol = txSearch.symbol
  if (txSearch.timeRange?.length === 2) params.timeRange = txSearch.timeRange
  return params
}

const getAddressTableData = () => {
  listChainWatchAddress(buildAddressListParams()).then((res) => {
    if (res.code === 0) {
      addressTableData.value = res.data.list || []
      addressSearch.total = res.data.total || 0
    }
  })
}

const getTxTableData = () => {
  listChainInboundTransaction(buildTxListParams()).then((res) => {
    if (res.code === 0) {
      txTableData.value = res.data.list || []
      txSearch.total = res.data.total || 0
    }
  })
}

const handleTabChange = (tab) => {
  if (tab === 'addresses') {
    getAddressTableData()
  } else if (tab === 'transactions') {
    getTxTableData()
  }
}

const searchAddresses = () => {
  addressSearch.page = 1
  getAddressTableData()
}

const resetAddressSearch = () => {
  addressSearch.chainType = ''
  addressSearch.address = ''
  addressSearch.enabled = undefined
  addressSearch.page = 1
  getAddressTableData()
}

const searchTransactions = () => {
  txSearch.page = 1
  getTxTableData()
}

const resetTxSearch = () => {
  txSearch.chainType = ''
  txSearch.toAddress = ''
  txSearch.transactionId = ''
  txSearch.symbol = ''
  txSearch.timeRange = null
  txSearch.page = 1
  getTxTableData()
}

const openAddAddressDialog = () => {
  addressForm.value = defaultAddressForm()
  addAddressDialogVisible.value = true
}

const handleAddAddress = () => {
  if (!addressForm.value.address?.trim()) {
    ElMessage.warning('Address is required')
    return
  }
  addingAddress.value = true
  addChainWatchAddress({
    chainType: addressForm.value.chainType,
    address: addressForm.value.address.trim(),
    contractAddress: addressForm.value.contractAddress?.trim() || '',
    watchTrx: addressForm.value.watchTrx,
    enabled: addressForm.value.enabled,
    remark: addressForm.value.remark?.trim() || ''
  }).then((res) => {
    addingAddress.value = false
    if (res.code === 0) {
      ElMessage.success('Address added')
      addAddressDialogVisible.value = false
      getAddressTableData()
    }
  }).catch(() => {
    addingAddress.value = false
  })
}

const handleDeleteAddress = (row) => {
  deleteChainWatchAddress({ id: row.ID }).then((res) => {
    if (res.code === 0) {
      ElMessage.success('Deleted')
      getAddressTableData()
    }
  })
}

const handleAddressPageSizeChange = (val) => {
  addressSearch.pageSize = val
  getAddressTableData()
}

const handleAddressPageChange = (val) => {
  addressSearch.page = val
  getAddressTableData()
}

const handleTxPageSizeChange = (val) => {
  txSearch.pageSize = val
  getTxTableData()
}

const handleTxPageChange = (val) => {
  txSearch.page = val
  getTxTableData()
}

onMounted(() => {
  getAddressTableData()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}

.form-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}
</style>
