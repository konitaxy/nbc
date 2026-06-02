<template>
  <div class="legal-page">
    <header class="site-nav">
      <router-link to="/" class="brand">
        <img :src="logo" alt="NEWBEECARD" />
        <div>
          <strong>NEWBEECARD</strong>
          <span>Privacy & acceptable use</span>
        </div>
      </router-link>
      <nav>
        <router-link to="/">首页</router-link>
        <router-link to="/about">About Us</router-link>
        <router-link to="/contact">Contact Us</router-link>
        <router-link to="/policy">隐私政策</router-link>
        <router-link to="/login" class="nav-cta">登录</router-link>
      </nav>
    </header>

    <main class="page-wrap">
      <section class="hero-panel">
        <p class="eyebrow">Privacy policy</p>
        <h1>隐私政策与合规使用说明</h1>
        <p>
          最近更新：2026-06-02。本政策说明 NEWBEECARD 在提供虚拟卡、充值、交易查询和客户支持服务时，可能如何收集、使用、保存、共享和保护信息。实际处理方式可能因地区、产品、合作方和监管要求而有所不同。
        </p>
      </section>

      <section class="policy-grid">
        <article>
          <h2>我们可能收集的信息</h2>
          <p>账户信息、联系方式、登录信息、设备与网络信息、交易记录、充值记录、客服沟通内容、KYC/KYB 或风控审核所需资料，以及法律法规或合作方要求的信息。</p>
        </article>
        <article>
          <h2>我们如何使用信息</h2>
          <p>用于账户注册与登录、身份核验、交易处理、风控监测、客户支持、账单展示、服务通知、争议处理、合规审查、反欺诈和履行法律义务。</p>
        </article>
        <article>
          <h2>我们如何共享信息</h2>
          <p>可能与支付、发卡、KYC、云服务、数据安全、客服、审计、法律顾问、监管机构或执法机构共享必要信息。我们不会出售用户个人信息。</p>
        </article>
      </section>

      <section class="split-section">
        <div>
          <p class="eyebrow">Data rights</p>
          <h2>用户权利与选择</h2>
          <p>
            你可以通过 support@newbeecard.com 请求访问、更正、删除或限制处理你的个人信息。部分请求可能因交易记录保存、反欺诈、争议处理、监管审计、税务或法律义务而无法立即完成。
          </p>
        </div>
        <ul class="check-list">
          <li>访问或获取账户相关个人信息副本</li>
          <li>更正不准确或过期的账户信息</li>
          <li>在适用法律允许范围内请求删除信息</li>
          <li>撤回非必要营销通信同意</li>
          <li>咨询跨境传输、数据保存和安全措施</li>
        </ul>
      </section>

      <section class="policy-grid">
        <article>
          <h2>安全措施</h2>
          <p>我们采用访问控制、加密传输、最小权限、日志监测、风控审核和供应商管理等措施保护信息。但任何互联网服务都无法保证绝对安全。</p>
        </article>
        <article>
          <h2>数据保存</h2>
          <p>我们仅在提供服务、处理争议、满足法律监管、审计、税务、反欺诈和安全目的所需期间保存信息。超过必要期限后会删除、匿名化或隔离保存。</p>
        </article>
        <article>
          <h2>Cookie 与分析</h2>
          <p>我们可能使用 Cookie、日志和类似技术维持登录状态、保障安全、分析服务表现并改进体验。你可通过浏览器设置管理 Cookie。</p>
        </article>
      </section>

      <section class="notice-panel">
        <p class="eyebrow">Acceptable use</p>
        <h2>禁止用途与风险声明</h2>
        <p>
          不得将 NEWBEECARD 用于任何违法、欺诈、洗钱、逃避制裁、侵犯知识产权、虚假交易、赌博、非法成人内容、盗刷、账号倒卖、规避平台限制或其他高风险活动。我们可根据法律、合作方要求或风控判断暂停、限制或终止服务，并保留必要记录用于调查。
        </p>
      </section>

      <section class="notice-panel">
        <p class="eyebrow">Disclaimer</p>
        <h2>免责声明</h2>
        <p>
          本政策为通用说明，不构成法律意见。若本政策与单独签署的服务协议、特定产品条款或适用法律冲突，以更严格或更具体的要求为准。我们可能不时更新本政策，重大变更会通过网站或账户通知展示。
        </p>
      </section>
    </main>

    <footer class="site-footer">
      <span>© 2026 NEWBEECARD. All rights reserved.</span>
      <div>
        <router-link to="/about">About Us</router-link>
        <router-link to="/contact">联系我们</router-link>
      </div>
    </footer>
  </div>
</template>

<script setup>
import logo from '@/assets/NEWBEECARD-logo.png'
</script>

<style scoped>
@import './legal-pages.css';
</style>
