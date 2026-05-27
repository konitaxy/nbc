<template>
    <div>
    <div>
        <el-form inline label-position="top">
            <el-form-item label=" ">
                <el-input v-model="search.name" placeholder="请输入卡组名称"></el-input>
            </el-form-item>
            <el-form-item label=" ">
                <el-button type="primary" @click="getTableData">{{ $t('lang.search') }}</el-button>
            </el-form-item>
            <el-form-item label=" ">
              <el-button icon="Plus" type="primary" @click="handleAddCardGroup">{{ $t('lang.add_card_group') }}</el-button>
            </el-form-item>
        </el-form>
        
    </div>
    <div>
        <el-table :data="tableData" show-overflow-tooltip>
            <el-table-column prop="name" :label="$t('lang.card_group_name')"></el-table-column>
            <el-table-column  :label="$t('lang.update_time')">
                <template v-slot="scope">
                    {{ formatDate(scope.row.UpdatedAt) }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('lang.operations')">
            <template #default="{row}" :show-overflow-tooltip="false">
                <div>
                    <el-button type="primary" size="small" @click="handleEditGroup(row)">{{$t('lang.edit')}}</el-button>
                    <el-button type="danger" size="small" @click="handleDelGroup(row)">{{$t('lang.delete')}}</el-button>
                </div>
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

    </div>
    <el-dialog
        :title="$t('lang.add_card_group')"
        width="50%"
        v-model="dialogs.addGroupDailogVisible"
    >
            <div>
                <el-form>
                    <el-form-item  prop="name" :label="$t('lang.card_group_name')">
                        <el-input v-model="form.name"></el-input>
                    </el-form-item>
                </el-form>
            </div>
            <template #footer>
            <div class="d-flex justify-content-end">
                <el-button 
                type="info" 
                @click="dialogs.addGroupDailogVisible = false">
                {{ $t('lang.cancel') }}
                </el-button>
                <el-button 
                type="primary" 
                @click.prevent="handleAddCardGroupSubmit">
                {{ $t('lang.confirm') }}
                </el-button>
            </div>
            </template>
</el-dialog>
</div>
</template>
<script setup>
import {ref,reactive,onMounted} from 'vue'
import {listCardGroup,addCardGroup,delCardGroup} from '@/api/finance'
import {formatDate} from '@/utils/format'
import { ElMessage,ElMessageBox } from 'element-plus';
import { useI18n } from 'vue-i18n';
const { t, locale } = useI18n();
const search = reactive({
  name:'',
  page:1,
  pageSize:10,
  total:0
})
const dialogs = reactive({
  addGroupDailogVisible:false,
})

const tableData = ref([])
const getTableData = () => {
  listCardGroup(search).then(res =>{
    if(res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
    }
  })
}
onMounted(()=>{
  getTableData()
})
const form= ref({
    name:''
})
const handleAddCardGroup = () => {
    dialogs.addGroupDailogVisible = true
    form.value.name = ''
}
const handleAddCardGroupSubmit = ()=>{
    addCardGroup(form.value).then(res => {
    if(res.code === 0){
      ElMessage.success('Success')
      getTableData()
      dialogs.addGroupDailogVisible = false
    }
  })
}
const handleEditGroup =(row) =>{
  form.value = row
  dialogs.addGroupDailogVisible = true
}
const handleDelGroup = (row) =>{
    ElMessageBox.confirm(
    t('lang.del_card_group_warning'),
    t('lang.warning'),
    {
      confirmButtonText: t('lang.confirm'),
      cancelButtonText: t('lang.cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      delCardGroup({
        id:row.ID
      }).then(res =>{
        if(res.code === 0){
            ElMessage.success("Success")
            getTableData()
        }
      })
    })
    .catch(() => {
    })

}
const handleSizeChange = (val) => {
    search.pageSize = val
    getTableData()
  }
  
const handleCurrentChange = (val) => {
  search.page = val
  getTableData()
};
</script>
