<template>
  <el-sub-menu ref="subMenu" :index="routerInfo.name">
    <template #title>
      <div>
        <i v-if="routerInfo.meta.bootstrapIcon" class="item-label-icon" style="height: 0.9em;width: 0.9em;margin-bottom:8px;margin-right: 15px;" :class="buildIcon(routerInfo.meta.bootstrapIcon)"></i>
        <el-icon v-else="routerInfo.meta.icon">
          <component :is="routerInfo.meta.icon" />
        </el-icon>
        <span>{{ $t('lang.'+routerInfo.meta.title) }}</span>
      </div>
    </template>
    <slot />
  </el-sub-menu>
</template>

<script>
export default {
  name: 'AsyncSubmenu',
}
</script>

<script setup>
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
</script>

<style lang="scss" scoped>
.bootstrap-icon {
  display: inline-flex;
  justify-content: center;
  font-size: 18px;
  vertical-align: middle;
}
</style>
