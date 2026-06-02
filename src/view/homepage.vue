<template>
  <div class="homepage">
    <div class="aurora aurora-one"></div>
    <div class="aurora aurora-two"></div>

    <header class="site-nav" :class="{ 'is-menu-open': isMobileMenuOpen }">
      <router-link to="/" class="brand" aria-label="NEWBEECARD home">
        <img :src="brandLogo" alt="NEWBEECARD" class="brand-logo" />
        <div class="brand-copy">
          <strong>NEWBEECARD</strong>
          <span>Global subscription cards</span>
        </div>
      </router-link>

      <button
        class="menu-toggle"
        type="button"
        :aria-expanded="isMobileMenuOpen"
        aria-label="Toggle navigation menu"
        @click="isMobileMenuOpen = !isMobileMenuOpen"
      >
        <span></span>
        <span></span>
        <span></span>
      </button>

      <nav class="nav-links" aria-label="Primary navigation">
        <a href="#solutions" @click="closeMobileMenu">应用场景</a>
        <a href="#control" @click="closeMobileMenu">卡片管理</a>
        <a href="#pricing" @click="closeMobileMenu">价格明细</a>
        <a href="#faq" @click="closeMobileMenu">常见问题</a>
      </nav>

      <div class="nav-actions">
        <router-link to="/login" class="nav-login" @click="closeMobileMenu">登录</router-link>
        <router-link to="/login" class="nav-cta" @click="closeMobileMenu">注册开卡</router-link>
      </div>
    </header>

    <main>
      <section class="hero">
        <div class="hero-copy">
          <p class="eyebrow">Virtual cards for global tools</p>
          <h1>海外订阅、广告充值、AI 工具付款，一张虚拟卡就够了</h1>
          <p class="hero-text">
            NEWBEECARD 为个人创作者、跨境团队和广告代理提供可控额度、快速开卡、透明账单的虚拟卡服务。覆盖 ChatGPT、App Store、Google Play、TikTok Ads、Meta Ads、Netflix、Spotify 等常见海外平台。
          </p>

          <div class="hero-actions">
            <router-link to="/login" class="button button-primary">立即开卡</router-link>
            <a href="#pricing" class="button button-secondary">查看费率</a>
          </div>

          <div class="hero-metrics" aria-label="Service highlights">
            <div>
              <strong>50万+</strong>
              <span>用户与团队使用</span>
            </div>
            <div>
              <strong>50+</strong>
              <span>可选卡头资源</span>
            </div>
            <div>
              <strong>24/7</strong>
              <span>在线客服支持</span>
            </div>
          </div>
        </div>

        <div class="hero-visual" aria-label="NEWBEECARD product preview">
          <div class="card-orbit card-orbit-one">
            <span>AI Subscription</span>
            <strong>$20.00</strong>
            <small>Auto renewed</small>
          </div>
          <div class="card-orbit card-orbit-two">
            <span>Ads Budget</span>
            <strong>$500.00</strong>
            <small>Daily limit ready</small>
          </div>
          <img :src="heroIllustration" alt="NEWBEECARD virtual card dashboard" />
          <div class="settlement-panel">
            <div>
              <span>本月消费</span>
              <strong>$8,426.32</strong>
            </div>
            <div class="progress-row">
              <i style="width: 72%"></i>
            </div>
            <p>AI 工具、广告投放、应用商店订阅分账清晰可追踪</p>
          </div>
        </div>
      </section>

      <section class="ticker" aria-label="Supported platforms">
        <div class="ticker-track">
          <div class="ticker-group">
            <span><i class="bi bi-robot"></i>ChatGPT</span>
            <span><i class="platform-mark">C</i>Claude</span>
            <span><i class="bi bi-stars"></i>Midjourney</span>
            <span><i class="bi bi-apple"></i>App Store</span>
            <span><i class="bi bi-google-play"></i>Google Play</span>
            <span><i class="bi bi-tiktok"></i>TikTok Ads</span>
            <span><i class="bi bi-meta"></i>Meta Ads</span>
            <span><i class="platform-mark">N</i>Netflix</span>
            <span><i class="bi bi-spotify"></i>Spotify</span>
          </div>
          <div class="ticker-group" aria-hidden="true">
            <span><i class="bi bi-robot"></i>ChatGPT</span>
            <span><i class="platform-mark">C</i>Claude</span>
            <span><i class="bi bi-stars"></i>Midjourney</span>
            <span><i class="bi bi-apple"></i>App Store</span>
            <span><i class="bi bi-google-play"></i>Google Play</span>
            <span><i class="bi bi-tiktok"></i>TikTok Ads</span>
            <span><i class="bi bi-meta"></i>Meta Ads</span>
            <span><i class="platform-mark">N</i>Netflix</span>
            <span><i class="bi bi-spotify"></i>Spotify</span>
          </div>
        </div>
      </section>

      <section class="section solutions" id="solutions">
        <div class="section-heading">
          <p class="eyebrow">Where it works</p>
          <h2>把高频海外付款场景拆成可管理的卡片组合</h2>
          <p>不同服务使用不同卡片、额度和账单标签，团队成员不用共用一张实体信用卡，财务也不用在截图里找记录。</p>
        </div>

        <div class="solution-grid">
          <article class="solution-card solution-large">
            <div class="card-index">01</div>
            <h3>AI 工具订阅</h3>
            <p>适合 ChatGPT、Claude、Midjourney、Cursor、Kiro 等工具续费，支持按项目分配卡片和预算。</p>
            <ul>
              <li>订阅续费提醒</li>
              <li>失败交易快速排查</li>
              <li>按成员独立记账</li>
            </ul>
          </article>

          <article class="solution-card">
            <div class="card-index">02</div>
            <h3>广告投放充值</h3>
            <p>TikTok、Meta、Google Ads 广告账户可按客户、渠道、项目拆分付款。</p>
          </article>

          <article class="solution-card">
            <div class="card-index">03</div>
            <h3>应用商店付款</h3>
            <p>覆盖 App Store、Google Play、应用内购买和开发者工具订阅。</p>
          </article>

          <article class="solution-card solution-accent">
            <div class="card-index">04</div>
            <h3>流媒体与会员</h3>
            <p>Netflix、Spotify、YouTube Premium 等服务可用独立小额卡控制风险。</p>
          </article>

          <article class="solution-card">
            <div class="card-index">05</div>
            <h3>企业团队用卡</h3>
            <p>给运营、设计、开发、投放团队分别开卡，权限、额度、账单互不干扰。</p>
          </article>
        </div>
      </section>

      <section class="section control" id="control">
        <div class="control-panel">
          <div class="control-copy">
            <p class="eyebrow">Control center</p>
            <h3>开卡、充值、限额、冻结和对账，都在同一个控制台完成</h3>
            <p>
              页面不只是展示卡号。NEWBEECARD 把消费控制、交易追踪、团队协作和客服处理放进同一个工作流，让海外付款从临时操作变成稳定流程。
            </p>
          </div>

          <div class="control-list">
            <div class="control-item">
              <span>01</span>
              <strong>秒级创建虚拟卡</strong>
              <p>选择卡头、国家、用途标签和初始额度，创建后即可复制卡信息使用。</p>
            </div>
            <div class="control-item">
              <span>02</span>
              <strong>额度与风控可调整</strong>
              <p>按订阅周期或广告预算设置上限，异常消费可以即时冻结。</p>
            </div>
            <div class="control-item">
              <span>03</span>
              <strong>账单自动归档</strong>
              <p>按平台、成员、项目、客户维度筛选交易，减少人工对账成本。</p>
            </div>
          </div>
        </div>
      </section>

      <section class="section workflow">
        <div class="section-heading">
          <p class="eyebrow">Launch fast</p>
          <h2 class="compact-title">三步开始使用</h2>
        </div>

        <div class="workflow-grid">
          <article>
            <span>1</span>
            <h3>注册并进入控制台</h3>
            <p>登录后完善基础信息，选择个人或团队使用方式。</p>
          </article>
          <article>
            <span>2</span>
            <h3>充值并选择卡头</h3>
            <p>支持 USDT 充值，按使用场景选择美国或香港虚拟卡。</p>
          </article>
          <article>
            <span>3</span>
            <h3>复制卡信息付款</h3>
            <p>在目标平台完成订阅或充值，后续交易自动进入账单。</p>
          </article>
        </div>
      </section>

      <section class="section pricing" id="pricing">
        <div class="section-heading">
          <p class="eyebrow">Transparent pricing</p>
          <h2 class="compact-title">费率清楚，适合从个人试用到团队规模化使用</h2>
        </div>

        <div class="pricing-grid">
          <article class="price-card price-highlight">
            <span class="price-label">开卡费</span>
            <strong>$1 起</strong>
            <p>低门槛体验虚拟卡，适合首次订阅或小额测试。</p>
            <ul>
              <li>无限开卡</li>
              <li>50+ 卡头选择</li>
              <li>美国 / 香港虚拟卡</li>
            </ul>
          </article>
          <article class="price-card">
            <span class="price-label">充值手续费</span>
            <strong>1%</strong>
            <p>大额和小额充值都按清晰费率计算，便于预算管理。</p>
          </article>
          <article class="price-card">
            <span class="price-label">稳定币充值</span>
            <strong>USDT</strong>
            <p>支持主流充值方式，到账后即可用于开卡和补充余额。</p>
          </article>
          <article class="price-card">
            <span class="price-label">团队方案</span>
            <strong>定制</strong>
            <p>适合广告代理、跨境团队、批量订阅和多客户对账。</p>
          </article>
        </div>
      </section>

      <section class="section faq" id="faq">
        <div class="section-heading">
          <p class="eyebrow">Questions</p>
          <h2 class="compact-title">常见问题</h2>
        </div>

        <div class="faq-grid">
          <article>
            <h3>可以用于哪些平台？</h3>
            <p>常见 AI 工具、应用商店、广告平台、流媒体和 SaaS 订阅都可以按场景尝试使用。</p>
          </article>
          <article>
            <h3>如何控制订阅风险？</h3>
            <p>每张卡都可以单独设置额度、用途和状态，订阅异常时可冻结或换卡。</p>
          </article>
          <article>
            <h3>团队如何对账？</h3>
            <p>建议按成员、客户或项目创建不同卡片，交易记录天然分组，后续筛选更清楚。</p>
          </article>
          <article>
            <h3>充值后多久能用？</h3>
            <p>余额到账后即可开卡，具体到账时间以充值渠道和后台状态为准。</p>
          </article>
        </div>
      </section>

      <section class="final-cta">
        <div>
          <p class="eyebrow">Ready for global payments</p>
          <h2>把海外订阅和广告付款交给一套可控的虚拟卡系统</h2>
        </div>
        <router-link to="/login" class="button button-primary">马上开卡</router-link>
      </section>
    </main>

    <footer class="site-footer">
      <div class="footer-brand">
        <img :src="brandLogo" alt="NEWBEECARD" />
        <div>
          <strong>NEWBEECARD</strong>
          <p>一站式海外订阅、AI 工具和广告充值虚拟卡服务。</p>
        </div>
      </div>
      <div class="footer-links">
        <router-link to="/">首页</router-link>
        <router-link to="/about">About Us</router-link>
        <router-link to="/contact">Contact Us</router-link>
        <router-link to="/policy">隐私政策</router-link>
      </div>
      <p class="copyright">© 2026 NEWBEECARD. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import heroIllustration from '@/assets/hero-visual.png'
import brandLogo from '@/assets/NEWBEECARD-logo.png'

const isMobileMenuOpen = ref(false)

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}
</script>

<style scoped>
.homepage {
  --bg: #020812;
  --panel: rgba(5, 16, 29, 0.84);
  --panel-strong: rgba(4, 13, 24, 0.96);
  --line: rgba(139, 214, 255, 0.16);
  --line-strong: rgba(139, 214, 255, 0.34);
  --text: #f4fbff;
  --muted: rgba(225, 242, 255, 0.74);
  --dim: rgba(225, 242, 255, 0.56);
  --cyan: #44d5ff;
  --blue: #2f7dff;
  --green: #7dffcc;
  --gold: #ffd166;
  position: relative;
  min-height: 100vh;
  overflow-x: hidden;
  background:
    radial-gradient(circle at 12% 8%, rgba(68, 213, 255, 0.09), transparent 28%),
    radial-gradient(circle at 84% 18%, rgba(47, 125, 255, 0.11), transparent 26%),
    linear-gradient(180deg, #020812 0%, #030b16 52%, #04101d 100%);
  color: var(--text);
  font-family: "Avenir Next", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.aurora {
  position: fixed;
  width: 34vw;
  height: 34vw;
  min-width: 320px;
  min-height: 320px;
  border-radius: 999px;
  filter: blur(70px);
  opacity: 0.2;
  pointer-events: none;
  z-index: 0;
  animation: floatGlow 12s ease-in-out infinite;
}

.aurora-one {
  top: 14%;
  left: -8%;
  background: #1fa6d6;
}

.aurora-two {
  right: -10%;
  bottom: 8%;
  background: #1d55c9;
  animation-delay: -5s;
}

@keyframes floatGlow {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(24px, -28px, 0) scale(1.08); }
}

.site-nav,
.hero,
.ticker,
.section,
.final-cta,
.site-footer {
  position: relative;
  z-index: 1;
}

.site-nav {
  width: min(1180px, calc(100% - 36px));
  margin: 0 auto;
  padding: 22px 0;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 24px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: var(--text);
  text-decoration: none;
}

.brand-logo {
  width: 58px;
  height: 58px;
  object-fit: contain;
  display: block;
}

.brand-copy {
  display: grid;
  gap: 2px;
}

.brand-copy strong {
  font-size: 17px;
  letter-spacing: 0.08em;
}

.brand-copy span {
  color: var(--dim);
  font-size: 12px;
}

.nav-links {
  justify-self: center;
  display: flex;
  gap: 24px;
  padding: 12px 18px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: rgba(4, 15, 29, 0.5);
  backdrop-filter: blur(14px);
}

.nav-links a,
.footer-links a {
  color: var(--muted);
  text-decoration: none;
  transition: color 0.2s ease;
}

.nav-links a {
  font-size: 14px;
}

.nav-links a:hover,
.footer-links a:hover {
  color: var(--cyan);
}

.menu-toggle {
  display: none;
  width: 46px;
  height: 46px;
  padding: 0;
  border: 1px solid var(--line);
  border-radius: 16px;
  background: rgba(4, 15, 29, 0.68);
  cursor: pointer;
  align-content: center;
  gap: 5px;
  place-items: center;
}

.menu-toggle span {
  display: block;
  width: 19px;
  height: 2px;
  border-radius: 999px;
  background: var(--text);
  transition: transform 0.22s ease, opacity 0.22s ease;
}

.is-menu-open .menu-toggle span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}

.is-menu-open .menu-toggle span:nth-child(2) {
  opacity: 0;
}

.is-menu-open .menu-toggle span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.nav-login,
.nav-cta,
.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  border-radius: 999px;
  font-weight: 700;
  transition: transform 0.22s ease, border-color 0.22s ease, background 0.22s ease, box-shadow 0.22s ease;
}

.nav-login {
  color: var(--text);
  padding: 11px 16px;
  border: 1px solid var(--line);
}

.nav-cta {
  color: #04111e;
  padding: 12px 18px;
  background: linear-gradient(135deg, var(--green), var(--cyan));
  box-shadow: 0 16px 34px rgba(68, 213, 255, 0.22);
}

.nav-login:hover,
.nav-cta:hover,
.button:hover {
  transform: translateY(-2px);
}

.hero {
  width: min(1180px, calc(100% - 36px));
  min-height: 690px;
  margin: 0 auto;
  padding: 72px 0 54px;
  display: grid;
  grid-template-columns: minmax(0, 1.02fr) minmax(420px, 0.98fr);
  align-items: center;
  gap: 50px;
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 18px;
  padding: 8px 13px;
  border: 1px solid rgba(125, 255, 204, 0.24);
  border-radius: 999px;
  background: rgba(125, 255, 204, 0.08);
  color: var(--green);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.eyebrow::before {
  content: "";
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 18px currentColor;
}

.hero h1,
.section-heading h2,
.control-copy h3,
.final-cta h2 {
  margin: 0;
  letter-spacing: -0.055em;
}

.hero h1 {
  max-width: 720px;
  font-size: clamp(2.45rem, 4.7vw, 5.1rem);
  line-height: 1.02;
}

.hero-text {
  max-width: 650px;
  margin: 28px 0 0;
  color: var(--muted);
  font-size: 17px;
  line-height: 1.9;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 32px;
}

.button {
  min-height: 54px;
  padding: 0 24px;
}

.button-primary {
  color: #04111e;
  background: linear-gradient(135deg, var(--green), var(--cyan) 58%, #83a8ff);
  box-shadow: 0 22px 48px rgba(68, 213, 255, 0.24);
}

.button-secondary {
  color: var(--text);
  border: 1px solid var(--line-strong);
  background: rgba(8, 24, 43, 0.68);
}

.hero-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  margin-top: 38px;
  max-width: 650px;
}

.hero-metrics div {
  padding: 18px;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: rgba(8, 24, 43, 0.58);
  backdrop-filter: blur(12px);
}

.hero-metrics strong {
  display: block;
  color: #ffffff;
  font-size: 28px;
  line-height: 1;
}

.hero-metrics span {
  display: block;
  margin-top: 8px;
  color: var(--dim);
  font-size: 13px;
}

.hero-visual {
  position: relative;
  min-height: 590px;
  display: grid;
  place-items: center;
}

.hero-visual::before {
  content: none;
}

.hero-visual img {
  position: relative;
  width: min(105%, 640px);
  max-height: 520px;
  object-fit: contain;
  filter: drop-shadow(0 42px 72px rgba(0, 0, 0, 0.62));
  transform: translateX(12px) scale(1.06);
  z-index: 2;
}

.card-orbit,
.settlement-panel {
  position: absolute;
  z-index: 3;
  border: 1px solid rgba(139, 214, 255, 0.18);
  background: rgba(5, 17, 31, 0.42);
  backdrop-filter: blur(10px);
  box-shadow: 0 22px 48px rgba(0, 0, 0, 0.2);
}

.card-orbit {
  width: 190px;
  padding: 18px;
  border-radius: 22px;
}

.card-orbit span,
.settlement-panel span {
  color: var(--dim);
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.12em;
}

.card-orbit strong,
.settlement-panel strong {
  display: block;
  margin: 8px 0 4px;
  font-size: 28px;
}

.card-orbit small {
  color: var(--green);
}

.card-orbit-one {
  top: 60px;
  left: 4px;
  animation: cardDrift 6s ease-in-out infinite;
}

.card-orbit-two {
  right: 0;
  top: 135px;
  animation: cardDrift 7s ease-in-out infinite reverse;
}

@keyframes cardDrift {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-16px); }
}

.settlement-panel {
  left: 36px;
  right: 34px;
  bottom: 42px;
  padding: 20px;
  border-radius: 26px;
}

.progress-row {
  height: 10px;
  margin: 14px 0;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
}

.progress-row i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, var(--green), var(--cyan), var(--blue));
}

.settlement-panel p {
  margin: 0;
  color: var(--muted);
  line-height: 1.6;
}

.ticker {
  width: min(1180px, calc(100% - 36px));
  margin: 0 auto;
  padding: 18px;
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 24px;
  background: rgba(4, 15, 29, 0.58);
  mask-image: linear-gradient(90deg, transparent 0%, #000 9%, #000 91%, transparent 100%);
}

.ticker-track {
  display: flex;
  width: max-content;
  animation: platformMarquee 28s linear infinite;
}

.ticker:hover .ticker-track {
  animation-play-state: paused;
}

.ticker-group {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-right: 12px;
}

.ticker span {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 999px;
  color: var(--muted);
  background: rgba(255, 255, 255, 0.045);
  white-space: nowrap;
}

@keyframes platformMarquee {
  from {
    transform: translateX(0);
  }

  to {
    transform: translateX(-50%);
  }
}

@media (prefers-reduced-motion: reduce) {
  .ticker-track {
    animation: none;
  }
}

.ticker i {
  display: inline-grid;
  width: 24px;
  height: 24px;
  place-items: center;
  border-radius: 50%;
  color: #dff8ff;
  background: rgba(68, 213, 255, 0.1);
  font-size: 13px;
  font-style: normal;
  line-height: 1;
}

.ticker .bi-apple {
  background: rgba(255, 255, 255, 0.1);
}

.ticker .bi-google-play {
  color: #7dffcc;
}

.ticker .bi-tiktok {
  color: #ff6f9f;
}

.ticker .bi-meta {
  color: #7fb3ff;
}

.ticker .bi-spotify {
  color: #74f58c;
}

.ticker .platform-mark {
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0;
}

.section {
  width: min(1180px, calc(100% - 36px));
  margin: 0 auto;
  padding: 92px 0 0;
}

.section-heading {
  max-width: 780px;
  margin-bottom: 32px;
}

.section-heading h2,
.final-cta h2 {
  font-size: clamp(2.2rem, 4vw, 4rem);
  line-height: 1.05;
}

.section-heading .compact-title {
  font-size: clamp(1.75rem, 2.6vw, 2.75rem);
  line-height: 1.12;
}

.section-heading p:not(.eyebrow),
.control-copy p {
  margin: 18px 0 0;
  color: var(--muted);
  font-size: 16px;
  line-height: 1.85;
}

.solution-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.solution-card,
.price-card,
.faq article,
.workflow-grid article {
  border: 1px solid var(--line);
  background: var(--panel);
  backdrop-filter: blur(14px);
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.18);
}

.solution-card {
  min-height: 250px;
  padding: 26px;
  border-radius: 28px;
  display: flex;
  flex-direction: column;
}

.solution-large {
  grid-column: span 2;
}

.solution-accent {
  background:
    linear-gradient(135deg, rgba(125, 255, 204, 0.15), transparent 52%),
    var(--panel-strong);
  border-color: rgba(125, 255, 204, 0.28);
}

.card-index,
.price-label {
  color: var(--green);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.solution-card h3 {
  margin: 42px 0 12px;
  font-size: 24px;
}

.solution-card p,
.solution-card li,
.price-card p,
.price-card li,
.faq p,
.workflow-grid p {
  color: var(--muted);
  line-height: 1.78;
}

.solution-card ul,
.price-card ul {
  display: grid;
  gap: 10px;
  margin: auto 0 0;
  padding: 0;
  list-style: none;
}

.solution-card ul {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.solution-card li,
.price-card li {
  position: relative;
  padding-left: 18px;
}

.solution-card li {
  padding: 8px 12px 8px 24px;
  border: 1px solid rgba(139, 214, 255, 0.14);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.045);
  font-size: 13px;
}

.solution-card li::before,
.price-card li::before {
  content: "";
  position: absolute;
  left: 0;
  top: 12px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--cyan);
}

.solution-card li::before {
  left: 11px;
  top: 15px;
  width: 6px;
  height: 6px;
}

.control-panel {
  display: grid;
  grid-template-columns: 0.9fr 1.1fr;
  gap: 28px;
  padding: 34px;
  border: 1px solid var(--line);
  border-radius: 36px;
  background:
    linear-gradient(120deg, rgba(68, 213, 255, 0.1), transparent 42%),
    rgba(4, 15, 29, 0.76);
}

.control-copy {
  align-self: center;
}

.control-copy .eyebrow {
  margin-bottom: 26px;
}

.control-copy h3 {
  margin-top: 0;
  font-size: clamp(1.55rem, 2.45vw, 2.55rem);
  line-height: 1.14;
  letter-spacing: -0.035em;
}

.control-list {
  display: grid;
  gap: 16px;
}

.control-item {
  padding: 24px;
  border: 1px solid var(--line);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.045);
}

.control-item span {
  color: var(--cyan);
  font-weight: 800;
}

.control-item strong {
  display: block;
  margin-top: 12px;
  font-size: 20px;
}

.control-item p {
  margin: 10px 0 0;
  color: var(--muted);
  line-height: 1.75;
}

.workflow-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.workflow-grid article {
  padding: 26px;
  border-radius: 28px;
}

.workflow-grid span {
  display: grid;
  width: 54px;
  height: 54px;
  place-items: center;
  border-radius: 18px;
  color: #04111e;
  background: linear-gradient(135deg, var(--gold), var(--green));
  font-size: 24px;
  font-weight: 900;
}

.workflow-grid h3 {
  margin: 24px 0 10px;
  font-size: 22px;
}

.pricing-grid {
  display: grid;
  grid-template-columns: 1.3fr repeat(3, 1fr);
  gap: 18px;
}

.price-card {
  min-height: 290px;
  padding: 28px;
  border-radius: 30px;
}

.price-highlight {
  background:
    linear-gradient(140deg, rgba(125, 255, 204, 0.16), rgba(68, 213, 255, 0.08)),
    var(--panel-strong);
  border-color: rgba(125, 255, 204, 0.32);
}

.price-card strong {
  display: block;
  margin: 18px 0 14px;
  font-size: clamp(2.2rem, 4vw, 4.4rem);
  line-height: 1;
  letter-spacing: -0.06em;
}

.faq-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.faq {
  padding-top: 124px;
}

.faq article {
  padding: 26px;
  border-radius: 26px;
}

.faq h3 {
  margin: 0 0 10px;
  font-size: 21px;
}

.faq p {
  margin: 0;
}

.final-cta {
  width: min(1180px, calc(100% - 36px));
  margin: 96px auto 0;
  padding: 42px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 26px;
  border: 1px solid rgba(125, 255, 204, 0.28);
  border-radius: 36px;
  background:
    radial-gradient(circle at 15% 20%, rgba(125, 255, 204, 0.16), transparent 28%),
    linear-gradient(135deg, rgba(8, 24, 43, 0.96), rgba(5, 17, 31, 0.92));
}

.final-cta h2 {
  max-width: 780px;
}

.site-footer {
  width: min(1180px, calc(100% - 36px));
  margin: 38px auto 0;
  padding: 34px 0 42px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 24px;
  align-items: center;
  color: var(--dim);
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.footer-brand img {
  width: 54px;
  height: 54px;
  object-fit: contain;
  display: block;
}

.footer-brand strong {
  color: var(--text);
  letter-spacing: 0.08em;
}

.footer-brand p {
  margin: 4px 0 0;
}

.footer-links {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 18px;
}

.copyright {
  grid-column: 1 / -1;
  margin: 0;
  padding-top: 22px;
  border-top: 1px solid var(--line);
  color: var(--dim);
}

@media (max-width: 1040px) {
  .site-nav {
    grid-template-columns: 1fr auto;
  }

  .menu-toggle {
    display: grid;
  }

  .nav-links {
    display: none;
  }

  .nav-actions {
    display: none;
  }

  .site-nav.is-menu-open .nav-links,
  .site-nav.is-menu-open .nav-actions {
    grid-column: 1 / -1;
    display: grid;
    width: 100%;
    border-radius: 22px;
    background: rgba(4, 15, 29, 0.9);
    backdrop-filter: blur(18px);
  }

  .site-nav.is-menu-open .nav-links {
    justify-self: stretch;
    gap: 4px;
    margin-top: 6px;
    padding: 10px;
    border: 1px solid var(--line);
  }

  .site-nav.is-menu-open .nav-links a {
    padding: 13px 14px;
    border-radius: 14px;
  }

  .site-nav.is-menu-open .nav-links a:hover {
    background: rgba(68, 213, 255, 0.08);
  }

  .site-nav.is-menu-open .nav-actions {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
    padding: 0;
    background: transparent;
    backdrop-filter: none;
  }

  .hero,
  .control-panel {
    grid-template-columns: 1fr;
  }

  .hero {
    padding-top: 42px;
  }

  .hero-visual {
    min-height: 520px;
  }

  .pricing-grid,
  .solution-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .site-nav,
  .hero,
  .ticker,
  .section,
  .final-cta,
  .site-footer {
    width: min(100% - 28px, 1180px);
  }

  .site-nav {
    grid-template-columns: 1fr auto;
  }

  .site-nav.is-menu-open .nav-actions {
    width: 100%;
  }

  .site-nav.is-menu-open .nav-actions {
    grid-template-columns: 1fr;
  }

  .brand-copy span {
    display: none;
  }

  .hero h1 {
    font-size: clamp(2.8rem, 16vw, 4.8rem);
  }

  .hero-metrics,
  .solution-grid,
  .workflow-grid,
  .pricing-grid,
  .faq-grid {
    grid-template-columns: 1fr;
  }

  .solution-large {
    grid-column: auto;
  }

  .hero-visual {
    min-height: 460px;
  }

  .hero-visual img {
    width: 112%;
    transform: translateX(0) scale(1);
  }

  .card-orbit {
    width: 155px;
    padding: 14px;
  }

  .card-orbit strong {
    font-size: 22px;
  }

  .card-orbit-one {
    left: 0;
    top: 34px;
  }

  .card-orbit-two {
    right: 0;
    top: 112px;
  }

  .settlement-panel {
    left: 8px;
    right: 8px;
    bottom: 18px;
  }

  .control-panel,
  .final-cta {
    padding: 24px;
    border-radius: 28px;
  }

  .final-cta {
    align-items: flex-start;
    flex-direction: column;
  }

  .site-footer {
    grid-template-columns: 1fr;
  }

  .footer-links {
    justify-content: flex-start;
  }
}
</style>
