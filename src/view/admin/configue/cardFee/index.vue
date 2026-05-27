<template>
  <div class="p-2">
    <el-tabs v-model="activeTab" class="mb-4" @tab-change="handleTabChange">
      <el-tab-pane label="Global fee config" name="globalFeeConfig">
        <div>
          <el-form :inline="true" class="mb-4" label-position="top">
            <el-form-item label="Fee Type">
              <el-select v-model="search.feeType" placeholder="select fee type" clearable>
                <el-option v-for="item of feeTypes" :label="item.label" :value="item.value"></el-option>
            </el-select>

            </el-form-item>
            <el-form-item label="Available">
              <el-select v-model="search.available" clearable placeholder="All">
                <el-option label="Enable" :value="true"></el-option>
                <el-option label="Disable" :value="false"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="Card Bin">
              <el-input v-model="search.cardBin" placeholder="Please enter card bin"></el-input>
            </el-form-item>
            <el-form-item label="Card Type">
              <el-input v-model="search.cardType" placeholder="Please enter card type"></el-input>
            </el-form-item>
            <el-form-item label=" ">
              <el-button type="primary" @click="handleSearch">Search</el-button>
              <el-button @click="resetSearch">Reset</el-button>
            </el-form-item>
          </el-form>
          <div style="float:right;margin-bottom: 10px;"><el-button type="primary" icon="plus" @click="dialogs.addFeeGlobalCfgDialog = true;resetFeeCfg()">Add</el-button></div>
          <el-table :data="tableData" style="width: 100%">
            <el-table-column prop="feeType" label="Fee Type">
              <template #default="{row}">{{ formatFeeType(row.feeType) }}</template>
            </el-table-column>
            <el-table-column prop="cardBin" label="Card Bin"></el-table-column>
            <el-table-column prop="cardModel" label="Card Model"></el-table-column>
            <el-table-column prop="fee" label="Fee">
              <template #default="{row}">
                  <span v-if="row.calType == 2">{{ row.fee }}</span>
                  <span v-else>{{ row.fee }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="available" label="Available">
              <template #default="{ row }">
                <el-tag v-if="row.available === true" type="success">Active</el-tag>
                <el-tag v-else type="danger">Disable</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Operate">
              <template #default="{row}">
                <el-button size="small" type="text" @click="handleEdit(row)">Edit</el-button>
                <el-button size="small" type="text" @click="handleChangeFeeConfigStatus(row, row.available === true ? false : true)">{{ row.available === true ? 'Disable' : 'Active' }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="mt-4">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="search.total"
              :page-sizes="[10, 50, 100]"
              :page-size="search.pageSize"
              :current-page="search.page"
              @size-change="handlePageSizeChange"
              @current-change="handlePageChange"
            ></el-pagination>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="User fee config" name="userFeeConfig">
        <div>
          <el-form :inline="true" label-position="top" :model="search" class="demo-form-inline">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Client No">
                  <el-input v-model="search.clientNo" placeholder="Please enter client No"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label=" ">
                  <el-button type="primary" @click="getTableData">Search</el-button>
                  <el-button @click="resetSearch">Reset</el-button>               
               </el-form-item>
              </el-col>
            </el-row>
          </el-form>
          <div style="float:right;margin-bottom: 10px;"><el-button type="primary" icon="plus" @click="dialogs.addFeeUserCfgDialog = true;resetFeeCfg()">Add</el-button></div>
          <!-- Table -->
          <el-table :data="tableData" style="width: 100%" border>
            <el-table-column prop="clientNo" label="Client No" width="180"></el-table-column>
            <el-table-column prop="cardBin" label="Card Bin"></el-table-column>
            <el-table-column prop="feeType" label="Fee Type" ><template #default="{row}">{{ formatFeeType(row.feeType) }}</template></el-table-column>
            <el-table-column prop="fee" label="Fee">
              <template #default="{row}">
                  <span v-if="row.calType == 2">{{ row.fee }}</span>
                  <span v-else>{{ row.fee }}%</span>
              </template>
            </el-table-column>
            <el-table-column label="Operate" width="100">
              <template #default="{ row }">
                <el-button size="small" type="text" @click="handleEditMulti(row,2)">View</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="mt-4">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="search.total"
              :page-sizes="[10, 50, 100]"
              :page-size="search.pageSize"
              :current-page="search.page"
              @size-change="handlePageSizeChange"
              @current-change="handlePageChange"
            ></el-pagination>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="Month fee config" name="monthFeeConfig">
        <div>
          <el-form :inline="true" label-position="top" :model="search" class="demo-form-inline">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Card Bin ID">
                  <el-input v-model="search.cardBinId" placeholder="Please enter card bin id"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Client No">
                  <el-input v-model="search.clientNo" placeholder="Please enter client No"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="24">
                <el-form-item label=" ">
                  <el-button type="primary" @click="handleSearch">Search</el-button>
                <el-button @click="resetSearch">Reset</el-button>
              </el-form-item>
              </el-col>
            </el-row>
          </el-form>
          <div>
            <hr>
            <div><el-button class="float-end" type="primary" icon="plus" @click="dialogs.addFeeMonthCfgDialog = true;resetFeeCfg()">Add</el-button></div>
          </div>
          <el-table :data="tableData" style="width: 100%" border>
            <el-table-column prop="cardBinId" label="Card Bin ID" width="180"></el-table-column>
            <el-table-column prop="cardBin" label="Card Bin"></el-table-column>
            <el-table-column prop="userId" label="User ID" ></el-table-column>
            <el-table-column label="Operate" width="100">
              <template #default="{ row }">
                <el-button size="small" type="text" @click="handleEdit(row)">View</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="mt-4">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="search.total"
              :page-sizes="[10, 50, 100]"
              :page-size="search.pageSize"
              :current-page="search.page"
              @size-change="handlePageSizeChange"
              @current-change="handlePageChange"
            ></el-pagination>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <el-dialog v-model="dialogs.addFeeGlobalCfgDialog" title="Set global Cfg" width="80%" style="max-width: 600px;">
    <el-form :model="feeCfg" label-position="top" label-width="160px">
      
      <el-row :gutter="20"> 
            <el-col :span="12">
          <el-form-item label="* Fee Type">
            <el-select v-model="feeCfg.feeType" placeholder="select fee type" clearable>
              <el-option v-for="item of feeTypes" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Card Bin">
            <el-select v-model="feeCfg.cardBin">
              <el-option label="All" value="All"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Fee Cal Method">
            <el-radio-group v-model="feeCfg.calType">
              <el-radio :value="1" >rate</el-radio>
              <el-radio :value="2">fixed</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Fixed Fee" v-if="feeCfg.calType == 2">
            <el-input-number v-model="feeCfg.fee" placeholder="Please enter Fee Rate" :min="0" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item label="* Fee Rate" v-else>
            <el-input-number v-model="feeCfg.fee" placeholder="Please enter Fee Rate" :min="0" :controls="false">
              <template #suffix>%</template>
            </el-input-number>
          </el-form-item>
        </el-col>
        <!-- <el-col :span="12">
          <el-form-item label="* Min Fee">
            <el-input-number v-model="feeCfg.minFee" placeholder="Please enter min fee" :min="0" :controls="false"></el-input-number>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Max Fee">
            <el-input-number v-model="feeCfg.maxFee" placeholder="Please enter max fee" :min="feeCfg.minFee" :controls="false"></el-input-number>
          </el-form-item>
        </el-col> -->
      </el-row>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogs.addFeeGlobalCfgDialog = false">Cancel</el-button>
        <el-button type="primary" @click="handleAddGlobalFeeCfg">Confirm</el-button>
      </span>
    </template>
  </el-dialog>

  <el-dialog v-model="dialogs.addFeeUserCfgDialog" title="Set user config" width="80%" style="max-width: 600px;">
    <el-form :model="feeCfg" label-position="top" label-width="160px">
      
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Client No">
            <el-input v-model="feeCfg.clientNo"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Fee Type">
            <el-select v-model="feeCfg.feeType" placeholder="select fee type" clearable>
              <el-option v-for="item of feeTypes" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Card Bin">
            <el-select v-model="feeCfg.cardBin">
              <el-option label="All" value="All"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Fee Cal Method">
            <el-radio-group v-model="feeCfg.calType">
              <el-radio :value="1" >rate</el-radio>
              <el-radio :value="2">fixed</el-radio>
            </el-radio-group>
          </el-form-item>
        <!-- </el-col>
        <el-col :span="12"> -->
          <el-form-item label="* Fixed Fee" v-if="feeCfg.calType == 2">
            <el-input-number v-model="feeCfg.fee" placeholder="Please enter Fee Rate" :min="0" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item label="* Fee Rate" v-else>
            <el-input-number v-model="feeCfg.fee" placeholder="Please enter Fee Rate" :min="0" :controls="false">
              <template #suffix>%</template>
            </el-input-number>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogs.addFeeUserCfgDialog = false">Cancel</el-button>
        <el-button type="primary" @click="handleAddUserFeeCfg">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="dialogs.addFeeUserCfgMultiDialog" title="Add Fee User Cfg" width="80%">
    <el-tabs v-model="multiActiveTab" type="card" tab-position="left" class="mb-2 text-secondary" @tab-change="handleMultiCfgTabChange">
      <el-tab-pane v-for="item of feeTypes" :label="item.label" :name="item.value">
        <div class="row bg-body-secondary p-2 mx-2">
            <div class="col-6 text-center">
              <p><span class="fs-6 fw-bold">Current Config:</span> <el-tag type="success">{{ currTabCfg.ID && currTabCfg.available?currClient.clientNo:'Global'}}</el-tag></p>
            </div>
            <div class="col-4 offset-2"> 
                <div>
                  <button v-if="currTabCfg.ID > 0 && currTabCfg.available" class="btn btn-danger me-1" @click="handleSetFeeCfgGlobal">Set As Global</button>
                  <button class="btn btn-primary" @click="handleEdit(currTabCfg.ID?currTabCfg:{clientNo:currClient.clientNo},2)">Edit</button>
              </div>
              <div>Available: {{ currTabCfg.available === true ? 'Y' : 'N' }}</div>
            </div>
        </div>
        <div class="row p-4 mt-4 mx-2" style="border: 1px solid #ddd; border-top: none;">
            
            <div class="col-4 offset-1" v-if="currTabCfg.calType == 1"> 
              <p class="fw-bold">Fee Rate:&nbsp;&nbsp; {{currTabCfg.fee}}%</p>
            </div>
            <div class="col-4 offset-1" v-else> 
              <p class="fw-bold">Fixed Fee:&nbsp;&nbsp; {{currTabCfg.fee}}</p>
            </div>
            <!-- <div class="col-4 offset-1"> 
              <p class="fw-bold">Min Fee:&nbsp;&nbsp; {{currTabCfg.minFee}}</p>
            </div>
            <div class="col-4 offset-1"> 
              <p class="fw-bold">Max Fee:&nbsp;&nbsp; {{currTabCfg.maxFee}}</p>
            </div> -->
            <!-- <div class="col-4 offset-1"> 
              <p class="fw-bold">Decline Fee:&nbsp;&nbsp; {{currTabCfg.declineFee}}</p>
            </div> -->
        </div>
        <div class="row p-2 mx-2 mt-4">
          <el-table :data="feeCfgOpLogs" >
            <el-table-column label="Operator" prop="name"></el-table-column>
            <el-table-column label="Operate Time">
              <template #default="{row}">
                {{formatDate(row.CreatedAt)}}
              </template>
            </el-table-column>
            <el-table-column label="detail">
              
            </el-table-column>
          </el-table>
          <div class="mt-1">
            <el-pagination
              style="padding-top: 0px;float:right"
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="logSearch.total"
              :page-size="logSearch.pageSize"
              :current-page="logSearch.page"
              @size-change="handleLogPageSizeChange"
              @current-change="handleLogPageChange"
            ></el-pagination>
          </div>
        </div>
      </el-tab-pane>
    
    </el-tabs>
  </el-dialog>
  <el-dialog title="Add month fee config"
    v-model="dialogs.addFeeMonthCfgDialog"
    width="50%"
  >
  <el-form :model="feeCfg" label-position="top" label-width="160px" class="p-2">
      <el-form-item label="Client No">
        <el-input v-model="feeCfg.clientNo"></el-input>
      </el-form-item>
      <el-form-item label="CardBin">
        <el-select v-model="feeCfg.cardBin">
            <el-option label="All" value="All"></el-option>
            <el-option v-for="item in cardBinList" :key="item.id" :label="item.cardBin" :value="item.cardBin"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="Fixed Fee">
        <el-input-number v-model="feeCfg.fee" :controls="false" :min="0">
          <template #suffix>
            USD
          </template>
        </el-input-number>
      </el-form-item>
  </el-form>
  <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogs.addFeeMonthCfgDialog = false">Cancel</el-button>
        <el-button type="primary" @click="handleAddMonthFeeCfg">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive,onMounted } from 'vue';
import { addFeeGlobalCfg,addFeeUserCfg,addFeeMonthCfg,setUserCfgGlobal,listFeeUserCfg,listFeeGlobalCfg,listLogs } from '@/api/finance';
import {listCardBin} from '@/api/card'
import { formatDate } from '@/utils/format';
import { ElMessage } from 'element-plus'
const feeTypes = [
  { label: "Create Card", value: "CREATE_CARD" },
  { label: "Recharge Card", value: "RECHARGE_CARD" },
  { label: "Terminate Card", value: "TERMINATE_CARD" },
  { label: "Authorization Transaction", value: "AUTHORIZATION_TRANSACTION" },
  { label: "Refund Transaction", value: "REFUND_TRANSACTION" },
  { label: "Settlement Transaction", value: "SETTLEMENT_TRANSACTION" },
  { label: "Chargeback", value: "CHARGEBACK" },
  { label: "Cross Board", value: "CROSS_BOARD" },
  { label: "Auth Reversal Transaction", value: "AUTH_REVERSAL_TRANSACTION" },
  { label: "Refund Reversal Transaction", value: "REFUND_REVERSAL_TRANSACTION" },
  { label: "Authorization Query", value: "AUTHORIZATION_QUERY" },
  { label: "Withdraw Card", value: "WITHDRAW_CARD" },
  { label: "Atm Transaction", value: "ATM_TRANSACTION" }
]
const activeTab = ref('userFeeConfig');
const dialogVisible = ref(false);
const search =reactive({
  cardBin: '',
  cfgType: 1,
  feeType:'',
  cardBin:'All',
  cardType:'',
  cardBinId:'',
  available: null,
  page:1,
  pageSize: 10,
  total: 0,
})
const feeCfg = ref({
  calType:1,
  feeType:'',
  cardBin:'',
  cardType:'',
  minFee:0,
  maxFee:0,
  feeRate:0,
  fixedFee:0,

 })
 const resetSearch =() =>{
  search.cardBin = ''
  search.feeType = ''
  search.cfgType = 1,
  search.cardBin = 'All'
  search.cardType = ''
  search.cardBinId = ''
  search.available = null
}
const resetFeeCfg = ()=>{
  feeCfg.value = {
  calType:1,
  feeType:'',
  cardBin:'',
  cardType:'',
  minFee:0,
  maxFee:0,
  feeRate:0,
  fixedFee:0,

 }
}
const handleSetFeeCfgGlobal = () =>{
  setUserCfgGlobal({
    id: currTabCfg.value.ID 
  }).then(res =>{
    ElMessage.success("Success")
    currTabCfg.value.available = false
  })
}
 const handleAddGlobalFeeCfg = () =>{
  addFeeGlobalCfg(feeCfg.value).then(res =>{
    if(res.code === 0){
      ElMessage.success('Success')
      dialogs.addFeeGlobalCfgDialog = false
      getTableData()
    }
  })
 }

 const handleAddUserFeeCfg = () =>{
  addFeeUserCfg(feeCfg.value).then(res =>{
    if(res.code === 0){
      ElMessage.success('Success')
      dialogs.addFeeUserCfgDialog = false
      getTableData()
      currTabCfg.value.available = true
    }
  })
 }
 const cardBinList = ref([])
onMounted(()=>{
  getTableData()
  listCardBin({
    page: 1,
    pageSize: 1000
  }).then(res =>{
    cardBinList.value = res.data.list
  })
})
const handleAddMonthFeeCfg = ()=>{
  var cfg ={
    ...feeCfg.value,
    feeType: 'CARD_MONTH_FEE'
  }
  addFeeMonthCfg(cfg).then(res =>{
    if(res.code === 0){
      ElMessage.success('Success')
      dialogs.addFeeMonthCfgDialog = false
      getTableData()
    }
  })
}
const userFeeCfgs = ref([])
const multiActiveTab = ref('CREATE_CARD')
const feeCfgOpLogs = ref([])
const currClient =ref()
const logSearch = reactive({
  page:1,
  pageSize:10,
  total:0,
})
const getFeeCfgOpLogs = (row) =>{
  if (row == null){
    logSearch.total = 0
    logSearch.page = 1
    return
  }
  logSearch.objId = row.ID
  logSearch.opType = 4
  listLogs(logSearch).then(res =>{
    if(res.code === 0){
      feeCfgOpLogs.value = res.data.list
      logSearch.total = res.data.total
    }
  })
}
const currTabCfg = ref({})
const handleMultiCfgTabChange = (tab) => {
  multiActiveTab.value = tab
  const row = filterUserFeeCfg(tab)
  currTabCfg.value = {}
  feeCfgOpLogs.value = []
  if (row.length > 0 ){
    currTabCfg.value = row[0]
    getFeeCfgOpLogs(row[0])
  }else {
    getFeeCfgOpLogs()

  }
}
const filterUserFeeCfg = (type)=>{
  return userFeeCfgs.value.filter(item => item.feeType === type)
}
const handleEditMulti = (row) =>{ 
   currClient.value = row
  listFeeUserCfg({
    clientNo: row.clientNo,
    page:1,
    pageSize:999,
  }).then(res =>{ 
    if (res.code === 0){
      multiActiveTab.value = row.feeType
      userFeeCfgs.value = res.data.list
      handleMultiCfgTabChange(row.feeType)
      dialogs.addFeeUserCfgMultiDialog = true
    }
  })
  
}
const dialogs  = reactive({
  addFeeGlobalCfgDialog:false,
  addFeeUserCfgDialog:false,
  addFeeUserCfgMultiDialog:false,
  addFeeMonthCfgDialog:false,
})
const tableData = ref([])
const getTableData = () => {
  if(activeTab.value == 'globalFeeConfig'){
    listFeeGlobalCfg(search).then(res => {
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
        currTabCfg.value = filterUserFeeCfg(multiActiveTab.value)
      }
    })
  }else if (activeTab.value == 'userFeeConfig'){
    listFeeUserCfg(search).then(res => {
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }else {
    var newSearch = {
      ...search,
      cardBin:'',
      feeType: 'CARD_MONTH_FEE'
    }

    listFeeUserCfg(newSearch).then(res => {
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }
}
const handleTabChange = (tab) => {
  activeTab.value = tab
  getTableData()
}
const formatFeeType =(str) =>{
  return str
    .toLowerCase()
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}
const handleSearch = () => {
  getTableData()
}

const handlePageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.page = val
  getTableData()
}
const handleLogPageSizeChange =(val) =>{
  logSearch.pageSize = val
  getFeeCfgOpLogs(currClient.value)
}
const handleLogPageChange =(val) =>{
  logSearch.page = val
  getTableData()
}


const handleEdit = (row,type) => {
  feeCfg.value = row
  if (type == 2){
    dialogs.addFeeUserCfgDialog = true
  }else {
    dialogs.addFeeGlobalCfgDialog = true
  }
  
};

const handleFeeCfgAdd = ()=>{
  dialogs.addFeeGlobalCfgDialog = true
  feeCfg.value.type = 'global'
}
const handleChangeFeeConfigStatus = (row, available) => {
  row.available = available
  addFeeGlobalCfg(row).then((res) => {
    getTableData()
  })
}
const onSearch = () => {
  // Implement search logic here
  ElMessage.success('Search button clicked');
};

const handleEditBlockList = (row) => {
  // Implement edit logic here
  console.log('Editing row:', row);
};
</script>

<style scoped>
/* :deep(.el-select__wrapper) {
    min-width: 150px;
} */
</style>