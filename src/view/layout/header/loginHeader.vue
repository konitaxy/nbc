<template>
    <div class="h-100 w-100">
      <el-row justify="start">
        <el-col :span="10">
          <img class="logo" :src="logo"></img>
        </el-col>
        <el-col :xs="14" :sm="14" class="hstack gap-2 justify-content-end">
          
          <a class="text-dark hover-pointer fw-medium d-none d-sm-block me-2" href="https://metalposter.com" target="_blank">
            Store
          </a>
          <a class="text-dark hover-pointer fw-medium d-none d-sm-block me-2" href="https://metalposter.com/pages/about-us" target="_blank">
            About
          </a>
          <a class="text-dark hover-pointer fw-medium d-none d-sm-block me-4" href="https://metalposter.com/blogs/news" target="_blank">
            Blog
          </a>
          <el-dropdown class="d-sm-none" style="--el-dropdown-menuItem-hover-color:aliceblue">
            <span class="el-dropdown-link" style="line-height: 3em;">
                Store
              <el-icon style="top:0.2em">
                <arrow-down />
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <a class="dropdown-item text-dark fw-medium" href="https://metalposter.com" target="_blank">
                    Store
                  </a>
              </el-dropdown-item>
                <el-dropdown-item>
                  <a class="dropdown-item text-dark fw-medium" href="https://metalposter.com/pages/about-us" target="_blank">
                      About
                  </a>
                </el-dropdown-item>
                <el-dropdown-item>
                  <a class="dropdown-item text-dark fw-medium" href="https://metalposter.com/blogs/news" target="_blank">
                    Blog
                  </a>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <div class="btn btn-danger btn-sm" @click="handleSignIn">
            Sign in
          </div>
          <div class="btn btn-outline-secondary btn-sm" @click="handleSignUp">
            Apply
          </div>
      </el-col>
      </el-row>
      </div>
  </template>
  
  <script>
export default {
  name: 'LoginHeader',
}
</script>
<script setup>
import {ref} from 'vue'
import { emitter } from '@/utils/bus.js'
import logoArt from '@/assets/logo.png'
  // No script needed in this example.

const emit = defineEmits(['signIn','signUp']);
const isMobile = ref(false)
const logo = logoArt

const handleSignIn = () => {
  emit('signIn')
}
const handleSignUp = () => {
  emit('signUp')
}
  // No script needed in this example.
</script>
  
  <style scoped>
  :deep(.el-dropdown-link:focus-visible) {
    outline:none !important
  }
  .header {
    padding: 0;
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  }
  .right-icons-container {
    width: Fixed (205.22px)px;
    height: Hug (25px)px;
    top: 20px;
    left: 1366.78px;
    gap: 32px;
    opacity: 0px;

  }
  .right-icons {
    width: 8.13px;
    height: 13.46px;
    top:5.23px;
    left:7.88px;
    margin-bottom: 10px;
    color: var(--bs-secondary-color)
  }
  .right-icons-mobile {
    width: 8.13px;
    height: 13.46px;
    top:5.23px;
    left:7.88px;
    gap:32px;
    margin-bottom: 10px;
    color: var(--bs-secondary-color)
  }
  .logo {
    padding:10px 10px 10px 20px;
    max-width: 186px; height:65px;object-fit:fill;
    /* margin-left:10px */
  }
  @media (max-width: 767px) {
    .logo {
      padding:19px 10px 15px 10px;
    }
  }
  </style>