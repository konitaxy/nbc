<template>
  <div class="container form-white" style="min-height: 600px">
    <el-alert v-if="user.clientReviewStatus === 2 && dd.needEnhancedKYB" type="error" :description="dd.tip"></el-alert>

    <div class="certification-choice-container" v-if="certLevel === 1">
      <!-- 用户未提交认证 -->
      <el-card shadow="never" body-style="background-color:#FAFAFA;" class="box-card" v-if="user.clientReviewStatus !== 3">
        <div class="card-content" @click="goToCert('individual')">
          <h3>{{ $t('lang.individual_certification') }}</h3>
          <p>{{ $t('lang.select_individual_cert') }}</p>
        </div>
      </el-card>

      <el-card shadow="never" body-style="background-color:#FAFAFA;" class="box-card" v-if="user.clientReviewStatus !== 3">
        <div class="card-content" @click="goToCert('enterprise')">
          <h3>{{ $t('lang.enterprise_certification') }}</h3>
          <p>{{ $t('lang.select_enterprise_cert') }}</p>
        </div>
      </el-card>

      <!-- 用户审核中 -->
      <el-card class="box-card center-card" v-else-if="user.clientReviewStatus === 2 && !dd.needEnhancedKYB">
        <div class="card-content">
          <h3>{{ $t('lang.under_review') }}</h3>
          <p>{{ $t('lang.review_in_progress') }}</p>
        </div>
      </el-card>

      <!-- 用户已认证 -->
      <el-card class="box-card center-card" v-else-if="user.clientReviewStatus === 3">
        <div class="card-content">
          <h3>{{ $t('lang.certification_successful') }}</h3>
          <p>{{ $t('lang.certification_passed') }}</p>
        </div>
      </el-card>
    </div>
    <div>
    <div v-if="certLevel > 2">
      <h2 class="p-3">{{ $t('lang.enterprise_certification_form') }}</h2>
      <hr>
        <el-form :model="dd" label-width="140px" label-position="right">
          
          <el-row :gutter="10">
            <!-- 企业类型 -->
            <el-col :span="12">
              <el-form-item :label="$t('lang.enterprise_type')">
                <el-input v-model="dd.entEnterpriseType" />
              </el-form-item>
            </el-col>

            <!-- 企业中文名 -->
            <el-col :span="12">
              <el-form-item :label="$t('lang.enterprise_chinese_name')" :error="errors.entEnterpriseChineseName">
                <el-input v-model="dd.entEnterpriseChineseName" />
              </el-form-item>
            </el-col>

            <!-- 企业英文名 -->
            <el-col :span="12">
              <el-form-item :label="$t('lang.enterprise_english_name')">
                <el-input v-model="dd.entEnterpriseEnglishName" />
              </el-form-item>
            </el-col>

            <!-- 营业执照形式 -->
            <el-col :span="12">
              <el-form-item :label="$t('lang.business_registration_form')">
                <el-input v-model="dd.entBusinessRegistrationForm" type="textarea" :rows="3" />
              </el-form-item>
            </el-col>

            <!-- 营业执照编号 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.business_registration_no')" :error="errors.entBusinessRegistrationNo">
            <el-input v-model="dd.entBusinessRegistrationNo" />
          </el-form-item>
        </el-col>

        <!-- 营业地址证明 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.business_address_proof')">
            <el-upload
              class="avatar-uploader"
              action=""
              accept="image/*"
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleBusinessAddressProofChange"
              :before-upload="beforeUpload"
            >
              <div v-if="dd.entBusinessAddressProof" class="avatar-wrapper">
                <el-image
                  :src="dd.entBusinessAddressProof"
                  fit="scale-down"
                  class="img-thumbnail"
                  style="width: 100px; height: 100px; border-radius: 4px;"
                />
              </div>
              <el-button v-else type="primary">{{ $t('lang.upload_image') }}</el-button>
            </el-upload>
          </el-form-item>
        </el-col>

        <!-- 成立日期 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.date_of_establishment')">
            <el-date-picker v-model="dd.entDateOfEstablishment" type="date" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>

        <!-- 营业执照过期日期 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.date_of_expiration')">
            <el-date-picker v-model="dd.entDateOfExpiration" type="date" value-format="YYYY-MM-DD" />
          </el-form-item>
        </el-col>

        <!-- 本地经营场所 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.local_business_premise')" :error="errors.entLocalBusinessPremise">
            <el-input v-model="dd.entLocalBusinessPremise" />
          </el-form-item>
        </el-col>

        <!-- 所在省 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.province')">
            <el-input v-model="dd.entProvince" />
          </el-form-item>
        </el-col>

        <!-- 所在城市 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.city')">
            <el-input v-model="dd.entCity" />
          </el-form-item>
        </el-col>

        <!-- 是否上市公司 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.listed_company')">
            <el-switch v-model="dd.entListedCompany" />
          </el-form-item>
        </el-col>

        <!-- 是否国企 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.state_owned')">
            <el-switch v-model="dd.entStateOwned" />
          </el-form-item>
        </el-col>

        <!-- 是否外资企业 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.foreign_invested')">
            <el-switch v-model="dd.entForeignInvested" />
          </el-form-item>
        </el-col>

        <!-- 股东结构 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.shareholder_structure')">
            <el-input v-model="dd.entShareholderStructure" />
          </el-form-item>
        </el-col>

        <!-- 是否 B2B 企业 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.is_b2b')">
            <el-switch v-model="dd.entIsB2B" />
          </el-form-item>
        </el-col>

        <!-- 经营地址 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.operation_address')">
            <el-input v-model="dd.entOperationAddress" />
          </el-form-item>
        </el-col>

        <!-- 注册资本 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.registered_capital')">
            <el-input v-model="dd.entRegisteredCapital" />
          </el-form-item>
        </el-col>

        <!-- 拟从事行业 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.intended_business_industry')">
            <el-input v-model="dd.entIntendedBusinessIndustry" />
          </el-form-item>
            </el-col>
          </el-row>
        </el-form>

      </div>
      
      <div v-if="certLevel > 1">
        <div>
        <h2 class="p-3">{{ $t('lang.individual_certification_form') }}</h2>
        <hr>
      </div>
        <el-form :model="individual" label-width="140px" label-position="right">
          <el-row :gutter="10">

            <!-- 所属国家/地区 -->
            <el-col :span="12">
          <el-form-item :label="$t('lang.country_or_region')" :error="errors.indCountryOrRegion">
            <el-input v-model="dd.indCountryOrRegion" />
          </el-form-item>
        </el-col>

        <!-- 职位 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.position')">
            <el-input v-model="dd.indPosition" />
          </el-form-item>
        </el-col>

        <!-- 中文名 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.chinese_name')" :error="errors.indChineseName">
            <el-input v-model="dd.indChineseName" />
          </el-form-item>
        </el-col>

        <!-- 英文名 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.english_name')">
            <el-input v-model="dd.indEnglishName" />
          </el-form-item>
        </el-col>

        <!-- ID类型 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.id_type')" :error="errors.indIDType">
            <el-select v-model="dd.indIDType" clearable>
              <el-option :label="$t('lang.id_card')" value="ID card" />
              <el-option :label="$t('lang.driving_license')" value="driving license" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('lang.identification_no')" :error="errors.indIdentificationNo">
            <el-input v-model="dd.indIdentificationNo" />
          </el-form-item>
        </el-col>

        <!-- 身份证正面 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.id_front')">
            <el-upload
              class="avatar-uploader"
              action=""
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleIDFrontEndChange"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <div v-if="dd.indIDFrontEnd" class="avatar-wrapper">
                <el-image
                  :src="dd.indIDFrontEnd"
                  fit="scale-down"
                  class="img-thumbnail"
                  style="width: 100px; height: 100px; border-radius: 4px;"
                />
              </div>
              <el-button v-else type="primary">{{ $t('lang.upload_front_image') }}</el-button>
            </el-upload>
          </el-form-item>
        </el-col>

        <!-- 身份证背面 -->
        <el-col :span="12">
          <el-form-item :label="$t('lang.id_back')">
            <el-upload
              class="avatar-uploader ms-1"
              action=""
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleIDBackEndChange"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <div v-if="dd.indIDBackEnd" class="avatar-wrapper">
                <el-image
                  :src="dd.indIDBackEnd"
                  fit="scale-down"
                  class="img-thumbnail"
                  style="width: 100px; height: 100px; border-radius: 4px;"
                />
              </div>
              <el-button v-else type="primary">{{ $t('lang.upload_back_image') }}</el-button>
            </el-upload>
          </el-form-item>
            </el-col>

        
<!-- 发行日期 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.issue_date')">
    <el-date-picker v-model="dd.indIssueDate" type="date" value-format="YYYY-MM-DD" />
  </el-form-item>
</el-col>

<!-- 过期日期 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.expiration_date')">
    <el-date-picker v-model="dd.indExpirationDate" type="date" value-format="YYYY-MM-DD" />
  </el-form-item>
</el-col>

<!-- 出生日期 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.date_of_birth')">
    <el-date-picker v-model="dd.indDateOfBirth" type="date" value-format="YYYY-MM-DD" />
  </el-form-item>
</el-col>

<!-- 省/州 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.province_or_state')">
    <el-input v-model="dd.indProvinceOrState" />
  </el-form-item>
</el-col>

<!-- 城市 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.city')">
    <el-input v-model="dd.indCity" />
  </el-form-item>
</el-col>

<!-- 居住地址 -->
<el-col :span="12">
  <el-form-item :label="$t('lang.residential_address')">
    <el-input v-model="dd.indResidentialAddress" />
  </el-form-item>
</el-col>
          </el-row>
        </el-form>

      </div>
      <!-- 提交按钮 -->
      <el-col :span="24" v-if="certLevel > 1" class="d-flex justify-content-center">
        <el-button type="info" @click="certLevel = 1">{{ $t('lang.back') }}</el-button>
        <el-button :loading="loading.submitLoading" type="primary" @click="handleSubmitDD">{{ $t('lang.submit') }}</el-button>
        </el-col>
    </div>
  </div>
</template>

<script setup>
import { ref,onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import { getDueDiligence,setDueDiligence} from '@/api/profile'
import { ElMessage } from 'element-plus'
const router = useRouter()
const userStore = useUserStore()

const user = ref(userStore.userInfo)
const dd = ref(userStore.userInfo.dueDiligence)
onMounted(() => {
  getDueDiligence().then(res=>{
    dd.value = res.data
    if (userStore.userInfo.clientReviewStatus === 2){
        if(res.code === 0){
          if(dd.value.needEnhancedKYB){
            goToCert(userStore.userInfo.clientType)
          }
        }  
  }
  })
})
const type = ref('individual') 
const certLevel = ref(1)

const goToCert = (type) => {
  if (type === 'individual') {
    // router.push('/certification/individual')
    certLevel.value = 2
  } else if (type === 'enterprise') {
    certLevel.value = 3
  }
}
const loading = reactive({
  submitLoading: false
})
const handleSubmitDD = ()=>{
  var form ={
    ...dd.value,
    type: type.value,
  }
  loading.submitLoading = true
  setDueDiligence(form).then(res => {
    loading.submitLoading = false
    if (res.code === 0) {
      ElMessage.success('submit success')
    }
  })
}
const errors = reactive({
  entEnterpriseChineseName: '',
  entBusinessRegistrationNo: '',
  entLocalBusinessPremise: '',
  indCountryOrRegion: ''
})
const beforeUpload = (rawFile) => {
  const isValid = ['image/jpeg', 'image/png', 'image/gif'].includes(rawFile.type)
  if (!isValid) {
    ElMessage.error('只能上传 JPG/PNG/GIF 格式的图片!')
  }
  return isValid
}

// 当文件状态改变时的钩子，添加文件、上传成功和上传失败时都会被调用
const handleBusinessAddressProofChange = (uploadFile, uploadFiles) => {
  const file = uploadFile.raw;
  if (file) {
    if (file.size > 1024 * 1024 * 5) {
      ElMessage.error('upload file size must be less than 5M')
      return false
    }
    const reader = new FileReader();
    reader.onload = () => {
      // 将Base64字符串赋值给dd.entBusinessAddressProof
      dd.value.entBusinessAddressProof = reader.result;
    };
    reader.readAsDataURL(file);
  }
}
const removeBusinessAddressProof = () => {
  dd.value.entBusinessAddressProof = ''
}

// 身份证正面处理逻辑
const handleIDFrontEndChange = (uploadFile, uploadFiles) => {
  const file = uploadFile.raw
  if (file) {
    if (file.size > 1024 * 1024 * 5) {
      ElMessage.error('upload file size must be less than 5M')
      return false
    }
    const reader = new FileReader()
    reader.onload = () => {
      dd.value.indIDFrontEnd = reader.result // base64
    }
    reader.readAsDataURL(file)
  }
}

// 移除身份证正面
const removeIDFrontEnd = () => {
  dd.value.indIDFrontEnd = ''
}

// 身份证背面处理逻辑
const handleIDBackEndChange = (uploadFile, uploadFiles) => {
  const file = uploadFile.raw
  if (file) {
    if (file.size > 1024 * 1024 * 5) {
      ElMessage.error('upload file size must be less than 5M')
      return false
    }
    const reader = new FileReader()
    reader.onload = () => {
      dd.value.indIDBackEnd = reader.result // base64
    }
    reader.readAsDataURL(file)
  }
}
function compressImage(file, maxSize ,target) {
  const image = new Image();
  const reader = new FileReader();

  reader.onload = (e) => {
    image.src = e.target.result;
  };

  image.onload = () => {
    let canvas = document.createElement('canvas');
    let ctx = canvas.getContext('2d');

    // 设置初始画布尺寸（保持原始宽高比）
    let { width, height } = image;
    canvas.width = width;
    canvas.height = height;

    // 初始质量
    let quality = 0.9;

    // 持续压缩直到 size 合理
    function compress() {
      // 清空画布并绘制
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.drawImage(image, 0, 0, width, height);

      // 转为 blob（jpeg 格式压缩）
      canvas.toBlob(
        (blob) => {
          const blobSize = blob.size;
          console.log('压缩后 Blob 大小:', (blobSize / 1024).toFixed(2), 'KB');

          if (blobSize > maxSize && quality > 0.5) {
            // 仍太大，降低质量继续压缩
            quality -= 0.1;
            compress();
          } else {
            // 达标，转为 base64
            const reader = new FileReader();
            reader.onload = (e) => {
              const base64 = e.target.result;
              target = base64;
              console.log('身份证背面（压缩后）base64:', base64.length, '字符');
              ElMessage.success(`压缩完成，大小: ${(blobSize / 1024).toFixed(2)} KB`);
            };
            reader.readAsDataURL(blob);
          }
        },
        'image/jpeg', // 统一输出为 jpeg（压缩率高），PNG 会转为 JPG
        quality
      );
    }

    compress();
  };

  image.onerror = () => {
    ElMessage.error('图片加载失败，请检查文件是否损坏');
  };

  // 开始读取文件
  reader.readAsDataURL(file);
}

// 移除身份证背面
const removeIDBackEnd = () => {
  dd.value.indIDBackEnd = ''
}

</script>

<style scoped>
.certification-choice-container {
  display: flex;
  justify-content: center;
  padding: 40px;
  flex-wrap: wrap;
}

.box-card {
  width: 360px;
  margin: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #198754;

}

.box-card:hover {
  box-shadow: 0 4px 12px rgba(33, 150, 4, 0.1);
  border-color: #1e9a60 ;
}

.center-card {
  text-align: center;
}

.card-content {
  padding: 20px;
}
</style>