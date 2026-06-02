<template>
  <div class="legal-page">
    <header class="site-nav">
      <router-link to="/" class="brand">
        <img :src="logo" alt="NEWBEECARD" />
        <div>
          <strong>NEWBEECARD</strong>
          <span>Support & compliance desk</span>
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
        <p class="eyebrow">Contact us</p>
        <h1>联系 NEWBEECARD 团队</h1>
        <p>
          如果你需要账户支持、交易查询、商务合作、合规反馈或投诉处理，请通过下方渠道联系我们。为了保护账户安全，请不要在非官方渠道发送完整卡号、密码、验证码、私钥或身份文件。
        </p>
      </section>

      <section class="contact-grid">
        <article>
          <span>Support</span>
          <h2>客服支持</h2>
          <p>support@newbeecard.com</p>
          <small>账户、登录、充值、卡片状态、交易查询</small>
        </article>
        <article>
          <span>Business</span>
          <h2>商务合作</h2>
          <p>bd@newbeecard.com</p>
          <small>渠道合作、企业方案、批量用卡、代理合作</small>
        </article>
        <article>
          <span>Compliance</span>
          <h2>合规与投诉</h2>
          <p>compliance@newbeecard.com</p>
          <small>滥用举报、身份核验、合规请求、执法协助</small>
        </article>
      </section>

      <section class="split-section">
        <div>
          <p class="eyebrow">Before contacting</p>
          <h2>提交工单前请准备这些信息</h2>
          <p>完整、准确的信息可以帮助我们更快定位问题。请仅通过官方邮箱提交敏感材料，且不要发送登录密码或验证码。</p>
        </div>
        <ul class="check-list">
          <li>注册邮箱或用户 ID</li>
          <li>问题发生时间和相关订单号</li>
          <li>交易金额、币种、商户名称和卡片后四位</li>
          <li>错误提示截图，但请遮挡完整卡号和敏感身份信息</li>
          <li>如果是投诉或合规举报，请提供可核实的事实说明</li>
        </ul>
      </section>

      <section class="policy-grid">
        <article>
          <h2>响应时间</h2>
          <p>普通咨询通常在 1-2 个工作日内回复；账户安全、异常交易、合规举报会优先处理。复杂问题可能需要更多核验时间。</p>
        </article>
        <article>
          <h2>安全提醒</h2>
          <p>NEWBEECARD 不会通过私聊索要你的密码、验证码或完整卡片信息。任何自称客服并要求转账或提供敏感凭证的行为都应立即举报。</p>
        </article>
        <article>
          <h2>法律与执法请求</h2>
          <p>执法、监管或法律请求应发送至 compliance@newbeecard.com，并包含正式文件、联系人、机构信息和合法依据。</p>
        </article>
      </section>

      <section class="notice-panel">
        <p class="eyebrow">Office information</p>
        <h2>服务区域说明</h2>
        <p>
          NEWBEECARD 可能根据服务可用性、监管要求、风控政策或合作方要求限制部分地区、行业、账户或交易。页面展示的联系方式不代表我们可在所有司法辖区提供服务。
        </p>
      </section>
    </main>

    <footer class="site-footer">
      <span>© 2026 NEWBEECARD. All rights reserved.</span>
      <div>
        <router-link to="/about">About Us</router-link>
        <router-link to="/policy">隐私政策</router-link>
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
