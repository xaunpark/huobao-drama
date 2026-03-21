<template>
  <div class="prompt-templates-page">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.back()" :icon="ArrowLeft" text>Quay lại</el-button>
        <h2>📋 Quản lý Prompt Templates</h2>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">
        Tạo Template mới
      </el-button>
    </div>

    <el-table :data="templates" v-loading="loading" stripe style="width: 100%">
      <el-table-column prop="name" label="Tên Template" min-width="200" />
      <el-table-column prop="description" label="Mô tả" min-width="300">
        <template #default="{ row }">
          {{ row.description || '—' }}
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="Cập nhật" width="180">
        <template #default="{ row }">
          {{ formatDate(row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column label="Thao tác" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" text @click="openEditDialog(row)">Sửa</el-button>
          <el-button size="small" type="info" text @click="handleDuplicate(row)">Nhân bản</el-button>
          <el-popconfirm title="Xóa template này?" @confirm="handleDelete(row)">
            <template #reference>
              <el-button size="small" type="danger" text>Xóa</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? 'Sửa Template' : 'Tạo Template mới'"
      width="80%"
      top="5vh"
      destroy-on-close
    >
      <el-form :model="form" label-position="top">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Tên Template" required>
              <el-input v-model="form.name" placeholder="VD: Explainer Video Finance" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Mô tả">
              <el-input v-model="form.description" placeholder="Mô tả ngắn gọn" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider>Prompt Types</el-divider>
        <p class="hint-text">
          💡 Để trống = Hệ thống sẽ dùng Prompt mặc định. Bạn chỉ cần nhập những phần muốn tùy chỉnh.
        </p>

        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane
            v-for="group in PROMPT_TYPE_GROUPS"
            :key="group.key"
            :label="group.label"
            :name="group.key"
          >
            <div v-for="pt in group.types" :key="pt.key" class="prompt-type-block">
              <div class="prompt-type-header">
                <span class="prompt-type-label">{{ pt.label }}</span>
                <el-button
                  size="small"
                  text
                  type="primary"
                  @click="loadDefault(pt.key)"
                  :loading="loadingDefault === pt.key"
                >
                  Tải Prompt mặc định
                </el-button>
              </div>
              <el-input
                v-model="(form.prompts as any)[pt.key]"
                type="textarea"
                :rows="8"
                :placeholder="'Để trống = dùng Prompt mặc định của hệ thống'"
                resize="vertical"
              />
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">Hủy</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          {{ isEditing ? 'Cập nhật' : 'Tạo' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { promptTemplateAPI } from '../../api/prompt-template'
import type { PromptTemplate, PromptTemplatePrompts } from '../../types/prompt-template'
import { PROMPT_TYPE_GROUPS } from '../../types/prompt-template'

const templates = ref<PromptTemplate[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const activeTab = ref('script')
const loadingDefault = ref<string | null>(null)
const defaultPrompts = ref<PromptTemplatePrompts | null>(null)

const form = ref<{
  name: string
  description: string
  prompts: PromptTemplatePrompts
}>({
  name: '',
  description: '',
  prompts: {}
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('vi-VN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await promptTemplateAPI.list()
    // Axios response interceptor extracts res.data.data -> res
    templates.value = Array.isArray(res) ? res : ((res as any)?.data || [])
  } catch (e: any) {
    ElMessage.error('Không thể tải danh sách templates')
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', description: '', prompts: {} }
  activeTab.value = 'script'
  dialogVisible.value = true
}

const openEditDialog = (tpl: PromptTemplate) => {
  isEditing.value = true
  editingId.value = tpl.id
  form.value = {
    name: tpl.name,
    description: tpl.description || '',
    prompts: { ...(tpl.prompts || {}) }
  }
  activeTab.value = 'script'
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.value.name.trim()) {
    ElMessage.warning('Vui lòng nhập tên Template')
    return
  }
  saving.value = true
  try {
    if (isEditing.value && editingId.value) {
      await promptTemplateAPI.update(editingId.value, {
        name: form.value.name,
        description: form.value.description,
        prompts: form.value.prompts
      })
      ElMessage.success('Cập nhật thành công')
    } else {
      await promptTemplateAPI.create({
        name: form.value.name,
        description: form.value.description,
        prompts: form.value.prompts
      })
      ElMessage.success('Tạo template thành công')
    }
    dialogVisible.value = false
    fetchTemplates()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.error || 'Có lỗi xảy ra')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (tpl: PromptTemplate) => {
  try {
    await promptTemplateAPI.delete(tpl.id)
    ElMessage.success('Đã xóa template')
    fetchTemplates()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.error || 'Không thể xóa')
  }
}

const handleDuplicate = async (tpl: PromptTemplate) => {
  try {
    await promptTemplateAPI.duplicate(tpl.id)
    ElMessage.success('Đã nhân bản template')
    fetchTemplates()
  } catch (e: any) {
    ElMessage.error('Không thể nhân bản')
  }
}

const loadDefault = async (promptType: string) => {
  loadingDefault.value = promptType
  try {
    if (!defaultPrompts.value) {
      const res = await promptTemplateAPI.getDefaults()
      defaultPrompts.value = res || {}
    }
    const defaultValue = (defaultPrompts.value as any)?.[promptType] || ''
    if (defaultValue) {
      ;(form.value.prompts as any)[promptType] = defaultValue
      ElMessage.success('Đã tải prompt mặc định')
    } else {
      ElMessage.info('Không tìm thấy prompt mặc định cho loại này')
    }
  } catch {
    ElMessage.error('Không thể tải prompt mặc định')
  } finally {
    loadingDefault.value = null
  }
}

onMounted(fetchTemplates)
</script>

<style scoped>
.prompt-templates-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.header-left h2 {
  margin: 0;
  font-size: 20px;
}
.hint-text {
  color: #909399;
  font-size: 13px;
  margin-bottom: 16px;
}
.prompt-type-block {
  margin-bottom: 24px;
}
.prompt-type-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.prompt-type-label {
  font-weight: 600;
  font-size: 14px;
}
</style>
