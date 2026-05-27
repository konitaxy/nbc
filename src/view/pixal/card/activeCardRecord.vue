<template>
    <div class="container form-white">
      <div class="row">
        <div class="col-md-6 mb-3">
          <el-input v-model="searchUserId" placeholder="输入用户ID" class="me-2" style="width: 150px;"></el-input>
          <el-button type="primary" @click="handleSearch">{{$t('lang.search')}}</el-button>
        </div>
      </div>
      <el-table :data="tableData" style="width: 100%" stripe>
        <el-table-column prop="userId" label="用户ID" width="180"></el-table-column>
        <el-table-column prop="cardCount" label="开卡数量" width="180"></el-table-column>
        <el-table-column prop="successCount" label="成功数量" width="180"></el-table-column>
        <el-table-column prop="failCount" label="失败数量" width="180"></el-table-column>
        <el-table-column prop="remark" :label="$t('lang.notes')" width="180"></el-table-column>
        <el-table-column prop="startTime" :label="$t('lang.start_time')" width="180"></el-table-column>
        <el-table-column prop="endTime" :label="$t('lang.end_time')" width="180"></el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button  type="text" @click="handleDetail(scope.$index, scope.row)">详情</el-button>
          </template>
        </el-table-column>
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
    </div>
  </template>
  
  <script setup>
  import { reactive, ref } from 'vue';
  
  const searchUserId = ref('');
  const tableData = ref([
    // 示例数据
    {
      userId: '1332',
      cardCount: 1,
      successCount: 1,
      failCount: 0,
      remark: '',
      startTime: '2025-07-11 17:10:19',
      endTime: '2025-07-11 17:10:23'
    },
    // 更多数据...
  ]);
  const search = reactive({
    searchUserId: '',
    total: 0,
    pageSize: 10,
    page:1
  });
  const handleSizeChange = (val) => {
    search.pageSize = val;
    // 更新数据
  };
  
  const handleCurrentChange = (val) => {
    search.page = val;
    // 更新数据
  };
  const handleSearch = () => {
    // 实现搜索逻辑
  };
  
  const handleDetail = (index, row) => {
    console.log('详情:', index, row);
  };
  
  
  </script>
  
  <style scoped>
  /* 自定义样式 */
  </style>