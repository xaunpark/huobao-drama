<template>
  <el-dialog
    v-model="visible"
    :title="$t('professionalEditor.batch.title')"
    width="900px"
    :close-on-click-modal="false"
    class="batch-dialog"
  >
    <div class="batch-container">
      <!-- 顶部配置栏 -->
      <div class="batch-config-bar">
        <el-alert
          :title="$t('professionalEditor.batch.title')"
          type="info"
          :description="$t('professionalEditor.batch.instructions')"
          show-icon
          :closable="false"
          style="margin-bottom: 20px"
        />

        <div class="config-row">
          <div class="config-item">
            <span class="label">{{ $t('professionalEditor.batch.imageModel') }}</span>
            <el-tag size="small" type="info">{{ $t('professionalEditor.batch.defaultModel') }}</el-tag>
          </div>
          <div class="config-item">
            <span class="label">{{ $t('professionalEditor.batch.videoModel') }}</span>
            <el-select v-model="selectedVideoModel" :placeholder="$t('video.selectVideoModel')" size="small" style="width: 200px">
              <el-option
                v-for="model in videoModels"
                :key="model.id"
                :label="model.name"
                :value="model.id"
              />
            </el-select>
          </div>
          <div class="config-item">
            <span class="label">{{ $t('professionalEditor.batch.generationMode') }}</span>
            <el-select v-model="generationMode" size="small" style="width: 150px">
              <el-option :label="$t('professionalEditor.batch.keyframeMode')" value="key" />
              <el-option :label="$t('professionalEditor.batch.r2vMode')" value="action" />
            </el-select>
          </div>
        </div>

        <div class="action-row" style="margin-top: 15px">
          <el-button type="primary" :loading="isBatching" @click="startFullBatch">
            <el-icon><MagicStick /></el-icon> {{ $t('professionalEditor.batch.runAll') }}
          </el-button>
          <el-button-group>
            <el-button :disabled="isBatching" @click="startStep('prompt')">{{ $t('professionalEditor.batch.onlyPrompt') }}</el-button>
            <el-button :disabled="isBatching" @click="startStep('image')">{{ $t('professionalEditor.batch.onlyImage') }}</el-button>
            <el-button :disabled="isBatching" @click="startStep('video')">{{ $t('professionalEditor.batch.onlyVideo') }}</el-button>
          </el-button-group>
          <el-button type="success" plain :disabled="isBatching" @click="startUpscaleAll">
            <el-icon><MagicStick /></el-icon> Upscale All Videos
          </el-button>
          <el-button type="info" plain :loading="isDownloadingZip" @click="downloadAllVideos">
             <el-icon><Download /></el-icon> {{ $t('professionalEditor.batch.downloadAll') }}
          </el-button>
          <el-button type="danger" plain v-if="isBatching" @click="stopBatch">
            {{ $t('professionalEditor.batch.stop') }}
          </el-button>
        </div>
      </div>

      <!-- 分镜列表状态 -->
      <div class="shot-progress-list">
        <el-table :data="storyboards" style="width: 100%" height="400px" v-loading="isBatching && localStoryboards.length === 0">
          <el-table-column :label="$t('professionalEditor.batch.shot')" width="80" property="storyboard_number">
            <template #default="scope">
              {{ $t('professionalEditor.batch.shot') }} {{ scope.row.storyboard_number }}
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.description')" show-overflow-tooltip>
            <template #default="scope">
              {{ scope.row.action || scope.row.description || $t('professionalEditor.batch.noDescription') }}
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.prompt')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.id, 'prompt')" size="small">
                {{ getStatusText(scope.row.id, 'prompt') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.image')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.id, 'image')" size="small">
                {{ getStatusText(scope.row.id, 'image') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.video')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.id, 'video')" size="small">
                {{ getStatusText(scope.row.id, 'video') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.progress')" width="150">
            <template #default="scope">
              <el-progress 
                :percentage="getProgress(scope.row.id)" 
                :status="getProgressStatus(scope.row.id)"
                :stroke-width="10"
              />
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">{{ $t('common.close') || 'Close' }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { MagicStick, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import JSZip from 'jszip'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { dramaAPI } from '@/api/drama'
import { generateFramePrompt } from '@/api/frame'
import { imageAPI } from '@/api/image'
import { videoAPI } from '@/api/video'
import { taskAPI } from '@/api/task'
import { getVideoUrl } from '@/utils/image'
import type { Storyboard } from '@/types/drama'
import { useAISettings } from '@/composables/useAISettings'
// Local mutable copy of storyboards that can be refreshed from DB
const localStoryboards = ref<Storyboard[]>([])

const props = defineProps<{
  modelValue: boolean
  storyboards: Storyboard[]
  episodeId: number
  dramaId: number
  style?: string
  customStyle?: string
  videoModels: any[]
  defaultVideoModel?: string
}>()

const emit = defineEmits(['update:modelValue', 'completed', 'refresh'])
const { t } = useI18n()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const selectedVideoModel = ref(props.defaultVideoModel || '')
const isBatching = ref(false)
const isDownloadingZip = ref(false)
const shouldStop = ref(false)
const generationMode = ref<'key'|'action'>('action')
const { maxConcurrentThreads } = useAISettings()

// 任务状态追踪 map: shotId -> { step: 'prompt'|'image'|'video', status: 'pending'|'loading'|'done'|'error', progress: number }
const taskStates = reactive<Record<string, any>>({})

// Refresh storyboard data from DB to get latest status
const refreshStoryboardsFromDB = async () => {
  if (!props.episodeId) {
    // Fallback: use props data
    localStoryboards.value = [...props.storyboards]
    return
  }
  try {
    const res = await dramaAPI.getStoryboards(props.episodeId.toString())
    const freshStoryboards: Storyboard[] = res?.storyboards || []
    if (freshStoryboards.length > 0) {
      localStoryboards.value = freshStoryboards
    } else {
      localStoryboards.value = [...props.storyboards]
    }
  } catch (e) {
    console.warn('Failed to refresh storyboards from DB, using props data', e)
    localStoryboards.value = [...props.storyboards]
  }
}

const initTaskStates = async () => {
  // Fetch fresh data from DB before determining status
  await refreshStoryboardsFromDB()
  
  // Also check image/video records per storyboard from DB for accurate status
  for (const sb of localStoryboards.value) {
    // Check if prompt exists in DB ONLY. 
    // AND it must NOT be the backend-generated default fallback ("first frame" or "continuous movement progression")
    const isFallback = typeof sb.image_prompt === 'string' && (
      sb.image_prompt.trim().endsWith('first frame') || 
      sb.image_prompt.trim().endsWith('continuous movement progression')
    )
    const hasPrompt = sb.image_prompt && String(sb.image_prompt) !== 'null' && String(sb.image_prompt) !== 'undefined' && String(sb.image_prompt).trim().length > 0 && !isFallback;
    
    // Check image status from DB (not just from storyboard fields)
    let hasImage = !!sb.composed_image || !!sb.image_url
    if (!hasImage) {
      try {
        const imgRes = await imageAPI.listImages({ storyboard_id: Number(sb.id), frame_type: generationMode.value as any })
        const completedImg = imgRes.items?.find((i: any) => i.status === 'completed' && (i.image_url || i.local_path))
        if (completedImg) hasImage = true
      } catch { /* ignore */ }
    }
    
    // Check video status from DB
    let videoState = 'pending'
    try {
      const vidRes = await videoAPI.listVideos({ storyboard_id: String(sb.id) })
      // Find latest relevant video
      const completedVidHd = vidRes.items?.find((v: any) => v.status === 'completed' && v.is_upscaled)
      const upscalingVid = vidRes.items?.find((v: any) => v.status === 'upscaling')
      const completedVid = vidRes.items?.find((v: any) => v.status === 'completed' && !v.is_upscaled)
      const failedWithBase = vidRes.items?.find((v: any) => (v.status === 'failed' || v.status === 'error') && (v.video_url || v.local_path))

      if (completedVidHd) {
        videoState = 'hd'
      } else if (upscalingVid) {
        videoState = 'upscaling'
      } else if (completedVid) {
        videoState = 'done'
      } else if (failedWithBase) {
        videoState = 'done'
      } else if (sb.video_url) {
        videoState = 'done' // fallback
      }
    } catch { 
      if (sb.video_url) videoState = 'done'
    }
    
    let progress = 0
    if (videoState === 'hd') progress = 100
    else if (videoState === 'done') progress = 100
    else if (videoState === 'upscaling') progress = 90
    else if (hasImage) progress = 60
    else if (hasPrompt) progress = 30
    
    taskStates[sb.id] = {
      prompt: hasPrompt ? 'done' : 'pending',
      image: hasImage ? 'done' : 'pending',
      video: videoState,
      progress: progress
    }
  }
}

const getStatusTag = (id: string, type: string) => {
  const state = taskStates[id]?.[type]
  if (state === 'loading' || state === 'upscaling') return ''
  if (state === 'done') return 'success'
  if (state === 'hd') return 'warning'
  if (state === 'error') return 'danger'
  return 'info'
}

const getStatusText = (id: string, type: string) => {
  const state = taskStates[id]?.[type]

  if (state === 'upscaling') return 'Upscaling...'
  if (state === 'hd') return 'HD'
  if (state === 'loading') return t('professionalEditor.batch.status.loading')
  if (state === 'done') return t('professionalEditor.batch.status.done')
  if (state === 'error') return t('professionalEditor.batch.status.failed')
  return t('professionalEditor.batch.status.pending')
}

const getProgress = (id: string) => taskStates[id]?.progress || 0
const getProgressStatus = (id: string) => {
  const s = taskStates[id]
  if (s?.prompt === 'error' || s?.image === 'error' || s?.video === 'error') return 'exception'
  if (s?.video === 'done' || s?.video === 'hd') return 'success'
  if (s?.video === 'upscaling') return 'warning'
  return undefined
}

const stopBatch = () => {
  shouldStop.value = true
  isBatching.value = false
  ElMessage.warning(t('professionalEditor.batch.stopping'))
}

// 核心逻辑：提取提示词并保存
const processPrompt = async (sb: Storyboard) => {
  taskStates[sb.id].prompt = 'loading'
  taskStates[sb.id].progress = 10
  try {
    const { task_id } = await generateFramePrompt(Number(sb.id), { frame_type: generationMode.value as any })
    
    // 轮询
    let result = null
    while (true) {
      if (shouldStop.value) throw new Error('Stopped')
      const task = await taskAPI.getStatus(task_id)
      if (task.status === 'completed') {
        let res = task.result
        if (typeof res === 'string') res = JSON.parse(res)
        result = res.response
        break
      } else if (task.status === 'failed') {
        throw new Error(task.message || 'Prompt failed')
      }
      await new Promise(r => setTimeout(r, 2000))
    }

    let finalPrompt = ""
    if (result.single_frame) finalPrompt = result.single_frame.prompt
    else if (result.multi_frame?.frames) {
      finalPrompt = result.multi_frame.frames.map((f: any) => f.prompt).join('\n\n')
    }

    // Backend LLM already incorporates style when generating the image prompt.
    // We no longer manually prepend the large style prefix to `finalPrompt` to prevent 
    // confusing downstream API limits and causing token dilution.
    if (props.style && finalPrompt.startsWith(`${props.style}, `)) {
      // Clean up legacy simple prefix if it was prepended by mistake
      finalPrompt = finalPrompt.substring(props.style.length + 2)
    }

    // 保存到DB
    await dramaAPI.updateStoryboard(sb.id.toString(), { image_prompt: finalPrompt })
    
    // ALSO save to sessionStorage to keep it in sync with ProfessionalEditor UI (Action sequence)
    sessionStorage.setItem(`frame_prompt_${sb.id}_${generationMode.value}`, finalPrompt)
    
    taskStates[sb.id].prompt = 'done'
    taskStates[sb.id].progress = 30
    return finalPrompt
  } catch (e) {
    taskStates[sb.id].prompt = 'error'
    throw e
  }
}

// 核心逻辑：生成格点图
const processImage = async (sb: Storyboard, prompt: string) => {
  taskStates[sb.id].image = 'loading'
  try {
    // 收集参考图片（角色+场景），与 ProfessionalEditor 保持一致
    const referenceImages: string[] = []

    // 1. 添加场景图片
    if ((sb as any).background?.local_path) {
      referenceImages.push((sb as any).background.local_path)
    }

    // 2. 添加当前镜头登场的角色图片
    if (sb.characters && Array.isArray(sb.characters)) {
      sb.characters.forEach((char: any) => {
        if (char.local_path) {
          referenceImages.push(char.local_path)
        }
      })
    }

    // 3. 添加当前镜头中的道具图片
    if ((sb as any).props && Array.isArray((sb as any).props)) {
      (sb as any).props.forEach((prop: any) => {
        if (prop.local_path) {
          referenceImages.push(prop.local_path)
        }
      })
    }

    const result = await imageAPI.generateImage({
      drama_id: props.dramaId.toString(),
      prompt: prompt,
      storyboard_id: Number(sb.id),
      image_type: 'storyboard',
      frame_type: generationMode.value as any,
      reference_images: referenceImages.length > 0 ? referenceImages : undefined
    })

    // 轮询图片直到完成
    while (true) {
      if (shouldStop.value) throw new Error('Stopped')
      const res = await imageAPI.listImages({ storyboard_id: Number(sb.id), frame_type: generationMode.value as any })
      const img = res.items?.find((i: any) => i.id === result.id)
      if (img?.status === 'completed') {
        taskStates[sb.id].image = 'done'
        taskStates[sb.id].progress = 60
        return img
      } else if (img?.status === 'failed') {
        throw new Error(img.error_msg || 'Image generation failed')
      }
      await new Promise(r => setTimeout(r, 3000))
    }
  } catch (e) {
    taskStates[sb.id].image = 'error'
    throw e
  }
}


// 从模型名称提取provider (copied from ProfessionalEditor for consistency)
const extractProviderFromModel = (modelName: string): string => {
  if (modelName.startsWith("doubao-") || modelName.startsWith("seedance")) {
    return "doubao";
  }
  if (modelName.startsWith("runway")) {
    return "runway";
  }
  if (modelName.startsWith("pika")) {
    return "pika";
  }
  if (
    modelName.startsWith("MiniMax-") ||
    modelName.toLowerCase().startsWith("minimax") ||
    modelName.startsWith("hailuo")
  ) {
    return "minimax";
  }
  if (modelName.startsWith("sora")) {
    return "openai";
  }
  if (modelName.startsWith("kling")) {
    return "kling";
  }
  return "doubao";
};

// 核心逻辑：生成视频
const processVideo = async (sb: Storyboard, image: any) => {
  if (!selectedVideoModel.value) throw new Error('No video model')
  taskStates[sb.id].video = 'loading'
  try {
    const provider = extractProviderFromModel(selectedVideoModel.value)
    
    // 构建 R2V 请求
    const result = await videoAPI.generateVideo({
      drama_id: props.dramaId.toString(),
      storyboard_id: Number(sb.id),
      prompt: sb.video_prompt || sb.action || "Cinematic video",
      duration: 5,
      provider: provider,
      model: selectedVideoModel.value,
      reference_mode: 'multiple',
      reference_image_urls: [image.local_path || image.image_url]
    })
    
    // 轮询视频直到完成 (视频生成通常耗时较长)
    while (true) {
      if (shouldStop.value) throw new Error('Stopped')
      const videoTask = await videoAPI.getVideo(result.id)
      if (videoTask.status === 'completed') {
        taskStates[sb.id].video = 'done'
        taskStates[sb.id].progress = 100
        return videoTask
      } else if (videoTask.status === 'failed') {
        throw new Error(videoTask.error_msg || 'Video generation failed')
      }
      await new Promise(r => setTimeout(r, 5000))
    }
  } catch (e) {
    taskStates[sb.id].video = 'error'
    throw e
  }
}

const runConcurrently = async <T>(items: T[], limit: number, worker: (item: T) => Promise<void>) => {
  const executing: Promise<void>[] = []
  for (const item of items) {
    if (shouldStop.value) break
    const p = worker(item).finally(() => {
      executing.splice(executing.indexOf(p), 1)
    })
    executing.push(p)
    if (executing.length >= limit) {
      await Promise.race(executing)
    }
  }
  await Promise.all(executing)
}

const startFullBatch = async () => {
  if (!selectedVideoModel.value) {
    ElMessage.warning(t('professionalEditor.batch.selectVideoModelFirst'))
    return
  }
  isBatching.value = true
  shouldStop.value = false
  await initTaskStates()

  const promptCache = new Map<string, string>()

  // Phase 1: All Prompts
  await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
    if (taskStates[sb.id]?.prompt === 'done') return
    try {
      const p = await processPrompt(sb)
      promptCache.set(sb.id, p)
    } catch (e: any) {
      console.error(`Shot ${sb.storyboard_number} prompt failed:`, e)
      ElMessage.error(`${t('professionalEditor.batch.shot')} ${sb.storyboard_number} ${t('professionalEditor.batch.status.failed')}: ${e.message || 'Unknown'}`)
    }
  })

  // Phase 2: All Images
  if (!shouldStop.value) {
    await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
      if (taskStates[sb.id]?.prompt !== 'done' || taskStates[sb.id]?.image === 'done') return
      try {
        const p = promptCache.get(sb.id) || sb.image_prompt || ""
        await processImage(sb, p)
      } catch (e: any) {
        console.error(`Shot ${sb.storyboard_number} image failed:`, e)
        ElMessage.error(`${t('professionalEditor.batch.shot')} ${sb.storyboard_number} ${t('professionalEditor.batch.status.failed')}: ${e.message || 'Unknown'}`)
      }
    })
  }

  // Phase 3: All Videos
  if (!shouldStop.value) {
    await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
      if (taskStates[sb.id]?.image !== 'done' || taskStates[sb.id]?.video === 'done') return
      try {
         const res = await imageAPI.listImages({ storyboard_id: Number(sb.id), frame_type: generationMode.value as any })
         const img = res.items?.find((i: any) => i.status === 'completed' && (i.image_url || i.local_path))
         if (img) await processVideo(sb, img)
         else ElMessage.warning(t('professionalEditor.batch.lackActionImage', { number: sb.storyboard_number }))
      } catch (e: any) {
        console.error(`Shot ${sb.storyboard_number} video failed:`, e)
        ElMessage.error(`${t('professionalEditor.batch.shot')} ${sb.storyboard_number} ${t('professionalEditor.batch.status.failed')}: ${e.message || 'Unknown'}`)
      }
    })
  }

  isBatching.value = false
  if (!shouldStop.value) {
    ElMessage.success(t('professionalEditor.batch.completed'))
    emit('completed')
  }
}

const startStep = async (step: string) => {
  if (step === 'video' && !selectedVideoModel.value) {
    ElMessage.warning(t('professionalEditor.batch.selectVideoModelFirst'))
    return
  }
  isBatching.value = true
  shouldStop.value = false
  await initTaskStates()

  await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
    try {
      if (step === 'prompt') {
        if (taskStates[sb.id]?.prompt !== 'done') await processPrompt(sb)
      }
      else if (step === 'image') {
        if (taskStates[sb.id]?.image !== 'done') {
           const p = sb.image_prompt || await processPrompt(sb)
           await processImage(sb, p)
        }
      }
      else if (step === 'video') {
         if (taskStates[sb.id]?.video !== 'done' && taskStates[sb.id]?.video !== 'hd' && taskStates[sb.id]?.video !== 'upscaling') {
           const res = await imageAPI.listImages({ storyboard_id: Number(sb.id), frame_type: generationMode.value as any })
           const img = res.items?.find((i: any) => i.status === 'completed' && (i.image_url || i.local_path))
           if (img) await processVideo(sb, img)
           else ElMessage.warning(t('professionalEditor.batch.lackActionImage', { number: sb.storyboard_number }))
         }
      }
    } catch (e: any) {
      console.error(`Shot ${sb.storyboard_number} step ${step} failed:`, e)
      const stepName = step === 'prompt' ? t('professionalEditor.batch.prompt') : (step === 'image' ? t('professionalEditor.batch.image') : t('professionalEditor.batch.video'))
      ElMessage.error(`${t('professionalEditor.batch.shot')} ${sb.storyboard_number} ${stepName} ${t('professionalEditor.batch.status.failed')}: ${e.message || 'Unknown'}`)
    }
  })

  isBatching.value = false
  if (!shouldStop.value) {
    ElMessage.success(t('professionalEditor.batch.completed'))
    emit('completed')
  }
}

const startUpscaleAll = async () => {
  if (isBatching.value) return
  isBatching.value = true
  shouldStop.value = false
  
  try {
    ElMessage.info('Bắt đầu quy trình Upscale tất cả video đã hoàn thành...')
    await initTaskStates()
    
    await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
      // Process shots that have a base video (status 'done') or are already 'upscaling'
      if (!['done', 'upscaling'].includes(taskStates[sb.id]?.video)) return
      if (taskStates[sb.id]?.video === 'hd') return

      try {
        const vidRes = await videoAPI.listVideos({ storyboard_id: String(sb.id), page_size: 10 })
        const targetVideo = vidRes.items?.find((v: any) => (v.status === 'completed' && !v.is_upscaled) || (v.status === 'failed' && (v.video_url || v.local_path)))
        const currentlyUpscaling = vidRes.items?.find((v: any) => v.status === 'upscaling')

        if (currentlyUpscaling) {
          taskStates[sb.id].video = 'upscaling'
          let isUpscaled = false
          while (!isUpscaled && !shouldStop.value) {
            await new Promise(r => setTimeout(r, 5000))
            const check = await videoAPI.getVideo(currentlyUpscaling.id)
            if (check.status === 'completed' && check.is_upscaled) {
              isUpscaled = true
            } else if (check.status === 'failed') {
              throw new Error(check.error_msg || 'Upscaling failed')
            }
          }
          taskStates[sb.id].video = 'hd'
          taskStates[sb.id].progress = 100
          return
        }
        
        if (targetVideo) {
          taskStates[sb.id].video = 'upscaling'
          await videoAPI.upscaleVideo(targetVideo.id)
          let isUpscaled = false
          while (!isUpscaled && !shouldStop.value) {
            await new Promise(r => setTimeout(r, 5000))
            const check = await videoAPI.getVideo(targetVideo.id)
            if (check.status === 'completed' && check.is_upscaled) {
              isUpscaled = true
            } else if (check.status === 'failed') {
              throw new Error(check.error_msg || 'Upscaling failed')
            }
          }
          taskStates[sb.id].video = 'hd'
          taskStates[sb.id].progress = 100
        } else {
          // Check if it already has an upscaled version we missed
          const alreadyHd = vidRes.items?.find((v: any) => v.status === 'completed' && v.is_upscaled)
          if (alreadyHd) {
             taskStates[sb.id].video = 'hd'
             taskStates[sb.id].progress = 100
          } else {
             taskStates[sb.id].video = 'done'
          }
        }
      } catch (err: any) {
        console.error(`Shot ${sb.storyboard_number} auto upscale failed:`, err)
        taskStates[sb.id].video = 'error'
      }
    })
    
    if (!shouldStop.value) {
      ElMessage.success('Upscale All hoàn tất!')
    }
  } finally {
    isBatching.value = false
    emit('completed')
  }
}

const downloadAllVideos = async () => {
  if (isBatching.value) return
  isDownloadingZip.value = true
  
  try {
    // 1. Get ALL videos for this drama (not just completed, to catch those that failed upscale but have base URL)
    const allDramaVideos: any[] = []
    let currentPage = 1
    const pageSize = 100
    
    while (true) {
      const vidRes = await videoAPI.listVideos({ 
        drama_id: props.dramaId.toString(), 
        page: currentPage,
        page_size: pageSize
      })
      
      if (vidRes.items && vidRes.items.length > 0) {
        allDramaVideos.push(...vidRes.items)
      }
      
      if (!vidRes.pagination || vidRes.items.length < pageSize || allDramaVideos.length >= vidRes.pagination.total) {
        break
      }
      currentPage++
    }
    
    const usableVideos = allDramaVideos.filter(v => v.video_url || v.local_path)
    
    if (usableVideos.length === 0) {
      ElMessage.warning(t('professionalEditor.batch.noVideosToDownload'))
      return
    }

    // 2. Map storyboard_id to its best video (prefer HD, then latest usable)
    const bestVideos = new Map<number, any>()
    usableVideos.forEach((vid: any) => {
      const sbId = vid.storyboard_id
      if (!sbId) return
      
      const currentBest = bestVideos.get(sbId)
      if (!currentBest) {
        bestVideos.set(sbId, vid)
      } else {
        // Preference logic:
        // 1. HD (upscaled) vs non-HD
        if (vid.is_upscaled && !currentBest.is_upscaled) {
          bestVideos.set(sbId, vid)
        } else if (vid.is_upscaled === currentBest.is_upscaled) {
          // If both same HD status, pick latest
          if (new Date(vid.created_at) > new Date(currentBest.created_at)) {
            bestVideos.set(sbId, vid)
          }
        }
      }
    })

    // 3. Ensure we cover ALL storyboards if possible
    const itemsToDownload: { fileName: string; url: string }[] = []
    for (const sb of props.storyboards) {
      const vid = bestVideos.get(Number(sb.id))
      if (vid) {
        const sbNum = String(sb.storyboard_number).padStart(3, '0')
        const fileName = `Shot_${sbNum}.mp4`
        const url = getVideoUrl(vid)
        if (url) {
          itemsToDownload.push({ fileName, url })
        }
      }
    }

    if (itemsToDownload.length === 0) {
      ElMessage.warning(t('professionalEditor.batch.noVideosToDownload'))
      return
    }

    ElMessage.info(t('professionalEditor.batch.downloadingZip'))
    
    const zip = new JSZip()
    // No subfolders for the videos themselves, as requested
    
    // 4. Download each video and add to zip
    const downloadPromises = itemsToDownload.map(item => {
      const fullUrl = item.url.startsWith('http') ? item.url : `${window.location.origin}${item.url}`
      return axios.get(fullUrl, { responseType: 'blob' })
        .then(res => {
          zip.file(item.fileName, res.data)
        })
        .catch(err => {
          console.error(`Failed to download ${item.fileName}:`, err)
        })
    })
    
    await Promise.all(downloadPromises)
    
    // 5. Generate and download ZIP
    const content = await zip.generateAsync({ type: 'blob' })
    const zipUrl = URL.createObjectURL(content)
    const link = document.createElement('a')
    link.href = zipUrl
    link.download = `drama_${props.dramaId}_videos.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(zipUrl)
    
    ElMessage.success(t('professionalEditor.batch.completed'))
  } catch (err: any) {
    console.error('Download all failed:', err)
    ElMessage.error(t('message.operationFailed'))
  } finally {
    isDownloadingZip.value = false
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal) initTaskStates()
})

// Also expose storyboards for the table display
const storyboards = computed(() => {
  return localStoryboards.value.length > 0 ? localStoryboards.value : props.storyboards
})
</script>

<style scoped>
.batch-dialog :deep(.el-dialog__body) {
  padding-top: 10px;
}
.batch-container {
  padding: 0 10px;
}
.batch-config-bar {
  background: var(--el-fill-color-light);
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}
.config-row {
  display: flex;
  gap: 30px;
  align-items: center;
  flex-wrap: wrap;
}
.config-item .label {
  font-weight: bold;
  margin-right: 12px;
  color: var(--el-text-color-regular);
}
.action-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.shot-progress-list {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  overflow: hidden;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding-right: 10px;
}

/* 动画和过渡 */
.el-progress--line {
  margin-bottom: 0;
}
</style>
