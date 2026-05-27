<template>
    <div class="col col-12 page-title"><h5>Transfer Management</h5></div>
    <div class="form-list-container">
      <!-- 操作栏 -->
      <div class="operation-bar">
        <el-input v-model="search.orderNumber" :placeholder="$t('lang.enter_order_number')" style="max-width: 140px;"></el-input>
        <el-select v-model="search.transcationType" :placeholder="$t('lang.please_select_audit_status')" @change="handleSearch" style="width: 150px;" clearable>
          <el-option label="转账" :value="1"></el-option>
          <el-option label="收款" :value="2"></el-option>
        </el-select>
        <el-select v-model="search.auditStatus" :placeholder="$t('lang.please_select_audit_status')" @change="handleSearch" style="width: 150px;" clearable>
          <el-option label="待审核" :value="1"></el-option>
          <el-option label="通过" :value="2"></el-option>
          <el-option label="驳回" :value="3"></el-option>
        </el-select>
        <el-button class="ms-2" type="primary" icon="search" @click="handleSearch">{{$t('lang.search')}}</el-button>
        <el-button type="primary" icon="plus" @click="handleRecharge">转账</el-button>
        <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="handleExport">{{$t('lang.export')}}</el-button>
      </div>
  
      <!-- 充值记录表格 -->
      <el-table :data="data" style="width: 100%">
        <el-table-column prop="orderNumber" :label="$t('lang.order_number')" width="280"></el-table-column>
        <el-table-column prop="rechargeAmount" label="充值金额"></el-table-column>
        <el-table-column prop="arrivalAmount" label="到账金额"></el-table-column>
        <el-table-column prop="rechargeTime" label="充值时间"></el-table-column>
        <el-table-column prop="auditStatus" label="审核状态"></el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
  
      <!-- 分页 -->
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
      <el-dialog title="查看" v-model="dialogs.rechargeDetailDialog" width="30%"  align-center>
        <div class="d-flex flex-column gap-2">
            <p><span>充值单号:</span> {{ currentRecord.orderNumber }}</p>
            <p><span>审核状态:</span> {{ currentRecord.auditStatus }}</p>
            <p><span>充值金额:</span> {{ currentRecord.rechargeAmount }}</p>
            <p><span>到账金额:</span> {{ currentRecord.arrivalAmount }}</p>
            <p><span>充值时间:</span> {{ currentRecord.rechargeTime }}</p>
            <p><span>审核备注:</span> {{ currentRecord.auditRemark }}</p>
            <p><span>充值备注:</span> {{ currentRecord.rechargeRemark }}</p>
        </div>
        </el-dialog>
        <el-dialog title="充值" v-model="dialogs.rechargeDialogVisible" width="40%" align-center>
            <el-alert
                title="不接受任何个人账户转账"
                type="warning"
                show-icon
                :closable="false"
            ></el-alert>
            <el-form :model="rechargeForm" ref="rechargeFormRef" :rules="rechargeRules" class="mt-2">
                <el-form-item label="账户信息" prop="accountInfo">
                <el-input v-model="rechargeForm.accountInfo" disabled></el-input>
                </el-form-item>
                <el-form-item label="银行账号" prop="bankAccount">
                <el-input v-model="rechargeForm.bankAccount" disabled></el-input>
                </el-form-item>
                <el-form-item label="SWIFT code/BIC" prop="swiftCode">
                <el-input v-model="rechargeForm.swiftCode" disabled></el-input>
                </el-form-item>
                <el-form-item label="银行代码" prop="bankCode">
                <el-input v-model="rechargeForm.bankCode" disabled></el-input>
                </el-form-item>
                <el-form-item label="分行代码" prop="branchCode">
                <el-input v-model="rechargeForm.branchCode" disabled></el-input>
                </el-form-item>
                <el-form-item label="银行所在地" prop="bankLocation">
                <el-input v-model="rechargeForm.bankLocation" disabled></el-input>
                </el-form-item>
                <el-form-item :label="$t('lang.amount')" prop="amount">
                <el-input v-model.number="rechargeForm.amount" placeholder="请输入充值金额"></el-input>
                </el-form-item>
                <el-form-item label="充值凭证" prop="voucher">
                    <el-upload ref="elUpladRef" action="#" @on-change="handleFileChange" list-type="picture-card" :auto-upload="false" :limit="1">
                            <el-icon><Plus/></el-icon>

                            <template #file="{ file }">
                            <div>
                                <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
                                <span class="el-upload-list__item-actions">
                                
                                <span
                                    class="el-upload-list__item-delete"
                                    @click="handleFileRemove(file)"
                                >
                                    <el-icon><Delete /></el-icon>
                                </span>
                                </span>
                            </div>
                            </template>
                        </el-upload>
                </el-form-item>
                <el-form-item :label="$t('lang.notes')" prop="remark">
                <el-input type="textarea" v-model="rechargeForm.remark" placeholder="请输入备注"></el-input>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="dialogs.rechargeDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="submitRecharge">提交</el-button>
                </span>
            </template>
            </el-dialog>
    </div>
  </template>
  
  <script setup>
  import {reactive, ref} from 'vue'
  const dialogs = reactive({
    rechargeDetailDialog:false,
    rechargeDialogVisible:false
  })
  const search = reactive({
    auditStatus: null,
    transcationType: null,
    total: 0,
    pageSize: 10,
    page:1
  });

 const data = ref([
          {
            orderNumber: '12025051120440112162292226070271',
            rechargeAmount: '500.00',
            arrivalAmount: '500.00',
            rechargeTime: '2025-05-11 20:44:01',
            auditStatus: '通过'
          },
          {
            orderNumber: '12025031114594127062295155617015',
            rechargeAmount: '500.00',
            arrivalAmount: '500.00',
            rechargeTime: '2025-03-11 14:59:41',
            auditStatus: '通过'
          }
        ])

    const rechargeForm = ref({
        accountInfo: 'HT International Trade Service Limited',
        bankAccount: '79821010003964',
        swiftCode: 'DHBKHKHH',
        bankCode: '016',
        branchCode: '478',
        bankLocation: '中国香港',
        amount: null
        });
        const rechargeRules = {
        amount: [
            { required: true, message: '请输入充值金额', trigger: 'blur' },
            { type: 'number', message: '充值金额必须为数字', trigger: 'blur' }
        ]
        };
      const handleSearch =() =>{
        // 处理搜索逻辑
        console.log('搜索条件:', search);
      }
      const handleRecharge =() =>{
        // 处理充值逻辑
        dialogs.rechargeDialogVisible = true;
      }
      const handleExport =() =>{
        // 处理导出逻辑
        console.log('点击了导出按钮');
      }
      const currentRecord = ref()
      const handleView =(val) => {
        // 处理查看逻辑

        currentRecord.value = val
        dialogs.rechargeDetailDialog = true
      }
      const handleSizeChange=(val) => {
        // 处理每页显示条数变化
        console.log(`每页 ${val} 条`);
        this.pageSize = val;
      }
      const  handleCurrentChange =(val) => {
        // 处理当前页码变化
        console.log(`当前页: ${val}`);
        this.currentPage = val;
      }
      const submitRecharge =() =>{

      }
      const elUpladRef = ref(null);
      const handleFileChange = (file,fl) => {
        rechargeForm.value.voucher = file.raw
        console.log(rechargeForm.value.voucher)
        };
        const handleFileRemove =(file) =>{
            elUpladRef.value.clearFiles();
        }

        const beforeUpload = (file) => {
        const isJPG = file.type === 'image/jpeg';
        const isPNG = file.type === 'image/png';
        const isLt2M = file.size / 1024 / 1024 < 2;

        if (!isJPG && !isPNG) {
            ElMessage.error('上传图片只能是 JPG 或 PNG 格式!');
            return false;
        }
        if (!isLt2M) {
            ElMessage.error('上传图片大小不能超过 2MB!');
            return false;
        }
        return true;
        };
  </script>
  
  <style scoped>
  
  .operation-bar {
    margin-bottom: 20px;
    display: flex;
    /* display: inline; */
    gap: 10px;
  }
 
  </style>