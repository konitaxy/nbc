<template>
  <div class="p-2">
    <!-- Filter Form -->
    <el-form :inline="true"  label-position="top" :model="search" class="demo-form-inline">
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client No"></el-input>
          </el-form-item>
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
          <el-form-item label="Manager">
            <el-input v-model="search.manager" placeholder="Please enter account manager"></el-input>
          </el-form-item>
          <el-form-item label="Card No">
            <el-input v-model="search.cardNo" placeholder="card no or last 4 digits"></el-input>
          </el-form-item>
          <el-form-item label="Card Status">
            <el-select v-model="search.cardStatus" placeholder="Please select card status" style="width: 150px;" clearable>
              <el-option label="Pending" value="Pending"></el-option>
              <el-option label="Active" value="Active"></el-option>
              <el-option label="Frozen" value="Frozen"></el-option>
              <el-option label="Suspend" value="Suspend"></el-option>
              <el-option label="Closed" value="Closed"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label=" ">
            <el-button type="primary" @click="getTableData">Search</el-button>
            <el-button @click="reset">Reset</el-button>          
          </el-form-item>
          
    </el-form>
        
    <div  style="overflow-x: auto;">
      <el-table :data="tableData" show-overflow-tooltip	>
        <!-- <el-table-column type="selection" width="55"></el-table-column> -->
        <el-table-column prop="cardNo" label="卡号" width="160"></el-table-column>
        <el-table-column prop="client.clientNo" label="client No" width="100"></el-table-column>
        <el-table-column prop="client.email" label="Email" width="160"></el-table-column>
        <el-table-column prop="currency" label="卡币种" width="80"></el-table-column>
        <el-table-column prop="cardBrand" label="卡品牌" width="120"></el-table-column>
        <el-table-column prop="balance" label="卡片余额" width="120"></el-table-column>
        <el-table-column prop="cardStatus" label="卡状态" width="120">
          <template #default="{row}">
            <el-tag v-if="row.cardStatus==='Pending'" type="info">{{ row.cardStatus }}</el-tag>
            <el-tag v-else-if="row.cardStatus==='Active'" type="success">{{ row.cardStatus }}</el-tag>
            <el-tag v-else-if="row.cardStatus==='Frozen'" type="warning">{{ row.cardStatus }}</el-tag>
            <el-tag v-else-if="row.cardStatus==='Suspend'" type="warning">{{ row.cardStatus }}</el-tag>
            <el-tag v-else type="danger">{{ row.cardStatus }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="activeDate" label="创建时间" width="120"></el-table-column>
        <el-table-column label="操作" width="300">
            <template #default="scope">
              <el-button icon="Refresh"  size="small" @click="handleCardRefresh(scope.row)" class="me-1">refresh</el-button>
              <el-button 
                v-if="scope.row.cardStatus !== 'Frozen' && scope.row.cardStatus !== 'Suspend'"
                type="warning" 
                size="small" 
                @click="handleFrozen(scope.row)"
                class="me-1"
              >
                冻结
              </el-button>
              <el-button 
                v-if="scope.row.cardStatus === 'Frozen' || scope.row.cardStatus === 'Suspend'"
                type="success" 
                size="small" 
                @click="handleUnfrozen(scope.row)"
                class="me-1"
              >
                解冻
              </el-button>
              <el-popconfirm
                width="220"
                icon-color="#626AEF"
                title="Sure to Terminate?"
                @confirm="handleTerminateCard(scope.row)"
              >
                <template #reference>
                  <el-button type="danger" size="small">Temerated</el-button>

                </template>
              </el-popconfirm>
            <!-- <el-dropdown>
                <el-button size="small" type="primary">
                更多操作<i class="el-icon-arrow-down el-icon--right"></i>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item @click="handleRecharge(scope.row)">充值</el-dropdown-item>
                        <el-dropdown-item @click="handleCancel(scope.row)">销卡</el-dropdown-item>  
                    </el-dropdown-menu>
                </template>
            </el-dropdown> -->
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
    
    <!-- 冻结/解冻对话框 -->
    <el-dialog 
      v-model="frozenDialogVisible" 
      :title="frozenForm.action === 'frozen' ? '冻结卡片' : '解冻卡片'" 
      width="500px"
    >
      <el-form :model="frozenForm" label-width="100px">
        <el-form-item label="卡号">
          <el-input v-model="frozenForm.cardNo" disabled></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input 
            v-model="frozenForm.remark" 
            type="textarea" 
            :rows="4"
            :placeholder="frozenForm.action === 'frozen' ? '请输入冻结原因' : '请输入解冻原因'"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="frozenDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleFrozenConfirm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from '@/utils/format';
import { listCard,rechargeCard,cancelCard,syncCard,frozenCard} from '@/api/card';
import { ElMessage, ElMessageBox } from 'element-plus';
import {writeText} from 'clipboard-polyfill'

const dialogs=reactive({
  clientDDDetailDialogVisible:false
})
const search = reactive({
  clientNo: '',
  cardNo: '',
  nickName: '',
  manager: '',
  cardStatus: '',
  email:'',
  page:1,
  pageSize:10,
  total:0
});

const tableData = ref([])
const frozenDialogVisible = ref(false)
const frozenForm = reactive({
  id: null,
  action: 'frozen', // 'frozen' 或 'unfrozen'
  remark: '',
  cardNo: ''
})

const getTableData = () => {
  listCard(search).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.total = res.data.total
    }
    
  })
};

const handleTerminateCard = (row)=>{
  cancelCard({
    id:row.ID,
    clientId:row.clientId,
    cardId:row.cardId,
    cardbinId:row.cardbinId
  }).then(res=>{
    if (res.code === 0){
      ElMessage.success('Success')
      getTableData()
    }
  })
}
const handleCopy = (text) => {
  writeText(text)
  ElMessage.success('Copy Success')
};
const handleClientStatusChange = (row,status)=>{
  changeClientStatus({
    id: row.ID,
    clientStatus: status
  }).then(res =>{
    if (res.code === 0){
      getTableData()
      ElMessage.success('Success')
    }
  })
}
const reset = () => {
  search.clientNo =''
    search.cardStatus='',
    search.email='',
    search.manager='',
    search.page=1,
    search.pageSize=10,
    search.total=0
  
}
const handleCardRefresh = (row) => {
    row.loading = true
    console.log(row)
    syncCard({id:row.ID,clientId:row.client.ID}).then(res => { 
      if (res.code ===0){
        handleListCard()
      } 
      row.loading = false
    }) 
  }
onMounted(()=>{
  getTableData()
})
const handlePageSizeChange =(val) =>{
  search.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.page = val
  getTableData()
}

const handleFrozen = (row) => {
  frozenForm.id = row.ID
  frozenForm.action = 'frozen'
  frozenForm.remark = ''
  frozenForm.cardNo = row.cardNo
  frozenDialogVisible.value = true
}

const handleUnfrozen = (row) => {
  frozenForm.id = row.ID
  frozenForm.action = 'unfrozen'
  frozenForm.remark = ''
  frozenForm.cardNo = row.cardNo
  frozenDialogVisible.value = true
}

const handleFrozenConfirm = () => {
  frozenCard({
    id: frozenForm.id,
    action: frozenForm.action,
    remark: frozenForm.remark || ''
  }).then(res => {
    if (res.code === 0) {
      ElMessage.success(frozenForm.action === 'frozen' ? '冻结成功' : '解冻成功')
      frozenDialogVisible.value = false
      getTableData()
    }
  })
}
</script>
<style scoped>
.container {
  padding: 10px;
}
/* :deep(.el-select__wrapper) {
    min-width: 150px;
} */
:deep(.el-dialog .el-dialog__body){
  padding: 0px;
}
</style>