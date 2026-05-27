<template>
    <div>
        <el-tabs>
            <el-tab-pane :label="$t('lang.wallet_recharge')">
                <RechargeManager/>
            </el-tab-pane>
            <el-tab-pane :label="$t('lang.wallet_withdrawal')">
                <WithdrawManager/>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
<script setup>
import WithdrawManager from '@/view/pixal/wallet/withdrawManager/index.vue'
import RechargeManager from '@/view/pixal/wallet/rechargeManager/index.vue'

</script>
