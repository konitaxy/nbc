<template>
  <div class="legal-page">
    <header class="site-nav">
      <router-link to="/" class="brand">
        <img :src="logo" alt="NEWBEECARD" />
        <div>
          <strong>NEWBEECARD</strong>
          <span>Global subscription cards</span>
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
        <p class="eyebrow">About NEWBEECARD</p>
        <h1>为海外订阅、广告充值与团队用卡提供更清晰的虚拟卡管理方式</h1>
        <p>
          NEWBEECARD 面向个人创作者、跨境团队、广告代理和 SaaS 使用者，提供虚拟卡开卡、余额管理、交易记录和客服支持等服务。我们关注的是可控、透明和可追踪，而不是鼓励规避任何平台、银行或地区规则。
        </p>
      </section>

      <section class="split-section">
        <div class="copy-block">
          <p class="eyebrow">What we do</p>
          <h2>我们解决的核心问题</h2>
          <p>
            跨境订阅经常涉及多平台、多币种、多成员和多张卡的管理。NEWBEECARD 帮助用户把开卡、充值、限额、冻结、交易查询和账单归档集中到一个控制台，减少人工记录和共用卡片带来的风险。
          </p>
        </div>
        <div class="feature-grid">
          <article>
            <span>01</span>
            <h3>虚拟卡管理</h3>
            <p>按用途、项目或成员创建卡片，设置额度并追踪交易。</p>
          </article>
          <article>
            <span>02</span>
            <h3>订阅与广告场景</h3>
            <p>适用于 AI 工具、应用商店、广告账户和 SaaS 订阅等常见付款需求。</p>
          </article>
          <article>
            <span>03</span>
            <h3>账单可追踪</h3>
            <p>以交易记录和用途标签降低团队对账成本。</p>
          </article>
          <article>
            <span>04</span>
            <h3>支持与风控</h3>
            <p>提供客服入口和异常交易处理流程，帮助用户及时处理问题。</p>
          </article>
        </div>
      </section>

      <section class="policy-grid">
        <article>
          <h2>合规使用承诺</h2>
          <p>
            用户应遵守所在地法律法规、发卡机构规则、支付网络规则以及目标平台条款。NEWBEECARD 不支持洗钱、欺诈、规避制裁、虚假交易、盗刷、账号倒卖、赌博、成人非法内容、侵犯知识产权或任何被法律或合作方禁止的用途。
          </p>
        </article>
        <article>
          <h2>服务边界</h2>
          <p>
            虚拟卡是否可用于某个第三方平台，可能受到平台政策、地区、风控、卡组织、发卡机构和商户规则影响。我们不保证所有平台均可持续成功扣款，也不代表任何第三方平台、银行或卡组织作出承诺。
          </p>
        </article>
        <article>
          <h2>资金与交易责任</h2>
          <p>
            用户应妥善保管账号、密码、验证码和卡片信息，并对其账号下发起的授权交易负责。发现异常交易、账号泄露或疑似欺诈时，应立即联系我们并配合必要的身份和交易核验。
          </p>
        </article>
      </section>

      <section class="notice-panel">
        <p class="eyebrow">Legal notice</p>
        <h2>重要说明</h2>
        <p>
          本页面仅用于介绍服务和风险边界，不构成法律、税务、金融、投资或合规建议。不同地区对支付、隐私、KYC、AML、税务和消费者保护的要求可能不同，请在实际运营前咨询具备资质的专业顾问。
        </p>
      </section>
    </main>

    <footer class="site-footer">
      <span>© 2026 NEWBEECARD. All rights reserved.</span>
      <div>
        <router-link to="/policy">隐私政策</router-link>
        <router-link to="/contact">联系我们</router-link>
      </div>
    </footer>
  </div>
</template>

<script setup>
import logo from '@/assets/NEWBEECARD-logo.png'
</script>

<style scoped>
.legal-page {
  --bg: #020812;
  --panel: rgba(5, 16, 29, 0.84);
  --line: rgba(139, 214, 255, 0.16);
  --text: #f4fbff;
  --muted: rgba(225, 242, 255, 0.74);
  --dim: rgba(225, 242, 255, 0.56);
  --cyan: #44d5ff;
  --green: #7dffcc;
  min-height: 100vh;
  padding: 24px;
  color: var(--text);
  background:
    radial-gradient(circle at 12% 8%, rgba(68, 213, 255, 0.1), transparent 28%),
    radial-gradient(circle at 88% 20%, rgba(47, 125, 255, 0.12), transparent 26%),
    linear-gradient(180deg, #020812 0%, #030b16 52%, #04101d 100%);
}

.site-nav,
.page-wrap,
.site-footer {
  width: min(1180px, calc(100% - 24px));
  margin: 0 auto;
}

.site-nav,
.brand,
.site-nav nav,
.site-footer,
.site-footer div {
  display: flex;
  align-items: center;
}

.site-nav {
  justify-content: space-between;
  gap: 24px;
  padding: 4px 0 46px;
}

.brand {
  gap: 12px;
  color: var(--text);
  text-decoration: none;
}

.brand img {
  width: 54px;
  height: 54px;
  object-fit: contain;
}

.brand strong {
  display: block;
  letter-spacing: 0.08em;
}

.brand span {
  color: var(--dim);
  font-size: 12px;
}

.site-nav nav,
.site-footer div {
  flex-wrap: wrap;
  gap: 14px;
}

.site-nav a,
.site-footer a {
  color: var(--muted);
  text-decoration: none;
}

.site-nav nav a {
  padding: 10px 14px;
  border-radius: 999px;
}

.site-nav nav a:hover,
.site-footer a:hover {
  color: var(--cyan);
}

.nav-cta {
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.055);
}

.page-wrap {
  display: grid;
  gap: 22px;
}

.hero-panel,
.split-section,
.policy-grid article,
.notice-panel {
  border: 1px solid var(--line);
  border-radius: 32px;
  background: var(--panel);
  backdrop-filter: blur(14px);
}

.hero-panel {
  padding: clamp(30px, 5vw, 58px);
}

.eyebrow {
  display: inline-flex;
  margin: 0 0 18px;
  color: var(--green);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

h1,
h2,
h3,
p {
  margin: 0;
}

h1 {
  max-width: 900px;
  font-size: clamp(2.3rem, 5vw, 4.8rem);
  line-height: 1.02;
  letter-spacing: -0.055em;
}

h2 {
  font-size: clamp(1.45rem, 2.4vw, 2.4rem);
  line-height: 1.12;
}

h3 {
  font-size: 1.08rem;
}

p {
  color: var(--muted);
  line-height: 1.82;
}

.hero-panel p:last-child,
.copy-block p {
  max-width: 780px;
  margin-top: 22px;
}

.split-section {
  display: grid;
  grid-template-columns: 0.85fr 1.15fr;
  gap: 28px;
  padding: 28px;
}

.copy-block {
  align-self: center;
}

.feature-grid,
.policy-grid {
  display: grid;
  gap: 16px;
}

.feature-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.feature-grid article {
  padding: 22px;
  border: 1px solid var(--line);
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.045);
}

.feature-grid span {
  color: var(--cyan);
  font-weight: 800;
  letter-spacing: 0.12em;
}

.feature-grid h3 {
  margin: 28px 0 10px;
}

.policy-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.policy-grid article,
.notice-panel {
  padding: 28px;
}

.policy-grid h2,
.notice-panel h2 {
  margin-bottom: 14px;
}

.notice-panel {
  border-color: rgba(125, 255, 204, 0.24);
}

.site-footer {
  justify-content: space-between;
  gap: 20px;
  padding: 34px 0 10px;
  color: var(--dim);
}

@media (max-width: 900px) {
  .site-nav,
  .site-footer {
    align-items: flex-start;
    flex-direction: column;
  }

  .split-section,
  .policy-grid {
    grid-template-columns: 1fr;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }
}
</style>
