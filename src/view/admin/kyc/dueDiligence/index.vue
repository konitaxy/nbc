<template>
  <div class="p-2">
    <!-- Filter Form -->
    <el-form :inline="true"  label-position="top" :model="search" class="demo-form-inline">
      
          <el-form-item label="Client No">
            <el-input v-model="search.clientNo" placeholder="Please enter client no"></el-input>
          </el-form-item>
        
          <el-form-item label="Location Name">
            <el-input v-model="search.location" placeholder="Please enter location name"></el-input>
          </el-form-item>
        
          <!-- <el-form-item label="English Name">
            <el-input v-model="search.enName" placeholder="Please enter english name"></el-input>
          </el-form-item> -->
      
          <el-form-item label="Account Manager">
            <el-input v-model="search.accountManager" placeholder="Please enter account manager"></el-input>
          </el-form-item>
<!--         
          <el-form-item label="Client Status">
           
            <el-select v-model="search.clientStatus" placeholder="Please select client status">
              <el-option label="Review" :value="1"></el-option>
              <el-option label="Active" :value="2"></el-option>
              <el-option label="Suspend" :value="3"></el-option>

            </el-select>
          </el-form-item> -->
        
          <el-form-item label="Review Status">
            <el-select v-model="search.clientReviewStatus" placeholder="Please select review status">
              <el-option label="Unreview" :value="1"></el-option>
              <el-option label="Reviewing" :value="2"></el-option>
              <el-option label="Completed" :value="3"></el-option>

            </el-select>
          </el-form-item>
        
          <!-- <el-form-item label="Client Risk Level">
            <el-select v-model="search.riskLevel" placeholder="Please select client risk level">
              <el-option label="Low" value="low"></el-option>
              <el-option label="Medium" value="medium"></el-option>
              <el-option label="High" value="high"></el-option>
            </el-select>
          </el-form-item> -->
          <el-form-item label="Email">
            <el-input v-model="search.email" placeholder="Please enter email"></el-input>
          </el-form-item>
          <el-form-item label=" ">
            <el-button type="primary" @click="getTableData">Search</el-button>
            <el-button @click="reset">Reset</el-button>          
          </el-form-item>
      <!-- Add more rows as needed -->
          
    </el-form>
        
       

    <!-- Table -->
    <div  style="overflow-x: auto;">
      <el-table :data="tableData" show-overflow-tooltip	>
        <el-table-column prop="clientNo" label="Client No" min-width="120">
          <template #default="{row}">
            <el-link style="cursor: pointer;font-weight: 600;" @click="handleOpenDDDetal(row)"><i class="bi bi-link-45deg"></i>{{ row.clientNo }}</el-link>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="enName" label="English Name" width="120"></el-table-column> -->
        <el-table-column prop="email" label="Email" min-width="120"></el-table-column>
        <el-table-column prop="accountManager" label="Manager"  min-width="120"></el-table-column>
        <el-table-column prop="clientType" label="Client Type" min-width="120"></el-table-column>
        <el-table-column prop="location" label="Location" min-width="90"></el-table-column>
        <el-table-column prop="clientStatus" label="Client Status" min-width="120">
          <template #default="{row}">
            <el-tag :type="row.clientStatus == 1?'info':row.clientStatus == 2?'success':'danger'">{{row.clientStatus == 1?'Review':row.clientStatus == 2?'Active':'Suspend'}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reviewStatus" label="Review Status" min-width="120">
          <template #default="{row}">
            <el-tag :type="row.clientReviewStatus == 1?'info':row.clientReviewStatus == 2?'primary':'success'">{{row.clientReviewStatus == 1?'UnReview':row.clientReviewStatus == 2?'Reviewing':'Completed'}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="Register Time" min-width="120">
          <template #default="{row}">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="Remark" min-width="120">
          <template #default="{row}">
            <i @click="handleEditRemark(row)" class="bi bi-pencil-square"></i>{{ row.remark }}
          </template>

        </el-table-column>
        <el-table-column prop="action" label="Action" fixed="right" min-width="150">
          <template #default="{row}">           
              <el-button type="primary" size="small" @click="handleOpenDDDetal(row)">Audit</el-button>
        </template>
        </el-table-column>
      </el-table>
      <div class="mt-4">
      <el-pagination
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="search.total"
        :page-sizes="[20, 50, 100]"
        :page-size="search.pageSize"
        :current-page="search.page"
        @size-change="handlePageSizeChange"
        @current-change="handlePageChange"
      ></el-pagination>
    </div>
    </div>
    <el-dialog
      title="Client DD Detail"
      width="80%"
      v-model="dialogs.clientDDDetailDialogVisible"
      align-center
    >
      
      <div class="ps-2">
        <!-- Basic Information Section -->
        <el-collapse v-model="activeNames">
          <el-collapse-item title="Basic Information" name="1">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Client No:</strong></el-col>
                  <el-col :span="16">{{ ddDetail.clientNo }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Status:</strong></el-col>
                  <el-col :span="16"><el-tag :type="ddDetail.clientStatus == 1?'info':ddDetail.clientStatus == 2?'success':'danger'">{{ ddDetail.clientStatus == 1?'Review':ddDetail.clientStatus == 2?'Active':'Suspence' }}</el-tag></el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Account Manager:</strong></el-col>
                  <el-col :span="16">{{ ddDetail.accountManager }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Risk Result:</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Client IP Address:</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Risk Level:</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Type:</strong></el-col>
                  <el-col :span="16"><el-tag type="success">{{ ddDetail.clientType }}</el-tag></el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>DD Times:</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.ddTimes }}</el-col>
                </el-row>
              </el-col>
            </el-row>
          </el-collapse-item>
        </el-collapse>

        <!-- Customer Identity Program Section -->
        <el-collapse v-model="activeNames">
          <el-collapse-item title="Customer Identity Program" name="2">
            <el-row :gutter="20" v-if="ddDetail.clientType == 'enterprise'">
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Enterprise Type</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entEnterpriseType }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Enterprise Chinese Name</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entEnterpriseChineseName }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Enterprise English Name</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entEnterpriseEnglishName }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Business Registration Form</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entEnterpriseEnglishName }}</el-col>
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Business Registration NO.</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entBusinessRegistrationNo }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Date of Establishment</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entDateOfEstablishment }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Date of Expiration</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entDateOfExpiration }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Type:</strong></el-col>
                  <el-col :span="16"><el-tag type="success">Enterprise</el-tag></el-col>
                </el-row>
              </el-col>
                <el-col :span="12">
                  <el-row>
                    <el-col :span="8"><strong>Local Business Premise</strong></el-col>
                    <el-col :span="16">{{ ddDetail.dueDiligence.entLocalBusinessPremise }}</el-col>
                    </el-row>
                    <el-row>
                      <el-col :span="8"><strong>Business Address Proof</strong></el-col>
                      <el-col :span="16"><el-image
                          style="width: 100px; height: 100px"
                          :src="ddDetail.dueDiligence.entBusinessAddressProof"
                          :zoom-rate="1.6"
                          :max-scale="2"
                          :min-scale="0.8"
                          :preview-src-list="[ddDetail.dueDiligence.entBusinessAddressProof]"
                          show-progress
                          :initial-index="4"
                          fit="cover"
                        /></el-col>
                    </el-row>
                    <el-row>
                      <el-col :span="8"><strong>Province</strong></el-col>
                      <el-col :span="16">{{ ddDetail.dueDiligence.entProvince }}</el-col>
                    </el-row>
                    <el-row>
                      <el-col :span="8"><strong>City</strong></el-col>
                      <el-col :span="16">{{ ddDetail.city }}</el-col>
                    </el-row>
                </el-col>
                <el-col :span="12">
                  <el-row>
                    <el-col :span="8"><strong>Listed Company</strong></el-col>
                    <el-col :span="16">{{ ddDetail.dueDiligence.entListedCompany }}</el-col>
                  </el-row>
                  <el-row>
                    <el-col :span="8"><strong>State Owned</strong></el-col>
                    <el-col :span="16">{{ ddDetail.dueDiligence.entStateOwned }}</el-col>
                  </el-row>
                  <el-row>
                    <el-col :span="8"><strong>Foreign-invested</strong></el-col>
                    <el-col :span="16">{{ ddDetail.dueDiligence.entForeignInvested }}</el-col>
                  </el-row>
                  <el-row>
                    <el-col :span="8"><strong>Shareholder Structure</strong></el-col>
                    <el-col :span="16">{{ ddDetail.dueDiligence.entShareholderStructure }}</el-col>
                  </el-row>
                </el-col>
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Operation Address</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.entOperationAddress }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Risk Level:</strong></el-col>
                  <el-col :span="16">{{ ddDetail.riskLevel }}</el-col>
                </el-row>
                
              </el-col>
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Client IP Address:</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Client Risk Level:</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
                
                
              </el-col>
              <div>
              <hr>
            </div>
            </el-row>
            <!-- <el-button class="btn btn-danger" style="min-width: 90px;">Edit</el-button> -->
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>Country or Region</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indCountryOrRegion }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Position</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indPosition }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>English Name</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indEnglishName }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>ID Back End</strong></el-col>
                  <el-col :span="16"><el-image
                    style="width: 100px; height: 100px"
                    :src="ddDetail.dueDiligence.indIDBackEnd"
                    :zoom-rate="1.6"
                    :max-scale="2"
                    :min-scale="0.8"
                    :preview-src-list="[ddDetail.dueDiligence.indIDBackEnd]"
                    show-progress
                    :initial-index="4"
                    fit="cover"
                  /></el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Issue Date</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indIssueDate }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Date of Birth</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indDateOfBirth }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>City</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indCity  }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Adverse Media</strong></el-col>
                  <el-col :span="16">2</el-col>
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-row>
                  <el-col :span="8"><strong>CID Type</strong></el-col>
                  <el-col :span="16">-</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Chinese Name</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indChineseName }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Identification No.</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indIdentificationNo }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>ID Front End</strong></el-col>
                  <el-col :span="16">
                    <el-image
                    style="width: 100px; height: 100px"
                    :src="ddDetail.dueDiligence.indIDFrontEnd"
                    :zoom-rate="1.6"
                    :max-scale="2"
                    :min-scale="0.8"
                    :preview-src-list="[ddDetail.dueDiligence.indIDFrontEnd]"
                    show-progress
                    :initial-index="4"
                    fit="cover"
                  /></el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Expiration Date</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indExpirationDate }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Province or state</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indProvinceOrState }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Residential Address</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indResidentialAddress }}</el-col>
                </el-row>
                <el-row>
                  <el-col :span="8"><strong>Reliability of Documents</strong></el-col>
                  <el-col :span="16">{{ ddDetail.dueDiligence.indReliabilityOfDocuments }}</el-col>
                </el-row>
              </el-col>
            </el-row>
            <!-- <el-button class="btn btn-danger" style="min-width: 90px;">Edit</el-button> -->

          </el-collapse-item>
        </el-collapse>

        <!-- CDD Final Result Section -->
        <!-- <el-collapse v-model="activeNames">
          <el-collapse-item title="CDD Final Result" name="3">
            <el-tag type="success">{{ ddDetail.clientReviewStatus = }}</el-tag>
          </el-collapse-item>
        </el-collapse> -->

        <!-- Buttons Section -->
        <div style="text-align: center; margin-top: 20px;">
          <el-button type="primary" @click="dialogs.clientDDDetailDialogVisible = false">Go Back</el-button>
          <el-button type="info" @click="needEnhancedKYB">Need Enhanced KYB</el-button>
          <el-button type="success" @click="handleChangeReviewStatus(3)">Accept</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script setup>
import { reactive, ref,onMounted } from 'vue';
import {formatDate} from '@/utils/format';
import { listClient,setClientName,remarkClient,ddClient,reviewClient,changeClientStatus,enhancedKYB,getDueDiligence} from '@/api/client';
import { ElMessageBox,ElMessage } from 'element-plus';
const dialogs=reactive({
  clientDDDetailDialogVisible:false
})
const activeNames = ref(['1','2']);
const search = ref({
  clientNo: '',
  location: '',
  enName: '',
  accountManager: '',
  reviewStatus: '',
  email:'',
  page:1,
  pageSize:20,
  total:0
});
onMounted(()=>{
  getTableData()
})
const tableData = ref([])

const getTableData = () => {
  listClient(search.value).then(res => {
    if (res.code === 0){
      tableData.value = res.data.list
      search.value.total = res.data.total
    }
    
  })
};
const handleChangeReviewStatus = (status) => {
  console.log(ddDetail.value)
  reviewClient({
    id: ddDetail.value.ID,
    clientReviewStatus: status
  }).then(res => {
    if (res.code === 0){
      ElMessage.success('操作成功')
      getTableData()
    }
  })
}

const reset = () => {
  search.value = {
    clientId: '',
    clientNo: '',
    locationName: '',
    englishName: '',
    accountManager: '',
    clientStatus: '',
    reviewStatus: '',
    riskLevel: ''
  };
};
const handleEditRemark = (row) => {
  ElMessageBox.prompt('Please enter your remark', 'Edit Remark', {
    confirmButtonText: 'Save',
    cancelButtonText: 'Cancel',
    inputValue: row.remark,
  }
  ).then((val)=>{
    remarkClient({
      id: row.ID,
      remark: val.value
    }).then(res =>{
      if (res.code === 0){
        ElMessage({
          type: 'success',
          message: 'Edit Success'
        });
      }
    })
  })
}
const ddDetail = ref()
const handleOpenDDDetal =(row) =>{
  getDueDiligence({
    id: row.ID
  }).then(res =>{ 
    if (res.code === 0){
      console.log(res.data)
      ddDetail.value = row
      ddDetail.value.dueDiligence = res.data
      dialogs.clientDDDetailDialogVisible = true
    }
  })

}
const needEnhancedKYB = () =>{
  ElMessageBox.prompt('Please enter some tip', 'kyb enhanced', {
    confirmButtonText: 'Confirm',
    cancelButtonText: 'Cancel',
    inputValue: ddDetail.value.tip,
  }
  ).then((val)=>{
    enhancedKYB({
      id: ddDetail.value.ID,
      tip: val.value,
    }).then(res =>{
      if (res.code === 0){
        ElMessage({
          type: 'success',
          message: 'Success'
        });
      }
    })
  })
}
const handlePageSizeChange =(val) =>{
  search.value.pageSize = val
  getTableData()
}
const handlePageChange =(val) =>{
  search.value.page = val
  getTableData()
}
</script>
<style scoped>
.container {
  padding: 10px;
}
:deep(.el-select__wrapper) {
    min-width: 150px;
}
:deep(.el-dialog .el-dialog__body){
  padding: 0px;
}
</style>