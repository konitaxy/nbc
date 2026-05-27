<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style lang="scss">
@use '@/style/main.scss';
@use '@/style/base.scss';
@use '@/style/mobile.scss';
#app {
  background: #eee;
  height: 100vh;
  // overflow: hidden;
  font-weight: 400 !important;
}
.el-button{
  font-weight: 400 !important;
}
</style>
