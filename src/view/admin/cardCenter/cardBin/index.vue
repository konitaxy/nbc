<template>
  <div class="p-2">
    <el-tabs v-model="activeTab" class="mb-4" @tab-change="getTableData">
      <el-tab-pane label="Card Bin List" name="cardBinList">
        <div>
          <el-form :inline="true" class="mb-4">
            <el-form-item label="Card Bin">
              <el-input v-model="search.cardBin" placeholder="Please enter card bin"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">Search</el-button>
              <el-button @click="resetSearch">Reset</el-button>
            </el-form-item>
          </el-form>
          <div style="float:right;margin-bottom: 10px;"><el-button type="primary" icon="plus" @click="dialogVisible = true">Add</el-button></div>
          <el-table :data="tableData" style="width: 100%">
            <el-table-column prop="cardBinId" label="Card Bin ID"></el-table-column>
            <el-table-column prop="cardBin" label="Card Bin"></el-table-column>
            <el-table-column prop="cardModel" label="Card Model"></el-table-column>
            <el-table-column prop="cardBrand" label="Brand"></el-table-column>
            <el-table-column prop="currency" label="Currency"></el-table-column>
            <el-table-column prop="qtyIssuanceLimitCardholder" label="QTY Limit PER Card Holder"></el-table-column>
            <el-table-column prop="channel" label="Channel"></el-table-column>
            <el-table-column label="Operate">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.$index, scope.row)">Edit</el-button>
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
      <el-tab-pane label="Card Bin Block List" name="cardBinBlockList">
        <div>
          <el-form :inline="true" label-position="top" :model="searchForm" class="demo-form-inline">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Card Bin ID">
                  <el-input v-model="searchForm.cardBinId" placeholder="Please enter card bin id"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="User ID">
                  <el-input v-model="searchForm.userId" placeholder="Please enter user id"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="24">
                <el-form-item label=" ">
                  <el-button type="primary" @click="onSearch">Search</el-button>
                <el-button @click="onReset">Reset</el-button>                </el-form-item>
                
              </el-col>
            </el-row>
          </el-form>

          <!-- Table -->
          <el-table :data="tableData" style="width: 100%" border>
            <el-table-column prop="cardBinId" label="Card Bin ID" width="180"></el-table-column>
            <el-table-column prop="cardBin" label="Card Bin"></el-table-column>
            <el-table-column prop="userId" label="User ID" ></el-table-column>
            <el-table-column label="Operate" width="100">
              <template #default="{ row }">
                <el-button size="small" type="text" @click="handleEdit(row)">Edit</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="mt-4">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="search.total"
              :page-size="search.pageSize"
              :current-page="search.page"
              @size-change="handlePageSizeChange"
              @current-change="handlePageChange"
            ></el-pagination>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <el-dialog v-model="dialogVisible" title="Add Card Bin" width="80%" align-center style="max-width: 600px;" @close="resetCardBinForm">
    <el-form :model="cardBinForm" label-position="top" label-width="160px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Card Bin ID">
            <el-input v-model="cardBinForm.cardBinId" placeholder="Please enter card bin id"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Card Bin">
            <el-input v-model="cardBinForm.cardBin" placeholder="Please enter card bin"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="Card Model">
            <el-select v-model="cardBinForm.cardModel" placeholder="Select card model">
              <el-option label="Share" value="SHARE"></el-option>
              <el-option label="Card" value="CARD"></el-option>
              <!-- Add more options as needed -->
            </el-select>          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Card Brand">
            <el-select v-model="cardBinForm.cardBrand" placeholder="Select card brand">
              <el-option label="Visa" value="Visa"></el-option>
              <el-option label="MasterCard" value="MasterCard"></el-option>
              <!-- Add more options as needed -->
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Currency">
            <el-select v-model="cardBinForm.currency" placeholder="Select card currency">
              <el-option label="USD" value="USD"></el-option>
              <el-option label="USDT" value="USDT"></el-option>
              <el-option label="RMB" value="RMB"></el-option>
              <!-- Add more options as needed -->
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Region">
            <el-select v-model="cardBinForm.region" placeholder="Select card regon">
              <el-option label="USA" value="USA"></el-option>
              <el-option label="CHN" value="CHN"></el-option>
              <el-option label="VN" value="VN"></el-option>
              <!-- Add more options as needed -->
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="Channel">
            <el-select v-model="cardBinForm.channel" placeholder="Select card Channel">
              <el-option label="SLASH" value="SLASH"></el-option>
              <!-- Add more options as needed -->
            </el-select>
            </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Description">
            <el-input v-model="cardBinForm.description" placeholder="Please enter description"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* QTY issuance limit for client">
            <el-input-number v-model="cardBinForm.qtyIssuanceLimitCardholder" :min="1" :controls="false"></el-input-number>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Auth amount limit">
            <el-input-number v-model="cardBinForm.authAmountLimit" :min="1" :controls="false"></el-input-number>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Create Recharge Limit">
            <el-input-number v-model="cardBinForm.createRechargeLimit" placeholder="Enter recharge limit" :controls="false"></el-input-number>
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="* Min Recharge Balance">
            <el-input-number v-model="cardBinForm.minBalance" :min="1" :controls="false"></el-input-number>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="20">
          <el-form-item label="Support Platform">
            <el-input v-model="cardBinForm.supportPlatform" type="textarea" :rows="5" placeholder="Enter support platform details"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Bin Status">
            <el-radio-group v-model="cardBinForm.binStatus">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Issuer Available">
            <el-radio-group v-model="cardBinForm.issuerAvailable">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Cancel Card">
            <el-radio-group v-model="cardBinForm.cancelCard">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Top-up">
            <el-radio-group v-model="cardBinForm.topUp">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Withdrawal">
            <el-radio-group v-model="cardBinForm.withdrawal">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Customer Available">
            <el-radio-group v-model="cardBinForm.customerAvailable">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Support freezing">
            <el-radio-group v-model="cardBinForm.supportFreezing">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="* Cardholder required">
            <el-radio-group v-model="cardBinForm.cardholderRequired">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="* Channel Auto Cancel">
            <el-radio-group v-model="cardBinForm.channelAutoCancel">
              <el-radio :label="true">Y</el-radio>
              <el-radio :label="false">N</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <!-- Add more form items as per the image -->
      <!-- ... -->
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleAddCardBinConfirm">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, reactive,onMounted } from 'vue';
import { addCardBin,blockCardBin,listCardBin } from '@/api/card';
import { ElMessage } from 'element-plus';

const activeTab = ref('cardBinList');
const dialogVisible = ref(false);
const search =reactive({
  cardBin: '',
  cardBinId:'',
  page:1,
  pageSize: 10,
  total: 0,
})
const searchForm =reactive({
  cardBin: '',
  cardBinId:'',
  page:1,
  pageSize: 10,
  total: 0,
})
const initialCardBinForm = {
  cardBinId: '',
  cardBin: '',
  channelCardBinId: '',
  cardBrand: '',
  cardModel: '',
  currency: '',
  region: '',
  channel: '',
  // Add other fields as per the image
  qtyIssuanceLimitCardholder: 0,
  authAmountLimit:0,
  description: '',
  createRechargeLimit: '',
  supportPlatform: '',
  minBalance: 0,
  binStatus: true,
  issuerAvailable: true,
  cancelCard: true,
  topUp: true,
  withdrawal: true,
  customerAvailable: true,
  supportFreezing: true,
  cardholderRequired: false,
  channelAutoCancel: true,
};

const cardBinForm = ref({ ...initialCardBinForm });
onMounted(()=>{
  getTableData()
})
const resetCardBinForm = () => {
  cardBinForm.value = { ...initialCardBinForm };
}

const handleAddCardBinConfirm = ()=>{
  addCardBin(cardBinForm.value).then(res=>{
    if (res.code === 0){
      dialogVisible.value = false
      ElMessage.success('Add Success')
      getTableData()
    }
    
  })
}
const tableData = ref([])
const getTableData = () => {
  search.blocked = activeTab.value == 'cardBinBlockList'
  listCardBin(search).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.total = res.data.total
    }
  })
};
const handleSearch = () => {
  getTableData()
}

const resetSearch = () => {
  search.cardBin = '';
  search.cardBinId = '';
}
const handlePageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.page = val
  getTableData()
}

const handleEdit = (index, row) => {
  cardBinForm.value = row
  dialogVisible.value = true
};

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
/* Add any custom styles here */
</style>