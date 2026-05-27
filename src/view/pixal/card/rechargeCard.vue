<template>
    <div class="card-management form-white">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('lang.virtual_card')" name="active">
        
<div class="search-bar">
  <!-- 搜索栏 - 卡号后四位 -->
    <el-input 
      :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')" 
      v-model="search.cardNoSuffix" 
      :sm="2" :md="2" :lg="2"
      style="max-width: 150px;">
    </el-input>
    <el-input 
      :placeholder="$t('lang.please_enter_notes')" 
      v-model="search.remark" 
      :sm="2" :md="2" :lg="2"
      style="max-width: 180px;">
    </el-input>
    <el-select 
     v-if="$userStore.hasRole(100)"
      v-model="search.iamUserId" 
      @change="handleListCard"
      :placeholder="$t('lang.select_belong_to')" 
      filterable
      clearable>
      <el-option v-for="item in iamUsers" :key="item.ID" :label="item.nickname || item.email" :value="item.ID"></el-option>
    </el-select>
    <!-- 搜索栏 - 卡状态 -->
    <el-select 
      v-model="search.cardStatus" 
      @change="handleListCard"
      :placeholder="$t('lang.select_card_status')" 
      clearable>
      <el-option :label="$t('lang.pending')" value="Pending"></el-option>
      <el-option :label="$t('lang.active')" value="Active"></el-option>
      <el-option :label="$t('lang.suspend')" value="Suspend"></el-option>
      <el-option :label="$t('lang.failure')" value="Failure"></el-option>
      <el-option :label="$t('lang.terminated')" value="Closed"></el-option>
    </el-select>

    <el-select 
      v-model="search.groupId" 
      @change="handleListCard"
      :placeholder="$t('lang.card_group_name')" 
      clearable>
      <el-option v-for="item in cardGroups" :label="item.name" :value="item.ID"></el-option>
    </el-select>
    <el-select 
      v-model="search.cardBrand" 
      @change="handleListCard"
      :placeholder="$t('lang.card_brand')" 
      clearable>
      <el-option value="Visa" label="Visa"></el-option>
      <el-option value="MasterCard" label="MasterCard"></el-option>
    </el-select>


  <!-- 搜索栏 - 日期范围 -->
  <el-date-picker
    v-model="search.timeRange"
    type="daterange"
    format="YYYY-MM-DD"
    value-format="YYYY-MM-DD"
    :range-separator="$t('lang.to')" 
    :start-placeholder="$t('lang.start_time')" 
    :end-placeholder="$t('lang.end_time')"
    @change="handleListCard"
    class="col col-lg-3">
  </el-date-picker>




  <!-- 搜索栏 - 余额区间 -->
  <div style="flex-wrap: nowrap; display: flex; align-items: center; overflow: auto;">
    <span style="padding-right: 5px;">{{ $t('lang.balance_range') }}:</span> <!-- "余额区间" -->
    <el-input-number 
      v-model="search.minBalance" 
      :min="0" 
      :max="100000" 
      :controls="false">
    </el-input-number>
    <span>&nbsp;{{ $t('lang.to') }}&nbsp;</span> <!-- "至" -->
    <el-input-number 
      v-model="search.maxBalance" 
      :min="search.minBalance" 
      :max="100000"  
      :controls="false">
    </el-input-number>
  </div>
 

  <!-- 搜索栏 - 按钮组 -->
  <div class="col col-12 form-actions">
    <el-button type="primary" icon="search" @click="handleListCard">{{ $t('lang.search') }}</el-button>
    <el-button v-if="$userStore.hasRole(5)" type="primary" icon="plus" @click="handleActiveCard">{{ $t('lang.open_card') }}</el-button>
    <el-button v-if="$userStore.hasRole(5)" type="danger" icon="delete" @click="handleBatchCancelCard">{{ $t('lang.batch_cancel_cards') }}</el-button>
    <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="onExport">{{ $t('lang.export') }}</el-button>
  </div>
</div>

<!-- 表格 -->
<el-table ref="cardListTableRef" :data="tableData" style="width: 100%" border show-overflow-tooltip row-key="ID" @selection-change="onCancelSelectionChange">
  <el-table-column v-if="$userStore.hasRole(5)" type="selection" width="55" :reserve-selection="false"></el-table-column>
  
  <!-- 表格列 - 卡号 -->
  <el-table-column prop="cardNo" :label="$t('lang.card_number')" width="200">
    <template #default="{row}">
      <el-link v-if="row.cardNo" type="primary" icon="view" @click="handleViewDetail(row)">{{ row.cardNo }}</el-link>
    </template>
  </el-table-column>

  <!-- 表格列 - CVV -->
  <!-- <el-table-column :label="$t('lang.cvv')" width="80">
    <template #default="{ row }">
      <span v-if="row.showCVV">{{ row.cvv }}</span>
      <span v-else>***</span>
    </template>
  </el-table-column> -->

  <!-- 表格列 - 有效期 -->
  <!-- <el-table-column :label="$t('lang.expiry_date')" width="120">
    <template #default="{ row }">
      {{ row.activeDate?addYear(row.activeDate):'' }}
    </template>
  </el-table-column> -->
  
  <!-- 表格列 - 卡状态 -->
  <el-table-column prop="cardStatus" :label="$t('lang.card_status')" width="120">
    <template #default="{ row }">
      <el-tag 
        :type="row.cardStatus === 'Active' ? 'success' : row.cardStatus === 'Failure' ?'danger': row.cardStatus === 'Suspend' ? 'warning' : 'info'">
        {{ 
          row.cardStatus === 'Active' ? $t('lang.active') : row.cardStatus === 'Failure' ? $t('lang.failure') : row.cardStatus === 'Closed' ? $t('lang.terminated') : row.cardStatus === 'Suspend' ? $t('lang.suspend') : row.cardStatus 
        }}
      </el-tag>
    </template>
  </el-table-column>

  <!-- 表格列 - 卡片余额 -->
  <el-table-column prop="balance" :label="$t('lang.card_balance')" width="120"></el-table-column>
  <!-- 表格列 - 卡币种 -->
  <el-table-column prop="currency" :label="$t('lang.card_currency')" width="100"></el-table-column>
  <!-- 表格列 - 创建时间 -->
  <el-table-column prop="activeDate" :label="$t('lang.creation_time')" width="150"></el-table-column>

  <!-- 表格列 - 卡品牌 -->
  <el-table-column prop="cardBrand" :label="$t('lang.card_brand')" width="120"></el-table-column>

  <!-- 表格列 - 备注 -->
  <el-table-column prop="remark" :label="$t('lang.notes')" width="200">
    <template #default="scope">
      <div class="text-nowrap">
        <i class="bi bi-pencil-square me-1" style="cursor: pointer;" @click="handleEditRemark(scope.row)"></i>{{ scope.row.remark }}
      </div>
    </template>
  </el-table-column>
  <!-- 表格列 - 备注 -->
  <el-table-column prop="group.name" :label="$t('lang.card_group_name')" width="200">
    <template #default="scope">
      <div class="text-nowrap">
        <i class="bi bi-pencil-square me-1" style="cursor: pointer;" @click="handleEditGroup(scope.row)"></i>{{ scope.row.group?.name }}
      </div>
    </template>
  </el-table-column>

  <!-- 表格列 - 归属 -->
  <el-table-column prop="iamUserName" :label="$t('lang.belong_to')" width="150"></el-table-column>

  <!-- 表格列 - 操作 -->
  <el-table-column :label="$t('lang.actions')"  :width="isMobile ? 110 : (locale === 'zh' ? 270 : 350)"  fixed="right" :show-overflow-tooltip="false">
    <template #default="scope">
      <el-dropdown class="d-block d-sm-none">
        <el-button size="mini" type="primary">
          {{ $t('lang.actions') }}<i class="el-icon-arrow-down el-icon--right"></i>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item v-if="$userStore.hasRole(5) && scope.row.cardLevel !== 'SubCard'" @click="handleRechargeCard(scope.row)">{{ $t('lang.recharge') }}</el-dropdown-item>
            <el-dropdown-item v-if="$userStore.hasRole(5) && scope.row.cardLevel === 'SubCard'" @click="handleAdjustLimit(scope.row)">{{ $t('lang.adjust_limit') }}</el-dropdown-item>
            <el-dropdown-item v-if="$userStore.hasRole(5) && scope.row.cardLevel !== 'SubCard'" @click="handleWithdrawCard(scope.row)">{{ $t('lang.withdraw_card') }}</el-dropdown-item>
            <el-dropdown-item v-if="$userStore.hasRole(5)" @click="handleCancelCard(scope.row)">{{ $t('lang.terminate_card') }}</el-dropdown-item>
            <el-dropdown-item @click="handleRefresh(scope.row)">{{ $t('lang.refresh') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <div class="action-buttons d-none d-sm-flex display-flex"> 
        <el-button v-if="$userStore.hasRole(5) && scope.row.cardLevel !== 'SubCard'" type="primary" size="small" @click="handleRechargeCard(scope.row)">{{ $t('lang.top_up') }}</el-button>
        <el-button v-if="$userStore.hasRole(5) && scope.row.cardLevel === 'SubCard'" type="primary" size="small" @click="handleAdjustLimit(scope.row)">{{ $t('lang.adjust_limit') }}</el-button>
        <el-button v-if="$userStore.hasRole(5) && scope.row.cardLevel !== 'SubCard'" type="secondary" size="small" @click="handleWithdrawCard(scope.row)">{{ $t('lang.withdraw_card') }}</el-button>
        <el-button v-if="$userStore.hasRole(5)" type="danger" size="small" @click="handleCancelCard(scope.row)">{{ $t('lang.terminate') }}</el-button>
        <el-button :loading="scope.row.loading" type="success" size="small" @click="handleRefresh(scope.row)">{{ $t('lang.refresh') }}</el-button>
      </div>
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
        :pager-count="5"
        :current-page="search.page"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      ></el-pagination>
    </el-tab-pane>
    <el-tab-pane :label="$t('lang.block_card_list')" name="block">
             
      <div class="search-bar">
        <!-- 搜索栏 - 卡号后四位 -->
          <el-input 
            :placeholder="$t('lang.please_enter_last_four_digits_of_card_number')" 
            v-model="search2.cardNoSuffix" 
            :sm="2" :md="2" :lg="2"
            style="max-width: 200px;">
          </el-input>

          

        <!-- 搜索栏 - 日期范围 -->
        <el-date-picker
          v-model="search2.timeRange"
          type="daterange"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          :range-separator="$t('lang.to')" 
          :start-placeholder="$t('lang.start_time')" 
          :end-placeholder="$t('lang.end_time')"
          @change="handleListCard"
          class="col col-lg-3">
        </el-date-picker>

        <!-- 搜索栏 - 按钮组 -->
        <div class="col col-12 form-actions">
          <el-button type="primary" icon="search" @click="getTableData2">{{ $t('lang.search') }}</el-button>
          <!-- <el-button type="primary" icon="plus" @click="handleActiveCard">{{ $t('lang.open_card') }}</el-button> -->
          <!-- <el-button type="danger" icon="delete" @click="batchCancelCards">{{ $t('lang.batch_terminate_cards') }}</el-button> -->
          <el-button v-if="$userStore.hasRole(6)" type="primary" icon="download" @click="onExport">{{ $t('lang.export') }}</el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData2" style="width: 100%" border show-overflow-tooltip>
        <!-- <el-table-column type="selection" width="55"></el-table-column> -->
        
        <!-- 表格列 - 卡号 -->
        <el-table-column prop="cardNo" :label="$t('lang.card_number')" width="200">
          <template #default="{row}">
            <el-link v-if="row.cardNo" type="primary" icon="view" @click="handleViewDetail(row)">{{ row.cardNo }}</el-link>
          </template>
        </el-table-column>

        <!-- 表格列 - CVV -->
        <!-- <el-table-column :label="$t('lang.cvv')" width="80">
          <template #default="{ row }">
            <span v-if="row.showCVV">{{ row.cvv }}</span>
            <span v-else>***</span>
          </template>
        </el-table-column> -->

        <!-- 表格列 - 有效期 -->
        <el-table-column :label="$t('lang.expiry_date')" width="120">
          <template #default="{ row }">
            {{ row.activeDate?addYear(row.activeDate):'' }}
          </template>
        </el-table-column>
        
        <!-- 表格列 - 卡状态 -->
        <el-table-column prop="cardStatus" :label="$t('lang.card_status')" width="120">
          <template #default="{ row }">
            <el-tag 
              :type="row.cardStatus === 'Active' ? 'success' : row.cardStatus === 'Failure' ?'danger': row.cardStatus === 'Suspend' ? 'warning' : 'info'">
              {{ 
                row.cardStatus === 'Active' ? $t('lang.active') : row.cardStatus === 'Failure' ? $t('lang.failure') : row.cardStatus === 'Closed' ? $t('lang.terminated') : row.cardStatus === 'Suspend' ? $t('lang.suspend') : row.cardStatus 
              }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 表格列 - 卡片余额 -->
        <el-table-column prop="balance" :label="$t('lang.card_balance')" width="120"></el-table-column>

        <!-- 表格列 - 创建时间 -->
        <el-table-column prop="activeDate" :label="$t('lang.creation_time')" width="150"></el-table-column>

        <!-- 表格列 - 卡币种 -->
        <el-table-column prop="currency" :label="$t('lang.card_currency')" width="100"></el-table-column>

        <!-- 表格列 - 卡品牌 -->
        <el-table-column prop="cardBrand" :label="$t('lang.card_brand')" width="120"></el-table-column>

        <!-- 表格列 - 备注 -->
        <el-table-column prop="remark" :label="$t('lang.notes')" width="200">
          <template #default="scope">
            <div class="text-nowrap">
              <i class="bi bi-pencil-square me-1" style="cursor: pointer;" @click="handleEditRemark(scope.row)"></i>{{ scope.row.remark }}
            </div>
          </template>
        </el-table-column>

        <!-- <el-table-column :label="$t('lang.actions')"  :width="isMobile ? 110 : (locale === 'zh' ? 270 : 350)"  fixed="right" :show-overflow-tooltip="false">
          <template #default="scope">
            <el-dropdown class="d-block d-sm-none">
              <el-button size="mini" type="primary">
                {{ $t('lang.actions') }}<i class="el-icon-arrow-down el-icon--right"></i>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="scope.row.cardLevel !== 'SubCard'" @click="handleRechargeCard(scope.row)">{{ $t('lang.recharge') }}</el-dropdown-item>
                  <el-dropdown-item v-if="scope.row.cardLevel === 'SubCard'" @click="handleAdjustLimit(scope.row)">{{ $t('lang.adjust_limit') }}</el-dropdown-item>
                  <el-dropdown-item v-if="scope.row.cardLevel !== 'SubCard'" @click="handleWithdrawCard(scope.row)">{{ $t('lang.withdraw_card') }}</el-dropdown-item>
                  <el-dropdown-item @click="handleCancelCard(scope.row)">{{ $t('lang.terminate_card') }}</el-dropdown-item>
                  <el-dropdown-item @click="handleRefresh(scope.row)">{{ $t('lang.refresh') }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <div class="action-buttons d-none d-sm-flex display-flex"> 
              <el-button v-if="scope.row.cardLevel !== 'SubCard'" type="primary" size="small" @click="handleRechargeCard(scope.row)">{{ $t('lang.top_up') }}</el-button>
              <el-button v-if="scope.row.cardLevel === 'SubCard'" type="primary" size="small" @click="handleAdjustLimit(scope.row)">{{ $t('lang.adjust_limit') }}</el-button>
              <el-button v-if="scope.row.cardLevel !== 'SubCard'" type="secondary" size="small" @click="handleWithdrawCard(scope.row)">{{ $t('lang.withdraw_card') }}</el-button>
              <el-button type="danger" size="small" @click="handleCancelCard(scope.row)">{{ $t('lang.terminate') }}</el-button>
              <el-button :loading="scope.row.loading" type="success" size="small" @click="handleRefresh(scope.row)">{{ $t('lang.refresh') }}</el-button>
            </div>
          </template>
        </el-table-column> -->
      </el-table>
            <!-- 分页 -->
            <el-pagination
              background
              layout="total, sizes, prev, pager, next"
              :total="search2.total"
              :page-size="search2.pageSize"
              :page-sizes="[10, 50, 100]"
              :pager-count="5"
              :current-page="search2.page"
              @size-change="handleSizeChange2"
              @current-change="handleCurrentChange2"
            ></el-pagination>
    </el-tab-pane>
  </el-tabs>
  
    <!-- 开卡对话框 -->
    <el-dialog :title="$t('lang.open_card')" v-model="dialogs.activeCardDialogVisible" width="50%" align-center>
        <el-form :model="activeCardForm" label-width="140px">
            <div class="row">
            <div class="col">
            
            <el-form-item :label="$t('lang.card_type')">
              <el-radio-group v-model="filters.cardModel" @change="filterResult"> 
                <el-radio v-for="item in filterParams.cardModel" :value="item">{{ getCardModelLabel(item) }}</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item :label="$t('lang.card_currency')">
              <el-radio-group v-model="filters.currency" @change="filterResult"> 
                <el-radio v-for="item in filterParams.currency" :value="item">{{item}}</el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- <el-form-item :label="$t('lang.card_brand')">
              <el-radio-group v-model="filters.brand" @change="filterResult"> 
                <el-radio v-for="item in filterParams.brand" :value="item">{{item}}</el-radio>
              </el-radio-group>
            </el-form-item> -->


            <el-form-item :label="$t('lang.card_brand')">
              <el-radio-group v-model="filters.brand" @change="filterResult">
                <el-radio-button :class="`card-bg card-${item}`" v-for="item in filterParams.brand" :key="item" :value="item">
                </el-radio-button>
              </el-radio-group>
            </el-form-item>
            
              <el-form-item :label="$t('lang.card_bin')">
                <el-radio-group v-model="cardForm.card"> 
                  <el-radio v-for="item in cardList" :value="item"><el-tag type="success">{{item.cardBin}}</el-tag></el-radio>
                </el-radio-group>
              </el-form-item>
            
            <div v-if="cardForm.card">
              <el-form-item>
                <el-input type="textarea" autosize disabled v-model="cardForm.card.supportPlatform" class="bg-body-secondary text-secondary note-bg">
                </el-input>
              </el-form-item>
              
              <el-form-item :label="$t('lang.wallet_balance')">
                <p>{{ userStore.userInfo.wallet == null ? '0' : userStore.userInfo.wallet.balance }} USD</p>
              </el-form-item>
              <template v-if="cardForm.card && cardForm.card.cardModel === 'SHARE'">
                <el-form-item :label="$t('lang.share_card')">
                  <el-radio-group v-model="cardForm.shareCardType" @change="handleShareCardTypeChange">
                    <el-radio value="MasterCard">{{ $t('lang.master_card') }}</el-radio>
                    <el-radio value="SubCard">{{ $t('lang.sub_card') }}</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="cardForm.shareCardType === 'SubCard'" :label="$t('lang.primary_card')">
                  <el-select v-model="cardForm.primaryCardId" clearable :placeholder="$t('lang.select_primary_card')">
                    <el-option v-for="item in filteredMasterCards" :key="item.ID" :value="item.cardId" :label="item.cardNo"></el-option>
                  </el-select>
                </el-form-item>
                <template v-if="cardForm.shareCardType === 'SubCard' && cardForm.primaryCardId">
                  <el-form-item :label="$t('lang.auth_limit_flag')">
                    <el-checkbox :model-value="cardForm.authLimitFlag === 'N'" @change="(val) => cardForm.authLimitFlag = val ? 'N' : 'Y'">{{ $t('lang.no_limit') }}</el-checkbox>
                  </el-form-item>
                  <el-form-item :label="$t('lang.total_auth_limit')">
                    <el-input-number v-model="cardForm.totalAuthLimit" :min="0" :controls="false" :placeholder="$t('lang.please_enter_total_auth_limit')" :disabled="cardForm.authLimitFlag === 'N'"></el-input-number>
                  </el-form-item>
                </template>
              </template>
              <el-form-item v-if="!(cardForm.card && cardForm.card.cardModel === 'SHARE' && cardForm.shareCardType === 'SubCard' && cardForm.primaryCardId)" :label="$t('lang.recharge_amount')" required>
                  <el-input-number v-model="cardForm.amount" :min="cardForm.card.createRechargeLimit" :placeholder="`>${cardForm.card.createRechargeLimit} ${cardForm.card.currency}`" :controls="false">
                  </el-input-number>
              </el-form-item>
              <el-form-item :label="$t('lang.notes')">
                <el-input v-model="cardForm.remark" type="textarea" :rows="2" :placeholder="$t('lang.please_enter_notes')"></el-input>
              </el-form-item>
              <el-form-item :label="$t('lang.open_card_number')">
                <el-input-number v-model="cardForm.number" :max="10" :min="1" controls-position="right" :placeholder="$t('lang.open_card_number')"></el-input-number>
              </el-form-item>
              <el-form-item :label="$t('lang.card_group_name')">
                <el-select v-model="cardForm.groupId" clearable>
                    <el-option v-for="item in cardGroups" :label="item.name" :value="item.ID"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item v-if="!(cardForm.card && cardForm.card.cardModel === 'SHARE' && cardForm.shareCardType === 'MasterCard')" :label="$t('lang.card_holder')">
                <el-select v-model="cardForm.cardHolderId">
                  <el-option v-for="item in holders" :value="item.cardHolderId" :label="`${item.firstName} ${item.lastName}`"></el-option>
                </el-select>
                <el-link @click="dialogs.addCardHolderDialogVisible = true;cardHolder={mobilePrefix: '+1'}" icon="plus" style="color:#01ad5a; margin-left: 5px; ">{{ $t('lang.add_new_cardholder') }}</el-link>
              </el-form-item>
            </div>
            </div>
        </div>
        </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogs.activeCardDialogVisible = false">{{ $t('lang.cancel') }}</el-button>
          <el-button 
            :disabled="!cardForm.card"  
            :loading="loading.addCardLoading" 
            type="primary" 
            @click="handleActiveCardConfirm">
            {{ $t('lang.create') }}
          </el-button>
        </span>
      </template>
    </el-dialog>

  <!-- 添加持卡人弹窗 -->
  <el-dialog 
    v-model="dialogs.addCardHolderDialogVisible" 
    :title="$t('lang.add_cardholder')" 
    width="50%" 
    align-center>
    <el-form  label-width="auto">
      <el-form-item :label="$t('lang.nationality') + ' *'">
        <el-select v-model="cardHolder.region" :placeholder="$t('lang.please_select_nationality')" clearable style="min-width: 180px;">
            <!-- <el-option :label="$t('lang.china')" value="CHN"></el-option> -->
            <el-option :label="$t('lang.united_states')" value="USA"></el-option>
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
    
    <el-form-item :label="$t('lang.cardholder_first_name') + ' *'">
      <el-input v-model="cardHolder.lastName" :placeholder="$t('lang.please_enter_cardholder_first_name')"></el-input>
    </el-form-item>

      <el-form-item :label="$t('lang.cardholder_last_name') + ' *'">
        <el-input v-model="cardHolder.firstName" :placeholder="$t('lang.please_enter_cardholder_last_name')"></el-input>
      </el-form-item>

    <!-- 账单地国家 -->
    <el-form-item :label="$t('lang.billing_country') + ' *'">
      <el-select v-model="cardHolder.countryCode" :placeholder="$t('lang.please_select_country')" clearable style="min-width: 200px;">
          <el-option :label="$t('lang.united_states')" value="USA"></el-option>
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

  <!-- 充值弹窗 -->
  <el-dialog
    :title="$t('lang.card_recharge')"
    width="50%"
    v-model="dialogs.rechargeCardDialogVisible"
    style="max-width: 500px;"
    align-center
  >
    <el-form label-position="top" class="p-4">
      <el-form-item :label="$t('lang.card_number')">
        <el-input v-model="currCard.cardNo" disabled>
          <template #prepend>
            <span style="color: #1e73be;padding:0 10px; font-weight: bold;">{{ currCard.cardBrand }}</span>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.card_currency')">
        <el-input v-model="currCard.currency" disabled></el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.available_balance')">
        <el-input v-model="userStore.userInfo.wallet.balance" disabled></el-input>
      </el-form-item>

      <el-form-item 
        :label="$t('lang.recharge_amount')" 
        required 
        :error="$t('lang.please_enter_amount_min', { min: currCard.belongCardbin.minBalance })" 
        :validate-status="currCard.rechargeAmount >= currCard.belongCardbin.minBalance ? 'success' : 'error'">
        <el-input v-model="currCard.rechargeAmount" :placeholder="$t('lang.please_enter_amount')">
          <template #append>USD</template>
        </el-input>
        <PreRechargeSummary :summary="preRechargeSummary" :loading="preRechargeLoading" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="d-flex justify-content-end">
        <el-button 
          :loading="loading.rechargeCardLoading" 
          type="info" 
          @click="dialogs.rechargeCardDialogVisible = false">
          {{ $t('lang.cancel') }}
        </el-button>
        <el-button 
          :loading="loading.rechargeCardLoading"
          :disabled="preRechargeLoading"
          type="primary" 
          @click="handleRechargeCardConfirm">
          {{ $t('lang.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 调整额度弹窗 -->
  <el-dialog
    :title="$t('lang.adjust_limit')"
    width="50%"
    v-model="dialogs.adjustLimitDialogVisible"
    style="max-width: 500px;"
    align-center
    destroy-on-close
  >
    <el-form label-position="top" class="p-4">
      <el-form-item :label="$t('lang.card_number')">
        <el-input v-model="currCard.cardNo" disabled>
          <template #prepend>
            <span style="color: #1e73be;padding:0 10px; font-weight: bold;">{{ currCard.cardBrand }}</span>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.card_currency')">
        <el-input v-model="currCard.currency" disabled></el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.current_limit')">
        <el-input :model-value="formatCurrentAuthLimitDisplay(currCard)" disabled>
          <template v-if="currCard && showCurrentAuthLimitCurrencyAppend(currCard)" #append>{{ currCard.currency }}</template>
        </el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.total_auth_limit')">
        <div class="d-flex align-items-center w-100" style="gap: 10px; flex-wrap: nowrap">
          <el-input-number
            v-model="currCard.newLimit"
            :min="0"
            :controls="false"
            :placeholder="$t('lang.please_enter_total_auth_limit')"
            :disabled="currCard.authLimitFlag === 'N'"
            class="flex-grow-1"
            style="min-width: 0"
          ></el-input-number>
          <el-checkbox
            :model-value="currCard.authLimitFlag === 'N'"
            class="flex-shrink-0 m-0"
            @change="onAdjustAuthLimitFlagChange"
          >{{ $t('lang.no_limit') }}</el-checkbox>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="d-flex justify-content-end">
        <el-button 
          :loading="loading.adjustLimitLoading" 
          type="info" 
          @click="dialogs.adjustLimitDialogVisible = false">
          {{ $t('lang.cancel') }}
        </el-button>
        <el-button 
          :loading="loading.adjustLimitLoading" 
          type="primary" 
          @click="handleAdjustLimitConfirm">
          {{ $t('lang.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 退款弹窗 -->
  <el-dialog
    :title="$t('lang.card_withdrawal')"
    width="50%"
    v-model="dialogs.withdrawCardDialogVisible"
    style="max-width: 500px;"
    align-center
  >
    <el-form label-position="top" class="p-4">
      <el-form-item :label="$t('lang.card_number')">
        <el-input v-model="currCard.cardNo" disabled>
          <template #prepend>
            <span style="color: #1e73be;padding:0 10px; font-weight: bold;">{{ currCard.cardBrand }}</span>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.card_currency')">
        <el-input v-model="currCard.currency" disabled></el-input>
      </el-form-item>

      <el-form-item :label="$t('lang.available_balance')">
        <el-input v-model="currCard.balance" disabled>
          <template #append>{{ currCard.currency }}</template>
        </el-input>
      </el-form-item>

      <el-form-item 
        :label="$t('lang.withdrawal_amount')" 
        required 
        :error="$t('lang.please_enter_amount_max', { max: currCard.balance })" 
        :validate-status="currCard.withdrawAmount <= currCard.balance ? 'success' : 'error'">
        <el-input v-model="currCard.withdrawAmount" :placeholder="$t('lang.please_enter_amount')">
          <template #append>USD</template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="d-flex justify-content-end">
        <el-button 
          :loading="loading.withdrawCardLoading" 
          type="info" 
          @click="dialogs.withdrawCardDialogVisible = false">
          {{ $t('lang.cancel') }}
        </el-button>
        <el-button 
          :loading="loading.withdrawCardLoading" 
          type="primary" 
          @click.prevent="handleWithdrawCardConfirm">
          {{ $t('lang.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 卡详情弹窗 -->
  <el-dialog
    v-model="dialogs.cardDetailDialogVisible"
    width="80%"
    destroy-on-close
    align-center
    class="card-detail-dialog"
  >
    <CardDetail :card="cardView" />
  </el-dialog>
  <el-dialog v-model="dialogs.setCardGroupDialogVisible" :title="$t('lang.set_card_group')" width="30%">
    <el-select v-model="groupForm.groupId" :placeholder="$t('lang.please_select')" style="width: 100%;">
      <el-option
        v-for="item in cardGroups"
        :key="item.ID"
        :label="item.name"
        :value="item.ID"
      />
    </el-select>
    <template #footer>
      <el-button @click="dialogs.setCardGroupDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSetGroupConfirm">确定</el-button>
    </template>
  </el-dialog>
</div>
  </template>
  
  <script setup>
  import { reactive, ref,onMounted, computed,h, watch } from 'vue';
  import { ElMessage,ElMessageBox,ElSelect,ElOption } from 'element-plus';
  import { formatDate,addYear} from '@/utils/format';
  import { setCardGroup,listCardGroup,listCardHolder,remarkCard,syncCard,addCardHolder,listCardBin,createCard,cancelCard,listCard,rechargeCard,withdrawCard,adjustSubCardLimit } from '@/api/finance';
  import { getIamUserList } from '@/api/iam';
  import { useUserStore } from '@/pinia/modules/user'
  import CardDetail from './cardDetail.vue'
  import { randomEmailAndName, randomUsPhoneForBilling, randomVerifiedUsBillingAddress, randomBirth } from '@/utils/random'
  import {buildExcel} from '@/utils/excel'
  import { buildCancelListPayload, applyCancelCardResult } from '@/utils/cancelCard'
  import { useCardPreRecharge } from '@/composables/useCardPreRecharge'
  import PreRechargeSummary from '@/components/card/PreRechargeSummary.vue'
  import { useI18n } from 'vue-i18n';
  const isMobile = computed(() => window.innerWidth < 768)
  const { t, locale } = useI18n();
  const rawTotalAuthLimitR = (c) => c?.totalAuthLimit ?? c?.total_auth_limit
  const isNoLimitTotalAuthR = (c) => {
    const raw = rawTotalAuthLimitR(c)
    if (raw == null || raw === '') return false
    return Number(raw) === 0
  }
  const isUnlimitedAuthLimitRowR = (c) => c && (c.authLimitFlag === 'N' || isNoLimitTotalAuthR(c))
  const formatCurrentAuthLimitDisplay = (c) => {
    if (!c) return ''
    if (isUnlimitedAuthLimitRowR(c)) return t('lang.no_limit_balance')
    return c.totalAuthLimit
  }
  const showCurrentAuthLimitCurrencyAppend = (c) => c && !isUnlimitedAuthLimitRowR(c)
  const userStore = useUserStore()

  const dialogs =reactive({
    activeCardDialogVisible:false,
    addCardHolderDialogVisible:false,
    rechargeCardDialogVisible:false,
    withdrawCardDialogVisible:false,
    adjustLimitDialogVisible:false,
    cardDetailDialogVisible:false,
    remarkCardDialogVisible:false,
    setCardGroupDialogVisible:false,
  })

  const currCard = ref()

  const { summary: preRechargeSummary, loading: preRechargeLoading, reset: resetPreRecharge } = useCardPreRecharge({
    enabled: () => dialogs.rechargeCardDialogVisible,
    getCardId: () => currCard.value?.cardId,
    getRechargeAmount: () => currCard.value?.rechargeAmount,
  })

  const loading = reactive({
    addCardHolderLoading:false,
    addCardLoading:false,
    rechargeCardLoading:false,
    withdrawCardLoading:false,
    adjustLimitLoading:false,
    syncCardLoading:false,
  })
  const activeTab = ref('active')
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
    cardNumber: '',
    clientId: 0,
    dateRange: [],
    cardStatus: '',
    remark: '',
    minBalance: 0,
    maxBalance: 0,
    iamUserId: '',
    cardModel: 'CARD',
    total: 0,
    pageSize: 10,
    page:1
  })
  const search2 = reactive({
    cardNumber: '',
    clientId: 0,
    dateRange: [],
    cardStatus: 'Closed',
    remark: '',
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
    cardBin:'',
    number: 1,
    shareCardType: '',
    primaryCardId: '',
    totalAuthLimit: '',
    authLimitFlag: 'Y'
  })

  const allMasterCards = ref([])
  const handleListAllMasterCards = () => {
    listCard({
      cardLevel: 'MasterCard',
      cardStatus: 'Active',
      cardModel: 'SHARE',
      pageSize: 10000,
      page: 1
    }).then((res) => {
      if (res.code === 0) {
        allMasterCards.value = res.data.list || []
      }
    })
  }
  // 开卡子卡主卡列表：在 allMasterCards 上按当前卡 BIN 过滤
  const filteredMasterCards = computed(() => {
    if (!cardForm.value.card || !cardForm.value.card.cardBin) {
      return allMasterCards.value
    }
    return allMasterCards.value.filter((card) => {
      return card.cardBin === cardForm.value.card.cardBin ||
        card.belongCardbin?.cardBin === cardForm.value.card.cardBin
    })
  })
  const handleListCard = () => {
    listCard(search).then(res => { 
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }
  const cardGroups = ref([])
  const handleListCardGroup = () => {
    listCardGroup(holderSearch).then(res => { 
      if (res.code === 0){
        cardGroups.value = res.data.list
      }
    })
  }
  const iamUsers = ref([])
  const handleListIamUsers = () => {
    getIamUserList({ page: 1, pageSize: 100 }).then(res => {
      if (res.code === 0) {
        iamUsers.value = res.data.list
      }
    })
  }
  const getTableData = ()=>{
    listCard(search).then(res => { 
      if (res.code === 0){
        tableData.value = res.data.list
        search.total = res.data.total
      }
    })
  }
  const tableData2 = ref([])
  const getTableData2 = ()=>{
    listCard(search2).then(res => { 
      if (res.code === 0){
        tableData2.value = res.data.list
        search2.total = res.data.total
      }
    })
  }
  onMounted(() => {
    handleListCardBin()
    handleListCardHolder()
    handleListCardGroup()
    handleListIamUsers()
    handleListAllMasterCards()
    handleListCard()
  })
  // 监听卡BIN变化，清空主卡选择
  watch(() => cardForm.value.card?.cardBin, (newCardBin, oldCardBin) => {
    if (newCardBin && newCardBin !== oldCardBin) {
      cardForm.value.primaryCardId = ''
    }
  })
  const handleTabChange = ()=>{
    if(activeTab.value === 'active'){
      getTableData()
    }else {
      getTableData2()
    }
    
  }
  const randomCardHolder = () => {
    const emailAndName = randomEmailAndName()
    const billing = randomVerifiedUsBillingAddress()
    cardHolder.value = {
      region: 'USA',
      countryCode: 'USA',
      firstName: emailAndName[0],
      lastName: emailAndName[1],
      email: emailAndName[2],
      mobilePrefix: '+1',
      mobile: randomUsPhoneForBilling(billing),
      birthDate: randomBirth(),
      state: billing.state,
      city: billing.city,
      postcode: billing.postcode,
      address: billing.address
    }
  }
  const handleRefresh = (row) => {
    row.loading = true
    loading.syncCardLoading = true
    syncCard({id:row.ID}).then(res => { 
      if (res.code ===0){
        handleListCard()
      } 
      loading.syncCardLoading = false
      row.loading = false
    }) 
  }
  const cardFormRef = ref(null)
  const handleListCardHolder = () => {
    listCardHolder(holderSearch).then((res) => {
      holders.value = res.data.list;
    })
  };
  const cardBins = ref([])
  const handleListCardBin = () => {
    listCardBin({ ...holderSearch, cardModel: 'CARD' }).then((res) => {
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
    cardForm.value.card = {}
    // 当卡BIN改变时，清空已选择的主卡
    cardForm.value.primaryCardId = ''
  }
  const getCardTypeLabel = (cardLevel) => {
    if (cardLevel === 'SubCard') {
      return t('lang.sub_card')
    } else if (cardLevel === 'MasterCard') {
      return t('lang.master_card')
    } else {
      return t('lang.recharge_card')
    }
  }
  const getCardModelLabel = (cardModel) => {
    if (cardModel === 'SHARE') {
      return t('lang.share_card')
    } else if (cardModel === 'CARD') {
      return t('lang.recharge_card')
    } else {
      return cardModel
    }
  }
  
  const activeCardForm = ref({
    
  })
  const cardHolder = ref({
    mobilePrefix: '+1',
  })
  
  const handleAddCardHoderConfirm = ()=>{
    loading.addCardHolderLoading = true
    addCardHolder(cardHolder.value).then(res=>{
      loading.addCardHolderLoading = false
      if(res.code === 0){
        ElMessage.success('Success')
        dialogs.addCardHolderDialogVisible = false
        handleListCardHolder()
      
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
      cardModel: cardForm.value.card.cardModel,
      amount: cardForm.value.amount,
      number:cardForm.value.number,
      remark:cardForm.value.remark,
      groupId:cardForm.value.groupId
    }
    // 如果选择的是共享卡主卡，不需要传递cardHolderId
    if (!(cardForm.value.card && cardForm.value.card.cardModel === 'SHARE' && cardForm.value.shareCardType === 'MasterCard')) {
      form.cardHolderId = cardForm.value.cardHolderId
    }
    // 如果选择了共享卡，并且选择了子卡，添加相关字段
    if (cardForm.value.card && cardForm.value.card.cardModel === 'SHARE' && cardForm.value.shareCardType === 'SubCard' && cardForm.value.primaryCardId) {
      form.primaryCardId = cardForm.value.primaryCardId
      form.totalAuthLimit = cardForm.value.totalAuthLimit
      form.authLimitFlag = cardForm.value.authLimitFlag
    }
    if(form.cardBin == null || form.cardBinId == null){
      ElMessage.error('Please select card bin')
      loading.addCardLoading = false
      return
    }
    createCard(form).then(res =>{
      if(res.code === 0){
        ElMessage.success('Create Card Success')
        dialogs.activeCardDialogVisible = false
        getTableData()
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
  
  const handleShareCardTypeChange = (value) => {
    // 当选择类型改变时，清空主卡相关字段
    if (value === 'MasterCard') {
      cardForm.value.primaryCardId = ''
      cardForm.value.totalAuthLimit = ''
      cardForm.value.authLimitFlag = 'Y'
      cardForm.value.cardHolderId = '' // 选择主卡时清空持卡人
    }
  }
  const handleActiveCard = () => {
    // 重置主卡相关字段
    cardForm.value.shareCardType = ''
    cardForm.value.primaryCardId = ''
    cardForm.value.totalAuthLimit = ''
    cardForm.value.authLimitFlag = 'Y'
    // 查找默认的 cardBin 并设置为默认选中
    const defaultCard = cardBins.value.find(item => item.isDefault)
    if (defaultCard) {
      // 设置筛选条件为默认卡的属性
      filters.cardModel = defaultCard.cardModel
      filters.currency = defaultCard.currency
      filters.brand = defaultCard.cardBrand
      // 根据筛选条件过滤 cardList
      filterResult()
      // 设置默认选中的卡
      cardForm.value.card = defaultCard
    }
    handleListAllMasterCards()
    dialogs.activeCardDialogVisible = true
  }
  const batchFreezeCards = () => {
    // 处理批量冻结逻辑
    console.log('点击了批量冻结按钮')
  };
  const handleRechargeCardConfirm = ()=>{
    if (!preRechargeSummary.value?.quotationRequestId) {
      ElMessage.warning(t('lang.pre_recharge_quote_required'))
      return
    }
    loading.rechargeCardLoading = true
    rechargeCard({
      id: currCard.value.ID,
      currency: userStore.userInfo.wallet.currency,
      amount: currCard.value.rechargeAmount,
      quotationRequestId: preRechargeSummary.value.quotationRequestId,
    }).then(res => {
      loading.rechargeCardLoading = false
      if(res.code === 0){
        ElMessage.success('Rechage success')
        dialogs.rechargeCardDialogVisible = false
        resetPreRecharge()
        getTableData()
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
        ElMessage.success('withdraw success')
        getTableData()
      }
      userStore.GetUserInfo()
    })
  }
  const handleCancelCard = (row)=>{
    ElMessageBox.confirm(
    t('lang.cancel_card_warning'),
    t('lang.warning'),
    {
      confirmButtonText: t('lang.terminate_card'),
      cancelButtonText: t('lang.cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      handleCancelConfirm(row)
    })
    .catch(() => {
    })
  }
  const cardListTableRef = ref()
  const selectedCancelRows = ref([])

  const onCancelSelectionChange = (rows) => {
    selectedCancelRows.value = rows
  }

  const runCancelWithPayload = (payload) => {
    cancelCard(payload).then((res) => {
      applyCancelCardResult(res, {
        t,
        ElMessage,
        onAfter: (ok) => {
          if (ok) {
            getTableData()
            cardListTableRef.value?.clearSelection()
          }
        }
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
    runCancelWithPayload(payload)
  }

  const handleBatchCancelCard = () => {
    const rows = selectedCancelRows.value
    if (!rows?.length) {
      ElMessage.warning(t('lang.cancel_batch_empty'))
      return
    }
    let payload
    try {
      payload = buildCancelListPayload(rows)
    } catch (e) {
      if (e.message === 'cancel_list_too_many') ElMessage.warning(t('lang.cancel_list_too_many'))
      else ElMessage.warning(t('lang.cancel_list_invalid'))
      return
    }
    ElMessageBox.confirm(
      t('lang.cancel_card_warning'),
      t('lang.warning'),
      {
        confirmButtonText: t('lang.terminate_card'),
        cancelButtonText: t('lang.cancel'),
        type: 'warning',
      }
    )
      .then(() => {
        runCancelWithPayload(payload)
      })
      .catch(() => {})
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

const handleSizeChange2 = (val) => {
    search2.pageSize = val
    getTableData2()
  }
  
const handleCurrentChange2 = (val) => {
  search2.page = val
  getTableData2()
};

const cardDetail = ref()
const handleRechargeCard = (row) => {
  if (row.cardLevel === 'SubCard') {
    ElMessage.warning(t('lang.sub_card_cannot_recharge'))
    return
  }
  row.belongCardbin = cardBins.value.find(item => item.cardBinId === row.cardBinId)
  currCard.value = row
  resetPreRecharge()
  dialogs.rechargeCardDialogVisible = true
};

const handleWithdrawCard = (row) => {
  if (row.cardLevel === 'SubCard') {
    ElMessage.warning(t('lang.sub_card_cannot_withdraw'))
    return
  }
  currCard.value = row
  dialogs.withdrawCardDialogVisible = true

};

const handleAdjustLimit = (row) => {
  const unl = isUnlimitedAuthLimitRowR(row)
  const n = row.totalAuthLimit != null && row.totalAuthLimit !== '' ? Number(row.totalAuthLimit) : 0
  currCard.value = { ...row }
  currCard.value.authLimitFlag = unl ? 'N' : 'Y'
  currCard.value.newLimit = unl ? 0 : n
  currCard.value._revertNewLimit = unl ? 0 : n
  dialogs.adjustLimitDialogVisible = true
}

const onAdjustAuthLimitFlagChange = (checked) => {
  const c = currCard.value
  if (!c) return
  c.authLimitFlag = checked ? 'N' : 'Y'
  if (checked) {
    c._revertNewLimit = c.newLimit
    c.newLimit = 0
  } else {
    const r = Number(c._revertNewLimit)
    c.newLimit = Number.isFinite(r) && r > 0 ? r : 1
  }
}

const handleAdjustLimitConfirm = () => {
  if (currCard.value.authLimitFlag !== 'N') {
    if (currCard.value.newLimit == null || Number(currCard.value.newLimit) <= 0) {
      ElMessage.warning(t('lang.please_enter_new_limit'))
      return
    }
  }
  const unlimited = currCard.value.authLimitFlag === 'N'
  loading.adjustLimitLoading = true
  adjustSubCardLimit({
    id: currCard.value.ID,
    totalAuthLimit: unlimited ? 0 : currCard.value.newLimit,
    authLimitFlag: unlimited ? 'N' : 'Y'
  }).then(res => {
    loading.adjustLimitLoading = false
    if(res.code === 0){
      ElMessage.success(t('lang.adjust_limit_success'))
      dialogs.adjustLimitDialogVisible = false
      getTableData()
    }
  })
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
            "Status":x.cardStatus,
            "Balance":parseFloat(x.balance),
            "Currency":x.currency,

            "Creation Time":x.activeDate,
            "Brand":x.cardBrand,
            "Notes":x.remark,
            "Group Name":x.group?.name
          }
        })
        buildExcel(list,"card_list")
      }
    })
  }
  const handleEditRemark = (row) =>{
    ElMessageBox.prompt(
    '',
    t('lang.notes'),
    {
      confirmButtonText: t('lang.confirm'),
      cancelButtonText: t('lang.cancel'),
      inputType: 'textarea',
      inputValue: row.remark,
    }
  )
    .then(({ value }) => {
      remarkCard({
        id:row.ID,
        remark:value
      }).then(res=>{
        if(res.code === 0){
          handleListCard()
        }
      })
    })
  }
  const groupForm = ref({})
  const handleEditGroup = (row) =>{
    dialogs.setCardGroupDialogVisible = true
    groupForm.value.id = row.ID
    groupForm.value.groupId = row.groupId==0?null:row.groupId
    
  }
const handleSetGroupConfirm = () =>{
  setCardGroup(groupForm.value).then(res =>{
    if(res.code === 0){
      ElMessage.success("Success")
      handleListCard()
      dialogs.setCardGroupDialogVisible = false
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

.el-input-number .el-input__inner {
  text-align: left ;
}

.note-bg  {
  border-radius: 16px;
}
:deep(.card-detail-dialog .el-dialog__header) {
  border-bottom: none !important;
}

.action-buttons .el-button {
  margin-left: 5px;
}

</style>