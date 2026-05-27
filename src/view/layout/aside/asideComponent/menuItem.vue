<template>
  <el-menu-item :index="routerInfo.name">
    <template #title>
      <div v-if="routerInfo.meta.title == 'Logout'"  class="item-label d-flex align-items-center" style="color:#7a7186">
        <el-popconfirm
          title="Confirm Log Out?"
          icon=""
          placement="bottom-start"
          @confirm="loginout"
        >
        <template #reference>
            <div>
              <i v-if="routerInfo.meta.bootstrapIcon" style="height: 0.9em;width: 0.9em;margin-bottom:8px;margin-right: 15px;" :class="buildIcon(routerInfo.meta.bootstrapIcon)"></i>
              <el-icon v-else-if="routerInfo.meta.icon">
                <component :is="routerInfo.meta.icon" />
              </el-icon>
              <span class="gva-menu-item-title">{{ $t(routerInfo.meta.title) }}</span>
            </div>
        </template>
        <template #actions="{confirm,cancel}">
          <div class="gap-1 justify-items-center d-flex">
            <button class="btn btn-sm w-75 btn-outline-secondary" @click="cancel">No</button>
            <button class="btn btn-sm w-75 btn-danger" @click="confirm">Yes</button>
          </div>
        </template>
    </el-popconfirm>
    </div>
    <div v-else  class="item-label d-flex align-items-center" >
      <!-- <i v-if="routerInfo.meta.bootstrapIcon" style="height: 0.9em;width: 0.9em;margin-bottom:8px;margin-right: 15px;" :class="buildIcon(routerInfo.meta.bootstrapIcon)"></i> -->
      <!-- <FontAwesomeIcon v-if="routerInfo.meta.bootstrapIcon" style="height: 0.9em;width: 0.9em;margin-bottom:8px;margin-right: 15px;" :icon="fa-solid fa-check-square" /> -->
      <i v-if="routerInfo.meta.bootstrapIcon" class="item-label-icon" style="height: 0.9em;width: 0.9em;margin-bottom:8px;margin-right: 15px;" :class="buildIcon(routerInfo.meta.bootstrapIcon)"></i>
      <el-icon v-else-if="routerInfo.meta.icon">
        <component :is="routerInfo.meta.icon" />
      </el-icon>
      
      <span class="gva-menu-item-title">{{ $t('lang.'+routerInfo.meta.title) }}</span>
    </div>
    </template>
  </el-menu-item>
</template>

<script>
export default {
  name: 'MenuItem',
}
</script>


<script setup>
import { useUserStore } from '@/pinia/modules/user'
const userStore = useUserStore()

defineProps({
  routerInfo: {
    default: function() {
      return null
    },
    type: Object
  }
})
const buildIcon = (icon) => {
  return "bootstrap-icon "+ icon
}
const loginout = ()=>{
    userStore.LoginOut()
}
</script>

<style lang="scss" scoped>
.gva-menu-item-title {
  max-width: 160px;
}
.bootstrap-icon{
    margin: 5px;
    // width: 2px;
    // text-align: center;
    font-size: 18px;
    vertical-align: middle;
    // align-items: center;
    display: inline-flex;
    justify-content: end;
}
</style>
