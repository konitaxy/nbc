<template>
  <el-container class="layout-cont">
    <el-header style="padding: 0;position: sticky; height:70px" class="mb-0 header">
      <Header v-on:expland="changeShadow"/>
    </el-header>
    <el-container :class="[isSider?'openside':'hideside',isMobile ? 'mobile': '']" style="flex-grow: 1;overflow-y: auto;">
      <el-row :class="[isShadowBg?'shadowBg':'']" @click="changeShadow()" />

      <el-aside class="main-cont main-left bg-white" >
        <!-- <div class="logo-container">
          <img src="@/assets/logo.png" class="logo-img" alt="logo" />
        </div> -->
        <Aside class="aside bg-light d-flex justify-content-start" />
      </el-aside>

      <!-- 分块滑动功能 -->
      <el-main class="main-cont main-right">
        <!-- <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
        </transition> -->
        
        <router-view v-if="reloadFlag" v-slot="{ Component,route }" :key="route.fullPath" v-loading="loadingFlag" element-loading-text="on loading">
          <transition mode="out-in" name="el-fade-in-linear">
            <keep-alive :include="routerStore.keepAliveRouters">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>      
      </el-main>
    </el-container>
    
  </el-container>
</template>

<script>
export default {
  name: 'Layout',
}
</script>

<script setup>
import Aside from '@/view/layout/aside/index.vue'
import Header from '@/view/layout/header/index.vue'
import HistoryComponent from '@/view/layout/aside/historyComponent/history.vue'
import BottomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import CustomPic from '@/components/customPic/index.vue'
import Setting from './setting/index.vue'
import { setUserAuthority } from '@/api/user'
import { emitter } from '@/utils/bus.js'
import { computed, ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import {setAvatar} from '@/api/profile'
import { ElMessage } from 'element-plus'
const router = useRouter()
const route = useRoute()
const routerStore = useRouterStore()
// 三种窗口适配
const isCollapse = ref(false)
const isSider = ref(true)
const isMobile = ref(false)
const initPage = () => {
  const screenWidth = document.body.clientWidth
  if (screenWidth < 992) {
    isMobile.value = true
    isSider.value = false
    isCollapse.value = true
  } else if (screenWidth >= 992 && screenWidth < 1200) {
    isMobile.value = false
    isSider.value = true
    isCollapse.value = false
  } else {
    isMobile.value = false
    isSider.value = true
    isCollapse.value = false
  }
}

initPage()

const loadingFlag = ref(false)
onMounted(() => {
  // 挂载一些通用的事件
  emitter.emit('collapse', isCollapse.value)
  emitter.emit('mobile', isMobile.value)
  emitter.on('reload', reload)
  emitter.on('showLoading', () => {
    loadingFlag.value = true
  })
  emitter.on('closeLoading', () => {
    loadingFlag.value = false
  })
  emitter.on('closeExpland',() =>{
    if(isMobile.value){
      isSider.value = false
      isShadowBg.value = false
      isCollapse.value = true
      emitter.emit('collapse', true)
    }
  })
  window.onresize = () => {
    return (() => {
      initPage()
      emitter.emit('collapse', isCollapse.value)
      emitter.emit('mobile', isMobile.value)
    })()
  }
  if (userStore.loadingInstance) {
    userStore.loadingInstance.close()
  }
})

const userStore = useUserStore()
const textColor = computed(() => {
  
  if (userStore.sideMode === 'dark') {
    return '#fff'
  } else if (userStore.sideMode === 'light') {
    return '#191a23'
  } else {
    return userStore.baseColor
  }
})

const backgroundColor = computed(() => {
  if (userStore.sideMode === 'dark') {
    return '#191a23'
  } else if (userStore.sideMode === 'light') {
    return '#fff'
  } else {
    return userStore.sideMode
  }
})

const matched = computed(() => route.meta.matched)

const changeUserAuth = async(id) => {
  const res = await setUserAuthority({
    authorityId: id
  })
  if (res.code === 0) {
    emitter.emit('closeAllPage')
    setTimeout(() => {
      window.location.reload()
    }, 1)
  }
}

const reloadFlag = ref(true)
const reload = async() => {
  if (route.meta.keepAlive) {
    reloadFlag.value = false
    await nextTick()
    reloadFlag.value = true
  } else {
    const title = route.meta.title
    router.push({ name: 'Reload', params: { title }})
  }
}
const isShadowBg = ref(false)
const totalCollapse = () => {
  isCollapse.value = !isCollapse.value
  isSider.value = !isCollapse.value
  isShadowBg.value = !isCollapse.value
  emitter.emit('collapse', isCollapse.value)
}

const toPerson = () => {
  router.push({ name: 'person' })
}
const changeShadow = () => {
  isShadowBg.value = !isShadowBg.value
  isSider.value = !!isCollapse.value
  totalCollapse()
}
const fileInput = ref(null)
const toggleFileInput =() =>{
  fileInput.value.click();
}
const handleFileChange = (event) =>{
  const file = event.target.files[0];
  if (!file) return;
  if (file) {
        if(file.size > 1024*1024){
          ElMessage.error('Image should be less than 1MB.')
          return
        }
        const reader = new FileReader();
        reader.onload = function(e) {
            const img = new Image();
            
            img.src = e.target.result;

            img.onload = function() {
                const { width, height } = img;
                if(width<300 || height < 300){
                    ElMessage.error('Image should be at least 300 by 300.')
                    event.target.value = null
                    return
                }
                const formData = new FormData();
                formData.append('file', file);

                // 上传头像
                setAvatar(formData).then(res =>{
                  if (res.code === 0){
                    userStore.GetUserInfo()
                    
                  }
                })
            };
        };

        // 将文件读取为 Data URL
        reader.readAsDataURL(file);
    }
}

</script>

<style lang="scss">
@use '@/style/mobile.scss';

.dark{
  background-color: #191a23 !important;
  color: #fff !important;
}
.light{
  background-color: #fff !important;
  color: #000 !important;
}
.user-info {
  // display: flex;
  align-items: center;
}

.edit-icon {
  position: absolute;
    top: 25px;
    right: calc(50% - 70px);
    cursor: pointer;
    height: 30px;
    width: 30px;
    border-radius: 50%;
    padding: 6px;
  // box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
}
.user-card {
  // width: 189px;
  height: 175px;
  top: 37px;
  left: 101px;
  gap: 0px;
  opacity: 0px;
  .card-name {
    font-family: Poppins;
    font-size: 18px;
    font-weight: 500;
    line-height: 27px;
    text-align: center;

  }
  .card-email {
    font-family: Poppins;
    font-size: 16px;
    font-weight: 300;
    line-height: 24px;
    text-align: center;
    color:var(--bs-tertiary-color)
  }

}

</style>
