<template>
<!-- Header -->
<header>
  <nav class="navbar navbar-expand-lg sticky-top">
    <div class="container">
      <a class="navbar-brand d-flex align-items-center fw-bold" href="#">
        <!-- <span class="icon-circle me-2"><i class="bi bi-credit-card-2-front"></i></span> PIXEL CARD -->
        <div class="brand"><img src="@/assets/logo-black-white.png" :alt="$t('lang.homepage.common.logo')" /></div>
      </a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#nav">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="nav">
        <ul class="navbar-nav mx-auto">
          <li class="nav-item">
            <a class="nav-link" href="#features">{{ $t('lang.homepage.nav.advertising') }}</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#platform">{{ $t('lang.homepage.nav.platform') }}</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#support">{{ $t('lang.homepage.nav.contact_us') }}</a>
          </li>
        </ul>
        <div class="d-flex gap-2 links-container">
          <nav class="links">
            <!-- 使用 $t 翻译链接文本 -->
            <!-- <a href="#services">{{ $t('lang.homepage.header.services') }}</a> -->
            <!-- btn-ghost 类的按钮文本 -->
            <el-dropdown popper-class="px-2">
              <a href="#" class="link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover" style="line-height: normal;">
                {{ langStore.currentLocale ==='zh'?'中文':langStore.currentLocale ==='en'?'English':'Tiếng Việt' }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </a>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLangSet('zh')" class="">中文</el-dropdown-item>
                  <el-dropdown-item @click="handleLangSet('en')" class="">English</el-dropdown-item>
                  <el-dropdown-item @click="handleLangSet('vn')" class="">Tiếng Việt</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <a class="btn-ghost btn-white login-ghost" href="/login">{{ $t('lang.homepage.header.signIn') }}</a>
            <a class="btn-ghost primary-ghost" href="/login">{{ $t('lang.homepage.get_started') }}</a>
            <el-dropdown class="d-none" popper-class="px-2">
              <a href="#" class="link-success link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover" style="line-height: normal;">
                {{ langStore.currentLocale ==='zh'?'中文':langStore.currentLocale ==='en'?'English':'Tiếng Việt' }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </a>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLangSet('zh')" class="">中文</el-dropdown-item>
                  <el-dropdown-item @click="handleLangSet('en')" class="">English</el-dropdown-item>
                  <el-dropdown-item @click="handleLangSet('vn')" class="">Tiếng Việt</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

          </nav>
        </div>
      </div>
    </div>
  </nav>
</header> 

<!-- HERO -->
<section class="hero">
  <div class="container">
    <div class="row align-items-center gy-4">
      <div class="col-lg-6">
        <h1 class="display-5" style="font-weight: 500;">{{ $t('lang.homepage.reliable_virtual_cards') }}</h1>
        <p class="mb-3 mt-3 fs-5">{{ $t('lang.homepage.instant_issuance') }}</p>
        <div class="d-flex gap-3">
          <a class="btn btn-primary btn-lg mr-2" href="/login">{{ $t('lang.homepage.get_started') }}</a>
          <a class="btn btn-outline-primary btn-lg" href="#support">{{ $t('lang.homepage.contact_us') }}</a>
        </div>
      </div>
      <div class="col-lg-6">
        <div class="mock rounded-24 shadow-soft p-5 text-center">
          <img class="img-fluid" :alt="$t('lang.homepage.cards_mockup')" src="@/assets/homepage/banner-cards.png">
        </div>
      </div>
    </div>

    <div class="row g-3 mt-5">
      <div class="col-md-4">
        <div class="stat-card" style="text-align: center;">
          <div class="align-items-center mb-2">
            <img src="@/assets/homepage/section-card-2.png" class="me-2" alt="">
            <p class="fs-6 pd-1">{{ $t('lang.homepage.trusted_by') }}</p>
          </div>
          <h4 class="mb-0 fs-5 fw-bold">{{ $t('lang.homepage.users_count', { count: 15000 }) }}</h4>
        </div>
      </div>
      <div class="col-md-4">
        <div class="stat-card">
          <div class="align-items-center mb-2">
            <img src="@/assets/homepage/section-card-1.png" class="me-2" alt="">
            <p class="fs-6 pd-1">{{ $t('lang.homepage.optimized_for') }}</p>
          </div>
          <h4 class="mb-0 fs-5 fw-bold">{{ $t('lang.homepage.ad_platforms') }}</h4>
        </div>
      </div>
      <div class="col-md-4">
        <div class="stat-card">
          <div class="align-items-center mb-2">
            <img src="@/assets/homepage/section-card-3.png" class="me-2" alt="">
            <p class="fs-6 pd-1">{{ $t('lang.homepage.issued') }}</p>
          </div>
          <h4 class="mb-0">{{ $t('lang.homepage.cards_issued', { count: 1000000 }) }}</h4>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- TRUSTED PLATFORM -->
<section class="py-5" id="platform" style="background-color: #fff;">
  <div class="container">
    <!-- 主标题 -->
    <div class="text-center mb-4 fs-1">
      <h3 class="fs-1 fw-6">{{ $t('lang.homepage.trusted_platform_title') }}</h3>
    </div>

    <!-- 特性 1: Full control -->
    <div class="row g-4 align-items-stretch">
      <div class="col-lg-12 col-md-12">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0">
            <div class="row no-gutters h-100">
              <div class="col-lg-6 col-md-6 left p-5">
                <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.full_control_title') }}</h3>
                <p class="fs-4 mb-2">{{ $t('lang.homepage.full_control_desc1') }}</p>
                <p class="fs-4 mb-3">{{ $t('lang.homepage.full_control_desc2') }}</p>
                <button style="max-width: 165px;" class="btn btn-primary btn-lg mr-2" @click="$router.push('/login')">
                  {{ $t('lang.homepage.get_started') }}
                </button>
              </div>
              <div class="col-lg-6 col-md-6">
                <div class="feature-card h-100">
                  <div class="card-body pd-0">
                    <div class="mock rounded-20 mb-3 text-center" style="padding: 40px 0;">
                      <img
                        src="@/assets/homepage/section-analytic.png"
                        class="img-fluid rounded-20"
                        :alt="$t('lang.homepage.alt_analytic')"
                      >
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 特性 2: No hidden fees -->
    <div class="row g-4 align-items-stretch mt-4">
      <div class="col-lg-12">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0">
            <div class="row no-gutters h-100">
              <div class="col-md-6 col-lg-6 p-5 left">
                <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.no_fees_title') }}</h3>
                <p class="fs-4 mb-3">{{ $t('lang.homepage.no_fees_desc') }}</p>
                <button style="max-width: 165px;" class="btn btn-primary btn-lg mr-2" @click="$router.push('/login')">
                  {{ $t('lang.homepage.get_started') }}
                </button>
              </div>
              <div class="col-md-6 col-lg-6" style="padding: 20px 20px; overflow: hidden;">
                <div class="mock rounded-20 mb-3 text-center" style="margin-bottom:0px">
                  <img
                    src="@/assets/homepage/section-topup.png"
                    class="img-fluid rounded-20"
                    :alt="$t('lang.homepage.alt_topup')"
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 特性 3 & 4: Best for retail & Unlimited spending -->
    <div class="row g-5 align-items-stretch mt-4">
      <!-- Best for retail -->
      <div class="col-lg-6 col-md-6">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0">
            <div class="col-md-12 col-lg-12 p-5 left">
              <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: -30px;">
                <img
                  src="@/assets/homepage/section-retail.png"
                  class="img-fluid rounded-20"
                  style="max-width: 400px; width: 100%;"
                  :alt="$t('lang.homepage.alt_retail')"
                >
              </div>
              <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.retail_title') }}</h3>
              <p class="fs-4 mb-3">{{ $t('lang.homepage.retail_desc') }}</p>
            </div>
          </div>
        </div>
      </div>
      <!-- Unlimited spending -->
      <div class="col-lg-6 col-md-6">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0 d-flex">
            <div class="col-md-12 col-lg-12" style="padding: 20px 20px; overflow: hidden;">
              <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: -30px;">
                <img
                  src="@/assets/homepage/section-spending.png"
                  class="img-fluid rounded-20"
                  style="max-width: 400px;"
                  :alt="$t('lang.homepage.alt_spending')"
                >
              </div>
              <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.unlimited_title') }}</h3>
              <p class="fs-4 mb-3">{{ $t('lang.homepage.unlimited_desc') }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 特性 5: Single balance -->
    <div class="row g-4 align-items-stretch mt-4">
      <div class="col-lg-12">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0 d-flex">
            <div class="row no-gutters h-100">
              <div class="col-md-6 col-lg-6 p-5 left">
                <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.balance_title') }}</h3>
                <p class="fs-4 mb-3">{{ $t('lang.homepage.balance_desc') }}</p>
                <button class="btn btn-primary btn-lg mr-2" @click="$router.push('/login')">
                  {{ $t('lang.homepage.get_started') }}
                </button>
              </div>
              <div class="col-md-6 col-lg-6" style="padding: 20px 20px; overflow: hidden;">
                <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: -30px;">
                  <img
                    src="@/assets/homepage/section-balance.png"
                    class="img-fluid rounded-20"
                    :alt="$t('lang.homepage.alt_balance')"
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Optimized for advertising -->
<!-- Optimized for advertising -->
<section class="py-5" id="features">
  <div class="container">
    <div class="text-center mb-4 fs-1">
      <h3 class="fs-1 fw-6">{{ $t('lang.homepage.features.optimized_advertising_title') }}</h3>
      <p>{{ $t('lang.homepage.features.optimized_advertising_desc') }}</p>
    </div>
    <div class="row g-5 align-items-stretch mt-4">
      <div class="col-lg-6 col-md-6">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0 d-flex">
            <div class="col-md-12 col-lg-12 p-5 left">
              <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: 0px;">
                <img src="@/assets/homepage/section-media.png" class="img-fluid rounded-20" style="max-width: 400px; width: 100%;" :alt="$t('lang.homepage.features.alt_media')">
              </div>
              <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.features.cards_for_ads_title') }}</h3>
              <p class="fs-4 mb-3">{{ $t('lang.homepage.features.cards_for_ads_desc') }}</p>
            </div>
          </div>
        </div>
      </div>
      <div class="col-lg-6 col-md-6">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0 d-flex">
            <div class="col-md-12 col-lg-12" style="padding: 20px 20px; overflow: hidden;">
              <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: 0px;">
                <img src="@/assets/homepage/section-channelpay.png" class="img-fluid rounded-20" style="max-width: 400px;" :alt="$t('lang.homepage.features.alt_channelpay')">
              </div>
              <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.features.high_attachment_rate_title') }}</h3>
              <p class="fs-4 mb-3">{{ $t('lang.homepage.features.high_attachment_rate_desc') }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row g-4 align-items-stretch mt-4">
      <div class="col-lg-12">
        <div class="feature-card h-100 pd-0">
          <div class="card-body pd-0 d-flex">
            <div class="row no-gutters h-100">
              <div class="col-md-6 col-lg-6 p-5 left">
                <h3 class="fs-1 fw-bold mb-4">{{ $t('lang.homepage.features.work_with_cards_title') }}</h3>
                <p class="fs-4 mb-3">{{ $t('lang.homepage.features.work_with_cards_desc') }}</p>
              </div>
              <div class="col-md-6 col-lg-6" style="padding: 20px 20px; overflow: hidden;">
                <div class="mock rounded-20 mb-3 text-center" style="margin-bottom: -30px;">
                  <img src="@/assets/homepage/section-setup.png" class="img-fluid rounded-20" :alt="$t('lang.homepage.features.alt_setup')">
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- All-in-one platform -->
<section class="py-5 all-in-one">
  <div class="container">
    <div class="text-center mb-4">
      <h3 class="fs-1 fw-6">{{ $t('lang.homepage.all_in_one.title') }}</h3>
      <p class="muted fs-5">{{ $t('lang.homepage.all_in_one.desc') }}</p>
    </div>
    <div class="row g-4">
      <div class="col-6 col-md-6 col-lg-3">
        <div class="service-card h-100 text-center ">
          <img class="service-card__icon" src="@/assets/homepage/icon-Readyforuse.png" :alt="$t('lang.homepage.all_in_one.ready_for_use')" />
          <h6 class="service-card__title">{{ $t('lang.homepage.all_in_one.ready_for_use') }}</h6>
          <p class="small muted mb-0">{{ $t('lang.homepage.all_in_one.instant_issuance') }}</p>
        </div>
      </div>
      <div class="col-6 col-md-6 col-lg-3">
        <div class="service-card h-100 text-center">
          <img class="service-card__icon" src="@/assets/homepage/icon-premium-support.png" :alt="$t('lang.homepage.all_in_one.premium_support')" />
          <h6 class="service-card__title">{{ $t('lang.homepage.all_in_one.premium_support') }}</h6>
          <p class="small muted mb-0">{{ $t('lang.homepage.all_in_one.human_assistance') }}</p>
        </div>
      </div>
      <div class="col-6 col-md-6 col-lg-3">
        <div class="service-card h-100 text-center">
          <img class="service-card__icon" src="@/assets/homepage/icon-teamwork-and-management.png" :alt="$t('lang.homepage.all_in_one.teamwork_roles')" />
          <h6 class="service-card__title">{{ $t('lang.homepage.all_in_one.teamwork_roles') }}</h6>
          <p class="small muted mb-0">{{ $t('lang.homepage.all_in_one.invite_teammates') }}</p>
        </div>
      </div>
      <div class="col-6 col-md-6 col-lg-3">
        <div class="service-card h-100 text-center">
          <img class="service-card__icon" src="@/assets/homepage/icon-detailed-analytics.png" :alt="$t('lang.homepage.all_in_one.detailed_analytics')" />
          <h6 class="service-card__title">{{ $t('lang.homepage.all_in_one.detailed_analytics') }}</h6>
          <p class="small muted mb-0">{{ $t('lang.homepage.all_in_one.track_spending') }}</p>
        </div>
      </div>
    </div>
    <div class="d-flex flex-wrap gap-3 justify-content-center mt-4">
      <!-- <a class="btn btn-primary btn-lg" href="/login">{{ $t('lang.homepage.get_started') }}</a> -->
      <a class="btn btn-outline-primary btn-lg" href="#support">{{ $t('lang.homepage.nav.contact_us') }}</a>
    </div>
  </div>
</section>

<!-- Go Big band -->
<section class="py-5 go-big">
  <div class="container">
    <div class="go-big__panel p-4 p-lg-5 rounded-24 text-center">
      <h3 class="mb-3 fs-1 fw-6">{{ $t('lang.homepage.go_big.title') }}</h3>
      <p class="muted fs-5 mb-4">{{ $t('lang.homepage.go_big.subtitle') }}</p>
      <div class="row row-cols-1 row-cols-md-3 g-4 go-big__features">
        <div class="col">
          <div class="go-big__item h-100">
            <span class="go-big__icon">
              <i class="bi bi-braces-asterisk"></i>
            </span>
            <h6 class="go-big__item-title">{{ $t('lang.homepage.go_big.api') }}</h6>
            <p class="small muted mb-0">{{ $t('lang.homepage.go_big.integrate_into_stack') }}</p>
          </div>
        </div>
        <div class="col">
          <div class="go-big__item h-100">
            <span class="go-big__icon">
              <i class="bi bi-people"></i>
            </span>
            <h6 class="go-big__item-title">{{ $t('lang.homepage.go_big.growing_together') }}</h6>
            <p class="small muted mb-0">{{ $t('lang.homepage.go_big.scale_with_confidence') }}</p>
          </div>
        </div>
        <div class="col">
          <div class="go-big__item h-100">
            <span class="go-big__icon">
              <i class="bi bi-server"></i>
            </span>
            <h6 class="go-big__item-title">{{ $t('lang.homepage.go_big.p2p_infrastructure') }}</h6>
            <p class="small muted mb-0">{{ $t('lang.homepage.go_big.run_on_secure_rails') }}</p>
          </div>
        </div>
      </div>
      <p class="muted mt-4 mb-4">{{ $t('lang.homepage.go_big.support_line') }}</p>
      <div class="d-flex justify-content-center">
        <a href="#support" class="btn btn-outline-primary btn-lg">{{ $t('lang.homepage.nav.contact_us') }}</a>
      </div>
    </div>
  </div>
</section>

<!-- Individual & Instant Support -->
<section class="py-5 support" id="support">
  <div class="container">
    <div class="text-center"><h3 class="fs-1 mb-5 mt-5">Invidual & Instant Support</h3></div>
    <div class="row align-items-center gy-4 mt-2 support-wrapper">
      <div class="col-lg-6 order-2 order-lg-1">
        <div class="support__card fs-1 p-4 p-lg-5 rounded-24 h-100">
          <span class="badge-soft text-uppercase small fw-semibold">{{ $t('lang.homepage.support.subtitle') }}</span>
          <!-- <h4 class="mt-3 mb-3 fs-2">{{ $t('lang.homepage.support.title') }}</h4> -->
          <p class="muted mb-4 fs-4">{{ $t('lang.homepage.support.desc') }}</p>
          <a class="btn btn-telegram btn-lg" href="https://t.me/charlin0067" target="_blank" rel="noopener">
            <i class="bi bi-telegram me-2"></i>{{ $t('lang.homepage.support.cta') }}
          </a>
        </div>
      </div>
      <div class="col-lg-6 order-1 order-lg-2 text-center">
        <div class="support__visual rounded-24">
          <img src="@/assets/homepage/section-tg-chat.png" class="img-fluid" :alt="$t('lang.homepage.support.alt_image')" />
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Final CTA -->
<section class="py-5 final-cta">
  <div class="container">
    <div class="row align-items-center gy-4">
      <div class="col-lg-6">
        <h3 class="fs-1 fw-bold">{{ $t('lang.homepage.footer.trusted_cards_title') }}</h3>
        <p class="muted fs-5 mb-4">{{ $t('lang.homepage.footer.cards_subtitle') }}</p>
        <div class="d-flex flex-wrap gap-3">
          <a class="btn btn-primary btn-lg" href="/login">{{ $t('lang.homepage.footer.get_started') }}</a>
          <a class="btn btn-outline-primary btn-lg" href="#support">{{ $t('lang.homepage.nav.contact_us') }}</a>
        </div>
      </div>
      <div class="col-lg-6 text-center">
        <div class="final-cta__image">
          <img
            src="@/assets/homepage/card-footer.png"
            class="img-fluid"
            :alt="$t('lang.homepage.footer.alt_card_footer')"
          >
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Footer -->
<footer class="footer py-4">
  <div class="container footer-inner">
    <div class="footer-nav small">
      <a class="link" href="#advertising">{{ $t('lang.homepage.footer.advertising') }}</a>
      <a class="link" href="#platform">{{ $t('lang.homepage.footer.platform') }}</a>
      <a class="link" href="#support">{{ $t('lang.homepage.footer.contact_us') }}</a>
    </div>
    <div class="footer-divider"></div>
    <div class="footer-copy small">
      {{ $t('lang.homepage.footer.copyright') }}
    </div>
  </div>
</footer>



</template>
<script setup>
import { useLanguageStore } from '@/pinia/modules/language';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const langStore = useLanguageStore()
const handleLangSet = (lang) =>{
  langStore.changeLocale(lang)
}
</script>

<style lang="scss" scoped >
body, h1, h2, h3, p, span, div {
  font-family: Arial, sans-serif;
}
:root{
  --pc-green:#37c27f;
  --pc-dark:#0f1f1a;
  --pc-mint:#e8f7ef;
  --pc-soft:#f5faf7;
  --pc-gray:#6b7280;
}
#app{overflow: auto;}
body{color:#0b1720;background:linear-gradient(180deg,#f3fbf6 0%,#ffffff 30%,#f6faf7 100%);}
.navbar{backdrop-filter:blur(6px);box-shadow:0 2px 10px rgba(16,24,40,.06);}


.badge-soft{background:var(--pc-mint);color:#0a5a3b;}
.hero{
  background: linear-gradient(to bottom, #D8F2E5 0%, #fff 40%, white 100%);
  padding: 72px 0 40px;
}

header{background: linear-gradient( 233deg, #00AD5A 0%, #00572D 100%);}
.brand img{height:50px; padding: 7px;}

.links { display: flex; gap: 30px; align-items: center; }
.links a { color: white; text-decoration: none; font-weight: 500; }

@media (max-width: 1024px) {
  .navbar-nav { text-align: center;}
  .links { flex-direction: column; gap: 10px; justify-content: center;
    .btn-ghost { width: 100%; text-align: center; }
  }
  .links-container { justify-content: center;}
  .navbar-collapse {
    margin: 0;
    padding: 20px 0;
    flex-direction: column;
    width: 100%;
    a {
      color: black !important;
    }
    .btn-ghost {
      width: 100%;
      height: 44px;
      border-radius: 22px;
      font-weight: 600;
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }
    .btn-ghost.login-ghost {
      background-color: transparent;
      color: #fff;
      border: 2px solid rgba(255,255,255,0.6);
    }
    .btn-ghost.primary-ghost {
      background: linear-gradient(256deg,#0c5d30 0%,#0a3f24 100%);
      color: #fff;
      border: none;
      box-shadow: 0 10px 20px rgba(9,49,25,0.26);
    }
  }
  
}
.btn-ghost{
  background-color: #163D03;
  height: 34px;
  border-radius: 17px;
  padding: 5px 25px;
  transition: all .2s ease;
}
.btn-white{
  background-color: #00A656;
  color: black;
  border-color: #fff;
}
.btn-ghost.login-ghost{
  background:transparent;
  color:#fff;
  border:2px solid rgba(255,255,255,0.6);
}
.btn-ghost.primary-ghost{
  background:linear-gradient(256deg,#0c5d30 0%,#0a3f24 100%);
  color:#fff;
  border:none;
  box-shadow:0 8px 18px rgba(8,47,25,0.25);
}
.btn-primary {
  background: linear-gradient( 256deg, #00AD5A 0%, #02783E 49%, #016334 100%);
  height: 44px;
  border-radius: 22px;
  border: none;
}
.btn-outline-primary {
  background: #fff;
  color: #00AD5A;
  width: 200px;
  max-width: 100%;
  height: 44px;
  border-radius: 22px;
  border: 1px solid #00AD5A;
}
.hero h1{font-weight:800;letter-spacing:.2px;}
.stat-card{
  border:1px solid #e7efe9;background:#fff;border-radius:16px;padding:18px 16px;box-shadow:0 8px 24px rgba(16,24,40,.06);
  height:100%; text-align: center; display: flex; flex-direction: column; justify-content: center;position: relative;
  background-image: url('@/assets/homepage/section-card-bg.png'); background-repeat: no-repeat; background-position: center bottom; background-size: contain;
}
.stat-card span{width: 100%;}
.stat-card img{width: 80%; object-fit: contain; margin-top: -50px; padding-bottom: 20px;}
.section-head{margin-bottom:16px}
.feature-card{
  background:#E7F6EF;border-radius:18px;
}
.feature-card .card-body {padding:22px;}
.feature-card .card-body img{width:100%;height:auto;object-fit:fit;}
.feature-card .card-body .left{vertical-align: middle;display: flex;flex-direction: column;justify-content: center;}
.icon-circle{
  width:40px;height:40px;border-radius:10px;display:inline-flex;align-items:center;justify-content:center;
  background:var(--pc-mint);color:#0a5a3b;font-size:20px;
}
.muted{color:var(--pc-gray);}
.soft-strip{background:var(--pc-soft);border:1px solid #e5efe9;border-radius:18px;}
.service-card{
  background:#ffffff;
  border-radius:18px;
  border:1px solid #dbeee4;
  box-shadow:0 18px 40px rgba(16,24,40,0.08);
  padding:30px 20px;
  transition: transform .2s ease, box-shadow .2s ease;
  background-image: url('@/assets/homepage/section-card-bg.png'); background-repeat: no-repeat; background-position: center bottom; background-size: contain;
}
.service-card:hover{
  transform: translateY(-4px);
  box-shadow:0 22px 48px rgba(16,24,40,0.12);
}
.service-card__icon{
  width:68px;
  height:auto;
  margin-bottom:16px;
}
.service-card__title{
  font-weight:600;
  margin-bottom:8px;
}
.go-big{
  background-image: url('@/assets/homepage/section-go-big-with-pixelcard-bg.png');
  background-repeat: no-repeat;
  background-size: cover;
}

@media (max-width: 768px) {
  .go-big{
    background: none;
    background-color: #E7F6Ef;
  }
}
// .go-big__panel{
//   // background: rgba(255,255,255,0.95);
//   // border:1px solid #d5ebdf;
//   // box-shadow:0 20px 50px rgba(15,31,26,0.08);
// }
.go-big__features{
  align-items: stretch;
}
.go-big__icon{
  width:64px;
  height:64px;
  border-radius:18px;
  display:inline-flex;
  align-items:center;
  justify-content:center;
  background:linear-gradient(180deg,#00ad5a 0%,#02924d 100%);
  color:#fff;
  font-size:28px;
  margin-bottom:18px;
}
.go-big__item{
  padding:0 12px;
}
.go-big__item-title{
  font-weight:600;
  margin-bottom:8px;
  color:#0b1720;
}
.support{background: #fff;}
.support__card{
  background:none;
}
.support-wrapper {
  background-color: #E7F6EF;
  border-radius: 24px;
}
.btn-telegram{
  background: linear-gradient(256deg, #00AD5A 0%, #02783E 49%, #016334 100%);
  color:#fff;
  border:none;
  border-radius:22px;
  padding:10px 32px;
  display:inline-flex;
  align-items:center;
  gap:10px;
}
.btn-telegram:hover{
  color:#fff;
  opacity:0.9;
}
.support__visual{
  // padding:32px;
  display:inline-block;
}
.final-cta{background: linear-gradient(180deg,#ffffff 0%,#dff3e8 100%);}
.final-cta__image{
  display:inline-block;
}
.final-cta__image img{
  max-width:520px;
}
.rounded-20{border-radius:20px;}
.rounded-24{border-radius:24px;}
.shadow-soft{box-shadow:0 12px 30px rgba(16,24,40,.08);}
/* 图片占位 */
.mock{ border: none; box-shadow: none;}
footer {
  border:none;
  padding: 0;
  background-color: #205306;
  border-top:1px solid #d4ebde;
}
.footer-inner{
  width:100%;
  padding:24px 0;
  display:flex;
  flex-direction:column;
  align-items:center;
  gap:16px;
  text-align:center;
}
.footer-nav{
  display:flex;
  gap:32px;
  justify-content:center;
}
.footer-nav .link{
  color:#f6fdf9;
  text-decoration:none;
}
.footer-divider{
  width:100%;
  max-width:320px;
  height:1px;
  background:rgba(255,255,255,0.8);
}
.footer-copy{
  color:#f6fdf9;
}

@media (max-width: 768px) {
  .btn-primary,
  .btn-outline-primary {
    width: 100%;
  }
  .service-card{
    padding:24px 16px;
  }
  .support__visual{
    padding:20px;
  }
  .final-cta__image{
    padding:16px;
  }
  .final-cta__image img{
    max-width:100%;
  }
}

@media (min-width: 992px) {
  .footer-inner{
    flex-direction:row;
    justify-content:space-between;
    align-items:center;
    text-align:left;
  }
  .footer-divider{
    display:none;
  }
  .footer-nav{
    gap:28px;
  }
}
</style>
