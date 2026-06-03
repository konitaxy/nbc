<template>
  <div class="console-sidebar-inner text-center">
    <div class="console-sidebar-scroll">
      <transition
        :duration="{ enter: 800, leave: 100 }"
        mode="out-in"
        name="el-fade-in-linear"
      >
        <el-menu
          :collapse="isCollapse"
          :collapse-transition="false"
          :default-active="active"
          :active-text-color="userStore.activeColor"
          class="el-menu-vertical"
          unique-opened
          @select="selectMenuItem"
        >
          <template v-for="item in routerStore.asyncRouters[0].children">
            <aside-component
              v-if="!item.hidden"
              :key="item.name"
              :router-info="item"
            />
          </template>
        </el-menu>
      </transition>
    </div>
  
  </div>
</template>

<script>
export default {
  name: 'Aside',
}
</script>

<script setup>
import AsideComponent from '@/view/layout/aside/asideComponent/index.vue'
import { emitter } from '@/utils/bus.js'
import { ref, watch, onUnmounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'

const route = useRoute()
const router = useRouter()

const userStore = useUserStore()
const routerStore = useRouterStore()

const active = ref('')

const showDrawer = reactive({
  loginout:false,
  dialog:false
})
watch(route, () => {
  active.value = route.name
})
const isCollapse = ref(false)
const isMobile = ref(false)
const initPage = () => {
  active.value = route.name
  const screenWidth = document.body.clientWidth
  if (screenWidth < 992) {
    isCollapse.value = !isCollapse.value
  }

  emitter.on('collapse', (item) => {
    isCollapse.value = item
  })
  emitter.on('mobile', (item) => {
    isMobile.value = item
  })
}

initPage()

// onUnmounted(() => {
//   emitter.off('collapse')
// })

const selectMenuItem = (index, _, ele, aaa) => {
  const query = {}
  const params = {}
 routerStore.routeMap[index]?.parameters &&
    routerStore.routeMap[index]?.parameters.forEach((item) => {
      if (item.type === 'query') {
        query[item.key] = item.value
      } else {
        params[item.key] = item.value
      }
    })
 if (index === route.name) return
 if (index.indexOf('http://') > -1 || index.indexOf('https://') > -1) {
   window.open(index)
 } else if(index.indexOf('loginout')> -1){
  if (isMobile.value){
    showDrawer.loginout = true
  }else {
    showDrawer.dialog = true
  }
  
  // loginout()
 }else {
   router.push({ name: index, query, params })
  if(isMobile.value){
    //  isCollapse.value = true
     emitter.emit('closeExpland')
  }
 }
}
const loginout = ()=>{
    userStore.LoginOut()
}
const cancelLoginout = ()=>{
  document.querySelector('.el-popover').remove()
}
</script>

<style lang="scss">
.el-sub-menu__title,
.el-menu-item {
  .item-label {
      transition-property:  color;
      transition-timing-function: ease-in-out;
      transition-duration: 0.3s;
      
    } 
}

.el-scrollbar {
  .el-scrollbar__view {
    height: 100%;
  }
}
.menu-info {
  .menu-contorl {
    line-height: 27px;
    font-size: 20px;
    display: table-cell;
    vertical-align: middle;
  }
}

.el-menu-item.is-active{
    color: var(--el-color-primary) ;
    // padding-left:100px;
    // margin-top:7px;
    // margin-bottom:7px;
    .item-label {
      color: var(--el-color-primary) ;
      transition-property: color;
      transition-timing-function: fade-in-out;
      transition-duration: 0.3s;
      .item-label-icon {
        color: var(--el-color-primary);
      }
    }    
}


.modal-loginout{
  width: 400px;
  height: 200px;
  opacity: 100;
  position: absolute;
  bottom: 0;
  left: 0;
}
</style>
