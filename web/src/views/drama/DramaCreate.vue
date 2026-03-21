<template>
  <!-- Drama Create Page / 创建短剧页面 -->
  <div class="page-container">
    <div class="content-wrapper animate-fade-in">
      <!-- Header / 头部 -->
      <AppHeader :fixed="false" :show-logo="false">
        <template #left>
          <el-button text @click="goBack" class="back-btn">
            <el-icon><ArrowLeft /></el-icon>
            <span>返回</span>
          </el-button>
          <div class="page-title">
            <h1>创建新项目</h1>
            <span class="subtitle">填写基本信息来创建你的短剧项目</span>
          </div>
        </template>
      </AppHeader>

      <!-- Form Card / 表单卡片 -->
      <div class="form-card">

        <el-form 
          ref="formRef" 
          :model="form" 
          :rules="rules" 
          label-position="top"
          class="create-form"
          @submit.prevent="handleSubmit"
        >
          <el-form-item label="项目标题" prop="title" required>
            <el-input 
              v-model="form.title" 
              placeholder="给你的短剧起个名字"
              size="large"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="项目描述" prop="description">
            <el-input 
              v-model="form.description" 
              type="textarea" 
              :rows="5"
              placeholder="简要描述你的短剧内容、风格或创意（可选）"
              maxlength="500"
              show-word-limit
              resize="none"
            />
          </el-form-item>

          <el-form-item label="Phong cách" prop="style" required>
            <el-select
              v-model="form.style"
              placeholder="Chọn phong cách"
              size="large"
              style="width: 100%"
            >
              <el-option label="Studio Ghibli" value="ghibli" />
              <el-option label="Chinese Anime" value="guoman" />
              <el-option label="Wasteland" value="wasteland" />
              <el-option label="Nostalgia" value="nostalgia" />
              <el-option label="Pixel Art" value="pixel" />
              <el-option label="Voxel" value="voxel" />
              <el-option label="Urban" value="urban" />
              <el-option label="Chinese 3D" value="guoman3d" />
              <el-option label="Chibi 3D" value="chibi3d" />
              <el-option label="Custom" value="custom" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="form.style === 'custom'" label="Mô tả phong cách" prop="custom_style" required>
            <el-input
              v-model="form.custom_style"
              type="textarea"
              :rows="3"
              placeholder="Mô tả chi tiết phong cách bạn muốn, VD: Kurzgesagt flat vector illustration..."
              maxlength="200"
              show-word-limit
              resize="none"
            />
          </el-form-item>

          <el-form-item label="Prompt Template">
            <el-select
              v-model="form.prompt_template_id"
              placeholder="Mặc định hệ thống"
              clearable
              size="large"
              style="width: 100%"
              @change="onTemplateChange"
            >
              <el-option
                v-for="tpl in promptTemplates"
                :key="tpl.id"
                :label="tpl.name"
                :value="tpl.id"
              >
                <span>{{ tpl.name }}</span>
                <span v-if="tpl.description" style="color: #999; font-size: 12px; margin-left: 8px;">{{ tpl.description }}</span>
              </el-option>
            </el-select>
            <div class="form-tip">
              <router-link to="/settings/prompt-templates">Quản lý Templates →</router-link>
            </div>
            <div v-if="templateHasStyle" class="template-style-hint">
              <el-icon><WarningFilled /></el-icon>
              <span>Template này có <strong>style_prompt</strong> riêng — sẽ override phong cách đã chọn ở trên khi sinh ảnh.</span>
            </div>
          </el-form-item>

          <div class="form-actions">
            <el-button size="large" @click="goBack">取消</el-button>
            <el-button 
              type="primary" 
              size="large"
              :loading="loading"
              @click="handleSubmit"
            >
              <el-icon v-if="!loading"><Plus /></el-icon>
              创建项目
            </el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft, Plus, WarningFilled } from '@element-plus/icons-vue'
import { dramaAPI } from '@/api/drama'
import { promptTemplateAPI } from '@/api/prompt-template'
import type { CreateDramaRequest } from '@/types/drama'
import type { PromptTemplate } from '@/types/prompt-template'
import { AppHeader } from '@/components/common'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const promptTemplates = ref<PromptTemplate[]>([])

const form = reactive<CreateDramaRequest>({
  title: '',
  description: '',
  style: 'ghibli'
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入项目标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  style: [
    { required: true, message: '请选择风格', trigger: 'change' }
  ],
  custom_style: [
    {
      validator: (_rule, value, callback) => {
        if (form.style === 'custom' && !value) {
          callback(new Error('请输入自定义风格描述'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// Template conflict detection
const selectedTemplate = computed(() => {
  if (!form.prompt_template_id) return null
  return promptTemplates.value.find(t => t.id === form.prompt_template_id) || null
})

const templateHasStyle = computed(() => {
  const tpl = selectedTemplate.value
  return tpl?.prompts?.style_prompt && tpl.prompts.style_prompt.trim().length > 0
})

const onTemplateChange = (templateId: number | undefined) => {
  if (!templateId) return
  const tpl = promptTemplates.value.find(t => t.id === templateId)
  if (tpl?.prompts?.style_prompt && tpl.prompts.style_prompt.trim()) {
    if (form.style !== 'custom') {
      ElMessage.info('Template có style riêng, sẽ override phong cách đã chọn khi sinh ảnh.')
    }
  }
}

// Submit form / 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const drama = await dramaAPI.create(form)
        ElMessage.success('创建成功')
        router.push(`/dramas/${drama.id}`)
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// Go back / 返回上一页
const goBack = () => {
  router.back()
}

// Load prompt templates
const loadPromptTemplates = async () => {
  try {
    const res = await promptTemplateAPI.list()
    promptTemplates.value = Array.isArray(res) ? res : []
  } catch {
    // silent - dropdown will just be empty
  }
}

onMounted(loadPromptTemplates)
</script>

<style scoped>
/* ========================================
   Page Layout / 页面布局 - 紧凑边距
   ======================================== */
.page-container {
  min-height: 100vh;
  background-color: var(--bg-primary);
  padding: var(--space-2) var(--space-3);
  transition: background-color var(--transition-normal);
}

@media (min-width: 768px) {
  .page-container {
    padding: var(--space-3) var(--space-4);
  }
}

.content-wrapper {
  max-width: 640px;
  margin: 0 auto;
}

/* ========================================
   Form Card / 表单卡片
   ======================================== */
.form-card {
  background: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-card);
}

/* ========================================
   Form Styles / 表单样式 - 紧凑内边距
   ======================================== */
.create-form {
  padding: var(--space-4);
}

.create-form :deep(.el-form-item) {
  margin-bottom: var(--space-4);
}

/* ========================================
   Form Actions / 表单操作区
   ======================================== */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-primary);
  margin-top: var(--space-2);
}

.form-actions .el-button {
  min-width: 100px;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
}
.form-tip a {
  color: var(--el-color-primary);
  text-decoration: none;
}

.template-style-hint {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-top: 8px;
  padding: 8px 12px;
  background: var(--el-color-warning-light-9);
  border: 1px solid var(--el-color-warning-light-5);
  border-radius: 6px;
  font-size: 12px;
  color: var(--el-color-warning-dark-2);
  line-height: 1.5;
}
.template-style-hint .el-icon {
  margin-top: 2px;
  flex-shrink: 0;
}
</style>
